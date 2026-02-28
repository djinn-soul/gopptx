package text

import (
	"testing"
	"github.com/djinn-soul/gopptx/pkg/pptx/action"
)

func TestText_RunCreation(t *testing.T) {
	r := NewRun("Hello").
		WithBold(true).
		WithItalic(true).
		WithUnderline(true).
		WithColor("FF0000").
		WithSizePt(12).
		WithStrikethrough(true).
		WithSubscript(true).
		WithHighlight("00FF00").
		WithFont("Arial").
		WithCode(true).
		WithAllCaps(true).
		WithHyperlink(action.Hyperlink{}).
		WithHoverAction(action.Hyperlink{})

	if !r.Bold || !r.Italic || r.Underline != UnderlineStyleSingle || r.Color != "FF0000" || r.SizePt != 12 {
		t.Error("Run creation basic properties failed")
	}
	if r.Strikethrough != StrikethroughStyleSingle {
		t.Error("Strikethrough failed")
	}
	if !r.Subscript || r.Superscript {
		t.Error("Subscript/Superscript logic failed")
	}
	if r.Font != "Arial" || !r.Code || !r.AllCaps {
		t.Error("Run creation extra properties failed")
	}

	// Test mutually exclusive
	r = r.WithSuperscript(true)
	if !r.Superscript || r.Subscript {
		t.Error("Superscript override failed")
	}
	r = r.WithSmallCaps(true)
	if !r.SmallCaps || r.AllCaps {
		t.Error("SmallCaps override failed")
	}

	r = r.WithUnderlineStyle(UnderlineStyleDouble)
	if r.Underline != UnderlineStyleDouble { t.Error("UnderlineStyle failed") }
	r = r.WithStrikethroughStyle(StrikethroughStyleDouble)
	if r.Strikethrough != StrikethroughStyleDouble { t.Error("StrikethroughStyle failed") }
}

func TestText_ParagraphCreation(t *testing.T) {
	style := NewParagraphStyle().
		WithAlign(TextAlignCenter).
		WithNumbered().
		WithLevel(1).
		WithBulletColor("0000FF").
		WithBulletSize(80).
		WithSpaceBeforePt(10).
		WithSpaceAfterPt(10).
		WithLineSpacingPct(120).
		WithRTL(true)

	p := NewParagraph().
		WithStyle(style).
		AddRun(NewRun("R1"))

	if p.Style.Align != TextAlignCenter || p.Style.BulletStyle != BulletStyleNumber || len(p.Runs) != 1 {
		t.Error("Paragraph creation failed")
	}

	// Test other alignments
	s2 := NewParagraphStyle().WithAlignLeft().WithAlignRight().WithAlignJustify()
	if s2.Align != TextAlignJustify { t.Error("Align chaining failed") }

	s3 := NewParagraphStyle().WithLetteredLower().WithRomanUpper().WithNoBullet().WithCustomBullet("*")
	if s3.BulletStyle != BulletStyleCustom || s3.BulletChar != "*" { t.Error("Bullet style failed") }
}

func TestText_Validate(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		if err := NewRun("X").WithSizePt(-1).Validate(); err == nil {
			t.Error("expected error for negative size")
		}
		if err := NewRun("X").WithColor("invalid").Validate(); err == nil {
			t.Error("expected error for invalid color")
		}
	})

	t.Run("ParagraphStyle", func(t *testing.T) {
		s := NewParagraphStyle()
		if err := s.WithBulletColor("invalid").Validate(); err == nil {
			t.Error("expected error for invalid bullet color")
		}
		s = NewParagraphStyle()
		if err := s.WithSpaceBeforePt(-1).Validate(); err == nil {
			t.Error("expected error for negative space")
		}
		s.Align = "invalid"
		if err := s.Validate(); err == nil {
			t.Error("expected error for invalid align")
		}
	})
}
