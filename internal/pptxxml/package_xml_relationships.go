package pptxxml

import (
	"strconv"
	"strings"
)

// RootRelationships renders _rels/.rels.
func RootRelationships(hasCustomProps, hasSignatures bool) string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="ppt/presentation.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>`)
	rID := 4
	if hasCustomProps {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(rID))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/custom-properties" Target="docProps/custom.xml"/>`,
		)
		rID++
	}
	if hasSignatures {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(rID))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/package/2006/relationships/digital-signature/origin" Target="_xmlsignatures/origin.sigs"/>`,
		)
	}
	b.WriteString("\n</Relationships>")
	return b.String()
}

// PresentationRelationships renders ppt/_rels/presentation.xml.rels.
//
//nolint:funlen // Relationship writer enumerates all optional package relationships explicitly.
func PresentationRelationships(
	slideCount int,
	includeNotesMaster bool,
	customXMLCount int,
	masterCount int,
	hasSections bool,
	hasCommentAuthors bool,
	hasVBA bool,
	hasHandoutMaster bool,
	embeddedFontsCount int,
) string {
	if masterCount < 1 {
		masterCount = 1
	}
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)

	nextRid := 1
	for i := range masterCount {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			"\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster\" Target=\"slideMasters/slideMaster",
		)
		b.WriteString(strconv.Itoa(i + 1))
		b.WriteString(".xml\"/>")
		nextRid++
	}

	b.WriteString(`
<Relationship Id="rId`)
	b.WriteString(strconv.Itoa(nextRid))
	b.WriteString(
		`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/>`,
	)
	nextRid++

	for i := 1; i <= slideCount; i++ {
		b.WriteString("\n<Relationship Id=\"rId")
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			"\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide\" Target=\"slides/slide",
		)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(".xml\"/>")
		nextRid++
	}

	if includeNotesMaster {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesMaster" Target="notesMasters/notesMaster1.xml"/>`,
		)
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
	}

	if hasSections {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(`" Type="`)
		b.WriteString(sectionListRelationshipType)
		b.WriteString(`" Target="sectionList.xml"/>`)
		nextRid++
	}

	if hasCommentAuthors {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/commentAuthors" Target="commentAuthors.xml"/>`,
		)
		nextRid++
	}

	if hasVBA {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			`" Type="http://schemas.microsoft.com/office/2006/relationships/vbaProject" Target="vbaProject.bin"/>`,
		)
		nextRid++
	}

	if hasHandoutMaster {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/handoutMaster" Target="handoutMasters/handoutMaster1.xml"/>`,
		)
		nextRid++
	}

	for i := 1; i <= embeddedFontsCount; i++ {
		b.WriteString(`
<Relationship Id="rId`)
		b.WriteString(strconv.Itoa(nextRid))
		b.WriteString(
			`" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/font" Target="fonts/font`,
		)
		b.WriteString(strconv.Itoa(i))
		b.WriteString(`.fntdata"/>`)
		nextRid++
	}

	b.WriteString(`
</Relationships>`)
	return b.String()
}
