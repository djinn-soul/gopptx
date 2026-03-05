"""Notes slide proxy for gopptx slide API."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import override

if TYPE_CHECKING:
    from .slide import SlideBase


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

    @override
    def __repr__(self) -> str:
        """Return debug representation for notes proxy."""
        return f"<NotesSlide slide_index={self._slide.index}>"
