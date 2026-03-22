"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from collections.abc import Mapping, Sequence
    from ..constants import ConnectorType, ShapeType
    from ..presentation.slides.master import SlideMaster
    from ..schemas import (
        Shape,
        ImageMetadata,
        ShapeUpdate,
        SlideChartRef,
        TableCellInfo,
        TableInfo,
        TextRun,
        ChartState,
    )
    from .shapes.freeform_builder import FreeformBuilder
    from ..slide.chart.data import CategoryChartData, XyChartData
    from .slide import Slide


class SlidePresentationProtocol(Protocol):
    """Operations exposed by the presentation bridge to slide proxies."""

    _batch_active: bool

    def execute(
        self, op: str, payload: dict[str, object] | None = None
    ) -> dict[str, object]: ...

    def invalidate_cache(self) -> None: ...

    def slide_index_for_id(self, slide_id: int) -> int: ...

    @property
    def slides(self) -> list[Slide]: ...

    @property
    def slide_masters(self) -> list[SlideMaster]: ...

    def list_shapes(self, slide_index: int) -> list[Shape]: ...

    def get_slide_text_states(self, slide_index: int) -> list[dict[str, object]]: ...

    def get_shape_text_state(
        self, slide_index: int, shape_id: int
    ) -> dict[str, object]: ...

    def get_shape_runs(self, slide_index: int, shape_id: int) -> list[TextRun]: ...

    def set_shape_runs(
        self, slide_index: int, shape_id: int, runs: list[TextRun]
    ) -> None: ...

    def update_shape_run_text(
        self, slide_index: int, shape_id: int, run_index: int, text: str
    ) -> None: ...

    def update_slide_run_texts(
        self, slide_index: int, updates: list[dict[str, object]]
    ) -> None: ...

    def append_shape_run(
        self, slide_index: int, shape_id: int, run: TextRun
    ) -> None: ...

    def add_shape(
        self,
        slide_index: int,
        shape_type: ShapeType,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int: ...

    def add_textbox(
        self,
        slide_index: int,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        text: str = "",
        **kwargs: object,
    ) -> int: ...

    def add_textboxes(
        self, slide_index: int, textboxes: Sequence[Mapping[str, object]]
    ) -> list[int]: ...

    def add_connector(
        self,
        slide_index: int,
        connector_type: ConnectorType,
        begin_x: float,
        begin_y: float,
        end_x: float,
        end_y: float,
        **kwargs: object,
    ) -> int: ...

    def add_connectors(
        self, slide_index: int, connectors: Sequence[Mapping[str, object]]
    ) -> list[int]: ...

    def add_group_shape(self, slide_index: int, shapes: list[int] | None = None) -> int: ...

    def build_freeform(
        self,
        slide_index: int,
        *,
        start_x: float = 0,
        start_y: float = 0,
        scale: tuple[float, float] | float = 1.0,
    ) -> FreeformBuilder: ...

    def add_image(
        self,
        slide_index: int,
        path: str | None,
        bounds: tuple[float, float, float, float],
        **kwargs: object,
    ) -> int: ...

    def get_image_metadata(self, slide_index: int, shape_id: int) -> ImageMetadata: ...

    def add_video(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int: ...

    def add_audio(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        poster_frame: str | bytes | None = None,
        mime_type: str | None = None,
    ) -> int: ...

    def add_ole_object(
        self,
        slide_index: int,
        source: str | bytes,
        bounds: tuple[float, float, float, float],
        *,
        name: str | None = None,
        prog_id: str | None = None,
        icon: str | bytes | None = None,
    ) -> int: ...

    def remove_shape(self, slide_index: int, shape_id: int) -> None: ...

    def group_shapes(self, slide_index: int, shape_ids: list[int]) -> int: ...

    def ungroup_shapes(self, slide_index: int, shape_id: int) -> int: ...

    def move_shape_to_front(self, slide_index: int, shape_id: int) -> None: ...

    def move_shape_to_back(self, slide_index: int, shape_id: int) -> None: ...

    def update_shape(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None: ...

    def add_mermaid(self, slide_index: int, diagram: str, *, theme: str = "") -> tuple[int, int]: ...

    def list_placeholders(self, slide_index: int) -> list[dict[str, object]]: ...

    def get_slide_layout_ref(self, slide_index: int) -> tuple[str, str]: ...

    def set_placeholder_content(
        self,
        slide_index: int,
        ph_index: int,
        ph_type: str = "",
        **kwargs: object,
    ) -> None: ...

    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: Sequence[str] | CategoryChartData | XyChartData,
        values_or_series: Sequence[float]
        | Sequence[dict[str, object]]
        | None = None,
        **kwargs: object,
    ) -> int: ...

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]: ...

    def get_chart_state_by_index(self, slide_index: int, chart_index: int) -> ChartState: ...

    def get_chart_state_by_rel_id(self, slide_index: int, rel_id: str) -> ChartState: ...

    def update_chart_data(
        self,
        slide_index: int,
        chart_selector: dict[str, object] | list[str],
        data: dict[str, object] | list[dict[str, object]],
    ) -> None: ...

    def update_chart_formatting(
        self,
        slide_index: int,
        chart_selector: dict[str, object],
        fmt: dict[str, object],
    ) -> None: ...

    def add_table(
        self,
        slide: int | None = None,
        slide_index: int | None = None,
        rows: int | None = None,
        cols: int | None = None,
        **kwargs: object,
    ) -> int: ...

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo: ...

    def set_table_flags(self, slide_index: int, shape_id: int, flags: dict[str, bool]) -> None: ...

    def set_table_cell_text(
        self, slide_index: int, shape_id: int, row: int, col: int, text: str
    ) -> None: ...

    def get_table_cell(
        self, slide_index: int, shape_id: int, row: int, col: int
    ) -> TableCellInfo: ...

    def set_table_style(self, slide_index: int, shape_id: int, style_guid: str) -> None: ...

    def merge_table_cells(
        self, slide_index: int, shape_id: int, cell_range: tuple[int, int, int, int]
    ) -> None: ...

    def split_table_cell(self, slide_index: int, shape_id: int, row: int, col: int) -> None: ...

    def set_table_row_height(
        self, slide_index: int, shape_id: int, row: int, height: int
    ) -> None: ...

    def set_table_column_width(
        self, slide_index: int, shape_id: int, col: int, width: int
    ) -> None: ...

    def get_notes(self, slide_index: int) -> str: ...

    def get_notes_payload(self, slide_index: int) -> dict[str, object]: ...

    def set_notes(self, slide_index: int, text: str) -> None: ...

    def set_notes_shape_text(self, slide_index: int, shape_id: int, text: str) -> None: ...

    def set_notes_shape_props(
        self, slide_index: int, shape_id: int, updates: ShapeUpdate
    ) -> None: ...

    def remove_slide(self, index: int) -> None: ...

    def update_slide(
        self,
        index: int,
        *,
        title: str | None = None,
        layout: str | None = None,
        bullets: list[str] | None = None,
    ) -> None: ...

    def set_slide_title(self, index: int, title: str) -> None: ...

    def duplicate_slide(self, index: int, insert_at: int | None = None) -> int: ...

    def flush_pending_textbox_adds(self) -> None: ...

    def has_pending_textbox_adds(self, slide_index: int) -> bool: ...

    def queue_textbox_add(self, slide_index: int, payload: dict[str, object]) -> int: ...

    def flush_pending_slide_run_text_updates(self, slide_index: int) -> None: ...

    def has_pending_slide_run_text_updates(self, slide_index: int) -> bool: ...

    def queue_shape_runs_replace(
        self, slide_index: int, shape_id: int, runs: list[dict[str, object]]
    ) -> None: ...

    def queue_shape_run_text_update(
        self, slide_index: int, shape_id: int, run_index: int, text: str
    ) -> None: ...

    def flush_pending_shape_runs_replacements(self, slide_index: int, shape_id: int) -> None: ...

    def has_pending_shape_runs_replace(self, slide_index: int, shape_id: int) -> bool: ...
