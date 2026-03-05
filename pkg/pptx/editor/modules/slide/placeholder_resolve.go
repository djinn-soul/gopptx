package slide

import (
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

type PlaceholderMeta struct {
	Name  string
	Type  string
	Index int
}

type ParsePlaceholdersFn func(content []byte) []PlaceholderMeta

func LookupSlidePlaceholders(
	slidePart string,
	getPart func(string) ([]byte, bool),
	parsePlaceholders ParsePlaceholdersFn,
) ([]PlaceholderMeta, error) {
	if strings.TrimSpace(slidePart) == "" {
		return nil, nil
	}
	content, ok := getPart(slidePart)
	if !ok {
		return nil, fmt.Errorf("slide part %q missing for placeholder resolution", slidePart)
	}
	return parsePlaceholders(content), nil
}

func ResolvePlaceholderTarget(
	override shapes.PlaceholderContent,
	placeholders []PlaceholderMeta,
) (string, int, error) {
	targetType := strings.TrimSpace(override.Type)
	targetIndex := override.Index
	target := override.Target
	if target == nil {
		return targetType, targetIndex, nil
	}

	if t := strings.TrimSpace(target.Type); t != "" {
		return t, target.Index, nil
	}

	name := strings.TrimSpace(target.Name)
	if name == "" {
		return targetType, targetIndex, nil
	}

	matches := make([]PlaceholderMeta, 0, 1)
	for _, ph := range placeholders {
		if strings.EqualFold(strings.TrimSpace(ph.Name), name) {
			matches = append(matches, ph)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0].Type, matches[0].Index, nil
	case 0:
		return "", 0, fmt.Errorf("placeholder name %q not found", name)
	default:
		return "", 0, fmt.Errorf("placeholder name %q is ambiguous (%d matches)", name, len(matches))
	}
}

func EnsureSlideRelsExist(
	hasPart func(string) bool,
	slidePart string,
) error {
	relsPath := common.SlideRelsPartName(slidePart)
	if hasPart(relsPath) {
		return nil
	}
	return fmt.Errorf("missing slide relationships part %q", relsPath)
}
