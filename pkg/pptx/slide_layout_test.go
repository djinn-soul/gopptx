package pptx

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestCreateWithSlidesUsesTitleOnlyLayout(t *testing.T) {
	slide := NewSlide("Title Only").WithTitleOnlyLayout()

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if !strings.Contains(slideXML, `name="Title"`) {
		t.Fatalf("expected title shape in title-only layout")
	}
	if strings.Contains(slideXML, `name="Content"`) {
		t.Fatalf("did not expect content shape in title-only layout")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout2.xml"`) {
		t.Fatalf("expected title-only slide layout target")
	}
}

func TestCreateWithSlidesUsesBlankLayout(t *testing.T) {
	slide := NewSlide("").WithBlankLayout()

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	if strings.Contains(slideXML, `name="Title"`) {
		t.Fatalf("did not expect title shape in blank layout")
	}
	if strings.Contains(slideXML, `name="Content"`) {
		t.Fatalf("did not expect content shape in blank layout")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout3.xml"`) {
		t.Fatalf("expected blank slide layout target")
	}
}

func TestCreateWithSlidesUsesCenteredTitleLayout(t *testing.T) {
	slide := NewSlide("Centered").WithCenteredTitleLayout()

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`name="Title"`,
		`<a:off x="457200" y="2743200"/>`,
		`<a:pPr algn="ctr"/>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in centered-title layout", needle)
		}
	}
	if strings.Contains(slideXML, `name="Content"`) {
		t.Fatalf("did not expect content shape in centered-title layout")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout4.xml"`) {
		t.Fatalf("expected centered-title slide layout target")
	}
}

func TestCreateWithSlidesUsesTitleAndBigContentLayout(t *testing.T) {
	slide := NewSlide("Big").WithTitleAndBigContentLayout().
		AddBullet("One").
		AddBullet("Two")

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`name="Content"`,
		`<a:off x="457200" y="1189200"/>`,
		`<a:ext cx="8230200" cy="5668800"/>`,
		`<a:t>One</a:t>`,
		`<a:t>Two</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in title-and-big-content layout", needle)
		}
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout5.xml"`) {
		t.Fatalf("expected title-and-big-content slide layout target")
	}
}

func TestCreateWithSlidesUsesTwoColumnLayoutAndSplitsBullets(t *testing.T) {
	slide := NewSlide("Two Col").WithTwoColumnLayout().
		AddBullet("A").
		AddBullet("B").
		AddBullet("C").
		AddBullet("D").
		AddBullet("E")

	data, err := CreateWithSlides("Demo", []SlideContent{slide})
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	slideXML := readZipFile(t, zr, "ppt/slides/slide1.xml")
	checks := []string{
		`name="Left Content"`,
		`name="Right Content"`,
		`<a:off x="457200" y="1189200"/>`,
		`<a:off x="4572300" y="1189200"/>`,
		`<a:t>A</a:t>`,
		`<a:t>C</a:t>`,
		`<a:t>D</a:t>`,
		`<a:t>E</a:t>`,
	}
	for _, needle := range checks {
		if !strings.Contains(slideXML, needle) {
			t.Fatalf("expected %q in two-column layout", needle)
		}
	}

	leftPos := strings.Index(slideXML, `name="Left Content"`)
	rightPos := strings.Index(slideXML, `name="Right Content"`)
	if leftPos == -1 || rightPos == -1 || leftPos > rightPos {
		t.Fatalf("expected left column before right column")
	}
	withinLeft := slideXML[leftPos:rightPos]
	if !strings.Contains(withinLeft, `<a:t>A</a:t>`) || !strings.Contains(withinLeft, `<a:t>C</a:t>`) {
		t.Fatalf("expected first half bullets in left column")
	}
	if strings.Contains(withinLeft, `<a:t>D</a:t>`) {
		t.Fatalf("did not expect right-half bullets in left column")
	}

	relsXML := readZipFile(t, zr, "ppt/slides/_rels/slide1.xml.rels")
	if !strings.Contains(relsXML, `Target="../slideLayouts/slideLayout6.xml"`) {
		t.Fatalf("expected two-column slide layout target")
	}
}

func TestCreateWithSlidesEmbedsAllLayoutFiles(t *testing.T) {
	slides := []SlideContent{
		NewSlide("A"),
		NewSlide("B").WithTitleOnlyLayout(),
		NewSlide("").WithBlankLayout(),
		NewSlide("C").WithCenteredTitleLayout(),
		NewSlide("D").WithTitleAndBigContentLayout(),
		NewSlide("E").WithTwoColumnLayout(),
	}

	data, err := CreateWithSlides("Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	required := []string{
		"ppt/slideLayouts/slideLayout1.xml",
		"ppt/slideLayouts/slideLayout2.xml",
		"ppt/slideLayouts/slideLayout3.xml",
		"ppt/slideLayouts/slideLayout4.xml",
		"ppt/slideLayouts/slideLayout5.xml",
		"ppt/slideLayouts/slideLayout6.xml",
		"ppt/slideLayouts/_rels/slideLayout1.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout2.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout3.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout4.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout5.xml.rels",
		"ppt/slideLayouts/_rels/slideLayout6.xml.rels",
	}
	for _, name := range required {
		if !zipHasFile(zr, name) {
			t.Fatalf("missing %s", name)
		}
	}
}

func TestCreateWithSlidesRejectsTitleOnlyWithBullets(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("No Bullets").WithTitleOnlyLayout().AddBullet("x"),
	})
	if err == nil {
		t.Fatalf("expected title-only bullets validation error")
	}
	if !strings.Contains(err.Error(), "title_only layout does not support bullets") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsBlankLayoutWithTitle(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Should be Empty").WithBlankLayout(),
	})
	if err == nil {
		t.Fatalf("expected blank-layout title validation error")
	}
	if !strings.Contains(err.Error(), "blank layout requires empty title") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsInvalidSlideLayout(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Bad").WithLayout("hero"),
	})
	if err == nil {
		t.Fatalf("expected invalid layout validation error")
	}
	if !strings.Contains(err.Error(), "layout must be one of title_and_content|title_only|blank|centered_title|title_and_big_content|two_column") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWithSlidesRejectsCenteredTitleWithBullets(t *testing.T) {
	_, err := CreateWithSlides("Demo", []SlideContent{
		NewSlide("Centered").WithCenteredTitleLayout().AddBullet("x"),
	})
	if err == nil {
		t.Fatalf("expected centered-title bullets validation error")
	}
	if !strings.Contains(err.Error(), "centered_title layout does not support bullets") {
		t.Fatalf("unexpected error: %v", err)
	}
}
