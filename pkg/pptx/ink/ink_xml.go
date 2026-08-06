package ink

import (
	"strconv"
	"strings"
)

const (
	// PartContentType is the content-type override for an InkML part.
	PartContentType = "application/inkml+xml"
	// RelationshipType is the relationship type a slide uses to reach the part.
	// PowerPoint stores ink as a custom XML relationship.
	RelationshipType = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/customXml"

	inkmlNamespace = "http://www.w3.org/2003/InkML"
	p14Namespace   = "http://schemas.microsoft.com/office/powerpoint/2010/main"
	mcNamespace    = "http://schemas.openxmlformats.org/markup-compatibility/2006"

	xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`
)

// InkML renders the `ppt/ink/inkN.xml` part: one brush per stroke and one
// trace per stroke, with coordinates in HiMetric units.
func (a *Annotation) InkML() string {
	var b strings.Builder
	b.WriteString(xmlHeader)
	b.WriteString(`<inkml:ink xmlns:inkml="` + inkmlNamespace + `">`)
	b.WriteString(`<inkml:definitions>`)
	b.WriteString(inkContextXML())
	for i, stroke := range a.Strokes {
		if stroke.IsEmpty() {
			continue
		}
		writeBrush(&b, i, stroke.Pen)
	}
	b.WriteString(`</inkml:definitions>`)
	for i, stroke := range a.Strokes {
		if stroke.IsEmpty() {
			continue
		}
		writeTrace(&b, i, stroke)
	}
	b.WriteString(`</inkml:ink>`)
	return b.String()
}

// inkContextXML describes the coordinate space traces are written in: X and Y
// as HiMetric integers, which is what PowerPoint writes for mouse and pen ink.
func inkContextXML() string {
	return `<inkml:context xml:id="ctx0">` +
		`<inkml:inkSource xml:id="inkSrc0">` +
		`<inkml:traceFormat>` +
		`<inkml:channel name="X" type="integer" units="himetric"/>` +
		`<inkml:channel name="Y" type="integer" units="himetric"/>` +
		`</inkml:traceFormat>` +
		`<inkml:channelProperties>` +
		`<inkml:channelProperty channel="X" name="resolution" value="1" units="1/himetric"/>` +
		`<inkml:channelProperty channel="Y" name="resolution" value="1" units="1/himetric"/>` +
		`</inkml:channelProperties>` +
		`</inkml:inkSource>` +
		`</inkml:context>`
}

func writeBrush(b *strings.Builder, index int, pen Pen) {
	width := emuToHiMetric(pen.WidthEmu)
	height := emuToHiMetric(pen.HeightEmu)
	if height <= 0 {
		height = width
	}

	b.WriteString(`<inkml:brush xml:id="br`)
	b.WriteString(strconv.Itoa(index))
	b.WriteString(`">`)
	writeBrushProperty(b, "width", itoa(width), "himetric")
	writeBrushProperty(b, "height", itoa(height), "himetric")
	writeBrushProperty(b, "color", "#"+normalizeColor(pen.Color), "")
	writeBrushProperty(b, "tip", pen.Tip.XMLValue(), "")
	if pen.Transparency > fullyOpaque {
		writeBrushProperty(b, "transparency", strconv.Itoa(int(pen.Transparency)), "")
	}
	if pen.IgnorePressure {
		writeBrushProperty(b, "ignorePressure", "1", "")
	}
	b.WriteString(`</inkml:brush>`)
}

func writeBrushProperty(b *strings.Builder, name, value, units string) {
	b.WriteString(`<inkml:brushProperty name="`)
	b.WriteString(name)
	b.WriteString(`" value="`)
	b.WriteString(escapeAttr(value))
	b.WriteString(`"`)
	if units != "" {
		b.WriteString(` units="`)
		b.WriteString(units)
		b.WriteString(`"`)
	}
	b.WriteString(`/>`)
}

func writeTrace(b *strings.Builder, index int, stroke Stroke) {
	b.WriteString(`<inkml:trace contextRef="#ctx0" brushRef="#br`)
	b.WriteString(strconv.Itoa(index))
	b.WriteString(`">`)
	for i, p := range stroke.Points {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(itoa(emuToHiMetric(p.X.Emu())))
		b.WriteString(" ")
		b.WriteString(itoa(emuToHiMetric(p.Y.Emu())))
	}
	b.WriteString(`</inkml:trace>`)
}

// ContentPartXML renders the shape-tree markup that points at the ink part.
// The ink element is wrapped in mc:AlternateContent because `p:contentPart` is
// a PowerPoint 2010 extension: readers that do not understand it skip the
// whole block instead of rejecting the slide.
func (a *Annotation) ContentPartXML(relID string, shapeID int) string {
	x, y, cx, cy := a.frame()
	name := a.Name
	if name == "" {
		name = "Ink"
	}

	var b strings.Builder
	b.WriteString(`<mc:AlternateContent xmlns:mc="` + mcNamespace + `">`)
	b.WriteString(`<mc:Choice xmlns:p14="` + p14Namespace + `" Requires="p14">`)
	b.WriteString(`<p:contentPart r:id="` + escapeAttr(relID) + `" p14:bwMode="auto">`)
	b.WriteString(`<p14:nvContentPartPr>`)
	b.WriteString(`<p14:cNvPr id="` + strconv.Itoa(shapeID) + `" name="` + escapeAttr(name) + `"/>`)
	b.WriteString(`<p14:cNvContentPartPr/>`)
	b.WriteString(`<p14:nvPr/>`)
	b.WriteString(`</p14:nvContentPartPr>`)
	b.WriteString(`<p14:xfrm>`)
	b.WriteString(`<a:off x="` + itoa(x) + `" y="` + itoa(y) + `"/>`)
	b.WriteString(`<a:ext cx="` + itoa(cx) + `" cy="` + itoa(cy) + `"/>`)
	b.WriteString(`</p14:xfrm>`)
	b.WriteString(`</p:contentPart>`)
	b.WriteString(`</mc:Choice>`)
	b.WriteString(`</mc:AlternateContent>`)
	return b.String()
}

func escapeAttr(v string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(v)
}
