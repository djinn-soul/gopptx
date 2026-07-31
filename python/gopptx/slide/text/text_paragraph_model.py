"""Paragraph proxies for the live text object model."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false, reportUnusedClass=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .text_run import Run
from .text_run_model import _ShapeRunCollection

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...schemas import RGBColor
    from .text_frame_protocol import ShapeTextFrameProtocol


class _ShapeParagraphProxy:
    """Live proxy for one indexed shape-text paragraph."""

    def __init__(
        self, text_frame: ShapeTextFrameProtocol, paragraph_index: int
    ) -> None:
        self._text_frame = text_frame
        self._paragraph_index = paragraph_index
        self._runs = _ShapeRunCollection(text_frame, paragraph_index)

    def _paragraph_payload(self) -> dict[str, object]:
        return self._text_frame.get_paragraph_payload(self._paragraph_index)

    def _set_paragraph_field(self, field: str, value: object) -> None:
        self._text_frame.set_paragraph_field(self._paragraph_index, field, value)

    @property
    def font(self) -> _ParagraphFontProxy:
        """Paragraph-level font proxy (Issue #168)."""
        return _ParagraphFontProxy(self)

    @property
    def runs(self) -> _ShapeRunCollection:
        return self._runs

    @property
    def text(self) -> str:
        return "".join(
            str(run.get("text", ""))
            for run in self._text_frame.get_paragraph_runs(self._paragraph_index)
        )

    @text.setter
    def text(self, value: str) -> None:
        self._text_frame.replace_paragraph_runs(
            self._paragraph_index, [Run(text=value).to_payload()]
        )

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
    def rtl(self) -> bool | None:
        """Return whether paragraph has right-to-left text layout enabled."""
        value = self._paragraph_payload().get("rtl")
        return bool(value) if isinstance(value, bool) else None

    @rtl.setter
    def rtl(self, value: bool | None) -> None:
        """Set right-to-left text layout."""
        self._set_paragraph_field("rtl", value)

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

    @property
    def right_to_left(self) -> bool | None:
        """Spelled-out alias for :attr:`rtl` (Issue #1080).

        None means the paragraph states nothing and inherits its direction.
        """
        return self.rtl

    @right_to_left.setter
    def right_to_left(self, value: bool | None) -> None:
        self.rtl = value

    def clear(self) -> None:
        """Clear paragraph text content while keeping the paragraph container."""
        self._text_frame.replace_paragraph_runs(self._paragraph_index, [])

    def remove(self) -> None:
        """Remove this paragraph from the text frame."""
        paragraphs = self._text_frame.get_paragraphs()
        if self._paragraph_index >= len(paragraphs):
            raise IndexError("paragraph index out of range")
        del paragraphs[self._paragraph_index]
        if not paragraphs:
            paragraphs = [cast("dict[str, object]", {"runs": [], "paragraph": {}})]
        self._text_frame.replace_paragraphs(paragraphs)

    def delete(self) -> None:
        """Alias for remove (Issue #144)."""
        self.remove()

    def remove_run(self, run: object) -> None:
        """Remove a run proxy or run index from this paragraph (Issue #144)."""
        if isinstance(run, int):
            self.runs[run].delete()
        else:
            delete_func = getattr(run, "delete", None)
            if callable(delete_func):
                delete_func()
            else:
                raise TypeError("run must be a ShapeRunProxy or run index")


class _ShapeParagraphCollection:
    """Collection facade for indexed shape-text paragraphs."""

    def __init__(self, text_frame: ShapeTextFrameProtocol) -> None:
        self._text_frame = text_frame
        self._paragraphs: dict[int, _ShapeParagraphProxy] = {}

    def __len__(self) -> int:
        return len(self._text_frame.get_paragraphs())

    def __getitem__(self, index: int) -> _ShapeParagraphProxy:
        if index < 0:
            index += len(self)
        if index < 0 or index >= len(self):
            raise IndexError("paragraph index out of range")
        proxy = self._paragraphs.get(index)
        if proxy is None:
            proxy = _ShapeParagraphProxy(self._text_frame, index)
            self._paragraphs[index] = proxy
        return proxy

    def __iter__(self) -> Iterator[_ShapeParagraphProxy]:
        for index in range(len(self)):
            yield self[index]

    def add_paragraph(self, text: str = "") -> _ShapeParagraphProxy:
        """Add a paragraph to the text frame (Issue #1135)."""
        paragraphs = self._text_frame.get_paragraphs()
        paragraphs.append({"runs": [Run(text=text).to_payload()], "paragraph": {}})
        self._text_frame.replace_paragraphs(paragraphs)
        return self[len(paragraphs) - 1]

    @staticmethod
    def remove(paragraph: _ShapeParagraphProxy) -> None:
        """Remove one paragraph from this collection."""
        paragraph.remove()

    def __delitem__(self, index: int) -> None:
        """Delete paragraph at index (Issue #144)."""
        self[index].remove()


class _ParagraphFontProxy:
    """Paragraph-level font proxy for formatting runs in a paragraph (Issue #168)."""

    def __init__(self, paragraph: _ShapeParagraphProxy) -> None:
        self._paragraph = paragraph

    @property
    def name(self) -> str | None:
        runs = list(self._paragraph.runs)
        return runs[0].font.name if runs else None

    @name.setter
    def name(self, value: str | None) -> None:
        for run in self._paragraph.runs:
            run.font.name = value

    @property
    def size(self) -> float | None:
        """Font size in EMU, matching ``run.font.size`` (python-pptx parity)."""
        runs = list(self._paragraph.runs)
        return runs[0].font.size if runs else None

    @size.setter
    def size(self, value: float | None) -> None:
        for run in self._paragraph.runs:
            run.font.size = value

    @property
    def size_pt(self) -> float | None:
        """Font size in points, matching ``run.font.size_pt``."""
        runs = list(self._paragraph.runs)
        return runs[0].font.size_pt if runs else None

    @size_pt.setter
    def size_pt(self, value: float | None) -> None:
        for run in self._paragraph.runs:
            run.font.size_pt = value

    @property
    def bold(self) -> bool | None:
        runs = list(self._paragraph.runs)
        return runs[0].font.bold if runs else None

    @bold.setter
    def bold(self, value: bool | None) -> None:
        for run in self._paragraph.runs:
            run.font.bold = value

    @property
    def italic(self) -> bool | None:
        runs = list(self._paragraph.runs)
        return runs[0].font.italic if runs else None

    @italic.setter
    def italic(self, value: bool | None) -> None:
        for run in self._paragraph.runs:
            run.font.italic = value

    @property
    def underline(self) -> bool | str | None:
        runs = list(self._paragraph.runs)
        return runs[0].font.underline if runs else None

    @underline.setter
    def underline(self, value: bool | str | None) -> None:
        for run in self._paragraph.runs:
            run.font.underline = value

    @property
    def color(self) -> _ParagraphFontColorProxy:
        return _ParagraphFontColorProxy(self._paragraph)


class _ParagraphFontColorProxy:
    """Paragraph-level font color proxy (Issue #168)."""

    def __init__(self, paragraph: _ShapeParagraphProxy) -> None:
        self._paragraph = paragraph

    @property
    def rgb(self) -> RGBColor | None:
        runs = list(self._paragraph.runs)
        return runs[0].font.color.rgb if runs else None

    @rgb.setter
    def rgb(self, value: RGBColor | None) -> None:
        for run in self._paragraph.runs:
            run.font.color.rgb = value
