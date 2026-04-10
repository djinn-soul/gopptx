package pptxxml

import "strconv"

func splitBulletsForTwoColumns(bullets []string) ([]string, []string) {
	if len(bullets) == 0 {
		return nil, nil
	}
	mid := (len(bullets) + 1) / 2
	return bullets[:mid], bullets[mid:]
}

func splitBulletStylesForTwoColumns(
	styles []BulletParagraphSpec,
	leftCount int,
) ([]BulletParagraphSpec, []BulletParagraphSpec) {
	if len(styles) == 0 {
		return nil, nil
	}
	if leftCount > len(styles) {
		leftCount = len(styles)
	}
	return styles[:leftCount], styles[leftCount:]
}

func splitBulletRunsForTwoColumns(runs [][]TextRunSpec, leftCount int) ([][]TextRunSpec, [][]TextRunSpec) {
	if len(runs) == 0 {
		return nil, nil
	}
	if leftCount > len(runs) {
		leftCount = len(runs)
	}
	return runs[:leftCount], runs[leftCount:]
}

// defaultBulletParagraphPrefix is the precomputed XML before the bullet text
// for zero-value BulletParagraphSpec and ContentStyleSpec (sz=2800, no color, no bold/italic).
const defaultBulletParagraphPrefix = "\n<a:p>\n" + defaultBulletParagraphProps + "\n<a:r>\n" +
	`<a:rPr lang="en-US" sz="2800" b="0" i="0" u="none" dirty="0"></a:rPr>` + "\n<a:t>"

const defaultBulletParagraphSuffix = "</a:t>\n</a:r>\n</a:p>"

func bulletParagraph(text string, pStyle BulletParagraphSpec, style ContentStyleSpec) string {
	if pStyle.IsZero() && style == (ContentStyleSpec{}) {
		return defaultBulletParagraphPrefix + Escape(text) + defaultBulletParagraphSuffix
	}

	escaped := Escape(text)
	sz := 2800
	if style.SizePt > 0 {
		sz = style.SizePt * 100 //nolint:mnd // points->centipoints
	}
	colorXML := ""
	if style.Color != "" {
		colorXML = `<a:solidFill><a:srgbClr val="` + Escape(style.Color) + `"/></a:solidFill>`
	}

	return `
<a:p>
` + BulletParagraphPropsXML(pStyle) + `
<a:r>
<a:rPr lang="en-US" sz="` + strconv.Itoa(sz) + `" b="` + boolToFlag(style.Bold) + `" i="` + boolToFlag(style.Italic) + `" u="` + runUnderlineValue("", style.Underline) + `" dirty="0">` + colorXML + `</a:rPr>
<a:t>` + escaped + `</a:t>
</a:r>
</a:p>`
}
