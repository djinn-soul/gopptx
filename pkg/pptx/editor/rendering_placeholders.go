package editor

import (
	"fmt"
	"path"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	editorslide "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/slide"
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

// renderEditorPlaceholderSpecs converts SlideContent.PlaceholderOverrides into
// XML specs for the editor rendering path. It returns the specs, any additional
// image relationship targets, chart rels, and an error.
//
//nolint:gocognit // Placeholder rendering must keep explicit branching for supported override kinds.
func renderEditorPlaceholderSpecs(
	e *PresentationEditor,
	slide elements.SlideContent,
	slidePart string,
	slideNumber int,
	startRID int,
) ([]pptxxml.PlaceholderOverrideSpec, []string, []pptxxml.ChartRel, error) {
	if len(slide.PlaceholderOverrides) == 0 {
		return nil, nil, nil, nil
	}

	specs := make([]pptxxml.PlaceholderOverrideSpec, 0, len(slide.PlaceholderOverrides))
	var imageTargets []string
	var chartRels []pptxxml.ChartRel
	currentRID := startRID
	placeholders, err := editorslide.LookupSlidePlaceholders(
		slidePart,
		e.parts.Get,
		func(content []byte) []editorslide.PlaceholderMeta {
			parsed := parsePlaceholdersFromSlideXML(content)
			out := make([]editorslide.PlaceholderMeta, 0, len(parsed))
			for _, ph := range parsed {
				out = append(out, editorslide.PlaceholderMeta{
					Name:  ph.Name,
					Type:  ph.Type,
					Index: ph.Index,
				})
			}
			return out
		},
	)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, override := range slide.PlaceholderOverrides {
		targetType, targetIndex, err := editorslide.ResolvePlaceholderTarget(override, placeholders)
		if err != nil {
			return nil, nil, nil, err
		}
		spec := pptxxml.PlaceholderOverrideSpec{
			Index: targetIndex,
			Type:  targetType,
			Text:  override.Text,
		}

		if override.Override != nil {
			spec.X = editorslide.MapOptionalLength(override.Override.X)
			spec.Y = editorslide.MapOptionalLength(override.Override.Y)
			spec.CX = editorslide.MapOptionalLength(override.Override.CX)
			spec.CY = editorslide.MapOptionalLength(override.Override.CY)
			spec.TextStyle = editorslide.MapPlaceholderTextStyle(override.Override.TextStyle)
			spec.ForceRectGeometry = override.Override.ForceRect
		}

		if override.Image != nil {
			imageRef, imageTarget, imageErr := e.renderPlaceholderImageRef(override, currentRID)
			if imageErr != nil {
				return nil, nil, nil, imageErr
			}
			if imageRef != nil {
				spec.Image = imageRef
				imageTargets = append(imageTargets, imageTarget)
				currentRID++
			}
		}

		if override.Table != nil {
			tableSpec, err := override.Table.ToTableSpec(slideNumber)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("placeholder table %d: %w", targetIndex, err)
			}
			spec.Table = tableSpec
		}

		if override.Chart != nil {
			chartSpec := override.Chart.ToChartSpec()
			chartPath := fmt.Sprintf("ppt/charts/chart_ph_%d_%d.xml", slideNumber, targetIndex)
			e.parts.Set(chartPath, []byte(pptxxml.ChartPartXML(chartSpec)))

			rid := fmt.Sprintf("rId%d", currentRID)
			currentRID++
			spec.Chart = &pptxxml.ChartFrame{
				RelID: rid,
				X:     chartSpec.X,
				Y:     chartSpec.Y,
				CX:    chartSpec.CX,
				CY:    chartSpec.CY,
			}
			chartRels = append(chartRels, pptxxml.ChartRel{
				RID:    rid,
				Target: "../charts/" + path.Base(chartPath),
			})
		}

		specs = append(specs, spec)
	}

	return specs, imageTargets, chartRels, nil
}

func (e *PresentationEditor) renderPlaceholderImageRef(
	override shapes.PlaceholderContent,
	ridIndex int,
) (*pptxxml.ImageRef, string, error) {
	return editorslide.RenderPlaceholderImageRef(override, ridIndex, e.registerEditorImage)
}
