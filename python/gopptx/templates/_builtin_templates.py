"""Built-in gopptx templates: StatusTemplate and SimpleTemplate."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import TYPE_CHECKING

from typing_extensions import override

from gopptx import ops
from gopptx.presentation.presentation import Presentation

from ._template_utils import Template, apply_slides

if TYPE_CHECKING:
    from gopptx.presentation.theme.theme import Theme


@dataclass
class StatusTemplate(Template):
    """Builds a 4-slide status report template."""

    project: str
    okrs: list[str] = field(default_factory=list)
    risks: list[str] = field(default_factory=list)
    next_steps: list[str] = field(default_factory=list)
    theme: Theme | None = None

    @override
    def build(self) -> Presentation:
        """Build the status template presentation."""
        if not self.project:
            raise ValueError("project name cannot be empty")

        title = f"{self.project} - Status Update"
        prs = Presentation.new(title)

        try:
            result = prs.execute(
                ops.OP_BUILD_STATUS_TEMPLATE,
                {
                    "project": self.project,
                    "okrs": self.okrs,
                    "risks": self.risks,
                    "next_steps": self.next_steps,
                },
            )

            if self.theme:
                prs.apply_theme(self.theme)

            apply_slides(prs, result.get("slides", []))
            return prs
        except Exception:
            prs.close()
            raise


@dataclass
class SimpleTemplate(Template):
    """Builds a 2-slide simple template."""

    title: str
    content: str = ""
    theme: Theme | None = None

    @override
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

            apply_slides(prs, result.get("slides", []))
            return prs
        except Exception:
            prs.close()
            raise
