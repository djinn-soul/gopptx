"""Core properties shared by slide proxy implementations."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .background import SlideBackground
from .notes.notes_slide import NotesSlide

if TYPE_CHECKING:
    from ..schemas import SlideMetadata
    from .contracts import SlidePresentationProtocol


class SlideBase:
    """Base class providing core slide properties."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]
        _metadata: SlideMetadata  # pyright: ignore[reportUninitializedInstanceVariable]

    @property
    def presentation(self) -> SlidePresentationProtocol:
        """Return the owning presentation proxy."""
        return self._presentation

    @property
    def index(self) -> int:
        """Return the zero-based slide index."""
        return self._presentation.slide_index_for_id(self.slide_id)

    @property
    def slide_id(self) -> int:
        """Return the unique internal slide ID."""
        return self._metadata["SlideID"]

    @property
    def title(self) -> str:
        """Return the slide title."""
        return self._metadata["Title"]

    @title.setter
    def title(self, value: str) -> None:
        self._presentation.set_slide_title(self.index, value)
        self._metadata["Title"] = value

    @property
    def notes(self) -> str:
        """Return the speaker notes."""
        return self._presentation.get_notes(self.index)

    @notes.setter
    def notes(self, value: str) -> None:
        self._presentation.set_notes(self.index, value)

    @property
    def notes_slide(self) -> NotesSlide | None:
        """Return a notes-slide proxy when notes exist."""
        if self.index < 0:
            return None
        notes_payload = self._presentation.get_notes_payload(self.index)
        if notes_payload.get("notes_slide") is None:
            return None
        return NotesSlide(self)

    @property
    def background(self) -> SlideBackground:
        """Return the slide background proxy."""
        return SlideBackground(self)
