package export

import (
	"archive/zip"
	"regexp"
	"strings"

	editorshape "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/shape"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

type parsedSmartArt struct {
	ShapeID      int
	LayoutURI    string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	AltText      string
	IsDecorative bool
	QuickStyle   string
	ColorStyle   string
	Nodes        []smartart.Node
}

type smartArtFrameRef struct {
	ShapeID      int
	DataRelID    string
	LayoutRelID  string
	StyleRelID   string
	ColorRelID   string
	X            int64
	Y            int64
	CX           int64
	CY           int64
	AltText      string
	IsDecorative bool
}

var (
	reSmartArtGraphicFrame = regexp.MustCompile(`(?s)<p:graphicFrame\b.*?</p:graphicFrame>`)
	reSmartArtRelTag       = regexp.MustCompile(`<dgm:relIds\b[^>]*>`)
	reSmartArtRelID        = regexp.MustCompile(`r:(dm|lo|qs|cs)=["']([^"']+)["']`)
)

func extractSlideSmartArt(pptxPath string) ([][]parsedSmartArt, error) {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}
	slideOrder := resolveSlideOrder(fileMap)
	out := make([][]parsedSmartArt, len(slideOrder))
	for idx, slidePart := range slideOrder {
		slideXML := readZipBytes(fileMap, slidePart)
		if slideXML == nil {
			continue
		}
		frames := parseSmartArtFrames(slideXML)
		if len(frames) == 0 {
			continue
		}
		rels := readZipRelationships(fileMap, slideRelsPath(slidePart))
		row := make([]parsedSmartArt, 0, len(frames))
		for _, frame := range frames {
			dataTarget := rels[frame.DataRelID]
			layoutTarget := rels[frame.LayoutRelID]
			if dataTarget == "" || layoutTarget == "" {
				continue
			}
			styleTarget := rels[frame.StyleRelID]
			colorTarget := rels[frame.ColorRelID]
			dataPath := resolveRelPath(slidePart, dataTarget)
			layoutPath := resolveRelPath(slidePart, layoutTarget)
			if dataPath == "" || layoutPath == "" {
				continue
			}
			stylePath := resolveRelPath(slidePart, styleTarget)
			colorPath := resolveRelPath(slidePart, colorTarget)
			dataXML := readZipBytes(fileMap, dataPath)
			layoutXML := readZipBytes(fileMap, layoutPath)
			if dataXML == nil || layoutXML == nil {
				continue
			}
			layoutURI := smartart.ExtractLayoutURI(string(layoutXML))
			nodes, err := smartart.ParseDataModelNodes(dataXML)
			if err != nil {
				continue
			}
			quickStyle := extractSmartArtStyleID(fileMap, stylePath)
			colorStyle := extractSmartArtStyleID(fileMap, colorPath)
			if layoutURI == "" && len(nodes) == 0 {
				continue
			}
			row = append(row, parsedSmartArt{
				ShapeID:      frame.ShapeID,
				LayoutURI:    layoutURI,
				X:            frame.X,
				Y:            frame.Y,
				CX:           frame.CX,
				CY:           frame.CY,
				AltText:      frame.AltText,
				IsDecorative: frame.IsDecorative,
				QuickStyle:   quickStyle,
				ColorStyle:   colorStyle,
				Nodes:        nodes,
			})
		}
		out[idx] = row
	}
	return out, nil
}

func parseSmartArtFrames(slideXML []byte) []smartArtFrameRef {
	frames := reSmartArtGraphicFrame.FindAllString(string(slideXML), -1)
	out := make([]smartArtFrameRef, 0, len(frames))
	for _, frame := range frames {
		if !strings.Contains(frame, "<dgm:relIds") {
			continue
		}
		props, err := editorshape.ParseShapeProperties([]byte(frame))
		if err != nil {
			continue
		}
		meta, err := editorshape.ParseShapeReaderMetadata([]byte(frame))
		if err != nil {
			continue
		}
		relTag := reSmartArtRelTag.FindString(frame)
		if relTag == "" {
			continue
		}
		ref := smartArtFrameRef{
			ShapeID:      props.ID,
			X:            int64(props.X),
			Y:            int64(props.Y),
			CX:           int64(props.W),
			CY:           int64(props.H),
			AltText:      meta.AltText,
			IsDecorative: meta.IsDecorative,
		}
		for _, match := range reSmartArtRelID.FindAllStringSubmatch(relTag, -1) {
			switch match[1] {
			case "dm":
				ref.DataRelID = match[2]
			case "lo":
				ref.LayoutRelID = match[2]
			case "qs":
				ref.StyleRelID = match[2]
			case "cs":
				ref.ColorRelID = match[2]
			}
		}
		if ref.DataRelID == "" || ref.LayoutRelID == "" {
			continue
		}
		out = append(out, ref)
	}
	return out
}

func applyParsedSmartArt(slide *elements.SlideContent, diagrams []parsedSmartArt) {
	for _, parsed := range diagrams {
		diagram := smartart.NewSmartArt(smartart.CustomLayout(parsed.LayoutURI)).
			Position(styling.Emu(parsed.X), styling.Emu(parsed.Y)).
			Size(styling.Emu(parsed.CX), styling.Emu(parsed.CY))
		if parsed.AltText != "" {
			diagram = diagram.WithAltText(parsed.AltText)
		}
		if parsed.IsDecorative {
			diagram = diagram.WithDecorative(true)
		}
		if parsed.QuickStyle != "" {
			diagram = diagram.WithQuickStyle(parsed.QuickStyle)
		}
		if parsed.ColorStyle != "" {
			diagram = diagram.WithColorStyle(parsed.ColorStyle)
		}
		for _, node := range parsed.Nodes {
			diagram = diagram.AddNode(node)
		}
		slide.SmartArtDiagrams = append(slide.SmartArtDiagrams, diagram)
	}
}

func extractSmartArtStyleID(fileMap map[string]*zip.File, partPath string) string {
	if partPath == "" {
		return ""
	}
	partXML := readZipBytes(fileMap, partPath)
	if partXML == nil {
		return ""
	}
	return smartart.ExtractUniqueID(string(partXML))
}
