package text

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
)

// hyperlinksEqual compares two *action.Hyperlink values by content, not pointer identity.
func hyperlinksEqual(a, b *action.Hyperlink) bool {
	if a == b {
		return true // both nil, or same pointer
	}
	if a == nil || b == nil {
		return false
	}
	return a.Action == b.Action &&
		a.Tooltip == b.Tooltip &&
		a.HighlightClick == b.HighlightClick
}

// NormalizeTextRuns removes empty runs and merges adjacent runs with identical styling.
func NormalizeTextRuns(runs []TextRun) []TextRun {
	if len(runs) == 0 {
		return nil
	}
	result := make([]TextRun, 0, len(runs))
	for _, run := range runs {
		if run.Text == "" {
			continue
		}
		if len(result) > 0 {
			last := &result[len(result)-1]
			if last.Bold == run.Bold && last.Italic == run.Italic && last.Code == run.Code &&
				last.Color == run.Color && last.SizePt == run.SizePt && last.Underline == run.Underline &&
				last.Strikethrough == run.Strikethrough && last.Subscript == run.Subscript &&
				last.Superscript == run.Superscript && last.Highlight == run.Highlight &&
				last.Font == run.Font && hyperlinksEqual(last.Hyperlink, run.Hyperlink) {
				last.Text += run.Text
				continue
			}
		}
		result = append(result, run)
	}
	return result
}

// NormalizeTextParagraphStyle ensures all fields are within expected bounds.
func NormalizeTextParagraphStyle(style TextParagraphStyle) TextParagraphStyle {
	normalizedBulletStyle := NormalizeBulletStyle(style.BulletStyle)
	if normalizedBulletStyle == "" {
		normalizedBulletStyle = BulletStyleBullet
	}
	return TextParagraphStyle{
		Align:          NormalizeTextAlign(style.Align),
		SpaceBeforePt:  style.SpaceBeforePt,
		SpaceAfterPt:   style.SpaceAfterPt,
		LineSpacingPct: style.LineSpacingPct,
		BulletStyle:    normalizedBulletStyle,
		BulletChar:     strings.TrimSpace(style.BulletChar),
		BulletColor:    style.BulletColor,
		BulletSize:     style.BulletSize,
		Level:          style.Level,
	}
}

// RunsToPlainText converts a slice of TextRuns to a single string.
func RunsToPlainText(runs []TextRun) string {
	var sb strings.Builder
	for _, run := range runs {
		sb.WriteString(run.Text)
	}
	return sb.String()
}
