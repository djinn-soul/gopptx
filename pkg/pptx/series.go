package pptx

import (
	"fmt"
	"math"
	"strings"
)

// Series represents one named numeric series.
type Series struct {
	Name   string
	Values []float64
}

func copySeriesList(series []Series) []Series {
	out := make([]Series, len(series))
	for i := range series {
		vals := make([]float64, len(series[i].Values))
		copy(vals, series[i].Values)
		out[i] = Series{
			Name:   series[i].Name,
			Values: vals,
		}
	}
	return out
}

func validateSeriesList(series []Series, categoriesLen int, slideIndex int, label string) error {
	if len(series) == 0 {
		return fmt.Errorf("slide %d %s requires at least one series", slideIndex, label)
	}
	for i := range series {
		if strings.TrimSpace(series[i].Name) == "" {
			return fmt.Errorf("slide %d %s series %d name cannot be empty", slideIndex, label, i+1)
		}
		if len(series[i].Values) != categoriesLen {
			return fmt.Errorf(
				"slide %d %s series %d length mismatch (%d vs %d)",
				slideIndex,
				label,
				i+1,
				len(series[i].Values),
				categoriesLen,
			)
		}
		for j := range series[i].Values {
			if math.IsNaN(series[i].Values[j]) || math.IsInf(series[i].Values[j], 0) {
				return fmt.Errorf(
					"slide %d %s series %d value %d must be finite",
					slideIndex,
					label,
					i+1,
					j+1,
				)
			}
		}
	}
	return nil
}
