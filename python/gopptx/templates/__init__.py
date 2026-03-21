"""Template system for gopptx - declarative slide builders."""

from abc import ABC, abstractmethod
from dataclasses import dataclass, field

from gopptx import ops
from gopptx.presentation.presentation import Presentation
from gopptx.presentation.theme.theme import Theme


class Template(ABC):
    """Base class for presentation templates.

    Users can subclass this to create custom templates.
    """

    @abstractmethod
    def build(self) -> Presentation:
        """Build and return a Presentation with template slides."""
        pass


def _apply_slides(prs: Presentation, slides: list[dict]) -> None:
    """Remove the default blank slide and add template slides."""
    # Presentation.new() always starts with 1 blank slide — remove it
    prs.remove_slide(0)
    for i, slide_data in enumerate(slides):
        title_text = slide_data.get("title", "")
        layout = slide_data.get("layout") or None
        bullets = slide_data.get("bullets") or []
        notes = slide_data.get("notes") or ""
        table_data = slide_data.get("table")

        prs.add_slide(title_text, layout=layout, bullets=bullets or None)

        if notes:
            prs.set_notes(i, notes)

        if table_data:
            rows = table_data.get("rows") or []
            if rows:
                bounds = (
                    table_data.get("x", 457200),
                    table_data.get("y", 1600200),
                    table_data.get("cx", 8230200),
                    table_data.get("cy", 3200400),
                )
                prs.add_table_from_rows(i, rows, bounds, first_row=True, band_row=True)


@dataclass
class StatusTemplate(Template):
    """Builds a 4-slide status report template.

    Generates slides for:
    - Title slide with project name
    - OKR Status (with provided OKRs as bullets)
    - Risks & Blockers (with provided risks as bullets)
    - Next Steps (with provided next steps as bullets)
    """

    project: str
    okrs: list[str] = field(default_factory=list)
    risks: list[str] = field(default_factory=list)
    next_steps: list[str] = field(default_factory=list)
    theme: Theme | None = None

    def build(self) -> Presentation:
        """Build the status template presentation."""
        if not self.project:
            raise ValueError("project name cannot be empty")

        # Create presentation first
        title = f"{self.project} - Status Update"
        prs = Presentation.new(title)

        try:
            # Call Go bridge to get slide data
            result = prs.execute(
                ops.OP_BUILD_STATUS_TEMPLATE,
                {
                    "project": self.project,
                    "okrs": self.okrs,
                    "risks": self.risks,
                    "next_steps": self.next_steps,
                },
            )

            # Apply theme if provided
            if self.theme:
                prs.apply_theme(self.theme)

            _apply_slides(prs, result.get("slides", []))
            return prs
        except Exception:
            prs.close()
            raise


@dataclass
class SimpleTemplate(Template):
    """Builds a 2-slide simple template.

    Generates slides for:
    - Title slide
    - Content slide (with bullet if provided)
    """

    title: str
    content: str = ""
    theme: Theme | None = None

    def build(self) -> Presentation:
        """Build the simple template presentation."""
        if not self.title:
            raise ValueError("title cannot be empty")

        prs = Presentation.new(self.title)

        try:
            result = prs.execute(
                ops.OP_BUILD_SIMPLE_TEMPLATE,
                {
                    "title": self.title,
                    "content": self.content,
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
    """Builds a 5-slide proposal template.

    Generates slides for:
    - Title slide
    - Context (problem/background)
    - Solution
    - Pricing tiers table
    - Timeline/milestones table
    """

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
                {"date": m.date, "task": m.task, "status": m.status} for m in self.timeline
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
    """Builds a training template with agenda, concepts, and summary.

    Generates slides for:
    - Title slide
    - Agenda (with provided agenda items)
    - Concept slides (one per concept)
    - Summary slide
    """

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
    """Builds a 4-slide technical deep-dive template.

    Generates slides for:
    - Title slide
    - Architecture Overview
    - Technical Deep Dive
    - Performance & Benchmarks
    """

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


__all__ = [
    "Template",
    "StatusTemplate",
    "SimpleTemplate",
    "ProposalTemplate",
    "TrainingTemplate",
    "TechnicalTemplate",
    "PricingTier",
    "Milestone",
]
