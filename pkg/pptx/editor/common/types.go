package editorcommon

import (
	"fmt"
	"strings"
)

// Shared constants for editor logic.
const (
	PresentationRelPath = "ppt/_rels/presentation.xml.rels"
	PresentationXMLPath = "ppt/presentation.xml"
	ContentTypesPath    = "[Content_Types].xml"
	CorePropsPath       = "docProps/core.xml"

	// DCNamespace and related metadata XML namespaces.
	DCNamespace       = "http://purl.org/dc/elements/1.1/"
	DCTermsNamespace  = "http://purl.org/dc/terms/"
	DCMITypeNamespace = "http://purl.org/dc/dcmitype/"
	CPNamespace       = "http://schemas.openxmlformats.org/package/2006/metadata/core-properties"
	XSINamespace      = "http://www.w3.org/2001/XMLSchema-instance"

	RelTypeSlide          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	RelTypeSlideMaster    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
	RelTypeSlideLayout    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	RelTypeNotesSlide     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide"
	RelTypeNotesMaster    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster"
	RelTypeHyperlink      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	RelTypeImage          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	RelTypeChart          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
	RelTypeAudio          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/audio"
	RelTypeVideo          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/video"
	RelTypeMedia          = "http://schemas.microsoft.com/office/2007/relationships/media"
	RelTypeTheme          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	RelTypeSectionList    = "http://schemas.microsoft.com/office/2007/relationships/sectionList"
	RelTypePackage        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package"
	RelTypeCustomXML      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml"
	RelTypeCustomXMLProps = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps"

	RelationshipsXMLNS = "http://schemas.openxmlformats.org/package/2006/relationships"
	ContentTypesXMLNS  = "http://schemas.openxmlformats.org/package/2006/content-types"
	SlideContentType   = "application/vnd.openxmlformats-officedocument.presentationml.slide+xml"
)

// EditorRelationship describes an OOXML relationship entry.
type EditorRelationship struct {
	ID         string
	Type       string
	Target     string
	TargetMode string
}

// EditorSlideRef describes internal slide tracking data within the editor.
type EditorSlideRef struct {
	SlideID int64
	RelID   string
	Target  string
	Part    string
	Title   string
}

// xmlEscaper provides fast, basic XML escaping.
var xmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&quot;",
	"'", "&apos;",
)

// XMLEscape provides basic XML attribute escaping.
func XMLEscape(value string) string {
	return xmlEscaper.Replace(value)
}

// CanonicalPartPath cleans a path to use forward slashes and removes leading slash.
func CanonicalPartPath(target string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(target, "\\", "/"))
	if after, ok := strings.CutPrefix(clean, "/"); ok {
		return after
	}
	return clean
}

// SlideRelsPartName returns the relative path to a slide's relationships part.
func SlideRelsPartName(slidePart string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(slidePart, "\\", "/"))
	if clean == "" {
		return ""
	}
	lastSlash := strings.LastIndex(clean, "/")
	if lastSlash < 0 {
		return "_rels/" + clean + ".rels"
	}
	dir := clean[:lastSlash]
	file := clean[lastSlash+1:]
	return dir + "/_rels/" + file + ".rels"
}

// RelsPathFor returns the relationships part path for any given part path.
func RelsPathFor(partPath string) string {
	clean := strings.TrimSpace(strings.ReplaceAll(partPath, "\\", "/"))
	if clean == "" {
		return ""
	}
	lastSlash := strings.LastIndex(clean, "/")
	if lastSlash < 0 {
		return "_rels/" + clean + ".rels"
	}
	dir := clean[:lastSlash]
	file := clean[lastSlash+1:]
	return dir + "/_rels/" + file + ".rels"
}

// ParseRelationshipNumber extracts the numeric part of an rId string.
func ParseRelationshipNumber(id string) (int, bool) {
	if !strings.HasPrefix(id, "rId") {
		return 0, false
	}
	var num int
	_, err := fmt.Sscanf(id, "rId%d", &num)
	if err != nil {
		return 0, false
	}
	return num, true
}

// NextRelationshipNumber returns the next available relationship number
// (starting from 1) that is not present in the existing relationships list.
func NextRelationshipNumber(rels []EditorRelationship) int {
	nextNum := 1
	for _, r := range rels {
		if num, ok := ParseRelationshipNumber(r.ID); ok && num >= nextNum {
			nextNum = num + 1
		}
	}
	return nextNum
}

// ParseSlidePartNumber extracts the numeric part of a slide part name (e.g., slide1.xml).
func ParseSlidePartNumber(partPath string) (int, bool) {
	lastSlash := strings.LastIndex(partPath, "/")
	base := partPath
	if lastSlash >= 0 {
		base = partPath[lastSlash+1:]
	}
	if !strings.HasPrefix(base, "slide") || !strings.HasSuffix(base, ".xml") {
		return 0, false
	}
	var num int
	_, err := fmt.Sscanf(base, "slide%d.xml", &num)
	if err != nil {
		return 0, false
	}
	return num, true
}

// Shape represents a simplified view of a slide shape for editing.
type Shape struct {
	ID   int
	Name string
	Type string
	Text string
	X, Y int
	W, H int
}

// ShapeSearchQuery filters shapes for editor-wide search.
type ShapeSearchQuery struct {
	NameContains  string
	TypeEquals    string
	TextContains  string
	CaseSensitive bool
}

// ShapeSearchResult identifies one matched shape and its slide index.
type ShapeSearchResult struct {
	SlideIndex int
	Shape      Shape
}

// ChartSelector identifies a slide chart by index and/or relationship ID.
type ChartSelector struct {
	Index *int   `json:"index,omitempty"`
	RelID string `json:"rel_id,omitempty"`
}

// ChartSeriesData carries one chart series worth of input data.
type ChartSeriesData struct {
	Name       *string   `json:"name,omitempty"`
	Categories []string  `json:"categories,omitempty"`
	Values     []float64 `json:"values,omitempty"`
	XValues    []float64 `json:"x_values,omitempty"`
	YValues    []float64 `json:"y_values,omitempty"`
	Sizes      []float64 `json:"sizes,omitempty"`
}

// ChartDataUpdate is the complete chart update payload.
type ChartDataUpdate struct {
	Categories []string          `json:"categories,omitempty"`
	Series     []ChartSeriesData `json:"series,omitempty"`
}

// SlideChartRef describes a chart relationship discovered on a slide.
type SlideChartRef struct {
	Index     int
	RelID     string
	ChartPart string
}

// SlideLayoutInfo describes one available slide layout part.
type SlideLayoutInfo struct {
	Part       string
	Name       string
	MasterPart string
}

// SlideMasterInfo describes one available slide master part.
type SlideMasterInfo struct {
	Part string
}

// SlideMasterCloneResult summarizes an in-package layout/master clone operation.
type SlideMasterCloneResult struct {
	MasterPart string
	ThemePart  string
	LayoutMap  map[string]string
}

// Hyperlink holds properties for click or hover actions.
type Hyperlink struct {
	Address        *string `json:"address,omitempty"`
	Action         *string `json:"action,omitempty"`
	Tooltip        *string `json:"tooltip,omitempty"`
	TargetSlide    *int    `json:"target_slide,omitempty"`
	TargetJump     *string `json:"jump,omitempty"`
	Macro          *string `json:"macro,omitempty"`
	History        *bool   `json:"history,omitempty"`
	HighlightClick *bool   `json:"highlight_click,omitempty"`
	EndSound       *bool   `json:"end_sound,omitempty"`
}

// TextRun represents a single formatted text segment.
type TextRun struct {
	Text          string     `json:"text"`
	Bold          *bool      `json:"bold,omitempty"`
	Italic        *bool      `json:"italic,omitempty"`
	Underline     *string    `json:"underline,omitempty"`
	Strikethrough *string    `json:"strikethrough,omitempty"`
	Subscript     *bool      `json:"subscript,omitempty"`
	Superscript   *bool      `json:"superscript,omitempty"`
	Color         *string    `json:"color,omitempty"`
	Highlight     *string    `json:"highlight,omitempty"`
	Font          *string    `json:"font,omitempty"`
	SizePt        *int       `json:"size_pt,omitempty"`
	Code          *bool      `json:"code,omitempty"`
	AllCaps       *bool      `json:"all_caps,omitempty"`
	SmallCaps     *bool      `json:"small_caps,omitempty"`
	Hyperlink     *Hyperlink `json:"hyperlink,omitempty"`
	HoverAction   *Hyperlink `json:"hover_action,omitempty"`
}

// ShapeProps defines optional properties when creating a shape.
type ShapeProps struct {
	Name string    `json:"name,omitempty"`
	Runs []TextRun `json:"runs,omitempty"`
	// Add other properties as needed for Phase 1
}

// TextFrame defines formatting properties for the text body container within a shape.
type TextFrame struct {
	MarginTop     *int     `json:"margin_top,omitempty"`
	MarginBottom  *int     `json:"margin_bottom,omitempty"`
	MarginLeft    *int     `json:"margin_left,omitempty"`
	MarginRight   *int     `json:"margin_right,omitempty"`
	WordWrap      *bool    `json:"word_wrap,omitempty"`
	AutoFit       *bool    `json:"auto_fit,omitempty"`      // Deprecated: use auto_fit_type instead
	AutoFitType   *string  `json:"auto_fit_type,omitempty"` // "none", "normal", "shape"
	VerticalAlign *string  `json:"vertical_align,omitempty"`
	Orientation   *string  `json:"orientation,omitempty"`
	Columns       *int     `json:"columns,omitempty"`
	Rotation      *float64 `json:"rotation,omitempty"` // Degrees, converted to OOXML 1/60000 degree units.
}

// Paragraph defines paragraph-level formatting controls.
type Paragraph struct {
	Indent  *int `json:"indent,omitempty"`  // Left paragraph margin (EMU, maps to a:pPr marL).
	Hanging *int `json:"hanging,omitempty"` // Hanging indent amount (EMU, rendered as negative a:pPr indent).
}

// ShapeFill defines generic shape fill controls.
type ShapeFill struct {
	Solid      *string        `json:"solid,omitempty"`      // RGB hex (e.g., "FF0000")
	Gradient   *GradientFill  `json:"gradient,omitempty"`   // Linear gradient fill controls.
	Pattern    *PatternedFill `json:"pattern,omitempty"`    // Pattern fill controls.
	Background *bool          `json:"background,omitempty"` // True emits <a:noFill/> (python-pptx fill.background()).
}

// ShapeLine defines generic shape line controls.
type ShapeLine struct {
	Color     *string `json:"color,omitempty"`      // RGB hex (e.g., "00FF00")
	WidthEmu  *int    `json:"width_emu,omitempty"`  // Line width in EMU.
	DashStyle *string `json:"dash_style,omitempty"` // Preset dash style (e.g., "dash", "dashDot", "lgDash").
}

// ShapeShadow defines generic shape shadow controls.
type ShapeShadow struct {
	Inherit     *bool    `json:"inherit,omitempty"`      // Remove explicit effects and inherit when true.
	Color       *string  `json:"color,omitempty"`        // RGB hex.
	BlurEmu     *int     `json:"blur_emu,omitempty"`     // Shadow blur radius in EMU.
	DistanceEmu *int     `json:"distance_emu,omitempty"` // Shadow distance in EMU.
	AngleDeg    *float64 `json:"angle_deg,omitempty"`    // Shadow direction angle in degrees.
}

// ShapeGlow defines generic shape glow controls.
type ShapeGlow struct {
	Color     *string `json:"color,omitempty"`      // RGB hex.
	RadiusEmu *int    `json:"radius_emu,omitempty"` // Glow radius in EMU.
}

// ShapeBlur defines generic shape blur controls.
type ShapeBlur struct {
	RadiusEmu *int `json:"radius_emu,omitempty"` // Blur radius in EMU.
}

// ShapeSoftEdge defines generic shape soft-edge controls.
type ShapeSoftEdge struct {
	RadiusEmu *int `json:"radius_emu,omitempty"` // Soft-edge radius in EMU.
}

// ShapeReflection defines generic shape reflection controls.
type ShapeReflection struct {
	BlurEmu     *int `json:"blur_emu,omitempty"`     // Reflection blur radius in EMU.
	DistanceEmu *int `json:"distance_emu,omitempty"` // Reflection distance in EMU.
}

// GradientStop defines one gradient stop in a linear gradient.
type GradientStop struct {
	PositionPct *float64 `json:"position_pct,omitempty"` // 0.0 to 100.0
	Color       string   `json:"color"`
}

// GradientFill defines linear gradient controls.
type GradientFill struct {
	AngleDeg *float64       `json:"angle_deg,omitempty"`
	Stops    []GradientStop `json:"stops,omitempty"`
}

// PatternedFill defines patterned fill controls.
type PatternedFill struct {
	Preset  *string `json:"preset,omitempty"`   // e.g. "pct5", "diagCross"
	FgColor *string `json:"fg_color,omitempty"` // RGB hex
	BgColor *string `json:"bg_color,omitempty"` // RGB hex
}

// ImageMetadata describes basic image properties returned by the bridge.
type ImageMetadata struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Hash   string `json:"hash,omitempty"`
}

// ImageCrop defines cropping offsets (0.0 to 1.0).
type ImageCrop struct {
	Left   float64 `json:"left,omitempty"`
	Right  float64 `json:"right,omitempty"`
	Top    float64 `json:"top,omitempty"`
	Bottom float64 `json:"bottom,omitempty"`
}

type ShapeUpdate struct {
	Text        *string          `json:"text,omitempty"`
	Runs        *[]TextRun       `json:"runs,omitempty"`
	TextFrame   *TextFrame       `json:"text_frame,omitempty"`
	Paragraph   *Paragraph       `json:"paragraph,omitempty"`
	Fill        *ShapeFill       `json:"fill,omitempty"`
	Line        *ShapeLine       `json:"line,omitempty"`
	Shadow      *ShapeShadow     `json:"shadow,omitempty"`
	Glow        *ShapeGlow       `json:"glow,omitempty"`
	Blur        *ShapeBlur       `json:"blur,omitempty"`
	SoftEdge    *ShapeSoftEdge   `json:"soft_edge,omitempty"`
	Reflection  *ShapeReflection `json:"reflection,omitempty"`
	ClickAction *Hyperlink       `json:"click_action,omitempty"`
	HoverAction *Hyperlink       `json:"hover_action,omitempty"`

	// Image properties (valid if shape is a picture)
	Crop     *ImageCrop `json:"crop,omitempty"`
	Rotation *float64   `json:"rotation,omitempty"`
	FlipH    *bool      `json:"flip_h,omitempty"`
	FlipV    *bool      `json:"flip_v,omitempty"`

	X *int `json:"x,omitempty"`
	Y *int `json:"y,omitempty"`
	W *int `json:"w,omitempty"`
	H *int `json:"h,omitempty"`
}

// SlideImageRef describes one image relationship on a slide.
type SlideImageRef struct {
	Index  int
	RelID  string
	Target string
}
