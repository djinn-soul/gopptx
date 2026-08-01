"""Slides collection facade for presentation slide proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast, overload

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

    def remove(self, slide: Slide | int) -> None:
        """Remove slide proxy or index from the presentation (Issue #67)."""
        if isinstance(slide, int):
            index = slide
        elif hasattr(slide, "index"):
            index = slide.index
        else:
            raise TypeError("slide must be a Slide proxy or integer index")
        remove_func = getattr(self._owner, "remove_slide", None)
        if callable(remove_func):
            remove_func(index)
        else:
            raise AttributeError("owner does not support remove_slide")

    def delete(self, slide: Slide | int) -> None:
        """Alias for remove (Issue #67)."""
        self.remove(slide)

    def __delitem__(self, index: int) -> None:
        """Delete slide at index (Issue #67)."""
        self.remove(index)

    def move(self, slide_or_index: Slide | int, to_index: int) -> None:
        """Move slide from current position to to_index (Issue #68)."""
        if isinstance(slide_or_index, int):
            from_idx = slide_or_index
        elif hasattr(slide_or_index, "index"):
            from_idx = slide_or_index.index
        else:
            raise TypeError("slide must be a Slide proxy or integer index")
        move_func = getattr(self._owner, "move_slide", None)
        if callable(move_func):
            move_func(from_idx, to_index)
        else:
            raise AttributeError("owner does not support move_slide")

    def add_slide(
        self, layout: object = None, index: int | None = None, title: str = ""
    ) -> Slide:
        """Add or insert a slide with optional index (Issue #194)."""
        add_func = getattr(self._owner, "add_slide", None)
        if not callable(add_func):
            raise AttributeError("owner does not support add_slide")

        # A SlideLayout object names a concrete layout part. Routing its display
        # name ("Title Slide", or any custom layout name) through add_slide would
        # hit SlideLayoutType.validate, which only accepts the four built-in
        # tokens. Create the slide first, then rebind it to that exact part.
        layout_part = getattr(layout, "part_name", None) if layout is not None else None
        layout_name = (
            None
            if layout is None or layout_part is not None
            else getattr(layout, "name", str(layout))
        )
        slide_obj = cast("Slide", add_func(title=title, layout=layout_name))
        slide_idx = getattr(slide_obj, "index", -1)
        if layout_part and isinstance(slide_idx, int) and slide_idx >= 0:
            slide_obj.rebind_layout(cast("str", layout_part))
        if index is not None and isinstance(slide_idx, int) and index != slide_idx:
            self.move(slide_idx, index)
            slides_list = cast("list[Slide]", getattr(self._owner, "slides", []))
            return slides_list[index]
        return slide_obj

    def insert_slide(self, index: int, layout: object = None, title: str = "") -> Slide:
        """Insert a slide at index (Issue #194)."""
        return self.add_slide(layout=layout, index=index, title=title)
