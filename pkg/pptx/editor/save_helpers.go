package editor

import (
	"fmt"
	"unsafe"

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

func (e *PresentationEditor) buildManifestParts(
	hasSections bool,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (manifestParts, error) {
	// Run the four manifest-building tasks sequentially.
	// All tasks are CPU-bound string/map operations; zip writing dominates total
	// latency, so goroutine+sync overhead (~83 allocs/save) outweighs any parallel speedup.

	presentationXML, err := e.renderPresentationXMLWithSections()
	if err != nil {
		return manifestParts{}, err
	}

	presentationRelsXML, err := editorslide.RenderPresentationRelsXML(e.nonSlideRels, e.slides, hasSections, hasVBA)
	if err != nil {
		return manifestParts{}, err
	}

	corePropsXML, err := renderCoreProperties(e.metadata.CoreProperties)
	if err != nil {
		return manifestParts{}, fmt.Errorf("render core properties: %w", err)
	}

	contentTypesXML, err := e.renderContentTypesPart(
		hasSections,
		hasNotesMaster,
		hasCommentAuthors,
		hasVBA,
		hasHandoutMaster,
		customXMLPropsPaths,
	)
	if err != nil {
		return manifestParts{}, err
	}

	return manifestParts{
		presentationXML:     []byte(presentationXML),
		presentationRelsXML: []byte(presentationRelsXML),
		corePropsXML:        corePropsXML,
		contentTypesXML:     contentTypesXML,
	}, nil
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

	// Use a cached parse of [Content_Types].xml to avoid xml.Unmarshal on every save.
	// The cache is keyed by the backing-array pointer of the content-types bytes.
	// PartStore.Set() always creates a new slice, so a changed pointer means a cache miss.
	var base editorslide.ContentTypesBase
	if len(contentTypesData) > 0 {
		ptr := uintptr(unsafe.Pointer(&contentTypesData[0])) //nolint:gosec // intentional: staleness token, not dereferenced
		if e.ctBasePtr == ptr && e.ctBase != nil {
			base = e.ctBase
		} else {
			parsed, err := editorslide.ParseContentTypesBase(contentTypesData)
			if err != nil {
				return nil, fmt.Errorf("parse content types: %w", err)
			}
			e.ctBase = parsed
			e.ctBasePtr = ptr
			base = parsed
		}
	}

	var contentTypesXML string
	var err error
	if base != nil {
		contentTypesXML, err = editorslide.RewriteContentTypesFromBase(
			base,
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
	} else {
		contentTypesXML, err = editorslide.RewriteContentTypes(
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
	}
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
