package editor

import (
	"fmt"
	"path"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
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
// layoutImports remembers which source layout parts have already been imported
// and where they landed, so slides sharing a layout share one imported family.
func (e *PresentationEditor) deepCloneSlideAssets(
	srcEditor *PresentationEditor,
	srcSlidePart string,
	srcSlideRelsBytes []byte,
	dstSlidePart string,
	layoutImports map[string]string,
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
		case common.RelTypeSlideLayout:
			// Left alone, a copied slide binds to whatever local part happens to
			// carry the same name, or to none at all -- so it either looks wrong
			// or breaks the package. The import itself ran before the merge
			// started; this only retargets to what it produced.
			local, imported := layoutImports[common.CanonicalPartPath(srcTargetAbs)]
			if !imported {
				break
			}
			newTarget = local
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

// importSlideLayouts brings the layout family of every source slide into this
// presentation and returns the source-to-local layout part mapping.
//
// This runs before the slides are merged, not while they are being merged: the
// import allocates presentation relationship ids of its own, and doing that
// halfway through the merge made the master and the new slide claim the same
// rId, which cost the slide its relationship and left a package PowerPoint
// refuses to open.
//
// A whole family is imported at once, so a second slide bound to a sibling
// layout of the same master reuses it instead of duplicating the master.
func (e *PresentationEditor) importSlideLayouts(
	srcEditor *PresentationEditor,
	sourceSlides []common.EditorSlideRef,
) (map[string]string, error) {
	layoutImports := map[string]string{}
	if srcEditor == e {
		return layoutImports, nil
	}
	for _, slide := range sourceSlides {
		layoutPart, err := srcEditor.slideLayoutPartFor(slide.Part)
		if err != nil {
			return nil, err
		}
		if layoutPart == "" {
			continue
		}
		if _, done := layoutImports[layoutPart]; done {
			continue
		}
		result, importErr := e.ImportLayoutMasterFamily(srcEditor, layoutPart)
		if importErr != nil {
			return nil, fmt.Errorf("import layout %s: %w", layoutPart, importErr)
		}
		for sourceLayout, localLayout := range result.LayoutMap {
			layoutImports[common.CanonicalPartPath(sourceLayout)] = localLayout
		}
		if _, ok := layoutImports[layoutPart]; !ok {
			return nil, fmt.Errorf("imported layout family does not contain %s", layoutPart)
		}
	}
	return layoutImports, nil
}

// slideLayoutPartFor resolves the layout a slide is bound to, or "" when it has
// no layout relationship.
func (e *PresentationEditor) slideLayoutPartFor(slidePart string) (string, error) {
	relsPath := common.RelsPathFor(slidePart)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return "", nil
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return "", fmt.Errorf("parse slide rels %s: %w", relsPath, err)
	}
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideLayout {
			return common.CanonicalPartPath(
				common.ResolveRelationshipTarget(slidePart, rel.Target),
			), nil
		}
	}
	return "", nil
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
