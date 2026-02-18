package pptxxml

import (
	"fmt"
	"strings"
)

const (
	dgmDataNS     = "http://schemas.openxmlformats.org/drawingml/2006/diagram"
	dgmDrawingNS  = "http://schemas.microsoft.com/office/drawing/2008/diagram"
	drawingMainNS = "http://schemas.openxmlformats.org/drawingml/2006/main"
	dgmExtURI     = "http://schemas.microsoft.com/office/drawing/2008/diagram"
)

type smartArtNodePoint struct {
	modelID     string
	presModelID string
	parTransID  string
	sibTransID  string
	cxnID       string
	text        string
	order       int
}

type smartArtConnection struct {
	modelID    string
	srcID      string
	destID     string
	srcOrd     int
	destOrd    int
	kind       string
	presID     string
	parTransID string
	sibTransID string
}

const smartArtDocModelID = "{00000000-0000-0000-0000-000000000000}"

// SmartArtDataXML renders the dgm:dataModel XML (dataX.xml).
//
// This is the core data file containing the point list (nodes) and
// connection list (parent→child links). PowerPoint reads this to
// understand the diagram's logical structure.
func SmartArtDataXML(spec SmartArtSpec) string {
	return renderSmartArtDataFromTemplate(spec)
}

func buildSmartArtModel(nodes []SmartArtNodeSpec) ([]smartArtNodePoint, []smartArtConnection, string, int) {
	points := make([]smartArtNodePoint, 0)
	cxns := make([]smartArtConnection, 0)
	nextID := 0
	for _, node := range nodes {
		appendSmartArtModel(node, smartArtDocModelID, &nextID, &points, &cxns)
	}
	nextID++
	presDocID := smartArtModelID(nextID)
	for i := range points {
		nextID++
		points[i].presModelID = smartArtModelID(nextID)
	}
	return points, cxns, presDocID, nextID
}

func appendSmartArtModel(
	node SmartArtNodeSpec,
	parentID string,
	nextID *int,
	points *[]smartArtNodePoint,
	cxns *[]smartArtConnection,
) {
	*nextID++
	myID := smartArtModelID(*nextID)
	*nextID++
	cxnID := smartArtModelID(*nextID)
	*nextID++
	parTransID := smartArtModelID(*nextID)
	*nextID++
	sibTransID := smartArtModelID(*nextID)
	*points = append(*points, smartArtNodePoint{
		modelID:    myID,
		cxnID:      cxnID,
		parTransID: parTransID,
		sibTransID: sibTransID,
		text:       node.Text,
		order:      countChildrenForParent(parentID, *cxns),
	})

	// Parent→child connection.
	*cxns = append(*cxns, smartArtConnection{
		modelID:    cxnID,
		srcID:      parentID,
		destID:     myID,
		srcOrd:     countChildrenForParent(parentID, *cxns),
		destOrd:    0,
		parTransID: parTransID,
		sibTransID: sibTransID,
	})

	for _, child := range node.Children {
		appendSmartArtModel(child, myID, nextID, points, cxns)
	}
}

func countChildrenForParent(parentID string, cxns []smartArtConnection) int {
	count := 0
	for i := range cxns {
		if cxns[i].srcID == parentID {
			count++
		}
	}
	return count
}

// SmartArtLayoutXML renders a minimal dgm:layoutDef (layoutX.xml).
//
// PowerPoint uses the layout URI from the data to resolve the actual
// layout algorithm internally; we emit a stub.
func SmartArtLayoutXML(layoutURI, category string) string {
	return renderSmartArtLayoutFromTemplate(layoutURI)

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
  <dgm:layoutNode name="diagram">
    <dgm:alg type="lin"/>
    <dgm:shape type="rect"/>
    <dgm:presOf/>
  </dgm:layoutNode>
</dgm:layoutDef>`, dgmDataNS, Escape(layoutURI), category)
}

// SmartArtColorsXML renders a minimal dgm:colorsDef (colorsX.xml).
func SmartArtColorsXML(colorStyleID string) string {
	return renderSmartArtColorsFromTemplate(colorStyleID)

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
	return renderSmartArtStyleFromTemplate(quickStyleID)

	qsType := quickStyleID
	if qsType == "" {
		qsType = "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<dgm:styleDef xmlns:dgm="%s" xmlns:a="%s"
  uniqueId="%s" minVer="12.0">
  <dgm:title val=""/>
  <dgm:desc val=""/>
  <dgm:catLst>
    <dgm:cat type="simple" pri="100"/>
  </dgm:catLst>
  <dgm:styleLbl name="node0">
    <dgm:scene3d>
      <a:camera prst="orthographicFront"/>
      <a:lightRig rig="threePt" dir="t"/>
    </dgm:scene3d>
    <dgm:sp3d/>
    <dgm:txPr/>
    <dgm:style>
      <a:lnRef idx="2"><a:scrgbClr r="0" g="0" b="0"/></a:lnRef>
      <a:fillRef idx="1"><a:scrgbClr r="0" g="0" b="0"/></a:fillRef>
      <a:effectRef idx="0"><a:scrgbClr r="0" g="0" b="0"/></a:effectRef>
      <a:fontRef idx="minor"><a:schemeClr val="lt1"/></a:fontRef>
    </dgm:style>
  </dgm:styleLbl>
</dgm:styleDef>`, dgmDataNS, drawingMainNS, Escape(qsType))
}

// SmartArtDrawingXML renders a minimal dsp:drawing (drawingX.xml).
//
// PowerPoint regenerates this part when the file is opened, so we
// emit an empty placeholder.
func SmartArtDrawingXML(spec SmartArtSpec) string {
	return renderSmartArtDrawingFromTemplate(spec)
}

func renderSmartArtDrawingNode(pt smartArtNodePoint, idx int) string {
	const (
		nodeWidth  = int64(4572000)
		nodeHeight = int64(952500)
		nodeStartX = int64(0)
		nodeStartY = int64(482600)
		nodeGapY   = int64(317500)
	)
	y := nodeStartY + int64(idx)*(nodeHeight+nodeGapY)
	return fmt.Sprintf(`
    <dsp:sp modelId="%s">
      <dsp:nvSpPr>
        <dsp:cNvPr id="0" name=""/>
        <dsp:cNvSpPr/>
      </dsp:nvSpPr>
      <dsp:spPr>
        <a:xfrm>
          <a:off x="%d" y="%d"/>
          <a:ext cx="%d" cy="%d"/>
        </a:xfrm>
        <a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
        <a:solidFill>
          <a:schemeClr val="accent1"/>
        </a:solidFill>
        <a:ln w="19050" cap="flat" cmpd="sng" algn="ctr">
          <a:solidFill>
            <a:schemeClr val="lt1"/>
          </a:solidFill>
          <a:prstDash val="solid"/>
          <a:miter lim="800000"/>
        </a:ln>
        <a:effectLst/>
      </dsp:spPr>
      <dsp:style>
        <a:lnRef idx="2"><a:scrgbClr r="0" g="0" b="0"/></a:lnRef>
        <a:fillRef idx="1"><a:scrgbClr r="0" g="0" b="0"/></a:fillRef>
        <a:effectRef idx="0"><a:scrgbClr r="0" g="0" b="0"/></a:effectRef>
        <a:fontRef idx="minor"><a:schemeClr val="lt1"/></a:fontRef>
      </dsp:style>
      <dsp:txBody>
        <a:bodyPr spcFirstLastPara="0" vert="horz" wrap="square" lIns="247650" tIns="247650" rIns="247650" bIns="247650" numCol="1" spcCol="1270" anchor="ctr" anchorCtr="0">
          <a:noAutofit/>
        </a:bodyPr>
        <a:lstStyle/>
        <a:p>
          <a:pPr marL="0" lvl="0" indent="0" algn="ctr" defTabSz="2889250">
            <a:lnSpc><a:spcPct val="90000"/></a:lnSpc>
            <a:spcBef><a:spcPct val="0"/></a:spcBef>
            <a:spcAft><a:spcPct val="35000"/></a:spcAft>
            <a:buNone/>
          </a:pPr>
          <a:r>
            <a:rPr lang="en-US" sz="6500" kern="1200"/>
            <a:t>%s</a:t>
          </a:r>
          <a:endParaRPr lang="en-US" sz="6500" kern="1200"/>
        </a:p>
      </dsp:txBody>
      <dsp:txXfrm>
        <a:off x="%d" y="%d"/>
        <a:ext cx="%d" cy="%d"/>
      </dsp:txXfrm>
    </dsp:sp>`,
		Escape(pt.presModelID), nodeStartX, y, nodeWidth, nodeHeight, Escape(pt.text),
		nodeStartX, y, nodeWidth, nodeHeight)
}

func defaultColorStyleID(id string) string {
	if id != "" {
		return id
	}
	return "urn:microsoft.com/office/officeart/2005/8/colors/accent1_2"
}

func defaultQuickStyleID(id string) string {
	if id != "" {
		return id
	}
	return "urn:microsoft.com/office/officeart/2005/8/quickstyle/simple1"
}

func categoryFromLayoutURI(uri string) string {
	uri = strings.ToLower(uri)
	switch {
	case strings.Contains(uri, "process"):
		return "process"
	case strings.Contains(uri, "cycle"):
		return "cycle"
	case strings.Contains(uri, "hierarchy"), strings.Contains(uri, "orgchart"):
		return "hierarchy"
	case strings.Contains(uri, "venn"), strings.Contains(uri, "radial"), strings.Contains(uri, "target"):
		return "relationship"
	case strings.Contains(uri, "matrix"):
		return "matrix"
	case strings.Contains(uri, "pyramid"):
		return "pyramid"
	case strings.Contains(uri, "picture"):
		return "picture"
	default:
		return "list"
	}
}

func smartArtModelID(id int) string {
	return fmt.Sprintf("{00000000-0000-0000-0000-%012X}", id)
}
