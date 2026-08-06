package export

import (
	"archive/zip"
	"encoding/xml"
	"maps"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// A slide states only what differs from the layout it was made with. Its
// placeholders usually carry no geometry at all, its background is normally the
// master's, and any logo or rule the template draws lives on the layout or the
// master rather than on the slide.
//
// inheritedLook is what one slide picks up from that chain: where its
// placeholders sit, what colour the page is, and the shapes drawn beneath its
// own.
type inheritedLook struct {
	// placeholders is keyed by placeholder type and index, as the slide names
	// them in <p:ph>.
	placeholders map[placeholderKey]placeholderBox
	background   *elements.SlideBackground
	// shapes are the layout's and master's own shapes, master first, in the
	// order they paint.
	shapes []shapes.Shape
	// titleSizePt and bodySizePt are the master's txStyles defaults, which a
	// placeholder run that states no size of its own inherits.
	titleSizePt int
	bodySizePt  int
}

type placeholderKey struct {
	phType string
	idx    int
}

// placeholderBox is a placeholder's geometry in EMU, as the layout or master
// states it.
type placeholderBox struct {
	X, Y, CX, CY int64
	SizePt       int
}

// extractInheritedLooks returns, per slide (0-based), what that slide inherits
// from its layout and master.
func extractInheritedLooks(pptxPath string, theme deckTheme) []inheritedLook {
	zr, err := zip.OpenReader(pptxPath)
	if err != nil {
		return nil
	}
	defer zr.Close()

	fileMap := make(map[string]*zip.File, len(zr.File))
	for _, f := range zr.File {
		fileMap[canonicalZipPath(f.Name)] = f
	}
	slideOrder := resolveSlideOrder(fileMap)

	// A deck reuses a handful of layouts across many slides, so each one is
	// parsed once.
	cache := map[string]inheritedLook{}
	out := make([]inheritedLook, len(slideOrder))
	for i, slidePart := range slideOrder {
		layoutPart := relatedPart(fileMap, slidePart, "slideLayout")
		if layoutPart == "" {
			continue
		}
		look, ok := cache[layoutPart]
		if !ok {
			look = buildInheritedLook(fileMap, layoutPart, theme)
			cache[layoutPart] = look
		}
		if !slideShowsMasterShapes(fileMap, slidePart) {
			look.shapes = nil
		}
		out[i] = look
	}
	return out
}

// slideShowsMasterShapes reports whether the slide draws what its layout and
// master draw. A slide with showMasterSp="0" hides them.
func slideShowsMasterShapes(fileMap map[string]*zip.File, slidePart string) bool {
	data := readZipBytes(fileMap, slidePart)
	if data == nil {
		return true
	}
	var doc struct {
		ShowMasterSp *string `xml:"showMasterSp,attr"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return true
	}
	if doc.ShowMasterSp == nil {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(*doc.ShowMasterSp)) {
	case "0", "false":
		return false
	default:
		return true
	}
}

// inheritedLookFor returns what slide idx inherits, or an empty look.
func inheritedLookFor(looks []inheritedLook, idx int) inheritedLook {
	if idx < 0 || idx >= len(looks) {
		return inheritedLook{}
	}
	return looks[idx]
}

// applyInheritedLook fills in everything the slide left to its layout and
// master: the placeholder geometry its own <p:spPr> omits, the background, and
// the shapes drawn beneath its own.
func applyInheritedLook(sc *elements.SlideContent, order *slidePaintOrder, look inheritedLook) {
	if sc.Background == nil && look.background != nil {
		background := *look.background
		sc.Background = &background
	}
	applyInheritedPlaceholder(&sc.TitleBoundsEMU, &sc.TitleSize, look, titlePlaceholderKeys())
	applyInheritedPlaceholder(&sc.ContentBoundsEMU, &sc.ContentSize, look, bodyPlaceholderKeys())
	if sc.TitleSize <= 0 {
		sc.TitleSize = look.titleSizePt
	}
	if sc.ContentSize <= 0 {
		sc.ContentSize = look.bodySizePt
	}
	prependInheritedShapes(sc, order, look.shapes)
}

func titlePlaceholderKeys() []placeholderKey {
	return []placeholderKey{{phType: "title"}, {phType: "ctrTitle"}, {phType: placeholderCtrTitle}}
}

func bodyPlaceholderKeys() []placeholderKey {
	return []placeholderKey{
		{phType: placeholderBody, idx: 1},
		{phType: placeholderBody},
		{phType: placeholderSubtitle, idx: 1},
		{phType: placeholderSubtitle},
	}
}

// applyInheritedPlaceholder copies a layout placeholder's geometry onto the
// slide when the slide states none of its own.
func applyInheritedPlaceholder(bounds *[4]int64, sizePt *int, look inheritedLook, keys []placeholderKey) {
	if bounds[2] > 0 || bounds[3] > 0 {
		return
	}
	for _, key := range keys {
		box, ok := look.placeholders[placeholderKey{phType: strings.ToLower(key.phType), idx: key.idx}]
		if !ok {
			continue
		}
		*bounds = [4]int64{box.X, box.Y, box.CX, box.CY}
		if *sizePt <= 0 && box.SizePt > 0 {
			*sizePt = box.SizePt
		}
		return
	}
}

// prependInheritedShapes puts the layout's and master's shapes behind the
// slide's own.
//
// They only take part when the slide's shape tree was read: without it the
// renderer falls back to slice order, and a shape with no recorded position
// would paint over the slide instead of under it.
func prependInheritedShapes(sc *elements.SlideContent, order *slidePaintOrder, inherited []shapes.Shape) {
	if len(inherited) == 0 || order == nil || !order.known {
		return
	}
	behind := make([]shapes.Shape, 0, len(inherited)+len(sc.Shapes))
	for i, shape := range inherited {
		// Negative, so they sort below every element of the slide's own tree,
		// keeping master-before-layout order among themselves.
		shape.ZOrder = i - len(inherited)
		behind = append(behind, shape)
	}
	sc.Shapes = append(behind, sc.Shapes...)
}

// buildInheritedLook reads a layout and the master behind it.
func buildInheritedLook(fileMap map[string]*zip.File, layoutPart string, theme deckTheme) inheritedLook {
	look := inheritedLook{placeholders: map[placeholderKey]placeholderBox{}}
	masterPart := relatedPart(fileMap, layoutPart, "slideMaster")

	// The master is read first: its placeholders are the fallback for the ones
	// the layout does not restate, and its shapes paint furthest back.
	if masterPart != "" {
		if data := readZipBytes(fileMap, masterPart); data != nil {
			applyPartToLook(&look, data, theme)
		}
	}
	if data := readZipBytes(fileMap, layoutPart); data != nil {
		applyPartToLook(&look, data, theme)
	}
	return look
}

// applyPartToLook folds one layout or master part into the inherited look. A
// later part overrides the placeholders and background of an earlier one and
// paints its shapes above them.
func applyPartToLook(look *inheritedLook, data []byte, theme deckTheme) {
	part := parseLayoutPart(data, theme)
	maps.Copy(look.placeholders, part.placeholders)
	if part.background != nil {
		look.background = part.background
	}
	if part.titleSizePt > 0 {
		look.titleSizePt = part.titleSizePt
	}
	if part.bodySizePt > 0 {
		look.bodySizePt = part.bodySizePt
	}
	look.shapes = append(look.shapes, part.shapes...)
}

// layoutPart is one parsed layout or master.
