"""Paragraph proxies for the live text object model."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnusedClass=false

from __future__ import annotations

from typing import TYPE_CHECKING

from .text_run import Run
from .text_run_model import _ShapeRunCollection

if TYPE_CHECKING:
    from collections.abc import Iterator

    from .text_model import ShapeTextFrame


class _ShapeParagraphProxy:
    """Single-paragraph proxy; PPTX run API currently models one normalized paragraph."""

    def __init__(self, text_frame: ShapeTextFrame) -> None:
        self._text_frame = text_frame
        self._runs = _ShapeRunCollection(text_frame)

    def _paragraph_payload(self) -> dict[str, object]:
        return self._text_frame.get_paragraph_payload()

    def _set_paragraph_field(self, field: str, value: object) -> None:
        self._text_frame.set_paragraph_field(field, value)

    @property
    def runs(self) -> _ShapeRunCollection:
        return self._runs

    @property
    def text(self) -> str:
        return "".join(str(run.get("text", "")) for run in self._text_frame.get_runs())

    @text.setter
    def text(self, value: str) -> None:
        self._text_frame.replace_runs([Run(text=value).to_payload()])

    @property
    def alignment(self) -> str | None:
        value = self._paragraph_payload().get("alignment")
        return str(value) if isinstance(value, str) else None

    @alignment.setter
    def alignment(self, value: str | None) -> None:
        self._set_paragraph_field("alignment", value)

    @property
    def level(self) -> int | None:
        value = self._paragraph_payload().get("level")
        return int(value) if isinstance(value, int) else None

    @level.setter
    def level(self, value: int | None) -> None:
        self._set_paragraph_field("level", value)

    @property
    def bullet(self) -> bool:
        style = self._paragraph_payload().get("bullet_style")
        return isinstance(style, str) and style != "none"

    @bullet.setter
    def bullet(self, value: bool) -> None:
        self._set_paragraph_field("bullet_style", "bullet" if value else "none")
        if not value:
            self._set_paragraph_field("bullet_char", None)

    @property
    def bullet_style(self) -> str | None:
        value = self._paragraph_payload().get("bullet_style")
        return str(value) if isinstance(value, str) else None

    @bullet_style.setter
    def bullet_style(self, value: str | None) -> None:
        self._set_paragraph_field("bullet_style", value)

    @property
    def bullet_char(self) -> str | None:
        value = self._paragraph_payload().get("bullet_char")
        return str(value) if isinstance(value, str) else None

    @bullet_char.setter
    def bullet_char(self, value: str | None) -> None:
        self._set_paragraph_field("bullet_char", value)

    @property
    def line_spacing(self) -> float | int | None:
        payload = self._paragraph_payload()
        pct = payload.get("line_spacing_pct")
        if isinstance(pct, int):
            return float(pct) / 100000.0
        pts = payload.get("line_spacing_pts")
        if isinstance(pts, int):
            return pts
        return None

    @line_spacing.setter
    def line_spacing(self, value: float | None) -> None:
        if value is None:
            self._set_paragraph_field("line_spacing_pct", None)
            self._set_paragraph_field("line_spacing_pts", None)
            return
        if isinstance(value, float):
            self._set_paragraph_field("line_spacing_pts", None)
            self._set_paragraph_field("line_spacing_pct", round(value * 100000))
            return
        self._set_paragraph_field("line_spacing_pct", None)
        self._set_paragraph_field("line_spacing_pts", value)

    @property
    def space_before(self) -> int | None:
        value = self._paragraph_payload().get("space_before_pts")
        return int(value) if isinstance(value, int) else None

    @space_before.setter
    def space_before(self, value: int | None) -> None:
        self._set_paragraph_field("space_before_pts", value)

    @property
    def space_after(self) -> int | None:
        value = self._paragraph_payload().get("space_after_pts")
        return int(value) if isinstance(value, int) else None

    @space_after.setter
    def space_after(self, value: int | None) -> None:
        self._set_paragraph_field("space_after_pts", value)

    def clear(self) -> None:
        """Clear paragraph text content while keeping the paragraph container."""
        self._text_frame.replace_runs([])

    def remove(self) -> None:
        """Remove paragraph content (single-paragraph model maps this to clear)."""
        self.clear()


class _ShapeParagraphCollection:
    """Collection facade for the normalized single paragraph."""

    def __init__(self, text_frame: ShapeTextFrame) -> None:
        self._paragraph = _ShapeParagraphProxy(text_frame)

    def __len__(self) -> int:
        return 1

    def __getitem__(self, index: int) -> _ShapeParagraphProxy:
        if index not in {0, -1}:
            raise IndexError("paragraph index out of range")
        return self._paragraph

    def __iter__(self) -> Iterator[_ShapeParagraphProxy]:
        yield self[0]

    def remove(self, paragraph: _ShapeParagraphProxy) -> None:
        """Remove paragraph content in the normalized single-paragraph model."""
        if paragraph is not self._paragraph:
            raise ValueError("paragraph does not belong to this collection")
        paragraph.remove()
