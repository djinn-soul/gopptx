package common

import "github.com/djinn-soul/gopptx/pkg/pptx/common"

// SlideSize describes the dimensions of slides in a presentation in EMUs.
type SlideSize = common.SlideSize

var (
	// SlideSize4x3 is the standard 4:3 slide size (10x7.5 inches).
	SlideSize4x3 = common.SlideSize4x3
	// SlideSize16x9 is the standard 16:9 widescreen slide size (13.33x7.5 inches).
	SlideSize16x9 = common.SlideSize16x9
)

// PresentationMetadata describes summary information for a PPTX package.
type PresentationMetadata = common.PresentationMetadata

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
