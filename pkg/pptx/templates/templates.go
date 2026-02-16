package templates

import (
	"errors"
	"sync"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/styling"
)

// Template defines the interface for high-level presentation builders.
type Template interface {
	// Build generates the slides for the template.
	Build() ([]elements.SlideContent, error)
}

// SimpleTemplate creates a basic 2-slide deck: a title slide and a content slide.
type SimpleTemplate struct {
	Title   string // Main title for the first slide
	Content string // Bullet point content for the second slide
}

// Build generates slides for SimpleTemplate.
func (t SimpleTemplate) Build() ([]elements.SlideContent, error) {
	if t.Title == "" {
		return nil, errors.New("simple template title cannot be empty")
	}

	return buildParallel(
		func() elements.SlideContent {
			return elements.NewSlide(t.Title).WithCenteredTitleLayout()
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Content")
			if t.Content != "" {
				s = s.AddBullet(t.Content)
			}
			return s
		},
	), nil
}

// BrandingSpec defines visual branding for a template.
type BrandingSpec struct {
	Theme  *styling.Theme
	Header string
	Footer string
}

// ProposalTemplate creates a standard 5-slide proposal deck.
type ProposalTemplate struct {
	Title    string       // Main proposal title
	Subtitle string       // Optional subtitle
	Context  string       // Problem or background context
	Solution string       // Proposed solution details
	Pricing  []string     // List of pricing items or tiers
	Timeline string       // Project timeline or milestones
	Branding BrandingSpec // Optional branding/theme
}

// Apply applies branding settings to a slide.
func (b BrandingSpec) Apply(s elements.SlideContent) elements.SlideContent {
	if b.Footer != "" {
		s.FooterText = b.Footer
	}
	return s
}

// Build generates slides for ProposalTemplate.
func (t ProposalTemplate) Build() ([]elements.SlideContent, error) {
	if t.Title == "" {
		return nil, errors.New("proposal template title cannot be empty")
	}

	slides := buildParallel(
		func() elements.SlideContent {
			s := elements.NewSlide(t.Title).WithCenteredTitleLayout()
			if t.Subtitle != "" {
				s.Notes = t.Subtitle
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Context")
			if t.Context != "" {
				s = s.AddBullet(t.Context)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Solution")
			if t.Solution != "" {
				s = s.AddBullet(t.Solution)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Pricing")
			for _, item := range t.Pricing {
				s = s.AddBullet(item)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Timeline")
			if t.Timeline != "" {
				s = s.AddBullet(t.Timeline)
			}
			return s
		},
	)

	for i := range slides {
		slides[i] = t.Branding.Apply(slides[i])
	}
	return slides, nil
}

// TrainingTemplate creates an educational deck.
type TrainingTemplate struct {
	Title    string       // Title of the training session
	Agenda   []string     // List of topics to be covered
	Concepts []string     // Each concept will get its own slide
	Summary  string       // Closing summary or key takeaways
	Branding BrandingSpec // Optional branding/theme
}

// Build generates slides for TrainingTemplate.
func (t TrainingTemplate) Build() ([]elements.SlideContent, error) {
	if t.Title == "" {
		return nil, errors.New("training template title cannot be empty")
	}

	//nolint:mnd // Initial capacity for standard training slides (Title, Agenda, Summary)
	funcs := make([]func() elements.SlideContent, 0, 3+len(t.Concepts))
	funcs = append(funcs, func() elements.SlideContent {
		return elements.NewSlide(t.Title).WithCenteredTitleLayout()
	})
	funcs = append(funcs, func() elements.SlideContent {
		s := elements.NewSlide("Agenda")
		for _, item := range t.Agenda {
			s = s.AddBullet(item)
		}
		return s
	})

	for _, concept := range t.Concepts {
		c := concept
		funcs = append(funcs, func() elements.SlideContent {
			return elements.NewSlide(c).AddBullet("Details for " + c + "...")
		})
	}

	funcs = append(funcs, func() elements.SlideContent {
		s := elements.NewSlide("Summary")
		if t.Summary != "" {
			s = s.AddBullet(t.Summary)
		}
		return s
	})

	slides := buildParallel(funcs...)
	for i := range slides {
		slides[i] = t.Branding.Apply(slides[i])
	}
	return slides, nil
}

// StatusTemplate creates a 4-slide project status report.
type StatusTemplate struct {
	Project   string   // Name of the project
	OKRs      []string // Current status of key metrics or OKRs
	Risks     []string // Active risks or blocking issues
	NextSteps []string // Upcoming tasks or milestones
}

// Build generates slides for StatusTemplate.
func (t StatusTemplate) Build() ([]elements.SlideContent, error) {
	if t.Project == "" {
		return nil, errors.New("status template project name cannot be empty")
	}

	slides := buildParallel(
		func() elements.SlideContent {
			return elements.NewSlide(t.Project + " - Status Update").WithCenteredTitleLayout()
		},
		func() elements.SlideContent {
			s := elements.NewSlide("OKR Status")
			for _, okr := range t.OKRs {
				s = s.AddBullet(okr)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Risks & Blockers")
			for _, risk := range t.Risks {
				s = s.AddBullet(risk)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Next Steps")
			for _, step := range t.NextSteps {
				s = s.AddBullet(step)
			}
			return s
		},
	)
	return slides, nil
}

// TechnicalTemplate creates a 4-slide technical deep-dive.
type TechnicalTemplate struct {
	Title        string       // Subject of the technical presentation
	Architecture string       // High-level architecture description
	DeepDive     string       // Specific technical details
	Benchmarks   string       // Performance or comparison data
	Branding     BrandingSpec // Optional branding/theme
}

// Build generates slides for TechnicalTemplate.
func (t TechnicalTemplate) Build() ([]elements.SlideContent, error) {
	if t.Title == "" {
		return nil, errors.New("technical template title cannot be empty")
	}

	slides := buildParallel(
		func() elements.SlideContent {
			return elements.NewSlide(t.Title).WithCenteredTitleLayout()
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Architecture Overview")
			if t.Architecture != "" {
				s = s.AddBullet(t.Architecture)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Technical Deep Dive")
			if t.DeepDive != "" {
				s = s.AddBullet(t.DeepDive)
			}
			return s
		},
		func() elements.SlideContent {
			s := elements.NewSlide("Performance & Benchmarks")
			if t.Benchmarks != "" {
				s = s.AddBullet(t.Benchmarks)
			}
			return s
		},
	)

	for i := range slides {
		slides[i] = t.Branding.Apply(slides[i])
	}
	return slides, nil
}

// buildParallel executes slide generation functions in parallel.
func buildParallel(funcs ...func() elements.SlideContent) []elements.SlideContent {
	slides := make([]elements.SlideContent, len(funcs))
	var wg sync.WaitGroup
	wg.Add(len(funcs))

	for i, f := range funcs {
		go func(idx int, gen func() elements.SlideContent) {
			defer wg.Done()
			slides[idx] = gen()
		}(i, f)
	}

	wg.Wait()
	return slides
}
