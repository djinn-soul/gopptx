package chart

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

func validateAxisScaleAgainstXML(xml string, req common.ChartFormatUpdate) error {
	if err := validateRetainedAxisScales(
		"category axis",
		buildAxisStates(xml, []string{"catAx", "dateAx"}),
		req.CategoryAxisMinimumScale,
		req.CategoryAxisMaximumScale,
	); err != nil {
		return err
	}
	return validateRetainedAxisScales(
		"value axis",
		buildAxisStates(xml, []string{"valAx"}),
		req.ValueAxisMinimumScale,
		req.ValueAxisMaximumScale,
	)
}

func validateRetainedAxisScales(
	name string,
	states []common.ChartAxisState,
	requestedMinimum *float64,
	requestedMaximum *float64,
) error {
	if requestedMinimum == nil && requestedMaximum == nil {
		return nil
	}
	for _, state := range states {
		minimum, maximum := state.MinimumScale, state.MaximumScale
		if requestedMinimum != nil {
			minimum = requestedMinimum
		}
		if requestedMaximum != nil {
			maximum = requestedMaximum
		}
		if minimum != nil && maximum != nil && *minimum >= *maximum {
			return fmt.Errorf("%s minimum scale must be less than maximum scale", name)
		}
	}
	return nil
}
