package pptxxml

import (
	"fmt"
	"strconv"
	"strings"
)

// Escape replaces XML-sensitive characters with entity references.
func Escape(value string) string {
	return xmlEscapeReplacer.Replace(value)
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
			notesThemeIndex = 2
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
		`ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
</Types>`)

	return b.String()
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
func RootRelationships() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" ` +
		`Target="ppt/presentation.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" ` +
		`Target="docProps/core.xml"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" ` +
		`Target="docProps/app.xml"/>
</Relationships>`
}

// PresentationRelationships renders ppt/_rels/presentation.xml.rels.
func PresentationRelationships(slideCount int, includeNotesMaster bool, customXMLCount int, masterCount int) string {
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)

	// Master relationships: rId1..rIdN for N masters
	for i := range masterCount {
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster%d.xml"/>`, i+1, i+1))
	}

	// Theme relationship: rId(masterCount+1)
	b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>`, masterCount+1))

	// Slide relationships: rId(masterCount+2)..rId(masterCount+slideCount+1)
	for i := 1; i <= slideCount; i++ {
		rid := masterCount + 1 + i
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide%d.xml"/>`, rid, i))
	}
	if includeNotesMaster {
		//nolint:mnd // OOXML relationship ID offset
		rid := masterCount + slideCount + 2
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="notesMasters/notesMaster1.xml"/>`, rid))
	}

	//nolint:mnd // OOXML relationship ID offset
	baseRid := masterCount + slideCount + 2
	if includeNotesMaster {
		baseRid++
	}
	for i := 1; i <= customXMLCount; i++ {
		//nolint:mnd // Custom XML ID pair spacing
		rid := baseRid + (i-1)*2
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml" `+
			`Target="../customXml/item%d.xml"/>`, rid, i))
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" `+
			`Target="../customXml/itemProps%d.xml"/>`, rid+1, i))
	}

	b.WriteString(`
</Relationships>`)
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
) string {
	_ = title
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" saveSubsetFonts="1">
<p:sldMasterIdLst>`)
	for i := range masterCount {
		// Keep IDs globally unique across masters + layout IDs (block size: 1 master + 6 layouts).
		//nolint:mnd // OOXML master ID base
		masterID := int64(2147483648) + int64(i*7)
		rid := i + 1
		b.WriteString(fmt.Sprintf(`
<p:sldMasterId id="%d" r:id="rId%d"/>`, masterID, rid))
	}
	b.WriteString(`
</p:sldMasterIdLst>`)

	if includeNotesMaster {
		//nolint:mnd // OOXML notes master relationship ID offset
		rid := masterCount + slideCount + 2
		b.WriteString(fmt.Sprintf(`
<p:notesMasterIdLst>
<p:notesMasterId r:id="rId%d"/>
</p:notesMasterIdLst>`, rid))
	}

	b.WriteString(`
<p:sldIdLst>`)

	for i := 1; i <= slideCount; i++ {
		//nolint:mnd // OOXML slide ID base and rId offset
		slideID := 256 + i
		rid := masterCount + 1 + i
		b.WriteString(fmt.Sprintf(`
<p:sldId id="%d" r:id="rId%d"/>`, slideID, rid))
	}

	typeAttr := "custom"
	if width == 9144000 && height == 6858000 {
		typeAttr = "screen4x3"
	} else if width == 12192000 && height == 6858000 {
		typeAttr = "screen16x9"
	}

	b.WriteString(fmt.Sprintf(`
</p:sldIdLst>
<p:sldSz cx="%d" cy="%d" type="%s"/>
<p:notesSz cx="6858000" cy="9144000"/>`, width, height, typeAttr))

	if protection != nil {
		algSid := protection.HashAlgSID
		if algSid == 0 {
			algSid = 14 // SHA-512 in Office crypto SID mapping
		}
		b.WriteString(fmt.Sprintf(`
<p:modifyVerifier cryptProviderType="rsaAES" cryptAlgorithmClass="hash" cryptAlgorithmType="typeAny" cryptAlgorithmSid="%d" spinCount="%d" saltData="%s" hashData="%s"/>`,
			algSid, protection.SpinCount, protection.SaltData, protection.HashData))
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
