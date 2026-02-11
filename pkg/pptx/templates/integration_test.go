package templates_test

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx"
)

func TestTemplateIntegration(t *testing.T) {
	templatesList := []pptx.Template{
		pptx.SimpleTemplate{Title: "Integration Simple", Content: "Content"},
		pptx.ProposalTemplate{Title: "Integration Proposal", Context: "Context", Solution: "Solution", Pricing: []string{"Item 1"}, Timeline: "Soon"},
		pptx.TrainingTemplate{Title: "Integration Training", Agenda: []string{"Step 1"}, Concepts: []string{"Concept A"}, Summary: "End"},
		pptx.StatusTemplate{Project: "Integration Status", OKRs: []string{"OKR 1"}, Risks: []string{"Risk 1"}, NextSteps: []string{"Next"}},
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
