package tplx

import (
	"bytes"
	"encoding/xml"
	"strings"
)

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

// mergeTokenRuns merges all runs in the paragraph into one when their combined text
// contains a Jinja token that spans multiple runs.
func mergeTokenRuns(children []paraChild) []paraChild {
	runIdxs := make([]int, 0)
	for i, c := range children {
		if c.isRun {
			runIdxs = append(runIdxs, i)
		}
	}
	if len(runIdxs) < minSubmatchLen {
		return children
	}

	var sb strings.Builder
	for _, ri := range runIdxs {
		sb.WriteString(children[ri].run.text)
	}
	full := sb.String()
	if !needsMerge(full, children, runIdxs) {
		return children
	}

	var combined strings.Builder
	var firstRPr []byte
	for i, ri := range runIdxs {
		combined.WriteString(children[ri].run.text)
		if i == 0 {
			firstRPr = children[ri].run.rprBytes
		}
	}
	mergedRun := paraChild{isRun: true, run: runData{rprBytes: firstRPr, text: combined.String()}}

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
	}
	return out
}

// needsMerge returns true when a token crosses at least one run boundary.
func needsMerge(full string, children []paraChild, runIdxs []int) bool {
	if !tokenPattern.MatchString(full) {
		return false
	}
	for _, loc := range tokenPattern.FindAllStringIndex(full, -1) {
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
