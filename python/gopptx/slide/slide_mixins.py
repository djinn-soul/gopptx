"""Slide mixins for shape, chart, and table APIs."""
# ruff: noqa: D102

from __future__ import annotations

from typing import TYPE_CHECKING

from .table import Table

if TYPE_CHECKING:
    from ..presentation.presentation import Presentation
    from ..schemas import (
        ImageMetadata,
        Shape,
        ShapeProps,
        ShapeUpdate,
        SlideChartRef,
        TableCellInfo,
        TableInfo,
    )
    from .chart_data import CategoryChartData, XyChartData
    from .freeform_builder import FreeformBuilder


class SlideShapeMixin:
    """Mixin providing shape manipulation methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

        def _invalidate_text_state_cache_if_present(self) -> None: ...

    def _invalidate_shape_cache_if_present(self) -> None:
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()

    def add_shape(
        self,
        shape_type: str,
        bounds: tuple[float, float, float, float],
        **kwargs: str | ShapeProps,
    ) -> int:
        shape_id = self._presentation.add_shape(
            self.index, shape_type, bounds, **kwargs
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_textbox(
        self,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        text: str = "",
        **kwargs: str | ShapeProps,
    ) -> int:
        shape_id = self._presentation.add_textbox(
            self.index, left, top, width, height, text=text, **kwargs
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_connector(
        self,
        connector_type: str,
        begin_x: float,
        begin_y: float,
        end_x: float,
        end_y: float,
        **kwargs: str | ShapeProps,
    ) -> int:
        shape_id = self._presentation.add_connector(
            self.index,
            connector_type,
            begin_x,
            begin_y,
            end_x,
            end_y,
            **kwargs,
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_group_shape(self, shapes: list[int] | None = None) -> int:
        shape_id = self._presentation.add_group_shape(self.index, shapes=shapes)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def build_freeform(
        self,
        start_x: float = 0,
        start_y: float = 0,
        scale: tuple[float, float] | float = 1.0,
    ) -> FreeformBuilder:
        return self._presentation.build_freeform(
            self.index, start_x=start_x, start_y=start_y, scale=scale
        )

    def add_image(
        self,
        path: str | None,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int:
        shape_id = self._presentation.add_image(self.index, path, bounds, **kwargs)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def get_image_metadata(self, shape_id: int) -> ImageMetadata:
        return self._presentation.get_image_metadata(self.index, shape_id)

    def add_video(
        self,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        shape_id = self._presentation.add_video(
            self.index,
            source,
            bounds,
            name=name,
            poster_frame=poster_frame,
            mime_type=mime_type,
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_audio(
        self,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int:
        shape_id = self._presentation.add_audio(
            self.index,
            source,
            bounds,
            name=name,
            poster_frame=poster_frame,
            mime_type=mime_type,
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def add_ole_object(
        self,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        prog_id: str | None = None,
        icon: str | bytes | None = None,
    ) -> int:
        shape_id = self._presentation.add_ole_object(
            self.index, source, bounds, name=name, prog_id=prog_id, icon=icon
        )
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def remove_shape(self, shape_id: int) -> None:
        self._presentation.remove_shape(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def group_shapes(self, shape_ids: list[int]) -> int:
        shape_id = self._presentation.group_shapes(self.index, shape_ids)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return shape_id

    def ungroup_shapes(self, shape_id: int) -> int:
        group_id = self._presentation.ungroup_shapes(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return group_id

    def move_shape_to_front(self, shape_id: int) -> None:
        self._presentation.move_shape_to_front(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def move_shape_to_back(self, shape_id: int) -> None:
        self._presentation.move_shape_to_back(self.index, shape_id)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None:
        self._presentation.update_shape(self.index, shape_id, updates)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()

    def list_shapes(self) -> list[Shape]:
        return self._presentation.list_shapes(self.index)


class SlideChartMixin:
    """Mixin providing chart-related methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def list_charts(self) -> list[SlideChartRef]:
        return self._presentation.list_slide_charts(self.index)

    def add_chart(
        self,
        chart_type: str,
        categories: list[str] | CategoryChartData | XyChartData,
        values_or_series: list[float]
        | list[dict[str, str | list[float]]]
        | None = None,
        **kwargs: str | tuple[float, float, float, float],
    ) -> int:
        chart_id = self._presentation.add_chart(
            self.index, chart_type, categories, values_or_series, **kwargs
        )
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return chart_id


class SlideTableMixin:
    """Mixin providing table manipulation methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: Presentation  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int: ...

    def add_table(self, rows: int, cols: int, bounds: tuple[int, int, int, int]) -> int:
        shape_id = self._presentation.add_table(self.index, rows, cols, bounds)
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return shape_id

    def get_table(self, shape_id: int) -> TableInfo:
        return self._presentation.get_table(self.index, shape_id)

    def table(self, shape_id: int) -> Table:
        return Table(self._presentation, self.index, shape_id)

    def set_table_flags(self, shape_id: int, flags: dict[str, bool]) -> None:
        self._presentation.set_table_flags(self.index, shape_id, flags)

    def set_table_cell_text(self, shape_id: int, row: int, col: int, text: str) -> None:
        self._presentation.set_table_cell_text(self.index, shape_id, row, col, text)

    def get_table_cell(self, shape_id: int, row: int, col: int) -> TableCellInfo:
        return self._presentation.get_table_cell(self.index, shape_id, row, col)

    def merge_table_cells(
        self, shape_id: int, cell_range: tuple[int, int, int, int]
    ) -> None:
        self._presentation.merge_table_cells(self.index, shape_id, cell_range)

    def split_table_cell(self, shape_id: int, row: int, col: int) -> None:
        self._presentation.split_table_cell(self.index, shape_id, row, col)

    def set_table_style(self, shape_id: int, style_guid: str) -> None:
        self._presentation.set_table_style(self.index, shape_id, style_guid)

    def set_table_row_height(self, shape_id: int, row: int, height: int) -> None:
        self._presentation.set_table_row_height(self.index, shape_id, row, height)

    def set_table_column_width(self, shape_id: int, col: int, width: int) -> None:
        self._presentation.set_table_column_width(self.index, shape_id, col, width)
