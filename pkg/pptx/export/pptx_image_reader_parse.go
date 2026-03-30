package export

import (
	"encoding/xml"
	"strings"
)

const (
	imageCropScale = 100000.0
	imageRotScale  = 60000.0
)

type picRef struct {
	RelID        string
	X, Y         int64
	CX, CY       int64
	Rotation     float64
	CropLeft     float64
	CropRight    float64
	CropTop      float64
	CropBottom   float64
	FlipH        bool
	FlipV        bool
	Shadow       bool
	Reflection   bool
	AltText      string
	IsDecorative bool
}

type picReaderXML struct {
	NvPicPr struct {
		CNvPr struct {
			Descr *string `xml:"descr,attr"`
			Title *string `xml:"title,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvPicPr"`
	BlipFill struct {
		Blip struct {
			Embed string     `xml:"embed,attr"`
			Attrs []xml.Attr `xml:",any,attr"`
		} `xml:"blip"`
		SrcRect *struct {
			L *int `xml:"l,attr"`
			R *int `xml:"r,attr"`
			T *int `xml:"t,attr"`
			B *int `xml:"b,attr"`
		} `xml:"srcRect"`
	} `xml:"blipFill"`
	SpPr struct {
		Xfrm struct {
			Rot   *int       `xml:"rot,attr"`
			FlipH *bool      `xml:"flipH,attr"`
			FlipV *bool      `xml:"flipV,attr"`
			Attrs []xml.Attr `xml:",any,attr"`
			Off   struct {
				X int64 `xml:"x,attr"`
				Y int64 `xml:"y,attr"`
			} `xml:"off"`
			Ext struct {
				Cx int64 `xml:"cx,attr"`
				Cy int64 `xml:"cy,attr"`
			} `xml:"ext"`
		} `xml:"xfrm"`
		EffectLst *struct {
			OuterShdw  *struct{} `xml:"outerShdw"`
			Reflection *struct{} `xml:"reflection"`
		} `xml:"effectLst"`
	} `xml:"spPr"`
}

func parsePicElements(data []byte) []picRef {
	pics := make([]picRef, 0)
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	for {
		token, err := dec.Token()
		if err != nil {
			return pics
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "pic" {
			continue
		}
		var src picReaderXML
		if err := dec.DecodeElement(&src, &start); err != nil {
			continue
		}
		if pic, ok := picRefFromXML(&src); ok {
			pics = append(pics, pic)
		}
	}
}

func picRefFromXML(src *picReaderXML) (picRef, bool) {
	if src == nil {
		return picRef{}, false
	}
	ref := picRef{
		RelID: resolvePicRelID(src),
		X:     src.SpPr.Xfrm.Off.X,
		Y:     src.SpPr.Xfrm.Off.Y,
		CX:    src.SpPr.Xfrm.Ext.Cx,
		CY:    src.SpPr.Xfrm.Ext.Cy,
		FlipH: picFlipAttr(src.SpPr.Xfrm.FlipH, src.SpPr.Xfrm.Attrs, "flipH"),
		FlipV: picFlipAttr(src.SpPr.Xfrm.FlipV, src.SpPr.Xfrm.Attrs, "flipV"),
	}
	if ref.RelID == "" || ref.CX <= 0 || ref.CY <= 0 {
		return picRef{}, false
	}
	if src.SpPr.Xfrm.Rot != nil {
		ref.Rotation = float64(*src.SpPr.Xfrm.Rot) / imageRotScale
	}
	if src.BlipFill.SrcRect != nil {
		ref.CropLeft = cropFraction(src.BlipFill.SrcRect.L)
		ref.CropRight = cropFraction(src.BlipFill.SrcRect.R)
		ref.CropTop = cropFraction(src.BlipFill.SrcRect.T)
		ref.CropBottom = cropFraction(src.BlipFill.SrcRect.B)
	}
	if src.SpPr.EffectLst != nil {
		ref.Shadow = src.SpPr.EffectLst.OuterShdw != nil
		ref.Reflection = src.SpPr.EffectLst.Reflection != nil
	}
	ref.AltText, ref.IsDecorative = picAltText(src)
	return ref, true
}

func resolvePicRelID(src *picReaderXML) string {
	if src.BlipFill.Blip.Embed != "" {
		return strings.TrimSpace(src.BlipFill.Blip.Embed)
	}
	for _, attr := range src.BlipFill.Blip.Attrs {
		if attr.Name.Local == "embed" {
			return strings.TrimSpace(attr.Value)
		}
	}
	return ""
}

func picAltText(src *picReaderXML) (string, bool) {
	if src == nil || src.NvPicPr.CNvPr.Descr == nil {
		if src != nil && src.NvPicPr.CNvPr.Title != nil {
			return strings.TrimSpace(*src.NvPicPr.CNvPr.Title), false
		}
		return "", false
	}
	descr := strings.TrimSpace(*src.NvPicPr.CNvPr.Descr)
	if descr != "" {
		return descr, false
	}
	return "", true
}

func cropFraction(value *int) float64 {
	if value == nil {
		return 0
	}
	return float64(*value) / imageCropScale
}

func picFlipAttr(explicit *bool, attrs []xml.Attr, name string) bool {
	if explicit != nil {
		return *explicit
	}
	for _, attr := range attrs {
		if attr.Name.Local == name {
			value := strings.TrimSpace(strings.ToLower(attr.Value))
			return value == "1" || value == "true"
		}
	}
	return false
}
