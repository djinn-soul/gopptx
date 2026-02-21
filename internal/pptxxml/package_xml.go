package pptxxml

import (
	"strconv"
	"strings"
)

// Escape replaces XML-sensitive characters with entity references.
func Escape(value string) string {
	return xmlEscapeReplacer.Replace(value)
}

// FastEscapeRID is a specialized version of Escape for Relationship IDs (rIdN).
// Since RIDs are known to be alphanumeric, we can skip the expensive Replacer checks.
func FastEscapeRID(rid string) string {
	return rid
}

// WriteRID appends an escaped RID to a builder without extra allocations.
func WriteRID(b *strings.Builder, rid string) {
	b.WriteString(rid)
}

// ContentTypes renders [Content_Types].xml.
func ContentTypes(
	slideCount int,
	imageExtensions []string,
	chartCount int,
	smartArtCount int,
	notesSlides []int,
	includeNotesMaster bool,
	customXMLCount int,
	masterCount int,
	notesThemeIndex int,
	hasSections bool,
	commentSlides []int,
	hasSignatures bool,
) string {
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.Grow(4096 + slideCount*160 + chartCount*120 + smartArtCount*560 + len(notesSlides)*140 +
		customXMLCount*220 + masterCount*560 + len(imageExtensions)*96)
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ` +
		`ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`)

	if hasSections {
		b.WriteString(`
<Override PartName="/ppt/sectionList.xml" ContentType="application/vnd.microsoft.powerpoint.sectionList+xml"/>`)
	}

	writeInt := func(v int) {
		b.WriteString(strconv.Itoa(v))
	}

	for _, ext := range imageExtensions {
		contentType, ok := imageContentType(ext)
		if !ok {
			panic("unsupported image extension in content types: " + ext)
		}
		b.WriteString(`
<Default Extension="`)
		b.WriteString(ext)
		b.WriteString(`" ContentType="`)
		b.WriteString(contentType)
		b.WriteString(`"/>`)
	}

	for i := 1; i <= slideCount; i++ {
		b.WriteString(`
<Override PartName="/ppt/slides/slide`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`)
	}

	for i := 1; i <= chartCount; i++ {
		b.WriteString(`
<Override PartName="/ppt/charts/chart`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.chart+xml"/>`)
	}
	for i := 1; i <= smartArtCount; i++ {
		b.WriteString(`
<Override PartName="/ppt/diagrams/data`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.diagramData+xml"/>
<Override PartName="/ppt/diagrams/layout`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.diagramLayout+xml"/>
<Override PartName="/ppt/diagrams/colors`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.diagramColors+xml"/>
<Override PartName="/ppt/diagrams/quickStyle`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.diagramStyle+xml"/>
<Override PartName="/ppt/diagrams/drawing`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.ms-office.drawingml.diagramDrawing+xml"/>`)
	}

	for _, slideNumber := range notesSlides {
		b.WriteString(`
<Override PartName="/ppt/notesSlides/notesSlide`)
		writeInt(slideNumber)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"/>`)
	}
	if includeNotesMaster {
		b.WriteString(`
<Override PartName="/ppt/notesMasters/notesMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml"/>`)
	}

	for _, slideNumber := range commentSlides {
		b.WriteString(`
<Override PartName="/ppt/comments/comment`)
		writeInt(slideNumber)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.comments+xml"/>`)
	}
	if len(commentSlides) > 0 {
		b.WriteString(`
<Override PartName="/ppt/commentAuthors.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.commentAuthors+xml"/>`)
	}

	for i := 1; i <= customXMLCount; i++ {
		b.WriteString(`
<Override PartName="/customXml/item`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/xml"/>
<Override PartName="/customXml/itemProps`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.customXmlProperties+xml"/>`)
	}

	for i := 1; i <= masterCount*6; i++ {
		b.WriteString(`
<Override PartName="/ppt/slideLayouts/slideLayout`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>`)
	}
	for i := 1; i <= masterCount; i++ {
		b.WriteString(`
<Override PartName="/ppt/slideMasters/slideMaster`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>`)
	}
	for i := 1; i <= masterCount; i++ {
		b.WriteString(`
<Override PartName="/ppt/theme/theme`)
		writeInt(i)
		b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>`)
	}
	if includeNotesMaster {
		if notesThemeIndex < 1 {
			notesThemeIndex = 1
		}
		if notesThemeIndex > masterCount {
			b.WriteString(`
<Override PartName="/ppt/theme/theme`)
			writeInt(notesThemeIndex)
			b.WriteString(`.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>`)
		}
	}
	b.WriteString(`
<Override PartName="/docProps/core.xml" ` +
		`ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
<Override PartName="/docProps/app.xml" ` +
		`ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>`)

	if hasSignatures {
		b.WriteString(`
<Override PartName="/_xmlsignatures/origin.sigs" ContentType="application/vnd.openxmlformats-package.digital-signature-origin"/>`)
	}

	b.WriteString(`
</Types>`)

	return b.String()
}

// SignatureOrigin renders _xmlsignatures/origin.sigs.
func SignatureOrigin() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<SignatureOrigin xmlns="http://schemas.openxmlformats.org/package/2006/digital-signature"/>`
}

// NOTE: The use of a package-level variable here is intentional to avoid repeated [strings.Replacer] allocation.
// Do not move this to a local scope.
//
//nolint:gochecknoglobals // Reused for performance
var xmlEscapeReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&quot;",
	"'", "&apos;",
)

func imageContentType(ext string) (string, bool) {
	switch ext {
	case "png":
		return "image/png", true
	case "jpg", "jpeg":
		return "image/jpeg", true
	case "gif":
		return "image/gif", true
	case "bmp":
		return "image/bmp", true
	case "tif", "tiff":
		return "image/tiff", true
	case "wav":
		return "audio/wav", true
	case "mp3":
		return "audio/mpeg", true
	case "m4a":
		return "audio/mp4", true
	default:
		return "", false
	}
}

// RootRelationships renders _rels/.rels.
func RootRelationships(hasCustomProps, hasSignatures bool) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>`)
	rId := 4
	if hasCustomProps {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(rId))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/custom-properties" Target="docProps/custom.xml"/>`)
		rId++
	}
	if hasSignatures {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(rId))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/package/2006/relationships/digital-signature/origin" Target="_xmlsignatures/origin.sigs"/>`)
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}

// PresentationRelationships renders ppt/_rels/presentation.xml.rels.
func PresentationRelationships(slideCount int, includeNotesMaster bool, customXMLCount int, masterCount int, hasSections bool, hasCommentAuthors bool) string {
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)

	nextRid := 1

	// Master relationships: rId1..rIdN for N masters
	for i := range masterCount {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster\" Target=\"slideMasters/slideMaster")
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(".xml\"/>")
		nextRid++
	}

	// Theme relationship
	b.WriteString(`
<Relationship Id="rId`)
	b.WriteString(strconv.Itoa(nextRid))
	b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>`)
	nextRid++

	// Slide relationships
	for i := 1; i <= slideCount; i++ {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString("\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide\" Target=\"slides/slide")
		b.WriteString(strconv.Itoa(i))
		b.WriteString(".xml\"/>")
		nextRid++
	}

	if includeNotesMaster {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="notesMasters/notesMaster1.xml"/>`)
		nextRid++
	}

	for i := 1; i <= customXMLCount; i++ {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml" ` +
			`Target="../customXml/item`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`.xml"/>`)
		nextRid++

		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" ` +
			`Target="../customXml/itemProps`)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`.xml"/>`)
		nextRid++
	}

	if hasSections {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(`" Type="http://schemas.microsoft.com/office/2006/relationships/sectionList" Target="sectionList.xml"/>`)
		nextRid++
	}

	if hasCommentAuthors {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors" Target="commentAuthors.xml"/>`)
		nextRid++
	}

	b.WriteString(`
</Relationships>`)
	return b.String()
}

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
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	b.WriteString("\n<s:sectionLst xmlns:s=\"http://schemas.microsoft.com/office/powerpoint/2010/main\">")
	for _, s := range sections {
		// PPT uses braces for GUIDs, e.g. "{BA8A57BE-2A2B-4EF9-A77E-97BDDEBBA9AC}"
		b.WriteString("\n  <s:section name=\"")
		b.WriteString(Escape(s.Name))
		b.WriteString("\" id=\"")
		b.WriteString(s.GUID)
		b.WriteString("\">")
		b.WriteString("\n    <s:sldIdLst>")
		for _, slideID := range s.SlideIDs {
			b.WriteString("\n      <s:sldId id=\"")
			b.WriteString(strconv.FormatInt(slideID, 10))
			b.WriteString("\"/>")
		}
		b.WriteString("\n    </s:sldIdLst>")
		b.WriteString("\n  </s:section>")
	}
	b.WriteString("\n</s:sectionLst>")
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
func Presentation(
	title string,
	slideCount int,
	includeNotesMaster bool,
	width, height int64,
	masterCount int,
	protection *ProtectionInfo,
	sections []Section,
	rtl bool, // Note: rtl="1" only enables UI direction; content elements (text, etc.) may need individual alignment.
) string {
	_ = title
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
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

	if len(sections) > 0 {
		b.WriteString(`
<p:extLst>
<p:ext uri="{521415D9-36F7-43E2-AB2F-B90AF26B5E84}">
<p14:sectionLst xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main">`)
		for _, s := range sections {
			b.WriteString("\n<p14:section name=\"")
			b.WriteString(Escape(s.Name))
			b.WriteString("\" id=\"")
			b.WriteString(s.GUID)
			b.WriteString("\">")
			b.WriteString("\n<p14:sldIdLst>")
			for _, sid := range s.SlideIDs {
				b.WriteString("\n<p14:sldId id=\"")
				b.WriteString(strconv.FormatInt(sid, 10))
				b.WriteString("\"/>")
			}
			b.WriteString("\n</p14:sldIdLst>")
			b.WriteString("\n</p14:section>")
		}
		b.WriteString(`
</p14:sectionLst>
</p:ext>
</p:extLst>`)
	}

	b.WriteString(`
</p:presentation>`)
	return b.String()
}

// CustomProperties renders docProps/custom.xml.
func CustomProperties(markAsFinal bool) string {
	if !markAsFinal {
		return ""
	}
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/custom-properties" ` +
		`xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
<property fmtid="{D5CDD505-2E9C-101B-9397-08002B2CF9AE}" pid="2" name="_MarkAsFinal">
<vt:bool>true</vt:bool>
</property>
</Properties>`
}
