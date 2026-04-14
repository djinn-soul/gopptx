package slide

import (
	"errors"
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func ValidateMergeEditorsNil(dstNil, srcNil bool) error {
	if dstNil || srcNil {
		return errors.New("editors cannot be nil")
	}
	return nil
}

func SetSlideTitleInState(
	slides []common.EditorSlideRef,
	index int,
	title string,
	getPart GetPartFn,
	setPart SetPartFn,
) ([]common.EditorSlideRef, error) {
	if index < 0 || index >= len(slides) {
		return nil, fmt.Errorf("slide index %d out of range", index)
	}

	ref := slides[index]
	content, ok := getPart(ref.Part)
	if !ok {
		return nil, fmt.Errorf("slide part %q missing", ref.Part)
	}

	newContent, modified := ReplaceAllTitleTextRuns(content, common.XMLEscape(title))
	if !modified {
		return nil, fmt.Errorf("slide %d has no title text run to update", index)
	}

	setPart(ref.Part, newContent)
	slides[index].Title = title
	return slides, nil
}
