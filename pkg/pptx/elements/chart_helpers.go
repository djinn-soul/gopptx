package elements

func (s SlideContent) ChartKindCount() int {
	return s.directChartCount() + s.placeholderChartCount()
}

func (s SlideContent) directChartCount() int {
	count := 0
	if s.Chart != nil {
		count++
	}
	if s.BarHorizontal != nil {
		count++
	}
	if s.BarStacked != nil {
		count++
	}
	if s.BarStacked100 != nil {
		count++
	}
	if s.Line != nil {
		count++
	}
	if s.LineMarkers != nil {
		count++
	}
	if s.LineStacked != nil {
		count++
	}
	if s.Scatter != nil {
		count++
	}
	if s.Area != nil {
		count++
	}
	if s.AreaStacked != nil {
		count++
	}
	if s.AreaStacked100 != nil {
		count++
	}
	if s.Pie != nil {
		count++
	}
	if s.Doughnut != nil {
		count++
	}
	if s.Bubble != nil {
		count++
	}
	if s.Radar != nil {
		count++
	}
	if s.RadarFilled != nil {
		count++
	}
	if s.StockHLC != nil {
		count++
	}
	if s.StockOHLC != nil {
		count++
	}
	if s.Combo != nil {
		count++
	}
	return count
}

func (s SlideContent) placeholderChartCount() int {
	count := 0
	for _, override := range s.PlaceholderOverrides {
		if override.Chart != nil {
			count++
		}
	}
	return count
}
