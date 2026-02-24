package urlfetch

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/tables"
)

const (
	maxParaLen    = 200
	maxListLen    = 180
	maxQuoteLen   = 180
	maxCodeBullet = 150
)

// URLFetchConverter converts parsed web content into a PPTX byte slice.
type URLFetchConverter struct { //nolint:revive // keeping exported name for API compatibility
	cfg URLFetchConfig
}

// NewURLFetchConverter creates a converter with the default config.
func NewURLFetchConverter() *URLFetchConverter {
	return &URLFetchConverter{cfg: DefaultConfig()}
}

// NewURLFetchConverterWithConfig creates a converter with a custom config.
func NewURLFetchConverterWithConfig(cfg URLFetchConfig) *URLFetchConverter {
	return &URLFetchConverter{cfg: cfg}
}

// Web2Ppt is a compatibility alias for URLFetchConverter.
type Web2Ppt = URLFetchConverter

// NewWeb2Ppt creates a converter with the default config.
func NewWeb2Ppt() *URLFetchConverter { return NewURLFetchConverter() }

// NewWeb2PptWithConfig creates a converter with a custom config.
func NewWeb2PptWithConfig(cfg URLFetchConfig) *URLFetchConverter {
	return NewURLFetchConverterWithConfig(cfg)
}

// Convert transforms parsed web content into PPTX bytes.
func (c *URLFetchConverter) Convert(content *WebContent, opts *ConversionOptions) ([]byte, error) {
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
func (c *URLFetchConverter) buildSlides(content *WebContent, opts *ConversionOptions) ([]elements.SlideContent, error) {
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

func (c *URLFetchConverter) buildGroupedSlides(
	content *WebContent,
	slides []elements.SlideContent,
) ([]elements.SlideContent, error) {
	groups := content.GroupedByHeadings()
	if len(groups) == 0 {
		return c.buildLinearSlides(content, slides)
	}

	for _, group := range groups {
		if len(slides) >= c.cfg.MaxSlides {
			break
		}
		slide := elements.NewSlide(group.Heading.Text).WithTitleAndContentLayout()
		bulletCount := 0

		for _, block := range group.Children {
			if len(slides) >= c.cfg.MaxSlides {
				break
			}
			if bulletCount >= c.cfg.MaxBulletsPerSlide {
				if bulletCount > 0 {
					slides = append(slides, slide)
				}
				slide = elements.NewSlide(group.Heading.Text + " (cont.)").WithTitleAndContentLayout()
				bulletCount = 0
			}
			slide, bulletCount = c.appendBlock(slide, block, bulletCount)
		}

		if bulletCount > 0 {
			slides = append(slides, slide)
		}
	}
	return slides, nil
}

func (c *URLFetchConverter) buildLinearSlides(
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

	var cur *elements.SlideContent
	bulletCount := 0

	for _, block := range content.Blocks {
		if len(slides) >= c.cfg.MaxSlides {
			break
		}

		if block.IsHeading() {
			if cur != nil && bulletCount > 0 {
				slides = append(slides, *cur)
			}
			s := elements.NewSlide(block.Text).WithTitleAndContentLayout()
			cur = &s
			bulletCount = 0
			continue
		}

		if cur == nil {
			s := elements.NewSlide("Overview").WithTitleAndContentLayout()
			cur = &s
		}

		if bulletCount >= c.cfg.MaxBulletsPerSlide {
			slides = append(slides, *cur)
			cont := elements.NewSlide(cur.Title + " (cont.)").WithTitleAndContentLayout()
			cur = &cont
			bulletCount = 0
		}

		*cur, bulletCount = c.appendBlock(*cur, block, bulletCount)
	}

	if cur != nil && bulletCount > 0 {
		slides = append(slides, *cur)
	}
	return slides, nil
}

func (c *URLFetchConverter) appendBlock(
	slide elements.SlideContent,
	block ContentBlock,
	bulletCount int,
) (elements.SlideContent, int) {
	switch block.Kind {
	case KindTitle, KindHeading:
		// Headings are handled by slide grouping/creation logic.
	case KindParagraph:
		slide = slide.AddBullet(truncateText(block.Text, maxParaLen))
		bulletCount++
	case KindListItem:
		slide = slide.AddBullet("• " + truncateText(block.Text, maxListLen))
		bulletCount++
	case KindQuote:
		slide = slide.AddBullet(`"` + truncateText(block.Text, maxQuoteLen) + `"`)
		bulletCount++
	case KindCode:
		if c.cfg.IncludeCode {
			slide = slide.AddBullet("[Code] " + truncateText(block.Text, maxCodeBullet))
			bulletCount++
		}
	case KindTable:
		if c.cfg.IncludeTables && len(block.TableRows) > 0 {
			slide = slide.WithTable(buildTable(block.TableRows))
			bulletCount++
		}
	case KindImage:
		if c.cfg.IncludeImages && block.ImageAlt != "" {
			slide = slide.AddBullet("[Image: " + block.ImageAlt + "]")
			bulletCount++
		}
	case KindLink:
		if c.cfg.ExtractLinks {
			linkText := block.Text
			if block.LinkHref != "" && block.LinkHref != block.Text {
				linkText = linkText + " (" + block.LinkHref + ")"
			}
			slide = slide.AddBullet("[Link] " + truncateText(linkText, maxListLen))
			bulletCount++
		}
	}
	return slide, bulletCount
}

// buildTable converts raw HTML table rows to a tables.Table.
// Row 0 is rendered as a bold header row.
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
	return b&0xC0 != 0x80
}
