package urlfetch

import (
	"net/http"
	"strings"
	"time"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	maxParaLen    = 200
	maxListLen    = 180
	maxQuoteLen   = 180
	maxCodeBullet = 150
	utf8ContMask  = 0x80
	// defaultImageWidthEMU is the default width for embedded images (3 inches).
	defaultImageWidthEMU int64 = 2743200
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

// ConvertToSlides transforms parsed web content into slide content objects.
func (c *Converter) ConvertToSlides(content *WebContent, opts *ConversionOptions) ([]elements.SlideContent, error) {
	return c.buildSlides(content, opts)
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
	return slide.AddBullet(truncateText(block.Text, maxParaLen)), bulletCount + 1, imageCount
}

func (c *Converter) appendListItem(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
	imageCount int,
) (elements.SlideContent, int, int) {
	return slide.AddBullet("• " + truncateText(block.Text, maxListLen)), bulletCount + 1, imageCount
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
	return slide.AddBullet("[Link] " + truncateText(linkText, maxListLen)), bulletCount + 1, imageCount
}

// buildTable converts raw HTML table rows to a table with a bold header row.
func buildTable(rows [][]string) tables.Table {
	if len(rows) == 0 {
		return tables.NewTable(nil)
	}

	cols := 0
	for _, r := range rows {
		if len(r) > cols {
			cols = len(r)
		}
	}

	const totalWidthEMU int64 = 8230200
	colW := totalWidthEMU / int64(cols)
	colWidths := make([]styling.Length, cols)
	for i := range colWidths {
		colWidths[i] = styling.Emu(colW)
	}

	tbl := tables.NewTable(colWidths)
	for i, rawRow := range rows {
		norm := make([]string, cols)
		copy(norm, rawRow)

		if i == 0 {
			cells := make([]tables.TableCell, cols)
			for j, text := range norm {
				cells[j] = tables.NewTableCell(text).WithBold(true)
			}
			tbl = tbl.AddStyledRow(cells)
		} else {
			tbl = tbl.AddRow(norm)
		}
	}
	return tbl
}

// truncateText safely shortens text to maxLen chars, breaking at a word boundary.
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	end := maxLen
	for end > 0 && !isRuneStart(text[end]) {
		end--
	}
	if idx := strings.LastIndex(text[:end], " "); idx > maxLen/2 {
		end = idx
	}
	return strings.TrimRight(text[:end], " ") + "..."
}

func isRuneStart(b byte) bool {
	return b&0xC0 != utf8ContMask
}
