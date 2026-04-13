package urlfetch

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Converter converts parsed web content into a PPTX byte slice.
type Converter struct {
	cfg Config
}

// NewURLFetchConverter creates a converter with the default config.
func NewURLFetchConverter() *Converter {
	return &Converter{cfg: DefaultConfig()}
}

// NewURLFetchConverterWithConfig creates a converter with a custom config.
func NewURLFetchConverterWithConfig(cfg Config) *Converter {
	return &Converter{cfg: cfg}
}

// Web2Ppt is a compatibility alias for Converter.
type Web2Ppt = Converter

// NewWeb2Ppt creates a converter with the default config.
func NewWeb2Ppt() *Converter { return NewURLFetchConverter() }

// NewWeb2PptWithConfig creates a converter with a custom config.
func NewWeb2PptWithConfig(cfg Config) *Converter {
	return NewURLFetchConverterWithConfig(cfg)
}

// Convert transforms parsed web content into PPTX bytes.
func (c *Converter) Convert(content *WebContent, opts *ConversionOptions) ([]byte, error) {
	slides, err := c.buildSlides(content, opts)
	if err != nil {
		return nil, err
	}

	title := content.Title
	if opts != nil && opts.Title != nil {
		title = *opts.Title
	}

	if opts != nil && opts.AddPageNumbers {
		for i := range slides {
			slides[i] = slides[i].WithSlideNumber(true)
		}
	}

	creator := ""
	if opts != nil && opts.Author != nil {
		creator = *opts.Author
	}

	return presentationCreateWithMetadata(title, creator, slides)
}

// buildSlides constructs the slide list from extracted web content.
func (c *Converter) buildSlides(content *WebContent, opts *ConversionOptions) ([]elements.SlideContent, error) {
	var slides []elements.SlideContent

	titleText := content.Title
	if opts != nil && opts.Title != nil {
		titleText = *opts.Title
	}

	titleSlide := elements.NewSlide(titleText).WithCenteredTitleLayout()
	if content.Description != "" {
		titleSlide = titleSlide.AddBullet(content.Description)
	}
	if opts != nil && opts.IncludeSourceURL && content.URL != "" {
		titleSlide = titleSlide.AddBullet("Source: " + content.URL)
	}
	slides = append(slides, titleSlide)

	var err error
	if c.cfg.GroupByHeadings {
		slides, err = c.buildGroupedSlides(content, slides)
	} else {
		slides, err = c.buildLinearSlides(content, slides)
	}
	if err != nil {
		return nil, err
	}

	if len(slides) > c.cfg.MaxSlides {
		slides = slides[:c.cfg.MaxSlides]
	}
	return slides, nil
}
