"""Live text object model facades for shape text editing."""
# ruff: noqa: D102,D105,D107,TC003,PLR6201
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false, reportArgumentType=false, reportCallIssue=false, reportPropertyTypeMismatch=false, reportUnknownArgumentType=false, reportUnknownMemberType=false, reportAttributeAccessIssue=false

from __future__ import annotations

from collections.abc import Iterator
from typing import TYPE_CHECKING, cast

from .text_run import Run, RunHyperlink

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation


class ShapeRunProxy:
    """Live run proxy backed by bridge operations."""

    def __init__(
        self,
        presentation: Presentation,
        slide_index: int,
        shape_id: int,
        run_index: int,
    ) -> None:
        self._presentation = presentation
        self._slide_index = slide_index
        self._shape_id = shape_id
        self._run_index = run_index

    def _run_payload(self) -> dict[str, object]:
        runs = self._presentation.get_shape_runs(self._slide_index, self._shape_id)
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        return cast("dict[str, object]", runs[self._run_index])

    def _set_field(self, name: str, value: object) -> None:
        runs = self._presentation.get_shape_runs(self._slide_index, self._shape_id)
        if self._run_index < 0 or self._run_index >= len(runs):
            raise IndexError("run index out of range")
        run = dict(cast("dict[str, object]", runs[self._run_index]))
        run[name] = value
        runs[self._run_index] = run
        self._presentation.set_shape_runs(self._slide_index, self._shape_id, runs)

    @property
    def text(self) -> str:
        return str(self._run_payload().get("text", ""))

    @text.setter
    def text(self, value: str) -> None:
        self._presentation.update_shape_run_text(
            self._slide_index, self._shape_id, self._run_index, value
        )

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


class ShapeRunCollection:
    """Live run collection for a paragraph proxy."""

    def __init__(
        self, presentation: Presentation, slide_index: int, shape_id: int
    ) -> None:
        self._presentation = presentation
        self._slide_index = slide_index
        self._shape_id = shape_id

    def _count(self) -> int:
        return len(self._presentation.get_shape_runs(self._slide_index, self._shape_id))

    def __len__(self) -> int:
        return self._count()

    def __getitem__(self, index: int) -> ShapeRunProxy:
        if index < 0:
            index += self._count()
        if index < 0 or index >= self._count():
            raise IndexError("run index out of range")
        return ShapeRunProxy(
            self._presentation, self._slide_index, self._shape_id, index
        )

    def __iter__(self) -> Iterator[ShapeRunProxy]:
        for i in range(self._count()):
            yield ShapeRunProxy(
                self._presentation, self._slide_index, self._shape_id, i
            )

    def add_run(self, text: str = "") -> ShapeRunProxy:
        run = Run(text=text).to_payload()
        self._presentation.append_shape_run(self._slide_index, self._shape_id, run)
        return self[len(self) - 1]


class ShapeParagraphProxy:
    """Single-paragraph proxy; PPTX run API currently models one normalized paragraph."""

    def __init__(
        self,
        presentation: Presentation,
        slide_index: int,
        shape_id: int,
        paragraph_index: int,
    ) -> None:
        self._presentation = presentation
        self._slide_index = slide_index
        self._shape_id = shape_id
        self._paragraph_index = paragraph_index

    @property
    def runs(self) -> ShapeRunCollection:
        return ShapeRunCollection(self._presentation, self._slide_index, self._shape_id)

    @property
    def text(self) -> str:
        return "".join(run.text for run in self.runs)

    @text.setter
    def text(self, value: str) -> None:
        self._presentation.set_shape_runs(
            self._slide_index,
            self._shape_id,
            [Run(text=value).to_payload()],
        )


class ShapeParagraphCollection:
    """Paragraph collection facade for shape text.

    Current bridge semantics expose a single normalized paragraph object.
    """

    def __init__(
        self, presentation: Presentation, slide_index: int, shape_id: int
    ) -> None:
        self._presentation = presentation
        self._slide_index = slide_index
        self._shape_id = shape_id

    def __len__(self) -> int:
        return 1

    def __getitem__(self, index: int) -> ShapeParagraphProxy:
        if index not in (0, -1):
            raise IndexError("paragraph index out of range")
        return ShapeParagraphProxy(
            self._presentation, self._slide_index, self._shape_id, 0
        )

    def __iter__(self) -> Iterator[ShapeParagraphProxy]:
        yield self[0]


class ShapeTextFrame:
    """Live text-frame facade for one shape."""

    def __init__(
        self, presentation: Presentation, slide_index: int, shape_id: int
    ) -> None:
        self._presentation = presentation
        self._slide_index = slide_index
        self._shape_id = shape_id

    @property
    def paragraphs(self) -> ShapeParagraphCollection:
        return ShapeParagraphCollection(
            self._presentation,
            self._slide_index,
            self._shape_id,
        )
