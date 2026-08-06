package pptxxml

import (
	"strconv"
	"strings"
)

// Parts PowerPoint writes into every deck it saves, and which its own validator
// lists as required. gopptx generated neither, and the table one matters twice
// over: every table gopptx emits carries an <a:tableStyleId> that had no part
// to resolve against.
const (
	// ViewPropsPartName is the package path of the view properties part.
	ViewPropsPartName = "ppt/viewProps.xml"
	// ViewPropsContentType is the content type override for viewProps.xml.
	ViewPropsContentType = "application/vnd.openxmlformats-officedocument.presentationml.viewProps+xml"
	// ViewPropsRelationshipType is the relationship type from presentation.xml.
	ViewPropsRelationshipType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/viewProps"
	// ViewPropsRelationshipTarget is the target used from ppt/_rels/presentation.xml.rels.
	ViewPropsRelationshipTarget = "viewProps.xml"

	// TableStylesPartName is the package path of the table styles part.
	TableStylesPartName = "ppt/tableStyles.xml"
	// TableStylesContentType is the content type override for tableStyles.xml.
	TableStylesContentType = "application/vnd.openxmlformats-officedocument.presentationml.tableStyles+xml"
	// TableStylesRelationshipType is the relationship type from presentation.xml.
	TableStylesRelationshipType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/tableStyles"
	// TableStylesRelationshipTarget is the target used from ppt/_rels/presentation.xml.rels.
	TableStylesRelationshipTarget = "tableStyles.xml"

	// DefaultTableStyleGUID is "Medium Style 2 - Accent 1", the style a table
	// gets when it names none of its own.
	DefaultTableStyleGUID = "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"
)

// ViewProps renders ppt/viewProps.xml: the normal-view geometry, guides and
// zoom PowerPoint restores a deck into.
func ViewProps() string {
	return xmlHeader + `
<p:viewPr xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" ` +
		`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" ` +
		`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">
<p:normalViewPr>
<p:restoredLeft sz="15620"/>
<p:restoredTop sz="94660"/>
</p:normalViewPr>
<p:slideViewPr>
<p:cSldViewPr>
<p:cViewPr varScale="1">
<p:scale><a:sx n="76" d="100"/><a:sy n="76" d="100"/></p:scale>
<p:origin x="-1452" y="-96"/>
</p:cViewPr>
<p:guideLst/>
</p:cSldViewPr>
</p:slideViewPr>
<p:notesTextViewPr>
<p:cViewPr>
<p:scale><a:sx n="1" d="1"/><a:sy n="1" d="1"/></p:scale>
<p:origin x="0" y="0"/>
</p:cViewPr>
</p:notesTextViewPr>
<p:gridSpacing cx="76200" cy="76200"/>
</p:viewPr>`
}

// TableStyles renders ppt/tableStyles.xml, the part a table's
// <a:tableStyleId> resolves against.
func TableStyles() string {
	return xmlHeader + `
<a:tblStyleLst xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" def="` +
		DefaultTableStyleGUID + `"/>`
}

// defaultTextStyleLevels is how many <a:lvlNpPr> levels a presentation states.
const defaultTextStyleLevels = 9

// DefaultTextStyleXML renders <p:defaultTextStyle>, the nine outline levels a
// text box inherits from when it is not in a placeholder. Without it those
// boxes fall back to the application's defaults rather than the deck's.
func DefaultTextStyleXML() string {
	var b strings.Builder
	b.WriteString("\n<p:defaultTextStyle>")
	b.WriteString(`<a:defPPr><a:defRPr lang="en-US"/></a:defPPr>`)
	for level := 1; level <= defaultTextStyleLevels; level++ {
		//nolint:mnd // the level indents are the ones PowerPoint writes: 457200 EMU apart
		indent := 457200 * (level - 1)
		b.WriteString(`<a:lvl` + strconv.Itoa(level) + `pPr marL="` + strconv.Itoa(indent) +
			`" algn="l" defTabSz="914400" rtl="0" eaLnBrk="1" latinLnBrk="0" hangingPunct="1">` +
			`<a:defRPr sz="1800" kern="1200">` +
			`<a:solidFill><a:schemeClr val="tx1"/></a:solidFill>` +
			`<a:latin typeface="+mn-lt"/><a:ea typeface="+mn-ea"/><a:cs typeface="+mn-cs"/>` +
			`</a:defRPr></a:lvl` + strconv.Itoa(level) + `pPr>`)
	}
	b.WriteString("</p:defaultTextStyle>")
	return b.String()
}
