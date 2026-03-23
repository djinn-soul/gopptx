package presentation

import (
	"strings"

	"github.com/djinn-soul/gopptx/internal/pptxxml"
	"github.com/djinn-soul/gopptx/pkg/pptx/shapes"
)

const minPlaceholderOverridesForMerge = 2

func mergePlaceholderOverrides(overrides []shapes.PlaceholderContent) []shapes.PlaceholderContent {
	if len(overrides) < minPlaceholderOverridesForMerge {
		return overrides
	}

	merged := make([]shapes.PlaceholderContent, 0, len(overrides))
	for _, next := range overrides {
		idx := findMergeablePlaceholderIndex(merged, next)
		if idx < 0 {
			merged = append(merged, clonePlaceholderContent(next))
			continue
		}
		merged[idx] = mergePlaceholderContent(merged[idx], next)
	}
	return merged
}

func findMergeablePlaceholderIndex(current []shapes.PlaceholderContent, next shapes.PlaceholderContent) int {
	nextRawTypeEmpty := strings.TrimSpace(next.Type) == ""
	for i := range current {
		if current[i].Index != next.Index {
			continue
		}
		currentRawTypeEmpty := strings.TrimSpace(current[i].Type) == ""
		curType := pptxxml.NormalizePlaceholderType(current[i].Type)
		nextType := pptxxml.NormalizePlaceholderType(next.Type)
		if curType == nextType || currentRawTypeEmpty || nextRawTypeEmpty {
			return i
		}
	}
	return -1
}

func mergePlaceholderContent(base, next shapes.PlaceholderContent) shapes.PlaceholderContent {
	if next.Type != "" {
		base.Type = next.Type
	}
	if next.Text != "" {
		base.Text = next.Text
	}
	if next.Image != nil {
		base.Image = next.Image
		base.Table = nil
		base.Chart = nil
	}
	if next.Table != nil {
		base.Table = next.Table
		base.Image = nil
		base.Chart = nil
	}
	if next.Chart != nil {
		base.Chart = next.Chart
		base.Image = nil
		base.Table = nil
	}
	if next.Target != nil {
		base.Target = clonePlaceholderTarget(next.Target)
	}
	if next.Override != nil {
		base.Override = mergeOverrideOptions(base.Override, next.Override)
	}
	return base
}

func clonePlaceholderContent(in shapes.PlaceholderContent) shapes.PlaceholderContent {
	out := in
	out.Target = clonePlaceholderTarget(in.Target)
	out.Override = cloneOverrideOptions(in.Override)
	return out
}

func clonePlaceholderTarget(in *shapes.PlaceholderTarget) *shapes.PlaceholderTarget {
	if in == nil {
		return nil
	}
	copyVal := *in
	return &copyVal
}

func mergeOverrideOptions(base, next *shapes.PlaceholderOverrideOptions) *shapes.PlaceholderOverrideOptions {
	cloned := cloneOverrideOptions(base)
	if cloned == nil {
		cloned = &shapes.PlaceholderOverrideOptions{}
	}
	if next.X != nil {
		cloned.X = next.X
	}
	if next.Y != nil {
		cloned.Y = next.Y
	}
	if next.CX != nil {
		cloned.CX = next.CX
	}
	if next.CY != nil {
		cloned.CY = next.CY
	}
	if next.TextStyle != nil {
		cloned.TextStyle = mergeTextStyles(cloned.TextStyle, next.TextStyle)
	}
	if next.ForceRect != nil {
		cloned.ForceRect = next.ForceRect
	}
	return cloned
}

func cloneOverrideOptions(in *shapes.PlaceholderOverrideOptions) *shapes.PlaceholderOverrideOptions {
	if in == nil {
		return nil
	}
	out := *in
	out.TextStyle = cloneTextStyle(in.TextStyle)
	return &out
}

func mergeTextStyles(base, next *shapes.PlaceholderTextStyle) *shapes.PlaceholderTextStyle {
	cloned := cloneTextStyle(base)
	if cloned == nil {
		cloned = &shapes.PlaceholderTextStyle{}
	}
	if next.SizePt != nil {
		cloned.SizePt = next.SizePt
	}
	if next.Color != nil {
		cloned.Color = next.Color
	}
	if next.Bold != nil {
		cloned.Bold = next.Bold
	}
	if next.Italic != nil {
		cloned.Italic = next.Italic
	}
	if next.Underline != nil {
		cloned.Underline = next.Underline
	}
	if next.Align != nil {
		cloned.Align = next.Align
	}
	if next.Font != nil {
		cloned.Font = next.Font
	}
	return cloned
}

func cloneTextStyle(in *shapes.PlaceholderTextStyle) *shapes.PlaceholderTextStyle {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}
