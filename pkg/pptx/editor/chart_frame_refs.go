package editor

import "fmt"

type chartFrameRef struct {
	RelID   string
	ShapeID int
}

func extractChartFrameRefs(slideXML []byte) ([]chartFrameRef, error) {
	shapes, err := parseSlideShapes(slideXML)
	if err != nil {
		return nil, err
	}
	refs := make([]chartFrameRef, 0)
	for _, shape := range shapes {
		if shape.Type != shapeTypeGraphicFrame {
			continue
		}
		if shape.Start < 0 || shape.End > int64(len(slideXML)) || shape.Start >= shape.End {
			return nil, fmt.Errorf(
				"invalid graphic frame offsets: start=%d end=%d size=%d",
				shape.Start,
				shape.End,
				len(slideXML),
			)
		}
		match := reChartRelID.FindSubmatch(slideXML[shape.Start:shape.End])
		if len(match) != 2 {
			continue
		}
		refs = append(refs, chartFrameRef{
			RelID:   string(match[1]),
			ShapeID: shape.ID,
		})
	}
	return refs, nil
}
