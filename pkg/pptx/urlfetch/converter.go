package urlfetch

import (
	"strings"

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

// truncateText preserves full extracted text for slide conversion.
// maxLen is intentionally ignored to avoid silent content loss in URL fetch decks.
func truncateText(text string, _ int) string {
	return strings.TrimSpace(text)
}

func splitTextIntoChunks(text string, maxLen int) []string {
	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return nil
	}
	if len(normalized) <= maxLen {
		return []string{normalized}
	}

	words := strings.Fields(normalized)
	if len(words) == 0 {
		return nil
	}

	chunks := make([]string, 0, len(words)/8+1)
	var current strings.Builder
	for _, word := range words {
		if current.Len() == 0 {
			current.WriteString(word)
			continue
		}
		if current.Len()+1+len(word) <= maxLen {
			current.WriteByte(' ')
			current.WriteString(word)
			continue
		}
		chunks = append(chunks, current.String())
		current.Reset()
		current.WriteString(word)
	}
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}
	return chunks
}
