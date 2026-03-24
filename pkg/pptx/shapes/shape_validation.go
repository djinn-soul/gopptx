package shapes

import (
	"fmt"

	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

// Validate checks for validity of shape parameters.
func (s Shape) Validate(slideIndex, shapeIndex int) error {
	if !s.IsDecorative && len(s.AltText) > common.MaxAltTextLength {
		return fmt.Errorf(
			"shape %d on slide %d alt text exceeds %d characters",
			shapeIndex,
			slideIndex,
			common.MaxAltTextLength,
		)
	}

	if err := s.validateShapeBounds(slideIndex, shapeIndex); err != nil {
		return err
	}
	if !IsShapeType(s.Type) {
		return fmt.Errorf("shape %d type %q is invalid on slide %d", shapeIndex, s.Type, slideIndex)
	}

	if err := s.validateFills(slideIndex, shapeIndex); err != nil {
		return err
	}
	if err := s.validateLinesAndRotation(slideIndex, shapeIndex); err != nil {
		return err
	}
	return s.validateActions(slideIndex, shapeIndex)
}

func (s Shape) validateShapeBounds(slideIndex, shapeIndex int) error {
	if s.X < 0 || s.Y < 0 {
		return fmt.Errorf("shape %d on slide %d position cannot be negative", shapeIndex, slideIndex)
	}
	if s.CX <= 0 || s.CY <= 0 {
		return fmt.Errorf("shape %d on slide %d size must be > 0", shapeIndex, slideIndex)
	}
	return nil
}

func (s Shape) validateFills(slideIndex, shapeIndex int) error {
	// Check for conflicts between legacy and rich fill
	if s.RichFill != nil && (s.Fill != nil || s.GradientFill != nil) {
		return fmt.Errorf("shape %d (type %q) on slide %d cannot set both rich fill and legacy fill",
			shapeIndex, s.Type, slideIndex)
	}
	if s.Fill != nil && s.GradientFill != nil {
		return fmt.Errorf("shape %d (type %q) on slide %d cannot set both solid and gradient fill",
			shapeIndex, s.Type, slideIndex)
	}
	if s.Fill != nil {
		if err := s.Fill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid fill: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.GradientFill != nil {
		if err := s.GradientFill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid gradient fill: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RichFill != nil {
		if err := s.RichFill.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid rich fill: %w", shapeIndex, slideIndex, err)
		}
	}
	return nil
}

func (s Shape) validateLinesAndRotation(slideIndex, shapeIndex int) error {
	// Check for conflicts between legacy and rich line
	if s.RichLine != nil && s.Line != nil {
		return fmt.Errorf("shape %d (type %q) on slide %d cannot set both rich line and legacy line",
			shapeIndex, s.Type, slideIndex)
	}
	if s.Line != nil {
		if err := s.Line.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid line: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RichLine != nil {
		if err := s.RichLine.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid rich line: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RichShadow != nil {
		if err := s.RichShadow.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid rich shadow: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.RotationDeg != nil {
		if *s.RotationDeg < -360 || *s.RotationDeg > 360 {
			return fmt.Errorf("shape %d on slide %d rotation must be in [-360,360]", shapeIndex, slideIndex)
		}
	}
	return nil
}

func (s Shape) validateActions(slideIndex, shapeIndex int) error {
	if s.ClickAction != nil {
		if err := s.ClickAction.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid click action: %w", shapeIndex, slideIndex, err)
		}
	} else if s.Hyperlink != nil {
		if err := s.Hyperlink.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid hyperlink: %w", shapeIndex, slideIndex, err)
		}
	}
	if s.HoverAction != nil {
		if err := s.HoverAction.Validate(); err != nil {
			return fmt.Errorf("shape %d on slide %d has invalid hover action: %w", shapeIndex, slideIndex, err)
		}
	}
	return nil
}
