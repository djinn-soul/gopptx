package export

import "github.com/djinn-soul/gopptx/pkg/pptx/charts"

func applyPieLikeChart(ctx chartApplyCtx) {
	slide, pc, cats, vals, title, px, py, pw, ph :=
		ctx.slide, ctx.pc, ctx.cats, ctx.vals, ctx.title, ctx.px, ctx.py, ctx.pw, ctx.ph
	switch pc.Kind {
	case chartKindPie:
		c := charts.NewPieChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		ctx.applyLegend(&c.ShowLegend, &c.LegendPosition)
		slide.Pie = &c
	case chartKindPie3D:
		c := charts.NewPie3DChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		ctx.applyLegend(&c.ShowLegend, &c.LegendPosition)
		slide.Pie3D = &c
	case chartKindDoughnut:
		c := charts.NewDoughnutChart(cats, vals).WithTitle(title).Position(px, py).Size(pw, ph)
		c.AltText, c.IsDecorative = pc.AltText, pc.IsDecorative
		ctx.applyLegend(&c.ShowLegend, &c.LegendPosition)
		slide.Doughnut = &c
	}
}
