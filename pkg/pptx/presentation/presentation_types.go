package presentation

import (
	"crypto/rand"
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/fonts"
	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
	"github.com/djinn-soul/gopptx/pkg/pptx/vba"
)

const (
	minMasterCountWithNativeNotesTheme = 2
	protectionSaltBytes                = 16
	protectionHashAlgSIDSHA512         = 14
	guidRandomBytes                    = 16
	guidVersionMask                    = 0x0f
	guidVersionNibble                  = 0x40
	guidVariantMask                    = 0x3f
	guidVariantNibble                  = 0x80
	maxAuthorInitialRunes              = 2
	authorColorPaletteSize             = 10
	customXMLRelationshipPairCount     = 2
)

type Metadata struct {
	common.Metadata

	Theme         *styling.Theme
	Master        *elements.SlideMaster
	Masters       []*elements.SlideMaster
	NotesMaster   *elements.NotesMaster
	HandoutMaster *handout.HandoutMaster
	Sections      []Section
	RTL           bool
	VBA           *vba.VBAProject
	EmbeddedFonts []fonts.EmbeddedFont
}

type Section struct {
	Name         string
	SlideIndices []int // 0-based indices of slides in this section
}

type SlideSize = common.SlideSize

func GetSlideSize4x3() SlideSize {
	return common.GetSlideSize4x3()
}

func GetSlideSize16x9() SlideSize {
	return common.GetSlideSize16x9()
}

func convertShowSettings(s common.ShowSettings) *pptxxml.ShowSettings {
	if !s.Loop && s.Mode == common.ShowModePresent && !s.DisableTimings && !s.HideAnimation {
		return nil
	}
	return &pptxxml.ShowSettings{
		Loop:           s.Loop,
		Mode:           pptxxml.ShowMode(s.Mode),
		ShowScrollbar:  s.ShowScrollbar,
		DisableTimings: s.DisableTimings,
		HideAnimation:  s.HideAnimation,
	}
}

func generateGUID() (string, error) {
	b := make([]byte, guidRandomBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate random bytes for GUID: %w", err)
	}
	b[6] = (b[6] & guidVersionMask) | guidVersionNibble
	b[8] = (b[8] & guidVariantMask) | guidVariantNibble
	return fmt.Sprintf("{%08X-%04X-%04X-%04X-%012X}", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
