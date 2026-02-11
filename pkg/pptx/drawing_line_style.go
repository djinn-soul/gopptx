package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
)

const (
	LineDashSolid       = elements.LineDashSolid
	LineDashDash        = elements.LineDashDash
	LineDashDot         = elements.LineDashDot
	LineDashDashDot     = elements.LineDashDashDot
	LineDashDashDotDot  = elements.LineDashDashDotDot
	LineDashLongDash    = elements.LineDashLongDash
	LineDashLongDashDot = elements.LineDashLongDashDot
)

func normalizeDrawingLineDash(dash string) string {
	return elements.NormalizeDrawingLineDash(dash)
}
