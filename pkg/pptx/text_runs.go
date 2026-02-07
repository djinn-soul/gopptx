package pptx

import "strings"

// TextRun represents one inline text segment with optional formatting.
type TextRun struct {
	Text          string
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	Subscript     bool
	Superscript   bool
	Color         string
	Highlight     string
	Font          string
	SizePt        int
	Code          bool
}

// NewTextRun creates one text run with plain styling.
func NewTextRun(text string) TextRun {
	return TextRun{Text: text}
}

// WithBold sets bold style for this run.
func (r TextRun) WithBold(enabled bool) TextRun {
	r.Bold = enabled
	return r
}

// WithItalic sets italic style for this run.
func (r TextRun) WithItalic(enabled bool) TextRun {
	r.Italic = enabled
	return r
}

// WithUnderline sets underline style for this run.
func (r TextRun) WithUnderline(enabled bool) TextRun {
	r.Underline = enabled
	return r
}

// WithStrikethrough sets strikethrough style for this run.
func (r TextRun) WithStrikethrough(enabled bool) TextRun {
	r.Strikethrough = enabled
	return r
}

// WithSubscript sets subscript baseline for this run.
func (r TextRun) WithSubscript(enabled bool) TextRun {
	r.Subscript = enabled
	if enabled {
		r.Superscript = false
	}
	return r
}

// WithSuperscript sets superscript baseline for this run.
func (r TextRun) WithSuperscript(enabled bool) TextRun {
	r.Superscript = enabled
	if enabled {
		r.Subscript = false
	}
	return r
}

// WithColor sets RGB text color for this run.
func (r TextRun) WithColor(color string) TextRun {
	r.Color = normalizeHexColor(color)
	return r
}

// WithHighlight sets RGB highlight color for this run.
func (r TextRun) WithHighlight(color string) TextRun {
	r.Highlight = normalizeHexColor(color)
	return r
}

// WithFont sets typeface for this run.
func (r TextRun) WithFont(font string) TextRun {
	r.Font = strings.TrimSpace(font)
	return r
}

// WithSizePt sets font size in points for this run.
func (r TextRun) WithSizePt(size int) TextRun {
	r.SizePt = size
	return r
}

// WithCode toggles code style for this run.
func (r TextRun) WithCode(enabled bool) TextRun {
	r.Code = enabled
	return r
}

// AddBulletRuns appends one bullet built from rich text runs.
func (s SlideContent) AddBulletRuns(runs []TextRun) SlideContent {
	normalized := normalizeTextRuns(runs)
	s.Bullets = append(s.Bullets, runsToPlainText(normalized))
	s.BulletRuns = append(s.BulletRuns, normalized)
	s.BulletStyles = append(s.BulletStyles, s.DefaultBulletStyle)
	return s
}

func runsToPlainText(runs []TextRun) string {
	var b strings.Builder
	for _, run := range runs {
		b.WriteString(run.Text)
	}
	return b.String()
}

func normalizeTextRuns(runs []TextRun) []TextRun {
	if len(runs) == 0 {
		return nil
	}
	out := make([]TextRun, 0, len(runs))
	for _, run := range runs {
		text := strings.TrimSpace(run.Text)
		if text == "" {
			continue
		}
		out = append(out, TextRun{
			Text:          run.Text,
			Bold:          run.Bold,
			Italic:        run.Italic,
			Underline:     run.Underline,
			Strikethrough: run.Strikethrough,
			Subscript:     run.Subscript,
			Superscript:   run.Superscript,
			Color:         normalizeHexColor(run.Color),
			Highlight:     normalizeHexColor(run.Highlight),
			Font:          strings.TrimSpace(run.Font),
			SizePt:        run.SizePt,
			Code:          run.Code,
		})
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
