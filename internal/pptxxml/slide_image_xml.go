package pptxxml

import (
	"fmt"
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
)

func imageShape(image ImageRef, shapeID int) string {
	name := image.Name
	if name == "" {
		name = "Picture"
	}

	descrAttr := imageDescriptionAttr(image)
	srcRect := imageSrcRectXML(image.Crop)
	effectsXML := imageEffectsXML(image.Shadow, image.Reflection)
	xfrmAttrs := imageTransformAttrs(image.Rotation, image.FlipH, image.FlipV)

	return fmt.Sprintf(`
<p:pic>
<p:nvPicPr>
<p:cNvPr id="%d" name="%s"%s/>
<p:cNvPicPr><a:picLocks noChangeAspect="1"/></p:cNvPicPr>
<p:nvPr/>
</p:nvPicPr>
<p:blipFill>
<a:blip r:embed="%s"/>
%s
<a:stretch><a:fillRect/></a:stretch>
</p:blipFill>
<p:spPr>
<a:xfrm%s>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
%s
</p:spPr>
</p:pic>`,
		shapeID,
		Escape(name),
		descrAttr,
		Escape(image.RelID),
		srcRect,
		xfrmAttrs,
		image.X,
		image.Y,
		image.CX,
		image.CY,
		effectsXML,
	)
}

func imageDescriptionAttr(image ImageRef) string {
	if image.IsDecorative {
		return ` descr=""`
	}
	if image.AltText != "" {
		return fmt.Sprintf(` descr="%s"`, Escape(image.AltText))
	}
	return ""
}

func imageTransformAttrs(rot int64, flipH, flipV bool) string {
	attrs := ""
	if rot != 0 {
		attrs += fmt.Sprintf(` rot="%d"`, rot)
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
		attrs += fmt.Sprintf(` l="%d"`, c.Left)
	}
	if c.Top != 0 {
		attrs += fmt.Sprintf(` t="%d"`, c.Top)
	}
	if c.Right != 0 {
		attrs += fmt.Sprintf(` r="%d"`, c.Right)
	}
	if c.Bottom != 0 {
		attrs += fmt.Sprintf(` b="%d"`, c.Bottom)
	}
	if attrs != "" {
		return fmt.Sprintf(`<a:srcRect%s/>`, attrs)
	}
	return ""
}

func imageEffectsXML(shadow, reflection bool) string {
	if !shadow && !reflection {
		return ""
	}
	var b strings.Builder
	b.WriteString("<a:effectLst>")
	if shadow {
		b.WriteString(fmt.Sprintf(
			`<a:outerShdw blurRad="%d" dist="%d" dir="%d" rotWithShape="0">`+
				`<a:srgbClr val="000000"><a:alpha val="%d"/></a:srgbClr></a:outerShdw>`,
			defaultShadowBlurRad, defaultShadowDist,
			defaultShadowDir, defaultShadowAlpha,
		))
	}
	if reflection {
		b.WriteString(fmt.Sprintf(
			`<a:ref blurRad="%d" stA="%d" endA="%d" endPos="%d" dist="0" dir="%d" `+
				`sy="-100000" algn="bl" rotWithShape="0"/>`,
			defaultReflectionBlur, defaultReflectionStA,
			defaultReflectionEndA, defaultReflectionEndPos,
			defaultShadowDir,
		))
	}
	b.WriteString("</a:effectLst>")
	return b.String()
}
