package editor

import (
	"fmt"

	common "github.com/djinn-soul/gopptx/pkg/pptx/editor/common"
)

// validateShapeExtents rejects the extents OOXML forbids. cx/cy are typed
// ST_PositiveCoordinate, so a negative extent is schema-invalid and a zero
// extent is degenerate. Negative x/y is deliberately allowed: that is how a
// shape is placed partly off the slide.
//
// The builder path has always enforced this; validating here too means the
// error arrives at the call that caused it rather than at a later Validate().
func validateShapeExtents(w, h float64) error {
	if w <= 0 || h <= 0 {
		return fmt.Errorf("shape size must be > 0, got width=%v height=%v", w, h)
	}
	return nil
}

// validateShapeUpdateExtents applies the same rule to a partial update, where
// only the supplied dimension is checked.
func validateShapeUpdateExtents(updates common.ShapeUpdate) error {
	for _, dimension := range []struct {
		name  string
		value *int
	}{{"w", updates.W}, {"h", updates.H}} {
		if dimension.value != nil && *dimension.value <= 0 {
			return fmt.Errorf("shape %s must be > 0, got %d", dimension.name, *dimension.value)
		}
	}
	return nil
}
