package urlfetch

import "strings"

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
	minMainTextLen = 100
)
