"""Notes slide proxy for gopptx slide API."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

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

    @property
    def placeholders(self) -> list[dict[str, object]]:
        """Placeholder metadata from notes slide."""
        payload = self._slide._presentation.get_notes_payload(self._slide.index)  # noqa: SLF001
        raw = payload.get("notes_placeholders")
        return cast("list[dict[str, object]]", raw if isinstance(raw, list) else [])

    @property
    def shapes(self) -> list[dict[str, object]]:
        """Compatibility alias; notes currently expose placeholder-focused metadata."""
        return self.placeholders

    @override
    def __repr__(self) -> str:
        """Return debug representation for notes proxy."""
        return f"<NotesSlide slide_index={self._slide.index}>"
