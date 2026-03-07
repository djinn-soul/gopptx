package shape

import (
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func applyParsedShapeText(ps *ParsedShapeProperties, s *shapeXML) {
	var txt strings.Builder
	for pIdx, paragraph := range s.TxBody.P {
		applyParsedParagraph(ps, pIdx, paragraph.PPr)
		for _, runXML := range paragraph.R {
			txt.WriteString(runXML.T)
			if run, ok := parseTextRun(runXML); ok {
				ps.Runs = append(ps.Runs, run)
			}
		}
		if pIdx < len(s.TxBody.P)-1 {
			txt.WriteString("\n")
		}
	}
	ps.Text = txt.String()
}

func applyParsedParagraph(ps *ParsedShapeProperties, index int, pPr *struct {
	MarL   *int `xml:"marL,attr"`
	Indent *int `xml:"indent,attr"`
	TabLst *struct {
		Tabs []struct {
			Pos *int `xml:"pos,attr"`
		} `xml:"tab"`
	} `xml:"tabLst"`
}) {
	if index != 0 || pPr == nil {
		return
	}
	paragraph := &common.Paragraph{}
	if pPr.MarL != nil {
		paragraph.Indent = pPr.MarL
	}
	if pPr.Indent != nil && *pPr.Indent < 0 {
		hanging := -*pPr.Indent
		paragraph.Hanging = &hanging
	}
	if pPr.TabLst != nil {
		tabStops := make([]int, 0, len(pPr.TabLst.Tabs))
		for _, tab := range pPr.TabLst.Tabs {
			if tab.Pos == nil {
				continue
			}
			tabStops = append(tabStops, *tab.Pos)
		}
		if len(tabStops) > 0 {
			paragraph.TabStops = tabStops
		}
	}
	if paragraph.Indent != nil || paragraph.Hanging != nil || len(paragraph.TabStops) > 0 {
		ps.Paragraph = paragraph
	}
}

func parseTextRun(runXML struct {
	RPr *runPropsXML `xml:"rPr"`
	T   string       `xml:"t"`
}) (common.TextRun, bool) {
	if runXML.RPr == nil && runXML.T == "" {
		return common.TextRun{}, false
	}
	run := common.TextRun{Text: runXML.T}
	if runXML.RPr != nil {
		applyRunStyle(&run, runXML.RPr)
	}
	return run, true
}

func applyRunStyle(run *common.TextRun, rpr *runPropsXML) {
	if rpr.Bold != nil && *rpr.Bold {
		run.Bold = rpr.Bold
	}
	if rpr.Italic != nil && *rpr.Italic {
		run.Italic = rpr.Italic
	}
	if rpr.Underline != nil && *rpr.Underline != "" {
		run.Underline = rpr.Underline
	}
	if rpr.Strikethrough != nil && *rpr.Strikethrough != "" {
		run.Strikethrough = rpr.Strikethrough
	}
	applyRunBaseline(run, rpr)
	applyRunCaps(run, rpr)
	applyRunColors(run, rpr)
}

func applyRunBaseline(run *common.TextRun, rpr *runPropsXML) {
	switch {
	case ParseIntAttr(rpr.Baseline) < 0:
		v := true
		run.Subscript = &v
	case ParseIntAttr(rpr.Baseline) > 0:
		v := true
		run.Superscript = &v
	}
}

func applyRunCaps(run *common.TextRun, rpr *runPropsXML) {
	if rpr.Caps != nil {
		switch strings.ToLower(strings.TrimSpace(*rpr.Caps)) {
		case "all":
			v := true
			run.AllCaps = &v
		case "small":
			v := true
			run.SmallCaps = &v
		}
	}
	if ParseXMLBoolAttr(rpr.SmallCaps) {
		v := true
		run.SmallCaps = &v
	}
}

func applyRunColors(run *common.TextRun, rpr *runPropsXML) {
	if rpr.SolidFill.SrgbClr.Val != "" {
		val := rpr.SolidFill.SrgbClr.Val
		run.Color = &val
	}
	if rpr.Highlight.SrgbClr.Val != "" {
		val := rpr.Highlight.SrgbClr.Val
		run.Highlight = &val
	}
}
