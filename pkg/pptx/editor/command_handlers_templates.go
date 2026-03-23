package editor

import (
	"encoding/json"
	"errors"

	"github.com/djinn-soul/gopptx/pkg/pptx/templates"
)

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
		Title    string              `json:"title"`
		Subtitle string              `json:"subtitle"`
		Context  string              `json:"context"`
		Solution string              `json:"solution"`
		Pricing  []map[string]any    `json:"pricing"`
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
			if flist, ok := f.([]any); ok {
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
