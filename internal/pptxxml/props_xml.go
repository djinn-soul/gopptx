package pptxxml

import (
	"fmt"
	"strings"
	"time"
)

// CorePropertiesInfo contains metadata for docProps/core.xml. It mirrors the
// fields common.CoreProperties declares: the generator used to carry four of
// them and hardcode the rest, so anything else a caller set was dropped on the
// way to the file.
type CorePropertiesInfo struct {
	Title          string
	Subject        string
	Creator        string
	Description    string
	Keywords       string
	Category       string
	ContentStatus  string
	Identifier     string
	Language       string
	Version        string
	LastModifiedBy string
	Revision       string
	LastPrinted    string
	Created        string
	Modified       string
}

// CoreProperties renders docProps/core.xml.
func CoreProperties(info CorePropertiesInfo) string {
	creator := info.Creator
	if creator == "" {
		creator = "gopptx"
	}
	lastModifiedBy := info.LastModifiedBy
	if lastModifiedBy == "" {
		lastModifiedBy = creator
	}
	revision := info.Revision
	if revision == "" {
		revision = "1"
	}
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	created := info.Created
	if created == "" {
		created = now
	}
	modified := info.Modified
	if modified == "" {
		modified = now
	}

	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" ` +
		`xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:dcterms="http://purl.org/dc/terms/" ` +
		`xmlns:dcmitype="http://purl.org/dc/dcmitype/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">`)

	writeElement := func(tag, value string) {
		if value == "" {
			return
		}
		b.WriteString("\n<" + tag + ">" + Escape(value) + "</" + tag + ">")
	}

	// CT_CoreProperties is a sequence, so the order below is the schema's.
	writeElement("dc:title", info.Title)
	writeElement("dc:subject", info.Subject)
	writeElement("dc:creator", creator)
	writeElement("cp:keywords", info.Keywords)
	writeElement("dc:description", info.Description)
	writeElement("cp:lastModifiedBy", lastModifiedBy)
	writeElement("cp:revision", revision)
	if info.LastPrinted != "" {
		writeElement("cp:lastPrinted", info.LastPrinted)
	}
	b.WriteString("\n<dcterms:created xsi:type=\"dcterms:W3CDTF\">" + Escape(created) + "</dcterms:created>")
	b.WriteString("\n<dcterms:modified xsi:type=\"dcterms:W3CDTF\">" + Escape(modified) + "</dcterms:modified>")
	writeElement("cp:category", info.Category)
	writeElement("cp:contentStatus", info.ContentStatus)
	writeElement("dc:identifier", info.Identifier)
	writeElement("dc:language", info.Language)
	writeElement("cp:version", info.Version)

	b.WriteString("\n</cp:coreProperties>")
	return b.String()
}

// AppPropertiesInfo contains the knobs docProps/app.xml exposes. Application,
// Company and Manager used to be hardcoded or absent, and HiddenSlides was
// always zero even when the deck had hidden slides.
type AppPropertiesInfo struct {
	SlideCount   int
	NotesCount   int
	HiddenSlides int
	Width        int64
	Height       int64
	Application  string
	AppVersion   string
	Company      string
	Manager      string
	Words        int
	Paragraphs   int
	TotalTime    int
}

// AppProperties renders docProps/app.xml.
func AppProperties(info AppPropertiesInfo) string {
	format := "Custom"
	if info.Width == 9144000 && info.Height == 6858000 {
		format = "On-screen Show (4:3)"
	} else if info.Width == 12192000 && info.Height == 6858000 {
		format = "Widescreen"
	}
	application := info.Application
	if application == "" {
		application = "gopptx"
	}
	appVersion := info.AppVersion
	if appVersion == "" {
		appVersion = "1.0000"
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" `+
		`xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes">
<TotalTime>%d</TotalTime>
<Words>%d</Words>
<Application>%s</Application>
<PresentationFormat>%s</PresentationFormat>
<Paragraphs>%d</Paragraphs>
<Slides>%d</Slides>
<Notes>%d</Notes>
<HiddenSlides>%d</HiddenSlides>
<MMClips>0</MMClips>
<ScaleCrop>false</ScaleCrop>`,
		info.TotalTime, info.Words, Escape(application), format,
		info.Paragraphs, info.SlideCount, info.NotesCount, info.HiddenSlides)

	if info.Company != "" {
		b.WriteString("\n<Company>" + Escape(info.Company) + "</Company>")
	}
	if info.Manager != "" {
		b.WriteString("\n<Manager>" + Escape(info.Manager) + "</Manager>")
	}

	b.WriteString(`
<LinksUpToDate>false</LinksUpToDate>
<SharedDoc>false</SharedDoc>
<HyperlinksChanged>false</HyperlinksChanged>
<AppVersion>` + Escape(appVersion) + `</AppVersion>
</Properties>`)
	return b.String()
}
