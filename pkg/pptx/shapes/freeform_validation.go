package shapes

import (
	"errors"
	"fmt"
)

// Validate checks the freeform for validity.
func (f Freeform) Validate() error {
	if len(f.Points) < minFreeformPoints {
		return errors.New("freeform requires at least 2 points")
	}
	if err := f.validateFillCompatibility(); err != nil {
		return err
	}
	if err := f.validateRichFill(); err != nil {
		return err
	}
	if err := f.validateLegacyFill(); err != nil {
		return err
	}
	if err := f.validateGradientFill(); err != nil {
		return err
	}
	if err := f.validateRichLine(); err != nil {
		return err
	}
	if err := f.validateLegacyLine(); err != nil {
		return err
	}
	return f.validateRichShadow()
}

func (f Freeform) validateFillCompatibility() error {
	if f.RichFill != nil && (f.Fill != nil || f.GradientFill != nil) {
		return errors.New("cannot set both rich fill and legacy fill")
	}
	if f.Fill != nil && f.GradientFill != nil {
		return errors.New("cannot set both solid and gradient fill")
	}
	return nil
}

func (f Freeform) validateRichFill() error {
	if f.RichFill == nil {
		return nil
	}
	if err := f.RichFill.Validate(); err != nil {
		return fmt.Errorf("invalid rich fill: %w", err)
	}
	return nil
}

func (f Freeform) validateLegacyFill() error {
	if f.Fill == nil {
		return nil
	}
	if err := f.Fill.Validate(); err != nil {
		return fmt.Errorf("invalid fill: %w", err)
	}
	return nil
}

func (f Freeform) validateGradientFill() error {
	if f.GradientFill == nil {
		return nil
	}
	if err := f.GradientFill.Validate(); err != nil {
		return fmt.Errorf("invalid gradient fill: %w", err)
	}
	return nil
}

func (f Freeform) validateRichLine() error {
	if f.RichLine == nil {
		return nil
	}
	if err := f.RichLine.Validate(); err != nil {
		return fmt.Errorf("invalid rich line: %w", err)
	}
	return nil
}

func (f Freeform) validateLegacyLine() error {
	if f.Line == nil {
		return nil
	}
	if err := f.Line.Validate(); err != nil {
		return fmt.Errorf("invalid line: %w", err)
	}
	return nil
}

func (f Freeform) validateRichShadow() error {
	if f.RichShadow == nil {
		return nil
	}
	if err := f.RichShadow.Validate(); err != nil {
		return fmt.Errorf("invalid rich shadow: %w", err)
	}
	return nil
}
