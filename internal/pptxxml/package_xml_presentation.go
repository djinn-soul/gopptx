package pptxxml

import (
	"strconv"
	"strings"
)

// Section describes a presentation section grouping slides.
type Section struct {
	Name     string
	GUID     string
	SlideIDs []int64
}

// SectionListXML renders ppt/sectionList.xml.
func SectionListXML(sections []Section) string {
	if len(sections) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(sectionListBody("s", "http://schemas.microsoft.com/office/powerpoint/2010/main", sections))
	return b.String()
}

func sectionListBody(prefix, namespace string, sections []Section) string {
	var b strings.Builder
	b.WriteString("\n<")
	b.WriteString(prefix)
	b.WriteString(":sectionLst")
	if namespace != "" {
		b.WriteString(" xmlns:")
		b.WriteString(prefix)
		b.WriteString("=\"")
		b.WriteString(namespace)
		b.WriteString("\"")
	}
	b.WriteString(">")
	for _, s := range sections {
		b.WriteString("\n  <")
		b.WriteString(prefix)
		b.WriteString(":section name=\"")
		b.WriteString(Escape(s.Name))
		b.WriteString("\" id=\"")
		b.WriteString(s.GUID)
		b.WriteString("\">")
		b.WriteString("\n    <")
		b.WriteString(prefix)
		b.WriteString(":sldIdLst>")
		for _, slideID := range s.SlideIDs {
			b.WriteString("\n      <")
			b.WriteString(prefix)
			b.WriteString(":sldId id=\"")
			b.WriteString(strconv.FormatInt(slideID, 10))
			b.WriteString("\"/>")
		}
		b.WriteString("\n    </")
		b.WriteString(prefix)
		b.WriteString(":sldIdLst>")
		b.WriteString("\n  </")
		b.WriteString(prefix)
		b.WriteString(":section>")
	}
	b.WriteString("\n</")
	b.WriteString(prefix)
	b.WriteString(":sectionLst>")
	return b.String()
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
	// SlideRangeAll plays every slide (<p:sldAll/>), the default.
	SlideRangeAll SlideRangeKind = iota
	// SlideRangeRange plays a contiguous run (<p:sldRg st= end=>).
	SlideRangeRange
	// SlideRangeCustom plays a named custom show (<p:custShow id=>).
	SlideRangeCustom
)

// ShowSettings controls how the presentation is shown (p:showPr, which lives in
// presProps.xml — CT_PresentationProperties, not CT_Presentation).
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
	// RangeStart and RangeEnd are 1-based and inclusive, used by RangeKind
	// SlideRangeRange.
	RangeStart int
	RangeEnd   int
	// CustomShowID is the id of the custom show to play, used by RangeKind
	// SlideRangeCustom.
	CustomShowID int
	// PenColor is the annotation pen colour as a hex RGB.
	PenColor string
	// LaserColor is the laser-pointer colour as a hex RGB. PowerPoint carries
	// it in a p14 extension rather than an attribute.
	LaserColor string
	// ShowMediaControls displays the playback controls over media during the
	// show, also a p14 extension.
	ShowMediaControls bool
}

// IsZero reports whether all fields are at their default values (no showPr needed).
func (s ShowSettings) IsZero() bool {
	return !s.Loop && s.Mode == ShowModePresent && !s.DisableTimings && !s.HideAnimation &&
		!s.DisableNarration && s.RangeKind == SlideRangeAll && s.PenColor == "" &&
		s.LaserColor == "" && !s.ShowMediaControls
}

// ShowPrXML renders the <p:showPr> element, or empty string if all defaults.
// CT_ShowProperties orders the show mode, then the slide range, then the pen
// colour, then the extension list.
func ShowPrXML(s ShowSettings) string {
	if s.IsZero() {
		return ""
	}
	var b strings.Builder
	b.WriteString("<p:showPr")
	if s.Loop {
		b.WriteString(` loop="1"`)
	}
	if s.DisableNarration {
		b.WriteString(` showNarration="0"`)
	}
	if s.DisableTimings {
		b.WriteString(` useTimings="0"`)
	}
	if s.HideAnimation {
		b.WriteString(` showAnimation="0"`)
	}
	b.WriteString(">")

	switch s.Mode {
	case ShowModeKiosk:
		b.WriteString("<p:kiosk/>")
	case ShowModeBrowse:
		if !s.ShowScrollbar {
			b.WriteString(`<p:browse showScrollbar="0"/>`)
		} else {
			b.WriteString("<p:browse/>")
		}
	default:
		b.WriteString("<p:present/>")
	}

	b.WriteString(showSlideRangeXML(s))
	if s.PenColor != "" {
		b.WriteString(`<p:penClr><a:srgbClr val="` + Escape(s.PenColor) + `"/></p:penClr>`)
	}
	b.WriteString(showExtensionsXML(s))

	b.WriteString("</p:showPr>")
	return b.String()
}

func showSlideRangeXML(s ShowSettings) string {
	switch s.RangeKind {
	case SlideRangeRange:
		return `<p:sldRg st="` + strconv.Itoa(s.RangeStart) +
			`" end="` + strconv.Itoa(s.RangeEnd) + `"/>`
	case SlideRangeCustom:
		return `<p:custShow id="` + strconv.Itoa(s.CustomShowID) + `"/>`
	default:
		return "<p:sldAll/>"
	}
}

// showExtensionsXML carries the two settings PowerPoint keeps in p14
// extensions rather than attributes.
func showExtensionsXML(s ShowSettings) string {
	if s.LaserColor == "" && !s.ShowMediaControls {
		return ""
	}
	var b strings.Builder
	b.WriteString("<p:extLst>")
	if s.LaserColor != "" {
		b.WriteString(`<p:ext uri="{EC167BDD-8182-4AB7-AECC-EB403E3ABB37}">` +
			`<p14:laserClr xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main">` +
			`<a:srgbClr val="` + Escape(s.LaserColor) + `"/></p14:laserClr></p:ext>`)
	}
	if s.ShowMediaControls {
		b.WriteString(`<p:ext uri="{2FDB2607-1784-4EEB-B798-7EB5836EED8A}">` +
			`<p14:showMediaCtrls xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" val="1"/>` +
			`</p:ext>`)
	}
	b.WriteString("</p:extLst>")
	return b.String()
}

// CustomShow is a named subset of slides, played in its own order.
type CustomShow struct {
	Name string
	ID   int
	// SlideRelIDs are the presentation relationship ids of the slides, in play
	// order: <p:sld> points at the slide part, not at its p:sldId.
	SlideRelIDs []string
}

// CustomShowListXML renders <p:custShowLst>, which names each show and the
// slides it plays. PowerPoint has no other representation of a custom show.
func CustomShowListXML(shows []CustomShow) string {
	if len(shows) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n<p:custShowLst>")
	for _, show := range shows {
		b.WriteString(`
<p:custShow name="` + Escape(show.Name) + `" id="` + strconv.Itoa(show.ID) + `">
<p:sldLst>`)
		for _, relID := range show.SlideRelIDs {
			b.WriteString(`<p:sld r:id="` + Escape(relID) + `"/>`)
		}
		b.WriteString(`</p:sldLst>
</p:custShow>`)
	}
	b.WriteString("\n</p:custShowLst>")
	return b.String()
}

// ProtectionInfo defines the XML attributes for p:modifyVerifier.
type ProtectionInfo struct {
	HashAlgSID int
	HashData   string
	SaltData   string
	SpinCount  int
}

// Presentation renders ppt/presentation.xml.
//
//nolint:funlen // Presentation XML root contains many optional sections emitted in one ordered block.
func Presentation(
	title string,
	slideCount int,
	includeNotesMaster bool,
	width, height int64,
	masterCount int,
	protection *ProtectionInfo,
	sections []Section,
	rtl bool, // Note: rtl="1" only enables UI direction; content elements (text, etc.) may need individual alignment.
	embeddedFonts []EmbeddedFontRef,
	customShows []CustomShow,
	handoutMasterRelID string,
) string {
	_ = title
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" saveSubsetFonts="1"` +
		// With autoCompressPictures absent PowerPoint applies its own default on
		// save, which can re-encode images the caller supplied at full quality.
		` autoCompressPictures="0"`)
	if rtl {
		b.WriteString(` rtl="1"`)
	}
	b.WriteString(`>
<p:sldMasterIdLst>`)
	for i := range masterCount {
		// Keep IDs globally unique across masters + layout IDs (block size: one
		// master plus its layouts).
		//nolint:mnd // OOXML master ID base
		masterID := int64(2147483648) + int64(i*(LayoutsPerMaster+1))
		rid := i + 1
		b.WriteString("\n<p:sldMasterId id=\"")
		b.WriteString(strconv.FormatInt(masterID, 10))
		b.WriteString("\" r:id=\"rId")
		b.WriteString(strconv.Itoa(rid))
		b.WriteString("\"/>")
	}
	b.WriteString(`
</p:sldMasterIdLst>`)

	if includeNotesMaster {
		rid := masterCount + slideCount + 2
		b.WriteString(`
<p:notesMasterIdLst>
<p:notesMasterId r:id="rId`)
		b.WriteString(strconv.Itoa(rid))
		b.WriteString(`"/>
</p:notesMasterIdLst>`)
	}

	// CT_Presentation orders the handout master list right after the notes
	// master one. Without it the handout part is written and related but never
	// declared, so PowerPoint uses its own built-in handout instead.
	if handoutMasterRelID != "" {
		b.WriteString(`
<p:handoutMasterIdLst>
<p:handoutMasterId r:id="`)
		b.WriteString(Escape(handoutMasterRelID))
		b.WriteString(`"/>
</p:handoutMasterIdLst>`)
	}

	b.WriteString(`
<p:sldIdLst>`)
	for i := 1; i <= slideCount; i++ {
		//nolint:mnd // OOXML slide ID base and rId offset
		slideID := 256 + i
		rid := masterCount + 1 + i
		b.WriteString("\n<p:sldId id=\"")
		b.WriteString(strconv.Itoa(slideID))
		b.WriteString("\" r:id=\"rId")
		b.WriteString(strconv.Itoa(rid))
		b.WriteString("\"/>")
	}

	typeAttr := "custom"
	if width == 9144000 && height == 6858000 {
		typeAttr = "screen4x3"
	} else if width == 12192000 && height == 6858000 {
		typeAttr = "screen16x9"
	}

	b.WriteString(`
</p:sldIdLst>
<p:sldSz cx="`)
	b.WriteString(strconv.FormatInt(width, 10))
	b.WriteString(`" cy="`)
	b.WriteString(strconv.FormatInt(height, 10))
	b.WriteString(`" type="`)
	b.WriteString(typeAttr)
	b.WriteString(`"/>
<p:notesSz cx="6858000" cy="9144000"/>`)

	if len(embeddedFonts) > 0 {
		b.WriteString(EmbeddedFontsXML(embeddedFonts))
	}

	// CT_Presentation orders the custom-show list after the font list, then
	// defaultTextStyle, then the modify verifier.
	b.WriteString(CustomShowListXML(customShows))
	b.WriteString(DefaultTextStyleXML())

	if protection != nil {
		algSid := protection.HashAlgSID
		if algSid == 0 {
			algSid = 14 // SHA-512 in Office crypto SID mapping
		}
		b.WriteString(`
<p:modifyVerifier cryptProviderType="rsaAES" cryptAlgorithmClass="hash" cryptAlgorithmType="typeAny" cryptAlgorithmSid="`)
		b.WriteString(strconv.Itoa(algSid))
		b.WriteString(`" spinCount="`)
		b.WriteString(strconv.Itoa(protection.SpinCount))
		b.WriteString(`" saltData="`)
		b.WriteString(Escape(protection.SaltData))
		b.WriteString(`" hashData="`)
		b.WriteString(Escape(protection.HashData))
		b.WriteString(`"/>`)
	}

	// The extension list closes CT_Presentation, so it goes last.
	if len(sections) > 0 {
		b.WriteString(`
<p:extLst>
<p:ext uri="{521415D9-36F7-43E2-AB2F-B90AF26B5E84}">`)
		b.WriteString(sectionListBody("p14", "http://schemas.microsoft.com/office/powerpoint/2010/main", sections))
		b.WriteString(`
</p:ext>
</p:extLst>`)
	}

	b.WriteString(`
</p:presentation>`)
	return b.String()
}
