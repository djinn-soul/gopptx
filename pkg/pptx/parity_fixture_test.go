package pptx_test

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
	"github.com/djinn-soul/gopptx/pkg/pptx/internal/testutil"
)

func TestBasicParityFixtureAgainstPptRsSimpleDeck(t *testing.T) {
	reference := testutil.ReadAllSlidesXML(t, testutil.FixtureZipReader(t, "simple.pptx"))
	ours := generatedAllSlidesXML(t, []pptx.SlideContent{
		pptx.NewSlide("Welcome").
			AddBullet("Hello, World!").
			AddBullet("This is a simple presentation").
			AddBullet("Created with ppt-rs templates"),
		pptx.NewSlide("Features").
			AddBullet("Easy to use API").
			AddBullet("Multiple templates").
			AddBullet("Theme support"),
		pptx.NewSlide("Conclusion").
			AddBullet("Try ppt-rs today!").
			AddBullet("Visit github.com/yingkitw/ppt-rs"),
	})

	tokens := []string{
		`<a:t>Welcome</a:t>`,
		`<a:t>Features</a:t>`,
		`<a:t>Conclusion</a:t>`,
		`<a:t>Hello, World!</a:t>`,
		`<a:t>This is a simple presentation</a:t>`,
		`<a:t>Easy to use API</a:t>`,
		`<a:t>Theme support</a:t>`,
	}
	testutil.AssertContainsTokens(t, "ppt-rs simple fixture", reference, tokens)
	testutil.AssertContainsTokens(t, "gopptx basic parity deck", ours, tokens)
}

func TestLayoutParityFixtureAgainstPptRsDeck(t *testing.T) {
	reference := testutil.ReadAllSlidesXML(t, testutil.FixtureZipReader(t, "layout_demo.pptx"))
	ours := generatedAllSlidesXML(t, []pptx.SlideContent{
		pptx.NewSlide("Welcome to Layout Demo").WithTitleOnlyLayout(),
		pptx.NewSlide("Centered Title Slide").
			WithCenteredTitleLayout().
			WithTitleSize(60).
			WithTitleColor("4F81BD"),
		pptx.NewSlide("Standard Layout").
			AddBullet("Point 1: Title at top").
			AddBullet("Point 2: Content below").
			AddBullet("Point 3: Most common layout"),
		pptx.NewSlide("Big Content Area").
			WithTitleAndBigContentLayout().
			AddBullet("More space for content").
			AddBullet("Smaller title area").
			AddBullet("Good for detailed slides").
			AddBullet("Maximizes content space"),
		pptx.NewSlide("Two Column Layout").
			WithTwoColumnLayout().
			AddBullet("Left column content").
			AddBullet("Organized side by side").
			AddBullet("Great for comparisons"),
		pptx.NewSlide("").WithBlankLayout(),
		pptx.NewSlide("Summary").
			WithTitleSize(48).
			WithTitleBold(true).
			WithTitleColor("C0504D").
			AddBullet("Layout types implemented:").
			AddBullet("• TitleOnly - Just title").
			AddBullet("• CenteredTitle - Title centered").
			AddBullet("• TitleAndContent - Standard").
			AddBullet("• TitleAndBigContent - Large content").
			AddBullet("• TwoColumn - Side by side").
			AddBullet("• Blank - Empty slide").
			WithContentSize(20),
	})

	tokens := []string{
		`<a:t>Welcome to Layout Demo</a:t>`,
		`<a:t>Centered Title Slide</a:t>`,
		`sz="6000"`,
		`<a:srgbClr val="4F81BD"/>`,
		`<a:t>Standard Layout</a:t>`,
		`<a:t>Big Content Area</a:t>`,
		`<a:t>Two Column Layout</a:t>`,
		`<a:t>Summary</a:t>`,
		`sz="4800"`,
		`b="1"`,
		`<a:srgbClr val="C0504D"/>`,
		`sz="2000"`,
	}
	testutil.AssertContainsTokens(t, "ppt-rs layout fixture", reference, tokens)
	testutil.AssertContainsTokens(t, "gopptx layout parity deck", ours, tokens)
}

func TestTextFormattingParityFixtureAgainstPptRsProfessionalDeck(t *testing.T) {
	reference := testutil.ReadAllSlidesXML(t, testutil.FixtureZipReader(t, "professional.pptx"))
	ours := generatedAllSlidesXML(t, []pptx.SlideContent{
		pptx.NewSlide("Text Parity").AddBulletRuns([]pptx.TextRun{
			pptx.NewTextRun("Dark blue title style").WithBold(true).WithColor("003366"),
			pptx.NewTextRun("Orange italic style").WithItalic(true).WithColor("FF6600"),
			pptx.NewTextRun("Underlined emphasis").WithUnderline(true),
		}),
	})

	tokens := []string{
		`b="1"`,
		`i="1"`,
		`u="sng"`,
		`<a:srgbClr val="003366"/>`,
		`<a:srgbClr val="FF6600"/>`,
	}
	testutil.AssertContainsTokens(t, "ppt-rs professional fixture", reference, tokens)
	testutil.AssertContainsTokens(t, "gopptx text-format parity deck", ours, tokens)
}

func TestBulletStylesParityFixtureAgainstPptRsDeck(t *testing.T) {
	reference := testutil.ReadAllSlidesXML(t, testutil.FixtureZipReader(t, "bullet_styles.pptx"))
	ours := generatedAllSlidesXML(t, []pptx.SlideContent{
		pptx.NewSlide("Bullet Styles").
			AddBulletWithStyle("First step", pptx.NewTextParagraphStyle().WithNumbered()).
			AddBulletWithStyle("Option A", pptx.NewTextParagraphStyle().WithLetteredLower()).
			AddBulletWithStyle("Chapter I", pptx.NewTextParagraphStyle().WithRomanUpper()).
			AddBulletWithStyle("Nested", pptx.NewTextParagraphStyle().WithNumbered().WithLevel(1)).
			AddBulletWithStyle("Custom", pptx.NewTextParagraphStyle().WithCustomBullet("~")),
	})

	tokens := []string{
		`<a:buAutoNum type="arabicPeriod"/>`,
		`<a:buAutoNum type="alphaLcPeriod"/>`,
		`<a:buAutoNum type="romanUcPeriod"/>`,
		`<a:pPr lvl="1" marL="1371600" indent="-914400">`,
		`<a:buChar char="`,
	}
	testutil.AssertContainsTokens(t, "ppt-rs bullet fixture", reference, tokens)
	testutil.AssertContainsTokens(t, "gopptx bullet parity deck", ours, tokens)
}

func TestTextEnhancementsParityFixtureAgainstPptRsComprehensiveDemo(t *testing.T) {
	reference := testutil.ReadZipFile(t, testutil.FixtureZipReader(t, "comprehensive_demo.pptx"), "ppt/slides/slide29.xml")
	ours := generatedSlideXML(t,
		pptx.NewSlide("Text Enhancements - New Formatting").
			AddBulletRuns([]pptx.TextRun{pptx.NewTextRun("Strike").WithStrikethrough(true)}).
			AddBulletRuns([]pptx.TextRun{pptx.NewTextRun("Highlight").WithHighlight("FFFF00")}).
			AddBulletRuns([]pptx.TextRun{pptx.NewTextRun("H2O").WithSubscript(true)}).
			AddBulletRuns([]pptx.TextRun{pptx.NewTextRun("x2").WithSuperscript(true)}),
	)

	tokens := []string{
		`strike="sngStrike"`,
		`<a:highlight><a:srgbClr val="FFFF00"/></a:highlight>`,
		`baseline="-25000"`,
		`baseline="30000"`,
		`<a:buChar char="`,
	}
	testutil.AssertContainsTokens(t, "ppt-rs text-enhancement fixture", reference, tokens)
	testutil.AssertContainsTokens(t, "gopptx text-enhancement parity deck", ours, tokens)
}

func TestImageFormatParityCasesFromPptRsExamples(t *testing.T) {
	cases := []struct {
		name string
		ext  string
		mime string
		data []byte
	}{
		{name: "png", ext: "png", mime: "image/png", data: testutil.TinyPNG},
		{name: "jpg", ext: "jpg", mime: "image/jpeg", data: []byte{0xFF, 0xD8, 0xFF, 0xD9}},
		{name: "jpeg", ext: "jpeg", mime: "image/jpeg", data: []byte{0xFF, 0xD8, 0xFF, 0xD9}},
		{name: "gif", ext: "gif", mime: "image/gif", data: []byte("GIF89a")},
		{name: "bmp", ext: "bmp", mime: "image/bmp", data: []byte{'B', 'M', 0x00, 0x00}},
		{name: "tif", ext: "tif", mime: "image/tiff", data: []byte{'I', 'I', '*', 0x00}},
		{name: "tiff", ext: "tiff", mime: "image/tiff", data: []byte{'I', 'I', '*', 0x00}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			imgPath := filepath.Join(tmpDir, "sample."+tc.ext)
			if err := os.WriteFile(imgPath, tc.data, 0o600); err != nil {
				t.Fatalf("write %s image: %v", tc.ext, err)
			}

			data, err := pptx.CreateWithSlides("Demo", []pptx.SlideContent{
				pptx.NewSlide("Image").AddImage(pptx.NewImage(imgPath, 1200000, 1700000, 2400000, 1800000)),
			})
			if err != nil {
				t.Fatalf("CreateWithSlides error for %s: %v", tc.ext, err)
			}

			zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
			if err != nil {
				t.Fatalf("zip read error for %s: %v", tc.ext, err)
			}

			mediaPath := "ppt/media/image1." + tc.ext
			if !testutil.ZipHasFile(zr, mediaPath) {
				t.Fatalf("missing embedded media %s", mediaPath)
			}

			contentTypes := testutil.ReadZipFile(t, zr, "[Content_Types].xml")
			expectedType := fmt.Sprintf(`Extension="%s" ContentType="%s"`, tc.ext, tc.mime)
			if !strings.Contains(contentTypes, expectedType) {
				t.Fatalf("expected %q in content types", expectedType)
			}

			relsXML := testutil.ReadZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
			target := fmt.Sprintf(`Target="../media/image1.%s"`, tc.ext)
			if !strings.Contains(relsXML, target) {
				t.Fatalf("expected %q in slide rels", target)
			}
		})
	}
}

func generatedAllSlidesXML(t *testing.T, slides []pptx.SlideContent) string {
	t.Helper()
	data, err := pptx.CreateWithSlides("Parity", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	return testutil.ReadAllSlidesXML(t, zr)
}

func generatedSlideXML(t *testing.T, slide pptx.SlideContent) string {
	t.Helper()
	data, err := pptx.CreateWithSlides("Parity", []pptx.SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}
	return testutil.ReadZipFile(t, zr, "ppt/slides/slide1.xml")
}
