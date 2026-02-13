package pptxxml

import (
	"fmt"
	"strings"
)

var xmlEscapeReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&quot;",
	"'", "&apos;",
)

// Escape replaces XML-sensitive characters with entity references.
func Escape(value string) string {
	return xmlEscapeReplacer.Replace(value)
}

var imageContentTypes = map[string]string{
	"png":  "image/png",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"gif":  "image/gif",
	"bmp":  "image/bmp",
	"tif":  "image/tiff",
	"tiff": "image/tiff",
	"wav":  "audio/wav",
	"mp3":  "audio/mpeg",
	"m4a":  "audio/mp4",
}

// ContentTypes renders [Content_Types].xml.
func ContentTypes(slideCount int, imageExtensions []string, chartCount int, notesSlides []int, includeNotesMaster bool, customXMLCount int) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`)

	for _, ext := range imageExtensions {
		contentType, ok := imageContentTypes[ext]
		if !ok {
			panic(fmt.Sprintf("unsupported image extension in content types: %s", ext))
		}
		b.WriteString(fmt.Sprintf(`
<Default Extension="%s" ContentType="%s"/>`, ext, contentType))
	}

	for i := 1; i <= slideCount; i++ {
		b.WriteString(fmt.Sprintf(`
<Override PartName="/ppt/slides/slide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slide+xml"/>`, i))
	}

	for i := 1; i <= chartCount; i++ {
		b.WriteString(fmt.Sprintf(`
<Override PartName="/ppt/charts/chart%d.xml" ContentType="application/vnd.openxmlformats-officedocument.drawingml.chart+xml"/>`, i))
	}
	for _, slideNumber := range notesSlides {
		b.WriteString(fmt.Sprintf(`
<Override PartName="/ppt/notesSlides/notesSlide%d.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"/>`, slideNumber))
	}
	if includeNotesMaster {
		b.WriteString(`
<Override PartName="/ppt/notesMasters/notesMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml"/>`)
		b.WriteString(`
<Override PartName="/ppt/theme/theme2.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>`)
	}

	for i := 1; i <= customXMLCount; i++ {
		b.WriteString(fmt.Sprintf(`
<Override PartName="/customXml/item%d.xml" ContentType="application/xml"/>`, i))
		b.WriteString(fmt.Sprintf(`
<Override PartName="/customXml/itemProps%d.xml" ContentType="application/vnd.openxmlformats-officedocument.customXmlProperties+xml"/>`, i))
	}

	b.WriteString(`
<Override PartName="/ppt/slideLayouts/slideLayout1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideLayouts/slideLayout2.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideLayouts/slideLayout3.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideLayouts/slideLayout4.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideLayouts/slideLayout5.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideLayouts/slideLayout6.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>
<Override PartName="/ppt/slideMasters/slideMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>
<Override PartName="/ppt/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/>
<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
<Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
</Types>`)

	return b.String()
}

// RootRelationships renders _rels/.rels.
func RootRelationships() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`
}

// PresentationRelationships renders ppt/_rels/presentation.xml.rels.
func PresentationRelationships(slideCount int, includeNotesMaster bool, customXMLCount int) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster" Target="slideMasters/slideMaster1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>`)

	for i := 1; i <= slideCount; i++ {
		rid := i + 2
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide" Target="slides/slide%d.xml"/>`, rid, i))
	}
	if includeNotesMaster {
		rid := slideCount + 3
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="notesMasters/notesMaster1.xml"/>`, rid))
	}

	baseRid := slideCount + 3
	if includeNotesMaster {
		baseRid++
	}
	for i := 1; i <= customXMLCount; i++ {
		rid := baseRid + (i-1)*2
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml" Target="../customXml/item%d.xml"/>`, rid, i))
		b.WriteString(fmt.Sprintf(`
<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXmlProps" Target="../customXml/itemProps%d.xml"/>`, rid+1, i))
	}

	b.WriteString(`
</Relationships>`)
	return b.String()
}

// Presentation renders ppt/presentation.xml.
func Presentation(title string, slideCount int, includeNotesMaster bool, width, height int64) string {
	_ = title
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<p:presentation xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" saveSubsetFonts="1">
<p:sldMasterIdLst>
<p:sldMasterId id="2147483648" r:id="rId1"/>
</p:sldMasterIdLst>`)

	if includeNotesMaster {
		rid := slideCount + 3
		b.WriteString(fmt.Sprintf(`
<p:notesMasterIdLst>
<p:notesMasterId r:id="rId%d"/>
</p:notesMasterIdLst>`, rid))
	}

	b.WriteString(`
<p:sldIdLst>`)

	for i := 1; i <= slideCount; i++ {
		slideID := 256 + i
		rid := i + 2
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
<p:notesSz cx="6858000" cy="9144000"/>
</p:presentation>`, width, height, typeAttr))
	return b.String()
}
