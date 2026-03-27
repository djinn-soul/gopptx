"""Slides collection facade for presentation slide proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, overload

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...slide.slide import Slide


class _SlidesLookupProto(Protocol):
    def slide_index_for_id(self, slide_id: int) -> int:
        """Resolve a zero-based slide index for a stable slide ID."""
        ...


class Slides:
    """List-like collection with slide-ID lookup helper."""

    def __init__(self, owner: _SlidesLookupProto, items: list[Slide]) -> None:
        """Initialize a collection view bound to one presentation snapshot."""
        super().__init__()
        self._owner = owner
        self._items = items

    def __len__(self) -> int:
        """Return number of slides in this collection snapshot."""
        return len(self._items)

    def __iter__(self) -> Iterator[Slide]:
        """Iterate over slides in order."""
        return iter(self._items)

    @overload
    def __getitem__(self, index: int) -> Slide: ...

    @overload
    def __getitem__(self, index: slice) -> list[Slide]: ...

    def __getitem__(self, index: int | slice) -> Slide | list[Slide]:
        """Return one slide or a list of slides for a slice."""
        return self._items[index]

    def find_by_slide_id(self, slide_id: int) -> Slide | None:
        """Return the slide with ``slide_id`` or ``None`` when not found."""
        idx = self._owner.slide_index_for_id(slide_id)
        if idx < 0:
            return None
        if idx >= len(self._items):
            return None
        return self._items[idx]
