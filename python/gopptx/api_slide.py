from __future__ import annotations

from typing import TYPE_CHECKING, Optional

from .types import (
    Shape,
    ShapeProps,
    ShapeUpdate,
    SlideChartRef,
    SlideMetadata,
    TableCellInfo,
    TableInfo,
)

if TYPE_CHECKING:
    from .api_presentation import Presentation


class Slide:
    """Proxy object for a slide within a presentation."""

    def __init__(self, presentation: Presentation, metadata: SlideMetadata):
        self._presentation = presentation
        self._metadata = metadata

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
        return self._presentation.get_notes(self.index)

    @notes.setter
    def notes(self, value: str) -> None:
        self._presentation.set_notes(self.index, value)

    def add_shape(
        self,
        shape_type: str,
        x: float,
        y: float,
        w: float,
        h: float,
        text: Optional[str] = None,
        properties: Optional[ShapeProps] = None,
    ) -> int:
        return self._presentation.add_shape(
            self.index, shape_type, x, y, w, h, text=text, properties=properties
        )

    def add_image(self, path: str, x: float, y: float, w: float, h: float) -> int:
        return self._presentation.add_image(self.index, path, x, y, w, h)

    def add_table(
        self,
        rows: int,
        cols: int,
        x: int,
        y: int,
        cx: int,
        cy: int,
    ) -> int:
        return self._presentation.add_table(self.index, rows, cols, x, y, cx, cy)

    def get_table(self, shape_id: int) -> TableInfo:
        return self._presentation.get_table(self.index, shape_id)

    def set_table_flags(self, shape_id: int, flags: dict[str, bool]) -> None:
        self._presentation.set_table_flags(self.index, shape_id, flags)

    def set_table_cell_text(self, shape_id: int, row: int, col: int, text: str) -> None:
        self._presentation.set_table_cell_text(self.index, shape_id, row, col, text)

    def get_table_cell(self, shape_id: int, row: int, col: int) -> TableCellInfo:
        return self._presentation.get_table_cell(self.index, shape_id, row, col)

    def merge_table_cells(
        self, shape_id: int, row1: int, col1: int, row2: int, col2: int
    ) -> None:
        self._presentation.merge_table_cells(self.index, shape_id, row1, col1, row2, col2)

    def split_table_cell(self, shape_id: int, row: int, col: int) -> None:
        self._presentation.split_table_cell(self.index, shape_id, row, col)

    def remove_shape(self, shape_id: int) -> None:
        self._presentation.remove_shape(self.index, shape_id)

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None:
        self._presentation.update_shape(self.index, shape_id, updates)

    def list_shapes(self) -> list[Shape]:
        return self._presentation.list_shapes(self.index)

    def list_charts(self) -> list[SlideChartRef]:
        return self._presentation.list_slide_charts(self.index)

    def add_chart(
        self,
        chart_type: str,
        categories: list[str],
        values: list[float],
        title: str = "Chart",
        x: int = 0,
        y: int = 0,
        w: int = 0,
        h: int = 0,
    ) -> None:
        self._presentation.add_chart(
            self.index, chart_type, categories, values, title=title, x=x, y=y, w=w, h=h
        )

    def update(
        self,
        title: Optional[str] = None,
        layout: Optional[str] = None,
        bullets: Optional[list[str]] = None,
    ) -> None:
        self._presentation.update_slide(
            self.index, title=title, layout=layout, bullets=bullets
        )
        if title:
            self._metadata["Title"] = title

    def remove(self) -> None:
        self._presentation.remove_slide(self.index)

    def duplicate(self, insert_at: Optional[int] = None) -> Slide:
        new_idx = self._presentation.duplicate_slide(self.index, insert_at=insert_at)
        return self._presentation.slides[new_idx]

    def __repr__(self) -> str:
        return f"<Slide index={self.index} title='{self.title}'>"
