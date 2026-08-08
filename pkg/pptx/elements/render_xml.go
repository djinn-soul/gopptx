package elements

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/textspec"
)

func BuildSlideHyperlinkRels(
	slide SlideContent,
	firstRID int,
) (map[*action.Hyperlink]string, []pptxxml.HyperlinkRel, int) {
	hyperlinkRIDs := make(map[*action.Hyperlink]string)
	hyperlinks := make([]pptxxml.HyperlinkRel, 0)
	nextRID := firstRID

	// Dedup by value (target + action type), not pointer identity.
	seen := make(map[string]string) // value-key -> rId

	addHyperlink := func(h *action.Hyperlink) {
		if h == nil {
			return
		}
		if !h.RequiresRelationship() {
			return
		}

		valueKey := h.RelationshipTarget() + "|" + string(h.Action.Type)
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
			Target:   h.RelationshipTarget(),
			External: h.Action.IsExternal(),
			Type:     hyperlinkRelationshipType(h.Action.Type),
		})
	}

	for _, shape := range slide.Shapes {
		addShapeHyperlinks(shape, addHyperlink)
	}
	for _, connector := range slide.Connectors {
		addHyperlink(connector.ClickAction)
		addHyperlink(connector.HoverAction)
	}
	for _, runRow := range slide.BulletRuns {
		for _, run := range runRow {
			addHyperlink(run.Hyperlink)
			addHyperlink(run.HoverAction)
		}
	}

	return hyperlinkRIDs, hyperlinks, nextRID
}

func hyperlinkRelationshipType(actionType action.HyperlinkActionType) string {
	if actionType == action.HyperlinkActionSlide {
		return "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	}
	return "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
}

func ToXMLTextRunRows(rows [][]Run, hyperlinkRIDs map[*action.Hyperlink]string) [][]pptxxml.TextRunSpec {
	if len(rows) == 0 {
		return nil
	}
	out := make([][]pptxxml.TextRunSpec, len(rows))
	for i := range rows {
		if len(rows[i]) == 0 {
			continue
		}
		out[i] = textspec.ToXMLRunSpecs(rows[i], hyperlinkRIDs)
	}
	return out
}

func ToXMLBulletParagraphStyles(styles []ParagraphStyle) []pptxxml.BulletParagraphSpec {
	if len(styles) == 0 {
		return nil
	}
	out := make([]pptxxml.BulletParagraphSpec, len(styles))
	for i, style := range styles {
		out[i] = textspec.ToXMLParagraphStyleSpec(style)
	}
	return out
}

func ToXMLBackgroundSpec(bg *SlideBackground, imageRelID string) *pptxxml.SlideBackgroundSpec {
	if bg == nil || bg.Type == "" {
		return nil
	}
	spec := &pptxxml.SlideBackgroundSpec{
		Type: string(bg.Type),
	}
	switch bg.Type {
	case SlideBackgroundSolid:
		if bg.SolidFill != nil {
			spec.SolidFill = &pptxxml.ShapeFillSpec{
				Color:        common.NormalizeHexColor(bg.SolidFill.Color),
				Transparency: bg.SolidFill.Transparency,
			}
		}
	case SlideBackgroundGradient:
		if bg.GradientFill != nil {
			stops := make([]pptxxml.ShapeGradientStopSpec, 0, len(bg.GradientFill.Stops))
			for _, stop := range bg.GradientFill.Stops {
				stops = append(stops, pptxxml.ShapeGradientStopSpec{
					PositionPct:  stop.PositionPct,
					Color:        common.NormalizeHexColor(stop.Color),
					Transparency: stop.Transparency,
				})
			}
			spec.GradientFill = &pptxxml.ShapeGradientFillSpec{
				Type:     bg.GradientFill.Type,
				Stops:    stops,
				AngleDeg: bg.GradientFill.AngleDeg,
			}
		}
	case SlideBackgroundPicture:
		if bg.PictureFill != nil {
			spec.PictureFill = &pptxxml.ImageRef{
				RelID: imageRelID,
			}
		}
	}
	return spec
}

func MapTxStyles(styles *TxStyles) *pptxxml.TxStylesSpec {
	if styles == nil {
		return nil
	}
	return &pptxxml.TxStylesSpec{
		TitleStyle: MapTextLevelStyles(styles.TitleStyle),
		BodyStyle:  MapTextLevelStyles(styles.BodyStyle),
		OtherStyle: MapTextLevelStyles(styles.OtherStyle),
	}
}

func MapTextLevelStyles(levels []TextLevelStyle) []pptxxml.TextLevelStyle {
	if len(levels) == 0 {
		return nil
	}
	out := make([]pptxxml.TextLevelStyle, len(levels))
	for i, lvl := range levels {
		out[i] = pptxxml.TextLevelStyle{
			Level:      lvl.Level,
			Font:       lvl.Font,
			SizePt:     lvl.SizePt,
			Bold:       lvl.Bold,
			Italic:     lvl.Italic,
			Color:      common.NormalizeHexColor(lvl.Color),
			BulletChar: lvl.BulletChar,
			IndentEMU:  lvl.IndentEMU,
		}
	}
	return out
}

func MapNotesMasterToSpec(master *NotesMaster, backgroundRID string) *pptxxml.NotesMasterSpec {
	if master == nil {
		return nil
	}
	spec := &pptxxml.NotesMasterSpec{
		HeaderText:   master.HeaderText,
		FooterText:   master.FooterText,
		ShowDateTime: master.ShowDateTime,
		ShowSlideNum: master.ShowSlideNum,
		NotesStyle:   MapTextLevelStyles(master.BodyStyle),
	}
	if master.Background != nil {
		spec.Background = ToXMLBackgroundSpec(master.Background, backgroundRID)
	}
	return spec
}

// addShapeHyperlinks registers every link one shape can carry: its own click
// and hover actions, and the links on the runs of its rich text.
//
// The rich-text runs were missed entirely before, yet the renderer hands them
// the rID map to look themselves up in, so each link emitted an <a:hlinkClick>
// with no r:id and did nothing.
func addShapeHyperlinks(shape shapes.Shape, addHyperlink func(*action.Hyperlink)) {
	if shape.ClickAction != nil {
		addHyperlink(shape.ClickAction)
	} else if shape.Hyperlink != nil {
		addHyperlink(shape.Hyperlink)
	}
	addHyperlink(shape.HoverAction)

	for _, paragraph := range shape.TextParagraphs {
		for i := range paragraph.Runs {
			addHyperlink(paragraph.Runs[i].Hyperlink)
			addHyperlink(paragraph.Runs[i].HoverAction)
		}
	}
}
