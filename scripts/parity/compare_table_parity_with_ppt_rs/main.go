package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

type tableCase struct {
	Name       string
	FixturePPT string
	FixtureXML string
	Generated  pptx.SlideContent
}

type tableResult struct {
	Name      string
	RefRows   int
	OurRows   int
	RefCols   int
	OurCols   int
	Required  []string
	Missing   []string
	Pass      bool
	LoadError string
}

var tableCases = []tableCase{
	{
		Name:       "styled_table",
		FixturePPT: "comprehensive_demo.pptx",
		FixtureXML: "ppt/slides/slide7.xml",
		Generated: pptx.NewSlide("Table Demo").
			WithTable(
				pptx.NewTable([]pptx.Length{pptx.Emu(2200000), pptx.Emu(2200000), pptx.Emu(2200000)}).
					AddStyledRow([]pptx.TableCell{
						pptx.NewTableCell("Region").WithBold(true).WithBackgroundColor("2F5597").WithAlignCenter(),
						pptx.NewTableCell("Quarter").WithBold(true).WithBackgroundColor("2F5597").WithAlignCenter(),
						pptx.NewTableCell("Revenue").WithBold(true).WithBackgroundColor("2F5597").WithAlignCenter(),
					}).
					AddStyledRow([]pptx.TableCell{
						pptx.NewTableCell("North").WithAlignLeft(),
						pptx.NewTableCell("Q1").WithAlignCenter(),
						pptx.NewTableCell("$120k").WithAlignRight(),
					}).
					AddStyledRow([]pptx.TableCell{
						pptx.NewTableCell("South").WithAlignLeft(),
						pptx.NewTableCell("Q1").WithAlignCenter(),
						pptx.NewTableCell("$110k").WithAlignRight(),
					}).
					AddStyledRow([]pptx.TableCell{
						pptx.NewTableCell("West").WithAlignLeft(),
						pptx.NewTableCell("Q1").WithAlignCenter(),
						pptx.NewTableCell("$105k").WithAlignRight(),
					}),
			),
	},
	{
		Name:       "merged_cells",
		FixturePPT: "merged_cells.pptx",
		FixtureXML: "ppt/slides/slide1.xml",
		Generated: pptx.NewSlide("Merged Cells").WithTable(
			pptx.NewTable([]pptx.Length{pptx.Emu(2000000), pptx.Emu(2000000), pptx.Emu(2000000)}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Merged Header").WithRowSpan(2).WithColSpan(2),
					pptx.NewTableCell(""),
					pptx.NewTableCell("C1"),
				}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell(""),
					pptx.NewTableCell(""),
					pptx.NewTableCell("C2"),
				}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("R3C1"),
					pptx.NewTableCell("R3C2"),
					pptx.NewTableCell("R3C3"),
				}),
		),
	},
}

var tableSignatureTokens = []string{
	`<a:tbl>`,
	`<a:tblGrid>`,
	`<a:gridCol w="`,
	`<a:tr h="`,
	`<a:tc`,
	`<a:pPr algn="`,
	`<a:solidFill><a:srgbClr`,
	`rowSpan="`,
	`gridSpan="`,
	`hMerge="1"`,
	`vMerge="1"`,
}

func main() {
	results := make([]tableResult, 0, len(tableCases))
	for _, tc := range tableCases {
		result := compareTableCase(tc)
		results = append(results, result)
	}

	report := renderReport(results)
	if err := os.MkdirAll("reports", 0o755); err != nil {
		fail("create reports directory", err)
	}
	reportPath := filepath.Join("reports", "table_parity_report.md")
	if err := os.WriteFile(reportPath, []byte(report), 0o644); err != nil {
		fail("write parity report", err)
	}
	log.Printf("Wrote %s\n", reportPath)

	passed := 0
	for _, result := range results {
		if result.Pass {
			passed++
		}
	}
	log.Printf("Parity result: %d/%d table signatures matched ppt-rs fixture requirements.\n", passed, len(results))
	for _, result := range results {
		if result.Pass {
			continue
		}
		log.Printf("  - %s failed\n", result.Name)
	}
	if passed != len(results) {
		os.Exit(1)
	}
}

func compareTableCase(tc tableCase) tableResult {
	refXML, err := readFixtureSlideXML(tc.FixturePPT, tc.FixtureXML)
	if err != nil {
		return tableResult{Name: tc.Name, Pass: false, LoadError: err.Error()}
	}
	ourXML, err := generatedSlideXML(tc.Generated)
	if err != nil {
		return tableResult{Name: tc.Name, Pass: false, LoadError: err.Error()}
	}

	required := requiredTokens(refXML)
	missing := missingTokens(required, ourXML)
	refRows := strings.Count(refXML, "<a:tr ")
	ourRows := strings.Count(ourXML, "<a:tr ")
	refCols := strings.Count(refXML, "<a:gridCol ")
	ourCols := strings.Count(ourXML, "<a:gridCol ")
	pass := len(missing) == 0 && refRows == ourRows && refCols == ourCols
	return tableResult{
		Name:     tc.Name,
		RefRows:  refRows,
		OurRows:  ourRows,
		RefCols:  refCols,
		OurCols:  ourCols,
		Required: required,
		Missing:  missing,
		Pass:     pass,
	}
}

func readFixtureSlideXML(fixturePPT string, slidePath string) (string, error) {
	path := filepath.Join("testdata", "ppt_rs", fixturePPT)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read fixture %s: %w", path, err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("open fixture zip %s: %w", path, err)
	}
	return readZipEntry(zr, slidePath)
}

func generatedSlideXML(slide pptx.SlideContent) (string, error) {
	data, err := pptx.CreateWithSlides("Table Parity", []pptx.SlideContent{slide})
	if err != nil {
		return "", err
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	return readZipEntry(zr, "ppt/slides/slide1.xml")
}

func readZipEntry(zr *zip.Reader, name string) (string, error) {
	for _, f := range zr.File {
		if f.Name != name {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer func() { _ = rc.Close() }()
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(rc); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return "", fmt.Errorf("zip entry not found: %s", name)
}

func requiredTokens(referenceXML string) []string {
	required := make([]string, 0, len(tableSignatureTokens))
	for _, token := range tableSignatureTokens {
		if strings.Contains(referenceXML, token) {
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

func renderReport(results []tableResult) string {
	var b strings.Builder
	b.WriteString("# Table Parity Report (gopptx vs ppt-rs)\n\n")
	b.WriteString("| Case | Status | Rows (ref/our) | Cols (ref/our) | Missing tokens |\n")
	b.WriteString("|---|---|---:|---:|---|\n")
	for _, result := range results {
		status := "PASS"
		if !result.Pass {
			status = "FAIL"
		}
		missing := "-"
		if result.LoadError != "" {
			missing = "load error: " + result.LoadError
		} else if len(result.Missing) > 0 {
			missing = strings.Join(result.Missing, "<br>")
		}
		b.WriteString(fmt.Sprintf(
			"| `%s` | %s | %d/%d | %d/%d | %s |\n",
			result.Name,
			status,
			result.RefRows,
			result.OurRows,
			result.RefCols,
			result.OurCols,
			missing,
		))
	}
	b.WriteString("\nGenerated by `go run ./scripts/parity/compare_table_parity_with_ppt_rs`.\n")
	return b.String()
}

func fail(step string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "error: %s: %v\n", step, err)
	os.Exit(1)
}
