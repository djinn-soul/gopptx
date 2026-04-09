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

// ShowSettings controls how the presentation is shown (p:showPr in presentation.xml).
type ShowSettings struct {
	Loop           bool     // Loop presentation continuously when finished
	Mode           ShowMode // Present (default), Browse, or Kiosk
	ShowScrollbar  bool     // Show scrollbar in browse mode
	DisableTimings bool     // Ignore slide timings (useTimings="0")
	HideAnimation  bool     // Suppress animations (showAnimation="0")
}

// IsZero reports whether all fields are at their default values (no showPr needed).
func (s ShowSettings) IsZero() bool {
	return !s.Loop && s.Mode == ShowModePresent && !s.DisableTimings && !s.HideAnimation
}

// ShowPrXML renders the <p:showPr> element, or empty string if all defaults.
func ShowPrXML(s ShowSettings) string {
	if s.IsZero() {
		return ""
	}
	var b strings.Builder
	b.WriteString("<p:showPr")
	if s.Loop {
		b.WriteString(` loop="1"`)
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
	b.WriteString("</p:showPr>")
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
	show *ShowSettings,
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
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" saveSubsetFonts="1"`)
	if rtl {
		b.WriteString(` rtl="1"`)
	}
	b.WriteString(`>
<p:sldMasterIdLst>`)
	for i := range masterCount {
		// Keep IDs globally unique across masters + layout IDs (block size: 1 master + 6 layouts).
		//nolint:mnd // OOXML master ID base
		masterID := int64(2147483648) + int64(i*7)
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
		//nolint:mnd // OOXML notes master relationship ID offset
		rid := masterCount + slideCount + 2
		b.WriteString(`
<p:notesMasterIdLst>
<p:notesMasterId r:id="rId`)
		b.WriteString(strconv.Itoa(rid))
		b.WriteString(`"/>
</p:notesMasterIdLst>`)
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

	if len(embeddedFonts) > 0 {
		b.WriteString(EmbeddedFontsXML(embeddedFonts))
	}

	if len(sections) > 0 {
		b.WriteString(`
<p:extLst>
<p:ext uri="{521415D9-36F7-43E2-AB2F-B90AF26B5E84}">`)
		b.WriteString(sectionListBody("p14", "http://schemas.microsoft.com/office/powerpoint/2010/main", sections))
		b.WriteString(`
</p:ext>
</p:extLst>`)
	}

	if show != nil {
		if xml := ShowPrXML(*show); xml != "" {
			b.WriteString("\n")
			b.WriteString(xml)
		}
	}

	b.WriteString(`
</p:presentation>`)
	return b.String()
}
