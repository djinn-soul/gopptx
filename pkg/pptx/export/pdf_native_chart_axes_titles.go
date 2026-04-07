package export

import "github.com/signintech/gopdf"

//nolint:mnd // Axis title placement uses fixed offsets that match PPT defaults.
func drawCategoryAxisTitle(pdf *gopdf.GoPdf, px, py, pw, ph float64, title string) {
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(px + pw/2 - float64(len(title))*3)
	pdf.SetY(py + ph + 26)
	_ = pdf.Cell(nil, title)
	pdf.SetTextColor(0, 0, 0)
}

//nolint:mnd // Axis title placement uses fixed offsets that match PPT defaults.
func drawValueAxisTitle(pdf *gopdf.GoPdf, px, py, _ float64, ph float64, title string) {
	pdf.SetTextColor(60, 60, 60)
	titleX := px - 42
	titleY := py + ph/2
	pdf.Rotate(-90, titleX, titleY)
	pdf.SetX(titleX - float64(len(title))*3)
	pdf.SetY(titleY - 3)
	_ = pdf.Cell(nil, title)
	pdf.RotateReset()
	pdf.SetTextColor(0, 0, 0)
}
