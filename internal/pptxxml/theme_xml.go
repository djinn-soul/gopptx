package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

// LayoutsPerMaster is how many slide layouts each master ships.
const LayoutsPerMaster = 8

const layoutsPerMaster = LayoutsPerMaster

// ST_SlideLayoutType values. PowerPoint classifies a layout in the New Slide
// gallery by this attribute and matches placeholders through it; a layout
// without one is unclassified whatever its name says.
const (
	LayoutTypeTitle          = "title"
	LayoutTypeObject         = "obj"
	LayoutTypeTwoObject      = "twoObj"
	LayoutTypeSectionHeader  = "secHead"
	LayoutTypeTitleOnly      = "titleOnly"
	LayoutTypeBlank          = "blank"
	LayoutTypeVerticalText   = "vertTx"
	LayoutTypeVerticalTitle  = "vertTitleAndTx"
	LayoutTypeObjectAndText  = "objTx"
	LayoutTypePictureAndText = "picTx"
	LayoutTypeCustom         = "cust"
)

// ColorSchemeSpec defines the 12 colors in a theme.
type ColorSchemeSpec struct {
	Name     string
	Dk1      string
	Lt1      string
	Dk2      string
	Lt2      string
	Accent1  string
	Accent2  string
	Accent3  string
	Accent4  string
	Accent5  string
	Accent6  string
	Hlink    string
	FolHlink string
}

// FontSchemeSpec defines heading and body fonts.
type FontSchemeSpec struct {
	Name      string
	MajorFont string
	MinorFont string
}

// ThemeSpec defines a theme's color and font elements.
type ThemeSpec struct {
	Name   string
	Colors ColorSchemeSpec
	Fonts  FontSchemeSpec
}

// SlideMasterSpec defines the appearance of the slide master.
type SlideMasterSpec struct {
	MasterIndex  int
	Background   *SlideBackgroundSpec
	FooterText   string
	Shapes       []ShapeSpec
	Images       []ImageRef
	ColorMapping *ColorMappingSpec
	TxStyles     *TxStylesSpec
}

// TxStylesSpec defines the default text styles for a slide master.
// Each field holds up to 9 levels of text styling (Lvl1–Lvl9).
type TxStylesSpec struct {
	TitleStyle []TextLevelStyle
	BodyStyle  []TextLevelStyle
	OtherStyle []TextLevelStyle
}

// TextLevelStyle defines default text properties for one indent level.
type TextLevelStyle struct {
	Level      int    // 0-based (0=Lvl1, 8=Lvl9)
	Font       string // Typeface override
	SizePt     int    // Size in points
	Bold       bool
	Italic     bool
	Color      string // 6-digit hex RGB
	BulletChar string // Bullet character override
	IndentEMU  int64  // Left indent in EMU
}

// ColorMappingSpec describes how theme colors map to functional roles on slides.
type ColorMappingSpec struct {
	BG1 string
	TX1 string
}

// SlideLayout renders ppt/slideLayouts/slideLayout1.xml.
func SlideLayout() string {
	return SlideLayoutTitleAndContent()
}

// SlideLayoutTitleAndContent renders a title-and-content layout.
func SlideLayoutTitleAndContent() string {
	return slideLayout("Title and Content", LayoutTypeObject)
}

// SlideLayoutTitleOnly renders a title-only layout.
func SlideLayoutTitleOnly() string {
	return slideLayout("Title Only", LayoutTypeTitleOnly)
}

// SlideLayoutBlank renders a blank layout.
func SlideLayoutBlank() string {
	return slideLayout("Blank", LayoutTypeBlank)
}

// SlideLayoutCenteredTitle renders a centered-title layout.
func SlideLayoutCenteredTitle() string {
	return slideLayout("Centered Title", LayoutTypeTitle)
}

// SlideLayoutTitleAndBigContent renders a title-and-big-content layout.
func SlideLayoutTitleAndBigContent() string {
	return slideLayout("Title and Big Content", LayoutTypeObject)
}

// SlideLayoutTwoColumn renders a two-column layout.
func SlideLayoutTwoColumn() string {
	return slideLayout("Two Column", LayoutTypeTwoObject)
}

// SlideLayoutTitleAndVerticalText renders a title with a vertically written
// body.
func SlideLayoutTitleAndVerticalText() string {
	return slideLayout("Title and Vertical Text", LayoutTypeVerticalText)
}

// SlideLayoutVerticalTitleAndText renders a vertically written title and body.
func SlideLayoutVerticalTitleAndText() string {
	return slideLayout("Vertical Title and Text", LayoutTypeVerticalTitle)
}

func slideLayout(name, layoutType string) string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldLayout xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" type="` +
		Escape(layoutType) + `" preserve="1">
<p:cSld name="` + Escape(name) + `">
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
<p:clrMapOvr>
<a:masterClrMapping/>
</p:clrMapOvr>
</p:sldLayout>`
}

// slideLayoutIDListXML names every layout of one master, in relationship order.
func slideLayoutIDListXML(layoutIDBase int64) string {
	var b strings.Builder
	for i := range LayoutsPerMaster {
		b.WriteString("\n<p:sldLayoutId id=\"")
		b.WriteString(strconv.FormatInt(layoutIDBase+int64(i), 10))
		b.WriteString(`" r:id="rId`)
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(`"/>`)
	}
	return b.String()
}

func slideLayoutPartName(layoutIndex, masterIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	globalLayoutIndex := (masterIndex-1)*layoutsPerMaster + layoutIndex
	return fmt.Sprintf("slideLayout%d.xml", globalLayoutIndex)
}

// SlideLayoutRelationships renders ppt/slideLayouts/_rels/slideLayoutN.xml.rels.
func SlideLayoutRelationships(masterIndex int) string {
	if masterIndex < 1 {
		masterIndex = 1
	}
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" ` +
		`Target="../slideMasters/slideMaster` + strconv.Itoa(masterIndex) + `.xml"/>
</Relationships>`
}

// SlideMaster renders ppt/slideMasters/slideMaster1.xml.
//
//nolint:mnd // Contains specific ID base calculations required by the spec.
func SlideMaster(spec *SlideMasterSpec) string {
	bgXML := slideDefaultBackground
	if spec != nil && spec.Background != nil {
		bgXML = backgroundXML(spec.Background)
	}
	masterIndex := 1
	if spec != nil && spec.MasterIndex > 0 {
		masterIndex = spec.MasterIndex
	}
	masterAttrs := ``
	if masterIndex > 1 {
		// PowerPoint-authored packages mark additional master families as preserved.
		masterAttrs = ` preserve="1"`
	}
	// Allocate in blocks to avoid overlap with master IDs across multiple
	// masters: one master id plus its layouts.
	layoutIDBase := int64(2147483649 + (masterIndex-1)*(LayoutsPerMaster+1))

	footerXML := ""
	if spec != nil && spec.FooterText != "" {
		footerXML = fmt.Sprintf(`
<p:sp>
<p:nvSpPr>
<p:cNvPr id="10" name="Footer Placeholder"/>
<p:cNvSpPr>
<a:spLocks noGrp="1"/>
</p:cNvSpPr>
<p:nvPr>
<p:ph type="ftr" sz="quarter" idx="10"/>
</p:nvPr>
</p:nvSpPr>
<p:spPr>
<a:xfrm>
<a:off x="0" y="6350000"/>
<a:ext cx="9144000" cy="508000"/>
</a:xfrm>
</p:spPr>
<p:txBody>
<a:bodyPr/>
<a:lstStyle/>
<a:p>
<a:r>
<a:rPr lang="en-US" smtClean="0"/>
<a:t>%s</a:t>
</a:r>
</a:p>
</p:txBody>
</p:sp>`, Escape(spec.FooterText))
	}

	masterElementsXML := masterShapesAndImagesXML(spec)

	bg1, tx1 := "lt1", "dk1"
	if spec != nil && spec.ColorMapping != nil {
		if spec.ColorMapping.BG1 != "" {
			bg1 = spec.ColorMapping.BG1
		}
		if spec.ColorMapping.TX1 != "" {
			tx1 = spec.ColorMapping.TX1
		}
	}

	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:sldMaster xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main"` + masterAttrs + `>
<p:cSld>
` + bgXML + `
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
` + footerXML + masterElementsXML + `
</p:spTree>
</p:cSld>
<p:clrMap bg1="` + bg1 + `" tx1="` + tx1 + `" bg2="lt2" tx2="dk2" accent1="accent1" accent2="accent2" ` +
		`accent3="accent3" accent4="accent4" accent5="accent5" accent6="accent6" hlink="hlink" folHlink="folHlink"/>
<p:sldLayoutIdLst>` + slideLayoutIDListXML(layoutIDBase) + `
</p:sldLayoutIdLst>` + txStylesXML(spec) + `
</p:sldMaster>`
}

// masterShapesAndImagesXML renders shapes and images within the master spTree.
func masterShapesAndImagesXML(spec *SlideMasterSpec) string {
	if spec == nil {
		return ""
	}
	var b strings.Builder
	// Shape IDs start at 20 to avoid footer placeholder (id=10)
	nextID := 20
	for _, shape := range spec.Shapes {
		b.WriteString(customShapeXML(shape, nextID))
		nextID++
	}
	for _, img := range spec.Images {
		b.WriteString(imageShape(img, nextID))
		nextID++
	}
	return b.String()
}
