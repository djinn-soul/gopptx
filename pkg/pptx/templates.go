package pptx

import (
	"fmt"
	"sync"
)

// Template defines the interface for high-level presentation builders.
// Implementing this interface allows generating a slice of SlideContent
// that can be passed to CreateWithSlides.
type Template interface {
	// Build generates the slides for the template.
	// It uses parallel generation for independent slides where appropriate.
	Build() ([]SlideContent, error)
}

// SimpleTemplate creates a basic 2-slide deck: a title slide and a content slide.
// Use this for quick single-topic presentations.
type SimpleTemplate struct {
	Title   string // Main title for the first slide
	Content string // Bullet point content for the second slide
}

// Build generates slides for SimpleTemplate.
func (t SimpleTemplate) Build() ([]SlideContent, error) {
	if t.Title == "" {
		return nil, fmt.Errorf("simple template title cannot be empty")
	}

	return buildParallel(
		func() SlideContent {
			return NewSlide(t.Title).WithCenteredTitleLayout()
		},
		func() SlideContent {
			s := NewSlide("Content")
			if t.Content != "" {
				s = s.AddBullet(t.Content)
			}
			return s
		},
	), nil
}

// ProposalTemplate creates a standard 5-slide proposal deck:
// Title, Context, Solution, Pricing, and Timeline.
type ProposalTemplate struct {
	Title    string   // Main proposal title
	Subtitle string   // Optional subtitle (currently rendered in speaker notes)
	Context  string   // Problem or background context
	Solution string   // Proposed solution details
	Pricing  []string // List of pricing items or tiers
	Timeline string   // Project timeline or milestones
}

// Build generates slides for ProposalTemplate.
func (t ProposalTemplate) Build() ([]SlideContent, error) {
	if t.Title == "" {
		return nil, fmt.Errorf("proposal template title cannot be empty")
	}

	return buildParallel(
		func() SlideContent {
			s := NewSlide(t.Title).WithCenteredTitleLayout()
			if t.Subtitle != "" {
				s.Notes = t.Subtitle
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Context")
			if t.Context != "" {
				s = s.AddBullet(t.Context)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Solution")
			if t.Solution != "" {
				s = s.AddBullet(t.Solution)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Pricing")
			for _, item := range t.Pricing {
				s = s.AddBullet(item)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Timeline")
			if t.Timeline != "" {
				s = s.AddBullet(t.Timeline)
			}
			return s
		},
	), nil
}

// TrainingTemplate creates an educational deck with a title, agenda,
// dynamic concept slides, and a summary.
type TrainingTemplate struct {
	Title    string   // Title of the training session
	Agenda   []string // List of topics to be covered
	Concepts []string // Each concept will get its own slide
	Summary  string   // Closing summary or key takeaways
}

// Build generates slides for TrainingTemplate.
func (t TrainingTemplate) Build() ([]SlideContent, error) {
	if t.Title == "" {
		return nil, fmt.Errorf("training template title cannot be empty")
	}

	funcs := make([]func() SlideContent, 0, 3+len(t.Concepts))
	funcs = append(funcs, func() SlideContent {
		return NewSlide(t.Title).WithCenteredTitleLayout()
	})
	funcs = append(funcs, func() SlideContent {
		s := NewSlide("Agenda")
		for _, item := range t.Agenda {
			s = s.AddBullet(item)
		}
		return s
	})

	for _, concept := range t.Concepts {
		c := concept // capture for closure
		funcs = append(funcs, func() SlideContent {
			return NewSlide(c).AddBullet("Details for " + c + "...")
		})
	}

	funcs = append(funcs, func() SlideContent {
		s := NewSlide("Summary")
		if t.Summary != "" {
			s = s.AddBullet(t.Summary)
		}
		return s
	})

	return buildParallel(funcs...), nil
}

// StatusTemplate creates a 4-slide project status report:
// Title, OKR Status, Risks/Blockers, and Next Steps.
type StatusTemplate struct {
	Project   string   // Name of the project
	OKRs      []string // Current status of key metrics or OKRs
	Risks     []string // Active risks or blocking issues
	NextSteps []string // Upcoming tasks or milestones
}

// Build generates slides for StatusTemplate.
func (t StatusTemplate) Build() ([]SlideContent, error) {
	if t.Project == "" {
		return nil, fmt.Errorf("status template project name cannot be empty")
	}

	return buildParallel(
		func() SlideContent {
			return NewSlide(t.Project + " - Status Update").WithCenteredTitleLayout()
		},
		func() SlideContent {
			s := NewSlide("OKR Status")
			for _, okr := range t.OKRs {
				s = s.AddBullet(okr)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Risks & Blockers")
			for _, risk := range t.Risks {
				s = s.AddBullet(risk)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Next Steps")
			for _, step := range t.NextSteps {
				s = s.AddBullet(step)
			}
			return s
		},
	), nil
}

// TechnicalTemplate creates a 4-slide technical deep-dive:
// Title, Architecture, Deep Dive details, and Benchmarks.
type TechnicalTemplate struct {
	Title        string // Subject of the technical presentation
	Architecture string // High-level architecture description
	DeepDive     string // Specific technical details
	Benchmarks   string // Performance or comparison data
}

// Build generates slides for TechnicalTemplate.
func (t TechnicalTemplate) Build() ([]SlideContent, error) {
	if t.Title == "" {
		return nil, fmt.Errorf("technical template title cannot be empty")
	}

	return buildParallel(
		func() SlideContent {
			return NewSlide(t.Title).WithCenteredTitleLayout()
		},
		func() SlideContent {
			s := NewSlide("Architecture Overview")
			if t.Architecture != "" {
				s = s.AddBullet(t.Architecture)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Technical Deep Dive")
			if t.DeepDive != "" {
				s = s.AddBullet(t.DeepDive)
			}
			return s
		},
		func() SlideContent {
			s := NewSlide("Performance & Benchmarks")
			if t.Benchmarks != "" {
				s = s.AddBullet(t.Benchmarks)
			}
			return s
		},
	), nil
}

// buildParallel executes slide generation functions in parallel and returns them in order.
func buildParallel(funcs ...func() SlideContent) []SlideContent {
	slides := make([]SlideContent, len(funcs))
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	for i, f := range funcs {
		go func(idx int, gen func() SlideContent) {
			defer wg.Done()
			slides[idx] = gen()
		}(i, f)
	}

	wg.Wait()
	return slides
}
