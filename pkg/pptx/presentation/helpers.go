package presentation

import (
	"fmt"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/media"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

const masterImageRelIDStart = 8

func mapThemeToSpec(theme *styling.Theme) *pptxxml.ThemeSpec {
	if theme == nil {
		return nil
	}
	spec := &pptxxml.ThemeSpec{
		Name: theme.Name,
		Colors: pptxxml.ColorSchemeSpec{
			Name:     theme.Colors.Name,
			Dk1:      theme.Colors.Dk1,
			Lt1:      theme.Colors.Lt1,
			Dk2:      theme.Colors.Dk2,
			Lt2:      theme.Colors.Lt2,
			Accent1:  theme.Colors.Accent1,
			Accent2:  theme.Colors.Accent2,
			Accent3:  theme.Colors.Accent3,
			Accent4:  theme.Colors.Accent4,
			Accent5:  theme.Colors.Accent5,
			Accent6:  theme.Colors.Accent6,
			Hlink:    theme.Colors.Hlink,
			FolHlink: theme.Colors.FolHlink,
		},
		Fonts: pptxxml.FontSchemeSpec{
			Name:      theme.Fonts.Name,
			MajorFont: theme.Fonts.MajorFont,
			MinorFont: theme.Fonts.MinorFont,
		},
	}
	return spec
}

func mapMasterToSpec(master *elements.SlideMaster, imageRefs []pptxxml.ImageRef) *pptxxml.SlideMasterSpec {
	if master == nil {
		return nil
	}
	spec := &pptxxml.SlideMasterSpec{
		FooterText: master.FooterText,
		Images:     imageRefs,
		TxStyles:   elements.MapTxStyles(master.TxStyles),
	}
	if master.ColorMapping != nil {
		spec.ColorMapping = &pptxxml.ColorMappingSpec{
			BG1: master.ColorMapping.BG1,
			TX1: master.ColorMapping.TX1,
		}
	}
	if master.Background != nil {
		spec.Background = elements.ToXMLBackgroundSpec(master.Background, "")
	}
	// Map shapes (no hyperlinks on master shapes).
	masterShapes := make([]shapes.Shape, 0, len(master.Shapes))
	for _, sd := range master.Shapes {
		masterShapes = append(masterShapes, sd.ToShape())
	}
	spec.Shapes = shapes.ToXMLShapeSpecs(masterShapes, nil)
	return spec
}

// buildMasterImageInfo registers master images and returns relationship targets and ImageRef specs.
func buildMasterImageInfo(master *elements.SlideMaster, catalog *media.Catalog) ([]string, []pptxxml.ImageRef) {
	if master == nil || len(master.Images) == 0 {
		return nil, nil
	}
	targets := make([]string, 0, len(master.Images))
	refs := make([]pptxxml.ImageRef, 0, len(master.Images))
	for i, img := range master.Images {
		mediaName, err := catalog.RegisterImage(img)
		if err != nil {
			continue // skip unresolved master images
		}
		// Master image RIDs start at rId8 (rId1-6 are layouts, rId7 is theme).
		relID := fmt.Sprintf("rId%d", masterImageRelIDStart+i)
		targets = append(targets, fmt.Sprintf("../media/%s", mediaName))
		refs = append(refs, pptxxml.ImageRef{
			RelID: relID,
			Name:  fmt.Sprintf("Master Picture %d", i+1),
			X:     img.X.Emu(),
			Y:     img.Y.Emu(),
			CX:    img.CX.Emu(),
			CY:    img.CY.Emu(),

			AltText:      img.AltText,
			IsDecorative: img.IsDecorative,
		})
	}
	return targets, refs
}

func writeMediaFiles(pw *pptxxml.PackageWriter, catalog *media.Catalog) error {
	for _, asset := range catalog.Assets() {
		path := fmt.Sprintf("ppt/media/%s", asset.MediaName())
		pw.AddBinaryPart(path, asset.Data())
	}
	return nil
}
