package editor

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/djinn-soul/gopptx/pkg/pptx/elements"
	"github.com/djinn-soul/gopptx/pkg/pptx/templates"
	"github.com/noirbizarre/gonja"
)

// toString safely converts an interface{} to string.
func toString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// slideDataForJSON converts a SlideContent to a JSON-serializable map.
func slideDataForJSON(slide elements.SlideContent) map[string]any {
	data := map[string]any{
		"title":   slide.Title,
		"layout":  slide.Layout,
		"bullets": slide.Bullets,
	}
	return data
}

// handleBuildStatusTemplate builds a StatusTemplate and returns the slides as JSON.
func handleBuildStatusTemplate(_ *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Project   string   `json:"project"`
		OKRs      []string `json:"okrs"`
		Risks     []string `json:"risks"`
		NextSteps []string `json:"next_steps"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid status template payload")
	}

	template := templates.StatusTemplate{
		Project:   req.Project,
		OKRs:      req.OKRs,
		Risks:     req.Risks,
		NextSteps: req.NextSteps,
	}

	slides, err := template.Build()
	if err != nil {
		return nil, err
	}

	slidesData := make([]map[string]any, len(slides))
	for i, slide := range slides {
		slidesData[i] = slideDataForJSON(slide)
	}

	return map[string]any{"slides": slidesData}, nil
}

// handleBuildSimpleTemplate builds a SimpleTemplate and returns the slides as JSON.
func handleBuildSimpleTemplate(_ *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid simple template payload")
	}

	template := templates.SimpleTemplate{
		Title:   req.Title,
		Content: req.Content,
	}

	slides, err := template.Build()
	if err != nil {
		return nil, err
	}

	slidesData := make([]map[string]any, len(slides))
	for i, slide := range slides {
		slidesData[i] = slideDataForJSON(slide)
	}

	return map[string]any{"slides": slidesData}, nil
}

// handleBuildProposalTemplate builds a ProposalTemplate and returns the slides as JSON.
func handleBuildProposalTemplate(_ *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Title    string   `json:"title"`
		Subtitle string   `json:"subtitle"`
		Context  string   `json:"context"`
		Solution string   `json:"solution"`
		Pricing  []map[string]any `json:"pricing"`
		Timeline []map[string]string `json:"timeline"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid proposal template payload")
	}

	// Convert pricing and timeline to proper types
	pricing := make([]templates.PricingTier, len(req.Pricing))
	for i, p := range req.Pricing {
		features := []string{}
		if f, ok := p["features"]; ok {
			if flist, ok := f.([]interface{}); ok {
				for _, fi := range flist {
					if fs, ok := fi.(string); ok {
						features = append(features, fs)
					}
				}
			}
		}
		pricing[i] = templates.PricingTier{
			Name:     toString(p["name"]),
			Price:    toString(p["price"]),
			Features: features,
		}
	}

	timeline := make([]templates.Milestone, len(req.Timeline))
	for i, m := range req.Timeline {
		timeline[i] = templates.Milestone{
			Date:   m["date"],
			Task:   m["task"],
			Status: m["status"],
		}
	}

	template := templates.ProposalTemplate{
		Title:    req.Title,
		Subtitle: req.Subtitle,
		Context:  req.Context,
		Solution: req.Solution,
		Pricing:  pricing,
		Timeline: timeline,
	}

	slides, err := template.Build()
	if err != nil {
		return nil, err
	}

	slidesData := make([]map[string]any, len(slides))
	for i, slide := range slides {
		slidesData[i] = slideDataForJSON(slide)
	}

	return map[string]any{"slides": slidesData}, nil
}

// handleBuildTrainingTemplate builds a TrainingTemplate and returns the slides as JSON.
func handleBuildTrainingTemplate(_ *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Title    string   `json:"title"`
		Agenda   []string `json:"agenda"`
		Concepts []string `json:"concepts"`
		Summary  string   `json:"summary"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid training template payload")
	}

	template := templates.TrainingTemplate{
		Title:    req.Title,
		Agenda:   req.Agenda,
		Concepts: req.Concepts,
		Summary:  req.Summary,
	}

	slides, err := template.Build()
	if err != nil {
		return nil, err
	}

	slidesData := make([]map[string]any, len(slides))
	for i, slide := range slides {
		slidesData[i] = slideDataForJSON(slide)
	}

	return map[string]any{"slides": slidesData}, nil
}

// handleBuildTechnicalTemplate builds a TechnicalTemplate and returns the slides as JSON.
func handleBuildTechnicalTemplate(_ *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Title        string `json:"title"`
		Architecture string `json:"architecture"`
		DeepDive     string `json:"deep_dive"`
		Benchmarks   string `json:"benchmarks"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid technical template payload")
	}

	template := templates.TechnicalTemplate{
		Title:        req.Title,
		Architecture: req.Architecture,
		DeepDive:     req.DeepDive,
		Benchmarks:   req.Benchmarks,
	}

	slides, err := template.Build()
	if err != nil {
		return nil, err
	}

	slidesData := make([]map[string]any, len(slides))
	for i, slide := range slides {
		slidesData[i] = slideDataForJSON(slide)
	}

	return map[string]any{"slides": slidesData}, nil
}

// handleRenderTemplate renders all Jinja2 template expressions across every
// slide shape using the provided context map.  It supports the full Jinja2
// syntax (variables, filters, blocks, loops) via the gonja library.
//
// Each line of shape text that contains a Jinja2 expression is rendered
// independently so that run-level formatting (bold, colour, etc.) is
// preserved via find-and-replace.
func handleRenderTemplate(e *PresentationEditor, payload json.RawMessage) (any, error) {
	var req struct {
		Context map[string]any `json:"context"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, errors.New("invalid render_template payload: expected {\"context\": {...}}")
	}
	if len(req.Context) == 0 {
		return map[string]any{"replacements": 0}, nil
	}

	ctx := gonja.Context(req.Context)
	total := 0
	seen := make(map[string]string) // raw line -> rendered line

	for slideIdx := 0; slideIdx < e.SlideCount(); slideIdx++ {
		shapes, err := e.GetShapes(slideIdx)
		if err != nil {
			continue
		}
		for _, shape := range shapes {
			if !strings.Contains(shape.Text, "{{") && !strings.Contains(shape.Text, "{%") {
				continue
			}
			for _, line := range strings.Split(shape.Text, "\n") {
				if !strings.Contains(line, "{{") && !strings.Contains(line, "{%") {
					continue
				}
				if _, already := seen[line]; already {
					continue
				}
				tpl, err := gonja.FromString(line)
				if err != nil {
					seen[line] = line // keep on parse error
					continue
				}
				rendered, err := tpl.Execute(ctx)
				if err != nil {
					seen[line] = line // keep on render error
					continue
				}
				seen[line] = rendered
			}
		}
	}

	for raw, rendered := range seen {
		if raw != rendered {
			n, _ := e.FindAndReplaceInShapes(raw, rendered)
			total += n
		}
	}

	return map[string]any{"replacements": total}, nil
}
