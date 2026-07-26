"""Slide chart mixin scoped to chart-domain operations."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Sequence

    from ...presentation.charts.chart_types import ChartType
    from ...schemas import SlideChartRef
    from ..contracts.presentation import SlidePresentationProtocol
    from .data import CategoryChartData, XyChartData


class SlideChartMixin:
    """Mixin providing chart-related methods for Slide objects."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

    def list_charts(self) -> list[SlideChartRef]:
        """List chart references on this slide."""
        return self._presentation.list_slide_charts(self.index)

    def add_chart(
        self,
        chart_type: ChartType,
        categories: Sequence[str] | CategoryChartData | XyChartData,
        values_or_series: Sequence[float] | Sequence[dict[str, object]] | None = None,
        *,
        title: str = "Chart",
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
    ) -> int:
        """Add a chart to this slide.

        Args:
            chart_type: A ChartType member, e.g. ChartType.COLUMN,
                ChartType.LINE, ChartType.PIE. Bare strings still work but are
                deprecated and warn.
            categories: List of category labels or ChartData builder.
            values_or_series: List of values or list of series dicts.
            title: Chart title.
            bounds: (x, y, width, height) in EMU. Use Inches() or Pt() to build
                it; raw numbers are EMU, and there are 914400 EMU to the inch.

        Returns:
            Shape ID of the created chart.

        Raises:
            ValueError: If chart_type is invalid or bounds is not a 4-tuple.

        Examples:
            from gopptx.presentation.charts import ChartType
            from gopptx.schemas import Inches

            chart_id = slide.add_chart(
                ChartType.PIE,
                ["A", "B", "C"],
                [30.0, 40.0, 30.0],
                title="Distribution",
                bounds=(Inches(1), Inches(1), Inches(4), Inches(3)),
            )
        """
        chart_id = self._presentation.add_chart(
            self.index,
            chart_type,
            categories,
            values_or_series,
            title=title,
            bounds=bounds,
        )
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return chart_id

    def add_combo_chart(
        self,
        categories: list[str],
        bar_series: list[dict[str, object]],
        line_series: list[dict[str, object]],
        *,
        title: str = "Chart",
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
    ) -> int:
        """Add a combo (bar + line) chart to this slide.

        Args:
            categories: List of category labels.
            bar_series: List of bar series dicts with "name" and "values" keys.
            line_series: List of line series dicts with "name" and "values" keys.
            title: Chart title.
            bounds: (x, y, width, height) in EMU.

        Returns:
            Shape ID of the created chart.

        Example:
            chart_id = slide.add_combo_chart(
                ["Q1", "Q2", "Q3"],
                bar_series=[{"name": "Revenue", "values": [100, 200, 150]}],
                line_series=[{"name": "Growth %", "values": [10, 15, 12]}],
                title="Sales Overview",
                bounds=(Inches(1), Inches(1), Inches(8), Inches(5)),
            )
        """
        chart_id = self._presentation.add_combo_chart(
            self.index,
            categories,
            bar_series,
            line_series,
            title=title,
            bounds=bounds,
        )
        invalidate = getattr(self, "_invalidate_shape_cache", None)
        if callable(invalidate):
            invalidate()
        invalidate_text = getattr(self, "_invalidate_text_state_cache", None)
        if callable(invalidate_text):
            invalidate_text()
        return chart_id
