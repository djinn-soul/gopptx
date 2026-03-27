package editorcommon

import "github.com/djinn-soul/gopptx/pkg/pptx/common"

// SlideSize describes the dimensions of slides in a presentation in EMUs.
type SlideSize = common.SlideSize

// SlideSize4x3 returns the standard 4:3 slide size (10x7.5 inches).
func SlideSize4x3() SlideSize {
	return common.GetSlideSize4x3()
}

// SlideSize16x9 returns the standard 16:9 widescreen slide size (13.33x7.5 inches).
func SlideSize16x9() SlideSize {
	return common.GetSlideSize16x9()
}

// Metadata describes summary information for a PPTX package.
type Metadata struct {
	common.Metadata

	VBA any // *vba.VBAProject at runtime, loosely typed to avoid import cycles in common
}

// SlideMetadata describes one slide entry inside an editable presentation.
type SlideMetadata struct {
	Index          int
	SlideID        int64
	RelationshipID string
	PartName       string
	Title          string
}

// CoreProperties describes the docProps/core.xml metadata.
type CoreProperties = common.CoreProperties

// ShowSettings controls how a presentation is shown (maps to p:showPr in presentation.xml).
type ShowSettings = common.ShowSettings

// ShowMode defines the slide show presentation mode.
type ShowMode = common.ShowMode

const (
	ShowModePresent = common.ShowModePresent
	ShowModeBrowse  = common.ShowModeBrowse
	ShowModeKiosk   = common.ShowModeKiosk
)

// CustomXMLPart describes a Custom XML part exposed to the editor.
type CustomXMLPart = common.CustomXMLPart

// CustomXMLKV describes a key-value property for a CustomXMLPart.
type CustomXMLKV = common.CustomXMLKV
