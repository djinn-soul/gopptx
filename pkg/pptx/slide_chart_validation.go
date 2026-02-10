package pptx

func validateSlideCharts(s SlideContent, index int) error {
	// Primary charts (legacy fields)
	if s.Chart != nil {
		if err := s.Chart.Validate(index); err != nil {
			return err
		}
	}
	if s.BarHorizontal != nil {
		if err := s.BarHorizontal.Validate(index); err != nil {
			return err
		}
	}
	if s.BarStacked != nil {
		if err := s.BarStacked.Validate(index); err != nil {
			return err
		}
	}
	if s.BarStacked100 != nil {
		if err := s.BarStacked100.Validate(index); err != nil {
			return err
		}
	}
	if s.Line != nil {
		if err := s.Line.Validate(index); err != nil {
			return err
		}
	}
	if s.LineMarkers != nil {
		if err := s.LineMarkers.Validate(index); err != nil {
			return err
		}
	}
	if s.LineStacked != nil {
		if err := s.LineStacked.Validate(index); err != nil {
			return err
		}
	}
	if s.Scatter != nil {
		if err := s.Scatter.Validate(index); err != nil {
			return err
		}
	}
	if s.Area != nil {
		if err := s.Area.Validate(index); err != nil {
			return err
		}
	}
	if s.AreaStacked != nil {
		if err := s.AreaStacked.Validate(index); err != nil {
			return err
		}
	}
	if s.AreaStacked100 != nil {
		if err := s.AreaStacked100.Validate(index); err != nil {
			return err
		}
	}
	if s.Pie != nil {
		if err := s.Pie.Validate(index); err != nil {
			return err
		}
	}
	if s.Doughnut != nil {
		if err := s.Doughnut.Validate(index); err != nil {
			return err
		}
	}
	if s.Bubble != nil {
		if err := s.Bubble.Validate(index); err != nil {
			return err
		}
	}
	if s.Radar != nil {
		if err := s.Radar.Validate(index); err != nil {
			return err
		}
	}
	if s.RadarFilled != nil {
		if err := s.RadarFilled.Validate(index); err != nil {
			return err
		}
	}
	if s.StockHLC != nil {
		if err := s.StockHLC.Validate(index); err != nil {
			return err
		}
	}
	if s.StockOHLC != nil {
		if err := s.StockOHLC.Validate(index); err != nil {
			return err
		}
	}
	if s.Combo != nil {
		if err := s.Combo.Validate(index); err != nil {
			return err
		}
	}

	// Placeholder charts
	for _, override := range s.PlaceholderOverrides {
		if override.Chart != nil {
			if err := override.Chart.Validate(index); err != nil {
				return err
			}
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

	for _, override := range s.PlaceholderOverrides {
		if override.Chart != nil {
			count++
		}
	}
	return count
}
