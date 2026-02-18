package action

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// HyperlinkActionType defines the type of hyperlink action.
type HyperlinkActionType string

const (
	// HyperlinkActionURL links to an external URL.
	HyperlinkActionURL HyperlinkActionType = "url"
	// HyperlinkActionSlide links to a specific slide number.
	HyperlinkActionSlide HyperlinkActionType = "slide"
	// HyperlinkActionFirstSlide links to the first slide.
	HyperlinkActionFirstSlide HyperlinkActionType = "firstSlide"
	// HyperlinkActionLastSlide links to the last slide.
	HyperlinkActionLastSlide HyperlinkActionType = "lastSlide"
	// HyperlinkActionNextSlide links to the next slide.
	HyperlinkActionNextSlide HyperlinkActionType = "nextSlide"
	// HyperlinkActionPreviousSlide links to the previous slide.
	HyperlinkActionPreviousSlide HyperlinkActionType = "previousSlide"
	// HyperlinkActionEndShow ends the slideshow.
	HyperlinkActionEndShow HyperlinkActionType = "endShow"
	// HyperlinkActionEmail links to an email address.
	HyperlinkActionEmail HyperlinkActionType = "email"
	// HyperlinkActionFile links to a file.
	HyperlinkActionFile HyperlinkActionType = "file"
	// HyperlinkActionProgram links to an external program.
	HyperlinkActionProgram HyperlinkActionType = "program"
)

// HyperlinkAction defines the target of a hyperlink.
type HyperlinkAction struct {
	Type         HyperlinkActionType
	URL          string // For URL type
	SlideNumber  uint32 // For Slide type
	EmailAddress string // For Email type
	EmailSubject string // For Email type (optional)
	FilePath     string // For File type
	ProgramPath  string // For Program type
}

// Hyperlink represents a clickable hyperlink on a shape or text run.
type Hyperlink struct {
	Action         HyperlinkAction
	Tooltip        string
	HighlightClick bool
}

// NewHyperlink creates a new hyperlink with the given action.
func NewHyperlink(action HyperlinkAction) Hyperlink {
	return Hyperlink{
		Action:         action,
		HighlightClick: true,
	}
}

// HyperlinkURL creates a URL hyperlink action.
func HyperlinkURL(urlStr string) HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionURL, URL: urlStr}
}

// HyperlinkSlide creates a slide hyperlink action.
func HyperlinkSlide(slideNum uint32) HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionSlide, SlideNumber: slideNum}
}

// HyperlinkFirstSlide creates a first-slide hyperlink action.
func HyperlinkFirstSlide() HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionFirstSlide}
}

// HyperlinkLastSlide creates a last-slide hyperlink action.
func HyperlinkLastSlide() HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionLastSlide}
}

// HyperlinkNextSlide creates a next-slide hyperlink action.
func HyperlinkNextSlide() HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionNextSlide}
}

// HyperlinkPreviousSlide creates a previous-slide hyperlink action.
func HyperlinkPreviousSlide() HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionPreviousSlide}
}

// HyperlinkEndShow creates an end-show hyperlink action.
func HyperlinkEndShow() HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionEndShow}
}

// HyperlinkEmail creates an email hyperlink action.
func HyperlinkEmail(address string) HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionEmail, EmailAddress: address}
}

// HyperlinkEmailWithSubject creates an email hyperlink action with subject.
func HyperlinkEmailWithSubject(address, subject string) HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionEmail, EmailAddress: address, EmailSubject: subject}
}

// HyperlinkFile creates a file hyperlink action.
func HyperlinkFile(path string) HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionFile, FilePath: path}
}

// HyperlinkProgram creates a program hyperlink action.
func HyperlinkProgram(path string) HyperlinkAction {
	return HyperlinkAction{Type: HyperlinkActionProgram, ProgramPath: path}
}

// RelationshipTarget returns the target URL for the relationship.
func (a HyperlinkAction) RelationshipTarget() string {
	switch a.Type {
	case HyperlinkActionURL:
		return a.URL
	case HyperlinkActionSlide:
		return fmt.Sprintf("slide%d.xml", a.SlideNumber)
	case HyperlinkActionFirstSlide:
		return "ppaction://hlinkshowjump?jump=firstslide"
	case HyperlinkActionLastSlide:
		return "ppaction://hlinkshowjump?jump=lastslide"
	case HyperlinkActionNextSlide:
		return "ppaction://hlinkshowjump?jump=nextslide"
	case HyperlinkActionPreviousSlide:
		return "ppaction://hlinkshowjump?jump=previousslide"
	case HyperlinkActionEndShow:
		return "ppaction://hlinkshowjump?jump=endshow"
	case HyperlinkActionEmail:
		mailto := "mailto:" + a.EmailAddress
		if a.EmailSubject != "" {
			mailto += "?subject=" + url.QueryEscape(a.EmailSubject)
		}
		return mailto
	case HyperlinkActionFile:
		return "file:///" + strings.ReplaceAll(a.FilePath, "\\", "/")
	case HyperlinkActionProgram:
		return "file:///" + strings.ReplaceAll(a.ProgramPath, "\\", "/")
	default:
		return ""
	}
}

// IsExternal returns true if the hyperlink points to an external resource.
func (a HyperlinkAction) IsExternal() bool {
	switch a.Type {
	case HyperlinkActionURL,
		HyperlinkActionFirstSlide,
		HyperlinkActionLastSlide,
		HyperlinkActionNextSlide,
		HyperlinkActionPreviousSlide,
		HyperlinkActionEndShow,
		HyperlinkActionEmail,
		HyperlinkActionFile,
		HyperlinkActionProgram:
		return true
	default:
		return false
	}
}

// ActionType returns the ppaction string for internal navigation links.
func (a HyperlinkAction) ActionType() string {
	switch a.Type {
	case HyperlinkActionProgram:
		return "ppaction://program"
	case HyperlinkActionSlide:
		return "ppaction://hlinksldjump"
	case HyperlinkActionFirstSlide:
		return "ppaction://hlinkshowjump?jump=firstslide"
	case HyperlinkActionLastSlide:
		return "ppaction://hlinkshowjump?jump=lastslide"
	case HyperlinkActionNextSlide:
		return "ppaction://hlinkshowjump?jump=nextslide"
	case HyperlinkActionPreviousSlide:
		return "ppaction://hlinkshowjump?jump=previousslide"
	case HyperlinkActionEndShow:
		return "ppaction://hlinkshowjump?jump=endshow"
	default:
		return ""
	}
}

func (h Hyperlink) WithTooltip(tooltip string) Hyperlink {
	h.Tooltip = tooltip
	return h
}

func (h Hyperlink) WithHighlightClick(highlight bool) Hyperlink {
	h.HighlightClick = highlight
	return h
}

// Validate checks for invalid hyperlink properties.
func (h Hyperlink) Validate() error {
	switch h.Action.Type {
	case HyperlinkActionURL:
		if h.Action.URL == "" {
			return errors.New("hyperlink URL cannot be empty")
		}
	case HyperlinkActionFile:
		if h.Action.FilePath == "" {
			return errors.New("hyperlink file path cannot be empty")
		}
		if strings.Contains(h.Action.FilePath, "..") {
			return errors.New("hyperlink file path cannot contain directory traversal (..)")
		}
	case HyperlinkActionProgram:
		if h.Action.ProgramPath == "" {
			return errors.New("hyperlink program path cannot be empty")
		}
		if strings.Contains(h.Action.ProgramPath, "..") {
			return errors.New("hyperlink program path cannot contain directory traversal (..)")
		}
	case HyperlinkActionEmail:
		if h.Action.EmailAddress == "" {
			return errors.New("hyperlink email address cannot be empty")
		}
	case HyperlinkActionSlide, HyperlinkActionFirstSlide, HyperlinkActionLastSlide:
		// Internal slide-jump actions require no additional payload validation.
	case HyperlinkActionNextSlide, HyperlinkActionPreviousSlide, HyperlinkActionEndShow:
		// Navigation actions are fully described by the action type itself.
	default:
		// Unknown action type is handled by render-time fallback behavior.
	}
	return nil
}
