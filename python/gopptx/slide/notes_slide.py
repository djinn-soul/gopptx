"""Notes slide proxy for gopptx slide API."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .slide import Slide


class NotesSlide:
    """Proxy for slide notes content."""

    def __init__(self, slide: Slide) -> None:
        """Initialize notes proxy bound to a slide."""
        self._slide = slide

    @property
    def text(self) -> str:
        """Get notes text."""
        return self._slide.notes

    @text.setter
    def text(self, value: str) -> None:
        """Set notes text."""
        self._slide.notes = value

    def __repr__(self) -> str:
        """Return debug representation for notes proxy."""
        return f"<NotesSlide slide_index={self._slide.index}>"
