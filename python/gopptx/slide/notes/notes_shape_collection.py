"""Notes shape collection facade."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .notes_shape import NotesShape

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ._protocols import NotesSlideProto


class NotesShapeCollection:
    """Collection facade for notes placeholder/shape entries."""

    def __init__(self, notes_slide: NotesSlideProto) -> None:
        """Initialize the collection for a notes slide."""
        super().__init__()
        self._notes_slide = notes_slide

    def _items(self) -> list[NotesShape]:
        return [
            NotesShape(self._notes_slide, payload)
            for payload in self._notes_slide.shape_payloads()
        ]

    def __len__(self) -> int:
        """Return the number of notes shapes."""
        return len(self._notes_slide.shape_payloads())

    def __iter__(self) -> Iterator[NotesShape]:
        """Iterate over notes shapes."""
        return iter(self._items())

    def __getitem__(self, index: int) -> NotesShape:
        """Return one notes shape by zero-based index."""
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
        """Return the first notes shape matching the provided filters."""
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
        with_text_frame: bool | None = None,
    ) -> list[NotesShape]:
        """Return all notes shapes matching the provided filters."""
        out: list[NotesShape] = []
        for item in self:
            if shape_type is not None and item.shape_type != shape_type:
                continue
            if (
                placeholder_type is not None
                and item.placeholder_type != placeholder_type
            ):
                continue
            if (
                with_text_frame is not None
                and (item.text_frame() is not None) != with_text_frame
            ):
                continue
            out.append(item)
        return out

    def by_id(self, shape_id: int) -> NotesShape | None:
        """Return one notes shape by shape ID, if present."""
        return self.get(shape_id=shape_id)
