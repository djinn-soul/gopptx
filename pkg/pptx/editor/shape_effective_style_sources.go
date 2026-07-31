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
	source         string
	fillRGB        string
	fillScheme     string
	fillTransforms []colorTransform
	fontRGB        string
	fontScheme     string
	fontTransforms []colorTransform
	typeface       string
	sizePt         float64
	bold           *bool
	italic         *bool
	position       *common.EffectivePosition
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
	chain := styleChain{}
	layoutPart, masterPart, err := e.GetSlideLayoutRef(slideIndex)
	if err != nil {
		chain.theme = e.themeStyleContext("")
		return chain
	}
	chain.layoutPart = layoutPart
	chain.masterPart = masterPart
	chain.theme = e.themeStyleContext(masterPart)
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
	if color, found := solidFillColorDetails(shapeXML); found {
		level.fillRGB = color.rgb
		level.fillScheme = color.scheme
		level.fillTransforms = color.transforms
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
	if color, found := solidFillColorDetails([]byte(defRPr)); found {
		level.fontRGB = color.rgb
		level.fontScheme = color.scheme
		level.fontTransforms = color.transforms
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

func solidFillColorDetails(block []byte) (styleColor, bool) {
	scope := string(block)
	if spPr := spPrBlockPattern.FindString(scope); spPr != "" {
		scope = spPr
	}
	fill := solidFillPattern.FindStringSubmatch(scope)
	if len(fill) < 2 {
		return styleColor{}, false
	}
	if match := srgbClrValPattern.FindStringSubmatch(fill[1]); len(match) > 1 {
		return styleColor{
			rgb:        strings.ToUpper(match[1]),
			transforms: parseColorTransforms(fill[1]),
		}, true
	}
	if match := schemeClrPattern.FindStringSubmatch(fill[1]); len(match) > 1 {
		return styleColor{
			scheme:     match[1],
			transforms: parseColorTransforms(fill[1]),
		}, true
	}
	return styleColor{}, false
}
