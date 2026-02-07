package pptx

import "fmt"

func validateSlideTextRuns(s SlideContent, slideIndex int) error {
	if len(s.BulletRuns) == 0 {
		return nil
	}
	if len(s.BulletRuns) != len(s.Bullets) {
		return fmt.Errorf(
			"slide %d bullet runs count %d must match bullet count %d",
			slideIndex,
			len(s.BulletRuns),
			len(s.Bullets),
		)
	}

	for bulletIndex := range s.BulletRuns {
		runs := s.BulletRuns[bulletIndex]
		for runIndex, run := range runs {
			if run.Text == "" {
				return fmt.Errorf("slide %d bullet %d run %d text cannot be empty", slideIndex, bulletIndex+1, runIndex+1)
			}
			if run.SizePt < 0 {
				return fmt.Errorf("slide %d bullet %d run %d size must be >= 0", slideIndex, bulletIndex+1, runIndex+1)
			}
			if run.Color != "" && !isHexColor(run.Color) {
				return fmt.Errorf("slide %d bullet %d run %d color must be 6-digit RGB hex", slideIndex, bulletIndex+1, runIndex+1)
			}
		}
	}
	return nil
}
