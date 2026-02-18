package elements

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/animations"
	"github.com/djinn-soul/gopptx/pkg/pptx/common"
)

func validateSlideContent(s SlideContent, index int) error {
	if err := validateSlideObjects(s, index); err != nil {
		return err
	}
	if err := validateSlideCharts(s, index); err != nil {
		return err
	}
	if err := validateSlideSmartArt(s, index); err != nil {
		return err
	}
	if err := validateSlideTextStyles(s); err != nil {
		return err
	}
	if err := validateSlideTypography(s); err != nil {
		return err
	}
	if err := validateSlideTransitionAndTable(s, index); err != nil {
		return err
	}
	if err := validateSlideBulletAndAlignment(s); err != nil {
		return err
	}
	if s.Title == "" && s.Layout != SlideLayoutBlank {
		return errors.New("title cannot be empty")
	}
	return validateSlideAnimations(s)
}

func validateSlideObjects(s SlideContent, index int) error {
	for shapeIndex, shape := range s.Shapes {
		if err := shape.Validate(index, shapeIndex+1); err != nil {
			return err
		}
	}
	for connectorIndex, connector := range s.Connectors {
		if err := connector.ValidateWithShapes(s.Shapes, index, connectorIndex+1); err != nil {
			return err
		}
	}
	for imageIndex, image := range s.Images {
		if err := image.Validate(index, imageIndex+1); err != nil {
			return err
		}
	}
	return nil
}

func validateSlideCharts(s SlideContent, index int) error {
	for _, chart := range collectSlideCharts(s) {
		if err := chart.Validate(index); err != nil {
			return err
		}
	}
	return nil
}

func collectSlideCharts(s SlideContent) []ChartDefinition {
	const expectedChartTypeCount = 19
	// NOTE: Capacity derived from constant expectedChartTypeCount.
	chartsOnSlide := make([]ChartDefinition, 0, expectedChartTypeCount)
	if s.Chart != nil {
		chartsOnSlide = append(chartsOnSlide, s.Chart)
	}
	if s.BarHorizontal != nil {
		chartsOnSlide = append(chartsOnSlide, s.BarHorizontal)
	}
	if s.BarStacked != nil {
		chartsOnSlide = append(chartsOnSlide, s.BarStacked)
	}
	if s.BarStacked100 != nil {
		chartsOnSlide = append(chartsOnSlide, s.BarStacked100)
	}
	if s.Line != nil {
		chartsOnSlide = append(chartsOnSlide, s.Line)
	}
	if s.LineMarkers != nil {
		chartsOnSlide = append(chartsOnSlide, s.LineMarkers)
	}
	if s.LineStacked != nil {
		chartsOnSlide = append(chartsOnSlide, s.LineStacked)
	}
	if s.Scatter != nil {
		chartsOnSlide = append(chartsOnSlide, s.Scatter)
	}
	if s.Area != nil {
		chartsOnSlide = append(chartsOnSlide, s.Area)
	}
	if s.AreaStacked != nil {
		chartsOnSlide = append(chartsOnSlide, s.AreaStacked)
	}
	if s.AreaStacked100 != nil {
		chartsOnSlide = append(chartsOnSlide, s.AreaStacked100)
	}
	if s.Pie != nil {
		chartsOnSlide = append(chartsOnSlide, s.Pie)
	}
	if s.Doughnut != nil {
		chartsOnSlide = append(chartsOnSlide, s.Doughnut)
	}
	if s.Bubble != nil {
		chartsOnSlide = append(chartsOnSlide, s.Bubble)
	}
	if s.Radar != nil {
		chartsOnSlide = append(chartsOnSlide, s.Radar)
	}
	if s.RadarFilled != nil {
		chartsOnSlide = append(chartsOnSlide, s.RadarFilled)
	}
	if s.StockHLC != nil {
		chartsOnSlide = append(chartsOnSlide, s.StockHLC)
	}
	if s.StockOHLC != nil {
		chartsOnSlide = append(chartsOnSlide, s.StockOHLC)
	}
	if s.Combo != nil {
		chartsOnSlide = append(chartsOnSlide, s.Combo)
	}
	for _, override := range s.PlaceholderOverrides {
		if override.Chart != nil {
			chartsOnSlide = append(chartsOnSlide, override.Chart)
		}
	}
	return chartsOnSlide
}

func validateSlideTextStyles(s SlideContent) error {
	if err := s.DefaultBulletStyle.Validate(); err != nil {
		return err
	}
	for _, style := range s.BulletStyles {
		if err := style.Validate(); err != nil {
			return err
		}
	}
	for _, runs := range s.BulletRuns {
		for _, run := range runs {
			if err := run.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateSlideTypography(s SlideContent) error {
	if (s.TitleSize != 0 && s.TitleSize < 1) || s.TitleSize > 400 {
		return errors.New("title size must be between 1 and 400 pt (or 0 for default)")
	}
	if s.TitleColor != "" && !common.IsHexColor(s.TitleColor) {
		return errors.New("title color must be 6-digit RGB hex")
	}
	if (s.ContentSize != 0 && s.ContentSize < 1) || s.ContentSize > 400 {
		return errors.New("content size must be between 1 and 400 pt (or 0 for default)")
	}
	if s.ContentColor != "" && !common.IsHexColor(s.ContentColor) {
		return errors.New("content color must be 6-digit RGB hex")
	}
	if s.Background == nil {
		return nil
	}
	if err := s.Background.Validate(); err != nil {
		return fmt.Errorf("invalid background: %w", err)
	}
	return nil
}

func validateSlideTransitionAndTable(s SlideContent, index int) error {
	if err := ValidateSlideTransition(s, index); err != nil {
		return err
	}
	if s.Table == nil {
		return nil
	}
	return s.Table.Validate(index)
}

func validateSlideBulletAndAlignment(s SlideContent) error {
	if slices.Contains(s.Bullets, "") {
		return errors.New("bullet cannot be empty")
	}
	if err := validateTextAlignment(s.TitleAlign, "title alignment", []string{"l", "ctr", "r", "just"}); err != nil {
		return err
	}
	return validateTextAlignment(s.ContentVAlign, "content vertical alignment", []string{"t", "ctr", "b"})
}

func validateTextAlignment(actual string, field string, allowed []string) error {
	if actual == "" {
		return nil
	}
	if slices.Contains(allowed, actual) {
		return nil
	}
	return fmt.Errorf("invalid %s: %q (expected %s)", field, actual, joinAllowed(allowed))
}

func joinAllowed(values []string) string {
	return strings.Join(values, "|")
}

func validateSlideAnimations(s SlideContent) error {
	for i, anim := range s.Animations {
		if err := anim.Validate(); err != nil {
			return err
		}
		if i == 0 && isPrevBasedAnimation(anim.Trigger) {
			return errors.New("first animation trigger cannot be with/after previous")
		}
	}
	return nil
}

func isPrevBasedAnimation(trigger animations.AnimationTrigger) bool {
	return trigger == animations.AnimationWithPrevious || trigger == animations.AnimationAfterPrevious
}

func validateSlideSmartArt(s SlideContent, index int) error {
	for _, diagram := range s.SmartArtDiagrams {
		if err := diagram.Validate(index); err != nil {
			return err
		}
	}
	return nil
}
