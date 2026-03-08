"""Slide proxy class for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import override

from .chart_model import Chart, ChartCollection
from .notes_slide import NotesSlide
from .placeholder_mixin import SlidePlaceholderMixin
from .shape_proxy import ShapeCollection, ShapeProxy
from .slide_mixins import SlideChartMixin, SlideShapeMixin, SlideTableMixin
from .slide_shape_batch_mixin import SlideShapeBatchMixin
from .slide_text_mixin import SlideTextMixin

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from ..schemas import Shape, SlideMetadata


class SlideBase:
    """Base class providing core slide properties (index, title, notes)."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]
        _metadata: SlideMetadata  # pyright: ignore[reportUninitializedInstanceVariable]

    @property
    def index(self) -> int:
        """The zero-based index of this slide."""
        return self._presentation.slide_index_for_id(self.slide_id)

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
    SlideTextMixin,
    SlideShapeBatchMixin,
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
        self._shape_records_cache: list[Shape] | None = None
        self._shape_record_map: dict[int, Shape] | None = None
        self._shape_text_state_cache: dict[int, dict[str, object]] | None = None

    def _shape_records(self) -> list[Shape]:
        if self._shape_records_cache is None:
            self._shape_records_cache = self._presentation.list_shapes(self.index)
            self._shape_record_map = None
        return self._shape_records_cache

    def _shape_record_by_id(self, shape_id: int) -> Shape:
        if self._shape_record_map is None:
            record_map: dict[int, Shape] = {}
            for shape in self._shape_records():
                raw_id = shape.get("ID", shape.get("id"))
                if raw_id is not None:
                    record_map[int(str(raw_id))] = shape
            self._shape_record_map = record_map
        return self._shape_record_map[shape_id]

    def _invalidate_shape_cache(self) -> None:
        self._shape_records_cache = None
        self._shape_record_map = None

    def _shape_text_states(self) -> dict[int, dict[str, object]]:
        if self._shape_text_state_cache is None:
            state_map: dict[int, dict[str, object]] = {}
            for state in self._presentation.get_slide_text_states(self.index):
                raw_id = state.get("shape_id", state.get("ShapeID"))
                if raw_id is None:
                    continue
                state_map[int(str(raw_id))] = dict(state)
            self._shape_text_state_cache = state_map
        return self._shape_text_state_cache

    def _invalidate_text_state_cache(self) -> None:
        self._shape_text_state_cache = None

    @property
    def shapes(self) -> ShapeCollection:
        """python-pptx-style shape collection."""
        if self._shapes_collection is None:
            self._shapes_collection = ShapeCollection(self)
        return self._shapes_collection

    @override
    def list_shapes(self) -> list[Shape]:
        """List slide shapes using a slide-local snapshot cache."""
        return self._shape_records()

    @override
    def get_shape_text_state(self, shape_id: int) -> dict[str, object]:
        """Get shape text state using slide-local bulk-prefetched cache."""
        state = self._shape_text_states().get(shape_id)
        if state is None:
            state = self._presentation.get_shape_text_state(self.index, shape_id)
            self._shape_text_states()[shape_id] = dict(state)
        return dict(state)

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
