package editor

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
	editormodchart "github.com/djinn-soul/gopptx/pkg/pptx/editor/modules/chart"
)

// Chart data-source kinds reported by GetChartDataSource.
const (
	ChartDataSourceEmbedded = "embedded"
	ChartDataSourceExternal = "external"
	ChartDataSourceNone     = "none"
)

var (
	externalDataPattern = regexp.MustCompile(
		`(?s)<c:externalData\b[^>]*>.*?</c:externalData>|<c:externalData\b[^>]*/>`,
	)
	externalDataRelID     = regexp.MustCompile(`\br:id="([^"]*)"`)
	externalAutoUpdatePat = regexp.MustCompile(`<c:autoUpdate\b[^>]*\bval="([^"]*)"`)
)

// UpdateChartCachedValues refreshes the numbers a chart displays without
// touching its workbook link.
//
// UpdateChartData regenerates an embedded workbook and repoints the chart at
// it, which turns a chart linked to an external workbook into an embedded one.
// A deck rebuilt weekly from a linked workbook wants the opposite: refresh the
// cache PowerPoint draws from, and leave the link alone, so the numbers are
// current without opening each chart and clicking Refresh (upstream #115).
func (e *PresentationEditor) UpdateChartCachedValues(
	slideIndex int,
	selector common.ChartSelector,
	req common.ChartDataUpdate,
) error {
	chartRef, err := e.resolveChartRef(slideIndex, selector)
	if err != nil {
		return err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}

	kind := editormodchart.DetectChartKind(chartXML)
	if validateErr := editormodchart.ValidateChartUpdatePayload(kind, req); validateErr != nil {
		return validateErr
	}

	patched, err := editormodchart.PatchChartDataCache(
		chartXML, kind, req,
		editormodchart.CachePatchOptions{KeepFormulas: true},
	)
	if err != nil {
		return err
	}
	e.parts.Set(chartRef.ChartPart, patched)
	return nil
}

// GetChartDataSource reports where a chart's data comes from, so a caller can
// tell a linked chart from an embedded one before deciding how to update it.
func (e *PresentationEditor) GetChartDataSource(
	slideIndex int,
	selector common.ChartSelector,
) (*common.ChartDataSource, error) {
	chartRef, err := e.resolveChartRef(slideIndex, selector)
	if err != nil {
		return nil, err
	}

	chartXML, ok := e.parts.Get(chartRef.ChartPart)
	if !ok {
		return nil, fmt.Errorf("chart part %s not found", chartRef.ChartPart)
	}

	source := &common.ChartDataSource{
		ChartPart: chartRef.ChartPart,
		Kind:      ChartDataSourceNone,
	}
	block := externalDataPattern.FindString(string(chartXML))
	if block == "" {
		// A chart written without <c:externalData> still names its workbook
		// through a package relationship, which is how gopptx itself embeds one.
		return e.describeChartPackageRelationship(chartRef.ChartPart, source)
	}
	if match := externalAutoUpdatePat.FindStringSubmatch(block); len(match) > 1 {
		autoUpdate := match[1] == "1" || strings.EqualFold(match[1], "true")
		source.AutoUpdate = &autoUpdate
	}

	relMatch := externalDataRelID.FindStringSubmatch(block)
	if len(relMatch) < 2 {
		return source, nil
	}
	source.RelID = relMatch[1]
	return e.describeChartDataRelationship(chartRef.ChartPart, source)
}

// describeChartPackageRelationship finds the workbook a chart carries when the
// chart XML declares no <c:externalData> element.
func (e *PresentationEditor) describeChartPackageRelationship(
	chartPart string,
	source *common.ChartDataSource,
) (*common.ChartDataSource, error) {
	relsData, ok := e.parts.Get(common.RelsPathFor(chartPart))
	if !ok {
		return source, nil
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil, fmt.Errorf("parse chart rels: %w", err)
	}
	for _, rel := range rels {
		if rel.Type != common.RelTypePackage {
			continue
		}
		source.RelID = rel.ID
		source.Target = rel.Target
		if strings.EqualFold(rel.TargetMode, "External") {
			source.Kind = ChartDataSourceExternal
			return source, nil
		}
		source.Kind = ChartDataSourceEmbedded
		source.PartPath = common.ResolveRelationshipTarget(chartPart, rel.Target)
		return source, nil
	}
	return source, nil
}

func (e *PresentationEditor) describeChartDataRelationship(
	chartPart string,
	source *common.ChartDataSource,
) (*common.ChartDataSource, error) {
	relsData, ok := e.parts.Get(common.RelsPathFor(chartPart))
	if !ok {
		return source, nil
	}
	rels, err := parseRelationshipsXML(relsData)
	if err != nil {
		return nil, fmt.Errorf("parse chart rels: %w", err)
	}
	for _, rel := range rels {
		if rel.ID != source.RelID {
			continue
		}
		source.Target = rel.Target
		if strings.EqualFold(rel.TargetMode, "External") {
			source.Kind = ChartDataSourceExternal
			return source, nil
		}
		source.Kind = ChartDataSourceEmbedded
		source.PartPath = common.ResolveRelationshipTarget(chartPart, rel.Target)
		return source, nil
	}
	return source, nil
}

func (e *PresentationEditor) resolveChartRef(
	slideIndex int,
	selector common.ChartSelector,
) (common.SlideChartRef, error) {
	if e == nil || e.parts == nil {
		return common.SlideChartRef{}, errors.New("editor cannot be nil")
	}
	refs, err := e.ListSlideCharts(slideIndex)
	if err != nil {
		return common.SlideChartRef{}, err
	}
	return editormodchart.ResolveChartSelector(refs, selector, slideIndex)
}
