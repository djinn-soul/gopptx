"""Advanced presentation templates: Proposal, Training, Technical."""

from __future__ import annotations

from dataclasses import dataclass, field

from gopptx import ops
from gopptx.presentation.presentation import Presentation
from gopptx.presentation.theme.theme import Theme

from ._template_utils import Template, _apply_slides


@dataclass
class PricingTier:
    """Pricing tier for proposal template."""

    name: str
    price: str
    features: list[str] = field(default_factory=list)


@dataclass
class Milestone:
    """Milestone for proposal template."""

    date: str
    task: str
    status: str = ""


@dataclass
class ProposalTemplate(Template):
    """Builds a 5-slide proposal template."""

    title: str
    subtitle: str = ""
    context: str = ""
    solution: str = ""
    pricing: list[PricingTier] = field(default_factory=list)
    timeline: list[Milestone] = field(default_factory=list)
    theme: Theme | None = None

    def build(self) -> Presentation:
        """Build the proposal template presentation."""
        if not self.title:
            raise ValueError("title cannot be empty")

        prs = Presentation.new(self.title)

        try:
            pricing_data = [
                {"name": p.name, "price": p.price, "features": p.features}
                for p in self.pricing
            ]
            timeline_data = [
                {"date": m.date, "task": m.task, "status": m.status}
                for m in self.timeline
            ]

            result = prs.execute(
                ops.OP_BUILD_PROPOSAL_TEMPLATE,
                {
                    "title": self.title,
                    "subtitle": self.subtitle,
                    "context": self.context,
                    "solution": self.solution,
                    "pricing": pricing_data,
                    "timeline": timeline_data,
                },
            )

            if self.theme:
                prs.apply_theme(self.theme)

            _apply_slides(prs, result.get("slides", []))
            return prs
        except Exception:
            prs.close()
            raise


@dataclass
class TrainingTemplate(Template):
    """Builds a training template with agenda, concepts, and summary."""

    title: str
    agenda: list[str] = field(default_factory=list)
    concepts: list[str] = field(default_factory=list)
    summary: str = ""
    theme: Theme | None = None

    def build(self) -> Presentation:
        """Build the training template presentation."""
        if not self.title:
            raise ValueError("title cannot be empty")

        prs = Presentation.new(self.title)

        try:
            result = prs.execute(
                ops.OP_BUILD_TRAINING_TEMPLATE,
                {
                    "title": self.title,
                    "agenda": self.agenda,
                    "concepts": self.concepts,
                    "summary": self.summary,
                },
            )

            if self.theme:
                prs.apply_theme(self.theme)

            _apply_slides(prs, result.get("slides", []))
            return prs
        except Exception:
            prs.close()
            raise


@dataclass
class TechnicalTemplate(Template):
    """Builds a 4-slide technical deep-dive template."""

    title: str
    architecture: str = ""
    deep_dive: str = ""
    benchmarks: str = ""
    theme: Theme | None = None

    def build(self) -> Presentation:
        """Build the technical template presentation."""
        if not self.title:
            raise ValueError("title cannot be empty")

        prs = Presentation.new(self.title)

        try:
            result = prs.execute(
                ops.OP_BUILD_TECHNICAL_TEMPLATE,
                {
                    "title": self.title,
                    "architecture": self.architecture,
                    "deep_dive": self.deep_dive,
                    "benchmarks": self.benchmarks,
                },
            )

            if self.theme:
                prs.apply_theme(self.theme)

            _apply_slides(prs, result.get("slides", []))
            return prs
        except Exception:
            prs.close()
            raise
