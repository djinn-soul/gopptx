package urlfetch

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const (
	maxParaLen    = 200
	maxListLen    = 180
	maxQuoteLen   = 180
	maxCodeBullet = 150
	// defaultImageWidthEMU is the default width for embedded images (3 inches).
	defaultImageWidthEMU int64 = 2743200
)

func (c *Converter) appendBlock(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
	imageFetcher *ImageFetcher,
) (elements.SlideContent, int, int) {
	maxImages := c.cfg.MaxImagesPerSlide
	if maxImages <= 0 {
		maxImages = defaultMaxImagesPerSlide
	}

	switch block.Kind {
	case KindTitle, KindHeading:
		return slide, bulletCount, imageCount
	case KindParagraph:
		return c.appendParagraph(slide, block, bulletCount, imageCount)
	case KindListItem:
		return c.appendListItem(slide, block, bulletCount, imageCount)
	case KindQuote:
		return c.appendQuote(slide, block, bulletCount, imageCount)
	case KindCode:
		return c.appendCode(slide, block, bulletCount, imageCount)
	case KindTable:
		return c.appendTable(slide, block, bulletCount, imageCount)
	case KindImage:
		return c.appendImage(slide, block, bulletCount, imageCount, maxImages, imageFetcher)
	case KindLink:
		return c.appendLink(slide, block, bulletCount, imageCount)
	}
	return slide, bulletCount, imageCount
}

// fetchAndCreateImage downloads an image and creates a shapes.Image.
func (c *Converter) fetchAndCreateImage(fetcher *ImageFetcher, src, _ string) (shapes.Image, error) {
	fetched, err := fetcher.FetchImage(src)
	if err != nil {
		return shapes.Image{}, err
	}

	widthEMU, heightEMU := CalculateImageDimensions(
		fetched.Width,
		fetched.Height,
		defaultImageWidthEMU,
	)

	img := shapes.NewImageFromBytes(
		fetched.Data,
		fetched.Format,
		styling.Emu(0),
		styling.Emu(0),
		styling.Emu(widthEMU),
		styling.Emu(heightEMU),
	)

	return img, nil
}

func (c *Converter) appendGroupedBlock(
	heading string,
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
	maxImages int,
	imageFetcher *ImageFetcher,
	slides []elements.SlideContent,
) (elements.SlideContent, int, int, []elements.SlideContent) {
	if bulletCount >= c.cfg.MaxBulletsPerSlide || imageCount >= maxImages {
		if bulletCount > 0 || imageCount > 0 {
			slides = append(slides, slide)
		}
		slide = elements.NewSlide(heading + " (cont.)").WithTitleAndContentLayout()
		bulletCount = 0
		imageCount = 0
	}
	slide, bulletCount, imageCount = c.appendBlock(slide, block, bulletCount, imageCount, imageFetcher)
	return slide, bulletCount, imageCount, slides
}

func (c *Converter) handleLinearHeading(
	block ContentBlock,
	cur *elements.SlideContent,
	slides []elements.SlideContent,
	bulletCount int,
	imageCount int,
) (*elements.SlideContent, []elements.SlideContent, int, int, bool) {
	if !block.IsHeading() {
		return cur, slides, bulletCount, imageCount, false
	}
	if cur != nil && (bulletCount > 0 || imageCount > 0) {
		slides = append(slides, *cur)
	}
	s := elements.NewSlide(block.Text).WithTitleAndContentLayout()
	return &s, slides, 0, 0, true
}

func (c *Converter) ensureLinearSlideCapacity(
	cur *elements.SlideContent,
	bulletCount int,
	imageCount int,
	maxImages int,
	slides []elements.SlideContent,
) (*elements.SlideContent, int, int, []elements.SlideContent) {
	if cur == nil {
		s := elements.NewSlide("Overview").WithTitleAndContentLayout()
		cur = &s
	}
	if bulletCount < c.cfg.MaxBulletsPerSlide && imageCount < maxImages {
		return cur, bulletCount, imageCount, slides
	}
	slides = append(slides, *cur)
	cont := elements.NewSlide(cur.Title + " (cont.)").WithTitleAndContentLayout()
	return &cont, 0, 0, slides
}

func (c *Converter) appendParagraph(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	added := 0
	for _, chunk := range splitTextIntoChunks(block.Text, maxParaLen) {
		slide = slide.AddBullet(chunk)
		added++
	}
	if added == 0 {
		return slide, bulletCount, imageCount
	}
	return slide, bulletCount + added, imageCount
}

func (c *Converter) appendListItem(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	added := 0
	for _, chunk := range splitTextIntoChunks(block.Text, maxListLen) {
		slide = slide.AddBullet("• " + chunk)
		added++
	}
	if added == 0 {
		return slide, bulletCount, imageCount
	}
	return slide, bulletCount + added, imageCount
}

func (c *Converter) appendQuote(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	return slide.AddBullet(`"` + truncateText(block.Text, maxQuoteLen) + `"`), bulletCount + 1, imageCount
}

func (c *Converter) appendCode(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	if !c.cfg.IncludeCode {
		return slide, bulletCount, imageCount
	}
	return slide.AddBullet("[Code] " + truncateText(block.Text, maxCodeBullet)), bulletCount + 1, imageCount
}

func (c *Converter) appendTable(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	if !c.cfg.IncludeTables || len(block.TableRows) == 0 {
		return slide, bulletCount, imageCount
	}
	return slide.WithTable(buildTable(block.TableRows)), bulletCount + 1, imageCount
}

func (c *Converter) appendImage(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
	maxImages int,
	imageFetcher *ImageFetcher,
) (elements.SlideContent, int, int) {
	if imageFetcher != nil && block.ImageSrc != "" && imageCount < maxImages {
		img, err := c.fetchAndCreateImage(imageFetcher, block.ImageSrc, block.ImageAlt)
		if err == nil {
			return slide.AddImage(img), bulletCount, imageCount + 1
		}
	}
	if c.cfg.IncludeImages && block.ImageAlt != "" {
		return slide.AddBullet("[Image: " + block.ImageAlt + "]"), bulletCount + 1, imageCount
	}
	return slide, bulletCount, imageCount
}

func (c *Converter) appendLink(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	if !c.cfg.ExtractLinks {
		return slide, bulletCount, imageCount
	}
	linkText := block.Text
	if block.LinkHref != "" && block.LinkHref != block.Text {
		linkText = linkText + " (" + block.LinkHref + ")"
	}
	added := 0
	for _, chunk := range splitTextIntoChunks(linkText, maxListLen) {
		slide = slide.AddBullet("[Link] " + chunk)
		added++
	}
	if added == 0 {
		return slide, bulletCount, imageCount
	}
	return slide, bulletCount + added, imageCount
}
