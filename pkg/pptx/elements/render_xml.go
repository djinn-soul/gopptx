package elements

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

func BuildSlideHyperlinkRels(slide SlideContent, firstRID int) (map[*action.Hyperlink]string, []pptxxml.HyperlinkRel, int) {
	hyperlinkRIDs := make(map[*action.Hyperlink]string)
	hyperlinks := make([]pptxxml.HyperlinkRel, 0)
	nextRID := firstRID

	// Dedup by value (target + action type), not pointer identity.
	seen := make(map[string]string) // value-key -> rId

	addHyperlink := func(h *action.Hyperlink) {
		if h == nil {
			return
		}

		valueKey := h.Action.RelationshipTarget() + "|" + string(h.Action.Type)
		if rid, exists := seen[valueKey]; exists {
			// Same value via a different pointer: reuse the existing RID.
			hyperlinkRIDs[h] = rid
			return
		}

		// Already mapped this exact pointer (redundant safety check).
		if _, exists := hyperlinkRIDs[h]; exists {
			return
		}

		rid := fmt.Sprintf("rId%d", nextRID)
		seen[valueKey] = rid
		hyperlinkRIDs[h] = rid
		nextRID++

		hyperlinks = append(hyperlinks, pptxxml.HyperlinkRel{
			RID:      rid,
			Target:   h.Action.RelationshipTarget(),
			External: h.Action.IsExternal(),
		})
	}

	for _, shape := range slide.Shapes {
		addHyperlink(shape.Hyperlink)
	}
	for _, runRow := range slide.BulletRuns {
		for _, run := range runRow {
			addHyperlink(run.Hyperlink)
		}
	}

	return hyperlinkRIDs, hyperlinks, nextRID
}

func ToXMLTextRunRows(rows [][]TextRun, hyperlinkRIDs map[*action.Hyperlink]string) [][]pptxxml.TextRunSpec {
	if len(rows) == 0 {
		return nil
	}
	out := make([][]pptxxml.TextRunSpec, len(rows))
	for i := range rows {
		if len(rows[i]) == 0 {
			continue
		}
		runs := make([]pptxxml.TextRunSpec, 0, len(rows[i]))
		for _, run := range rows[i] {
			spec := pptxxml.TextRunSpec{
				Text:          run.Text,
				Bold:          run.Bold,
				Italic:        run.Italic,
				Underline:     run.Underline,
				Strikethrough: run.Strikethrough,
				Subscript:     run.Subscript,
				Superscript:   run.Superscript,
				Color:         common.NormalizeHexColor(run.Color),
				Highlight:     common.NormalizeHexColor(run.Highlight),
				Font:          run.Font,
				SizePt:        run.SizePt,
				Code:          run.Code,
			}
			if run.Hyperlink != nil {
				if rid, ok := hyperlinkRIDs[run.Hyperlink]; ok {
					spec.Hyperlink = &pptxxml.HyperlinkSpec{
						RelID:          rid,
						Tooltip:        run.Hyperlink.Tooltip,
						HighlightClick: run.Hyperlink.HighlightClick,
						Action:         run.Hyperlink.Action.ActionType(),
					}
				}
			}
			runs = append(runs, spec)
		}
		out[i] = runs
	}
	return out
}

func ToXMLBulletParagraphStyles(styles []TextParagraphStyle) []pptxxml.BulletParagraphSpec {
	if len(styles) == 0 {
		return nil
	}
	out := make([]pptxxml.BulletParagraphSpec, len(styles))
	for i, style := range styles {
		out[i] = pptxxml.BulletParagraphSpec{
			Align:          text.NormalizeTextAlign(style.Align),
			SpaceBeforePt:  style.SpaceBeforePt,
			SpaceAfterPt:   style.SpaceAfterPt,
			LineSpacingPct: style.LineSpacingPct,
			BulletStyle:    text.NormalizeBulletStyle(style.BulletStyle),
			BulletChar:     style.BulletChar,
			BulletColor:    common.NormalizeHexColor(style.BulletColor),
			BulletSize:     style.BulletSize,
			Level:          style.Level,
		}
	}
	return out
}
