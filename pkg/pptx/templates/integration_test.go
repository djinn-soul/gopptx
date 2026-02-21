package templates_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestTemplateIntegration(t *testing.T) {
	templatesList := []pptx.Template{
		pptx.SimpleTemplate{Title: "Integration Simple", Content: "Content"},
		pptx.ProposalTemplate{
			Title:    "Integration Proposal",
			Context:  "Context",
			Solution: "Solution",
			Pricing: []pptx.PricingTier{
				{Name: "Item 1", Price: "$10", Features: []string{"Basic"}},
			},
			Timeline: []pptx.Milestone{
				{Date: "2026-01-01", Task: "Start", Status: "Done"},
			},
		},
		pptx.TrainingTemplate{
			Title:    "Integration Training",
			Agenda:   []string{"Step 1"},
			Concepts: []string{"Concept A"},
			Summary:  "End",
		},
		pptx.StatusTemplate{
			Project:   "Integration Status",
			OKRs:      []string{"OKR 1"},
			Risks:     []string{"Risk 1"},
			NextSteps: []string{"Next"},
		},
		pptx.TechnicalTemplate{Title: "Integration Tech", Architecture: "Arch", DeepDive: "Deep", Benchmarks: "Fast"},
	}

	for _, tmpl := range templatesList {
		slides, err := tmpl.Build()
		if err != nil {
			t.Errorf("Build error for template %T: %v", tmpl, err)
			continue
		}

		data, err := pptx.CreateWithSlides("Integration Test", slides)
		if err != nil {
			t.Errorf("CreateWithSlides error for template %T: %v", tmpl, err)
			continue
		}

		if len(data) == 0 {
			t.Errorf("generated zero bytes for template %T", tmpl)
		}
	}
}

func TestProposalTemplateStructuredTableEmission(t *testing.T) {
	tmpl := pptx.ProposalTemplate{
		Title: "Structured Proposal",
		Pricing: []pptx.PricingTier{
			{Name: "Basic", Price: "$100", Features: []string{"Feature 1", "Feature 2"}},
		},
		Timeline: []pptx.Milestone{
			{Date: "2026-03-01", Task: "Kickoff", Status: "Planned"},
		},
	}

	slides, err := tmpl.Build()
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	if got, want := len(slides), 5; got != want {
		t.Fatalf("expected %d slides, got %d", want, got)
	}

	if slides[3].Table == nil {
		t.Fatalf("expected pricing slide to include table")
	}
	if slides[4].Table == nil {
		t.Fatalf("expected timeline slide to include table")
	}

	pricingRows := slides[3].Table.Rows
	if pricingRows[0][0] != "Tier" || pricingRows[1][0] != "Basic" || pricingRows[1][1] != "$100" {
		t.Fatalf("unexpected pricing table rows: %#v", pricingRows)
	}

	timelineRows := slides[4].Table.Rows
	if timelineRows[0][0] != "Date" || timelineRows[1][0] != "2026-03-01" || timelineRows[1][1] != "Kickoff" {
		t.Fatalf("unexpected timeline table rows: %#v", timelineRows)
	}
}
