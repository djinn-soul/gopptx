package pptxxml

import (
	"strconv"
	"strings"
)

// ContentTypes renders [Content_Types].xml.
//
//nolint:gocognit,funlen // OPC content-type emission branches over many optional package parts by design.
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
	hasCustomProps bool,
	hasSignatures bool,
	hasVBA bool,
	hasHandoutMaster bool,
	hasEmbeddedFonts bool,
) string {
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.Grow(4096 + slideCount*160 + chartCount*120 + smartArtCount*560 + len(notesSlides)*140 +
		customXMLCount*220 + masterCount*560 + len(imageExtensions)*96)
	b.WriteString(xmlHeader)
	b.WriteString(`
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>`)
	if hasVBA {
		b.WriteString(`
<Default Extension="bin" ContentType="application/vnd.ms-office.vbaProject"/>`)
	}
	if hasEmbeddedFonts {
		b.WriteString(`
<Default Extension="fntdata" ContentType="application/x-fontdata"/>`)
	}
	if hasVBA {
		b.WriteString(`
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.ms-powerpoint.presentation.macroEnabled.main+xml"/>`)
	} else {
		b.WriteString(`
<Override PartName="/ppt/presentation.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"/>`)
	}

	if hasSections {
		b.WriteString(`
<Override PartName="/ppt/sectionList.xml" ContentType="application/vnd.microsoft.powerpoint.sectionList+xml"/>`)
	}

	writeInt := func(v int) {
		b.WriteString(strconv.Itoa(v))
	}

	for _, rawExt := range imageExtensions {
		ext := strings.TrimPrefix(strings.ToLower(rawExt), ".")
		contentType, ok := imageContentType(ext)
		if !ok {
			contentType = "application/octet-stream"
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
		b.WriteString(
			`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"/>`,
		)
	}
	if includeNotesMaster {
		b.WriteString(`
<Override PartName="/ppt/notesMasters/notesMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml"/>`)
	}
	if hasHandoutMaster {
		b.WriteString(`
<Override PartName="/ppt/handoutMasters/handoutMaster1.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.handoutMaster+xml"/>`)
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

	if hasVBA {
		b.WriteString(`
<Override PartName="/ppt/vbaProject.bin" ContentType="application/vnd.ms-office.vbaProject"/>`)
	}

	for i := 1; i <= masterCount*6; i++ {
		b.WriteString(`
<Override PartName="/ppt/slideLayouts/slideLayout`)
		writeInt(i)
		b.WriteString(
			`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"/>`,
		)
	}
	for i := 1; i <= masterCount; i++ {
		b.WriteString(`
<Override PartName="/ppt/slideMasters/slideMaster`)
		writeInt(i)
		b.WriteString(
			`.xml" ContentType="application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"/>`,
		)
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

	if hasCustomProps {
		b.WriteString(`
<Override PartName="/docProps/custom.xml" ContentType="application/vnd.openxmlformats-officedocument.custom-properties+xml"/>`)
	}

	if hasSignatures {
		b.WriteString(`
<Override PartName="/_xmlsignatures/origin.sigs" ContentType="application/vnd.openxmlformats-package.digital-signature-origin"/>`)
	}

	b.WriteString(`
</Types>`)
	return b.String()
}

func imageContentType(ext string) (string, bool) {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
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
	case "emf":
		return "image/x-emf", true
	case "wmf":
		return "image/x-wmf", true
	case "wdp", "hdp":
		return "image/vnd.ms-photo", true
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
