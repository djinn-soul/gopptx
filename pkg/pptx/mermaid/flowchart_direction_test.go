package mermaid

import "testing"

// nodeOrigins maps each node label to where it was placed.
func nodeOrigins(t *testing.T, code string) map[string][2]int64 {
	t.Helper()
	elements := renderFlowchart(code, GetTheme("default"))
	origins := make(map[string][2]int64, len(elements.Shapes))
	for _, s := range elements.Shapes {
		if s.Text == "" {
			continue
		}
		origins[s.Text] = [2]int64{int64(s.X), int64(s.Y)}
	}
	return origins
}

// The parsed Direction reached the connector router but never the node
// placement, so `flowchart TD` came out laid left-to-right like `flowchart LR`.
func TestFlowchartDirectionDrivesNodePlacement(t *testing.T) {
	tests := []struct {
		name   string
		header string
		// check receives the origins of A (first rank) and C (last rank).
		check func(t *testing.T, a, c [2]int64)
	}{
		{
			name:   "LR ranks run left to right",
			header: "flowchart LR",
			check: func(t *testing.T, a, c [2]int64) {
				t.Helper()
				if c[0] <= a[0] {
					t.Errorf("C x=%d should be right of A x=%d", c[0], a[0])
				}
				if c[1] != a[1] {
					t.Errorf("LR should keep both on one row, got y %d and %d", a[1], c[1])
				}
			},
		},
		{
			name:   "TD ranks run top to bottom",
			header: "flowchart TD",
			check: func(t *testing.T, a, c [2]int64) {
				t.Helper()
				if c[1] <= a[1] {
					t.Errorf("C y=%d should be below A y=%d", c[1], a[1])
				}
				if c[0] != a[0] {
					t.Errorf("TD should keep both in one column, got x %d and %d", a[0], c[0])
				}
			},
		},
		{
			name:   "RL ranks run right to left",
			header: "flowchart RL",
			check: func(t *testing.T, a, c [2]int64) {
				t.Helper()
				if c[0] >= a[0] {
					t.Errorf("C x=%d should be left of A x=%d", c[0], a[0])
				}
			},
		},
		{
			name:   "BT ranks run bottom to top",
			header: "flowchart BT",
			check: func(t *testing.T, a, c [2]int64) {
				t.Helper()
				if c[1] >= a[1] {
					t.Errorf("C y=%d should be above A y=%d", c[1], a[1])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origins := nodeOrigins(t, tt.header+"\n    A[Start] --> B[Middle]\n    B --> C[End]")
			a, okA := origins["Start"]
			c, okC := origins["End"]
			if !okA || !okC {
				t.Fatalf("missing nodes, got %v", origins)
			}
			tt.check(t, a, c)
		})
	}
}

// TB is the Mermaid default and must behave exactly like TD.
func TestFlowchartDefaultsToTopBottom(t *testing.T) {
	origins := nodeOrigins(t, "flowchart TB\n    A[Start] --> B[End]")
	a, b := origins["Start"], origins["End"]
	if b[1] <= a[1] {
		t.Errorf("TB should stack downward, A y=%d B y=%d", a[1], b[1])
	}
}
