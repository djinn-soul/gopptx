package pptx

import (
	"testing"
)

func TestSimpleTemplate(t *testing.T) {
	tmpl := SimpleTemplate{
		Title:   "Simple Deck",
		Content: "Hello World",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("SimpleTemplate Build error: %v", err)
	}

	if len(slides) != 2 {
		t.Errorf("expected 2 slides, got %d", len(slides))
	}

	if slides[0].Title != "Simple Deck" {
		t.Errorf("expected first slide title 'Simple Deck', got %q", slides[0].Title)
	}
	if slides[1].Title != "Content" {
		t.Errorf("expected second slide title 'Content', got %q", slides[1].Title)
	}
}

func TestProposalTemplate(t *testing.T) {
	tmpl := ProposalTemplate{
		Title:    "Big Project",
		Context:  "The problem",
		Solution: "The solution",
		Pricing:  []string{"$100", "$200"},
		Timeline: "Next month",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("ProposalTemplate Build error: %v", err)
	}

	if len(slides) != 5 {
		t.Errorf("expected 5 slides, got %d", len(slides))
	}
}

func TestTrainingTemplate(t *testing.T) {
	tmpl := TrainingTemplate{
		Title:    "Go Basics",
		Agenda:   []string{"Syntax", "Types"},
		Concepts: []string{"Variables", "Functions", "Structs"},
		Summary:  "Done",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("TrainingTemplate Build error: %v", err)
	}

	// Title + Agenda + 3 Concepts + Summary = 6 slides
	expected := 6
	if len(slides) != expected {
		t.Errorf("expected %d slides, got %d", expected, len(slides))
	}
}

func TestStatusTemplate(t *testing.T) {
	tmpl := StatusTemplate{
		Project:   "Project X",
		OKRs:      []string{"Metric A", "Metric B"},
		Risks:     []string{"Risk 1"},
		NextSteps: []string{"Step 1"},
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("StatusTemplate Build error: %v", err)
	}

	if len(slides) != 4 {
		t.Errorf("expected 4 slides, got %d", len(slides))
	}
}

func TestTechnicalTemplate(t *testing.T) {
	tmpl := TechnicalTemplate{
		Title:        "Microservices",
		Architecture: "Diagram...",
		DeepDive:     "Code...",
		Benchmarks:   "Fast",
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("TechnicalTemplate Build error: %v", err)
	}

	if len(slides) != 4 {
		t.Errorf("expected 4 slides, got %d", len(slides))
	}
}

func TestTemplateValidation(t *testing.T) {
	tmpl := SimpleTemplate{} // missing title
	_, err := tmpl.Build()
	if err == nil {
		t.Error("expected error for empty title in SimpleTemplate")
	}

	tmpl2 := ProposalTemplate{}
	_, err = tmpl2.Build()
	if err == nil {
		t.Error("expected error for empty title in ProposalTemplate")
	}
}

func TestTemplateIntegration(t *testing.T) {
	templates := []Template{
		SimpleTemplate{Title: "Integration Simple", Content: "Content"},
		ProposalTemplate{Title: "Integration Proposal", Context: "Context", Solution: "Solution", Pricing: []string{"Item 1"}, Timeline: "Soon"},
		TrainingTemplate{Title: "Integration Training", Agenda: []string{"Step 1"}, Concepts: []string{"Concept A"}, Summary: "End"},
		StatusTemplate{Project: "Integration Status", OKRs: []string{"OKR 1"}, Risks: []string{"Risk 1"}, NextSteps: []string{"Next"}},
		TechnicalTemplate{Title: "Integration Tech", Architecture: "Arch", DeepDive: "Deep", Benchmarks: "Fast"},
	}

	for _, tmpl := range templates {
		slides, err := tmpl.Build()
		if err != nil {
			t.Errorf("Build error for template %T: %v", tmpl, err)
			continue
		}

		data, err := CreateWithSlides("Integration Test", slides)
		if err != nil {
			t.Errorf("CreateWithSlides error for template %T: %v", tmpl, err)
			continue
		}

		if len(data) == 0 {
			t.Errorf("generated zero bytes for template %T", tmpl)
		}
	}
}
