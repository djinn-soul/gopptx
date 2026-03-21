"""Template system for gopptx - declarative slide builders."""

from __future__ import annotations

from ._advanced_templates import (
    Milestone,
    PricingTier,
    ProposalTemplate,
    TechnicalTemplate,
    TrainingTemplate,
)
from ._builtin_templates import SimpleTemplate, StatusTemplate
from ._template_utils import Template, _apply_slides

__all__ = [
    "Milestone",
    "PricingTier",
    "ProposalTemplate",
    "SimpleTemplate",
    "StatusTemplate",
    "TechnicalTemplate",
    "Template",
    "TrainingTemplate",
    "_apply_slides",
]
