package editor

import (
	"fmt"
	"path"

	"github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// deepCloneSlideAssets walks through the relationships of a source slide and copies
// all referenced assets (images, charts, etc.) to the target editor.
// It returns a modified relationships XML byte slice where targets are remapped to the new locations.
func (e *PresentationEditor) deepCloneSlideAssets(srcEditor *PresentationEditor, srcSlidePart string, srcSlideRelsBytes []byte, dstSlidePart string) ([]byte, error) {
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
			// Notes slides are handled by specific logic in ops.go usually, or we can generalize here.
			// The original implementation in ops.go handled this.
			// For now, let's leave it as "pass through" if we don't move that logic here yet,
			// OR we can implement `copyNotesSlideAsset` if we want to be fully generic.
			// The task focuses on Media and Charts.
			// Let's implement basics first.
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
		rendered, err := renderRelationshipsXML(rels)
		if err != nil {
			return nil, err
		}
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

	if hasRels {
		rels, err := parseRelationshipsXML(srcRelsData)
		if err != nil {
			return "", fmt.Errorf("parse source chart rels: %w", err)
		}

		changed := false
		for i, rel := range rels {
			srcTargetAbs := common.ResolveRelationshipTarget(srcPath, rel.Target)

			// Check for Excel embedding
			if rel.Type == "http://schemas.openxmlformats.org/officeDocument/2006/relationships/package" {
				newExcelPath, err := e.copyExcelAsset(srcEditor, srcTargetAbs)
				if err != nil {
					return "", err
				}

				// Relink
				relPath := common.MakeRelativePath(newChartPath, newExcelPath)
				rels[i].Target = relPath

				// Also need to update the rId validation in the chart XML if strictly required,
				// but usually rId stays same if we preserve order.
				// However, if we filter rels, rIDs might change?
				// Actually, we are copying the rels file, so we keep IDs stable.
				changed = true
			}
			// TODO: Handle chart colors/styles if external?
		}

		if changed {
			newRelsData, err := renderRelationshipsXML(rels)
			if err != nil {
				return "", err
			}
			e.parts.Set(common.SlideRelsPartName(newChartPath), []byte(newRelsData))
		} else {
			e.parts.Set(common.SlideRelsPartName(newChartPath), srcRelsData)
		}
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

	newExcelNum := e.nextExcelNum
	e.nextExcelNum++
	// Keep extension
	ext := path.Ext(srcPath)
	if ext == "" {
		ext = ".xlsx"
	}
	newPath := fmt.Sprintf("ppt/embeddings/Microsoft_Excel_Worksheet%d%s", newExcelNum, ext)

	e.parts.Set(newPath, data)
	return newPath, nil
}
