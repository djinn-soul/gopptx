package text

import (
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
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

// NormalizeRuns removes empty runs and merges adjacent runs with identical styling.
func NormalizeRuns(runs []Run) []Run {
	if len(runs) == 0 {
		return nil
	}
	result := make([]Run, 0, len(runs))
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

// NormalizeParagraphStyle ensures all fields are within expected bounds.
func NormalizeParagraphStyle(style ParagraphStyle) ParagraphStyle {
	normalizedBulletStyle := NormalizeBulletStyle(style.BulletStyle)
	if normalizedBulletStyle == "" {
		normalizedBulletStyle = BulletStyleBullet
	}
	return ParagraphStyle{
		Align:          NormalizeTextAlign(style.Align),
		SpaceBeforePt:  style.SpaceBeforePt,
		SpaceAfterPt:   style.SpaceAfterPt,
		LineSpacingPct: style.LineSpacingPct,
		LineSpacingPts: style.LineSpacingPts,
		BulletStyle:    normalizedBulletStyle,
		BulletChar:     strings.TrimSpace(style.BulletChar),
		BulletColor:    style.BulletColor,
		BulletSize:     style.BulletSize,
		TabStops:       append([]styling.Length(nil), style.TabStops...),
		Level:          style.Level,
	}
}

// RunsToPlainText converts a slice of Runs to a single string.
func RunsToPlainText(runs []Run) string {
	var sb strings.Builder
	for _, run := range runs {
		sb.WriteString(run.Text)
	}
	return sb.String()
}
