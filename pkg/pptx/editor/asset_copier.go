package editor

import (
	"errors"
	"fmt"
	"path"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
)

func renderRelationshipsXML(rels []common.EditorRelationship) string {
	return editorslide.RenderRelationshipsXML(rels)
}

func rewriteChartExternalData(current []byte, newRelID string) []byte {
	return editorslide.RewriteChartExternalData(current, newRelID)
}

// deepCloneSlideAssets walks through the relationships of a source slide and copies
// all referenced assets (images, charts, etc.) to the target editor.
// It returns a modified relationships XML byte slice where targets are remapped to the new locations.
func (e *PresentationEditor) deepCloneSlideAssets(
	srcEditor *PresentationEditor,
	srcSlidePart string,
	srcSlideRelsBytes []byte,
	dstSlidePart string,
) ([]byte, error) {
	rels, err := parseRelationshipsXML(srcSlideRelsBytes)
	if err != nil {
		return nil, err
	}

	changed := false
	for i, rel := range rels {
		// Determine the absolute path of the target in the source package
		// slide relationships are usually relative to ppt/slides/slideN.xml
		// e.g. target="../media/image1.png" -> ppt/media/image1.png
		srcTargetAbs := common.ResolveRelationshipTarget(srcSlidePart, rel.Target)

		var newTarget string
		var handled bool

		switch rel.Type {
		case common.RelTypeImage:
			newTarget, err = e.copyImageAsset(srcEditor, srcTargetAbs)
			handled = true
		case common.RelTypeChart:
			newTarget, err = e.copyChartAsset(srcEditor, srcTargetAbs)
			handled = true
		case common.RelTypeNotesSlide:
			newTarget, err = e.copyNotesSlideAsset(srcEditor, srcTargetAbs, dstSlidePart)
			handled = true
		}

		if err != nil {
			return nil, fmt.Errorf("failed to copy asset %s (type %s): %w", srcTargetAbs, rel.Type, err)
		}

		if handled {
			// update target to be relative to the NEW slide location
			// We assume the new slide will be in ppt/slides/ just like the old one,
			// so relative paths like "../media/imageX.png" are standard.
			// But we need to construct the relative path from "ppt/slides/slideN.xml" to "ppt/media/imageM.png"

			relPath := common.MakeRelativePath(dstSlidePart, newTarget)
			rels[i].Target = relPath
			changed = true
		}
	}

	if changed {
		rendered := renderRelationshipsXML(rels)
		return []byte(rendered), nil
	}

	return srcSlideRelsBytes, nil
}

func (e *PresentationEditor) copyImageAsset(srcEditor *PresentationEditor, srcPath string) (string, error) {
	data, ok := srcEditor.parts.Get(srcPath)
	if !ok {
		return "", fmt.Errorf("source image part not found: %s", srcPath)
	}

	ext := path.Ext(srcPath)
	if len(ext) > 0 {
		ext = ext[1:] // remove dot
	}

	// RegisterImage handles deduplication via hash
	newPath, err := e.RegisterImage(data, ext)
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func (e *PresentationEditor) copyChartAsset(srcEditor *PresentationEditor, srcPath string) (string, error) {
	data, ok := srcEditor.parts.Get(srcPath)
	if !ok {
		return "", fmt.Errorf("source chart part not found: %s", srcPath)
	}

	// Create new chart part in target
	newChartNum := e.nextChartNum
	e.nextChartNum++
	newChartPath := fmt.Sprintf("ppt/charts/chart%d.xml", newChartNum)

	// We must also copy the chart's relationships (e.g. to Excel data or Colors)
	srcRelsPath := common.SlideRelsPartName(srcPath)
	srcRelsData, hasRels := srcEditor.parts.Get(srcRelsPath)

	if !hasRels {
		e.parts.Set(newChartPath, data)
		return newChartPath, nil
	}

	rels, err := parseRelationshipsXML(srcRelsData)
	if err != nil {
		return "", fmt.Errorf("parse source chart rels: %w", err)
	}

	changed := false
	for i, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}

		srcTargetAbs := common.ResolveRelationshipTarget(srcPath, rel.Target)
		newExcelPath, copyErr := e.copyExcelAsset(srcEditor, srcTargetAbs)
		if copyErr != nil {
			return "", copyErr
		}
		rels[i].Target = common.MakeRelativePath(newChartPath, newExcelPath)
		changed = true
	}

	if changed {
		newRelsData := renderRelationshipsXML(rels)
		e.parts.Set(common.SlideRelsPartName(newChartPath), []byte(newRelsData))
	} else {
		e.parts.Set(common.SlideRelsPartName(newChartPath), srcRelsData)
	}
	e.parts.Set(newChartPath, data)

	// Track embeddings if needed? e.chartEmbeddings
	// Not strictly required for simple copy, but good for bookkeeping.

	return newChartPath, nil
}

func (e *PresentationEditor) copyExcelAsset(srcEditor *PresentationEditor, srcPath string) (string, error) {
	data, ok := srcEditor.parts.Get(srcPath)
	if !ok {
		return "", fmt.Errorf("source excel part not found: %s", srcPath)
	}
	return e.registerExcelEmbedding(data)
}

func (e *PresentationEditor) copyNotesSlideAsset(
	srcEditor *PresentationEditor,
	srcPath, dstSlidePart string,
) (string, error) {
	data, ok := srcEditor.parts.Get(srcPath)
	if !ok {
		return "", fmt.Errorf("source notes part not found: %s", srcPath)
	}

	e.ensureNotesInfrastructure()
	if e.nextNotesNum < 1 {
		e.nextNotesNum = 1
	}

	newNotesPath := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", e.nextNotesNum)
	e.nextNotesNum++
	e.parts.Set(newNotesPath, editorslide.CloneBytes(data))

	srcNotesRelsPath := common.SlideRelsPartName(srcPath)
	if relsData, relsOK := srcEditor.parts.Get(srcNotesRelsPath); relsOK {
		rels, err := parseRelationshipsXML(relsData)
		if err != nil {
			return "", fmt.Errorf("parse source notes rels: %w", err)
		}
		for i, rel := range rels {
			switch rel.Type {
			case common.RelTypeSlide:
				rels[i].Target = common.MakeRelativePath(newNotesPath, dstSlidePart)
			case common.RelTypeNotesMaster:
				rels[i].Target = "../notesMasters/notesMaster1.xml"
			}
		}
		rendered := renderRelationshipsXML(rels)
		e.parts.Set(common.SlideRelsPartName(newNotesPath), []byte(rendered))
	}

	e.notesInventory[dstSlidePart] = newNotesPath
	return newNotesPath, nil
}

func (e *PresentationEditor) ListSlideCharts(slideIndex int) ([]common.SlideChartRef, error) {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return nil, fmt.Errorf("slide index %d out of range", slideIndex)
	}

	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return nil, fmt.Errorf("slide part %s not found", slideRef.Part)
	}
	relsPath := common.RelsPathFor(slideRef.Part)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return nil, fmt.Errorf("rels part %s not found", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil, fmt.Errorf("parse slide rels: %w", err)
	}

	relToChart := make(map[string]string)
	for _, rel := range rels {
		if rel.Type != common.RelTypeChart {
			continue
		}
		relToChart[rel.ID] = resolveRelationshipTarget(slideRef.Part, rel.Target)
	}

	rIDs := extractChartRelIDs(slideXML)
	out := make([]common.SlideChartRef, 0, len(rIDs))
	for i, relID := range rIDs {
		chartPart, found := relToChart[relID]
		if !found {
			continue
		}
		out = append(out, common.SlideChartRef{
			Index:     i,
			RelID:     relID,
			ChartPart: chartPart,
		})
	}
	return out, nil
}

func (e *PresentationEditor) UpdateChartData(
	slideIndex int,
	selector common.ChartSelector,
	req common.ChartDataUpdate,
) error {
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return err
	}
	chartRef, err := editormodchart.ResolveChartSelector(refs, selector, slideIndex)
	if err != nil {
		return err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}

	kind := editormodchart.DetectChartKind(chartXML)
	if validateErr := editormodchart.ValidateChartUpdatePayload(kind, req); validateErr != nil {
		return validateErr
	}

	workbook, err := editormodchart.GenerateExcelForChartUpdate(kind, req)
	if err != nil {
		return fmt.Errorf("generate excel: %w", err)
	}

	excelPartPath, err := e.resolveChartEmbeddingPath(chartRef.ChartPart)
	if err != nil {
		return err
	}
	e.parts.Set(excelPartPath, workbook)
	e.chartEmbeddings[chartRef.ChartPart] = excelPartPath

	patchedChartXML, err := editormodchart.PatchChartDataCache(chartXML, kind, req)
	if err != nil {
		return err
	}
	e.parts.Set(chartRef.ChartPart, patchedChartXML)
	return nil
}

func (e *PresentationEditor) resolveChartEmbeddingPath(chartPart string) (string, error) {
	if excelPath, ok := e.chartEmbeddings[chartPart]; ok {
		return excelPath, nil
	}

	relsPath := common.RelsPathFor(chartPart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return "", fmt.Errorf("chart rels part not found: %s", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return "", fmt.Errorf("parse chart rels: %w", err)
	}
	for _, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}
		excelPath := common.CanonicalPartPath(path.Join(path.Dir(chartPart), rel.Target))
		if !e.parts.Has(excelPath) {
			return "", fmt.Errorf("chart embedding part missing: %s", excelPath)
		}
		e.chartEmbeddings[chartPart] = excelPath
		return excelPath, nil
	}
	return "", errors.New("associated excel part not found for chart " + chartPart)
}
