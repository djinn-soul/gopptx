package pptx

func validateSlideCharts(s SlideContent, index int) error {
	if s.Chart != nil {
		if err := validateBarChart(*s.Chart, index); err != nil {
			return err
		}
	}
	if s.BarHorizontal != nil {
		if err := validateBarHorizontalChart(*s.BarHorizontal, index); err != nil {
			return err
		}
	}
	if s.BarStacked != nil {
		if err := validateBarStackedChart(*s.BarStacked, index); err != nil {
			return err
		}
	}
	if s.BarStacked100 != nil {
		if err := validateBarStacked100Chart(*s.BarStacked100, index); err != nil {
			return err
		}
	}
	if s.Line != nil {
		if err := validateLineChart(*s.Line, index); err != nil {
			return err
		}
	}
	if s.LineMarkers != nil {
		if err := validateLineMarkersChart(*s.LineMarkers, index); err != nil {
			return err
		}
	}
	if s.LineStacked != nil {
		if err := validateLineStackedChart(*s.LineStacked, index); err != nil {
			return err
		}
	}
	if s.Scatter != nil {
		if err := validateScatterChart(*s.Scatter, index); err != nil {
			return err
		}
	}
	if s.Area != nil {
		if err := validateAreaChart(*s.Area, index); err != nil {
			return err
		}
	}
	if s.AreaStacked != nil {
		if err := validateAreaStackedChart(*s.AreaStacked, index); err != nil {
			return err
		}
	}
	if s.AreaStacked100 != nil {
		if err := validateAreaStacked100Chart(*s.AreaStacked100, index); err != nil {
			return err
		}
	}
	if s.Pie != nil {
		if err := validatePieChart(*s.Pie, index); err != nil {
			return err
		}
	}
	if s.Dough != nil {
		if err := validateDoughnutChart(*s.Dough, index); err != nil {
			return err
		}
	}
	if s.Bubble != nil {
		if err := validateBubbleChart(*s.Bubble, index); err != nil {
			return err
		}
	}
	if s.Radar != nil {
		if err := validateRadarChart(*s.Radar, index); err != nil {
			return err
		}
	}
	if s.RadarFilled != nil {
		if err := validateRadarFilledChart(*s.RadarFilled, index); err != nil {
			return err
		}
	}
	if s.StockHLC != nil {
		if err := validateStockHLCChart(*s.StockHLC, index); err != nil {
			return err
		}
	}
	if s.StockOHLC != nil {
		if err := validateStockOHLCChart(*s.StockOHLC, index); err != nil {
			return err
		}
	}
	if s.Combo != nil {
		if err := validateComboChart(*s.Combo, index); err != nil {
			return err
		}
	}
	return nil
}

func chartKindCount(s SlideContent) int {
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
	if s.Dough != nil {
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
