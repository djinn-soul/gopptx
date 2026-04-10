// Package main demonstrates the native PDF export engine in gopptx.
// It creates one slide per feature category and exports directly to PDF
// without requiring LibreOffice or PowerPoint.
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/export"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

func main() {
	outDir := filepath.Join("examples", "output")
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		log.Fatalf("create output dir: %v", err)
	}
	pdfPath := filepath.Join(outDir, "81_native_pdf_showcase.pdf")

	slides := []elements.SlideContent{
		slideTitle(),
		slideTextFormatting(),
		slideBasicShapes(),
		slideGradientShapes(),
		slideArrowsAndCallouts(),
		slideConnectors(),
		slideTable(),
		slideBarChart(),
		slideLineAndPieCharts(),
		slideSmartArt(),
		slideSlideNumberAndFooter(),
	}

	opts := export.PDFOptions{Driver: export.PDFDriverNative}
	if err := export.PDFWithOptions("Native PDF Showcase", slides, pdfPath, opts); err != nil {
		log.Fatalf("PDF export failed: %v", err)
	}
	log.Printf("Written: %s\n", pdfPath)
}

// slideTitle is a centered title + subtitle slide.
func slideTitle() elements.SlideContent {
	s := elements.NewSlide("gopptx Native PDF Showcase")
	s.Layout = elements.SlideLayoutCenteredTitle
	s.TitleSize = 40
	s.TitleColor = "1F4E79"
	s.TitleBold = true
	s.Bullets = []string{
		"Pure-Go PDF rendering — no LibreOffice or PowerPoint required",
		"Shapes · Tables · Charts · SmartArt · Connectors",
	}

	bg := elements.NewSolidBackground("EBF3FB")
	s.Background = &bg
	return s
}

// slideTextFormatting shows per-bullet rich text styles: normal, bold, large, coloured, italic.
func slideTextFormatting() elements.SlideContent {
	s := elements.NewSlide("Text Formatting")
	s.TitleColor = "2E4057"
	s.TitleBold = true
	s.ContentColor = "333333" // slide-level default: dark grey
	s.ContentSize = 16        // slide-level default font size

	// Use Bullets as display labels and BulletRuns for the actual per-bullet style.
	s.Bullets = []string{
		"Normal bullet",
		"Bold bullet (ContentBold = true)",
		"Large content (ContentSize = 20)",
		"Custom color: deep blue title, dark-grey body",
		"Italic content text supported",
	}
	s.BulletRuns = [][]elements.Run{
		// 1 — normal: inherits slide ContentSize=16, ContentColor=333333
		{elements.NewRun("Normal bullet")},
		// 2 — bold
		{elements.NewRun("Bold bullet (ContentBold = true)").WithBold(true)},
		// 3 — larger font
		{{Text: "Large content (ContentSize = 20)", SizePt: 20}},
		// 4 — coloured text (deep blue on this bullet, dark-grey slide default on others)
		{{Text: "Custom color: deep blue title, dark-grey body", Color: "1F4E79"}},
		// 5 — italic
		{elements.NewRun("Italic content text supported").WithItalic(true)},
	}
	return s
}

// slideBasicShapes showcases core geometry shapes.
func slideBasicShapes() elements.SlideContent {
	s := elements.NewSlide("Basic Shapes")
	s.TitleColor = "1B4332"
	s.TitleBold = true

	// Row 1 — solid-filled shapes
	s.Shapes = append(s.Shapes,
		shapes.NewRectangle(0.5, 1.5, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("4CAF50")).
			WithLine(shapes.NewShapeLine("1B4332", styling.Points(1.5))).
			WithText("Rect"),

		shapes.NewRoundedRectangle(2.3, 1.5, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("2196F3")).
			WithLine(shapes.NewShapeLine("0D47A1", styling.Points(1.5))).
			WithText("Rounded"),

		shapes.NewEllipse(4.1, 1.5, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("FF9800")).
			WithLine(shapes.NewShapeLine("E65100", styling.Points(1.5))).
			WithText("Ellipse"),

		shapes.NewTriangle(5.9, 1.5, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("9C27B0")).
			WithLine(shapes.NewShapeLine("4A148C", styling.Points(1.5))).
			WithText("Triangle"),

		shapes.NewDiamond(7.7, 1.5, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("F44336")).
			WithLine(shapes.NewShapeLine("B71C1C", styling.Points(1.5))).
			WithText("Diamond"),
	)

	// Row 2 — more shapes
	s.Shapes = append(s.Shapes,
		shapes.NewHexagon(0.5, 3.2, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("00BCD4")).
			WithText("Hexagon"),

		shapes.NewPentagon(2.3, 3.2, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("FF5722")).
			WithText("Pentagon"),

		shapes.NewParallelogram(4.1, 3.2, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("607D8B")).
			WithText("Parallelgm"),

		shapes.NewTrapezoid(5.9, 3.2, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("795548")).
			WithText("Trapezoid"),

		shapes.NewOctagon(7.7, 3.2, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("009688")).
			WithText("Octagon"),
	)

	// Row 3 — stars, heart, wave
	s.Shapes = append(s.Shapes,
		shapes.NewStar(0.5, 4.9, 1.0).
			WithFill(shapes.NewShapeFill("FFD600")),

		shapes.NewStar8(2.3, 4.9, 1.0).
			WithFill(shapes.NewShapeFill("FF6F00")),

		shapes.NewHeart(4.2, 4.9, 1.0).
			WithFill(shapes.NewShapeFill("E91E63")),

		shapes.NewWave(5.9, 4.9, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("3F51B5")),

		shapes.NewCloud(7.7, 4.9, 1.5, 1.0).
			WithFill(shapes.NewShapeFill("B0BEC5")),
	)

	return s
}

// slideGradientShapes shows linear and radial gradient fills, plus rotation.
func slideGradientShapes() elements.SlideContent {
	s := elements.NewSlide("Gradient Fills & Rotation")
	s.TitleColor = "4A148C"
	s.TitleBold = true

	linearGrad := shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, []shapes.ShapeGradientStop{
		shapes.NewShapeGradientStop(0, "FF5252"),
		shapes.NewShapeGradientStop(100, "7C4DFF"),
	})

	radialGrad := shapes.NewShapeGradientFill(shapes.ShapeGradientTypeRadial, []shapes.ShapeGradientStop{
		shapes.NewShapeGradientStop(0, "FFFF00"),
		shapes.NewShapeGradientStop(100, "FF6D00"),
	})

	blueGrad := shapes.NewShapeGradientFill(shapes.ShapeGradientTypeLinear, []shapes.ShapeGradientStop{
		shapes.NewShapeGradientStop(0, "E3F2FD"),
		shapes.NewShapeGradientStop(100, "1565C0"),
	})

	s.Shapes = append(s.Shapes,
		// Linear gradient rectangle
		shapes.NewRectangle(0.5, 1.5, 2.5, 1.5).
			WithGradientFill(linearGrad).
			WithLine(shapes.NewShapeLine("333333", styling.Points(1))).
			WithText("Linear Gradient"),

		// Radial gradient ellipse
		shapes.NewEllipse(3.5, 1.5, 2.0, 1.5).
			WithGradientFill(radialGrad).
			WithText("Radial Gradient"),

		// Blue gradient rounded rect
		shapes.NewRoundedRectangle(6.0, 1.5, 3.0, 1.5).
			WithGradientFill(blueGrad).
			WithLine(shapes.NewShapeLine("0D47A1", styling.Points(1.5))).
			WithText("Gradient + Border"),

		// Rotated shapes
		shapes.NewRectangle(1.0, 3.5, 2.0, 1.0).
			WithFill(shapes.NewShapeFill("4CAF50")).
			WithRotation(30).
			WithText("30° Rotated"),

		shapes.NewRightArrow(4.5, 3.5, 2.5, 1.0).
			WithFill(shapes.NewShapeFill("FF9800")).
			WithRotation(45).
			WithText("45° Arrow"),

		shapes.NewStar(7.5, 3.5, 1.5).
			WithGradientFill(radialGrad).
			WithRotation(20),
	)

	return s
}

// slideArrowsAndCallouts shows arrow and callout shapes.
func slideArrowsAndCallouts() elements.SlideContent {
	s := elements.NewSlide("Arrows & Callouts")
	s.TitleColor = "B71C1C"
	s.TitleBold = true

	s.Shapes = append(s.Shapes,
		// Arrows
		shapes.NewRightArrow(0.5, 1.5, 2.0, 0.8).
			WithFill(shapes.NewShapeFill("F44336")).
			WithText("Right"),

		shapes.NewLeftArrow(3.0, 1.5, 2.0, 0.8).
			WithFill(shapes.NewShapeFill("2196F3")).
			WithText("Left"),

		shapes.NewUpArrow(5.5, 1.5, 0.8, 1.5).
			WithFill(shapes.NewShapeFill("4CAF50")).
			WithText("Up"),

		shapes.NewDownArrow(7.0, 1.5, 0.8, 1.5).
			WithFill(shapes.NewShapeFill("FF9800")).
			WithText("Down"),

		shapes.NewLeftRightArrow(0.5, 3.5, 2.5, 0.8).
			WithFill(shapes.NewShapeFill("9C27B0")).
			WithText("Left-Right"),

		shapes.NewUpDownArrow(3.5, 3.2, 0.8, 1.5).
			WithFill(shapes.NewShapeFill("00BCD4")).
			WithText("Up-Down"),

		shapes.NewChevron(4.8, 3.5, 2.0, 0.8).
			WithFill(shapes.NewShapeFill("FF5722")).
			WithText("Chevron"),

		// Callouts
		shapes.NewWedgeRectCallout(0.5, 5.0, 2.0, 1.0).
			WithFill(shapes.NewShapeFill("FFF9C4")).
			WithLine(shapes.NewShapeLine("F9A825", styling.Points(1))).
			WithText("Rect Callout"),

		shapes.NewWedgeEllipseCallout(3.0, 5.0, 2.0, 1.0).
			WithFill(shapes.NewShapeFill("E8F5E9")).
			WithLine(shapes.NewShapeLine("388E3C", styling.Points(1))).
			WithText("Ellipse Callout"),

		shapes.NewCloudCallout(5.5, 5.0, 2.5, 1.0).
			WithFill(shapes.NewShapeFill("E3F2FD")).
			WithLine(shapes.NewShapeLine("1976D2", styling.Points(1))).
			WithText("Cloud Callout"),
	)

	return s
}

// slideConnectors shows straight, elbow, and curved connectors.
func slideConnectors() elements.SlideContent {
	s := elements.NewSlide("Connectors")
	s.TitleColor = "1A237E"
	s.TitleBold = true

	// Anchor boxes
	box := func(x, y float64, color, label string) shapes.Shape {
		return shapes.NewRoundedRectangle(x, y, 1.4, 0.7).
			WithFill(shapes.NewShapeFill(color)).
			WithText(label)
	}

	s.Shapes = append(s.Shapes,
		box(0.5, 1.4, "1565C0", "Start A"),
		box(8.1, 1.4, "C62828", "End A"),
		box(0.5, 3.2, "2E7D32", "Start B"),
		box(8.1, 3.2, "AD1457", "End B"),
		box(0.5, 5.0, "4527A0", "Start C"),
		box(8.1, 5.0, "00695C", "End C"),
	)

	// Straight connector
	s.Connectors = append(s.Connectors,
		shapes.NewStraightConnector(
			styling.Inches(1.9), styling.Inches(1.75),
			styling.Inches(8.1), styling.Inches(1.75),
		).WithLine(shapes.NewShapeLine("1565C0", styling.Points(2))).
			WithLabel("Straight"),
	)

	// Elbow connector
	s.Connectors = append(s.Connectors,
		shapes.NewElbowConnector(
			styling.Inches(1.9), styling.Inches(3.55),
			styling.Inches(8.1), styling.Inches(3.55),
		).WithLine(shapes.NewShapeLine("2E7D32", styling.Points(2))).
			WithLabel("Elbow"),
	)

	// Curved connector
	s.Connectors = append(s.Connectors,
		shapes.NewCurvedConnector(
			styling.Inches(1.9), styling.Inches(5.35),
			styling.Inches(8.1), styling.Inches(5.35),
		).WithLine(shapes.NewShapeLine("4527A0", styling.Points(2))).
			WithLabel("Curved"),
	)

	return s
}

// slideTable shows a styled data table.
func slideTable() elements.SlideContent {
	s := elements.NewSlide("Table")
	s.TitleColor = "01579B"
	s.TitleBold = true

	colWidths := []styling.Length{
		styling.Inches(2.2),
		styling.Inches(1.8),
		styling.Inches(1.8),
		styling.Inches(1.8),
		styling.Inches(1.8),
	}
	tab := tables.NewTable(colWidths).Position(styling.Inches(0.3), styling.Inches(1.3))

	hdr := func(text string) tables.TableCell {
		return tables.TableCell{Text: text, Bold: true, BackgroundColor: "#01579B", Color: "#FFFFFF"}
	}
	tab = tab.AddStyledRow([]tables.TableCell{hdr("Product"), hdr("Q1"), hdr("Q2"), hdr("Q3"), hdr("Q4")})

	rows := [][]string{
		{"Alpha", "$1.2M", "$1.5M", "$1.8M", "$2.1M"},
		{"Beta", "$0.9M", "$1.1M", "$1.3M", "$1.6M"},
		{"Gamma", "$0.5M", "$0.7M", "$0.9M", "$1.2M"},
		{"Delta", "$0.3M", "$0.4M", "$0.6M", "$0.8M"},
	}
	altColors := []string{"#E3F2FD", "#FFFFFF"}
	for i, row := range rows {
		cells := make([]tables.TableCell, len(row))
		for j, v := range row {
			cells[j] = tables.TableCell{Text: v, BackgroundColor: altColors[i%2]}
		}
		tab = tab.AddStyledRow(cells)
	}

	s.Table = &tab
	return s
}

// slideBarChart demonstrates a vertical bar chart.
func slideBarChart() elements.SlideContent {
	s := elements.NewSlide("Bar Chart")
	s.TitleColor = "1B5E20"
	s.TitleBold = true

	bar := charts.NewBarChart(
		[]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
		[]float64{42, 58, 75, 63, 89, 97},
	)
	bar.Title = "Monthly Revenue (K)"
	bar.SeriesName = "Revenue"
	bar.BarColor = "43A047"
	bar.ShowLegend = true
	bar.LegendPosition = charts.LegendPositionBottom
	bar.ShowDataLabels = true
	bar.ShowMajorGridlines = true
	bar.CategoryAxisTitle = "Month"
	bar.ValueAxisTitle = "Revenue ($K)"

	s.Chart = &bar
	return s
}

// slideLineAndPieCharts shows a line chart and a pie chart side by side.
func slideLineAndPieCharts() elements.SlideContent {
	s := elements.NewSlide("Line Chart & Pie Chart")
	s.TitleColor = "4E342E"
	s.TitleBold = true

	// Line chart (left half)
	line := charts.NewLineChart(
		[]string{"2019", "2020", "2021", "2022", "2023"},
		[]float64{55, 68, 72, 85, 99},
	)
	line.Title = "Year-over-Year Growth"
	line.SeriesName = "Growth %"
	line.ShowLegend = true
	line.LegendPosition = charts.LegendPositionBottom
	line.ShowMajorGridlines = true
	line.X = styling.Inches(0.2)
	line.Y = styling.Inches(1.2)
	line.CX = styling.Inches(4.5)
	line.CY = styling.Inches(4.5)
	s.Line = &line

	// Pie chart (right half)
	pie := charts.NewPieChart(
		[]string{"North", "South", "East", "West"},
		[]float64{35, 25, 22, 18},
	)
	pie.Title = "Regional Split"
	pie.SeriesName = "Share"
	pie.ShowLegend = true
	pie.LegendPosition = charts.LegendPositionBottom
	pie.ShowDataLabels = true
	pie.X = styling.Inches(5.0)
	pie.Y = styling.Inches(1.2)
	pie.CX = styling.Inches(4.5)
	pie.CY = styling.Inches(4.5)
	s.Pie = &pie

	return s
}

// slideSmartArt shows a process SmartArt and an org-chart hierarchy.
func slideSmartArt() elements.SlideContent {
	s := elements.NewSlide("SmartArt Diagrams")
	s.TitleColor = "880E4F"
	s.TitleBold = true

	// Basic process list
	process := smartart.NewSmartArt(smartart.BasicProcess).
		AddItems([]string{"Discover", "Design", "Build", "Test", "Deploy"})
	process.X = styling.Inches(0.3)
	process.Y = styling.Inches(1.2)
	process.CX = styling.Inches(9.4)
	process.CY = styling.Inches(2.0)

	// Org-chart hierarchy
	orgchart := smartart.NewSmartArt(smartart.OrgChart).
		AddNode(
			smartart.NewNode("CEO").
				WithChild(smartart.NewNode("Engineering")).
				WithChild(smartart.NewNode("Marketing")).
				WithChild(smartart.NewNode("Sales")),
		)
	orgchart.X = styling.Inches(0.3)
	orgchart.Y = styling.Inches(3.4)
	orgchart.CX = styling.Inches(9.4)
	orgchart.CY = styling.Inches(2.8)

	s.SmartArtDiagrams = []smartart.SmartArt{process, orgchart}
	return s
}

// slideSlideNumberAndFooter shows slide numbering and footer text.
func slideSlideNumberAndFooter() elements.SlideContent {
	s := elements.NewSlide("Slide Numbers & Footer")
	s.TitleColor = "37474F"
	s.TitleBold = true
	s.ShowSlideNumber = true
	s.FooterText = "gopptx Native PDF Showcase  |  github.com/djinn-soul/gopptx"
	s.Bullets = []string{
		"This slide shows the slide-number indicator (bottom-right corner).",
		"The footer text appears centred at the bottom of the page.",
		"Both are rendered natively without external tools.",
	}

	bg := elements.NewSolidBackground("ECEFF1")
	s.Background = &bg
	return s
}
