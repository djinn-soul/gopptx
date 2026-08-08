package smartart

import "encoding/xml"

// The XML shapes below mirror only what the renderer reads. Field tags name
// local elements without a namespace, so they match both the dsp: and a:
// namespaces the drawing part mixes.

type dspDrawing struct {
	XMLName xml.Name  `xml:"drawing"`
	SpTree  dspSpTree `xml:"spTree"`
}

type dspSpTree struct {
	Sp []dspSp `xml:"sp"`
}

type dspSp struct {
	ModelID string   `xml:"modelId,attr"`
	SpPr    dspSpPr  `xml:"spPr"`
	TxBody  *aTxBody `xml:"txBody"`
}

type dspSpPr struct {
	Xfrm      *aXfrm        `xml:"xfrm"`
	PrstGeom  *aPrstGeom    `xml:"prstGeom"`
	SolidFill *aColorChoice `xml:"solidFill"`
	NoFill    *struct{}     `xml:"noFill"`
	Ln        *aLn          `xml:"ln"`
}

type aXfrm struct {
	Rot   int    `xml:"rot,attr"`
	FlipH string `xml:"flipH,attr"`
	FlipV string `xml:"flipV,attr"`
	Off   struct {
		X int64 `xml:"x,attr"`
		Y int64 `xml:"y,attr"`
	} `xml:"off"`
	Ext struct {
		CX int64 `xml:"cx,attr"`
		CY int64 `xml:"cy,attr"`
	} `xml:"ext"`
}

type aPrstGeom struct {
	Prst string `xml:"prst,attr"`
	Gds  []aGd  `xml:"avLst>gd"`
}

type aGd struct {
	Name string `xml:"name,attr"`
	Fmla string `xml:"fmla,attr"`
}

type aLn struct {
	W         int64         `xml:"w,attr"`
	SolidFill *aColorChoice `xml:"solidFill"`
	NoFill    *struct{}     `xml:"noFill"`
}

type aColorChoice struct {
	SrgbClr   *aColorNode `xml:"srgbClr"`
	SchemeClr *aColorNode `xml:"schemeClr"`
}

type aColorNode struct {
	Val    string `xml:"val,attr"`
	HueOff *aVal  `xml:"hueOff"`
	SatOff *aVal  `xml:"satOff"`
	LumOff *aVal  `xml:"lumOff"`
	LumMod *aVal  `xml:"lumMod"`
	SatMod *aVal  `xml:"satMod"`
	Shade  *aVal  `xml:"shade"`
	Tint   *aVal  `xml:"tint"`
	Alpha  *aVal  `xml:"alpha"`
}

type aVal struct {
	Val string `xml:"val,attr"`
}

type aTxBody struct {
	BodyPr aBodyPr `xml:"bodyPr"`
	Ps     []aP    `xml:"p"`
}

type aBodyPr struct {
	Anchor      string `xml:"anchor,attr"`
	LeftInset   *int64 `xml:"lIns,attr"`
	TopInset    *int64 `xml:"tIns,attr"`
	RightInset  *int64 `xml:"rIns,attr"`
	BottomInset *int64 `xml:"bIns,attr"`
}

type aP struct {
	PPr *aPPr `xml:"pPr"`
	Rs  []aR  `xml:"r"`
}

type aPPr struct {
	Algn string `xml:"algn,attr"`
}

type aR struct {
	RPr *aRPr  `xml:"rPr"`
	T   string `xml:"t"`
}

type aRPr struct {
	Sz        int           `xml:"sz,attr"`
	B         string        `xml:"b,attr"`
	I         string        `xml:"i,attr"`
	SolidFill *aColorChoice `xml:"solidFill"`
	Latin     *struct {
		Typeface string `xml:"typeface,attr"`
	} `xml:"latin"`
}
