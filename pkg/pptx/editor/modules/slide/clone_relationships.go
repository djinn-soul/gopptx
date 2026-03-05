package slide

import (
	"fmt"
	"path"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

type RenderRelationshipsFn func([]common.EditorRelationship) string
type RewriteChartExternalDataFn func([]byte, string) []byte

func CloneChartDependencies(
	srcChartPart string,
	newChartPart string,
	newChartData []byte,
	nextExcelNum int,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderRelationships RenderRelationshipsFn,
	rewriteChartExternalData RewriteChartExternalDataFn,
) ([]byte, int, string) {
	srcChartRelsPath := common.SlideRelsPartName(srcChartPart)
	relsData, relsOK := getPart(srcChartRelsPath)
	if !relsOK {
		return newChartData, nextExcelNum, ""
	}

	chartRels, _ := parseRelationships(relsData)
	embeddingPath := ""

	for i, cr := range chartRels {
		if cr.Type != common.RelTypePackage {
			continue
		}
		srcExcel := common.CanonicalPartPath(path.Join(path.Dir(srcChartPart), cr.Target))
		xdata, excelOK := getPart(srcExcel)
		if !excelOK {
			continue
		}

		newExcel := fmt.Sprintf("ppt/embeddings/Microsoft_Excel_Worksheet%d.xlsx", nextExcelNum)
		nextExcelNum++
		setPart(newExcel, CloneBytes(xdata))
		newChartData = rewriteChartExternalData(newChartData, cr.ID)
		chartRels[i].Target = "../embeddings/" + path.Base(newExcel)
		embeddingPath = newExcel
	}

	newChartRelsPath := common.SlideRelsPartName(newChartPart)
	rendered := renderRelationships(chartRels)
	setPart(newChartRelsPath, []byte(rendered))
	return newChartData, nextExcelNum, embeddingPath
}

func CloneNotesRelationships(
	srcNotesPart string,
	newNotesPart string,
	newSlidePart string,
	getPart GetPartFn,
	setPart SetPartFn,
	parseRelationships ParseRelationshipsFn,
	renderRelationships RenderRelationshipsFn,
) {
	srcNotesRelsPath := common.SlideRelsPartName(srcNotesPart)
	relsData, relsOK := getPart(srcNotesRelsPath)
	if !relsOK {
		return
	}

	notesRels, _ := parseRelationships(relsData)
	for i, nr := range notesRels {
		if nr.Type == common.RelTypeSlide {
			notesRels[i].Target = "../slides/" + path.Base(newSlidePart)
		}
	}

	newNotesRelsPath := common.SlideRelsPartName(newNotesPart)
	rendered := renderRelationships(notesRels)
	setPart(newNotesRelsPath, []byte(rendered))
}
