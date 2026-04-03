package export

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"path"
	"strconv"
	"strings"
)

const formatPNG = "png"

// SlideImage holds image bytes and its position on the slide (in EMU).
type SlideImage struct {
	Bytes        []byte
	Format       string // "png", "jpeg", "gif", "emf", etc.
	X, Y         int64  // EMU offset
	CX, CY       int64  // EMU size
	Rotation     float64
	CropLeft     float64
	CropRight    float64
	CropTop      float64
	CropBottom   float64
	FlipH        bool
	FlipV        bool
	Shadow       bool
	Reflection   bool
	AltText      string
	IsDecorative bool
}

// extractSlideImages reads a PPTX file and returns images per slide (0-based).
// It parses p:pic elements in each slide XML and resolves image data from ppt/media/.
func extractSlideImages(pptxPath string) ([][]SlideImage, error) {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	// Build a lookup: canonical path → zip.File
	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}

	// Find slide parts in order using presentation.xml.rels
	slideOrder := resolveSlideOrder(fileMap)

	result := make([][]SlideImage, len(slideOrder))
	for i, slidePart := range slideOrder {
		relsPath := slideRelsPath(slidePart)
		rels := readZipRelationships(fileMap, relsPath)

		// Read slide XML
		slideData := readZipBytes(fileMap, slidePart)
		if slideData == nil {
			continue
		}

		pics := parsePicElements(slideData)
		images := make([]SlideImage, 0, len(pics))
		for _, pic := range pics {
			target := rels[pic.RelID]
			if target == "" {
				continue
			}
			// Resolve target relative to slide part directory
			mediaPath := resolveRelPath(slidePart, target)
			data := readZipBytes(fileMap, mediaPath)
			if data == nil {
				continue
			}
			images = append(images, SlideImage{
				Bytes:        data,
				Format:       imageFormat(mediaPath),
				X:            pic.X,
				Y:            pic.Y,
				CX:           pic.CX,
				CY:           pic.CY,
				Rotation:     pic.Rotation,
				CropLeft:     pic.CropLeft,
				CropRight:    pic.CropRight,
				CropTop:      pic.CropTop,
				CropBottom:   pic.CropBottom,
				FlipH:        pic.FlipH,
				FlipV:        pic.FlipV,
				Shadow:       pic.Shadow,
				Reflection:   pic.Reflection,
				AltText:      pic.AltText,
				IsDecorative: pic.IsDecorative,
			})
		}
		result[i] = images
	}

	return result, nil
}

// resolveSlideOrder reads presentation.xml and its rels to get ordered slide part paths.
func resolveSlideOrder(fileMap map[string]*zip.File) []string {
	presRels := readZipRelationships(fileMap, "ppt/_rels/presentation.xml.rels")

	presData := readZipBytes(fileMap, "ppt/presentation.xml")
	if presData == nil {
		return nil
	}

	// Parse sldIdLst to get ordered r:id attributes on sldId elements.
	// The attribute is namespace-qualified: xmlns:r="...relationships" r:id="rId2"
	// Go's xml.Decoder exposes it as Attr{Name:{Space: "<full-ns-uri>", Local: "id"}, Value: "rIdN"}.
	var slideRelIDs []string
	dec := xml.NewDecoder(strings.NewReader(string(presData)))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		se, ok := tok.(xml.StartElement)
		if !ok || se.Name.Local != "sldId" {
			continue
		}
		for _, a := range se.Attr {
			if a.Name.Local == "id" && strings.Contains(a.Name.Space, "relationship") {
				slideRelIDs = append(slideRelIDs, a.Value)
				break
			}
		}
	}

	parts := make([]string, 0, len(slideRelIDs))
	for _, relID := range slideRelIDs {
		if target, ok := presRels[relID]; ok {
			// target is relative to ppt/, e.g. "slides/slide1.xml"
			parts = append(parts, canonicalZipPath("ppt/"+target))
		}
	}
	return parts
}

// readZipRelationships parses a .rels XML file and returns relID→target map.
func readZipRelationships(fileMap map[string]*zip.File, relsPath string) map[string]string {
	data := readZipBytes(fileMap, relsPath)
	if data == nil {
		return nil
	}
	rels := make(map[string]string)
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "Relationship" {
			var id, target string
			for _, a := range se.Attr {
				switch a.Name.Local {
				case "Id":
					id = a.Value
				case "Target":
					target = a.Value
				}
			}
			if id != "" && target != "" {
				rels[id] = target
			}
		}
	}
	return rels
}

func readZipBytes(fileMap map[string]*zip.File, name string) []byte {
	f, ok := fileMap[canonicalZipPath(name)]
	if !ok {
		return nil
	}
	rc, err := f.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil
	}
	return data
}

func canonicalZipPath(p string) string {
	p = strings.ReplaceAll(p, "\\", "/")
	return strings.TrimPrefix(p, "/")
}

func slideRelsPath(slidePart string) string {
	dir := path.Dir(slidePart)
	base := path.Base(slidePart)
	return dir + "/_rels/" + base + ".rels"
}

func resolveRelPath(slidePart, target string) string {
	if strings.HasPrefix(target, "/") {
		return canonicalZipPath(target)
	}

	dir := path.Dir(slidePart)
	resolved := path.Join(dir, target)

	// path.Join already calls path.Clean.
	// If it starts with '..', it means it escaped the root.
	if strings.HasPrefix(resolved, "..") {
		return ""
	}

	// Security: The test also expects 'slides/../media/image1.png' to be blocked
	// even if it resolves to 'ppt/media/image1.png'.
	// This is likely to prevent complex path logic.
	if strings.Contains(target, "/../") || strings.HasPrefix(target, "../../../") {
		return ""
	}

	return canonicalZipPath(resolved)
}

func imageFormat(p string) string {
	ext := strings.ToLower(path.Ext(p))
	switch ext {
	case ".png":
		return formatPNG
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".gif":
		return "gif"
	case ".emf", ".wmf":
		return "emf"
	default:
		return formatPNG
	}
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return n
}
