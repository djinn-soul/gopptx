package chart

import "testing"

func TestNormalizeDataLabelPosition(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		isBarChart bool
		expected   string
	}{
		{"Right on Bar Chart (Office 365 Web Fix)", "RIGHT", true, "outEnd"},
		{"r on Bar Chart (Office 365 Web Fix)", "r", true, "outEnd"},
		{"Left on Bar Chart", "LEFT", true, "inBase"},
		{"Right on Line/Scatter Chart", "RIGHT", false, "r"},
		{"Outside End", "outside_end", false, "outEnd"},
		{"Inside End", "inside_end", true, "inEnd"},
		{"Center", "center", false, "ctr"},
		{"Best Fit", "best_fit", false, "bestFit"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeDataLabelPosition(tc.raw, tc.isBarChart)
			if got != tc.expected {
				t.Errorf("normalizeDataLabelPosition(%q, %v) = %q; want %q", tc.raw, tc.isBarChart, got, tc.expected)
			}
		})
	}
}
