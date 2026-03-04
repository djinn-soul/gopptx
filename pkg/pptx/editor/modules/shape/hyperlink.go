package shape

import (
	"errors"
	"fmt"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

const (
	actionSlideJump      = "ppaction://hlinksldjump"
	actionShowJumpPrefix = "ppaction://hlinkshowjump?jump="
	actionMacroPrefix    = "ppaction://macro?name="
)

func DeriveActionURL(hl *common.Hyperlink) string {
	if hl == nil {
		return ""
	}
	if hl.TargetSlide != nil {
		return actionSlideJump
	}
	if hl.TargetJump != nil && *hl.TargetJump != "" {
		return actionShowJumpPrefix + strings.TrimSpace(*hl.TargetJump)
	}
	if hl.Macro != nil && *hl.Macro != "" {
		return actionMacroPrefix + strings.TrimSpace(*hl.Macro)
	}
	return ""
}

func ValidateHyperlinkAction(hl *common.Hyperlink) error {
	if hl == nil {
		return nil
	}
	selectorCount := 0
	hasAddress := strings.TrimSpace(GetStr(hl.Address)) != ""
	hasTargetSlide := hl.TargetSlide != nil
	hasJump := strings.TrimSpace(GetStr(hl.TargetJump)) != ""
	hasMacro := strings.TrimSpace(GetStr(hl.Macro)) != ""
	if hasAddress {
		selectorCount++
	}
	if hasTargetSlide {
		selectorCount++
	}
	if hasJump {
		selectorCount++
	}
	if hasMacro {
		selectorCount++
	}
	if selectorCount > 1 {
		return errors.New(
			"hyperlink selectors are mutually exclusive: use only one of address, target_slide, jump, or macro",
		)
	}
	if hasJump {
		jump := strings.ToLower(strings.TrimSpace(*hl.TargetJump))
		switch jump {
		case "nextslide", "previousslide", "firstslide", "lastslide", "lastslideviewed", "endshow":
			return nil
		default:
			return fmt.Errorf("unsupported jump target %q", *hl.TargetJump)
		}
	}
	return nil
}
