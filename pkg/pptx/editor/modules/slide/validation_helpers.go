package slide

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

type ParseRelationshipsFn func([]byte) ([]common.EditorRelationship, error)

func Relationships(
	slidePart string,
	getPart func(string) ([]byte, bool),
	parseRelationships ParseRelationshipsFn,
) ([]common.EditorRelationship, error) {
	relsPart := common.SlideRelsPartName(slidePart)
	data, ok := getPart(relsPart)
	if !ok {
		return nil, fmt.Errorf("missing slide relationships part %q", relsPart)
	}
	rels, err := parseRelationships(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", relsPart, err)
	}
	return rels, nil
}

func ScanSupportedSlideRels(rels []common.EditorRelationship) (string, error) {
	notesTarget := ""
	for _, rel := range rels {
		switch rel.Type {
		case common.RelTypeSlideLayout:
		case common.RelTypeNotesSlide:
			notesTarget = rel.Target
		case common.RelTypeHyperlink:
		case common.RelTypeImage, common.RelTypeChart, common.RelTypeAudio, common.RelTypeVideo, common.RelTypeTheme:
		default:
			return "", fmt.Errorf("unsupported relationship type %q", rel.Type)
		}
	}
	return notesTarget, nil
}

func HasSlideLayoutRelationship(rels []common.EditorRelationship) bool {
	for _, rel := range rels {
		if rel.Type == common.RelTypeSlideLayout {
			return true
		}
	}
	return false
}

func HasImageContent(slide elements.SlideContent) bool {
	if len(slide.Images) > 0 {
		return true
	}
	if slide.Background != nil && slide.Background.Type == elements.SlideBackgroundPicture &&
		slide.Background.PictureFill != nil {
		return true
	}
	for _, override := range slide.PlaceholderOverrides {
		if override.Image != nil {
			return true
		}
	}
	return false
}

func ValidateEditorSlideContent(slide elements.SlideContent) error {
	return slide.Validate(1)
}
