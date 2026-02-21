package pptx_test

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

func TestAccessibility(t *testing.T) {
	// Setup image data
	imgData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	tests := []struct {
		name       string
		setup      func() (pptx.SlideContent, string) // Returns slide and expected slide XML name
		expected   []string                           // Strings expected to be present
		unexpected []string                           // Strings expected to be absent
	}{
		// ---------------- Shapes ----------------
		{
			name: "Shape with AltText",
			setup: func() (pptx.SlideContent, string) {
				shape := pptx.NewRectangle(1, 1, 1, 1).
					WithAltText("Shape Alt").
					WithDecorative(false)
				return pptx.NewSlide("Slide 1").AddShape(shape), "ppt/slides/slide1.xml"
			},
			expected: []string{`descr="Shape Alt" title="Shape Alt"`},
		},
		{
			name: "Shape Decorative",
			setup: func() (pptx.SlideContent, string) {
				shape := pptx.NewRectangle(1, 1, 1, 1).
					WithAltText("Ignored Text").
					WithDecorative(true)
				return pptx.NewSlide("Slide 2").AddShape(shape), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
		{
			name: "Shape Empty AltText",
			setup: func() (pptx.SlideContent, string) {
				shape := pptx.NewRectangle(1, 1, 1, 1).
					WithAltText("").
					WithDecorative(false)
				return pptx.NewSlide("Slide 3").AddShape(shape), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},

		// ---------------- Images ----------------
		{
			name: "Image with AltText",
			setup: func() (pptx.SlideContent, string) {
				img := pptx.NewImageFromBytes(imgData, "png", 0, 0, 100, 100).
					WithAltText("Image Alt").
					WithDecorative(false)
				return pptx.NewSlide("Slide 4").AddImage(img), "ppt/slides/slide1.xml"
			},
			expected: []string{`descr="Image Alt" title="Image Alt"`},
		},
		{
			name: "Image Decorative",
			setup: func() (pptx.SlideContent, string) {
				img := pptx.NewImageFromBytes(imgData, "png", 0, 0, 100, 100).
					WithAltText("Ignored").
					WithDecorative(true)
				return pptx.NewSlide("Slide 5").AddImage(img), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
		{
			name: "Image Empty AltText",
			setup: func() (pptx.SlideContent, string) {
				img := pptx.NewImageFromBytes(imgData, "png", 0, 0, 100, 100).
					WithAltText("").
					WithDecorative(false)
				return pptx.NewSlide("Slide 6").AddImage(img), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},

		// ---------------- Connectors ----------------
		{
			name: "Connector with AltText",
			setup: func() (pptx.SlideContent, string) {
				conn := pptx.NewStraightConnector(1, 1, 2, 2).
					WithAltText("Connector Alt").
					WithDecorative(false)
				return pptx.NewSlide("Slide 7").AddConnector(conn), "ppt/slides/slide1.xml"
			},
			expected: []string{`descr="Connector Alt" title="Connector Alt"`},
		},
		{
			name: "Connector Decorative",
			setup: func() (pptx.SlideContent, string) {
				conn := pptx.NewStraightConnector(1, 1, 2, 2).
					WithAltText("Ignored").
					WithDecorative(true)
				return pptx.NewSlide("Slide 8").AddConnector(conn), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
		{
			name: "Connector Empty AltText",
			setup: func() (pptx.SlideContent, string) {
				conn := pptx.NewStraightConnector(1, 1, 2, 2).
					WithAltText("").
					WithDecorative(false)
				return pptx.NewSlide("Slide 9").AddConnector(conn), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},

		// ---------------- Tables ----------------
		{
			name: "Table with AltText",
			setup: func() (pptx.SlideContent, string) {
				table := pptx.NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).
					WithAltText("Table Alt").
					WithDecorative(false)
				table = table.AddRow([]string{"A", "B"})
				return pptx.NewSlide("Slide 10").WithTable(table), "ppt/slides/slide1.xml"
			},
			expected: []string{`descr="Table Alt" title="Table Alt"`},
		},
		{
			name: "Table Decorative",
			setup: func() (pptx.SlideContent, string) {
				table := pptx.NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).
					WithAltText("Ignored").
					WithDecorative(true)
				table = table.AddRow([]string{"A", "B"})
				return pptx.NewSlide("Slide 11").WithTable(table), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
		{
			name: "Table Empty AltText",
			setup: func() (pptx.SlideContent, string) {
				table := pptx.NewTable([]styling.Length{styling.Inches(1), styling.Inches(1)}).
					WithAltText("").
					WithDecorative(false)
				table = table.AddRow([]string{"A", "B"})
				return pptx.NewSlide("Slide 12").WithTable(table), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},

		// ---------------- Charts ----------------
		{
			name: "Chart with AltText",
			setup: func() (pptx.SlideContent, string) {
				chart := pptx.NewBarChart([]string{"Cat1"}, []float64{10}).
					WithAltText("Chart Alt").
					WithDecorative(false)
				return pptx.NewSlide("Slide 13").WithBarChart(chart), "ppt/slides/slide1.xml"
			},
			expected: []string{`descr="Chart Alt" title="Chart Alt"`},
		},
		{
			name: "Chart Decorative",
			setup: func() (pptx.SlideContent, string) {
				chart := pptx.NewBarChart([]string{"Cat1"}, []float64{10}).
					WithAltText("Ignored").
					WithDecorative(true)
				return pptx.NewSlide("Slide 14").WithBarChart(chart), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
		{
			name: "Chart Empty AltText",
			setup: func() (pptx.SlideContent, string) {
				chart := pptx.NewBarChart([]string{"Cat1"}, []float64{10}).
					WithAltText("").
					WithDecorative(false)
				return pptx.NewSlide("Slide 15").WithBarChart(chart), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},

		// ---------------- SmartArt ----------------
		{
			name: "SmartArt with AltText",
			setup: func() (pptx.SlideContent, string) {
				sa := smartart.NewSmartArt(smartart.BasicBlockList).
					WithAltText("SmartArt Alt").
					WithDecorative(false).
					AddNode(smartart.NewNode("Node 1"))
				return pptx.NewSlide("Slide 16").AddSmartArt(sa), "ppt/slides/slide1.xml"
			},
			expected: []string{`descr="SmartArt Alt" title="SmartArt Alt"`},
		},
		{
			name: "SmartArt Decorative",
			setup: func() (pptx.SlideContent, string) {
				sa := smartart.NewSmartArt(smartart.BasicBlockList).
					WithAltText("Ignored").
					WithDecorative(true).
					AddNode(smartart.NewNode("Node 1"))
				return pptx.NewSlide("Slide 17").AddSmartArt(sa), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
		{
			name: "SmartArt Empty AltText",
			setup: func() (pptx.SlideContent, string) {
				sa := smartart.NewSmartArt(smartart.BasicBlockList).
					WithAltText("").
					WithDecorative(false).
					AddNode(smartart.NewNode("Node 1"))
				return pptx.NewSlide("Slide 18").AddSmartArt(sa), "ppt/slides/slide1.xml"
			},
			expected:   []string{`descr=""`},
			unexpected: []string{`title="`},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pb := pptx.NewPresentationBuilder("Access Test")
			slide, xmlPath := tc.setup()
			pb.AddSlide(slide)

			pptxBytes, err := pb.Build()
			if err != nil {
				t.Fatalf("Build failed: %v", err)
			}

			zr, err := zip.NewReader(bytes.NewReader(pptxBytes), int64(len(pptxBytes)))
			if err != nil {
				t.Fatalf("Failed to create zip reader: %v", err)
			}

			slideXML := testutil.ReadZipFile(t, zr, xmlPath)

			for _, exp := range tc.expected {
				if !strings.Contains(slideXML, exp) {
					t.Errorf("Expected content %q not found in XML", exp)
				}
			}

			for _, unexp := range tc.unexpected {
				if strings.Contains(slideXML, unexp) {
					t.Errorf("Unexpected content %q found in XML", unexp)
				}
			}
		})
	}
}

func TestAccessibilityValidation(t *testing.T) {
	longText := strings.Repeat("a", 151)
	imgData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	tests := []struct {
		name          string
		setup         func() pptx.SlideContent
		expectedError string
	}{
		{
			name: "Shape AltText Too Long",
			setup: func() pptx.SlideContent {
				shape := pptx.NewRectangle(1, 1, 1, 1).
					WithAltText(longText).
					WithDecorative(false)
				return pptx.NewSlide("Slide 1").AddShape(shape)
			},
			expectedError: "alt text exceeds 150 characters",
		},
		{
			name: "Image AltText Too Long",
			setup: func() pptx.SlideContent {
				img := pptx.NewImageFromBytes(imgData, "png", 0, 0, 100, 100).
					WithAltText(longText).
					WithDecorative(false)
				return pptx.NewSlide("Slide 2").AddImage(img)
			},
			expectedError: "alt text exceeds 150 characters",
		},
		{
			name: "Connector AltText Too Long",
			setup: func() pptx.SlideContent {
				conn := pptx.NewStraightConnector(1, 1, 2, 2).
					WithAltText(longText).
					WithDecorative(false)
				return pptx.NewSlide("Slide 3").AddConnector(conn)
			},
			expectedError: "alt text exceeds 150 characters",
		},
		{
			name: "Table AltText Too Long",
			setup: func() pptx.SlideContent {
				table := pptx.NewTable([]styling.Length{styling.Inches(1)}).
					WithAltText(longText).
					WithDecorative(false)
				table = table.AddRow([]string{"A"})
				return pptx.NewSlide("Slide 4").WithTable(table)
			},
			expectedError: "alt text exceeds 150 characters",
		},
		{
			name: "Chart AltText Too Long",
			setup: func() pptx.SlideContent {
				chart := pptx.NewBarChart([]string{"Cat1"}, []float64{10}).
					WithAltText(longText).
					WithDecorative(false)
				return pptx.NewSlide("Slide 5").WithBarChart(chart)
			},
			expectedError: "alt text exceeds 150 characters",
		},
		{
			name: "SmartArt AltText Too Long",
			setup: func() pptx.SlideContent {
				sa := smartart.NewSmartArt(smartart.BasicBlockList).
					WithAltText(longText).
					WithDecorative(false).
					AddNode(smartart.NewNode("Node 1"))
				return pptx.NewSlide("Slide 6").AddSmartArt(sa)
			},
			expectedError: "alt text exceeds 150 characters",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pb := pptx.NewPresentationBuilder("Validation Test")
			pb.AddSlide(tc.setup())
			_, err := pb.Build()

			if err == nil {
				t.Error("Expected error but got nil")
			} else if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("Expected error containing %q, got %q", tc.expectedError, err.Error())
			}
		})
	}
}
