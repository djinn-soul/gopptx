package shape

import "encoding/xml"

// srcRectScale converts <a:srcRect> percent-thousandths into a fraction.
const srcRectScale = 100000.0

type custGeomXML struct {
	PathLst struct {
		Path []pathXML `xml:"path"`
	} `xml:"pathLst"`
}

type pathXML struct {
	W      int     `xml:"w,attr"`
	H      int     `xml:"h,attr"`
	Fill   string  `xml:"fill,attr"`
	Stroke *string `xml:"stroke,attr"`
	// Segments is a catch-all so moveTo/lnTo/cubicBezTo/quadBezTo/arcTo/close
	// keep their document order, which is the whole meaning of a path.
	Segments []pathSegmentXML `xml:",any"`
}

type pathSegmentXML struct {
	XMLName xml.Name
	Pt      []ptXML `xml:"pt"`
	WR      string  `xml:"wR,attr"`
	HR      string  `xml:"hR,attr"`
	StAng   string  `xml:"stAng,attr"`
	SwAng   string  `xml:"swAng,attr"`
}

// ptXML keeps coordinates as strings because DrawingML allows a guide formula
// (for example x="wd2") wherever a number is allowed.
type ptXML struct {
	X string `xml:"x,attr"`
	Y string `xml:"y,attr"`
}

type blipFillXML struct {
	Blip *struct {
		Embed string `xml:"embed,attr"`
		Link  string `xml:"link,attr"`
	} `xml:"blip"`
	SrcRect *struct {
		L *int `xml:"l,attr"`
		T *int `xml:"t,attr"`
		R *int `xml:"r,attr"`
		B *int `xml:"b,attr"`
	} `xml:"srcRect"`
	Stretch *struct{} `xml:"stretch"`
	Tile    *struct{} `xml:"tile"`
}
