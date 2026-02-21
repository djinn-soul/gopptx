package pptx

import (
	"github.com/djinn-soul/gopptx/pkg/pptx/templates"
)

type (
	Template          = templates.Template
	SimpleTemplate    = templates.SimpleTemplate
	ProposalTemplate  = templates.ProposalTemplate
	TrainingTemplate  = templates.TrainingTemplate
	StatusTemplate    = templates.StatusTemplate
	TechnicalTemplate = templates.TechnicalTemplate
	BrandingPreset    = templates.BrandingPreset
	BrandingSpec      = templates.BrandingSpec
	PricingTier       = templates.PricingTier
	Milestone         = templates.Milestone
)

const (
	PresetCorporate = templates.PresetCorporate
	PresetModern    = templates.PresetModern
	PresetCreative  = templates.PresetCreative
)
