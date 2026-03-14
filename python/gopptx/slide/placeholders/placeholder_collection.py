"""Placeholder collection proxy for slide placeholders."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .placeholder import Placeholder, create_placeholder

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ..slide import Slide


class PlaceholderCollection:
    """python-pptx-style placeholder collection.

    Supports:
    - iteration (`for ph in slide.placeholders`)
    - length (`len(slide.placeholders)`)
    - idx lookup (`slide.placeholders[idx]`)
    """

    def __init__(self, slide: Slide) -> None:
        """Bind collection to a single slide proxy."""
        super().__init__()
        self._slide = slide

    def _items(self) -> list[Placeholder]:
        ph_data = self._slide.list_placeholders()
        placeholders = [
            create_placeholder(
                self._slide,
                ph["index"],  # type: ignore[dict-item]
                ph.get("type", ""),  # type: ignore[dict-item]
                ph.get("name", ""),  # type: ignore[dict-item]
            )
            for ph in ph_data
        ]
        placeholders.sort(key=lambda p: p.idx)
        return placeholders

    def __iter__(self) -> Iterator[Placeholder]:
        """Iterate placeholders ordered by idx."""
        return iter(self._items())

    def __len__(self) -> int:
        """Return number of placeholders."""
        return len(self._items())

    def __getitem__(self, idx: int) -> Placeholder:
        """Return placeholder by idx."""
        for placeholder in self._items():
            if placeholder.idx == idx:
                return placeholder
        raise KeyError(f"no placeholder on this slide with idx == {idx}")

    def get(self, idx: int, default: Placeholder | None = None) -> Placeholder | None:
        """Return placeholder by idx, or default."""
        try:
            return self[idx]
        except KeyError:
            return default
