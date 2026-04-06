package export

import "github.com/signintech/gopdf"

// drawScatterGridlines draws a 4×4 grid inside the scatter plot frame.
func drawScatterGridlines(pdf *gopdf.GoPdf, px, py, pw, ph float64) {
	pdf.SetStrokeColor(90, 90, 90)
	for i := 1; i < 5; i++ {
		yg := py + float64(i)*ph/5
		pdf.Line(px, yg, px+pw, yg)
		xg := px + float64(i)*pw/5
		pdf.Line(xg, py, xg, py+ph)
	}
}
