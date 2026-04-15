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
	runIdxs := collectRunIndices(children)
	if len(runIdxs) < minSubmatchLen {
		return children
	}

	full, runStart := buildRunText(children, runIdxs)
	if !tokenPattern.MatchString(full) {
		return children
	}

	ranges := tokenCrossRunRanges(full, runStart)
	if len(ranges) == 0 {
		return children
	}

	merged := consolidateMergeRanges(ranges)
	mergedRuns := buildMergedRuns(children, runIdxs, merged)
	childToRunPos := buildChildToRunPos(runIdxs)
	runGroup := buildRunGroupMap(len(runIdxs), merged)
	return emitMergedChildren(children, childToRunPos, runGroup, mergedRuns)
}

type mergeRange struct {
	start int
	end   int
}

func collectRunIndices(children []paraChild) []int {
	runIdxs := make([]int, 0, len(children))
	for i, c := range children {
		if c.isRun {
			runIdxs = append(runIdxs, i)
		}
	}
	return runIdxs
}

func buildRunText(children []paraChild, runIdxs []int) (string, []int) {
	var sb strings.Builder
	runStart := make([]int, len(runIdxs))
	for i, ri := range runIdxs {
		runStart[i] = sb.Len()
		sb.WriteString(children[ri].run.text)
	}
	return sb.String(), runStart
}

func runPosForByte(runStart []int, charPos int) int {
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

func tokenCrossRunRanges(full string, runStart []int) []mergeRange {
	var ranges []mergeRange
	for _, loc := range tokenPattern.FindAllStringIndex(full, -1) {
		start := runPosForByte(runStart, loc[0])
		end := runPosForByte(runStart, loc[1]-1)
		if start != end {
			ranges = append(ranges, mergeRange{start: start, end: end})
		}
	}
	return ranges
}

func consolidateMergeRanges(ranges []mergeRange) []mergeRange {
	sort.Slice(ranges, func(i, j int) bool { return ranges[i].start < ranges[j].start })
	merged := []mergeRange{ranges[0]}
	for _, r := range ranges[1:] {
		last := &merged[len(merged)-1]
		if r.start > last.end+1 {
			merged = append(merged, r)
			continue
		}
		if r.end > last.end {
			last.end = r.end
		}
	}
	return merged
}

func buildMergedRuns(children []paraChild, runIdxs []int, merged []mergeRange) []paraChild {
	mergedRuns := make([]paraChild, len(merged))
	for groupIdx, mr := range merged {
		var tb strings.Builder
		for runPos := mr.start; runPos <= mr.end; runPos++ {
			tb.WriteString(children[runIdxs[runPos]].run.text)
		}
		mergedRuns[groupIdx] = paraChild{isRun: true, run: runData{
			rprBytes: children[runIdxs[mr.start]].run.rprBytes,
			text:     tb.String(),
		}}
	}
	return mergedRuns
}

func buildChildToRunPos(runIdxs []int) map[int]int {
	childToRunPos := make(map[int]int, len(runIdxs))
	for runPos, childIdx := range runIdxs {
		childToRunPos[childIdx] = runPos
	}
	return childToRunPos
}

func buildRunGroupMap(runCount int, merged []mergeRange) []int {
	runGroup := make([]int, runCount)
	for i := range runGroup {
		runGroup[i] = -1
	}
	for groupIdx, mr := range merged {
		for runPos := mr.start; runPos <= mr.end; runPos++ {
			runGroup[runPos] = groupIdx
		}
	}
	return runGroup
}

func emitMergedChildren(
	children []paraChild,
	childToRunPos map[int]int,
	runGroup []int,
	mergedRuns []paraChild,
) []paraChild {
	out := make([]paraChild, 0, len(children))
	emitted := make([]bool, len(mergedRuns))
	for childIdx, child := range children {
		runPos, isRun := childToRunPos[childIdx]
		if !isRun {
			out = append(out, child)
			continue
		}
		groupIdx := runGroup[runPos]
		if groupIdx < 0 {
			out = append(out, child)
			continue
		}
		if emitted[groupIdx] {
			continue
		}
		out = append(out, mergedRuns[groupIdx])
		emitted[groupIdx] = true
	}
	return out
}
