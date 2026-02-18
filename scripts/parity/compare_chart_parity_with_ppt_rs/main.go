package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func main() {
	referenceXML, err := loadReferenceXML()
	if err != nil {
		fail("load ppt-rs reference XML", err)
	}

	ourXML, err := loadGoPPTXXML()
	if err != nil {
		fail("generate gopptx chart XML", err)
	}

	results := compare(referenceXML, ourXML)
	report := renderReport(results)

	if mkdirErr := os.MkdirAll("reports", 0o755); mkdirErr != nil {
		fail("create reports directory", mkdirErr)
	}
	reportPath := filepath.Join("reports", "chart_parity_report.md")
	if writeErr := os.WriteFile(reportPath, []byte(report), 0o644); writeErr != nil {
		fail("write parity report", writeErr)
	}

	log.Printf("Wrote %s\n", reportPath)
	printSummary(results)

	for _, result := range results {
		if !result.Pass {
			os.Exit(1)
		}
	}
}

func loadReferenceXML() (map[string]string, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"cargo",
		"run",
		"--quiet",
		"--manifest-path",
		"scripts/parity/compare_chart_parity_with_ppt_rs/ppt_rs_chart_signatures/Cargo.toml",
	)
	cmd.Env = append(os.Environ(), "CARGO_TARGET_DIR=.tmp/cargo-target/ppt-rs-chart-signatures")
	output, err := cmd.Output()
	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("cargo run failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return nil, err
	}
	var out map[string]string
	if decodeErr := json.Unmarshal(output, &out); decodeErr != nil {
		return nil, fmt.Errorf("decode reference JSON: %w", decodeErr)
	}
	return out, nil
}

func loadGoPPTXXML() (map[string]string, error) {
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
		if _, readErr := buf.ReadFrom(r); readErr != nil {
			return "", readErr
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
	b.WriteString("# Chart Parity Report (gopptx vs ppt-rs)\n\n")
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
	b.WriteString("\nGenerated by `go run ./scripts/parity/compare_chart_parity_with_ppt_rs`.\n")
	return b.String()
}

func printSummary(results []compareResult) {
	passed := 0
	for _, r := range results {
		if r.Pass {
			passed++
		}
	}
	log.Printf("Parity result: %d/%d chart signatures matched ppt-rs reference requirements.\n", passed, len(results))
	for _, r := range results {
		if r.Pass {
			continue
		}
		log.Printf("  - %s failed\n", r.Chart)
	}
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", step, err)
	os.Exit(1)
}
