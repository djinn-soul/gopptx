package pptxxml

import (
	"fmt"
	"strings"
)

const (
	dgmDataNS     = "http://schemas.openxmlformats.org/drawingml/2006/diagram"
	dgmDrawingNS  = "http://schemas.microsoft.com/office/drawing/2008/diagram"
	drawingMainNS = "http://schemas.openxmlformats.org/drawingml/2006/main"
)

// SmartArtDataXML renders the dgm:dataModel XML (dataX.xml).
//
// This is the core data file containing the point list (nodes) and
// connection list (parent→child links). PowerPoint reads this to
// understand the diagram's logical structure.
func SmartArtDataXML(spec SmartArtSpec) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dgm:dataModel xmlns:dgm="` + dgmDataNS + `" xmlns:a="` + drawingMainNS + `">
  <dgm:ptLst>
    <dgm:pt modelId="0" type="doc"/>`)

	var nextID int
	var cxns []string
	for _, node := range spec.Nodes {
		writeNodePts(&b, &cxns, node, "0", &nextID)
	}

	b.WriteString(`
  </dgm:ptLst>
  <dgm:cxnLst>`)
	for _, cxn := range cxns {
		b.WriteString(cxn)
	}
	b.WriteString(`
  </dgm:cxnLst>
  <dgm:bg/>
  <dgm:whole/>
</dgm:dataModel>`)

	return b.String()
}

func writeNodePts(b *strings.Builder, cxns *[]string, node SmartArtNodeSpec, parentID string, nextID *int) {
	*nextID++
	myID := fmt.Sprintf("%d", *nextID)

	b.WriteString(fmt.Sprintf(`
    <dgm:pt modelId="%s" type="node">
      <dgm:prSet/>
      <dgm:spPr/>
      <dgm:t>
        <a:bodyPr/>
        <a:lstStyle/>
        <a:p><a:r><a:t>%s</a:t></a:r></a:p>
      </dgm:t>
    </dgm:pt>`, myID, Escape(node.Text)))

	// Parent→child connection.
	*nextID++
	cxnID := *nextID
	*cxns = append(*cxns, fmt.Sprintf(`
    <dgm:cxn modelId="%d" srcId="%s" destId="%s" type="parOf"/>`,
		cxnID, parentID, myID))

	for _, child := range node.Children {
		writeNodePts(b, cxns, child, myID, nextID)
	}
}

// SmartArtLayoutXML renders a minimal dgm:layoutDef (layoutX.xml).
//
// PowerPoint uses the layout URI from the data to resolve the actual
// layout algorithm internally; we emit a stub.
func SmartArtLayoutXML(layoutURI, category string) string {
	if category == "" {
		category = "list"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dgm:layoutDef xmlns:dgm="%s" uniqueId="%s" minVer="12.0">
  <dgm:title val=""/>
  <dgm:desc val=""/>
  <dgm:catLst>
    <dgm:cat type="%s" pri="100"/>
  </dgm:catLst>
</dgm:layoutDef>`, dgmDataNS, Escape(layoutURI), category)
}

// SmartArtColorsXML renders a minimal dgm:colorsDef (colorsX.xml).
func SmartArtColorsXML(colorStyleID string) string {
	csType := colorStyleID
	if csType == "" {
		csType = "urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dgm:colorsDef xmlns:dgm="%s" xmlns:a="%s"
  uniqueId="%s" minVer="12.0"/>`, dgmDataNS, drawingMainNS, Escape(csType))
}

// SmartArtStyleXML renders a minimal dgm:styleDef (quickStyleX.xml).
func SmartArtStyleXML(quickStyleID string) string {
	qsType := quickStyleID
	if qsType == "" {
		qsType = "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dgm:styleDef xmlns:dgm="%s" xmlns:a="%s"
  uniqueId="%s" minVer="12.0"/>`, dgmDataNS, drawingMainNS, Escape(qsType))
}

// SmartArtDrawingXML renders a minimal dsp:drawing (drawingX.xml).
//
// PowerPoint regenerates this part when the file is opened, so we
// emit an empty placeholder.
func SmartArtDrawingXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dsp:drawing xmlns:dsp="` + dgmDrawingNS + `">
  <dsp:spTree>
    <dsp:nvGrpSpPr>
      <dsp:cNvPr id="0" name=""/>
      <dsp:cNvGrpSpPr/>
    </dsp:nvGrpSpPr>
    <dsp:grpSpPr/>
  </dsp:spTree>
</dsp:drawing>`
}
