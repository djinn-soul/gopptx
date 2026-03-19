package main

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/djinn-soul/gopptx/pkg/stdlog"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

const (
	outputDir  = "examples/output"
	outputFile = "12_feature_showcase.pptx"
)

func main() {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		fail("create output directory", err)
	}

	slides, slideErr := buildShowcaseSlides()
	if slideErr != nil {
		fail("build slides", slideErr)
	}

	data, err := pptx.CreateWithSlides("gopptx Feature Showcase", slides)
	if err != nil {
		fail("create presentation", err)
	}

	path := filepath.Join(outputDir, outputFile)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		fail("write output", err)
	}

	log.Printf("Wrote %s\n", path)
}

func buildShowcaseSlides() ([]pptx.SlideContent, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("gopptx Feature Showcase").
			AddBullet("Chart parity (13 variants)").
			AddBullet("Slide layouts (title/content, title-only, blank)").
			AddBullet("Table styling + deep border semantics").
			AddBullet("Markdown inline rich text and text-run formatting").
			AddBullet("Shapes/connectors with gradient fill rendering"),

		pptx.NewSlide("Title Only Layout").WithTitleOnlyLayout(),
		pptx.NewSlide("").WithBlankLayout(),

		pptx.NewSlide("Run-Level Text Formatting").
			AddBulletRuns([]pptx.TextRun{
				pptx.NewTextRun("Bold ").WithBold(true),
				pptx.NewTextRun("Italic ").WithItalic(true),
				pptx.NewTextRun("Underline ").WithUnderline(true),
				pptx.NewTextRun("Color ").WithColor("1F4E78"),
				pptx.NewTextRun("Custom Font ").WithFont("Calibri"),
				pptx.NewTextRun("Sized").WithSizePt(20),
			}).
			AddBulletRuns([]pptx.TextRun{
				pptx.NewTextRun("Inline code style").WithCode(true),
			}),

		pptx.NewSlide("Table Styling + Borders").WithTable(
			pptx.NewTable([]pptx.Length{pptx.Inches(2.62), pptx.Inches(2.62), pptx.Inches(2.62)}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Feature").WithBold(true).WithBackgroundColor("D9E1F2"),
					pptx.NewTableCell("State").WithBold(true).WithBackgroundColor("D9E1F2"),
					pptx.NewTableCell("Notes").WithBold(true).WithBackgroundColor("D9E1F2"),
				}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Per-side border").
						WithAlignLeft().
						WithLeftBorderStyle(1.5, "AA0000", pptx.TableBorderDashDash).
						WithRightBorderStyle(1.5, "00AA00", pptx.TableBorderDashDot).
						WithTopBorderStyle(1.5, "0000AA", pptx.TableBorderDashLongDash).
						WithBottomBorderStyle(1.5, "777777", pptx.TableBorderDashSolid),
					pptx.NewTableCell("Done").WithAlignCenter().WithVAlignMiddle(),
					pptx.NewTableCell("XML emits only configured sides").WithAlignLeft(),
				}),
		),

		layoutHelpersGridSlide(),
		flowDiagramSlide(),

		pptx.NewSlide("Chart Sample (Combo)").WithComboChart(
			pptx.NewComboChart(
				[]string{"Q1", "Q2", "Q3"},
				[]pptx.Series{
					{Name: "Revenue", Values: []float64{12, 15, 14}},
				},
				[]pptx.Series{
					{Name: "Trend", Values: []float64{11, 14, 16}},
				},
			).WithTitle("Combo Parity"),
		),
	}

	markdown := `# Markdown Parsed Slide
- Plain bullet
- **Bold** + *italic* + ` + "`code`" + `
---
# Markdown Parsed Slide 2
1. Numbered item
2. Another item`

	mdSlides, err := pptx.SlidesFromMarkdown(markdown)
	if err != nil {
		return nil, err
	}

	return append(slides, mdSlides...), nil
}

func layoutHelpersGridSlide() pptx.SlideContent {
	slide := pptx.NewSlide("Layout Helpers (2x3 Grid)").WithTitleOnlyLayout()
	boxes, _ := pptx.GridInBox(
		2,
		3,
		pptx.Inches(0.4),
		pptx.Box{
			X:  pptx.Inches(0.5),
			Y:  pptx.Inches(1.8),
			CX: pptx.Inches(8.5),
			CY: pptx.Inches(4.7),
		},
	)
	for i, box := range boxes {
		slide = slide.AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, box.X, box.Y, box.CX, box.CY).
				WithText(fmt.Sprintf("Item %d", i+1)).
				WithFill(pptx.NewShapeFill("ED7D31")).
				WithLine(pptx.NewShapeLine("C65911", pptx.Points(1))),
		)
	}
	return slide
}

func flowDiagramSlide() pptx.SlideContent {
	return pptx.NewSlide("").WithBlankLayout().
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.9), pptx.Inches(1.8), pptx.Inches(2.6), pptx.Inches(1.1)).
				WithText("Input").
				WithFill(pptx.NewShapeFill("D9E1F2")).
				WithLine(pptx.NewShapeLine("5B9BD5", pptx.Points(1.5))).
				WithAutoFit(pptx.TextAutoFitNormal),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeDiamond, pptx.Inches(3.9), pptx.Inches(1.45), pptx.Inches(3.6), pptx.Inches(2.0)).
				WithText("Check").
				WithGradientFill(
					pptx.NewShapeGradientFill(
						pptx.ShapeGradientTypeLinear,
						[]pptx.ShapeGradientStop{
							pptx.NewShapeGradientStop(0, "4472C4"),
							pptx.NewShapeGradientStop(100, "8FB9E0").WithTransparency(0.2),
						},
					).WithLinearAngle(35),
				).
				WithLine(pptx.NewShapeLine("1F4E78", pptx.Points(1.5))).
				WithAutoFit(pptx.TextAutoFitNone).
				WithTextWrap(pptx.TextWrapNone),
		).
		AddConnector(
			pptx.NewElbowConnector(pptx.Inches(3.5), pptx.Inches(2.35), pptx.Inches(4.1), pptx.Inches(2.35)).
				WithLine(pptx.NewShapeLine("1F4E78", pptx.Points(1.1)).WithDash(pptx.LineDashDashDot)).
				WithArrows(pptx.ArrowTypeNone, pptx.ArrowTypeTriangle).
				WithArrowSize(pptx.ArrowSizeLarge).
				ConnectStartAuto(1).
				ConnectEndAuto(2),
		)
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", step, err)
	os.Exit(1)
}
