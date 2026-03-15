"""Live text object model facades for shape text editing."""
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportArgumentType=false, reportCallIssue=false, reportPropertyTypeMismatch=false, reportUnknownArgumentType=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ..shapes.shape_text_frame import ShapeTextFrame as BaseShapeTextFrame
from .text_run import Run, RunHyperlink

if TYPE_CHECKING:
    from collections.abc import Iterator


class ShapeTextFrame(BaseShapeTextFrame):
    """Public text-frame facade with paragraph collection access."""

    def __init__(self, slide: object, shape_id: int) -> None:
        """Create a text-frame proxy for a shape."""
        super().__init__(slide, shape_id)
        self._paragraphs: _ShapeParagraphCollection | None = None

    @property
    def paragraphs(self) -> _ShapeParagraphCollection:
        """Return paragraph collection facade."""
        if self._paragraphs is None:
            self._paragraphs = _ShapeParagraphCollection(self)
        return self._paragraphs


class _ShapeRunProxy:
    """Live run proxy backed by bridge operations."""

    def __init__(self, text_frame: ShapeTextFrame, run_index: int) -> None:
        self._text_frame = text_frame
        self._run_index = run_index

    def _run_payload(self) -> dict[str, object]:
        runs = self._text_frame.get_runs()
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        return cast("dict[str, object]", runs[self._run_index])

    def _set_field(self, name: str, value: object) -> None:
        runs = [dict(run) for run in self._text_frame.get_runs()]
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        runs[self._run_index][name] = value
        self._text_frame.replace_runs(runs)

    @property
    def text(self) -> str:
        return str(self._run_payload().get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._text_frame.update_run_text(self._run_index, value)

    @property
    def bold(self) -> bool | None:
        value = self._run_payload().get("bold")
        return bool(value) if isinstance(value, bool) else None

    @bold.setter
    def bold(self, value: bool | None) -> None:
        self._set_field("bold", value)

    @property
    def italic(self) -> bool | None:
        value = self._run_payload().get("italic")
        return bool(value) if isinstance(value, bool) else None

    @italic.setter
    def italic(self, value: bool | None) -> None:
        self._set_field("italic", value)

    @property
    def underline(self) -> str | None:
        value = self._run_payload().get("underline")
        return str(value) if isinstance(value, str) else None

    @underline.setter
    def underline(self, value: str | None) -> None:
        self._set_field("underline", value)

    @property
    def hyperlink(self) -> RunHyperlink:
        payload = self._run_payload().get("hyperlink")
        if isinstance(payload, dict):
            parsed = RunHyperlink.from_payload(payload)
            if parsed is not None:
                return parsed
        return RunHyperlink()

    @hyperlink.setter
    def hyperlink(self, value: RunHyperlink | dict[str, object] | None) -> None:
        if isinstance(value, RunHyperlink):
            self._set_field("hyperlink", value.to_payload())
            return
        self._set_field("hyperlink", value)


class _ShapeRunCollection:
    """Live run collection for a paragraph proxy."""

    def __init__(self, text_frame: ShapeTextFrame) -> None:
        self._text_frame = text_frame
        self._run_proxies: dict[int, _ShapeRunProxy] = {}

    def __len__(self) -> int:
        return len(self._text_frame.get_runs())

    def __getitem__(self, index: int) -> _ShapeRunProxy:
        if index < 0:
            index += len(self)
        if index < 0:
            raise IndexError("run index out of range")
        proxy = self._run_proxies.get(index)
        if proxy is None:
            proxy = _ShapeRunProxy(self._text_frame, index)
            self._run_proxies[index] = proxy
        return proxy

    def __iter__(self) -> Iterator[_ShapeRunProxy]:
        for index, _ in enumerate(self._text_frame.get_runs()):
            yield _ShapeRunProxy(self._text_frame, index)

    def add_run(self, text: str = "") -> _ShapeRunProxy:
        self._text_frame.append_run(Run(text=text).to_payload())
        return self[len(self) - 1]


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


ShapeRunProxy = _ShapeRunProxy
ShapeRunCollection = _ShapeRunCollection
ShapeParagraphProxy = _ShapeParagraphProxy
ShapeParagraphCollection = _ShapeParagraphCollection
