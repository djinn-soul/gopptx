package shapes

import (
	"math"
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
)

func toXMLTextFrameSpec(tf *TextFrame) *pptxxml.TextFrameSpec {
	var rotation *int64
	if tf.RotationDeg != nil {
		value := int64(math.Round(*tf.RotationDeg * float64(ooxmlAngleUnitsPerDegree)))
		rotation = &value
	}
	return &pptxxml.TextFrameSpec{
		MarginLeft:   tf.MarginLeft.Emu(),
		MarginRight:  tf.MarginRight.Emu(),
		MarginTop:    tf.MarginTop.Emu(),
		MarginBottom: tf.MarginBottom.Emu(),
		Anchor:       string(tf.Anchor),
		Wrap:         string(tf.Wrap),
		AutoFit:      string(tf.AutoFit),
		Orientation:  tf.Orientation,
		NumCol:       tf.Columns,
		Rotation:     rotation,
	}
}

func resolveActionSpec(primary, secondary *action.Hyperlink, rids map[*action.Hyperlink]string) *pptxxml.HyperlinkSpec {
	h := primary
	if h == nil {
		h = secondary
	}
	if h == nil {
		return nil
	}

	spec := &pptxxml.HyperlinkSpec{
		Tooltip:        h.Tooltip,
		HighlightClick: h.HighlightClick,
		History:        h.History,
		EndSound:       h.EndSound,
		Action:         h.ActionType(),
	}
	if rid, ok := rids[h]; ok {
		spec.RelID = rid
	}
	if spec.RelID == "" && strings.TrimSpace(spec.Tooltip) == "" && strings.TrimSpace(spec.Action) == "" &&
		spec.History == nil && spec.EndSound == nil && !spec.HighlightClick {
		return nil
	}
	return spec
}
