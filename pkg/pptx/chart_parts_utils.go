package pptx

import "github.com/vegito/goppt/internal/pptxxml"

func copyStringSlice(values []string) []string {
	out := make([]string, len(values))
	copy(out, values)
	return out
}

func copyFloat64Slice(values []float64) []float64 {
	out := make([]float64, len(values))
	copy(out, values)
	return out
}

func toXMLSeries(series []Series) []pptxxml.ChartSeries {
	out := make([]pptxxml.ChartSeries, len(series))
	for i := range series {
		out[i] = pptxxml.ChartSeries{
			Name:   series[i].Name,
			Values: copyFloat64Slice(series[i].Values),
		}
	}
	return out
}
