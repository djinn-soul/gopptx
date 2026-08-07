package common

import "time"

// Forward declarations or interface if needed?
// Actually, SlideMaster is in elements. elements imports common.
// styling doesn't import common.
// common should probably not import styling or elements if they import common.

// SlideSize describes the dimensions of slides in a presentation in EMUs.
type SlideSize struct {
	Width  int64
	Height int64
}

const (
	width4x3  = 9144000
	height4x3 = 6858000
	width16x9 = 12192000
)

// GetSlideSize4x3 returns the standard 4:3 slide size (10x7.5 inches).
func GetSlideSize4x3() SlideSize {
	return SlideSize{Width: width4x3, Height: height4x3}
}

// GetSlideSize16x9 returns the standard 16:9 widescreen slide size (13.33x7.5 inches).
func GetSlideSize16x9() SlideSize {
	return SlideSize{Width: width16x9, Height: height4x3}
}

// CustomXMLKV is a key-value property for a CustomXMLPart.
type CustomXMLKV struct {
	Key   string
	Value string
}

// CustomXMLPart represents an embedded custom XML document in the PPTX package.
// If RootElement is populated, the XML is generated structurally.
// Otherwise, Content must be a valid XML string for legacy passthrough.
type CustomXMLPart struct {
	ItemID      string        `json:"itemID,omitempty"`
	Content     string        `json:"content,omitempty"     xml:",innerxml"`
	RootElement string        `json:"rootElement,omitempty"`
	Namespace   string        `json:"namespace,omitempty"`
	Properties  []CustomXMLKV `json:"properties,omitempty"`
}

// Metadata describes summary information for a PPTX package.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Metadata struct {
	Title          string
	Subject        string
	Creator        string
	Description    string
	FooterText     string
	ShowDateTime   bool
	GeneratedDate  time.Time
	SlideSize      SlideSize
	SlideCount     int
	CustomXML      []CustomXMLPart
	CoreProperties CoreProperties
	AppProperties  AppProperties
	Protection     Protection
	ShowSettings   ShowSettings
	// CustomShows are the named slide subsets the deck offers.
	CustomShows []CustomShow
}

// AppProperties represents the writable part of docProps/app.xml. The
// generator used to hardcode Application and emit neither Company nor Manager.
type AppProperties struct {
	// Application names the producer. Empty writes "gopptx".
	Application string
	// AppVersion is the producer version. Empty writes "1.0000".
	AppVersion string
	Company    string
	Manager    string
}

// ShowMode defines the slide show presentation mode.
type ShowMode int

const (
	ShowModePresent ShowMode = iota // Standard presenter view (default)
	ShowModeBrowse                  // Browse in window
	ShowModeKiosk                   // Kiosk: full-screen, no controls
)

// SlideRangeKind selects which slides a show plays.
type SlideRangeKind int

const (
	// SlideRangeAll plays every slide, the default.
	SlideRangeAll SlideRangeKind = iota
	// SlideRangeRange plays a contiguous run of slides.
	SlideRangeRange
	// SlideRangeCustom plays one of the presentation's custom shows.
	SlideRangeCustom
)

// ShowSettings controls how a presentation is shown (maps to p:showPr, which
// lives in presProps.xml).
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type ShowSettings struct {
	Loop           bool     // Loop presentation continuously when finished
	Mode           ShowMode // Present (default), Browse, or Kiosk
	ShowScrollbar  bool     // Show scrollbar in browse mode
	DisableTimings bool     // Ignore slide timings (useTimings="0")
	HideAnimation  bool     // Suppress animations (showAnimation="0")
	// DisableNarration plays the show without its recorded narration.
	DisableNarration bool
	// RangeKind selects all slides, a range, or a custom show.
	RangeKind SlideRangeKind
	// RangeStart and RangeEnd are 1-based, inclusive, and used when RangeKind
	// is SlideRangeRange.
	RangeStart int
	RangeEnd   int
	// CustomShowName names the show to play when RangeKind is SlideRangeCustom.
	CustomShowName string
	// PenColor is the annotation pen colour, as hex RGB.
	PenColor string
	// LaserColor is the laser-pointer colour, as hex RGB.
	LaserColor string
	// ShowMediaControls displays playback controls over media during the show.
	ShowMediaControls bool
}

// CustomShow is a named subset of the deck's slides, played in its own order —
// the feature PowerPoint calls a custom slide show.
type CustomShow struct {
	// Name is what PowerPoint lists in the custom-show picker.
	Name string
	// SlideIndices are 0-based indices into the presentation's slides, in the
	// order the show plays them. A slide may appear more than once.
	SlideIndices []int
}

// Protection defines write-protection and suggested read-only settings.
//
//nolint:govet // Preserve public field order for source compatibility with positional literals.
type Protection struct {
	ModifyPassword    string
	MarkAsFinal       bool
	SignaturesEnabled bool
	EncryptPassword   string
}

// CoreProperties represents the docProps/core.xml metadata.
type CoreProperties struct {
	Title          string `json:"title"          xml:"http://purl.org/dc/elements/1.1/ title,omitempty"`
	Subject        string `json:"subject"        xml:"http://purl.org/dc/elements/1.1/ subject,omitempty"`
	Creator        string `json:"creator"        xml:"http://purl.org/dc/elements/1.1/ creator,omitempty"` // Also known as "author" in python-pptx
	Keywords       string `json:"keywords"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties keywords,omitempty"`
	Description    string `json:"description"    xml:"http://purl.org/dc/elements/1.1/ description,omitempty"` // Also known as "comments" in python-pptx
	LastModifiedBy string `json:"lastModifiedBy" xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties lastModifiedBy,omitempty"`
	Revision       string `json:"revision"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties revision,omitempty"` //nolint:lll // struct tags cannot be split
	Created        string `json:"created"        xml:"http://purl.org/dc/terms/ created,omitempty"`
	Modified       string `json:"modified"       xml:"http://purl.org/dc/terms/ modified,omitempty"`
	Category       string `json:"category"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties category,omitempty"`      //nolint:lll // struct tags cannot be split
	ContentStatus  string `json:"contentStatus"  xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties contentStatus,omitempty"` //nolint:lll // struct tags cannot be split
	Identifier     string `json:"identifier"     xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties identifier,omitempty"`    //nolint:lll // struct tags cannot be split
	Language       string `json:"language"       xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties language,omitempty"`      //nolint:lll // struct tags cannot be split
	LastPrinted    string `json:"lastPrinted"    xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties lastPrinted,omitempty"`   //nolint:lll // struct tags cannot be split
	Version        string `json:"version"        xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties version,omitempty"`       //nolint:lll // struct tags cannot be split
}
