package pptxxml

import "github.com/djinn-soul/gopptx/pkg/pptx/text"

func convertNotesStyle(s text.ParagraphStyle) BulletParagraphSpec {
	tabStops := make([]int64, 0, len(s.TabStops))
	for _, stop := range s.TabStops {
		tabStops = append(tabStops, stop.Emu())
	}
	return BulletParagraphSpec{
		Align:          s.Align,
		SpaceBeforePt:  s.SpaceBeforePt,
		SpaceAfterPt:   s.SpaceAfterPt,
		LineSpacingPct: s.LineSpacingPct,
		LineSpacingPts: s.LineSpacingPts,
		BulletStyle:    s.BulletStyle,
		BulletChar:     s.BulletChar,
		BulletColor:    s.BulletColor,
		BulletSize:     s.BulletSize,
		TabStops:       tabStops,
		Level:          s.Level,
		LeftIndent:     int64(s.LeftIndent),
		RightIndent:    int64(s.RightIndent),
		HangingIndent:  int64(s.HangingIndent),
	}
}

func convertNotesRun(r text.Run) TextRunSpec {
	return TextRunSpec{
		Text:          r.Text,
		Bold:          r.Bold,
		Italic:        r.Italic,
		Underline:     r.Underline,
		Strikethrough: r.Strikethrough,
		Subscript:     r.Subscript,
		Superscript:   r.Superscript,
		Color:         r.Color,
		Highlight:     r.Highlight,
		Font:          r.Font,
		SizePt:        r.SizePt,
		Code:          r.Code,
		AllCaps:       r.AllCaps,
		SmallCaps:     r.SmallCaps,
		Lang:          r.Lang,
	}
}
