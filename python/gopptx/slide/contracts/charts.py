"""Type-only bridge contract for slide proxy objects."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import Protocol

if TYPE_CHECKING:
    from collections.abc import Sequence

    from ...schemas import ChartState, SlideChartRef
    from ...slide.chart.data import CategoryChartData, XyChartData


class ChartOperationsProtocol(Protocol):
    """Chart management operations."""

    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: Sequence[str] | CategoryChartData | XyChartData,
        values_or_series: Sequence[float] | Sequence[dict[str, object]] | None = None,
        **kwargs: object,
    ) -> int:
        """Protocol member."""
        ...

    def add_combo_chart(
        self,
        slide_index: int,
        categories: list[str],
        bar_series: list[dict[str, object]],
        line_series: list[dict[str, object]],
        *,
        title: str = "Chart",
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
    ) -> int:
        """Protocol member."""
        ...

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]:
        """Protocol member."""
        ...

    def get_chart_state_by_index(
        self, slide_index: int, chart_index: int
    ) -> ChartState:
        """Protocol member."""
        ...

    def get_chart_state_by_rel_id(self, slide_index: int, rel_id: str) -> ChartState:
        """Protocol member."""
        ...

    def update_chart_data(
        self,
        slide_index: int,
        chart_selector: dict[str, object] | list[str],
        data: dict[str, object] | list[dict[str, object]],
    ) -> None:
        """Protocol member."""
        ...

    def update_chart_formatting(
        self,
        slide_index: int,
        chart_selector: dict[str, object],
        fmt: dict[str, object],
    ) -> None:
        """Protocol member."""
        ...
