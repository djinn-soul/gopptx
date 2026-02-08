package pptx

import (
	"archive/zip"
	"bytes"
	"testing"
)

func TestSlideLayoutParityFixturesAgainstPptRsLayoutDemo(t *testing.T) {
	slides := []SlideContent{
		NewSlide("Welcome to Layout Demo").WithTitleOnlyLayout(),
		NewSlide("Centered Title Slide").WithCenteredTitleLayout(),
		NewSlide("Standard Layout").WithTitleAndContentLayout().
			AddBullet("Point 1: Title at top").
			AddBullet("Point 2: Content below").
			AddBullet("Point 3: Most common layout"),
		NewSlide("Big Content Area").WithTitleAndBigContentLayout().
			AddBullet("More space for content").
			AddBullet("Smaller title area").
			AddBullet("Good for detailed slides").
			AddBullet("Maximizes content space"),
		NewSlide("Two Column Layout").WithTwoColumnLayout().
			AddBullet("Left column content").
			AddBullet("Organized side by side").
			AddBullet("Great for comparisons"),
		NewSlide("").WithBlankLayout(),
	}

	data, err := CreateWithSlides("Layout Demo", slides)
	if err != nil {
		t.Fatalf("CreateWithSlides error: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip read error: %v", err)
	}

	cases := []struct {
		referenceSlide string
		generatedSlide string
		tokens         []string
	}{
		{
			referenceSlide: "ppt/slides/slide1.xml",
			generatedSlide: "ppt/slides/slide1.xml",
			tokens: []string{
				`name="Title"`,
				`<a:off x="457200" y="274638"/>`,
				`<a:t>Welcome to Layout Demo</a:t>`,
			},
		},
		{
			referenceSlide: "ppt/slides/slide2.xml",
			generatedSlide: "ppt/slides/slide2.xml",
			tokens: []string{
				`name="Title"`,
				`<a:off x="457200" y="2743200"/>`,
				`<a:pPr algn="ctr"/>`,
			},
		},
		{
			referenceSlide: "ppt/slides/slide3.xml",
			generatedSlide: "ppt/slides/slide3.xml",
			tokens: []string{
				`name="Content"`,
				`<a:off x="457200" y="1600200"/>`,
				`<a:t>Point 1: Title at top</a:t>`,
			},
		},
		{
			referenceSlide: "ppt/slides/slide4.xml",
			generatedSlide: "ppt/slides/slide4.xml",
			tokens: []string{
				`name="Content"`,
				`<a:off x="457200" y="1189200"/>`,
				`<a:ext cx="8230200" cy="5668800"/>`,
			},
		},
		{
			referenceSlide: "ppt/slides/slide5.xml",
			generatedSlide: "ppt/slides/slide5.xml",
			tokens: []string{
				`name="Left Content"`,
				`name="Right Content"`,
				`<a:off x="457200" y="1189200"/>`,
				`<a:off x="4572300" y="1189200"/>`,
			},
		},
		{
			referenceSlide: "ppt/slides/slide6.xml",
			generatedSlide: "ppt/slides/slide6.xml",
			tokens: []string{
				`<p:spTree>`,
			},
		},
	}

	for _, tc := range cases {
		reference := fixtureSlideXML(t, "layout_demo.pptx", tc.referenceSlide)
		ours := readZipFile(t, zr, tc.generatedSlide)
		assertContainsTokens(t, "ppt-rs layout fixture "+tc.referenceSlide, reference, tc.tokens)
		assertContainsTokens(t, "gopptx layout parity "+tc.generatedSlide, ours, tc.tokens)
	}
}


