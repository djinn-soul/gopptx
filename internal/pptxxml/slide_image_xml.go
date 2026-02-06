package pptxxml

import "fmt"

// ImageRef describes one image reference in a slide.
type ImageRef struct {
	RelID string
	Name  string
	X     int64
	Y     int64
	CX    int64
	CY    int64
}

func imageShape(image ImageRef, shapeID int) string {
	name := image.Name
	if name == "" {
		name = "Picture"
	}
	return fmt.Sprintf(`
<p:pic>
<p:nvPicPr>
<p:cNvPr id="%d" name="%s"/>
<p:cNvPicPr><a:picLocks noChangeAspect="1"/></p:cNvPicPr>
<p:nvPr/>
</p:nvPicPr>
<p:blipFill>
<a:blip r:embed="%s"/>
<a:stretch><a:fillRect/></a:stretch>
</p:blipFill>
<p:spPr>
<a:xfrm>
<a:off x="%d" y="%d"/>
<a:ext cx="%d" cy="%d"/>
</a:xfrm>
<a:prstGeom prst="rect"><a:avLst/></a:prstGeom>
</p:spPr>
</p:pic>`,
		shapeID,
		Escape(name),
		Escape(image.RelID),
		image.X,
		image.Y,
		image.CX,
		image.CY,
	)
}
