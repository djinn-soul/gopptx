package editor

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

// Placeholder types found on a layout.
const (
	placeholderTypeSubTitle = "subTitle"
	placeholderTypeDate     = "dt"
	placeholderTypeFooter   = "ftr"
	placeholderTypeSlideNum = "sldNum"
)

// Placeholder types a fresh slide does not repeat from its layout. PowerPoint
// draws the date, footer and slide number from the layout unless the deck's
// header/footer settings put them on the slide, which is a separate call.
//
//nolint:gochecknoglobals // immutable lookup set
var layoutOnlyPlaceholderTypes = map[string]bool{
	placeholderTypeDate:     true,
	placeholderTypeFooter:   true,
	placeholderTypeSlideNum: true,
}

// AddSlideFromLayoutPart adds a slide whose placeholders come from a named
// layout part, rather than from one of the built-in layout modes.
//
// Building a default slide and then only retargeting its layout relationship is
// not enough: the slide keeps the default title and body placeholders whatever
// the layout says, so a blank layout still shows two boxes and a layout's own
// picture or extra content placeholders never appear at all. The placeholders
// are created here to match the layout instead (upstream issue #175).
func (e *PresentationEditor) AddSlideFromLayoutPart(
	layoutPart, title string,
	bullets []string,
) (int, error) {
	if e == nil {
		return 0, errors.New("editor cannot be nil")
	}
	layoutPart = common.CanonicalPartPath(layoutPart)
	if !e.parts.Has(layoutPart) {
		return 0, fmt.Errorf("layout part %s not found", layoutPart)
	}

	// A layout that declares no placeholders says nothing about what the slide
	// should hold, so the built-in title-and-content shapes are kept and only
	// the binding changes. Building strictly from such a layout would drop the
	// caller's title on the floor.
	content := elements.NewSlide("").WithLayout(elements.SlideLayoutBlank)
	placeholders := e.layoutPlaceholdersForNewSlide(layoutPart)
	if len(placeholders) == 0 {
		content = elements.NewSlide(title)
		for _, bullet := range bullets {
			content = content.AddBullet(bullet)
		}
	}

	index, err := e.AddSlide(content)
	if err != nil {
		return 0, err
	}
	if err := e.RebindSlideLayout(index, layoutPart); err != nil {
		return 0, err
	}
	if err := e.populateSlideFromLayout(index, placeholders, title, bullets); err != nil {
		return 0, err
	}
	return index, nil
}

// layoutPlaceholdersForNewSlide returns the layout placeholders a new slide
// should carry, dropping the date, footer and slide number.
func (e *PresentationEditor) layoutPlaceholdersForNewSlide(
	layoutPart string,
) []common.PlaceholderInfo {
	var kept []common.PlaceholderInfo
	for _, placeholder := range e.GetLayoutPlaceholders(layoutPart) {
		if layoutOnlyPlaceholderTypes[placeholder.Type] {
			continue
		}
		kept = append(kept, placeholder)
	}
	return kept
}

// populateSlideFromLayout writes one placeholder shape per layout placeholder,
// filling the title and body ones with the caller's text.
func (e *PresentationEditor) populateSlideFromLayout(
	slideIndex int,
	placeholders []common.PlaceholderInfo,
	title string,
	bullets []string,
) error {
	if len(placeholders) == 0 {
		return nil
	}
	partPath := e.slides[slideIndex].Part
	content, ok := e.parts.Get(partPath)
	if !ok {
		return fmt.Errorf("read slide part %s: not found", partPath)
	}
	nextID := e.maxObjectID(partPath, content) + 1

	bodyUsed := false
	for _, placeholder := range placeholders {
		text := ""
		switch {
		case isTitlePlaceholderType(placeholder.Type):
			text = title
		case !bodyUsed && isBodyPlaceholderType(placeholder.Type):
			text = strings.Join(bullets, "\n")
			bodyUsed = true
		}

		shapeXML := slidePlaceholderShapeXML(nextID, placeholder, text)
		updated, err := insertShapeIntoSlideTree(content, lastShapeEndOffset(content), shapeXML)
		if err != nil {
			return err
		}
		content = updated
		e.reserveObjectIDs(partPath, nextID)
		nextID++
	}
	e.parts.Set(partPath, content)
	return nil
}

// lastShapeEndOffset reports where the last shape in a slide's tree ends, which
// is where the next one is inserted. A tree with no shapes yet reports -1.
func lastShapeEndOffset(content []byte) int64 {
	shapeNodes, err := scanShapesWithOffsets(content, true)
	if err != nil {
		return -1
	}
	last := int64(-1)
	for _, shape := range shapeNodes {
		if shape.End > last {
			last = shape.End
		}
	}
	return last
}

func isTitlePlaceholderType(phType string) bool {
	return phType == placeholderTypeTitle || phType == placeholderTypeCenteredTitle
}

// isBodyPlaceholderType reports whether a placeholder takes the body text. An
// empty type is body: that is the schema default for p:ph.
func isBodyPlaceholderType(phType string) bool {
	return phType == "" ||
		phType == placeholderTypeBody ||
		phType == placeholderTypeSubTitle
}

// slidePlaceholderShapeXML renders one placeholder shape.
//
// No a:xfrm is written: without one the placeholder inherits its position and
// size from the layout, which is what keeps the slide following the layout when
// the layout is edited.
func slidePlaceholderShapeXML(id int, placeholder common.PlaceholderInfo, text string) string {
	name := placeholder.Name
	if strings.TrimSpace(name) == "" {
		name = fmt.Sprintf("Placeholder %d", id)
	}

	var ph strings.Builder
	ph.WriteString(`<p:ph`)
	if placeholder.Type != "" {
		ph.WriteString(` type="` + common.XMLEscape(placeholder.Type) + `"`)
	}
	if placeholder.Index > 0 {
		fmt.Fprintf(&ph, ` idx="%d"`, placeholder.Index)
	}
	ph.WriteString(`/>`)

	var body strings.Builder
	if text == "" {
		body.WriteString(`<a:p/>`)
	} else {
		for line := range strings.SplitSeq(text, "\n") {
			body.WriteString(`<a:p><a:r><a:rPr lang="en-US" dirty="0"/><a:t>`)
			body.WriteString(common.XMLEscape(line))
			body.WriteString(`</a:t></a:r></a:p>`)
		}
	}

	return fmt.Sprintf(
		`<p:sp xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" `+
			`xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main">`+
			`<p:nvSpPr><p:cNvPr id="%d" name="%s"/>`+
			`<p:cNvSpPr><a:spLocks noGrp="1"/></p:cNvSpPr>`+
			`<p:nvPr>%s</p:nvPr></p:nvSpPr>`+
			`<p:spPr/>`+
			`<p:txBody><a:bodyPr/><a:lstStyle/>%s</p:txBody></p:sp>`,
		id, common.XMLEscape(name), ph.String(), body.String(),
	)
}
