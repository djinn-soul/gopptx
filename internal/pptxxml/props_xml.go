package pptxxml

import (
	"fmt"
	"time"
)

// CorePropertiesInfo contains metadata for docProps/core.xml.
type CorePropertiesInfo struct {
	Title       string
	Subject     string
	Creator     string
	Description string
}

// CoreProperties renders docProps/core.xml.
func CoreProperties(info CorePropertiesInfo) string {
	creator := info.Creator
	if creator == "" {
		creator = "gopptx"
	}

	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" `+
		`xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" `+
		`xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<dc:title>%s</dc:title>
<dc:subject>%s</dc:subject>
<dc:creator>%s</dc:creator>
<cp:lastModifiedBy>%s</cp:lastModifiedBy>
<dc:description>%s</dc:description>
<cp:revision>1</cp:revision>
<dcterms:created xsi:type="dcterms:W3CDTF">%s</dcterms:created>
<dcterms:modified xsi:type="dcterms:W3CDTF">%s</dcterms:modified>
</cp:coreProperties>`, Escape(info.Title), Escape(info.Subject), Escape(creator),
		Escape(creator), Escape(info.Description), now, now)
}

// AppProperties renders docProps/app.xml.
func AppProperties(slideCount int, notesCount int, width, height int64) string {
	format := "Custom"
	if width == 9144000 && height == 6858000 {
		format = "On-screen Show (4:3)"
	} else if width == 12192000 && height == 6858000 {
		format = "Widescreen"
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" `+
		`xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
<TotalTime>0</TotalTime>
<Words>0</Words>
<Application>gopptx</Application>
<PresentationFormat>%s</PresentationFormat>
<Paragraphs>0</Paragraphs>
<Slides>%d</Slides>
<Notes>%d</Notes>
<HiddenSlides>0</HiddenSlides>
<MMClips>0</MMClips>
<ScaleCrop>false</ScaleCrop>
<LinksUpToDate>false</LinksUpToDate>
<SharedDoc>false</SharedDoc>
<HyperlinksChanged>false</HyperlinksChanged>
<AppVersion>1.0000</AppVersion>
</Properties>`, format, slideCount, notesCount)
}
