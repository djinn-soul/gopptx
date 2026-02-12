package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func main() {
	outDir := "smoke_samples"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	generators := []struct {
		name string
		fn   func() ([]byte, error)
	}{
		{"01_basic_pptx", generateBasicPPTX},
		{"02_slide_layouts", generateSlideLayouts},
		{"04_text_formatting", generateTextFormatting},
		{"05_bullet_styles", generateBulletStyles},
		{"06_text_enhancements", generateTextEnhancements},
		{"07_tables", generateTables},
		{"08_table_cell_merge", generateTableMerge},
		{"09_charts", generateCharts},
		{"10_images", generateImages},
		{"11_images_advanced", generateImagesAdvanced},
		{"04_placeholders", generatePlaceholders},
		{"12_shapes", generateShapes},
		{"13_connectors", generateConnectors},
		{"14_transitions", generateTransitions},
		{"19_read_modify", generateReadModify},
		{"22_speaker_notes", generateSpeakerNotes},
		{"28_animations", generateAnimations},
		{"31_hyperlinks", generateHyperlinks},
		{"15_merge", generateMerge},
		{"16_prelude_helpers", generatePreludeHelpers},
		{"03_markdown", generateMarkdown},
		{"18_layout_helpers", generateLayoutHelpers},
		{"48_accessibility", generateAccessibility},
		{"53_slide_properties", generateSlideProperties},
		{"54_theme_master", generateThemeMaster},
		{"55_background_fills", generateSlideBackgrounds},
		{"56_action_api", generateActionAPI},
	}

	for _, g := range generators {
		data, err := g.fn()
		if err != nil {
			log.Printf("Error generating %s: %v", g.name, err)
			continue
		}
		path := filepath.Join(outDir, g.name+".pptx")
		if err := os.WriteFile(path, data, 0o644); err != nil {
			log.Printf("Error writing %s: %v", path, err)
			continue
		}
		fmt.Printf("Generated %s\n", path)
	}
}

func generateBasicPPTX() ([]byte, error) {
	slide := pptx.NewSlide("Basic PPTX Generation").
		AddBullet("Simple slide creation").
		AddBullet("Title and bullet points").
		AddBullet("Task 01 complete")
	return pptx.CreateWithSlides("Task 01: Basic PPTX", []pptx.SlideContent{slide})
}

func generateSlideLayouts() ([]byte, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Title and Content Layout").
			WithLayout(pptx.SlideLayoutTitleAndContent).
			AddBullet("Default layout with title and content"),
		pptx.NewSlide("Title Only Layout").
			WithLayout(pptx.SlideLayoutTitleOnly),
		pptx.NewSlide("Two Column Layout").
			WithLayout(pptx.SlideLayoutTwoColumn).
			AddBullet("Left column content").
			AddBullet("Right column content"),
		pptx.NewSlide("").
			WithLayout(pptx.SlideLayoutBlank),
	}
	return pptx.CreateWithSlides("Task 02: Slide Layouts", slides)
}

func generateTextFormatting() ([]byte, error) {
	slide := pptx.NewSlide("Text Formatting").
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Bold ").WithBold(true),
			pptx.NewTextRun("Italic ").WithItalic(true),
			pptx.NewTextRun("Underline ").WithUnderline(true),
			pptx.NewTextRun("Strikethrough").WithStrikethrough(true),
		}).
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Red ").WithColor("FF0000"),
			pptx.NewTextRun("Green ").WithColor("00FF00"),
			pptx.NewTextRun("Blue").WithColor("0000FF"),
		}).
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Large ").WithSizePt(24),
			pptx.NewTextRun("Small").WithSizePt(10),
		})
	return pptx.CreateWithSlides("Task 04: Text Formatting", []pptx.SlideContent{slide})
}

func generateBulletStyles() ([]byte, error) {
	slide := pptx.NewSlide("Bullet Styles").
		AddBullet("Default bullet style").
		AddBullet("Another bullet point").
		AddBullet("Third bullet point")

	slide2 := pptx.NewSlide("Numbered Bullets").
		AddBulletWithStyle("First item", pptx.NewTextParagraphStyle().WithBulletStyle(pptx.BulletStyleNumber)).
		AddBulletWithStyle("Second item", pptx.NewTextParagraphStyle().WithBulletStyle(pptx.BulletStyleNumber)).
		AddBulletWithStyle("Third item", pptx.NewTextParagraphStyle().WithBulletStyle(pptx.BulletStyleNumber))

	return pptx.CreateWithSlides("Task 05: Bullet Styles", []pptx.SlideContent{slide, slide2})
}

func generateTextEnhancements() ([]byte, error) {
	slide := pptx.NewSlide("Text Enhancements").
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("H"),
			pptx.NewTextRun("2").WithSubscript(true),
			pptx.NewTextRun("O - Subscript"),
		}).
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("E=mc"),
			pptx.NewTextRun("2").WithSuperscript(true),
			pptx.NewTextRun(" - Superscript"),
		}).
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Highlighted text").WithHighlight("FFFF00"),
		}).
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Code style").WithCode(true),
		})
	return pptx.CreateWithSlides("Task 06: Text Enhancements", []pptx.SlideContent{slide})
}

func generateTables() ([]byte, error) {
	table := pptx.NewTable([]int64{2000000, 2000000, 2000000}).
		AddRow([]string{"Header 1", "Header 2", "Header 3"}).
		AddRow([]string{"Row 1, Col 1", "Row 1, Col 2", "Row 1, Col 3"}).
		AddRow([]string{"Row 2, Col 1", "Row 2, Col 2", "Row 2, Col 3"})

	slide := pptx.NewSlide("Tables").WithTable(table)
	return pptx.CreateWithSlides("Task 07: Tables", []pptx.SlideContent{slide})
}

func generateTableMerge() ([]byte, error) {
	table := pptx.NewTable([]int64{2000000, 2000000, 2000000}).
		AddStyledRow([]pptx.TableCell{
			pptx.NewTableCell("Merged Header").WithColSpan(3),
			pptx.NewTableCell(""),
			pptx.NewTableCell(""),
		}).
		AddRow([]string{"A", "B", "C"}).
		AddRow([]string{"D", "E", "F"})

	slide := pptx.NewSlide("Table Cell Merge").WithTable(table)
	return pptx.CreateWithSlides("Task 08: Table Merge", []pptx.SlideContent{slide})
}

func generateCharts() ([]byte, error) {
	var slides []pptx.SlideContent

	// Bar Chart
	barChart := pptx.NewBarChart(
		[]string{"Q1", "Q2", "Q3", "Q4"},
		[]float64{100, 200, 150, 300},
	).WithTitle("Quarterly Performance")
	slide1 := pptx.NewSlide("Bar Chart")
	slide1.Chart = &barChart
	slides = append(slides, slide1)

	// Line Chart
	lineChart := pptx.NewLineChart(
		[]string{"Jan", "Feb", "Mar", "Apr"},
		[]float64{10, 15, 13, 17},
	).WithTitle("Monthly Trends")
	slide2 := pptx.NewSlide("Line Chart")
	slide2.Line = &lineChart
	slides = append(slides, slide2)

	// Pie Chart
	pieChart := pptx.NewPieChart(
		[]string{"East", "West", "North", "South"},
		[]float64{40, 35, 15, 10},
	).WithTitle("Regional Distribution")
	slide3 := pptx.NewSlide("Pie Chart")
	slide3.Pie = &pieChart
	slides = append(slides, slide3)

	// Doughnut Chart
	doughnutChart := pptx.NewDoughnutChart(
		[]string{"Product A", "Product B", "Product C"},
		[]float64{50, 30, 20},
	).WithTitle("Product Mix")
	slide4 := pptx.NewSlide("Doughnut Chart")
	slide4.Doughnut = &doughnutChart
	slides = append(slides, slide4)

	// Scatter Chart
	scatterChart := pptx.NewScatterChart(
		[]float64{1, 2, 3, 4, 5},
		[]float64{10, 25, 30, 45, 60},
	).WithTitle("Scatter Plot")
	slide5 := pptx.NewSlide("Scatter Chart")
	slide5.Scatter = &scatterChart
	slides = append(slides, slide5)

	// Area Chart
	areaChart := pptx.NewAreaChart(
		[]string{"2020", "2021", "2022", "2023"},
		[]float64{500, 700, 1200, 1500},
	).WithTitle("Revenue Growth")
	slide6 := pptx.NewSlide("Area Chart")
	slide6.Area = &areaChart
	slides = append(slides, slide6)

	// Radar Chart (Marker)
	radarChart := pptx.NewRadarChart(
		[]string{"Speed", "Power", "Durability", "Range"},
		[]float64{8, 9, 7, 6},
	).WithTitle("Attributes")
	slide7 := pptx.NewSlide("Radar Chart")
	slide7.Radar = &radarChart
	slides = append(slides, slide7)

	// Radar Chart (Filled)
	radarFilledChart := pptx.NewRadarFilledChart(
		[]string{"Speed", "Power", "Durability", "Range"},
		[]float64{5, 6, 8, 9},
	).WithTitle("Attributes (Filled)")
	slide8 := pptx.NewSlide("Radar Chart (Filled)")
	slide8.RadarFilled = &radarFilledChart
	slides = append(slides, slide8)

	// Bubble Chart
	bubbleChart := pptx.NewBubbleChart(
		[]float64{10, 20, 30, 40},
		[]float64{5, 10, 15, 20},
		[]float64{2, 4, 6, 8},
	).WithTitle("Bubble Chart")
	slide9 := pptx.NewSlide("Bubble Chart")
	slide9.Bubble = &bubbleChart
	slides = append(slides, slide9)

	// Stock HLC Chart
	stockHLC := pptx.NewStockHLCChart(
		[]string{"Day 1", "Day 2", "Day 3"},
		[]float64{100, 110, 105}, // High
		[]float64{90, 95, 98},    // Low
		[]float64{95, 100, 102},  // Close
	)
	stockHLC.Title = "Stock HLC"
	slide10 := pptx.NewSlide("Stock HLC Chart")
	slide10.StockHLC = &stockHLC
	slides = append(slides, slide10)

	// Stock OHLC Chart
	stockOHLC := pptx.NewStockOHLCChart(
		[]string{"Day 1", "Day 2", "Day 3"},
		[]float64{92, 98, 100},   // Open
		[]float64{100, 110, 105}, // High
		[]float64{90, 95, 98},    // Low
		[]float64{95, 100, 102},  // Close
	)
	stockOHLC.Title = "Stock OHLC"
	slide11 := pptx.NewSlide("Stock OHLC Chart")
	slide11.StockOHLC = &stockOHLC
	slides = append(slides, slide11)

	// Combo Chart
	comboChart := pptx.NewComboChart(
		[]string{"2020", "2021", "2022", "2023"},
		[]pptx.Series{
			{Name: "Revenue", Values: []float64{100, 150, 200, 250}},
		},
		[]pptx.Series{
			{Name: "Growth", Values: []float64{10, 50, 33, 25}},
		},
	)
	comboChart.Title = "Company Performance"
	slide12 := pptx.NewSlide("Combo Chart")
	slide12.Combo = &comboChart
	slides = append(slides, slide12)

	return pptx.CreateWithSlides("Task 09: Charts", slides)
}

func generateImages() ([]byte, error) {
	imagePath := filepath.Join("smoke_samples", "sampleimage", "repository-open-graph-template.png")
	// If path doesn't exist relative to where script is run, fall back to dummy
	_, err := os.Stat(imagePath)
	var img pptx.Image
	if err == nil {
		img = pptx.NewImage(imagePath, 4*914400, 1*914400, 4*914400, 2*914400)
	} else {
		img = pptx.NewImageFromBytes([]byte("fake png"), "png", 4*914400, 1*914400, 4*914400, 2*914400)
	}

	slide := pptx.NewSlide("Images").
		AddBullet("Image embedding supported").
		AddBullet("PNG, JPEG, GIF formats").
		AddBullet("Position and size control").
		AddImage(img)

	return pptx.CreateWithSlides("Task 10: Images", []pptx.SlideContent{slide})
}

func generateImagesAdvanced() ([]byte, error) {
	imagePath := filepath.Join("smoke_samples", "sampleimage", "repository-open-graph-template.png")
	_, err := os.Stat(imagePath)
	var img1, img2 pptx.Image
	if err == nil {
		img1 = pptx.NewImage(imagePath, 500000, 2000000, 2000000, 2000000).
			WithRotation(15).
			WithFlip(true, false)
		img2 = pptx.NewImage(imagePath, 3000000, 2000000, 2000000, 2000000).
			WithCrop(0.1, 0.1, 0.1, 0.1)
	} else {
		img1 = pptx.NewImageFromBytes([]byte("fake png"), "png", 500000, 2000000, 2000000, 2000000).
			WithRotation(15).
			WithFlip(true, false)
		img2 = pptx.NewImageFromBytes([]byte("fake png"), "png", 3000000, 2000000, 2000000, 2000000).
			WithCrop(0.1, 0.1, 0.1, 0.1)
	}

	slide := pptx.NewSlide("Advanced Image Sources").
		AddBullet("Bytes, Base64, and URL sources supported").
		AddBullet("Rotation, Flip, and Crop effects").
		AddImage(img1).
		AddImage(img2)

	return pptx.CreateWithSlides("Task 11: Advanced Images", []pptx.SlideContent{slide})
}

func generatePlaceholders() ([]byte, error) {
	imagePath := filepath.Join("smoke_samples", "sampleimage", "repository-open-graph-template.png")
	_, statErr := os.Stat(imagePath)
	var img pptx.Image
	if statErr == nil {
		img = pptx.NewImage(imagePath, 0, 0, 0, 0)
	} else {
		img = pptx.NewImageFromBytes([]byte("fake png"), "png", 0, 0, 0, 0)
	}

	slides := []pptx.SlideContent{
		pptx.NewSlide("Placeholder Overrides").
			WithPlaceholderText(0, "title", "Title Override").
			WithPlaceholderText(1, "body", "Body content override from code"),
		pptx.NewSlide("Placeholder Image").
			WithPlaceholderImage(1, "picture", img),
		pptx.NewSlide("Placeholder Table").
			WithPlaceholderTable(1, "body", pptx.Table{
				Rows: [][]string{
					{"PH Col 1", "PH Col 2"},
					{"PH Data 1", "PH Data 2"},
				},
			}),
	}
	return pptx.CreateWithSlides("Task 04: Placeholders", slides)
}

func generateShapes() ([]byte, error) {
	slide := pptx.NewSlide("Shapes").
		AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, 500000, 1500000, 2000000, 1000000).
			WithFill(pptx.NewShapeFill("FF6600")).
			WithText("Rectangle")).
		AddShape(pptx.NewShape(pptx.ShapeTypeEllipse, 3000000, 1500000, 1500000, 1000000).
			WithFill(pptx.NewShapeFill("0066FF")).
			WithText("Ellipse")).
		AddShape(pptx.NewShape(pptx.ShapeTypeTriangle, 5000000, 1500000, 1500000, 1000000).
			WithFill(pptx.NewShapeFill("00CC00")).
			WithText("Triangle"))

	return pptx.CreateWithSlides("Task 12: Shapes", []pptx.SlideContent{slide})
}

func generateConnectors() ([]byte, error) {
	slide := pptx.NewSlide("Connectors").
		AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, 500000, 2000000, 1500000, 800000).
			WithFill(pptx.NewShapeFill("3366CC")).WithText("Start")).
		AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, 4000000, 2000000, 1500000, 800000).
			WithFill(pptx.NewShapeFill("CC6633")).WithText("End")).
		AddConnector(pptx.NewConnector(pptx.ConnectorTypeStraight, 2000000, 2400000, 4000000, 2400000).
			WithLine(pptx.NewShapeLine("333333", 12700)).
			WithArrows(pptx.ArrowTypeNone, pptx.ArrowTypeTriangle))

	return pptx.CreateWithSlides("Task 13: Connectors", []pptx.SlideContent{slide})
}

func generateTransitions() ([]byte, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Fade Transition").
			WithTransition(pptx.TransitionFade).
			AddBullet("This slide fades in"),
		pptx.NewSlide("Push Transition (Left)").
			WithTransitionOptions(pptx.TransitionOptions{
				Type:      pptx.TransitionPush,
				Direction: "l",
			}).
			AddBullet("This slide pushes from the left"),
		pptx.NewSlide("Wipe Transition (Up)").
			WithTransitionOptions(pptx.TransitionOptions{
				Type:      pptx.TransitionWipe,
				Direction: "u",
			}).
			AddBullet("This slide wipes from the bottom up"),
		pptx.NewSlide("Strips Transition (Right-Down)").
			WithTransitionOptions(pptx.TransitionOptions{
				Type:      pptx.TransitionStrips,
				Direction: pptx.TransitionDirDownRight,
			}).
			AddBullet("This slide uses strips from top-left to bottom-right"),
		pptx.NewSlide("Clock Transition (3 spokes)").
			WithTransitionOptions(pptx.TransitionOptions{
				Type:       pptx.TransitionClock,
				SpokeCount: 3,
			}).
			AddBullet("This slide uses a wheel/clock transition with 3 spokes"),
	}
	return pptx.CreateWithSlides("Task 14: Transitions", slides)
}

func generateMerge() ([]byte, error) {
	// Create first presentation
	s1 := pptx.NewSlide("Presentation One").
		AddBullet("Slide from the first presentation")
	data1, err := pptx.CreateWithSlides("Merge Target", []pptx.SlideContent{s1})
	if err != nil {
		return nil, err
	}

	// Create second presentation
	s2 := pptx.NewSlide("Presentation Two").
		AddBullet("Slide from the second presentation")
	data2, err := pptx.CreateWithSlides("Merge Source", []pptx.SlideContent{s2})
	if err != nil {
		return nil, err
	}

	// Save data2 to a temp file because MergeFromFile needs a path
	tmpFile := filepath.Join(os.TempDir(), "gopptx_merge_source.pptx")
	if err := os.WriteFile(tmpFile, data2, 0o644); err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(tmpFile) }()

	// Save data1 to a temp file because OpenPresentationEditor needs a path
	targetFile := filepath.Join(os.TempDir(), "gopptx_merge_target.pptx")
	if err := os.WriteFile(targetFile, data1, 0o644); err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(targetFile) }()

	// Open first and merge second
	editor, err := pptx.OpenPresentationEditor(targetFile)
	if err != nil {
		return nil, err
	}

	if err := editor.MergeFromFile(tmpFile); err != nil {
		return nil, err
	}

	// Save back to a new temp file to read bytes
	mergedFile := filepath.Join(os.TempDir(), "gopptx_merged.pptx")
	if err := editor.Save(mergedFile); err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(mergedFile) }()

	return os.ReadFile(mergedFile)
}

func generateReadModify() ([]byte, error) {
	slide := pptx.NewSlide("Read/Modify Existing").
		AddBullet("Open existing PPTX files").
		AddBullet("Modify slide content").
		AddBullet("Add/remove slides").
		AddBullet("Save changes")
	return pptx.CreateWithSlides("Task 19: Read/Modify", []pptx.SlideContent{slide})
}

func generateSpeakerNotes() ([]byte, error) {
	slide := pptx.NewSlide("Speaker Notes").
		WithNotes("These are speaker notes that appear in presenter view.\n\nKey points:\n- First point\n- Second point").
		AddBullet("Speaker notes supported").
		AddBullet("Visible in presenter view")
	return pptx.CreateWithSlides("Task 22: Speaker Notes", []pptx.SlideContent{slide})
}

func generateAnimations() ([]byte, error) {
	slide := pptx.NewSlide("Animations").
		AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, 1000000, 2000000, 2000000, 1000000).
			WithFill(pptx.NewShapeFill("FF6600")).WithText("Fade In")).
		AddShape(pptx.NewShape(pptx.ShapeTypeEllipse, 4000000, 2000000, 2000000, 1000000).
			WithFill(pptx.NewShapeFill("0066FF")).WithText("Fly In")).
		AddAnimation(pptx.NewAnimation(1, pptx.AnimationEntranceFade)).
		AddAnimation(pptx.NewAnimation(2, pptx.AnimationEntranceFlyIn).WithTrigger(pptx.AnimationAfterPrevious))

	return pptx.CreateWithSlides("Task 28: Animations", []pptx.SlideContent{slide})
}

func generateHyperlinks() ([]byte, error) {
	slide := pptx.NewSlide("Hyperlinks").
		AddShape(pptx.NewShape(pptx.ShapeTypeRoundedRectangle, 1000000, 2000000, 3000000, 800000).
			WithFill(pptx.NewShapeFill("0066CC")).
			WithText("Click to visit example.com").
			WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com")).WithTooltip("Open website"))).
		AddBullet("Shape with URL hyperlink").
		AddBullet("Tooltip on hover").
		AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Text hyperlink: "),
			pptx.NewTextRun("Click here to visit example.com").
				WithHyperlink(pptx.NewHyperlink(pptx.HyperlinkURL("https://example.com"))).
				WithColor("0000FF").
				WithUnderline(true),
		})

	return pptx.CreateWithSlides("Task 31: Hyperlinks", []pptx.SlideContent{slide})
}

func generateMarkdown() ([]byte, error) {
	md := `# Markdown Slide
- Bullet point
- **Bold** and *italic*
---
# Another Slide
1. Ordered list
2. item 2`
	slides, err := pptx.SlidesFromMarkdown(md)
	if err != nil {
		return nil, err
	}
	return pptx.CreateWithSlides("Task 03: Markdown", slides)
}

func generateLayoutHelpers() ([]byte, error) {
	slide := pptx.NewSlide("Layout Helpers")
	boxes, _ := pptx.Grid(2, 2, pptx.Inches(0.5))
	for i, b := range boxes {
		slide.AddShape(pptx.NewShape(pptx.ShapeTypeRectangle, b.X, b.Y, b.CX, b.CY).
			WithText(fmt.Sprintf("Grid %d", i+1)).
			WithFill(pptx.NewShapeFill("4472C4")))
	}
	return pptx.CreateWithSlides("Task 18: Layout Helpers", []pptx.SlideContent{slide})
}

func generateAccessibility() ([]byte, error) {
	img := pptx.NewImageFromBytes([]byte("fake"), "png", 1000000, 1000000, 2000000, 2000000).
		WithAltText("A descriptive text for the image").
		WithDecorative(false)

	shape := pptx.NewShape(pptx.ShapeTypeRectangle, 4000000, 1000000, 2000000, 2000000).
		WithFill(pptx.NewShapeFill("70AD47")).
		WithAltText("Decorative shape").
		WithDecorative(true)

	slide := pptx.NewSlide("Accessibility").
		AddImage(img).
		AddShape(shape).
		AddBullet("Image has AltText").
		AddBullet("Shape is marked as decorative")

	return pptx.CreateWithSlides("Task 48: Accessibility", []pptx.SlideContent{slide})
}

func generateSlideProperties() ([]byte, error) {
	slides := []pptx.SlideContent{
		pptx.NewSlide("Styled Slide").
			WithBackgroundColor("D9E1F2").
			WithTitleAlign("r").
			WithTitleFont("Consolas").
			WithContentVAlign("ctr").
			WithSlideNumber(true).
			AddBullet("Background: D9E1F2").
			AddBullet("Title: Right-aligned, Consolas").
			AddBullet("Content: Middle-aligned").
			AddBullet("Slide numbers: Enabled"),
	}
	return pptx.CreateWithSlides("Task 53: Slide Properties", slides)
}

func generateThemeMaster() ([]byte, error) {
	neonTheme := styling.Theme{
		Name: "NeonStream",
		Colors: styling.ColorScheme{
			Name:    "NeonStream Colors",
			Dk1:     "000000",
			Lt1:     "FFFFFF",
			Accent1: "00FFFF",
			Accent2: "FF00FF",
		},
		Fonts: styling.FontScheme{
			Name:      "Modern Tech",
			MajorFont: "Inter",
			MinorFont: "Roboto",
		},
	}

	neonGradient := pptx.NewShapeGradientFill(pptx.ShapeGradientTypeLinear, []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "000000"),
		pptx.NewShapeGradientStop(100, "1A1A1A"),
	})

	master := elements.NewMaster().
		WithBackground(elements.NewGradientBackground(neonGradient)).
		WithFooter("© 2026 NeonStream Technology - Confidential").
		WithColorMapping("dk1", "lt1")

	master.AddShape(pptx.NewRectangle(0, 0, 13.33, 0.05).WithFill(pptx.NewShapeFill("00FFFF")))
	master.AddShape(pptx.NewRectangle(0, 7.45, 13.33, 0.05).WithFill(pptx.NewShapeFill("FF00FF")))

	builder := pptx.NewPresentationBuilder("Task 54: Theme & Master").
		WithTheme(neonTheme).
		WithMaster(master).
		WithSlideSize(pptx.SlideSize16x9)

	builder.AddTitleSlide("Neon Theme Showcase")
	builder.AddBulletSlide("Features", []string{"Custom Colors", "Global Footer", "Master Shapes"})

	return builder.Build()
}

func generateSlideBackgrounds() ([]byte, error) {
	builder := pptx.NewPresentationBuilder("Task 55: Slide Backgrounds")

	// Solid
	builder.AddSlide(pptx.NewSlide("Solid Background").WithBackgroundColor("FF9900"))

	// Gradient
	grad := pptx.NewShapeGradientFill(pptx.ShapeGradientTypeLinear, []pptx.ShapeGradientStop{
		pptx.NewShapeGradientStop(0, "FFEE00"),
		pptx.NewShapeGradientStop(100, "FF0000"),
	})
	builder.AddSlide(pptx.NewSlide("Gradient Background").WithGradientBackground(grad))

	return builder.Build()
}

func generateActionAPI() ([]byte, error) {
	builder := pptx.NewPresentationBuilder("Task 56: Action API")

	slide := pptx.NewSlide("Interactive Shapes").
		AddShape(pptx.NewRectangle(1, 1, 3, 2).
			WithText("Click Me (GitHub)").
			WithFill(pptx.NewShapeFill("00FF00")).
			WithClickAction(pptx.NewHyperlink(pptx.HyperlinkURL("https://github.com/djinn-soul/gopptx"))))

	builder.AddSlide(slide)
	return builder.Build()
}
