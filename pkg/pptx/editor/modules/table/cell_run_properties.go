package table

import (
	"bytes"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// fontSzScale converts font size in points to OOXML hundredths-of-points.
const fontSzScale = 100

// rPrChildOrder is the child order of CT_TextCharacterProperties. PowerPoint
// refuses to open a file whose <a:rPr> children are out of schema order, so a
// merged rPr is always re-emitted through this list.
var rPrChildOrder = []string{
	"ln",
	"noFill", "solidFill", "gradFill", "blipFill", "pattFill", "grpFill",
	"effectLst", "effectDag",
	"highlight",
	"uLnTx", "uLn", "uFillTx", "uFill",
	"latin", "ea", "cs", "sym",
	"hlinkClick", "hlinkMouseOver",
	"rtl", "extLst",
}

// fillChildren is the mutually exclusive fill choice group: setting one must
// remove the others.
var fillChildren = []string{"noFill", "solidFill", "gradFill", "blipFill", "pattFill", "grpFill"}

// runProperties is a parsed <a:rPr> that keeps attribute order stable and
// re-emits children in schema order.
type runProperties struct {
	attrOrder []string
	attrs     map[string]string
	children  map[string]string
}

func newRunProperties() *runProperties {
	return &runProperties{attrs: map[string]string{}, children: map[string]string{}}
}

// parseRunProperties reads an existing <a:rPr .../> or <a:rPr ...>...</a:rPr>.
// A cell with no run properties yields an empty, valid value.
func parseRunProperties(rPrXML []byte) *runProperties {
	out := newRunProperties()
	if len(rPrXML) == 0 {
		return out
	}
	tagEnd := bytes.IndexByte(rPrXML, '>')
	if tagEnd == -1 {
		return out
	}

	head := rPrXML[len("<a:rPr"):tagEnd]
	head = bytes.TrimSuffix(bytes.TrimSpace(head), []byte("/"))
	values := parseXMLAttrs(head)
	for _, name := range attrOrderFrom(head) {
		out.setAttr(name, values[name])
	}

	if rPrXML[tagEnd-1] == '/' {
		return out
	}
	inner := rPrXML[tagEnd+1:]
	if closeIdx := bytes.Index(inner, []byte("</a:rPr>")); closeIdx != -1 {
		inner = inner[:closeIdx]
	}
	for _, tag := range rPrChildOrder {
		if el := extractXMLElement(inner, []byte("<a:"+tag)); len(el) > 0 {
			out.children[tag] = string(el)
		}
	}
	return out
}

// parseXMLAttrs pulls name="value" pairs out of a start-tag body.
func parseXMLAttrs(head []byte) map[string]string {
	attrs := map[string]string{}
	for _, name := range attrOrderFrom(head) {
		idx := bytes.Index(head, []byte(name+`="`))
		if idx == -1 {
			continue
		}
		rest := head[idx+len(name)+2:]
		if end := bytes.IndexByte(rest, '"'); end != -1 {
			attrs[name] = string(rest[:end])
		}
	}
	return attrs
}

// attrOrderFrom lists attribute names in the order they appear.
func attrOrderFrom(head []byte) []string {
	var names []string
	rest := head
	for {
		eq := bytes.Index(rest, []byte(`="`))
		if eq == -1 {
			return names
		}
		nameStart := eq
		for nameStart > 0 && !isXMLSpace(rest[nameStart-1]) {
			nameStart--
		}
		name := strings.TrimSpace(string(rest[nameStart:eq]))
		if name != "" {
			names = append(names, name)
		}
		valEnd := bytes.IndexByte(rest[eq+2:], '"')
		if valEnd == -1 {
			return names
		}
		rest = rest[eq+2+valEnd+1:]
	}
}

func isXMLSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func (r *runProperties) setAttr(name, value string) {
	if _, seen := r.attrs[name]; !seen {
		r.attrOrder = append(r.attrOrder, name)
	}
	r.attrs[name] = value
}

func (r *runProperties) setFill(xml string) {
	for _, tag := range fillChildren {
		delete(r.children, tag)
	}
	r.children["solidFill"] = xml
}

func (r *runProperties) render() string {
	var b strings.Builder
	b.WriteString("<a:rPr")
	seen := map[string]bool{}
	for _, name := range r.attrOrder {
		if seen[name] {
			continue
		}
		seen[name] = true
		value, ok := r.attrs[name]
		if !ok {
			continue
		}
		b.WriteString(` ` + name + `="` + value + `"`)
	}

	var body strings.Builder
	for _, tag := range rPrChildOrder {
		if xml, ok := r.children[tag]; ok {
			body.WriteString(xml)
		}
	}
	if body.Len() == 0 {
		b.WriteString("/>")
		return b.String()
	}
	b.WriteString(">" + body.String() + "</a:rPr>")
	return b.String()
}

// mergeCellRPr layers the requested style on top of the cell's existing run
// properties. Upstream issue #1037: updating a cell's data must not discard the
// font, size, colour or any other formatting the deck author set by hand, so
// every field the caller left unset is carried through untouched.
func mergeCellRPr(existing []byte, update CellContentUpdate) string {
	rPr := parseRunProperties(existing)

	if update.SizePt > 0 {
		rPr.setAttr("sz", strconv.Itoa(int(update.SizePt*fontSzScale)))
	}
	if update.Bold != nil {
		rPr.setAttr("b", boolAttr(*update.Bold))
	}
	if update.Italic != nil {
		rPr.setAttr("i", boolAttr(*update.Italic))
	}
	if update.Underline != nil {
		rPr.setAttr("u", underlineAttr(*update.Underline))
	}
	if color := strings.TrimSpace(update.Color); color != "" {
		rPr.setFill(`<a:solidFill><a:srgbClr val="` +
			common.XMLEscape(strings.TrimPrefix(color, "#")) + `"/></a:solidFill>`)
	}
	if font := strings.TrimSpace(update.FontName); font != "" {
		rPr.children["latin"] = `<a:latin typeface="` + common.XMLEscape(font) + `"/>`
	}
	return rPr.render()
}

func boolAttr(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func underlineAttr(v bool) string {
	if v {
		return "sng"
	}
	return "none"
}
