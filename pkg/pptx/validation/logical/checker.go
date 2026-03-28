package logical

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/validation/structural"
)

// Checker implements structural.Checker by running model-level validation.
type Checker struct{}

func (c *Checker) Check(ps structural.PartProvider) []structural.Issue {
	var issues []structural.Issue

	for _, p := range ps.Keys() {
		if !strings.HasPrefix(p, "ppt/slides/slide") || !strings.HasSuffix(p, ".xml") || strings.Contains(p, "_rels") {
			continue
		}

		index := parseSlideIndex(p)
		issues = append(issues, c.checkSlide(ps, p, index)...)
	}

	return issues
}

func (c *Checker) checkSlide(ps structural.PartProvider, slidePart string, index int) []structural.Issue {
	data, ok := ps.Get(slidePart)
	if !ok {
		return nil
	}

	// Extract slide content
	title := extractFirstAText(data)
	parsedShapes := parseSlideShapes(data)
	parsedImages := parseSlideImages(data)

	// Infer layout from slide content - if no title and no shapes, likely a blank layout
	layout := elements.SlideLayoutTitleAndContent
	if title == "" && len(parsedShapes) == 0 && len(parsedImages) == 0 {
		layout = elements.SlideLayoutBlank
	} else if title == "" {
		// If there's content but no title, could be title-only or custom layout
		// Be lenient and skip title validation by using blank layout
		layout = elements.SlideLayoutBlank
	}

	slide := elements.SlideContent{
		Title:  title,
		Layout: layout,
		Shapes: parsedShapes,
		Images: parsedImages,
	}

	if err := slide.Validate(index); err != nil {
		return []structural.Issue{{
			Code:        structural.CodeModelValidationError,
			Severity:    structural.SeverityWarning,
			Path:        slidePart,
			Description: fmt.Sprintf("Logical validation failed: %v", err),
			Repairable:  false,
		}}
	}

	return nil
}

func extractFirstAText(content []byte) string {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return ""
			}
			return ""
		}
		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "t" {
			continue
		}
		var value string
		if decodeErr := decoder.DecodeElement(&value, &start); decodeErr != nil {
			return ""
		}
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
}

type shapeNode struct {
	NvSpPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvSpPr"`
	NvPicPr struct {
		CNvPr struct {
			ID   int    `xml:"id,attr"`
			Name string `xml:"name,attr"`
		} `xml:"cNvPr"`
	} `xml:"nvPicPr"`
	BlipFill struct {
		Blip struct {
			Embed string `xml:"embed,attr"`
		} `xml:"blip"`
	} `xml:"blipFill"`
	SpPr struct {
		PrstGeom struct {
			Prst string `xml:"prst,attr"`
		} `xml:"prstGeom"`
		Xfrm struct {
			Off struct {
				X int64 `xml:"x,attr"`
				Y int64 `xml:"y,attr"`
			} `xml:"off"`
			Ext struct {
				Cx int64 `xml:"cx,attr"`
				Cy int64 `xml:"cy,attr"`
			} `xml:"ext"`
		} `xml:"xfrm"`
	} `xml:"spPr"`
	TxBody struct {
		P []struct {
			R []struct {
				T string `xml:"t"`
			} `xml:"r"`
		} `xml:"p"`
	} `xml:"txBody"`
}

func parseSlideShapes(content []byte) []shapes.Shape {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	result := make([]shapes.Shape, 0)

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result
			}
			return result
		}

		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "sp" {
			continue
		}

		var node shapeNode
		if decodeErr := decoder.DecodeElement(&node, &start); decodeErr != nil {
			continue
		}

		shape := shapes.Shape{
			Type: node.SpPr.PrstGeom.Prst,
			X:    styling.Emu(node.SpPr.Xfrm.Off.X),
			Y:    styling.Emu(node.SpPr.Xfrm.Off.Y),
			CX:   styling.Emu(node.SpPr.Xfrm.Ext.Cx),
			CY:   styling.Emu(node.SpPr.Xfrm.Ext.Cy),
			Name: extractShapeName(node),
			Text: extractShapeText(node),
		}
		result = append(result, shape)
	}
}

func parseSlideImages(content []byte) []shapes.Image {
	decoder := xml.NewDecoder(bytes.NewReader(content))
	result := make([]shapes.Image, 0)

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result
			}
			return result
		}

		start, ok := token.(xml.StartElement)
		if !ok || start.Name.Local != "pic" {
			continue
		}

		var node shapeNode
		if decodeErr := decoder.DecodeElement(&node, &start); decodeErr != nil {
			continue
		}

		relID := node.BlipFill.Blip.Embed
		if relID == "" {
			relID = "embedded"
		}
		img := shapes.Image{
			RelID: relID,
			X:     styling.Emu(node.SpPr.Xfrm.Off.X),
			Y:     styling.Emu(node.SpPr.Xfrm.Off.Y),
			CX:    styling.Emu(node.SpPr.Xfrm.Ext.Cx),
			CY:    styling.Emu(node.SpPr.Xfrm.Ext.Cy),
		}
		result = append(result, img)
	}
}

func extractShapeName(node shapeNode) string {
	if node.NvSpPr.CNvPr.ID != 0 {
		return node.NvSpPr.CNvPr.Name
	}
	return node.NvPicPr.CNvPr.Name
}

func extractShapeText(node shapeNode) string {
	var builder strings.Builder
	for pIdx, p := range node.TxBody.P {
		for _, r := range p.R {
			builder.WriteString(r.T)
		}
		if pIdx < len(node.TxBody.P)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func parseSlideIndex(partPath string) int {
	base := path.Base(partPath)
	numStr := strings.TrimSuffix(strings.TrimPrefix(base, "slide"), ".xml")
	num, _ := strconv.Atoi(numStr)
	return num
}
