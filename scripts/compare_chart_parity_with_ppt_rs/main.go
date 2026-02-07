package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/djinn09/goppt/pkg/pptx"
)

var chartOrder = []string{
	"bar",
	"barHorizontal",
	"barStacked",
	"barStacked100",
	"line",
	"lineMarkers",
	"lineStacked",
	"area",
	"areaStacked",
	"areaStacked100",
	"pie",
	"doughnut",
	"scatter",
	"scatterLines",
	"scatterSmooth",
	"bubble",
	"radar",
	"radarFilled",
	"stockHLC",
	"stockOHLC",
	"combo",
}

var signatureTokens = []string{
	`<c:barChart`,
	`<c:barDir val="bar"/>`,
	`<c:barDir val="col"/>`,
	`<c:grouping val="clustered"/>`,
	`<c:grouping val="stacked"/>`,
	`<c:grouping val="percentStacked"/>`,
	`<c:lineChart`,
	`<c:grouping val="standard"/>`,
	`<c:marker><c:symbol val="circle"/>`,
	`<c:areaChart`,
	`<c:pieChart`,
	`<c:doughnutChart`,
	`<c:holeSize`,
	`<c:scatterChart`,
	`<c:scatterStyle val="marker"/>`,
	`<c:scatterStyle val="lineMarker"/>`,
	`<c:scatterStyle val="smoothMarker"/>`,
	`<c:bubbleChart`,
	`<c:varyColors val="0"/>`,
	`<c:bubbleScale`,
	`<c:bubbleSize`,
	`<c:xVal`,
	`<c:yVal`,
	`<c:radarChart`,
	`<c:radarStyle val="marker"`,
	`<c:radarStyle val="filled"`,
	`<c:stockChart`,
	`<c:tx><c:v>Open</c:v></c:tx>`,
	`<c:tx><c:v>High</c:v></c:tx>`,
	`<c:tx><c:v>Low</c:v></c:tx>`,
	`<c:tx><c:v>Close</c:v></c:tx>`,
}

var requiredTokenOverrides = map[string]map[string]string{
	"bar": {
		`<c:barDir val="bar"/>`: `<c:barDir val="col"/>`,
	},
	"barStacked": {
		`<c:barDir val="bar"/>`:         `<c:barDir val="col"/>`,
		`<c:grouping val="clustered"/>`: `<c:grouping val="stacked"/>`,
	},
	"barStacked100": {
		`<c:barDir val="bar"/>`:         `<c:barDir val="col"/>`,
		`<c:grouping val="clustered"/>`: `<c:grouping val="percentStacked"/>`,
	},
}

type compareResult struct {
	Chart         string
	RefSeries     int
	OurSeries     int
	Required      []string
	Missing       []string
	Pass          bool
	ReferenceOnly bool
}

func main() {
	referenceXML, err := loadReferenceXML()
	if err != nil {
		fail("load ppt-rs reference XML", err)
	}

	ourXML, err := loadGoPPTXML()
	if err != nil {
		fail("generate goppt chart XML", err)
	}

	results := compare(referenceXML, ourXML)
	report := renderReport(results)

	if err := os.MkdirAll("reports", 0o755); err != nil {
		fail("create reports directory", err)
	}
	reportPath := filepath.Join("reports", "chart_parity_report.md")
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		fail("write parity report", err)
	}

	fmt.Printf("Wrote %s\n", reportPath)
	printSummary(results)

	for _, result := range results {
		if !result.Pass {
			os.Exit(1)
		}
	}
}

func loadReferenceXML() (map[string]string, error) {
	cmd := exec.Command(
		"cargo",
		"run",
		"--quiet",
		"--manifest-path",
		"scripts/compare_chart_parity_with_ppt_rs/ppt_rs_chart_signatures/Cargo.toml",
	)
	cmd.Env = append(os.Environ(), "CARGO_TARGET_DIR=.tmp/cargo-target/ppt-rs-chart-signatures")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("cargo run failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, err
	}
	var out map[string]string
	if err := json.Unmarshal(output, &out); err != nil {
		return nil, fmt.Errorf("decode reference JSON: %w", err)
	}
	return out, nil
}

func loadGoPPTXML() (map[string]string, error) {
	out := make(map[string]string, len(chartOrder))

	entries := map[string]pptx.SlideContent{
		"bar":            barSlide(),
		"barHorizontal":  barHorizontalSlide(),
		"barStacked":     barStackedSlide(),
		"barStacked100":  barStacked100Slide(),
		"line":           lineSlide(),
		"lineMarkers":    lineMarkersSlide(),
		"lineStacked":    lineStackedSlide(),
		"area":           areaSlide(),
		"areaStacked":    areaStackedSlide(),
		"areaStacked100": areaStacked100Slide(),
		"pie":            pieSlide(),
		"doughnut":       doughnutSlide(),
		"scatter":        scatterMarkerSlide(),
		"scatterLines":   scatterLinesSlide(),
		"scatterSmooth":  scatterSmoothSlide(),
		"bubble":         bubbleSlide(),
		"radar":          radarSlide(),
		"radarFilled":    radarFilledSlide(),
		"stockHLC":       stockHLCSlide(),
		"stockOHLC":      stockOHLCSlide(),
		"combo":          comboSlide(),
	}

	for key, slide := range entries {
		xml, err := chartXMLForSlide(slide)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		out[key] = xml
	}
	return out, nil
}

func chartXMLForSlide(slide pptx.SlideContent) (string, error) {
	data, err := pptx.CreateWithSlides("Parity", []pptx.SlideContent{slide})
	if err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	return readZipFile(zr, "ppt/charts/chart1.xml")
}

func readZipFile(zr *zip.Reader, name string) (string, error) {
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		r, err := f.Open()
		if err != nil {
			return "", err
		}
		defer func() { _ = r.Close() }()
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return "", fmt.Errorf("zip entry not found: %s", name)
}

func compare(reference map[string]string, ours map[string]string) []compareResult {
	results := make([]compareResult, 0, len(chartOrder))
	for _, key := range chartOrder {
		refXML, ok := reference[key]
		if !ok {
			results = append(results, compareResult{Chart: key, Pass: false, ReferenceOnly: true})
			continue
		}
		ourXML := ours[key]

		required := normalizeRequiredTokens(key, requiredTokensFromReference(refXML))
		missing := missingTokens(required, ourXML)

		refSeries := strings.Count(refXML, "<c:ser>")
		ourSeries := strings.Count(ourXML, "<c:ser>")

		pass := len(missing) == 0 && refSeries == ourSeries
		results = append(results, compareResult{
			Chart:     key,
			RefSeries: refSeries,
			OurSeries: ourSeries,
			Required:  required,
			Missing:   missing,
			Pass:      pass,
		})
	}
	return results
}

func requiredTokensFromReference(xml string) []string {
	required := make([]string, 0, len(signatureTokens))
	for _, token := range signatureTokens {
		if strings.Contains(xml, token) {
			required = append(required, token)
		}
	}
	sort.Strings(required)
	return required
}

func missingTokens(required []string, xml string) []string {
	missing := make([]string, 0)
	for _, token := range required {
		if !strings.Contains(xml, token) {
			missing = append(missing, token)
		}
	}
	return missing
}

func normalizeRequiredTokens(chart string, required []string) []string {
	overrides, ok := requiredTokenOverrides[chart]
	if !ok {
		return required
	}

	normalized := make([]string, len(required))
	copy(normalized, required)
	for i := range normalized {
		if replacement, exists := overrides[normalized[i]]; exists {
			normalized[i] = replacement
		}
	}

	seen := make(map[string]struct{}, len(normalized))
	unique := make([]string, 0, len(normalized))
	for _, token := range normalized {
		if _, exists := seen[token]; exists {
			continue
		}
		seen[token] = struct{}{}
		unique = append(unique, token)
	}
	sort.Strings(unique)
	return unique
}

func renderReport(results []compareResult) string {
	var b strings.Builder
	b.WriteString("# Chart Parity Report (goppt vs ppt-rs)\n\n")
	b.WriteString("| Chart | Status | Series (ref/our) | Missing tokens |\n")
	b.WriteString("|---|---|---:|---|\n")
	for _, r := range results {
		status := "PASS"
		if !r.Pass {
			status = "FAIL"
		}
		missing := "-"
		if len(r.Missing) > 0 {
			missing = strings.Join(r.Missing, "<br>")
		}
		if r.ReferenceOnly {
			missing = "reference XML not produced"
		}
		b.WriteString(fmt.Sprintf("| `%s` | %s | %d/%d | %s |\n", r.Chart, status, r.RefSeries, r.OurSeries, missing))
	}
	b.WriteString("\nGenerated by `go run ./scripts/compare_chart_parity_with_ppt_rs`.\n")
	return b.String()
}

func printSummary(results []compareResult) {
	passed := 0
	for _, r := range results {
		if r.Pass {
			passed++
		}
	}
	fmt.Printf("Parity result: %d/%d chart signatures matched ppt-rs reference requirements.\n", passed, len(results))
	for _, r := range results {
		if r.Pass {
			continue
		}
		fmt.Printf("  - %s failed\n", r.Chart)
	}
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", step, err)
	os.Exit(1)
}
