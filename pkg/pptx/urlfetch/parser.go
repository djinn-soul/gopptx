package urlfetch

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ContentKind identifies the type of an extracted content block.
type ContentKind int

const (
	KindTitle     ContentKind = iota // h1
	KindHeading                      // h2–h6
	KindParagraph                    // p
	KindListItem                     // li
	KindCode                         // pre / code
	KindImage                        // img
	KindTable                        // table
	KindQuote                        // blockquote
	KindLink                         // standalone a[href]
)

// ContentBlock is a single semantic unit extracted from an HTML document.
type ContentBlock struct {
	Kind         ContentKind
	Text         string
	Level        int // heading level (1=h1, 2=h2 …) or list nesting depth
	TableRows    [][]string
	ImageSrc     string
	ImageAlt     string
	LinkHref     string
	HeadingLevel int
}

// IsHeading reports whether this block is a title or section heading.
func (b ContentBlock) IsHeading() bool {
	return b.Kind == KindTitle || b.Kind == KindHeading
}

// HeadingGroup is a heading paired with the content blocks that follow it.
type HeadingGroup struct {
	Heading  ContentBlock
	Children []ContentBlock
}

// WebContent holds all content extracted from a web page.
type WebContent struct {
	Title       string
	URL         string
	Description string // may be empty
	Blocks      []ContentBlock
	Images      [][2]string // [src, alt] pairs
}

// IsEmpty reports whether no content blocks were found.
func (wc *WebContent) IsEmpty() bool {
	return len(wc.Blocks) == 0
}

// GroupedByHeadings groups blocks into sections led by a heading.
func (wc *WebContent) GroupedByHeadings() []HeadingGroup {
	var groups []HeadingGroup
	var cur *HeadingGroup
	for _, b := range wc.Blocks {
		if b.IsHeading() {
			if cur != nil {
				groups = append(groups, *cur)
			}
			cur = &HeadingGroup{Heading: b}
		} else if cur != nil {
			cur.Children = append(cur.Children, b)
		}
	}
	if cur != nil {
		groups = append(groups, *cur)
	}
	return groups
}

func mainContentSelectors() []string {
	return []string{
		"main article",
		"article",
		"main",
		"[role=main]",
		".content",
		".post-content",
		".article-content",
		".entry-content",
		".markdown-body",
		".prose",
		"#content",
		"#main",
		"#article",
		".article",
		"body",
	}
}

func shouldSkipTag(tag string) bool {
	switch tag {
	case "script", "style", "noscript", "svg",
		"form", "button", "input", "select",
		"textarea", "iframe":
		return true
	default:
		return false
	}
}

func shouldSkipClass(cls string) bool {
	lower := strings.ToLower(cls)
	return strings.Contains(lower, "advertisement") ||
		strings.Contains(lower, "ad-container") ||
		strings.Contains(lower, "social-share") ||
		strings.Contains(lower, "comment-section")
}

func isNoRecurseTag(tag string) bool {
	switch tag {
	case "p", "li", "pre", "code",
		"img", "table", "blockquote",
		"h1", "h2", "h3", "h4", "h5", "h6":
		return true
	default:
		return false
	}
}

const (
	minTextLen     = 10
	maxHeadingLen  = 300
	maxListItemLen = 500
	maxCodeLen     = 1000
	maxTableRows   = 30
	minMainTextLen = 100
)

// WebParser extracts structured content from an HTML document.
type WebParser struct {
	cfg Config
}

// NewWebParser creates a WebParser with default config.
func NewWebParser() *WebParser {
	return &WebParser{cfg: DefaultConfig()}
}

// NewWebParserWithConfig creates a WebParser with custom config.
func NewWebParserWithConfig(cfg Config) *WebParser {
	return &WebParser{cfg: cfg}
}

// Parse extracts structured content from raw HTML, attributing it to url.
func (p *WebParser) Parse(html, pageURL string) (*WebContent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Parse base URL for resolving relative URLs
	var baseURL *url.URL
	if pageURL != "" {
		baseURL, err = url.Parse(pageURL)
		if err != nil {
			return nil, fmt.Errorf("invalid page URL: %w", err)
		}
	}

	wc := &WebContent{URL: pageURL}
	wc.Title = p.extractTitle(doc)
	wc.Description = p.extractMetaDescription(doc)

	main := p.findMainContent(doc)
	if main == nil {
		return nil, ErrNoContent
	}

	p.walkSelection(main, wc, 0, baseURL)

	if wc.IsEmpty() {
		return nil, ErrNoContent
	}
	return wc, nil
}

func (p *WebParser) extractTitle(doc *goquery.Document) string {
	if t := strings.TrimSpace(doc.Find("title").First().Text()); t != "" {
		return t
	}
	if t := strings.TrimSpace(doc.Find("h1").First().Text()); t != "" {
		return t
	}
	if t, ok := doc.Find("meta[property='og:title']").First().Attr("content"); ok {
		if t = strings.TrimSpace(t); t != "" {
			return t
		}
	}
	return "Untitled"
}

func (p *WebParser) extractMetaDescription(doc *goquery.Document) string {
	if d, ok := doc.Find("meta[name='description']").First().Attr("content"); ok {
		if d = strings.TrimSpace(d); d != "" {
			return d
		}
	}
	if d, ok := doc.Find("meta[property='og:description']").First().Attr("content"); ok {
		if d = strings.TrimSpace(d); d != "" {
			return d
		}
	}
	return ""
}

func (p *WebParser) findMainContent(doc *goquery.Document) *goquery.Selection {
	// Use custom selectors if provided, otherwise use defaults
	selectors := p.cfg.ContentSelectors
	if len(selectors) == 0 {
		selectors = mainContentSelectors()
	}

	for _, sel := range selectors {
		found := doc.Find(sel).First()
		if found.Length() == 0 {
			continue
		}

		// Apply exclude selectors
		for _, excludeSel := range p.cfg.ExcludeSelectors {
			found.Find(excludeSel).Remove()
		}

		if len(strings.TrimSpace(found.Text())) >= minMainTextLen {
			return found
		}
	}
	return nil
}

//nolint:gocyclo,cyclop,gocognit,funlen // DOM-walk branching is intentionally linear and explicit.
func (p *WebParser) walkSelection(sel *goquery.Selection, wc *WebContent, depth int, baseURL *url.URL) {
	tag := goquery.NodeName(sel)

	if shouldSkipTag(tag) {
		return
	}

	// Check exclude selectors
	if slices.ContainsFunc(p.cfg.ExcludeSelectors, sel.Is) {
		return
	}

	if cls, exists := sel.Attr("class"); exists {
		if shouldSkipClass(cls) {
			return
		}
	}

	switch tag {
	case "h1":
		if t := cleanText(sel); t != "" && len(t) < maxHeadingLen {
			wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindTitle, Text: t, HeadingLevel: 1})
		}
	case "h2", "h3", "h4", "h5", "h6":
		if t := cleanText(sel); t != "" && len(t) < maxHeadingLen {
			lvl := int(tag[1] - '0')
			wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindHeading, Text: t, HeadingLevel: lvl})
		}
	case "p":
		if t := cleanText(sel); len(t) >= minTextLen {
			wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindParagraph, Text: t})
		}
	case "li":
		if t := cleanText(sel); t != "" && len(t) < maxListItemLen {
			wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindListItem, Text: t, Level: depth})
		}
	case "blockquote":
		if t := cleanText(sel); t != "" {
			wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindQuote, Text: t})
		}
	case "pre", "code":
		if p.cfg.IncludeCode {
			t := strings.TrimSpace(sel.Text())
			if t != "" && len(t) <= maxCodeLen {
				wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindCode, Text: t})
			}
		}
		return // no child recursion into code blocks
	case "img":
		if p.cfg.IncludeImages {
			src, _ := sel.Attr("src")
			alt, _ := sel.Attr("alt")
			alt = strings.TrimSpace(alt)
			if src != "" && !strings.HasPrefix(src, "data:") && alt != "" {
				// Resolve relative URLs
				resolvedSrc := p.resolveURL(src, baseURL)
				wc.Images = append(wc.Images, [2]string{resolvedSrc, alt})
				wc.Blocks = append(wc.Blocks, ContentBlock{
					Kind: KindImage, ImageSrc: resolvedSrc, ImageAlt: alt,
				})
			}
		}
	case "a":
		if p.cfg.ExtractLinks {
			href, _ := sel.Attr("href")
			t := cleanText(sel)
			if href != "" && t != "" {
				wc.Blocks = append(wc.Blocks, ContentBlock{
					Kind: KindLink, Text: t, LinkHref: href,
				})
			}
		}
		return // avoid double-counting link text via child recursion
	case "table":
		if p.cfg.IncludeTables {
			rows := extractTableRows(sel)
			if len(rows) > 0 && len(rows) <= maxTableRows {
				wc.Blocks = append(wc.Blocks, ContentBlock{Kind: KindTable, TableRows: rows})
			}
		}
		return // table handled fully; don't recurse
	}

	if isNoRecurseTag(tag) {
		return
	}

	sel.Children().Each(func(_ int, child *goquery.Selection) {
		p.walkSelection(child, wc, depth+1, baseURL)
	})
}

// resolveURL resolves a potentially relative URL against the base URL.
func (p *WebParser) resolveURL(ref string, baseURL *url.URL) string {
	if baseURL == nil {
		return ref
	}

	// If already absolute, return as-is
	if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
		return ref
	}

	// If protocol-relative, add https
	if strings.HasPrefix(ref, "//") {
		return "https:" + ref
	}

	// Resolve against base URL
	u, err := baseURL.Parse(ref)
	if err != nil {
		return ref
	}

	return u.String()
}

func cleanText(sel *goquery.Selection) string {
	return strings.Join(strings.Fields(sel.Text()), " ")
}

func extractTableRows(sel *goquery.Selection) [][]string {
	var rows [][]string
	sel.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		var cells []string
		tr.Find("th, td").Each(func(_ int, cell *goquery.Selection) {
			cells = append(cells, cleanText(cell))
		})
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	})
	return rows
}
