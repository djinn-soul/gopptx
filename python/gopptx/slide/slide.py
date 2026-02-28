"""Slide proxy class for gopptx library."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING

from .placeholder_mixin import SlidePlaceholderMixin
from .table import Table

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from ..schemas import (
        Shape,
        ShapeProps,
        ShapeUpdate,
        SlideChartRef,
        SlideMetadata,
        TableCellInfo,
        TableInfo,
    )


class SlideShapeMixin:
    """Mixin providing shape manipulation methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def add_shape(
        self,
        shape_type: str,
        bounds: tuple[float, float, float, float],
        **kwargs: str | ShapeProps,
    ) -> int:
        """Add a shape to this slide."""
        return self._presentation.add_shape(self.index, shape_type, bounds, **kwargs)

    def add_image(self, path: str, bounds: tuple[float, float, float, float]) -> int:
        """Add an image to this slide."""
        return self._presentation.add_image(self.index, path, bounds)

    def remove_shape(self, shape_id: int) -> None:
        """Remove a shape from this slide."""
        self._presentation.remove_shape(self.index, shape_id)

    def move_shape_to_front(self, shape_id: int) -> None:
        """Move a shape to the front of the z-order."""
        self._presentation.move_shape_to_front(self.index, shape_id)

    def move_shape_to_back(self, shape_id: int) -> None:
        """Move a shape to the back of the z-order."""
        self._presentation.move_shape_to_back(self.index, shape_id)

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None:
        """Update shape properties."""
        self._presentation.update_shape(self.index, shape_id, updates)

    def list_shapes(self) -> list[Shape]:
        """List all shapes on this slide."""
        return self._presentation.list_shapes(self.index)


class SlideBase:
    """Base class providing core slide properties (index, title, notes)."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]
        _metadata: SlideMetadata  # pyright: ignore[reportUninitializedInstanceVariable]

    @property
    def index(self) -> int:
        """The zero-based index of this slide."""
        for s in self._presentation.slides_metadata:
            if s["SlideID"] == self.slide_id:
                return int(s["Index"])
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


class SlideChartMixin:
    """Mixin providing chart-related methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def list_charts(self) -> list[SlideChartRef]:
        """List all charts on this slide."""
        return self._presentation.list_slide_charts(self.index)

    def add_chart(
        self,
        chart_type: str,
        categories: list[str],
        values_or_series: list[float] | list[dict[str, str | list[float]]],
        **kwargs: str | tuple[float, float, float, float],
    ) -> int:
        """Add a chart to this slide.

        Returns:
            int: The created chart shape ID.
        """
        return self._presentation.add_chart(
            self.index, chart_type, categories, values_or_series, **kwargs
        )


class SlideTableMixin:
    """Mixin providing table manipulation methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def add_table(self, rows: int, cols: int, bounds: tuple[int, int, int, int]) -> int:
        """Add a table to this slide."""
        return self._presentation.add_table(self.index, rows, cols, bounds)

    def get_table(self, shape_id: int) -> TableInfo:
        """Get table information for a table shape."""
        return self._presentation.get_table(self.index, shape_id)

    def table(self, shape_id: int) -> Table:
        """Returns a Table object for the given shape_id, providing a Pythonic grid API."""
        return Table(self._presentation, self.index, shape_id)

    def set_table_flags(self, shape_id: int, flags: dict[str, bool]) -> None:
        """Set table style flags."""
        self._presentation.set_table_flags(self.index, shape_id, flags)

    def set_table_cell_text(self, shape_id: int, row: int, col: int, text: str) -> None:
        """Set the text of a table cell."""
        self._presentation.set_table_cell_text(self.index, shape_id, row, col, text)

    def get_table_cell(self, shape_id: int, row: int, col: int) -> TableCellInfo:
        """Get information about a table cell."""
        return self._presentation.get_table_cell(self.index, shape_id, row, col)

    def merge_table_cells(
        self, shape_id: int, cell_range: tuple[int, int, int, int]
    ) -> None:
        """Merge a range of table cells."""
        self._presentation.merge_table_cells(self.index, shape_id, cell_range)

    def split_table_cell(self, shape_id: int, row: int, col: int) -> None:
        """Split a merged table cell."""
        self._presentation.split_table_cell(self.index, shape_id, row, col)


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

    def __repr__(self) -> str:  # pyright: ignore[reportImplicitOverride]
        """Return a string representation of this slide."""
        return f"<Slide index={self.index} title='{self.title}'>"
