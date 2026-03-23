"""Slide chart mixin scoped to chart-domain operations."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from collections.abc import Sequence

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
        chart_type: str,
        categories: Sequence[str] | CategoryChartData | XyChartData,
        values_or_series: Sequence[float] | Sequence[dict[str, object]] | None = None,
        **kwargs: str | tuple[float, float, float, float],
    ) -> int:
        """Add a chart to this slide.

        Args:
            chart_type: Chart type constant (ChartType enum value).
                Must be one of: ChartType.COLUMN, ChartType.LINE,
                ChartType.PIE, ChartType.SCATTER, ChartType.AREA,
                ChartType.RADAR, ChartType.BUBBLE, etc.
            categories: List of category labels or ChartData builder.
            values_or_series: List of values or list of series dicts.
            **kwargs: Additional options (bounds, title, etc.)

        Returns:
            Shape ID of the created chart.

        Examples:
            # Using ChartType constants (required)
            from gopptx.presentation.charts import ChartType

            chart_id = slide.add_chart(
                ChartType.COLUMN,
                ["Q1", "Q2", "Q3"],
                [100, 200, 150],
                title="Sales",
                bounds=(100, 100, 400, 300),
            )

            # Other chart types
            chart_id = slide.add_chart(
                ChartType.PIE,
                ["A", "B", "C"],
                [30.0, 40.0, 30.0],
                title="Distribution",
                bounds=(100, 100, 400, 300),
            )
        """
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
