"""Slide proxy class for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import override

from .chart_model import Chart, ChartCollection
from .notes_slide import NotesSlide
from .placeholder_mixin import SlidePlaceholderMixin
from .shape_proxy import ShapeCollection, ShapeProxy
from .slide_mixins import SlideChartMixin, SlideShapeMixin, SlideTableMixin

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from ..schemas import SlideMetadata


class SlideBase:
    """Base class providing core slide properties (index, title, notes)."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]
        _metadata: SlideMetadata  # pyright: ignore[reportUninitializedInstanceVariable]

    @property
    def index(self) -> int:
        """The zero-based index of this slide."""
        for slide_meta in self._presentation.slides_metadata:
            if slide_meta["SlideID"] == self.slide_id:
                return int(slide_meta["Index"])
        return -1

    @property
    def slide_id(self) -> int:
        """The unique internal ID of this slide."""
        return self._metadata["SlideID"]

    @property
    def title(self) -> str:
        """The title text of this slide."""
        return self._metadata["Title"]

    @title.setter
    def title(self, value: str) -> None:
        self._presentation.set_slide_title(self.index, value)
        self._metadata["Title"] = value

    @property
    def notes(self) -> str:
        """Get the speaker notes for this slide."""
        return self._presentation.get_notes(self.index)

    @notes.setter
    def notes(self, value: str) -> None:
        self._presentation.set_notes(self.index, value)

    @property
    def notes_slide(self) -> NotesSlide | None:
        """Return a notes-slide proxy, or None when notes slide is absent."""
        if self.index < 0:
            return None
        notes_payload = self._presentation.get_notes_payload(self.index)
        if notes_payload.get("notes_slide") is None:
            return None
        return NotesSlide(self)


class Slide(
    SlideTableMixin,
    SlideChartMixin,
    SlidePlaceholderMixin,
    SlideBase,
    SlideShapeMixin,
):
    """Proxy object for a slide within a presentation."""

    def __init__(self, presentation: Presentation, metadata: SlideMetadata) -> None:
        """Initialize the slide proxy."""
        super().__init__()
        self._presentation = presentation
        self._metadata = metadata
        self._shapes_collection: ShapeCollection | None = None
        self._charts_collection: ChartCollection | None = None

    @property
    def shapes(self) -> ShapeCollection:
        """python-pptx-style shape collection."""
        if self._shapes_collection is None:
            self._shapes_collection = ShapeCollection(self)
        return self._shapes_collection

    def shape(self, shape_id: int) -> ShapeProxy:
        """Return a live shape proxy by ID."""
        return self.shapes.by_id(shape_id)

    @property
    def charts(self) -> ChartCollection:
        """python-pptx-style chart collection."""
        if self._charts_collection is None:
            self._charts_collection = ChartCollection(self)
        return self._charts_collection

    def chart(self, index: int) -> Chart:
        """Return a chart proxy by slide-local chart index."""
        return self.charts[index]

    def update(
        self,
        title: str | None = None,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> None:
        """Update slide properties."""
        self._presentation.update_slide(
            self.index, title=title, layout=layout, bullets=bullets
        )
        if title:
            self._metadata["Title"] = title

    def remove(self) -> None:
        """Remove this slide from the presentation."""
        self._presentation.remove_slide(self.index)

    def duplicate(self, insert_at: int | None = None) -> Slide:
        """Duplicate this slide."""
        new_idx = self._presentation.duplicate_slide(self.index, insert_at=insert_at)
        return self._presentation.slides[new_idx]

    @override
    def __repr__(self) -> str:
        """Return a string representation of this slide."""
        return f"<Slide index={self.index} title='{self.title}'>"
