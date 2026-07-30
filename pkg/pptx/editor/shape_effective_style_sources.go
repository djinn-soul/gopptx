package editor

import (
	"regexp"
	"strconv"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// One level of the inheritance chain, flattened to the values a style lookup
// cares about. Empty or nil means the level says nothing about that value.
type inheritedStyle struct {
	source     string
	fillRGB    string
	fillScheme string
	fontRGB    string
	fontScheme string
	typeface   string
	sizePt     float64
	bold       *bool
	italic     *bool
	position   *common.EffectivePosition
}

type styleChain struct {
	layoutPart string
	masterPart string
	levels     []inheritedStyle
	theme      themeStyleContext
}

type themeStyleContext struct {
	scheme     ThemeColorScheme
	majorLatin string
	minorLatin string
}

var (
	spPrBlockPattern   = regexp.MustCompile(`(?s)<(?:\w+:)?spPr\b[^>]*>.*?</(?:\w+:)?spPr>`)
	solidFillPattern   = regexp.MustCompile(`(?s)<a:solidFill\b[^>]*>(.*?)</a:solidFill>`)
	srgbClrValPattern  = regexp.MustCompile(`<a:srgbClr\b[^>]*\bval="([0-9A-Fa-f]{6})"`)
	schemeClrPattern   = regexp.MustCompile(`<a:schemeClr\b[^>]*\bval="([^"]+)"`)
	lvl1PrPattern      = regexp.MustCompile(`(?s)<a:lvl1pPr\b[^>]*>.*?</a:lvl1pPr>`)
	defRPrPattern      = regexp.MustCompile(`(?s)<a:defRPr\b[^>]*>.*?</a:defRPr>|<a:defRPr\b[^>]*/>`)
	latinTypefacePat   = regexp.MustCompile(`<a:latin\b[^>]*\btypeface="([^"]*)"`)
	rPrSizeAttrPattern = regexp.MustCompile(`\bsz="(\d+)"`)
	rPrBoldAttrPattern = regexp.MustCompile(`\bb="([01])"`)
	rPrItalAttrPattern = regexp.MustCompile(`\bi="([01])"`)
	majorLatinPattern  = regexp.MustCompile(`(?s)<(?:\w+:)?majorFont>.*?<a:latin\b[^>]*typeface="([^"]*)"`)
	minorLatinPattern  = regexp.MustCompile(`(?s)<(?:\w+:)?minorFont>.*?<a:latin\b[^>]*typeface="([^"]*)"`)
)

const (
	// hundredthsOfPoint is the unit of the <a:defRPr sz=""> attribute.
	hundredthsOfPoint = 100.0

	placeholderTypeCenteredTitle = "ctrTitle"
	placeholderTypeBody          = "body"
)

// styleInheritanceChain collects the layout placeholder, the master
// placeholder and the master text styles that apply to one shape.
func (e *PresentationEditor) styleInheritanceChain(slideIndex int, shape *parsedShape) styleChain {
	chain := styleChain{theme: e.themeStyleContext()}
	layoutPart, masterPart, err := e.GetSlideLayoutRef(slideIndex)
	if err != nil {
		return chain
	}
	chain.layoutPart = layoutPart
	chain.masterPart = masterPart
	if shape.PhType == "" && shape.PhIndex < 0 {
		// Only a placeholder inherits from the layout and master.
		return chain
	}

	if level, ok := e.placeholderStyleFrom(layoutPart, shape, common.StyleSourceLayout); ok {
		chain.levels = append(chain.levels, level)
	}
	if level, ok := e.placeholderStyleFrom(masterPart, shape, common.StyleSourceMaster); ok {
		chain.levels = append(chain.levels, level)
	}
	if level, ok := e.masterTextStyleFor(masterPart, shape.PhType); ok {
		chain.levels = append(chain.levels, level)
	}
	return chain
}

// placeholderStyleFrom finds the placeholder in a layout or master that the
// slide shape inherits from, matching on index first and type second, which is
// the order PowerPoint uses.
func (e *PresentationEditor) placeholderStyleFrom(
	partPath string,
	shape *parsedShape,
	source string,
) (inheritedStyle, bool) {
	if partPath == "" {
		return inheritedStyle{}, false
	}
	content, ok := e.parts.Get(partPath)
	if !ok {
		return inheritedStyle{}, false
	}
	candidates, err := parseSlideShapes(content)
	if err != nil {
		return inheritedStyle{}, false
	}

	match := matchPlaceholder(candidates, shape)
	if match == nil {
		return inheritedStyle{}, false
	}
	return placeholderStyle(content[match.Start:match.End], match, source), true
}

func matchPlaceholder(candidates []parsedShape, shape *parsedShape) *parsedShape {
	var byType *parsedShape
	for i := range candidates {
		candidate := &candidates[i]
		if candidate.PhType == "" && candidate.PhIndex < 0 {
			continue
		}
		if shape.PhIndex >= 0 && candidate.PhIndex == shape.PhIndex && shape.PhType == candidate.PhType {
			return candidate
		}
		if byType == nil && shape.PhType != "" && candidate.PhType == shape.PhType {
			byType = candidate
		}
	}
	return byType
}

func placeholderStyle(shapeXML []byte, shape *parsedShape, source string) inheritedStyle {
	level := inheritedStyle{source: source}
	if rgb, slot, found := solidFillColor(shapeXML); found {
		level.fillRGB = rgb
		level.fillScheme = slot
	}
	if shape.W > 0 && shape.H > 0 {
		level.position = &common.EffectivePosition{X: shape.X, Y: shape.Y, W: shape.W, H: shape.H}
	}
	applyDefRPr(&level, firstLevelDefRPr(string(shapeXML)))
	return level
}

// masterTextStyleFor reads <p:txStyles>, the master-wide defaults that apply
// when no placeholder in the chain states a value.
func (e *PresentationEditor) masterTextStyleFor(
	masterPart, placeholderType string,
) (inheritedStyle, bool) {
	if masterPart == "" {
		return inheritedStyle{}, false
	}
	content, ok := e.parts.Get(masterPart)
	if !ok {
		return inheritedStyle{}, false
	}

	styleTag := masterStyleTagFor(placeholderType)
	block := regexp.MustCompile(
		`(?s)<p:` + styleTag + `\b[^>]*>.*?</p:` + styleTag + `>`,
	).FindString(string(content))
	if block == "" {
		return inheritedStyle{}, false
	}

	level := inheritedStyle{source: common.StyleSourceMaster}
	applyDefRPr(&level, firstLevelDefRPr(block))
	return level, true
}

func masterStyleTagFor(placeholderType string) string {
	switch placeholderType {
	case placeholderTypeTitle, placeholderTypeCenteredTitle:
		return "titleStyle"
	case placeholderTypeBody, "subTitle", "obj", "":
		return "bodyStyle"
	default:
		return "otherStyle"
	}
}

// firstLevelDefRPr returns the <a:defRPr> of the first outline level, which is
// the one a placeholder's first paragraph draws with.
func firstLevelDefRPr(block string) string {
	scope := block
	if lvl1 := lvl1PrPattern.FindString(block); lvl1 != "" {
		scope = lvl1
	}
	return defRPrPattern.FindString(scope)
}

func applyDefRPr(level *inheritedStyle, defRPr string) {
	if defRPr == "" {
		return
	}
	if rgb, slot, found := solidFillColor([]byte(defRPr)); found {
		level.fontRGB = rgb
		level.fontScheme = slot
	}
	if match := latinTypefacePat.FindStringSubmatch(defRPr); len(match) > 1 {
		level.typeface = match[1]
	}
	if match := rPrSizeAttrPattern.FindStringSubmatch(defRPr); len(match) > 1 {
		if size, err := strconv.Atoi(match[1]); err == nil {
			level.sizePt = float64(size) / hundredthsOfPoint
		}
	}
	if match := rPrBoldAttrPattern.FindStringSubmatch(defRPr); len(match) > 1 {
		value := match[1] == "1"
		level.bold = &value
	}
	if match := rPrItalAttrPattern.FindStringSubmatch(defRPr); len(match) > 1 {
		value := match[1] == "1"
		level.italic = &value
	}
}

// solidFillColor reads the first <a:solidFill> of a shape's <p:spPr>, or of the
// block itself when it carries no spPr (a <a:defRPr>, for instance). It reports
// the literal RGB and, separately, the scheme slot when the fill names one.
func solidFillColor(block []byte) (string, string, bool) {
	scope := string(block)
	if spPr := spPrBlockPattern.FindString(scope); spPr != "" {
		scope = spPr
	}
	fill := solidFillPattern.FindStringSubmatch(scope)
	if len(fill) < 2 {
		return "", "", false
	}
	if match := srgbClrValPattern.FindStringSubmatch(fill[1]); len(match) > 1 {
		return strings.ToUpper(match[1]), "", true
	}
	if match := schemeClrPattern.FindStringSubmatch(fill[1]); len(match) > 1 {
		return "", match[1], true
	}
	return "", "", false
}

func (e *PresentationEditor) themeStyleContext() themeStyleContext {
	ctx := themeStyleContext{}
	if scheme, err := e.GetThemeColorScheme(); err == nil {
		ctx.scheme = scheme
	}
	inv, err := e.GetThemeInventory()
	if err != nil || len(inv.ThemeParts) == 0 {
		return ctx
	}
	data, ok := e.parts.Get(inv.ThemeParts[0])
	if !ok {
		return ctx
	}
	if match := majorLatinPattern.FindStringSubmatch(string(data)); len(match) > 1 {
		ctx.majorLatin = match[1]
	}
	if match := minorLatinPattern.FindStringSubmatch(string(data)); len(match) > 1 {
		ctx.minorLatin = match[1]
	}
	return ctx
}

// resolveTypeface expands the "+mj-lt"/"+mn-lt" references a placeholder uses
// to defer its font to the theme.
func (t themeStyleContext) resolveTypeface(value string) (string, bool) {
	switch value {
	case "+mj-lt", "+mj-ea", "+mj-cs":
		return t.majorLatin, t.majorLatin != ""
	case "+mn-lt", "+mn-ea", "+mn-cs":
		return t.minorLatin, t.minorLatin != ""
	default:
		return "", false
	}
}

// resolveSchemeSlot maps a scheme colour reference onto the theme's palette.
// tx1/bg1/tx2/bg2 are the slide-facing aliases of dk1/lt1/dk2/lt2.
func (t themeStyleContext) resolveSchemeSlot(slot string) (string, bool) {
	byName := map[string]string{
		themeSlotDk1: t.scheme.Dk1, "tx1": t.scheme.Dk1,
		themeSlotLt1: t.scheme.Lt1, "bg1": t.scheme.Lt1,
		themeSlotDk2: t.scheme.Dk2, "tx2": t.scheme.Dk2,
		themeSlotLt2: t.scheme.Lt2, "bg2": t.scheme.Lt2,
		themeSlotAccent1: t.scheme.Accent1, themeSlotAccent2: t.scheme.Accent2,
		themeSlotAccent3: t.scheme.Accent3, themeSlotAccent4: t.scheme.Accent4,
		themeSlotAccent5: t.scheme.Accent5, themeSlotAccent6: t.scheme.Accent6,
		themeSlotHlink: t.scheme.Hlink, themeSlotFolHlink: t.scheme.FolHlink,
	}
	rgb, ok := byName[slot]
	return rgb, ok && rgb != ""
}
