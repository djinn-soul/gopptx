package slide

import (
	"encoding/xml"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const commentsPartType = "application/vnd.openxmlformats-officedocument.presentationml.comments+xml"

const (
	contentTypesBaseGrowCap = 80
	contentTypesItemGrowCap = 100
)

// ContentTypesBase is an opaque pre-parsed [Content_Types].xml representation.
// Callers obtain it via ParseContentTypesBase and pass it to RewriteContentTypesFromBase
// to skip xml.Unmarshal on repeated saves when the bytes have not changed.
type ContentTypesBase interface {
	// unexported method keeps the interface opaque; only this package can implement it.
	contentTypesBase() contentTypesDocument
}

type contentTypesBaseImpl struct {
	doc contentTypesDocument
}

func (c contentTypesBaseImpl) contentTypesBase() contentTypesDocument { return c.doc }

// ParseContentTypesBase parses [Content_Types].xml and returns an opaque value
// suitable for caching and passing to RewriteContentTypesFromBase.
func ParseContentTypesBase(current []byte) (ContentTypesBase, error) {
	doc, err := parseContentTypesDocument(current)
	if err != nil {
		return nil, err
	}
	return contentTypesBaseImpl{doc}, nil
}

// RewriteContentTypesFromBase is identical to RewriteContentTypes but accepts a
// pre-parsed base (from ParseContentTypesBase) to avoid re-parsing on each call.
func RewriteContentTypesFromBase(
	base ContentTypesBase,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	doc := base.contentTypesBase()
	// Copy Defaults slice before passing to ensureContentTypeDefaults, which may
	// append to it.  This prevents the cached base from growing across saves.
	if len(mediaPaths) > 0 || hasVBA {
		cp := make([]contentTypeDefault, len(doc.Defaults))
		copy(cp, doc.Defaults)
		doc.Defaults = cp
	}
	return rewriteContentTypesFromDoc(doc, slides, mediaPaths, hasSections, chartPaths,
		notesPaths, themePaths, layoutPaths, masterPaths, hasNotesMaster, hasCommentAuthors,
		commentPaths, hasVBA, hasHandoutMaster, customXMLPropsPaths)
}

func RewriteContentTypes(
	current []byte,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	doc, err := parseContentTypesDocument(current)
	if err != nil {
		return "", err
	}
	return rewriteContentTypesFromDoc(doc, slides, mediaPaths, hasSections, chartPaths,
		notesPaths, themePaths, layoutPaths, masterPaths, hasNotesMaster, hasCommentAuthors,
		commentPaths, hasVBA, hasHandoutMaster, customXMLPropsPaths)
}

func rewriteContentTypesFromDoc(
	doc contentTypesDocument,
	slides []common.EditorSlideRef,
	mediaPaths []string,
	hasSections bool,
	chartPaths []string,
	notesPaths []string,
	themePaths []string,
	layoutPaths []string,
	masterPaths []string,
	hasNotesMaster bool,
	hasCommentAuthors bool,
	commentPaths []string,
	hasVBA bool,
	hasHandoutMaster bool,
	customXMLPropsPaths []string,
) (string, error) {
	ensureContentTypeDefaults(&doc, mediaPaths, hasVBA)

	extraOverrides := len(slides) +
		len(chartPaths) +
		len(notesPaths) +
		len(themePaths) +
		len(layoutPaths) +
		len(masterPaths) +
		len(commentPaths) +
		len(customXMLPropsPaths)
	if hasSections {
		extraOverrides++
	}
	if hasNotesMaster {
		extraOverrides++
	}
	if hasHandoutMaster {
		extraOverrides++
	}
	if hasCommentAuthors {
		extraOverrides++
	}
	if hasVBA {
		extraOverrides++
	}

	overrides := filterDynamicOverrides(doc.Overrides, extraOverrides)
	overrides = appendSlideOverrides(overrides, slides)
	overrides = appendOptionalContentTypeOverride(overrides, hasSections, "/ppt/sectionList.xml",
		"application/vnd.microsoft.powerpoint.sectionList+xml")
	overrides = appendPathOverrides(overrides, chartPaths,
		"application/vnd.openxmlformats-officedocument.drawingml.chart+xml")
	overrides = appendPathOverrides(overrides, notesPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml")
	overrides = appendPathOverrides(overrides, themePaths,
		"application/vnd.openxmlformats-officedocument.theme+xml")
	overrides = appendPathOverrides(overrides, layoutPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml")
	overrides = appendPathOverrides(overrides, masterPaths,
		"application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasNotesMaster, "/ppt/notesMasters/notesMaster1.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.notesMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasHandoutMaster, "/ppt/handoutMasters/handoutMaster1.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.handoutMaster+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasCommentAuthors, "/ppt/commentAuthors.xml",
		"application/vnd.openxmlformats-officedocument.presentationml.commentAuthors+xml")
	overrides = appendOptionalContentTypeOverride(overrides, hasVBA, "/ppt/vbaProject.bin",
		"application/vnd.ms-office.vbaProject")
	overrides = appendPathOverrides(overrides, commentPaths, commentsPartType)
	overrides = appendPathOverrides(overrides, customXMLPropsPaths,
		"application/vnd.openxmlformats-officedocument.customXmlProperties+xml")
	overrides = dedupeContentTypeOverrides(overrides)

	sort.Slice(overrides, func(i, j int) bool { return overrides[i].PartName < overrides[j].PartName })
	doc.Overrides = overrides
	return renderContentTypesDocument(doc)
}

func parseContentTypesDocument(current []byte) (contentTypesDocument, error) {
	if len(current) == 0 {
		return contentTypesDocument{}, errors.New("missing content types content")
	}

	var doc contentTypesDocument
	if err := xml.Unmarshal(current, &doc); err != nil {
		return contentTypesDocument{}, fmt.Errorf("parse content types: %w", err)
	}
	if strings.TrimSpace(doc.XMLNS) == "" {
		doc.XMLNS = common.ContentTypesXMLNS
	}
	return doc, nil
}

func ensureContentTypeDefaults(doc *contentTypesDocument, mediaPaths []string, hasVBA bool) {
	exts := make(map[string]struct{}, len(doc.Defaults))
	for _, d := range doc.Defaults {
		exts[strings.ToLower(d.Extension)] = struct{}{}
	}
	for _, mediaPath := range mediaPaths {
		ext := strings.TrimPrefix(strings.ToLower(path.Ext(strings.TrimSpace(mediaPath))), ".")
		if _, ok := exts[ext]; ok {
			continue
		}
		contentType := contentTypeForExtension(ext)
		if contentType == "" {
			continue
		}
		doc.Defaults = append(doc.Defaults, contentTypeDefault{
			Extension:   ext,
			ContentType: contentType,
		})
		exts[ext] = struct{}{}
	}

	if hasVBA {
		if _, ok := exts["bin"]; !ok {
			doc.Defaults = append(doc.Defaults, contentTypeDefault{
				Extension:   "bin",
				ContentType: "application/vnd.ms-office.vbaProject",
			})
			exts["bin"] = struct{}{}
		}
	}
}

type contentTypesDocument struct {
	XMLName   xml.Name              `xml:"Types"`
	XMLNS     string                `xml:"xmlns,attr,omitempty"`
	Defaults  []contentTypeDefault  `xml:"Default"`
	Overrides []contentTypeOverride `xml:"Override"`
}

type contentTypeDefault struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

type contentTypeOverride struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}
