"""Shared utilities for gopptx template builders."""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from gopptx.presentation.presentation import Presentation


class Template(ABC):
    """Base class for presentation templates."""

    @abstractmethod
    def build(self) -> Presentation:
        """Build and return a Presentation with template slides."""


def _apply_slides(prs: Presentation, slides: list[dict]) -> None:
    """Remove the default blank slide and add template slides."""
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
