package shapes

import (
	"testing"
)

func TestResolveConnectorSiteIndices(t *testing.T) {
	rect := NewShape(ShapeTypeRectangle, 0, 0, 1000, 1000)
	shapes := []Shape{rect}

	tests := []struct {
		name      string
		connector Connector
		wantStart string
		wantEnd   string
	}{
		{
			name:      "auto site to center of shape 1 from a point to the right",
			connector: NewStraightConnector(2000, 500, 1000, 500).ConnectEndAuto(1),
			wantEnd:   ConnectionSiteRight,
		},
		{
			name:      "auto site to top of shape 1 from a point above",
			connector: NewStraightConnector(500, -1000, 500, 0).ConnectEndAuto(1),
			wantEnd:   ConnectionSiteTop,
		},
		{
			name:      "explicit site overrides auto",
			connector: NewStraightConnector(2000, 500, 1000, 500).ConnectEnd(1, ConnectionSiteLeft),
			wantEnd:   ConnectionSiteLeft,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startIdx, endIdx := ResolveConnectorSiteIndices(tt.connector, shapes)
			if tt.wantStart != "" {
				idx, _ := ConnectionSiteIndex(tt.wantStart)
				if startIdx == nil || *startIdx != idx {
					t.Errorf("start site index = %v, want %v", startIdx, idx)
				}
			}
			if tt.wantEnd != "" {
				idx, _ := ConnectionSiteIndex(tt.wantEnd)
				if endIdx == nil || *endIdx != idx {
					t.Errorf("end site index = %v, want %v", endIdx, idx)
				}
			}
		})
	}
}

func TestConnectorValidation(t *testing.T) {
	tests := []struct {
		name       string
		connector  Connector
		shapeCount int
		wantErr    bool
	}{
		{
			name:       "valid connector",
			connector:  NewStraightConnector(0, 0, 100, 100),
			shapeCount: 0,
			wantErr:    false,
		},
		{
			name:       "negative coordinate",
			connector:  NewStraightConnector(-1, 0, 100, 100),
			shapeCount: 0,
			wantErr:    true,
		},
		{
			name:       "same points",
			connector:  NewStraightConnector(100, 100, 100, 100),
			shapeCount: 0,
			wantErr:    true,
		},
		{
			name:       "invalid shape index",
			connector:  NewStraightConnector(0, 0, 100, 100).ConnectStartAuto(1),
			shapeCount: 0,
			wantErr:    true,
		},
		{
			name:       "valid shape index",
			connector:  NewStraightConnector(0, 0, 100, 100).ConnectStartAuto(1),
			shapeCount: 1,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.connector.Validate(tt.shapeCount, 1, 1); (err != nil) != tt.wantErr {
				t.Errorf("Connector.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
