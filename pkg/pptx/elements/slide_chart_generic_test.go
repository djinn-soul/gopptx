package elements

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/charts"
)

func newPie() charts.PieChart {
	return charts.NewPieChart([]string{"A", "B"}, []float64{1, 2})
}

// WithChart accepts a chart value and a pointer to one, and both must land in
// the same slide field.
func TestWithChartAcceptsValueAndPointer(t *testing.T) {
	pie := newPie()

	fromValue := NewSlide("Chart").WithChart(pie)
	fromPointer := NewSlide("Chart").WithChart(&pie)

	if fromValue.Pie == nil {
		t.Fatal("passing a chart value did not set the Pie field")
	}
	if fromPointer.Pie == nil {
		t.Fatal("passing a chart pointer did not set the Pie field")
	}
	if len(fromValue.Pie.Categories) != len(fromPointer.Pie.Categories) {
		t.Errorf("value and pointer produced different charts: %+v vs %+v",
			*fromValue.Pie, *fromPointer.Pie)
	}
}

// A nil interface and a typed nil pointer must both be no-ops rather than
// panicking or clearing an existing chart.
func TestWithChartIgnoresNilChartAndNilPointer(t *testing.T) {
	bar := charts.NewBarChart([]string{"A"}, []float64{1})
	base := NewSlide("Chart").WithBarChart(bar)

	if got := base.WithChart(nil); got.Chart == nil {
		t.Error("WithChart(nil) cleared the existing chart")
	}

	var nilPie *charts.PieChart
	got := base.WithChart(nilPie)
	if got.Chart == nil {
		t.Error("WithChart(typed nil pointer) cleared the existing chart")
	}
	if got.Pie != nil {
		t.Error("WithChart(typed nil pointer) set the Pie field")
	}
}

func TestDerefChart(t *testing.T) {
	pie := newPie()
	var nilPie *charts.PieChart

	t.Run("value passes through", func(t *testing.T) {
		got, ok := derefChart(pie)
		if !ok {
			t.Fatal("a chart value must be usable")
		}
		if _, isPie := got.(charts.PieChart); !isPie {
			t.Errorf("got %T, want charts.PieChart", got)
		}
	})

	t.Run("pointer is dereferenced", func(t *testing.T) {
		got, ok := derefChart(&pie)
		if !ok {
			t.Fatal("a chart pointer must be usable")
		}
		if _, isPie := got.(charts.PieChart); !isPie {
			t.Errorf("got %T, want charts.PieChart", got)
		}
	})

	t.Run("nil interface is rejected", func(t *testing.T) {
		if _, ok := derefChart(nil); ok {
			t.Error("a nil chart must not be usable")
		}
	})

	t.Run("nil pointer is rejected", func(t *testing.T) {
		if _, ok := derefChart(nilPie); ok {
			t.Error("a nil chart pointer must not be usable")
		}
	})
}
