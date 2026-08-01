package pptxxml

import "strconv"

// Axis identifiers referenced by c:axId inside each plot element and declared by
// the matching c:catAx / c:valAx. The values are arbitrary but must be stable
// within a chart part and distinct between axis pairs, which is what lets a plot
// bind to the secondary value axis instead of the primary one.
const (
	primaryCatAxID   = 48650112
	primaryValAxID   = 48672768
	secondaryCatAxID = 48694272
	secondaryValAxID = 48716928
)

// chartAxisIDRefsXML renders the c:axId pair a plot element carries. Every plot
// in a chart part must reference exactly the axes it is drawn against; pointing
// two plots at the same pair is what made a combo chart share one value axis.
func chartAxisIDRefsXML(catAxID, valAxID int) string {
	return `
<c:axId val="` + strconv.Itoa(catAxID) + `"/>
<c:axId val="` + strconv.Itoa(valAxID) + `"/>`
}

// primaryAxisIDRefsXML is the c:axId pair for the primary category/value axes.
func primaryAxisIDRefsXML() string {
	return chartAxisIDRefsXML(primaryCatAxID, primaryValAxID)
}

// secondaryAxisIDRefsXML is the c:axId pair for the secondary axes, used by a
// plot that should be measured against its own value scale.
func secondaryAxisIDRefsXML() string {
	return chartAxisIDRefsXML(secondaryCatAxID, secondaryValAxID)
}
