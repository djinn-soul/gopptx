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
	ShapeID      int
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

	InnerShadow       bool
	Glow              bool
	SoftEdges         bool
	Blur              bool
	GlowRadiusEmu     int
	SoftEdgeRadiusEmu int
	BlurRadiusEmu     int
}

// decorativeExtURI is the OOXML extension URI that Office uses to mark a
// picture as intentionally decorative (accessibility flag).
// Source: ECMA-376 / Microsoft Office Open XML extensions.
const decorativeExtURI = "{C183D7F6-72BE-476a-BEBA-66C5E2CAE503}"

type picCNvPrXML struct {
	ID     *int    `xml:"id,attr"`
	Descr  *string `xml:"descr,attr"`
	Title  *string `xml:"title,attr"`
	ExtLst *struct {
		Exts []struct {
			URI           string `xml:"uri,attr"`
			DecorativeExt *struct {
				Val *bool `xml:"val,attr"`
			} `xml:"decorative"`
		} `xml:"ext"`
	} `xml:"extLst"`
}

// picCNvPrIsDecorative returns true only when the cNvPr element carries the
// explicit OOXML decorative-image extension with val="1" (or val="true").
// An absent or empty descr attribute is NOT sufficient on its own.
func picCNvPrIsDecorative(c picCNvPrXML) bool {
	if c.ExtLst == nil {
		return false
	}
	for _, ext := range c.ExtLst.Exts {
		if ext.URI == decorativeExtURI && ext.DecorativeExt != nil {
			return ext.DecorativeExt.Val == nil || *ext.DecorativeExt.Val
		}
	}
	return false
}

type picReaderXML struct {
	NvPicPr struct {
		CNvPr picCNvPrXML `xml:"cNvPr"`
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
		EffectLst *picEffectListXML `xml:"effectLst"`
	} `xml:"spPr"`
}

// picEffectListXML is the `<a:effectLst>` of a picture. Only the presence of
// each effect and the radii the renderer needs are decoded.
type picEffectListXML struct {
	OuterShdw  *struct{}           `xml:"outerShdw"`
	Reflection *struct{}           `xml:"reflection"`
	InnerShdw  *struct{}           `xml:"innerShdw"`
	Glow       *picEffectRadiusXML `xml:"glow"`
	SoftEdge   *picEffectRadiusXML `xml:"softEdge"`
	Blur       *picEffectRadiusXML `xml:"blur"`
}

// picEffectRadiusXML is an effect carrying a single radius attribute.
type picEffectRadiusXML struct {
	Rad *int `xml:"rad,attr"`
}

// picGroupXML is one level of the shape tree as far as pictures are concerned:
// the pictures it holds directly, the groups nested inside it, and the
// transform that maps its children onto the slide.
type picGroupXML struct {
	GrpSpPr struct {
		Xfrm *groupXfrmXML `xml:"xfrm"`
	} `xml:"grpSpPr"`
	Pics   []picReaderXML `xml:"pic"`
	Groups []picGroupXML  `xml:"grpSp"`
}

// groupXfrmXML is a group's <a:xfrm>: where the group sits on the slide, and
// the coordinate space its children state their own geometry in.
type groupXfrmXML struct {
	Off struct {
		X int64 `xml:"x,attr"`
		Y int64 `xml:"y,attr"`
	} `xml:"off"`
	Ext struct {
		Cx int64 `xml:"cx,attr"`
		Cy int64 `xml:"cy,attr"`
	} `xml:"ext"`
	ChOff struct {
		X int64 `xml:"x,attr"`
		Y int64 `xml:"y,attr"`
	} `xml:"chOff"`
	ChExt struct {
		Cx int64 `xml:"cx,attr"`
		Cy int64 `xml:"cy,attr"`
	} `xml:"chExt"`
}

// groupTransform maps a point and a size out of a group's child space onto the
// slide. The identity transform is the zero value.
type groupTransform struct {
	offX, offY     int64
	scaleX, scaleY float64
}

func identityGroupTransform() groupTransform {
	return groupTransform{scaleX: 1, scaleY: 1}
}

// compose returns this transform followed by the one a nested group states.
func (t groupTransform) compose(x *groupXfrmXML) groupTransform {
	if x == nil || x.ChExt.Cx <= 0 || x.ChExt.Cy <= 0 || x.Ext.Cx <= 0 || x.Ext.Cy <= 0 {
		return t
	}
	scaleX := float64(x.Ext.Cx) / float64(x.ChExt.Cx)
	scaleY := float64(x.Ext.Cy) / float64(x.ChExt.Cy)
	// A child point p maps to off + (p - chOff) * scale; folding that into the
	// parent transform keeps nested groups a single multiply-and-add.
	return groupTransform{
		offX:   t.offX + int64(float64(x.Off.X)*t.scaleX) - int64(float64(x.ChOff.X)*t.scaleX*scaleX),
		offY:   t.offY + int64(float64(x.Off.Y)*t.scaleY) - int64(float64(x.ChOff.Y)*t.scaleY*scaleY),
		scaleX: t.scaleX * scaleX,
		scaleY: t.scaleY * scaleY,
	}
}

// apply maps a picture's geometry onto the slide.
func (t groupTransform) apply(ref *picRef) {
	if t.scaleX == 0 && t.scaleY == 0 {
		return
	}
	ref.X = t.offX + int64(float64(ref.X)*t.scaleX)
	ref.Y = t.offY + int64(float64(ref.Y)*t.scaleY)
	ref.CX = int64(float64(ref.CX) * t.scaleX)
	ref.CY = int64(float64(ref.CY) * t.scaleY)
}

// parsePicElements collects every picture on the slide, including those inside
// groups.
//
// A grouped picture states its geometry in the group's child space, so it has
// to be mapped onto the slide the way PowerPoint does — a group the user
// resized has a child space that no longer matches its own size, and the
// picture inside it was landing at the wrong place and the wrong size.
func parsePicElements(data []byte) []picRef {
	var doc struct {
		Tree picGroupXML `xml:"cSld>spTree"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil
	}
	pics := make([]picRef, 0)
	collectGroupPics(doc.Tree, identityGroupTransform(), &pics)
	return pics
}

func collectGroupPics(node picGroupXML, transform groupTransform, out *[]picRef) {
	for i := range node.Pics {
		pic, ok := picRefFromXML(&node.Pics[i])
		if !ok {
			continue
		}
		transform.apply(&pic)
		*out = append(*out, pic)
	}
	for _, child := range node.Groups {
		collectGroupPics(child, transform.compose(child.GrpSpPr.Xfrm), out)
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
	// The shape id ties this picture back to its position in the slide's shape
	// tree, which is what decides whether it paints above or below a shape.
	if src.NvPicPr.CNvPr.ID != nil {
		ref.ShapeID = *src.NvPicPr.CNvPr.ID
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
	applyPicEffects(&ref, src.SpPr.EffectLst)
	ref.AltText, ref.IsDecorative = picAltText(src)
	return ref, true
}

// applyPicEffects copies the decoded `<a:effectLst>` onto the picture. A nil
// list leaves every effect off.
func applyPicEffects(ref *picRef, effects *picEffectListXML) {
	if effects == nil {
		return
	}
	ref.Shadow = effects.OuterShdw != nil
	ref.Reflection = effects.Reflection != nil
	ref.InnerShadow = effects.InnerShdw != nil
	ref.Glow = effects.Glow != nil
	ref.SoftEdges = effects.SoftEdge != nil
	ref.GlowRadiusEmu = picEffectRadius(effects.Glow)
	ref.SoftEdgeRadiusEmu = picEffectRadius(effects.SoftEdge)
	if effects.Blur != nil {
		ref.Blur = true
		ref.BlurRadiusEmu = picEffectRadius(effects.Blur)
	}
}

// picEffectRadius is the effect's radius, or 0 when the effect or its
// attribute is absent.
func picEffectRadius(effect *picEffectRadiusXML) int {
	if effect == nil || effect.Rad == nil {
		return 0
	}
	return *effect.Rad
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
	if src == nil {
		return "", false
	}
	c := src.NvPicPr.CNvPr
	isDecorative := picCNvPrIsDecorative(c)
	if c.Descr != nil {
		if descr := strings.TrimSpace(*c.Descr); descr != "" {
			return descr, isDecorative
		}
	}
	if !isDecorative && c.Title != nil {
		if title := strings.TrimSpace(*c.Title); title != "" {
			return title, false
		}
	}
	return "", isDecorative
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
