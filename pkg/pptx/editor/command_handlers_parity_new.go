package editor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// UpdateSmartArt replaces text items in an existing SmartArt diagram.
// It finds the graphic frame by shapeID, resolves its r:dm rel to the data
// part, and rewrites <a:t> nodes inside <dgm:pt> elements.
func (e *PresentationEditor) UpdateSmartArt(slideIndex, shapeID int, items []string) error {
	if slideIndex < 0 || slideIndex >= len(e.slides) {
		return fmt.Errorf("slide index %d out of range", slideIndex)
	}
	slideRef := e.slides[slideIndex]
	slideXML, ok := e.parts.Get(slideRef.Part)
	if !ok {
		return fmt.Errorf("slide part %q not found", slideRef.Part)
	}

	dmRelID := extractSmartArtDataRelID(string(slideXML), shapeID)
	if dmRelID == "" {
		return fmt.Errorf("shapeID %d not found or is not a SmartArt graphic frame", shapeID)
	}

	relsPath := common.RelsPathFor(slideRef.Part)
	relsData, ok := e.parts.Get(relsPath)
	if !ok {
		return fmt.Errorf("rels part %q not found", relsPath)
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return fmt.Errorf("parse rels: %w", err)
	}
	var dataPartPath string
	for _, rel := range rels {
		if rel.ID == dmRelID {
			dataPartPath = common.ResolveRelationshipTarget(slideRef.Part, rel.Target)
			break
		}
	}
	if dataPartPath == "" {
		return fmt.Errorf("rel %q not resolved in %s", dmRelID, relsPath)
	}

	dataXML, ok := e.parts.Get(dataPartPath)
	if !ok {
		return fmt.Errorf("SmartArt data part %q not found", dataPartPath)
	}
	updated := rewriteSmartArtTextItems(string(dataXML), items)
	e.parts.Set(dataPartPath, []byte(updated))
	return nil
}

// extractSmartArtDataRelID locates the r:dm attribute for the given shapeID.
func extractSmartArtDataRelID(slideXML string, shapeID int) string {
	idPattern := fmt.Sprintf(`id="%d"`, shapeID)
	idx := strings.Index(slideXML, idPattern)
	if idx < 0 {
		return ""
	}
	dmIdx := strings.Index(slideXML[idx:], `r:dm=`)
	if dmIdx < 0 {
		return ""
	}
	re := regexp.MustCompile(`r:dm=["']([^"']+)["']`)
	m := re.FindStringSubmatch(slideXML[idx+dmIdx:])
	if len(m) < 2 { //nolint:mnd // index 1 is the capture group
		return ""
	}
	return m[1]
}

// rewriteSmartArtTextItems replaces <a:t> text content in sequence.
func rewriteSmartArtTextItems(dataXML string, items []string) string {
	itemIdx := 0
	re := regexp.MustCompile(`(<a:t>)[^<]*(</a:t>)`)
	return re.ReplaceAllStringFunc(dataXML, func(_ string) string {
		if itemIdx >= len(items) {
			return "<a:t></a:t>"
		}
		text := items[itemIdx]
		itemIdx++
		return "<a:t>" + xmlEscapeSimple(text) + "</a:t>"
	})
}

// xmlEscapeSimple escapes the minimal set of XML special characters.
func xmlEscapeSimple(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}

// handleUpdateSmartArt updates text items of an existing SmartArt diagram.
//
// Payload: {"slide_index": N, "shape_id": N, "items": ["text1", ...]}.
// Response: {"updated": true}.
func handleUpdateSmartArt(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	shapeID, ok := v.RequireInt(p, "shape_id")
	if !ok {
		return nil, v.Error()
	}
	items, _ := v.OptionalStringSlice(p, "items")
	if updateErr := e.UpdateSmartArt(slideIndex, shapeID, items); updateErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, updateErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}

// handleSetSlideBackground sets the background of a slide.
//
// Payload: {"slide_index": N, "type": "solid|gradient|image|theme", ...}.
// Response: {"updated": true}.
func handleSetSlideBackground(e *PresentationEditor, payload json.RawMessage) (any, error) {
	p, err := ParseRawPayload(payload)
	if err != nil {
		return nil, err
	}
	v := NewPayloadValidator()
	slideIndex, ok := requireSlideIndex(e, p, v)
	if !ok {
		return nil, v.Error()
	}
	bgType, ok := v.RequireString(p, "type")
	if !ok {
		return nil, v.Error()
	}

	bg := SlideBackground{Type: bgType}
	bg.Color = v.OptionalString(p, "color")
	bg.ColorRef = v.OptionalString(p, "color_ref")
	bg.ImagePath = v.OptionalString(p, "image_path")
	bg.ImageData = v.OptionalString(p, "image_data")
	bg.Angle, _ = v.OptionalInt(p, "angle")
	bg.Colors, _ = v.OptionalStringSlice(p, "colors")

	if setErr := e.SetSlideBackground(slideIndex, bg); setErr != nil {
		return nil, NewBridgeError(ErrCodeOpFailed, setErr.Error())
	}
	return map[string]bool{"updated": true}, nil
}
