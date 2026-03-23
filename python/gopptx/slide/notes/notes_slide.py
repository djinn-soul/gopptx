"""Notes slide proxy for gopptx slide API."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from typing_extensions import override

from .notes_shape import NotesShape
from .notes_shape_collection import NotesShapeCollection
from .notes_slide_style_mixin import NotesSlideStyleMixin

if TYPE_CHECKING:
    from ...schemas import ShapeUpdate
    from ._protocols import NotesSlideProto
    from .notes_text_model import NotesTextFrame


class _NotesBackingPresentationProto(Protocol):
    def get_notes_payload(self, slide_index: int) -> dict[str, object]: ...

    def set_notes_shape_text(
        self, slide_index: int, shape_id: int, text: str
    ) -> None: ...

    def set_notes_shape_props(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None: ...


class _NotesBackingSlideProto(Protocol):
    @property
    def notes(self) -> str: ...

    @notes.setter
    def notes(self, value: str) -> None: ...

    @property
    def presentation(self) -> _NotesBackingPresentationProto: ...

    @property
    def index(self) -> int: ...


class NotesSlide(NotesSlideStyleMixin):
    """Proxy for slide notes content."""

    def __init__(self, slide: _NotesBackingSlideProto) -> None:
        """Initialize notes proxy bound to a slide."""
        super().__init__()
        self._slide: _NotesBackingSlideProto = slide

    @property
    def text(self) -> str:
        """Get notes text."""
        return self._slide.notes

    @text.setter
    def text(self, value: str) -> None:
        """Set notes text."""
        self._slide.notes = value

    def _shape_payloads(self) -> list[dict[str, object]]:
        payload = self._slide.presentation.get_notes_payload(self._slide.index)
        raw = payload.get("notes_shapes")
        if isinstance(raw, list):
            return cast("list[dict[str, object]]", raw)
        return self._placeholder_payloads()

    def shape_payloads(self) -> list[dict[str, object]]:
        """Return raw notes shape payloads for collection traversal."""
        return self._shape_payloads()

    def _set_shape_text(self, shape_id: int, text: str) -> None:
        self._slide.presentation.set_notes_shape_text(self._slide.index, shape_id, text)

    @override
    def _set_shape_props(self, shape_id: int, updates: ShapeUpdate) -> None:
        if shape_id < 0:
            raise ValueError("notes shape id is unavailable for mutation")
        self._slide.presentation.set_notes_shape_props(
            self._slide.index,
            shape_id,
            updates,
        )

    def _placeholder_payloads(self) -> list[dict[str, object]]:
        payload = self._slide.presentation.get_notes_payload(self._slide.index)
        raw = payload.get("notes_placeholders")
        return cast("list[dict[str, object]]", raw if isinstance(raw, list) else [])

    @property
    def placeholders(self) -> list[dict[str, object]]:
        """Placeholder metadata from notes slide."""
        return self._placeholder_payloads()

    @property
    def shapes(self) -> NotesShapeCollection:
        """Typed notes shape collection for object-model traversal."""
        return NotesShapeCollection(cast("NotesSlideProto", self))

    def shape(self, shape_id: int) -> NotesShape | None:
        """Return one notes shape by id, if present."""
        return self.shapes.by_id(shape_id)

    def shape_by_name(self, name: str) -> NotesShape | None:
        """Return one notes shape by name, if present."""
        return self.shapes.get(name=name)

    def set_shape_bounds(
        self, shape_id: int, *, left: float, top: float, width: float, height: float
    ) -> None:
        """Set notes-shape geometry in EMU units."""
        self._set_shape_props(
            shape_id,
            {"x": int(left), "y": int(top), "w": int(width), "h": int(height)},
        )

    def set_shape_fill_background(self, shape_id: int) -> None:
        """Set notes-shape fill to background mode."""
        self._set_shape_props(shape_id, {"fill": {"background": True}})

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


__all__ = ["NotesShape", "NotesShapeCollection", "NotesSlide"]
