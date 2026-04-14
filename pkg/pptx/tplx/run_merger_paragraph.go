package tplx

import (
	"bytes"
	"encoding/xml"
	"sort"
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

// mergeTokenRuns merges only the minimal consecutive run ranges that are required to
// make each cross-boundary Jinja token contiguous. Runs that do not participate in any
// cross-boundary token are left untouched, preserving their original formatting.
func mergeTokenRuns(children []paraChild) []paraChild {
	// Collect positions of run children within the children slice.
	runIdxs := make([]int, 0, len(children))
	for i, c := range children {
		if c.isRun {
			runIdxs = append(runIdxs, i)
		}
	}
	if len(runIdxs) < minSubmatchLen {
		return children
	}

	// Concatenate all run texts and record where each run starts.
	var sb strings.Builder
	runStart := make([]int, len(runIdxs)) // byte offset of run[i] within the full string
	for i, ri := range runIdxs {
		runStart[i] = sb.Len()
		sb.WriteString(children[ri].run.text)
	}
	full := sb.String()

	if !tokenPattern.MatchString(full) {
		return children
	}

	// charToRunPos maps a byte offset in full to its run position index (within runIdxs).
	charToRunPos := func(charPos int) int {
		// Binary-search for the last run whose start <= charPos.
		lo, hi := 0, len(runStart)-1
		for lo < hi {
			mid := (lo + hi + 1) / 2
			if runStart[mid] <= charPos {
				lo = mid
			} else {
				hi = mid - 1
			}
		}
		return lo
	}

	// Find all token matches that cross a run boundary and record the run ranges.
	type mergeRange struct{ start, end int } // inclusive positions within runIdxs
	var ranges []mergeRange
	for _, loc := range tokenPattern.FindAllStringIndex(full, -1) {
		s := charToRunPos(loc[0])
		e := charToRunPos(loc[1] - 1)
		if s != e {
			ranges = append(ranges, mergeRange{s, e})
		}
	}
	if len(ranges) == 0 {
		return children
	}

	// Sort and consolidate overlapping / adjacent ranges so each run belongs to at most one group.
	sort.Slice(ranges, func(i, j int) bool { return ranges[i].start < ranges[j].start })
	merged := []mergeRange{ranges[0]}
	for _, r := range ranges[1:] {
		last := &merged[len(merged)-1]
		if r.start <= last.end+1 {
			if r.end > last.end {
				last.end = r.end
			}
		} else {
			merged = append(merged, r)
		}
	}

	// Pre-build one merged paraChild per consolidated range,
	// using the first run's rPr so at least the opening style is preserved.
	mergedRuns := make([]paraChild, len(merged))
	for gi, mr := range merged {
		var tb strings.Builder
		for pos := mr.start; pos <= mr.end; pos++ {
			tb.WriteString(children[runIdxs[pos]].run.text)
		}
		mergedRuns[gi] = paraChild{isRun: true, run: runData{
			rprBytes: children[runIdxs[mr.start]].run.rprBytes,
			text:     tb.String(),
		}}
	}

	// Map each child-slice index to its run position for O(1) lookup.
	childToRunPos := make(map[int]int, len(runIdxs))
	for pos, ci := range runIdxs {
		childToRunPos[ci] = pos
	}

	// Assign each run position to its merge group (-1 = not merged).
	runGroup := make([]int, len(runIdxs))
	for i := range runGroup {
		runGroup[i] = -1
	}
	for gi, mr := range merged {
		for pos := mr.start; pos <= mr.end; pos++ {
			runGroup[pos] = gi
		}
	}

	out := make([]paraChild, 0, len(children))
	emitted := make([]bool, len(merged))
	for i, c := range children {
		runPos, isRun := childToRunPos[i]
		if !isRun {
			out = append(out, c)
			continue
		}
		gi := runGroup[runPos]
		if gi < 0 {
			// Not involved in any cross-boundary token: keep as-is.
			out = append(out, c)
			continue
		}
		if !emitted[gi] {
			out = append(out, mergedRuns[gi])
			emitted[gi] = true
		}
		// Subsequent runs within the same group are absorbed; skip them.
	}
	return out
}
