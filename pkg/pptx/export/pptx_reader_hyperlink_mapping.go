package export

import (
	"net/url"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/action"
	editorcommon "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func editorHyperlinkToExportHyperlink(src *editorcommon.Hyperlink) *action.Hyperlink {
	if src == nil {
		return nil
	}

	var link action.Hyperlink
	switch {
	case src.TargetSlide != nil:
		link = action.NewHyperlink(action.HyperlinkSlide(uint32(*src.TargetSlide + 1)))
	case src.TargetJump != nil:
		jumpAction, ok := editorJumpToAction(*src.TargetJump)
		if !ok {
			return nil
		}
		link = action.NewHyperlink(jumpAction)
	case src.Address != nil:
		link = action.NewHyperlink(editorAddressToAction(strings.TrimSpace(*src.Address), strings.TrimSpace(getStr(src.Action))))
	default:
		return nil
	}

	if src.Tooltip != nil {
		link = link.WithTooltip(strings.TrimSpace(*src.Tooltip))
	}
	if src.HighlightClick != nil {
		link = link.WithHighlightClick(*src.HighlightClick)
	}
	return &link
}

func editorAddressToAction(address, actionValue string) action.HyperlinkAction {
	switch {
	case strings.HasPrefix(strings.ToLower(actionValue), "ppaction://program"):
		return action.HyperlinkProgram(decodeFileAddress(address))
	case strings.HasPrefix(strings.ToLower(address), "mailto:"):
		return editorMailtoToAction(address)
	case strings.HasPrefix(strings.ToLower(address), "file:///"):
		return action.HyperlinkFile(decodeFileAddress(address))
	default:
		return action.HyperlinkURL(address)
	}
}

func editorMailtoToAction(address string) action.HyperlinkAction {
	parsed, err := url.Parse(address)
	if err != nil {
		return action.HyperlinkURL(address)
	}
	subject := parsed.Query().Get("subject")
	mailbox := strings.TrimPrefix(address, "mailto:")
	if at := strings.Index(mailbox, "?"); at >= 0 {
		mailbox = mailbox[:at]
	}
	if subject != "" {
		return action.HyperlinkEmailWithSubject(mailbox, subject)
	}
	return action.HyperlinkEmail(mailbox)
}

func decodeFileAddress(address string) string {
	pathValue := strings.TrimPrefix(address, "file:///")
	if decoded, err := url.PathUnescape(pathValue); err == nil {
		pathValue = decoded
	}
	return strings.ReplaceAll(pathValue, "/", "\\")
}

func editorJumpToAction(jump string) (action.HyperlinkAction, bool) {
	switch strings.ToLower(strings.TrimSpace(jump)) {
	case "firstslide":
		return action.HyperlinkFirstSlide(), true
	case "lastslide":
		return action.HyperlinkLastSlide(), true
	case "previousslide":
		return action.HyperlinkPreviousSlide(), true
	case "nextslide":
		return action.HyperlinkNextSlide(), true
	case "endshow":
		return action.HyperlinkEndShow(), true
	default:
		return action.HyperlinkAction{}, false
	}
}

func getStr(src *string) string {
	if src == nil {
		return ""
	}
	return *src
}
