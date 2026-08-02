package pptxxml

import (
	"math"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/text"
)

const (
	defaultBulletRunSize = 2800
	// ptToEMU converts points to English Metric Units (1 pt = 12700 EMU).
	ptToEMU = emuPerPoint
)

// TextRunSpec describes one rich text run in a bullet paragraph.
type TextRunSpec struct {
	Text           string
	Bold           bool
	Italic         bool
	Underline      string
	Strikethrough  string
	Subscript      bool
	Superscript    bool
	Color          string
	Highlight      string
	Font           string
	SizePt         float64
	Code           bool
	AllCaps        bool
	SmallCaps      bool
	OutlineColor   string         // Character outline (stroke) color hex
	OutlineWidthPt float64        // Character outline width in points
	Lang           string         // BCP-47 language tag, e.g. "ar-SA", "fr-FR" (defaults to "en-US")
	Hyperlink      *HyperlinkSpec // Click action
	HoverAction    *HyperlinkSpec // Hover action
}

func bulletRunsAt(allRuns [][]TextRunSpec, index int) []TextRunSpec {
	if len(allRuns) == 0 || index < 0 || index >= len(allRuns) {
		return nil
	}
	return allRuns[index]
}

func bulletParagraphRuns(runs []TextRunSpec, style BulletParagraphSpec, contentStyle ContentStyleSpec) string {
	var b strings.Builder
	b.WriteString(`<a:p>` + BulletParagraphPropsXML(style))
	for _, run := range runs {
		if strings.TrimSpace(run.Text) == "" {
			continue
		}
		b.WriteString(richTextRun(run, contentStyle))
	}
	b.WriteString(`</a:p>`)
	return b.String()
}

func richTextRunLang(run TextRunSpec) string {
	if lang := strings.TrimSpace(run.Lang); lang != "" {
		return lang
	}
	return langEnUS
}

func richTextRun(run TextRunSpec, contentStyle ContentStyleSpec) string {
	var b strings.Builder
	b.WriteString(`<a:r><a:rPr lang="`)
	b.WriteString(Escape(richTextRunLang(run)))
	b.WriteString(`" sz="`)
	b.WriteString(runSizeValueWithDefault(run.SizePt, contentStyle.SizePt))
	b.WriteString(`" b="`)
	b.WriteString(boolToFlag(run.Bold || contentStyle.Bold))
	b.WriteString(`" i="`)
	b.WriteString(boolToFlag(run.Italic || contentStyle.Italic))
	b.WriteString(`" u="`)
	b.WriteString(runUnderlineValue(run.Underline, contentStyle.Underline))
	b.WriteString(`"`)
	if run.Strikethrough != "" && run.Strikethrough != valNone {
		b.WriteString(` strike="`)
		b.WriteString(Escape(run.Strikethrough))
		b.WriteString(`"`)
	}
	if run.Subscript {
		b.WriteString(` baseline="-25000"`)
	} else if run.Superscript {
		b.WriteString(` baseline="30000"`)
	}
	if run.AllCaps {
		b.WriteString(` cap="all"`)
	} else if run.SmallCaps {
		b.WriteString(` cap="small"`)
	}
	b.WriteString(` dirty="0">`)

	if highlight := strings.TrimSpace(run.Highlight); highlight != "" {
		b.WriteString(`<a:highlight><a:srgbClr val="`)
		b.WriteString(Escape(highlight))
		b.WriteString(`"/></a:highlight>`)
	}

	color := strings.TrimSpace(run.Color)
	if color == "" {
		color = strings.TrimSpace(contentStyle.Color)
	}

	if color != "" {
		b.WriteString(`<a:solidFill><a:srgbClr val="`)
		b.WriteString(Escape(color))
		b.WriteString(`"/></a:solidFill>`)
	}
	b.WriteString(runTypefaceXML(runFont(run), run.Lang))
	b.WriteString(richTextRunOutlineXML(run))
	if run.Hyperlink != nil {
		// A run that names its own colour keeps it: without the override
		// PowerPoint repaints the whole run in the theme's hyperlink colour
		// (upstream issue #940).
		link := *run.Hyperlink
		link.UseTextColor = link.UseTextColor || color != ""
		b.WriteString(HyperlinkXML(link, "a:hlinkClick"))
	}
	if run.HoverAction != nil {
		b.WriteString(HyperlinkXML(*run.HoverAction, "a:hlinkMouseOver"))
	}

	b.WriteString(`</a:rPr><a:t>`)
	b.WriteString(Escape(run.Text))
	b.WriteString(`</a:t></a:r>`)
	return b.String()
}

// RichTextRunXML renders one <a:r> node for a text run.
func RichTextRunXML(run TextRunSpec, contentStyle ContentStyleSpec) string {
	return richTextRun(run, contentStyle)
}

// runTypefaceXML writes the font into the typeface slot the run's language is
// actually rendered from. A Japanese or Arabic run whose font only appears in
// <a:latin> is drawn in the theme's fallback face instead, because PowerPoint
// picks <a:ea> for East Asian text and <a:cs> for complex scripts
// (upstream issue #172).
func runTypefaceXML(font, langTag string) string {
	trimmed := strings.TrimSpace(font)
	if trimmed == "" {
		return ""
	}
	escaped := Escape(trimmed)
	latin := `<a:latin typeface="` + escaped + `"/>`
	switch text.ScriptKindForLanguage(langTag) {
	case text.ScriptEastAsian:
		return latin + `<a:ea typeface="` + escaped + `"/>`
	case text.ScriptComplex:
		return latin + `<a:cs typeface="` + escaped + `"/>`
	case text.ScriptLatin:
		return latin
	default:
		return latin
	}
}

func richTextRunOutlineXML(run TextRunSpec) string {
	outlineColor := strings.TrimSpace(run.OutlineColor)
	if outlineColor == "" {
		return ""
	}
	widthEMU := int64(ptToEMU) // default 1pt
	if run.OutlineWidthPt > 0 {
		widthEMU = int64(run.OutlineWidthPt * ptToEMU)
	}
	var b strings.Builder
	b.WriteString(`<a:ln w="`)
	b.WriteString(strconv.FormatInt(widthEMU, 10))
	b.WriteString(`"><a:solidFill><a:srgbClr val="`)
	b.WriteString(Escape(outlineColor))
	b.WriteString(`"/></a:solidFill></a:ln>`)
	return b.String()
}

// runSizeValueWithDefault renders the a:rPr sz attribute in centipoints. The run
// size is a float so half-point sizes survive a read/write round-trip; it is
// rounded to the nearest centipoint because sz is an integer attribute.
func runSizeValueWithDefault(sizePt float64, defaultSizePt int) string {
	if sizePt > 0 {
		return strconv.Itoa(int(math.Round(sizePt * float64(ptFactor))))
	}
	if defaultSizePt > 0 {
		return strconv.Itoa(defaultSizePt * ptFactor)
	}
	return strconv.Itoa(defaultBulletRunSize)
}

func runFont(run TextRunSpec) string {
	if strings.TrimSpace(run.Font) != "" {
		return run.Font
	}
	if run.Code {
		return "Consolas"
	}
	return ""
}

const valNone = "none"

func runUnderlineValue(underline string, defaultUnderline bool) string {
	if underline != "" && underline != valNone {
		switch strings.ToLower(strings.TrimSpace(underline)) {
		case "single":
			return "sng"
		case "double":
			return "dbl"
		default:
			return underline
		}
	}
	if defaultUnderline {
		return "sng"
	}
	return valNone
}

func boolToFlag(enabled bool) string {
	if enabled {
		return "1"
	}
	return "0"
}
