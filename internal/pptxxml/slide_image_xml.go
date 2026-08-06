package pptxxml

import (
	"strconv"
	"strings"
)

// ImageRef describes one image reference in a slide.
type ImageRef struct {
	RelID        string
	Name         string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	Rotation     int64 // 60000ths of a degree
	FlipH        bool
	FlipV        bool
	Crop         *ImageCropRef
	Shadow       bool
	Reflection   bool
	AltText      string
	IsDecorative bool

	InnerShadow    bool
	Glow           bool
	SoftEdges      bool
	GlowSpec       *ShapeGlowSpec
	BlurSpec       *ShapeBlurSpec
	SoftEdgeSpec   *ShapeSoftEdgeSpec
	ReflectionSpec *ShapeReflectionSpec
}

// ImageCropRef defines cropping percentages (0-100000 range for OOXML).
type ImageCropRef struct {
	Left   int64
	Right  int64
	Top    int64
	Bottom int64
}

const (
	defaultShadowBlurRad    = 40000
	defaultShadowDist       = 20000
	defaultShadowDir        = 5400000
	defaultShadowAlpha      = 40000
	defaultReflectionBlur   = 6350
	defaultReflectionStA    = 50000
	defaultReflectionEndA   = 300
	defaultReflectionEndPos = 35000
	defaultGlowRadius       = 63500
	defaultGlowColor        = "4472C4"
	defaultGlowAlpha        = 35000
	defaultSoftEdgeRadius   = 38100
)

func imageShape(image ImageRef, shapeID int) string {
	name := image.Name
	if name == "" {
		name = "Picture"
	}

	descrAttr := imageDescriptionAttr(image)
	srcRect := imageSrcRectXML(image.Crop)
	effectsXML := imageEffectsXML(image)
	xfrmAttrs := imageTransformAttrs(image.Rotation, image.FlipH, image.FlipV)

	return `
<p:pic>
<p:nvPicPr>
<p:cNvPr id="` + strconv.Itoa(shapeID) + `" name="` + Escape(name) + `"` + descrAttr + `/>
<p:cNvPicPr><a:picLocks noChangeAspect="1"/></p:cNvPicPr>
<p:nvPr/>
</p:nvPicPr>
<p:blipFill>
<a:blip r:embed="` + FastEscapeRID(image.RelID) + `"/>
` + srcRect + `
<a:stretch><a:fillRect/></a:stretch>
</p:blipFill>
<p:spPr>
<a:xfrm` + xfrmAttrs + `>
<a:off x="` + strconv.FormatInt(image.X, 10) + `" y="` + strconv.FormatInt(image.Y, 10) + `"/>
<a:ext cx="` + strconv.FormatInt(image.CX, 10) + `" cy="` + strconv.FormatInt(image.CY, 10) + `"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
` + effectsXML + `
</p:spPr>
</p:pic>`
}

func imageDescriptionAttr(image ImageRef) string {
	if image.IsDecorative {
		return shapeDescrAttrEmpty
	}
	if image.AltText != "" {
		escaped := Escape(image.AltText)
		return ` descr="` + escaped + `" title="` + escaped + `"`
	}
	return shapeDescrAttrEmpty
}

func imageTransformAttrs(rot int64, flipH, flipV bool) string {
	attrs := ""
	if rot != 0 {
		attrs += ` rot="` + strconv.FormatInt(rot, 10) + `"`
	}
	if flipH {
		attrs += ` flipH="1"`
	}
	if flipV {
		attrs += ` flipV="1"`
	}
	return attrs
}

func imageSrcRectXML(c *ImageCropRef) string {
	if c == nil {
		return ""
	}
	attrs := ""
	if c.Left != 0 {
		attrs += ` l="` + strconv.FormatInt(c.Left, 10) + `"`
	}
	if c.Top != 0 {
		attrs += ` t="` + strconv.FormatInt(c.Top, 10) + `"`
	}
	if c.Right != 0 {
		attrs += ` r="` + strconv.FormatInt(c.Right, 10) + `"`
	}
	if c.Bottom != 0 {
		attrs += ` b="` + strconv.FormatInt(c.Bottom, 10) + `"`
	}
	if attrs != "" {
		return `<a:srcRect` + attrs + `/>`
	}
	return ""
}

// hasImageEffects reports whether any effect at all is configured, so the
// `<a:effectLst>` wrapper can be skipped entirely.
func hasImageEffects(image ImageRef) bool {
	return image.Shadow || image.Reflection || image.InnerShadow || image.Glow || image.SoftEdges ||
		image.GlowSpec != nil || image.BlurSpec != nil || image.SoftEdgeSpec != nil || image.ReflectionSpec != nil
}

// imageEffectsXML renders the picture's `<a:effectLst>`. Children are emitted
// in the CT_EffectList sequence order PowerPoint requires: blur, glow,
// innerShdw, outerShdw, reflection, softEdge.
func imageEffectsXML(image ImageRef) string {
	if !hasImageEffects(image) {
		return ""
	}
	var b strings.Builder
	b.WriteString("<a:effectLst>")
	appendImageBlurXML(&b, image)
	appendImageGlowXML(&b, image)
	appendImageInnerShadowXML(&b, image)
	shadow, reflection := image.Shadow, image.Reflection
	if shadow {
		b.WriteString(`<a:outerShdw blurRad="`)
		b.WriteString(strconv.Itoa(defaultShadowBlurRad))
		b.WriteString(`" dist="`)
		b.WriteString(strconv.Itoa(defaultShadowDist))
		b.WriteString(`" dir="`)
		b.WriteString(strconv.Itoa(defaultShadowDir))
		b.WriteString(`" rotWithShape="0"><a:srgbClr val="000000"><a:alpha val="`)
		b.WriteString(strconv.Itoa(defaultShadowAlpha))
		b.WriteString(`"/></a:srgbClr></a:outerShdw>`)
	}
	if image.ReflectionSpec != nil {
		b.WriteString(`<a:reflection blurRad="`)
		b.WriteString(strconv.Itoa(image.ReflectionSpec.BlurEmu))
		b.WriteString(`" dist="`)
		b.WriteString(strconv.Itoa(image.ReflectionSpec.DistanceEmu))
		b.WriteString(`" stA="`)
		b.WriteString(strconv.Itoa(defaultReflectionStA))
		b.WriteString(`" endA="`)
		b.WriteString(strconv.Itoa(defaultReflectionEndA))
		b.WriteString(`" endPos="`)
		b.WriteString(strconv.Itoa(defaultReflectionEndPos))
		b.WriteString(`" dir="`)
		b.WriteString(strconv.Itoa(defaultShadowDir))
		b.WriteString(`" sy="-100000" algn="bl" rotWithShape="0"/>`)
	} else if reflection {
		b.WriteString(`<a:reflection blurRad="`)
		b.WriteString(strconv.Itoa(defaultReflectionBlur))
		b.WriteString(`" stA="`)
		b.WriteString(strconv.Itoa(defaultReflectionStA))
		b.WriteString(`" endA="`)
		b.WriteString(strconv.Itoa(defaultReflectionEndA))
		b.WriteString(`" endPos="`)
		b.WriteString(strconv.Itoa(defaultReflectionEndPos))
		b.WriteString(`" dist="0" dir="`)
		b.WriteString(strconv.Itoa(defaultShadowDir))
		b.WriteString(`" sy="-100000" algn="bl" rotWithShape="0"/>`)
	}
	appendImageSoftEdgeXML(&b, image)
	b.WriteString("</a:effectLst>")
	return b.String()
}

func appendImageBlurXML(b *strings.Builder, image ImageRef) {
	if image.BlurSpec == nil {
		return
	}
	b.WriteString(`<a:blur rad="`)
	b.WriteString(strconv.Itoa(image.BlurSpec.RadiusEmu))
	b.WriteString(`" grow="1"/>`)
}

func appendImageGlowXML(b *strings.Builder, image ImageRef) {
	if image.GlowSpec != nil {
		b.WriteString(`<a:glow rad="`)
		b.WriteString(strconv.Itoa(image.GlowSpec.RadiusEmu))
		b.WriteString(`"><a:srgbClr val="`)
		b.WriteString(Escape(image.GlowSpec.Color))
		b.WriteString(`"/></a:glow>`)
		return
	}
	if image.Glow {
		b.WriteString(`<a:glow rad="`)
		b.WriteString(strconv.Itoa(defaultGlowRadius))
		b.WriteString(`"><a:srgbClr val="`)
		b.WriteString(defaultGlowColor)
		b.WriteString(`"><a:alpha val="`)
		b.WriteString(strconv.Itoa(defaultGlowAlpha))
		b.WriteString(`"/></a:srgbClr></a:glow>`)
	}
}

func appendImageInnerShadowXML(b *strings.Builder, image ImageRef) {
	if !image.InnerShadow {
		return
	}
	b.WriteString(`<a:innerShdw blurRad="`)
	b.WriteString(strconv.Itoa(defaultShadowBlurRad))
	b.WriteString(`" dist="`)
	b.WriteString(strconv.Itoa(defaultShadowDist))
	b.WriteString(`" dir="`)
	b.WriteString(strconv.Itoa(defaultShadowDir))
	b.WriteString(`"><a:srgbClr val="000000"><a:alpha val="`)
	b.WriteString(strconv.Itoa(defaultShadowAlpha))
	b.WriteString(`"/></a:srgbClr></a:innerShdw>`)
}

func appendImageSoftEdgeXML(b *strings.Builder, image ImageRef) {
	if image.SoftEdgeSpec != nil {
		b.WriteString(`<a:softEdge rad="`)
		b.WriteString(strconv.Itoa(image.SoftEdgeSpec.RadiusEmu))
		b.WriteString(`"/>`)
		return
	}
	if image.SoftEdges {
		b.WriteString(`<a:softEdge rad="`)
		b.WriteString(strconv.Itoa(defaultSoftEdgeRadius))
		b.WriteString(`"/>`)
	}
}
