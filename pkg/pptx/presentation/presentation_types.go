package presentation

import (
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/fonts"
	"github.com/djinn-soul/gopptx/pkg/pptx/handout"
	"github.com/djinn-soul/gopptx/pkg/pptx/printsettings"
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
	PrintSettings *printsettings.Settings
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

// convertShowSettings maps the public show settings onto the XML ones,
// resolving the named custom show to the id the XML refers to it by.
func convertShowSettings(s common.ShowSettings, shows []pptxxml.CustomShow) pptxxml.ShowSettings {
	out := pptxxml.ShowSettings{
		Loop:              s.Loop,
		Mode:              pptxxml.ShowMode(s.Mode),
		ShowScrollbar:     s.ShowScrollbar,
		DisableTimings:    s.DisableTimings,
		HideAnimation:     s.HideAnimation,
		DisableNarration:  s.DisableNarration,
		RangeKind:         pptxxml.SlideRangeKind(s.RangeKind),
		RangeStart:        s.RangeStart,
		RangeEnd:          s.RangeEnd,
		PenColor:          common.NormalizeHexColor(s.PenColor),
		LaserColor:        common.NormalizeHexColor(s.LaserColor),
		ShowMediaControls: s.ShowMediaControls,
	}
	if s.PenColor == "" {
		out.PenColor = ""
	}
	if s.LaserColor == "" {
		out.LaserColor = ""
	}
	if s.RangeKind == common.SlideRangeCustom {
		for _, show := range shows {
			if show.Name == s.CustomShowName {
				out.CustomShowID = show.ID
				break
			}
		}
	}
	return out
}

// convertCustomShows resolves each show's slide indices to the relationship ids
// presentation.xml.rels gives those slides.
func convertCustomShows(shows []common.CustomShow, slideCount, masterCount int) ([]pptxxml.CustomShow, error) {
	if len(shows) == 0 {
		return nil, nil
	}
	if masterCount < 1 {
		masterCount = 1
	}
	out := make([]pptxxml.CustomShow, 0, len(shows))
	for i, show := range shows {
		relIDs := make([]string, 0, len(show.SlideIndices))
		for _, idx := range show.SlideIndices {
			if idx < 0 || idx >= slideCount {
				return nil, fmt.Errorf(
					"custom show %q references slide index %d outside [0,%d)", show.Name, idx, slideCount)
			}
			// Slide relationships follow the masters and the theme.
			relIDs = append(relIDs, "rId"+strconv.Itoa(masterCount+1+idx+1))
		}
		out = append(out, pptxxml.CustomShow{Name: show.Name, ID: i, SlideRelIDs: relIDs})
	}
	return out, nil
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
