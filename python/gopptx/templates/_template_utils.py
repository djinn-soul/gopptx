"""Shared utilities for gopptx template builders."""

from __future__ import annotations

from abc import ABC, abstractmethod
from typing import TYPE_CHECKING, Protocol, cast

try:
    from typing import NotRequired, TypedDict
except ImportError:
    from typing_extensions import NotRequired, TypedDict

if TYPE_CHECKING:
    from gopptx.presentation.presentation import Presentation


class _TableData(TypedDict, total=False):
    rows: list[list[str]]
    x: int
    y: int
    cx: int
    cy: int


class SlideData(TypedDict, total=False):
    """Typed structure for a single slide definition in a template."""

    title: str
    layout: NotRequired[str | None]
    bullets: NotRequired[list[str] | None]
    notes: NotRequired[str]
    table: NotRequired[_TableData | None]


class Template(ABC):
    """Base class for presentation templates."""

    @abstractmethod
    def build(self) -> Presentation:
        """Build and return a Presentation with template slides."""


class TemplatePresentationProtocol(Protocol):
    """Presentation methods required by template slide application."""

    def remove_slide(self, index: int) -> None: ...

    def add_slide(
        self,
        title: str,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> object: ...

    def set_notes(self, slide_index: int, text: str) -> None: ...

    def add_table_from_rows(
        self,
        slide: int,
        rows: list[list[str]],
        bounds: tuple[int, int, int, int] | None = None,
        *,
        first_row: bool = True,
        band_row: bool = True,
        column_widths: list[int] | None = None,
        **kwargs: object,
    ) -> object: ...


def _require_slide_items(slides: object) -> list[SlideData]:
    """Validate raw template slide data before mutating the presentation."""
    if not isinstance(slides, list):
        raise TypeError("template slides must be a list")

    slide_items = cast("list[object]", slides)
    for index, raw_slide in enumerate(slide_items):
        if not isinstance(raw_slide, dict):
            raise TypeError(f"template slide {index} must be a mapping")
        slide = cast("dict[str, object]", raw_slide)
        table_data = slide.get("table")
        if table_data is not None and not isinstance(table_data, dict):
            raise TypeError(f"template slide {index} table must be a mapping")

    return cast("list[SlideData]", slide_items)


def require_template_slides(result: object) -> list[SlideData]:
    """Extract validated template slides from a builder result payload."""
    if not isinstance(result, dict):
        raise TypeError("template builder result must be a mapping")
    result_dict = cast("dict[str, object]", result)
    if "slides" not in result_dict:
        raise KeyError("template builder result missing 'slides'")
    return _require_slide_items(result_dict["slides"])


def apply_slides(prs: TemplatePresentationProtocol, slides: object) -> None:
    """Remove the default blank slide and add template slides."""
    slide_items = _require_slide_items(slides)
    prs.remove_slide(0)
    for i, slide_data in enumerate(slide_items):
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
