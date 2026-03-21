package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/editor"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	outputDir        = "examples/output"
	outputFile       = "43_presentation_props_editor.pptx"
	reskinOutputFile = "43_brand_reskin_theme_swap.pptx"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	tmpDir, tempErr := os.MkdirTemp("", "gopptx-props-example-*")
	if tempErr != nil {
		return fmt.Errorf("create temp directory: %w", tempErr)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp dir %s: %v\n", tmpDir, err)
		}
	}()

	inputPath := filepath.Join(tmpDir, "base_props_input.pptx")
	outputPath := filepath.Join(outputDir, outputFile)

	slides := buildPropsDemoSlides()
	if err := pptx.WriteFile(inputPath, "Presentation Props Base", slides); err != nil {
		return fmt.Errorf("create base presentation: %w", err)
	}

	editor, err := pptx.OpenPresentationEditor(inputPath)
	if err != nil {
		return fmt.Errorf("open base presentation: %w", err)
	}
	defer func() {
		if closeErr := editor.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close editor: %v\n", closeErr)
		}
	}()

	if err := editor.ApplyTheme(styling.ThemeCorporate); err != nil {
		return fmt.Errorf("apply theme: %w", err)
	}
	if err := editor.SetSlideSize(pptx.SlideSize16x9()); err != nil {
		return fmt.Errorf("set slide size: %w", err)
	}

	props := common.CoreProperties{
		Title:          "Presentation Properties Example",
		Subject:        "Editor metadata update",
		Creator:        "gopptx example",
		Description:    "Demonstrates theme, slide size, and core properties edits.",
		Keywords:       "gopptx, editor, metadata",
		LastModifiedBy: "gopptx example",
	}
	editor.SetCoreProperties(props)

	if err := editor.Save(outputPath); err != nil {
		return fmt.Errorf("save edited presentation: %w", err)
	}

	verificationEditor, err := pptx.OpenPresentationEditor(outputPath)
	if err != nil {
		return fmt.Errorf("reopen saved presentation: %w", err)
	}
	defer func() {
		if closeErr := verificationEditor.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close verification editor: %v\n", closeErr)
		}
	}()

	got := verificationEditor.GetCoreProperties()
	if got.Title != props.Title {
		return fmt.Errorf("core title mismatch: got %q want %q", got.Title, props.Title)
	}

	log.Printf("Generated presentation properties example: %s\n", outputPath)
	if err := runThemeReskinSmoke(); err != nil {
		return fmt.Errorf("run theme reskin smoke: %w", err)
	}
	return nil
}

func runThemeReskinSmoke() error {
	tmpDir, err := os.MkdirTemp("", "gopptx-theme-reskin-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "theme_reskin_input.pptx")
	outputPath := filepath.Join(outputDir, reskinOutputFile)

	slides := buildReskinSmokeSlides()
	if err := pptx.WriteFile(inputPath, "Theme Reskin Base", slides); err != nil {
		return fmt.Errorf("create base deck: %w", err)
	}

	ed, err := pptx.OpenPresentationEditor(inputPath)
	if err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	defer ed.Close()

	inv, err := ed.ThemeInventory()
	if err != nil {
		return fmt.Errorf("theme inventory: %w", err)
	}
	if len(inv.ThemeParts) == 0 {
		return errors.New("no theme parts discovered")
	}

	if err := ed.SetGlobalThemePreset("facet"); err != nil {
		return fmt.Errorf("apply preset: %w", err)
	}
	if err := ed.SetThemeFontScheme("Aptos Display", "Aptos"); err != nil {
		return fmt.Errorf("set font scheme: %w", err)
	}
	if err := ed.SetThemeColorScheme(editor.ThemeColorScheme{
		Accent1: "003366",
		Accent2: "00A3A3",
		Accent3: "E67E22",
		Hlink:   "0B66D0",
	}); err != nil {
		return fmt.Errorf("set color scheme: %w", err)
	}

	if err := ed.Save(outputPath); err != nil {
		return fmt.Errorf("save output: %w", err)
	}

	log.Printf("Generated theme reskin smoke example: %s\n", outputPath)
	return nil
}

func buildPropsDemoSlides() []pptx.SlideContent {
	hero := pptx.NewSlide("").WithBlankLayout().
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0), pptx.Inches(0), pptx.Inches(13.33), pptx.Inches(7.5)).
				WithFill(pptx.NewShapeFill("F7FAFC")),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.6), pptx.Inches(0.6), pptx.Inches(12.1), pptx.Inches(2.1)).
				WithGradientFill(
					pptx.NewShapeGradientFill(
						pptx.ShapeGradientTypeLinear,
						[]pptx.ShapeGradientStop{
							pptx.NewShapeGradientStop(0, "0B3C6D"),
							pptx.NewShapeGradientStop(100, "1F7A8C"),
						},
					).WithLinearAngle(25),
				).
				WithLine(pptx.NewShapeLine("0B3C6D", pptx.Points(1.2))).
				WithText("Presentation Properties + Theme Management"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.9), pptx.Inches(3.2), pptx.Inches(3.9), pptx.Inches(1.15)).
				WithFill(pptx.NewShapeFill("DCE6F2")).
				WithLine(pptx.NewShapeLine("9DB5CF", pptx.Points(1))).
				WithText("Metadata edits"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(4.95), pptx.Inches(3.2), pptx.Inches(3.9), pptx.Inches(1.15)).
				WithFill(pptx.NewShapeFill("D8F0EE")).
				WithLine(pptx.NewShapeLine("8BCAC4", pptx.Points(1))).
				WithText("Slide size + theme"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(9), pptx.Inches(3.2), pptx.Inches(3.4), pptx.Inches(1.15)).
				WithFill(pptx.NewShapeFill("FCE8D8")).
				WithLine(pptx.NewShapeLine("E3B082", pptx.Points(1))).
				WithText("Master-aware reskin"),
		)

	kpis := pptx.NewSlide("KPI Snapshot").WithBlankLayout().
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0), pptx.Inches(0), pptx.Inches(13.33), pptx.Inches(1.1)).
				WithFill(pptx.NewShapeFill("0D2E4F")).
				WithText("Quarterly Highlights"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.8), pptx.Inches(1.8), pptx.Inches(3.8), pptx.Inches(2.1)).
				WithFill(pptx.NewShapeFill("E8F1FB")).
				WithLine(pptx.NewShapeLine("9CB9D9", pptx.Points(1.2))).
				WithText("Pipeline\n$4.2M"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(4.8), pptx.Inches(1.8), pptx.Inches(3.8), pptx.Inches(2.1)).
				WithFill(pptx.NewShapeFill("E5F7F4")).
				WithLine(pptx.NewShapeLine("8BCBC4", pptx.Points(1.2))).
				WithText("Win Rate\n38%"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(8.8), pptx.Inches(1.8), pptx.Inches(3.8), pptx.Inches(2.1)).
				WithFill(pptx.NewShapeFill("FFF1E4")).
				WithLine(pptx.NewShapeLine("E6B17A", pptx.Points(1.2))).
				WithText("NPS\n61"),
		).
		AddShape(
			pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.8), pptx.Inches(4.35), pptx.Inches(11.8), pptx.Inches(1.2)).
				WithFill(pptx.NewShapeFill("F3F8FF")).
				WithLine(pptx.NewShapeLine("BFD2EA", pptx.Points(1))).
				WithText("• Cards make theme color changes easy to inspect.\n• This sample is intentionally shape-heavy for visual validation."),
		)

	table := pptx.NewSlide("Delivery Plan").WithTable(
		pptx.NewTable([]pptx.Length{pptx.Inches(2.8), pptx.Inches(2.3), pptx.Inches(2.7), pptx.Inches(2.7)}).
			AddStyledRow([]pptx.TableCell{
				pptx.NewTableCell("Workstream").WithBold(true).WithBackgroundColor("DCE6F2"),
				pptx.NewTableCell("Owner").WithBold(true).WithBackgroundColor("DCE6F2"),
				pptx.NewTableCell("Status").WithBold(true).WithBackgroundColor("DCE6F2"),
				pptx.NewTableCell("Next Milestone").WithBold(true).WithBackgroundColor("DCE6F2"),
			}).
			AddStyledRow([]pptx.TableCell{
				pptx.NewTableCell("Theme migration"),
				pptx.NewTableCell("Design Ops"),
				pptx.NewTableCell("On track").WithBackgroundColor("E6F4EA"),
				pptx.NewTableCell("Apr 12 review"),
			}).
			AddStyledRow([]pptx.TableCell{
				pptx.NewTableCell("Template refresh"),
				pptx.NewTableCell("Template Team"),
				pptx.NewTableCell("At risk").WithBackgroundColor("FDECEA"),
				pptx.NewTableCell("Resolve placeholder drift"),
			}),
	)

	return []pptx.SlideContent{hero, kpis, table}
}

func buildReskinSmokeSlides() []pptx.SlideContent {
	return []pptx.SlideContent{
		pptx.NewSlide("").WithBlankLayout().
			AddShape(
				pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0), pptx.Inches(0), pptx.Inches(13.33), pptx.Inches(7.5)).
					WithFill(pptx.NewShapeFill("F3F7FB")),
			).
			AddShape(
				pptx.NewShape(pptx.ShapeTypeRoundedRectangle, pptx.Inches(0.7), pptx.Inches(0.6), pptx.Inches(12), pptx.Inches(1.7)).
					WithGradientFill(
						pptx.NewShapeGradientFill(
							pptx.ShapeGradientTypeLinear,
							[]pptx.ShapeGradientStop{
								pptx.NewShapeGradientStop(0, "2B4C7E"),
								pptx.NewShapeGradientStop(100, "567EBB"),
							},
						).WithLinearAngle(35),
					).
					WithText("Brand Reskin Demo"),
			).
			AddShape(
				pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(0.9), pptx.Inches(3.0), pptx.Inches(2.2), pptx.Inches(1.5)).
					WithFill(pptx.NewShapeFill("003366")).
					WithText("Accent 1"),
			).
			AddShape(
				pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(3.35), pptx.Inches(3.0), pptx.Inches(2.2), pptx.Inches(1.5)).
					WithFill(pptx.NewShapeFill("00A3A3")).
					WithText("Accent 2"),
			).
			AddShape(
				pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(5.8), pptx.Inches(3.0), pptx.Inches(2.2), pptx.Inches(1.5)).
					WithFill(pptx.NewShapeFill("E67E22")).
					WithText("Accent 3"),
			).
			AddShape(
				pptx.NewShape(pptx.ShapeTypeRectangle, pptx.Inches(8.25), pptx.Inches(3.0), pptx.Inches(2.2), pptx.Inches(1.5)).
					WithFill(pptx.NewShapeFill("0B66D0")).
					WithText("Hyperlink"),
			).
			AddBullet("Theme preset + color/font schemes are applied in one pass."),
		pptx.NewSlide("Visual Consistency Checks").WithTable(
			pptx.NewTable([]pptx.Length{pptx.Inches(3.3), pptx.Inches(2.1), pptx.Inches(4.8)}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Surface").WithBold(true).WithBackgroundColor("DCE6F2"),
					pptx.NewTableCell("Result").WithBold(true).WithBackgroundColor("DCE6F2"),
					pptx.NewTableCell("Expectation").WithBold(true).WithBackgroundColor("DCE6F2"),
				}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Hero gradients + cards"),
					pptx.NewTableCell("Updated").WithBackgroundColor("E6F4EA"),
					pptx.NewTableCell("Swatches and banners adopt palette"),
				}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Typography"),
					pptx.NewTableCell("Updated").WithBackgroundColor("E6F4EA"),
					pptx.NewTableCell("Headings/body map to new theme fonts"),
				}).
				AddStyledRow([]pptx.TableCell{
					pptx.NewTableCell("Table header + status cells"),
					pptx.NewTableCell("Updated").WithBackgroundColor("E6F4EA"),
					pptx.NewTableCell("All fills stay readable after reskin"),
				}),
		),
	}
}
