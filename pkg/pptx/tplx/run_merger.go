// Package tplx provides a Jinja-style file-based template engine for PPTX files.
// Author a .pptx in PowerPoint with {{variable}} tokens in text boxes, table cells,
// or notes — then render it with a Go data map.
//
// Example:
//
//	result, err := tplx.Render("invoice_template.pptx", tplx.Context{
//	    "name":  "Acme Corp",
//	    "date":  "2026-03-01",
//	    "rows": []tplx.Row{
//	        {"item": "Widget A", "qty": "100", "price": "$1,200"},
//	    },
//	})
//	if err != nil { panic(err) }
//	result.Save("invoice_acme.pptx")
package tplx

import (
	"bytes"
	"encoding/xml"
	"regexp"
	"strings"
)

// tokenPattern matches a complete Jinja-style token like {{foo}} or {{#each items}}.
var tokenPattern = regexp.MustCompile(`\{\{[^{}]+\}\}`)

// ── run types ────────────────────────────────────────────────────────────────

// runData holds the parsed data for one <a:r> XML element.
type runData struct {
	rprBytes []byte   // raw re-encoded bytes for <a:rPr> child (nil if absent)
	text     string   // content of <a:t>
	extras   [][]byte // other child bytes (preserved verbatim, rare)
}

// paraChild is a polymorphic child inside a paragraph: either a run or raw XML bytes.
type paraChild struct {
	isRun bool
	run   runData
	raw   []byte
}

// ── public helpers ────────────────────────────────────────────────────────────

// mergeAdjacentRuns merges contiguous <a:r> runs within each <a:p> paragraph when
// those runs together form a Jinja token that PowerPoint split across multiple runs.
//
// PowerPoint routinely fragments "{{name}}" into 3-4 separate <a:r> elements due to
// spell-check and autocorrect boundaries. Template substitution requires the full
// token to be in a single run's <a:t> text.
func mergeAdjacentRuns(xmlBytes []byte) []byte {
	dec := xml.NewDecoder(bytes.NewReader(xmlBytes))
	dec.Strict = false

	var out bytes.Buffer
	out.Grow(len(xmlBytes))

	if err := streamRewrite(dec, &out, ""); err != nil {
		return xmlBytes
	}
	return out.Bytes()
}

// ── streaming rewriter ────────────────────────────────────────────────────────

// streamRewrite recursively re-emits tokens, rewriting <a:p> blocks.
func streamRewrite(dec *xml.Decoder, out *bytes.Buffer, parentLocal string) error {
	for {
		tok, err := dec.Token()
		if err != nil {
			return handleStreamReadError(parentLocal, err)
		}
		switch t := tok.(type) {
		case xml.ProcInst:
			continue
		case xml.StartElement:
			if err = handleStreamStartElement(dec, out, t); err != nil {
				return err
			}
		case xml.EndElement:
			if shouldStop := handleStreamEndElement(out, parentLocal, t); shouldStop {
				return nil
			}
		case xml.CharData:
			handleStreamCharData(out, parentLocal, t)
		case xml.Comment:
			writeXMLComment(out, t)
		}
	}
}

func handleStreamReadError(parentLocal string, err error) error {
	if parentLocal == "" {
		return nil // EOF at top level = done.
	}
	return err //nolint:wrapcheck // Streaming decoder errors are returned unchanged for callers.
}

func handleStreamStartElement(dec *xml.Decoder, out *bytes.Buffer, start xml.StartElement) error {
	if isDrawingMLLocal(start, "p") {
		return writeParagraph(dec, out, start)
	}
	out.Write(reencStart(start))
	return streamRewrite(dec, out, start.Name.Local)
}

func handleStreamEndElement(out *bytes.Buffer, parentLocal string, end xml.EndElement) bool {
	out.Write(reencEnd(end))
	return end.Name.Local == parentLocal
}

func handleStreamCharData(out *bytes.Buffer, parentLocal string, charData xml.CharData) {
	if parentLocal == "" {
		return
	}
	writeEscaped(out, string(charData))
}

func writeXMLComment(out *bytes.Buffer, comment xml.Comment) {
	out.WriteString("<!--")
	out.Write(comment)
	out.WriteString("-->")
}

// ── paragraph rewriter ────────────────────────────────────────────────────────

// writeParagraph reads a full <a:p>…</a:p> block, merges fragmented token runs, then writes it.
func writeParagraph(dec *xml.Decoder, out *bytes.Buffer, paraStart xml.StartElement) error {
	children, err := collectParaChildren(dec)
	if err != nil {
		return err
	}
	children = mergeTokenRuns(children)

	out.Write(reencStart(paraStart))
	for _, c := range children {
		if c.isRun {
			emitRun(out, c.run)
		} else {
			out.Write(c.raw)
		}
	}
	out.WriteString("</a:p>")
	return nil
}

// collectParaChildren reads all children of a paragraph until </a:p>.
func collectParaChildren(dec *xml.Decoder) ([]paraChild, error) {
	var children []paraChild
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return nil, err //nolint:wrapcheck // Preserve decoder errors from malformed paragraph XML.
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if isDrawingMLLocal(t, "r") {
				rd, err2 := readRun(dec)
				if err2 != nil {
					return nil, err2
				}
				children = append(children, paraChild{isRun: true, run: rd})
			} else {
				depth++
				var childBuf bytes.Buffer
				childBuf.Write(reencStart(t))
				if err2 := collectRawDepth(dec, &childBuf); err2 != nil {
					return nil, err2
				}
				depth--
				children = append(children, paraChild{raw: childBuf.Bytes()})
			}
		case xml.EndElement:
			depth--
		case xml.CharData:
			var buf bytes.Buffer
			writeEscaped(&buf, string(t))
			children = append(children, paraChild{raw: buf.Bytes()})
		}
	}
	return children, nil
}

// ── run merger ────────────────────────────────────────────────────────────────

// mergeTokenRuns merges all runs in the paragraph into one when their combined text
// contains a Jinja token that spans multiple runs.
func mergeTokenRuns(children []paraChild) []paraChild {
	// Collect positions of run children.
	runIdxs := make([]int, 0)
	for i, c := range children {
		if c.isRun {
			runIdxs = append(runIdxs, i)
		}
	}
	if len(runIdxs) < minSubmatchLen {
		return children
	}

	// Build concatenated text.
	var sb strings.Builder
	for _, ri := range runIdxs {
		sb.WriteString(children[ri].run.text)
	}
	full := sb.String()
	if !needsMerge(full, children, runIdxs) {
		return children
	}

	// Merge: combine all run texts into one run, inheriting first run's rPr.
	var combined strings.Builder
	var firstRPr []byte
	for i, ri := range runIdxs {
		combined.WriteString(children[ri].run.text)
		if i == 0 {
			firstRPr = children[ri].run.rprBytes
		}
	}
	mergedRun := paraChild{isRun: true, run: runData{rprBytes: firstRPr, text: combined.String()}}

	// Rebuild children: keep non-runs, replace first run slot with merged, drop rest.
	runSet := make(map[int]bool, len(runIdxs))
	for _, ri := range runIdxs {
		runSet[ri] = true
	}
	out := make([]paraChild, 0, len(children))
	firstEmitted := false
	for i, c := range children {
		if !runSet[i] {
			out = append(out, c)
			continue
		}
		if !firstEmitted {
			out = append(out, mergedRun)
			firstEmitted = true
		}
		// Subsequent runs dropped — their text is in mergedRun.
	}
	return out
}

// needsMerge returns true when a token crosses at least one run boundary.
func needsMerge(full string, children []paraChild, runIdxs []int) bool {
	if !tokenPattern.MatchString(full) {
		return false
	}
	for _, loc := range tokenPattern.FindAllStringIndex(full, -1) {
		// Map char positions to run indices.
		startRun := charToRunIdx(loc[0], children, runIdxs)
		endRun := charToRunIdx(loc[1]-1, children, runIdxs)
		if startRun != endRun {
			return true
		}
	}
	return false
}

// charToRunIdx returns the run-position-index (within runIdxs) for a character offset in full.
func charToRunIdx(charPos int, children []paraChild, runIdxs []int) int {
	offset := 0
	for ri, pos := range runIdxs {
		l := len(children[pos].run.text)
		if charPos < offset+l {
			return ri
		}
		offset += l
	}
	return len(runIdxs) - 1
}

// ── run reader ────────────────────────────────────────────────────────────────

// readRun reads a single <a:r>…</a:r> block, extracting rPr and text.
func readRun(dec *xml.Decoder) (runData, error) {
	var rd runData
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return rd, err //nolint:wrapcheck // Preserve tokenization errors while reading run content.
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			local := t.Name.Local
			switch local {
			case "rPr":
				var buf bytes.Buffer
				buf.Write(reencStart(t))
				if err2 := collectRawDepth(dec, &buf); err2 != nil {
					return rd, err2
				}
				depth--
				rd.rprBytes = buf.Bytes()
			case "t":
				text, err2 := readCharData(dec)
				if err2 != nil {
					return rd, err2
				}
				depth--
				rd.text = text
			default:
				var buf bytes.Buffer
				buf.Write(reencStart(t))
				if err2 := collectRawDepth(dec, &buf); err2 != nil {
					return rd, err2
				}
				depth--
				rd.extras = append(rd.extras, buf.Bytes())
			}
		case xml.EndElement:
			depth--
		}
	}
	return rd, nil
}

// readCharData reads text content until the matching end element.
func readCharData(dec *xml.Decoder) (string, error) {
	var sb strings.Builder
	for {
		tok, err := dec.Token()
		if err != nil {
			return sb.String(), err //nolint:wrapcheck // Preserve decoder errors while reading character data.
		}
		switch t := tok.(type) {
		case xml.CharData:
			sb.Write(t)
		case xml.EndElement:
			return sb.String(), nil
		}
	}
}

// collectRawDepth reads elements until depth returns to 0, appending to buf.
func collectRawDepth(dec *xml.Decoder, buf *bytes.Buffer) error {
	depth := 1
	for depth > 0 {
		tok, err := dec.Token()
		if err != nil {
			return err //nolint:wrapcheck // Preserve raw depth parsing failures for caller handling.
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			buf.Write(reencStart(t))
		case xml.EndElement:
			depth--
			if depth >= 0 {
				buf.Write(reencEnd(t))
			}
		case xml.CharData:
			writeEscaped(buf, string(t))
		}
	}
	return nil
}

// ── XML helpers ───────────────────────────────────────────────────────────────

// isDrawingMLLocal reports whether s is a DrawingML element with the given local name.
func isDrawingMLLocal(s xml.StartElement, local string) bool {
	if s.Name.Local != local {
		return false
	}
	return s.Name.Space == "a" ||
		s.Name.Space == "http://schemas.openxmlformats.org/drawingml/2006/main" ||
		s.Name.Space == ""
}

// reencStart re-encodes a StartElement, preserving namespace prefixes.
func reencStart(s xml.StartElement) []byte {
	var b bytes.Buffer
	b.WriteByte('<')
	writeQName(&b, s.Name)
	for _, a := range s.Attr {
		b.WriteByte(' ')
		writeQName(&b, a.Name)
		b.WriteString(`="`)
		if err := xml.EscapeText(&b, []byte(a.Value)); err != nil {
			b.WriteString(a.Value)
		}
		b.WriteByte('"')
	}
	b.WriteByte('>')
	return b.Bytes()
}

// reencEnd re-encodes an EndElement.
func reencEnd(e xml.EndElement) []byte {
	var b bytes.Buffer
	b.WriteString("</")
	writeQName(&b, e.Name)
	b.WriteByte('>')
	return b.Bytes()
}

func writeQName(b *bytes.Buffer, name xml.Name) {
	if name.Space != "" {
		if prefix := namespacePrefix(name.Space); prefix != "" {
			b.WriteString(prefix)
			b.WriteByte(':')
		}
	}
	b.WriteString(name.Local)
}

func namespacePrefix(space string) string {
	switch space {
	case "http://schemas.openxmlformats.org/presentationml/2006/main":
		return "p"
	case "http://schemas.openxmlformats.org/drawingml/2006/main":
		return "a"
	case "http://schemas.openxmlformats.org/officeDocument/2006/relationships":
		return "r"
	case "http://schemas.openxmlformats.org/markup-compatibility/2006":
		return "mc"
	case "http://schemas.microsoft.com/office/drawing/2010/main":
		return "a14"
	case "http://schemas.openxmlformats.org/package/2006/relationships":
		return "rel"
	case "http://www.w3.org/XML/1998/namespace":
		return "xml"
	case "xmlns":
		return "xmlns"
	default:
		if strings.Contains(space, "://") {
			return ""
		}
		return space
	}
}

// writeEscaped writes XML-escaped character data.
func writeEscaped(out *bytes.Buffer, s string) {
	if err := xml.EscapeText(out, []byte(s)); err != nil {
		out.WriteString(s)
	}
}

// emitRun writes a complete <a:r>…</a:r> element to out.
func emitRun(out *bytes.Buffer, rd runData) {
	out.WriteString("<a:r>")
	if len(rd.rprBytes) > 0 {
		out.Write(rd.rprBytes)
	}
	for _, x := range rd.extras {
		out.Write(x)
	}
	out.WriteString("<a:t>")
	writeEscaped(out, rd.text)
	out.WriteString("</a:t>")
	out.WriteString("</a:r>")
}
