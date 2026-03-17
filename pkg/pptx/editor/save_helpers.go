package editor

import (
	"fmt"
	"sync"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

type manifestParts struct {
	presentationXML     []byte
	presentationRelsXML []byte
	corePropsXML        []byte
	contentTypesXML     []byte
}

func (e *PresentationEditor) buildManifestPartsParallel(
	hasSections bool,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (manifestParts, error) {
	var (
		wg      sync.WaitGroup
		errOnce sync.Once
		outErr  error

		parts manifestParts
	)

	setErr := func(err error) {
		if err == nil {
			return
		}
		errOnce.Do(func() {
			outErr = err
		})
	}

	wg.Add(manifestBuildWorkers)

	go func() {
		defer wg.Done()
		presentationXML, err := e.renderPresentationXMLWithSections()
		if err != nil {
			setErr(err)
			return
		}
		parts.presentationXML = []byte(presentationXML)
	}()

	go func() {
		defer wg.Done()
		presentationRelsXML, err := editorslide.RenderPresentationRelsXML(e.nonSlideRels, e.slides, hasSections, hasVBA)
		if err != nil {
			setErr(err)
			return
		}
		parts.presentationRelsXML = []byte(presentationRelsXML)
	}()

	go func() {
		defer wg.Done()
		corePropsXML, err := renderCoreProperties(e.metadata.CoreProperties)
		if err != nil {
			setErr(fmt.Errorf("render core properties: %w", err))
			return
		}
		parts.corePropsXML = corePropsXML
	}()

	go func() {
		defer wg.Done()
		contentTypesXML, err := e.renderContentTypesPart(
			hasSections,
			hasNotesMaster,
			hasCommentAuthors,
			hasVBA,
			hasHandoutMaster,
			customXMLPropsPaths,
		)
		if err != nil {
			setErr(err)
			return
		}
		parts.contentTypesXML = contentTypesXML
	}()

	wg.Wait()
	if outErr != nil {
		return manifestParts{}, outErr
	}
	return parts, nil
}

func (e *PresentationEditor) renderContentTypesPart(
	hasSections bool,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) ([]byte, error) {
	mediaPaths := editorslide.MapValues(e.mediaInventory)
	filteredChartPaths := editorslide.FilterXMLPartPaths(e.parts.KeysWithPrefix("ppt/charts/chart"))
	notesPaths := editorslide.MapValues(e.notesInventory)

	themePaths := e.parts.KeysWithPrefix("ppt/theme/theme")
	layoutPaths := e.parts.KeysWithPrefix("ppt/slideLayouts/slideLayout")
	masterPaths := e.parts.KeysWithPrefix("ppt/slideMasters/slideMaster")

	filteredCommentPaths := editorslide.FilterXMLPartPaths(e.parts.KeysWithPrefix("ppt/comments/comment"))

	contentTypesData, _ := e.parts.Get(common.ContentTypesPath)
	contentTypesXML, err := editorslide.RewriteContentTypes(
		contentTypesData,
		e.slides,
		mediaPaths,
		hasSections,
		filteredChartPaths,
		notesPaths,
		themePaths,
		layoutPaths,
		masterPaths,
		hasNotesMaster,
		hasCommentAuthors,
		filteredCommentPaths,
		hasVBA,
		hasHandoutMaster,
		customXMLPropsPaths,
	)
	if err != nil {
		return nil, err
	}
	return []byte(contentTypesXML), nil
}

func (e *PresentationEditor) writeOptionalPresentationParts(
	out map[string][]byte,
	hasSections bool,
	hasNotesMaster bool,
	hasHandoutMaster bool,
	hasVBA bool,
	vbaProject *vba.VBAProject,
) {
	if hasSections {
		out["ppt/sectionList.xml"] = []byte(editorslide.BuildSectionListXML(e.sections))
	}

	if hasNotesMaster {
		// Ensure notes master rels are also persisted if they were injected
		if masterRels, ok := e.parts.Get("ppt/notesMasters/_rels/notesMaster1.xml.rels"); ok {
			out["ppt/notesMasters/_rels/notesMaster1.xml.rels"] = masterRels
		}
	}
	if hasHandoutMaster {
		if masterXML, ok := e.parts.Get("ppt/handoutMasters/handoutMaster1.xml"); ok {
			out["ppt/handoutMasters/handoutMaster1.xml"] = masterXML
		}
		if masterRels, ok := e.parts.Get("ppt/handoutMasters/_rels/handoutMaster1.xml.rels"); ok {
			out["ppt/handoutMasters/_rels/handoutMaster1.xml.rels"] = masterRels
		}
	}

	if hasVBA {
		out["ppt/vbaProject.bin"] = vbaProject.Data
	}
}

func (e *PresentationEditor) renderPresentationXMLWithSections() (string, error) {
	presentationXML, err := editorslide.RewritePresentationSlideList([]byte(e.presentationXML), e.slides)
	if err != nil {
		return "", err
	}

	hasNotesMaster := e.parts.Has("ppt/notesMasters/notesMaster1.xml")
	notesMasterRelID, err := editorslide.ResolveNotesMasterRelID(
		e.nonSlideRels,
		hasNotesMaster,
		common.RelTypeNotesMaster,
	)
	if err != nil {
		return "", err
	}

	presentationXML, err = editorslide.RewritePresentationNotesMasterList(
		[]byte(presentationXML),
		notesMasterRelID,
		hasNotesMaster,
	)
	if err != nil {
		return "", err
	}
	presentationXML, err = rewritePresentationModifyVerifier(
		[]byte(presentationXML),
		e.metadata.Protection.ModifyPassword,
	)
	if err != nil {
		return "", err
	}

	if len(e.sections) == 0 {
		return presentationXML, nil
	}

	presentationXML, err = editorslide.RewritePresentationSections([]byte(presentationXML), e.sections)
	if err != nil {
		return "", fmt.Errorf("rewrite sections: %w", err)
	}

	presentationXML, err = editorslide.RewritePresentationEmbeddedFonts([]byte(presentationXML), e.embeddedFontLst)
	if err != nil {
		return "", fmt.Errorf("rewrite embedded fonts: %w", err)
	}

	return presentationXML, nil
}
