package editor

import (
	"bytes"
	"fmt"
)

// applyOrDescend updates the target shape wherever it sits in the shape tree.
// The slide scan consumes each <p:grpSp> as one node, so a shape inside a group
// is never offered to apply directly; without the descent every update to a
// nested leaf fails with "shape id ... not found".
func (u *shapeUpdater) applyOrDescend(content []byte) func(int, *parsedShape) ([]byte, bool) {
	return func(i int, s *parsedShape) ([]byte, bool) {
		if updated, ok := u.apply(i, s); ok {
			return updated, true
		}
		if u.err != nil || u.found || !s.IsGroup {
			return nil, false
		}
		return u.applyWithinGroup(content[s.Start:s.End])
	}
}

// applyWithinGroup rewrites one group node when the target shape is somewhere
// beneath it, at any nesting depth.
func (u *shapeUpdater) applyWithinGroup(groupXML []byte) ([]byte, bool) {
	innerStart, innerEnd, ok := groupChildrenRange(groupXML)
	if !ok {
		return nil, false
	}
	inner := groupXML[innerStart:innerEnd]
	children, err := scanShapesWithOffsets(inner, false)
	if err != nil {
		u.err = fmt.Errorf("parse group children: %w", err)
		return nil, false
	}

	// Child offsets are relative to the group's own bytes, so the updater has
	// to slice against those bytes while the descent is in progress.
	outer := u.origSlide
	u.origSlide = inner
	newInner := replaceShapeNodes(inner, children, u.applyOrDescend(inner))
	u.origSlide = outer
	if u.err != nil || !u.found {
		return nil, false
	}

	out := make([]byte, 0, len(groupXML)-len(inner)+len(newInner))
	out = append(out, groupXML[:innerStart]...)
	out = append(out, newInner...)
	out = append(out, groupXML[innerEnd:]...)
	return out, true
}

// groupChildrenRange locates the bytes between a group's own start and end
// tags, which is the region its child shapes live in.
func groupChildrenRange(groupXML []byte) (int, int, bool) {
	name, ok := startTagQualifiedName(groupXML)
	if !ok {
		return 0, 0, false
	}
	startTagEnd := bytes.IndexByte(groupXML, '>')
	if startTagEnd < 0 {
		return 0, 0, false
	}
	// A self-closing group has no children to descend into.
	if startTagEnd > 0 && groupXML[startTagEnd-1] == '/' {
		return 0, 0, false
	}
	closing := []byte("</" + name + ">")
	closeStart := bytes.LastIndex(groupXML, closing)
	if closeStart <= startTagEnd {
		return 0, 0, false
	}
	return startTagEnd + 1, closeStart, true
}

// startTagQualifiedName reads the element name of a node's start tag, prefix
// included, so the matching end tag can be found without assuming "p:".
func startTagQualifiedName(nodeXML []byte) (string, bool) {
	if len(nodeXML) < 2 || nodeXML[0] != '<' {
		return "", false
	}
	for i := 1; i < len(nodeXML); i++ {
		switch nodeXML[i] {
		case ' ', '\t', '\r', '\n', '>', '/':
			if i == 1 {
				return "", false
			}
			return string(nodeXML[1:i]), true
		}
	}
	return "", false
}
