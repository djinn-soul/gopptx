package pptxxml

import "strings"

// SlideBackgroundSpec describes how the slide background should be filled.
type SlideBackgroundSpec struct {
	Type         string // "solid", "gradient", "picture"
	SolidFill    *ShapeFillSpec
	GradientFill *ShapeGradientFillSpec
	PictureFill  *ImageRef
}

const slideDefaultBackground = `
<p:bg>
<p:bgRef idx="1001">
<a:schemeClr val="bg1"/>
</p:bgRef>
</p:bg>`

func backgroundXML(bg *SlideBackgroundSpec) string {
	if bg == nil || bg.Type == "" {
		return ""
	}

	xml := `
<p:bg>
<p:bgPr>`

	switch bg.Type {
	case "solid":
		if bg.SolidFill != nil {
			xml += `
<a:solidFill>
<a:srgbClr val="` + strings.TrimPrefix(bg.SolidFill.Color, "#") + `"/>
</a:solidFill>`
		}
	case "gradient":
		if bg.GradientFill != nil {
			xml += shapeGradientFillXML(*bg.GradientFill)
		}
	case slideBackgroundPicture:
		if bg.PictureFill != nil {
			xml += `
<a:blipFill>
<a:blip r:embed="` + FastEscapeRID(bg.PictureFill.RelID) + `"/>
<a:stretch>
<a:fillRect/>
</a:stretch>
</a:blipFill>`
		}
	}

	xml += `
<a:effectLst/>
</p:bgPr>
</p:bg>`
	return xml
}
