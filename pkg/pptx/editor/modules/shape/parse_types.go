package shape

import common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"

const (
	rotationDegreeToOOXML = 60000.0
	gradientPositionScale = 1000.0
	ooxmlPercentScale     = 100000
	ptToEMU               = 12700
)

type solidFillXML struct {
	SrgbClr struct {
		Val   string `xml:"val,attr"`
		Alpha *struct {
			Val *int `xml:"val,attr"`
		} `xml:"alpha"`
	} `xml:"srgbClr"`
}

type runPropsXML struct {
	Bold          *bool        `xml:"b,attr"`
	Italic        *bool        `xml:"i,attr"`
	Underline     *string      `xml:"u,attr"`
	Strikethrough *string      `xml:"strike,attr"`
	Baseline      *string      `xml:"baseline,attr"`
	Caps          *string      `xml:"caps,attr"`
	SmallCaps     *string      `xml:"smCaps,attr"`
	Size          *int         `xml:"sz,attr"`
	SolidFill     solidFillXML `xml:"solidFill"`
	Highlight     solidFillXML `xml:"highlight"`
	Latin         *struct {
		Typeface string `xml:"typeface,attr"`
	} `xml:"latin"`
	Ln *struct {
		W         *int          `xml:"w,attr"`
		SolidFill *solidFillXML `xml:"solidFill"`
	} `xml:"ln"`
}

type gradientFillXML struct {
	Lin *struct {
		Ang *int `xml:"ang,attr"`
	} `xml:"lin"`
	GsLst struct {
		Gs []struct {
			Pos     *int `xml:"pos,attr"`
			SrgbClr *struct {
				Val string `xml:"val,attr"`
			} `xml:"srgbClr"`
		} `xml:"gs"`
	} `xml:"gsLst"`
}

type patternFillXML struct {
	Prst  *string `xml:"prst,attr"`
	FgClr *struct {
		SrgbClr *struct {
			Val string `xml:"val,attr"`
		} `xml:"srgbClr"`
	} `xml:"fgClr"`
	BgClr *struct {
		SrgbClr *struct {
			Val string `xml:"val,attr"`
		} `xml:"srgbClr"`
	} `xml:"bgClr"`
}

type bodyPrXML struct {
	LeftInset   *int      `xml:"lIns,attr"`
	RIns        *int      `xml:"rIns,attr"`
	TIns        *int      `xml:"tIns,attr"`
	BIns        *int      `xml:"bIns,attr"`
	Wrap        *string   `xml:"wrap,attr"`
	Anchor      *string   `xml:"anchor,attr"`
	Vert        *string   `xml:"vert,attr"`
	NumCol      *int      `xml:"numCol,attr"`
	Rot         *int      `xml:"rot,attr"`
	NoAutofit   *struct{} `xml:"noAutofit"`
	NormAutoFit *struct{} `xml:"normAutoFit"`
	SpAutoFit   *struct{} `xml:"spAutoFit"`
}

type spacingNodeXML struct {
	SpcPct *struct {
		Val *int `xml:"val,attr"`
	} `xml:"spcPct"`
	SpcPts *struct {
		Val *int `xml:"val,attr"`
	} `xml:"spcPts"`
}

type paragraphPropsXML struct {
	MarL   *int      `xml:"marL,attr"`
	Indent *int      `xml:"indent,attr"`
	Algn   *string   `xml:"algn,attr"`
	Lvl    *int      `xml:"lvl,attr"`
	BuNone *struct{} `xml:"buNone"`
	BuChar *struct {
		Char *string `xml:"char,attr"`
	} `xml:"buChar"`
	BuAutoNum *struct {
		Type *string `xml:"type,attr"`
	} `xml:"buAutoNum"`
	BuClr *struct {
		SrgbClr *struct {
			Val string `xml:"val,attr"`
		} `xml:"srgbClr"`
	} `xml:"buClr"`
	BuSzPct *struct {
		Val *int `xml:"val,attr"`
	} `xml:"buSzPct"`
	TabLst *struct {
		Tabs []struct {
			Pos *int `xml:"pos,attr"`
		} `xml:"tab"`
	} `xml:"tabLst"`
	LnSp   *spacingNodeXML `xml:"lnSp"`
	SpcBef *spacingNodeXML `xml:"spcBef"`
	SpcAft *spacingNodeXML `xml:"spcAft"`
}

type shapeXML struct {
	NvSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
		NvPr struct {
			Ph *struct {
				Idx  *int   `xml:"idx,attr"`
				Type string `xml:"type,attr"`
			} `xml:"ph"`
		} `xml:"nvPr"`
	} `xml:"nvSpPr"`
	NvPicPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
		NvPr struct {
			Ph *struct {
				Idx  *int   `xml:"idx,attr"`
				Type string `xml:"type,attr"`
			} `xml:"ph"`
		} `xml:"nvPr"`
	} `xml:"nvPicPr"`
	NvCxnSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
		CNvCxnSpPr struct {
			StCxn *struct {
				ID  *int `xml:"id,attr"`
				Idx *int `xml:"idx,attr"`
			} `xml:"stCxn"`
			EndCxn *struct {
				ID  *int `xml:"id,attr"`
				Idx *int `xml:"idx,attr"`
			} `xml:"endCxn"`
		} `xml:"cNvCxnSpPr"`
	} `xml:"nvCxnSpPr"`
	NvGrpSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvGrpSpPr"`
	NvGraphicFramePr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvGraphicFramePr"`
	Xfrm struct {
		Rot   *int    `xml:"rot,attr"`
		FlipH *string `xml:"flipH,attr"`
		FlipV *string `xml:"flipV,attr"`
		Off   struct {
			X int `xml:"x,attr"`
			Y int `xml:"y,attr"`
		} `xml:"off"`
		Ext struct {
			Cx int `xml:"cx,attr"`
			Cy int `xml:"cy,attr"`
		} `xml:"ext"`
	} `xml:"xfrm"`
	SpPr struct {
		PrstGeom *struct {
			Prst  string `xml:"prst,attr"`
			AvLst *struct {
				Gd []struct {
					Name string `xml:"name,attr"`
					Fmla string `xml:"fmla,attr"`
				} `xml:"gd"`
			} `xml:"avLst"`
		} `xml:"prstGeom"`
		NoFill    *struct{}        `xml:"noFill"`
		SolidFill *solidFillXML    `xml:"solidFill"`
		GradFill  *gradientFillXML `xml:"gradFill"`
		PattFill  *patternFillXML  `xml:"pattFill"`
		Ln        *struct {
			W         *int          `xml:"w,attr"`
			SolidFill *solidFillXML `xml:"solidFill"`
			PrstDash  *struct {
				Val string `xml:"val,attr"`
			} `xml:"prstDash"`
			HeadEnd *struct {
				Type *string `xml:"type,attr"`
				W    *string `xml:"w,attr"`
				Len  *string `xml:"len,attr"`
			} `xml:"headEnd"`
			TailEnd *struct {
				Type *string `xml:"type,attr"`
				W    *string `xml:"w,attr"`
				Len  *string `xml:"len,attr"`
			} `xml:"tailEnd"`
		} `xml:"ln"`
		EffectLst *struct {
			OuterShdw *struct {
				BlurRad *int `xml:"blurRad,attr"`
				Dist    *int `xml:"dist,attr"`
				Dir     *int `xml:"dir,attr"`
				SrgbClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"outerShdw"`
			Glow *struct {
				Rad     *int `xml:"rad,attr"`
				SrgbClr *struct {
					Val string `xml:"val,attr"`
				} `xml:"srgbClr"`
			} `xml:"glow"`
			Blur *struct {
				Rad *int `xml:"rad,attr"`
			} `xml:"blur"`
			SoftEdge *struct {
				Rad *int `xml:"rad,attr"`
			} `xml:"softEdge"`
			Reflection *struct {
				BlurRad *int `xml:"blurRad,attr"`
				Dist    *int `xml:"dist,attr"`
			} `xml:"reflection"`
		} `xml:"effectLst"`
		Xfrm struct {
			Rot   *int    `xml:"rot,attr"`
			FlipH *string `xml:"flipH,attr"`
			FlipV *string `xml:"flipV,attr"`
			Off   struct {
				X int `xml:"x,attr"`
				Y int `xml:"y,attr"`
			} `xml:"off"`
			Ext struct {
				Cx int `xml:"cx,attr"`
				Cy int `xml:"cy,attr"`
			} `xml:"ext"`
		} `xml:"xfrm"`
	} `xml:"spPr"`
	GrpSpPr struct {
		Xfrm struct {
			Off struct {
				X int `xml:"x,attr"`
				Y int `xml:"y,attr"`
			} `xml:"off"`
			Ext struct {
				Cx int `xml:"cx,attr"`
				Cy int `xml:"cy,attr"`
			} `xml:"ext"`
		} `xml:"xfrm"`
	} `xml:"grpSpPr"`
	TxBody struct {
		BodyPr *bodyPrXML `xml:"bodyPr"`
		P      []struct {
			PPr *paragraphPropsXML `xml:"pPr"`
			R   []struct {
				RPr *runPropsXML `xml:"rPr"`
				T   string       `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

type ParsedShapeProperties struct {
	ID         int
	Name       string
	Type       string
	Text       string
	Runs       []common.TextRun
	Paragraphs []common.ShapeTextParagraph
	TextFrame  *common.TextFrame
	Paragraph  *common.Paragraph
	Fill       *common.ShapeFill
	Line       *common.ShapeLine
	Shadow     *common.ShapeShadow
	Glow       *common.ShapeGlow
	Blur       *common.ShapeBlur
	SoftEdge   *common.ShapeSoftEdge
	Reflection *common.ShapeReflection
	Rotation   *float64
	Connector  *common.ConnectorInfo
	X, Y       int
	W, H       int
	PhIndex    int
	PhType     string

	Adjustments []common.ShapeAdjustment
}
