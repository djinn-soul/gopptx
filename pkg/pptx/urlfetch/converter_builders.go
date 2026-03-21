package urlfetch

import (
	"net/http"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// ConvertToSlides transforms parsed web content into slide content objects.
func (c *Converter) ConvertToSlides(content *WebContent, opts *ConversionOptions) ([]elements.SlideContent, error) {
	return c.buildSlides(content, opts)
}

func (c *Converter) buildGroupedSlides(
	content *WebContent,
	slides []elements.SlideContent,
) ([]elements.SlideContent, error) {
	groups := content.GroupedByHeadings()
	if len(groups) == 0 {
		return c.buildLinearSlides(content, slides)
	}

	var imageFetcher *ImageFetcher
	if c.cfg.DownloadImages {
		client := &http.Client{
			Timeout: time.Duration(c.cfg.TimeoutSecs) * time.Second,
		}
		imageFetcher = NewImageFetcher(client, c.cfg, content.URL)
	}

	for _, group := range groups {
		if len(slides) >= c.cfg.MaxSlides {
			break
		}
		slide := elements.NewSlide(group.Heading.Text).WithTitleAndContentLayout()
		bulletCount := 0
		imageCount := 0
		maxImages := c.cfg.MaxImagesPerSlide
		if maxImages <= 0 {
			maxImages = defaultMaxImagesPerSlide
		}

		for _, block := range group.Children {
			slide, bulletCount, imageCount, slides = c.appendGroupedBlock(
				group.Heading.Text,
				slide,
				block,
				bulletCount,
				imageCount,
				maxImages,
				imageFetcher,
				slides,
			)
			if len(slides) >= c.cfg.MaxSlides {
				break
			}
		}

		if bulletCount > 0 || imageCount > 0 {
			slides = append(slides, slide)
		}
	}
	return slides, nil
}

func (c *Converter) buildLinearSlides(
	content *WebContent,
	slides []elements.SlideContent,
) ([]elements.SlideContent, error) {
	if len(content.Blocks) == 0 {
		if content.Description != "" {
			s := elements.NewSlide("Content").WithTitleAndContentLayout().AddBullet(content.Description)
			slides = append(slides, s)
		}
		return slides, nil
	}

	var imageFetcher *ImageFetcher
	if c.cfg.DownloadImages {
		client := &http.Client{
			Timeout: time.Duration(c.cfg.TimeoutSecs) * time.Second,
		}
		imageFetcher = NewImageFetcher(client, c.cfg, content.URL)
	}

	var cur *elements.SlideContent
	bulletCount := 0
	imageCount := 0
	maxImages := c.cfg.MaxImagesPerSlide
	if maxImages <= 0 {
		maxImages = defaultMaxImagesPerSlide
	}

	for _, block := range content.Blocks {
		if len(slides) >= c.cfg.MaxSlides {
			break
		}

		var continueLoop bool
		cur, slides, bulletCount, imageCount, continueLoop = c.handleLinearHeading(
			block,
			cur,
			slides,
			bulletCount,
			imageCount,
		)
		if continueLoop {
			continue
		}
		cur, bulletCount, imageCount, slides = c.ensureLinearSlideCapacity(
			cur,
			bulletCount,
			imageCount,
			maxImages,
			slides,
		)
		*cur, bulletCount, imageCount = c.appendBlock(
			*cur,
			block,
			bulletCount,
			imageCount,
			imageFetcher,
		)
	}

	if cur != nil && (bulletCount > 0 || imageCount > 0) {
		slides = append(slides, *cur)
	}
	return slides, nil
}
