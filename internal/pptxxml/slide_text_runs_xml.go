package pptxxml

import (
	"strconv"
	"strings"
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
	SizePt         int
	Code           bool
	AllCaps        bool
	SmallCaps      bool
	OutlineColor   string         // Character outline (stroke) color hex
	OutlineWidthPt float64        // Character outline width in points
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
	b.WriteString(`<a:p>` + bulletParagraphPropsXML(style))
	for _, run := range runs {
		if strings.TrimSpace(run.Text) == "" {
			continue
		}
		b.WriteString(richTextRun(run, contentStyle))
	}
	b.WriteString(`</a:p>`)
	return b.String()
}

func richTextRun(run TextRunSpec, contentStyle ContentStyleSpec) string {
	var b strings.Builder
	b.WriteString(`<a:r><a:rPr lang="en-US" sz="`)
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
	if font := strings.TrimSpace(runFont(run)); font != "" {
		b.WriteString(`<a:latin typeface="`)
		b.WriteString(Escape(font))
		b.WriteString(`"/>`)
	}
	b.WriteString(richTextRunOutlineXML(run))
	if run.Hyperlink != nil {
		b.WriteString(HyperlinkXML(*run.Hyperlink, "a:hlinkClick"))
	}
	if run.HoverAction != nil {
		b.WriteString(HyperlinkXML(*run.HoverAction, "a:hlinkMouseOver"))
	}

	b.WriteString(`</a:rPr><a:t>`)
	b.WriteString(Escape(run.Text))
	b.WriteString(`</a:t></a:r>`)
	return b.String()
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

func runSizeValueWithDefault(sizePt int, defaultSizePt int) string {
	if sizePt > 0 {
		return strconv.Itoa(sizePt * ptFactor)
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
		return underline
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
