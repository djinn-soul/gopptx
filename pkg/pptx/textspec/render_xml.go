// Package textspec converts the public run and paragraph model into the XML
// specs the renderer consumes. It is its own package because internal/pptxxml
// already imports pkg/pptx/text, so the conversion cannot live beside the model
// without a cycle. Slide bullets, shape text bodies and anything else carrying
// runs all render from this one mapping.
package textspec

import (
	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

// ToXMLHyperlinkSpec converts one hyperlink, resolving its relationship id from
// the map the slide builder assigned. It returns nil when the link carries
// nothing worth emitting.
func ToXMLHyperlinkSpec(h *action.Hyperlink, hyperlinkRIDs map[*action.Hyperlink]string) *pptxxml.HyperlinkSpec {
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
	if rid, ok := hyperlinkRIDs[h]; ok {
		spec.RelID = rid
	}
	if spec.RelID == "" && spec.Tooltip == "" && spec.Action == "" && spec.History == nil && spec.EndSound == nil &&
		!spec.HighlightClick {
		return nil
	}
	return spec
}

// ToXMLRunSpec converts one styled run.
func ToXMLRunSpec(run text.Run, hyperlinkRIDs map[*action.Hyperlink]string) pptxxml.TextRunSpec {
	spec := pptxxml.TextRunSpec{
		Text:           run.Text,
		Bold:           run.Bold,
		Italic:         run.Italic,
		Underline:      run.Underline,
		Strikethrough:  run.Strikethrough,
		Subscript:      run.Subscript,
		Superscript:    run.Superscript,
		Color:          common.NormalizeHexColor(run.Color),
		Highlight:      common.NormalizeHexColor(run.Highlight),
		Font:           run.Font,
		SizePt:         float64(run.SizePt),
		Code:           run.Code,
		AllCaps:        run.AllCaps,
		SmallCaps:      run.SmallCaps,
		OutlineColor:   common.NormalizeHexColor(run.OutlineColor),
		OutlineWidthPt: run.OutlineWidthPt,
		Lang:           run.Lang,
	}
	if run.Hyperlink != nil {
		spec.Hyperlink = ToXMLHyperlinkSpec(run.Hyperlink, hyperlinkRIDs)
	}
	if run.HoverAction != nil {
		spec.HoverAction = ToXMLHyperlinkSpec(run.HoverAction, hyperlinkRIDs)
	}
	return spec
}

// ToXMLRunSpecs converts a run of runs, keeping their order.
func ToXMLRunSpecs(runs []text.Run, hyperlinkRIDs map[*action.Hyperlink]string) []pptxxml.TextRunSpec {
	if len(runs) == 0 {
		return nil
	}
	out := make([]pptxxml.TextRunSpec, 0, len(runs))
	for _, run := range runs {
		out = append(out, ToXMLRunSpec(run, hyperlinkRIDs))
	}
	return out
}

// ToXMLParagraphStyleSpec converts paragraph-level formatting.
func ToXMLParagraphStyleSpec(style text.ParagraphStyle) pptxxml.BulletParagraphSpec {
	return pptxxml.BulletParagraphSpec{
		Align:          text.NormalizeTextAlign(style.Align),
		SpaceBeforePt:  style.SpaceBeforePt,
		SpaceAfterPt:   style.SpaceAfterPt,
		LineSpacingPct: style.LineSpacingPct,
		LineSpacingPts: style.LineSpacingPts,
		BulletStyle:    text.NormalizeBulletStyle(style.BulletStyle),
		BulletChar:     style.BulletChar,
		BulletColor:    common.NormalizeHexColor(style.BulletColor),
		BulletSize:     style.BulletSize,
		TabStops:       toXMLTabStops(style.TabStops),
		Level:          style.Level,
		LeftIndent:     style.LeftIndent.Emu(),
		RightIndent:    style.RightIndent.Emu(),
		HangingIndent:  style.HangingIndent.Emu(),
		RTL:            style.RTL,
	}
}

// ToXMLParagraphSpecs converts whole paragraphs — style plus runs — which is
// what a shape's text body needs, as opposed to the parallel bullet slices a
// placeholder renders from.
func ToXMLParagraphSpecs(
	paragraphs []text.Paragraph,
	hyperlinkRIDs map[*action.Hyperlink]string,
) []pptxxml.ParagraphSpec {
	if len(paragraphs) == 0 {
		return nil
	}
	out := make([]pptxxml.ParagraphSpec, 0, len(paragraphs))
	for _, p := range paragraphs {
		out = append(out, pptxxml.ParagraphSpec{
			Style: ToXMLParagraphStyleSpec(p.Style),
			Runs:  ToXMLRunSpecs(p.Runs, hyperlinkRIDs),
		})
	}
	return out
}

func toXMLTabStops(stops []styling.Length) []int64 {
	if len(stops) == 0 {
		return nil
	}
	converted := make([]int64, 0, len(stops))
	for _, stop := range stops {
		converted = append(converted, stop.Emu())
	}
	return converted
}
