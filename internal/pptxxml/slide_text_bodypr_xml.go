package pptxxml

import "strconv"

// TextBodyPrXML renders <a:bodyPr> with the same defaults used by shape text bodies.
func TextBodyPrXML(textFrame *TextFrameSpec) string {
	autoFitXML := `<a:spAutoFit/>`
	bodyPrAttr := ` wrap="square" rtlCol="0" anchor="ctr" lIns="` + strconv.Itoa(
		defaultMargin,
	) + `" tIns="` + strconv.Itoa(
		defaultMargin,
	) + `" rIns="` + strconv.Itoa(
		defaultMargin,
	) + `" bIns="` + strconv.Itoa(
		defaultMargin,
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
		case "spAutoFit":
			autoFitXML = `<a:spAutoFit/>`
		case normAutoFitToken:
			// The public API token remains "normAutoFit", but the OOXML element
			// name is schema-valid only as <a:normAutofit/>.
			autoFitXML = `<a:normAutofit/>`
		case "none":
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
