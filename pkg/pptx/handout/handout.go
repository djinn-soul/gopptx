package handout

import "fmt"

// HandoutLayout specifies how many slides appear per handout page.
//
//nolint:revive // Public type name kept for compatibility with existing consumers.
type HandoutLayout int

const (
	// Layout1Up prints 1 slide per page.
	Layout1Up HandoutLayout = 1
	// Layout2Up prints 2 slides per page.
	Layout2Up HandoutLayout = 2
	// Layout3Up prints 3 slides per page.
	Layout3Up HandoutLayout = 3
	// Layout4Up prints 4 slides per page.
	Layout4Up HandoutLayout = 4
	// Layout6Up prints 6 slides per page.
	Layout6Up HandoutLayout = 6
	// Layout9Up prints 9 slides per page.
	Layout9Up HandoutLayout = 9
	// LayoutOutline prints a text-only outline view.
	LayoutOutline HandoutLayout = 0
)

// SlidesPerPage returns the number of slides per page for this layout.
// Returns 0 for Outline.
func (l HandoutLayout) SlidesPerPage() int {
	return int(l)
}

// String returns a human-readable description.
func (l HandoutLayout) String() string {
	switch l {
	case LayoutOutline:
		return "Outline"
	default:
		return fmt.Sprintf("%dUp", int(l))
	}
}

// HandoutMaster configures the handout master page included in a presentation.
//
//nolint:revive // Public type name kept for compatibility with existing consumers.
type HandoutMaster struct {
	Layout         HandoutLayout
	ShowHeader     bool
	ShowFooter     bool
	ShowDate       bool
	ShowPageNumber bool
	HeaderText     string
	FooterText     string
}

// New returns a HandoutMaster with sensible defaults (1-up, all items visible).
func New() *HandoutMaster {
	return &HandoutMaster{
		Layout:         Layout1Up,
		ShowHeader:     true,
		ShowFooter:     true,
		ShowDate:       true,
		ShowPageNumber: true,
	}
}

// WithLayout sets the slides-per-page layout.
func (h *HandoutMaster) WithLayout(layout HandoutLayout) *HandoutMaster {
	h.Layout = layout
	return h
}

// WithHeader sets the header text and ensures it is visible.
func (h *HandoutMaster) WithHeader(text string) *HandoutMaster {
	h.HeaderText = text
	h.ShowHeader = true
	return h
}

// WithFooter sets the footer text and ensures it is visible.
func (h *HandoutMaster) WithFooter(text string) *HandoutMaster {
	h.FooterText = text
	h.ShowFooter = true
	return h
}

// HideHeader hides the header placeholder.
func (h *HandoutMaster) HideHeader() *HandoutMaster {
	h.ShowHeader = false
	return h
}

// HideFooter hides the footer placeholder.
func (h *HandoutMaster) HideFooter() *HandoutMaster {
	h.ShowFooter = false
	return h
}

// HideDate hides the date placeholder.
func (h *HandoutMaster) HideDate() *HandoutMaster {
	h.ShowDate = false
	return h
}

// HidePageNumber hides the page number placeholder.
func (h *HandoutMaster) HidePageNumber() *HandoutMaster {
	h.ShowPageNumber = false
	return h
}

// boolAttr converts a bool to "1" or "0" for OOXML boolean attributes.
func boolAttr(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// GenerateXML produces the `ppt/handoutMasters/handoutMaster1.xml` content.
func (h *HandoutMaster) GenerateXML() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:handoutMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
  <p:cSld>
    <p:spTree>
      <p:nvGrpSpPr>
        <p:cNvPr id="1" name=""/>
        <p:cNvGrpSpPr/>
        <p:nvPr/>
      </p:nvGrpSpPr>
      <p:grpSpPr>
        <a:xfrm>
          <a:off x="0" y="0"/>
          <a:ext cx="0" cy="0"/>
          <a:chOff x="0" y="0"/>
          <a:chExt cx="0" cy="0"/>
        </a:xfrm>
      </p:grpSpPr>
    </p:spTree>
  </p:cSld>
  <p:clrMap bg1="lt1" tx1="dk1" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
  <p:hf hdr="%s" ftr="%s" dt="%s" sldNum="%s"/>
</p:handoutMaster>`,
		boolAttr(h.ShowHeader),
		boolAttr(h.ShowFooter),
		boolAttr(h.ShowDate),
		boolAttr(h.ShowPageNumber),
	)
}

// RelationshipsXML produces `ppt/handoutMasters/_rels/handoutMaster1.xml.rels`.
func RelationshipsXML(themeIndex int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="../theme/theme%d.xml"/>
</Relationships>`, themeIndex)
}
