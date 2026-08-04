package pptxxml

import (
	"math"
	"strconv"
)

// OOXML stores autofit percentages in thousandths of a percent, and bounds them
// (ECMA-376 §20.1.10.42/§20.1.10.43): the font may shrink to 25% at the least
// and line spacing may be cut by at most 20%. A value outside the range makes
// PowerPoint reject the part, so the writer clamps rather than emits it raw.
const (
	autofitPercentScale     = 1000.0
	minFontScalePercent     = 25.0
	maxFontScalePercent     = 100.0
	maxLineSpaceReductionPc = 20.0
)

// normAutofitXML renders <a:normAutofit>, carrying the shrink amounts when the
// caller supplied them. PowerPoint recomputes these when the text is edited,
// but it honours what is in the file until then, which is what makes a deck
// generated headlessly look shrunk on first open.
func normAutofitXML(textFrame *TextFrameSpec) string {
	attrs := ""
	if textFrame.FontScale != nil {
		scale := clampPercent(*textFrame.FontScale, minFontScalePercent, maxFontScalePercent)
		attrs += ` fontScale="` + strconv.Itoa(int(math.Round(scale*autofitPercentScale))) + `"`
	}
	if textFrame.LineSpaceReduction != nil {
		reduction := clampPercent(*textFrame.LineSpaceReduction, 0, maxLineSpaceReductionPc)
		attrs += ` lnSpcReduction="` + strconv.Itoa(
			int(math.Round(reduction*autofitPercentScale)),
		) + `"`
	}
	return `<a:normAutofit` + attrs + `/>`
}

func clampPercent(value, low, high float64) float64 {
	return math.Min(math.Max(value, low), high)
}

// TextBodyPrXML renders <a:bodyPr> with the same defaults used by shape text bodies.
func TextBodyPrXML(textFrame *TextFrameSpec) string {
	autoFitXML := textAutoFitElement
	bodyPrAttr := ` wrap="square" rtlCol="0" anchor="ctr" lIns="` + strconv.Itoa(
		defaultInsetLR,
	) + `" tIns="` + strconv.Itoa(
		defaultInsetTB,
	) + `" rIns="` + strconv.Itoa(
		defaultInsetLR,
	) + `" bIns="` + strconv.Itoa(
		defaultInsetTB,
	) + `"`

	if textFrame != nil {
		bodyPrAttr = ` wrap="` + Escape(
			textFrame.Wrap,
		) + `" rtlCol="0" anchor="` + Escape(
			textFrame.Anchor,
		) + `" lIns="` + strconv.FormatInt(
			textFrame.MarginLeft,
			10,
		) + `" tIns="` + strconv.FormatInt(
			textFrame.MarginTop,
			10,
		) + `" rIns="` + strconv.FormatInt(
			textFrame.MarginRight,
			10,
		) + `" bIns="` + strconv.FormatInt(
			textFrame.MarginBottom,
			10,
		) + `"`
		if textFrame.Rotation != nil {
			bodyPrAttr += ` rot="` + strconv.FormatInt(*textFrame.Rotation, 10) + `"`
		}
		if textFrame.Orientation != "" {
			bodyPrAttr += ` vert="` + Escape(textFrame.Orientation) + `"`
		}
		if textFrame.NumCol > 0 {
			bodyPrAttr += ` numCol="` + strconv.Itoa(textFrame.NumCol) + `"`
		}
		switch textFrame.AutoFit {
		case textAutoFitTag:
			autoFitXML = textAutoFitElement
		case normAutoFitToken:
			// The public API token remains "normAutoFit", but the OOXML element
			// name is schema-valid only as <a:normAutofit/>.
			autoFitXML = normAutofitXML(textFrame)
		case valNone:
			autoFitXML = `<a:noAutofit/>`
		default:
			autoFitXML = ""
		}
	}

	bodyPrChildren := autoFitXML
	if textFrame != nil && textFrame.AutoFit == normAutoFitToken {
		bodyPrChildren = `<a:prstTxWarp prst="textNoShape"><a:avLst/></a:prstTxWarp>` + "\n" + autoFitXML
	}

	return `<a:bodyPr` + bodyPrAttr + `>
` + bodyPrChildren + `
</a:bodyPr>`
}
