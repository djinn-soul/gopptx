package export

import (
	"archive/zip"
	"path"
	"regexp"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

// A diagram's data part points at the drawing PowerPoint cached for it through
// a dsp:dataModelExt extension, which names a relationship of the data part.
// gopptx's own writer deliberately drops that extension so PowerPoint lays the
// diagram out from the data model rather than trusting a template's cache (see
// dropSmartArtDrawingCacheLink), which leaves the drawing part in the package
// with nothing pointing at it. So the lookup falls back to the part naming
// convention — dataN.xml is cached by drawingN.xml — before giving up.
var (
	reSmartArtDataModelExt = regexp.MustCompile(`<dsp:dataModelExt\b[^>]*\brelId="([^"]+)"`)
	reSmartArtDataPartNum  = regexp.MustCompile(`^data(\d+)\.xml$`)
)

// readSmartArtDrawingCache loads the cached layout of one diagram, or nil when
// the package has none to offer.
func readSmartArtDrawingCache(fileMap map[string]*zip.File, dataPath string, dataXML []byte) []smartart.DrawingShape {
	drawingPath := resolveSmartArtDrawingPath(fileMap, dataPath, dataXML)
	if drawingPath == "" {
		return nil
	}
	drawingXML := readZipBytes(fileMap, drawingPath)
	if drawingXML == nil {
		return nil
	}
	return smartart.ParseDrawingShapes(drawingXML)
}

func resolveSmartArtDrawingPath(fileMap map[string]*zip.File, dataPath string, dataXML []byte) string {
	rels := readZipRelationships(fileMap, slideRelsPath(dataPath))
	if relID := smartArtDataModelExtRelID(dataXML); relID != "" {
		if target := rels[relID]; target != "" {
			if resolved := resolveRelPath(dataPath, target); resolved != "" {
				return resolved
			}
		}
	}
	// No extension, or one that points nowhere: take a related part that is
	// named like a drawing, then fall back to the naming convention.
	for _, target := range rels {
		if !strings.HasPrefix(path.Base(target), "drawing") {
			continue
		}
		if resolved := resolveRelPath(dataPath, target); resolved != "" {
			return resolved
		}
	}
	return conventionalSmartArtDrawingPath(dataPath)
}

func smartArtDataModelExtRelID(dataXML []byte) string {
	match := reSmartArtDataModelExt.FindSubmatch(dataXML)
	if len(match) < 2 {
		return ""
	}
	return string(match[1])
}

// conventionalSmartArtDrawingPath is ppt/diagrams/drawingN.xml for a data part
// named ppt/diagrams/dataN.xml.
func conventionalSmartArtDrawingPath(dataPath string) string {
	match := reSmartArtDataPartNum.FindStringSubmatch(path.Base(dataPath))
	if len(match) < 2 {
		return ""
	}
	return path.Join(path.Dir(dataPath), "drawing"+match[1]+".xml")
}
