"""Notes slide proxy for gopptx slide API."""
# ruff: noqa: D102, D105, D107, SLF001, TC003
# pyright: reportPrivateUsage=false

from __future__ import annotations

from collections.abc import Iterator
from typing import TYPE_CHECKING, cast

from typing_extensions import override

from .notes_text_model import NotesTextFrame

if TYPE_CHECKING:
    from .slide import SlideBase


_TEXT_PLACEHOLDER_TYPES = {"body", "title", "ctrTitle", "subTitle"}


class NotesShape:
    """Proxy for one notes placeholder/shape entry."""

    def __init__(self, notes_slide: NotesSlide, payload: dict[str, object]) -> None:
        super().__init__()
        self._notes_slide = notes_slide
        self._payload = payload

    @property
    def placeholder_type(self) -> str:
        value = self._payload.get("placeholder_type", self._payload.get("Type"))
        return str(value) if isinstance(value, str) else ""

    @property
    def idx(self) -> int:
        value = self._payload.get("placeholder_index", self._payload.get("Index"))
        if isinstance(value, int):
            return value
        return -1

    @property
    def shape_id(self) -> int:
        value = self._payload.get("id", self._payload.get("ID"))
        if isinstance(value, int):
            return value
        return -1

    @property
    def name(self) -> str:
        value = self._payload.get("name", self._payload.get("Name"))
        return str(value) if isinstance(value, str) else ""

    @property
    def shape_type(self) -> str:
        value = self._payload.get("type", self._payload.get("Type"))
        return str(value) if isinstance(value, str) else ""

    @property
    def is_placeholder(self) -> bool:
        return self.idx >= 0

    @property
    def has_text_frame(self) -> bool:
        value = self._payload.get("has_text_frame")
        if isinstance(value, bool):
            return value
        return self.placeholder_type in _TEXT_PLACEHOLDER_TYPES

    @property
    def x(self) -> float:
        value = self._payload.get("x", self._payload.get("X"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def y(self) -> float:
        value = self._payload.get("y", self._payload.get("Y"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def cx(self) -> float:
        value = self._payload.get("cx", self._payload.get("CX"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def cy(self) -> float:
        value = self._payload.get("cy", self._payload.get("CY"))
        if isinstance(value, (int, float)):
            return float(value)
        return 0.0

    @property
    def left(self) -> float:
        return self.x

    @property
    def top(self) -> float:
        return self.y

    @property
    def width(self) -> float:
        return self.cx

    @property
    def height(self) -> float:
        return self.cy

    @property
    def text(self) -> str:
        if self.placeholder_type in _TEXT_PLACEHOLDER_TYPES:
            return self._notes_slide.text
        value = self._payload.get("text", self._payload.get("Text"))
        return str(value) if isinstance(value, str) else ""

    @text.setter
    def text(self, value: str) -> None:
        if self.placeholder_type in _TEXT_PLACEHOLDER_TYPES:
            self._notes_slide.text = value
            return
        if not self.has_text_frame:
            raise ValueError("target notes shape has no text frame")
        if self.shape_id < 0:
            raise ValueError("notes shape id is unavailable for mutation")
        self._notes_slide._set_shape_text(self.shape_id, value)
        self._payload["text"] = value

    @property
    def text_frame(self) -> NotesTextFrame | None:
        if self.placeholder_type not in _TEXT_PLACEHOLDER_TYPES:
            return None
        return NotesTextFrame(self)

    @override
    def __repr__(self) -> str:
        return (
            f"<NotesShape id={self.shape_id} idx={self.idx}"
            f" type={self.shape_type!r} ph={self.placeholder_type!r}"
            f" name={self.name!r}>"
        )


class NotesShapeCollection:
    """Collection facade for notes placeholder/shape entries."""

    def __init__(self, notes_slide: NotesSlide) -> None:
        super().__init__()
        self._notes_slide = notes_slide

    def _items(self) -> list[NotesShape]:
        return [
            NotesShape(self._notes_slide, payload)
            for payload in self._notes_slide._shape_payloads()
        ]

    def __len__(self) -> int:
        return len(self._notes_slide._shape_payloads())

    def __iter__(self) -> Iterator[NotesShape]:
        return iter(self._items())

    def __getitem__(self, index: int) -> NotesShape:
        items = self._items()
        if index < 0:
            index += len(items)
        if index < 0 or index >= len(items):
            raise IndexError("notes shape index out of range")
        return items[index]

    def get(
        self,
        *,
        shape_id: int | None = None,
        idx: int | None = None,
        shape_type: str | None = None,
        placeholder_type: str | None = None,
        name: str | None = None,
    ) -> NotesShape | None:
        for item in self:
            if shape_id is not None and item.shape_id != shape_id:
                continue
            if idx is not None and item.idx != idx:
                continue
            if shape_type is not None and item.shape_type != shape_type:
                continue
            if (
                placeholder_type is not None
                and item.placeholder_type != placeholder_type
            ):
                continue
            if name is not None and item.name != name:
                continue
            return item
        return None

    def find_all(
        self,
        *,
        shape_type: str | None = None,
        placeholder_type: str | None = None,
        has_text_frame: bool | None = None,
    ) -> list[NotesShape]:
        out: list[NotesShape] = []
        for item in self:
            if shape_type is not None and item.shape_type != shape_type:
                continue
            if (
                placeholder_type is not None
                and item.placeholder_type != placeholder_type
            ):
                continue
            if has_text_frame is not None and item.has_text_frame != has_text_frame:
                continue
            out.append(item)
        return out

    def by_id(self, shape_id: int) -> NotesShape | None:
        return self.get(shape_id=shape_id)


class NotesSlide:
    """Proxy for slide notes content."""

    def __init__(self, slide: SlideBase) -> None:
        """Initialize notes proxy bound to a slide."""
        super().__init__()
        self._slide = slide

    @property
    def text(self) -> str:
        """Get notes text."""
        return self._slide.notes

    @text.setter
    def text(self, value: str) -> None:
        """Set notes text."""
        self._slide.notes = value

    def _shape_payloads(self) -> list[dict[str, object]]:
        payload = self._slide._presentation.get_notes_payload(self._slide.index)
        raw = payload.get("notes_shapes")
        if isinstance(raw, list):
            return cast("list[dict[str, object]]", raw)
        return self._placeholder_payloads()

    def _set_shape_text(self, shape_id: int, text: str) -> None:
        self._slide._presentation.set_notes_shape_text(
            self._slide.index, shape_id, text
        )

    def _placeholder_payloads(self) -> list[dict[str, object]]:
        payload = self._slide._presentation.get_notes_payload(self._slide.index)
        raw = payload.get("notes_placeholders")
        return cast("list[dict[str, object]]", raw if isinstance(raw, list) else [])

    @property
    def placeholders(self) -> list[dict[str, object]]:
        """Placeholder metadata from notes slide."""
        return self._placeholder_payloads()

    @property
    def shapes(self) -> NotesShapeCollection:
        """Typed notes shape collection for object-model traversal."""
        return NotesShapeCollection(self)

    def shape(self, shape_id: int) -> NotesShape | None:
        """Return one notes shape by id, if present."""
        return self.shapes.by_id(shape_id)

    @property
    def body_shape(self) -> NotesShape | None:
        """Return the notes body placeholder proxy when available."""
        return self.shapes.get(placeholder_type="body")

    @property
    def text_shapes(self) -> list[NotesShape]:
        """Return notes shapes that expose text-frame behavior."""
        return self.shapes.find_all(has_text_frame=True)

    @property
    def placeholder_shapes(self) -> list[NotesShape]:
        """Return notes shapes marked as placeholders."""
        return [shape for shape in self.shapes if shape.is_placeholder]

    @property
    def notes_text_frame(self) -> NotesTextFrame | None:
        """python-pptx-like notes text-frame alias for the body placeholder."""
        body = self.body_shape
        if body is None:
            return None
        return body.text_frame

    @override
    def __repr__(self) -> str:
        """Return debug representation for notes proxy."""
        return f"<NotesSlide slide_index={self._slide.index}>"
