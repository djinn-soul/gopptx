"""Chart mutation mixin for the Presentation API."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from .chart_mixin_updates import PresentationChartUpdatesMixin
from .chart_types import ChartType
from .state_mixin import PresentationChartStateMixin

if TYPE_CHECKING:
    from collections.abc import Sequence

    from ...slide.chart.data import CategoryChartData, XyChartData


class PresentationChartMixin(
    PresentationChartUpdatesMixin, PresentationChartStateMixin
):
    """Mixin providing chart creation and manipulation methods."""

    _BOUNDS_LEN = 4

    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: Sequence[str] | CategoryChartData | XyChartData,
        values_or_series: Sequence[float] | Sequence[dict[str, object]] | None = None,
        **kwargs: str | tuple[float, float, float, float],
    ) -> int:
        """Add a chart to a slide.

        Args:
            slide_index: Zero-based slide index.
            chart_type: Chart type constant (ChartType enum value).
                Must be one of: ChartType.COLUMN, ChartType.LINE,
                ChartType.PIE, ChartType.SCATTER, ChartType.AREA,
                ChartType.RADAR, ChartType.BUBBLE, etc.
            categories: List of category labels or ChartData builder.
            values_or_series: List of values or list of series dicts.
            **kwargs: Additional options including:
                - bounds: (x, y, width, height) tuple
                - title: Chart title string

        Returns:
            Shape ID of the created chart.

        Raises:
            ValueError: If chart_type is invalid or bounds format is wrong.

        Examples:
            # Using ChartType enum constants (required)
            from gopptx.presentation.charts import ChartType

            chart_id = prs.add_chart(
                0,
                ChartType.COLUMN,
                ["Q1", "Q2", "Q3"],
                [100, 200, 150],
                title="Sales",
                bounds=(100, 100, 400, 300),
            )

            # Other chart types
            chart_id = prs.add_chart(
                0,
                ChartType.PIE,
                ["Product A", "Product B", "Product C"],
                [25.0, 35.0, 40.0],
                title="Sales by Product",
                bounds=(100, 100, 400, 300),
            )
        """
        # Validate chart type - only constant string values accepted
        if not chart_type:
            raise ValueError(
                f"chart_type must be a ChartType constant (e.g., ChartType.COLUMN). Got: {chart_type!r}. Use ChartType.get_all() to see available options."
            )

        # Check if it's a valid chart type value
        valid_types = set(ChartType.get_all().values())
        if chart_type not in valid_types:
            valid_values = ", ".join(sorted(valid_types))
            raise ValueError(
                f"Invalid chart_type {chart_type!r}. Use ChartType constants like ChartType.COLUMN, ChartType.LINE, ChartType.PIE. Available raw values: {valid_values}"
            )

        if hasattr(categories, "to_add_chart_args"):
            categories, values_or_series = cast(
                "CategoryChartData | XyChartData", categories
            ).to_add_chart_args()
        if values_or_series is None:
            values_or_series = []
        bounds = cast(
            "tuple[float, float, float, float]", kwargs.get("bounds", (0, 0, 0, 0))
        )
        if len(bounds) != self._BOUNDS_LEN:
            raise ValueError("bounds must be a tuple of (x, y, w, h)")
        x, y, w, h = bounds
        title = kwargs.get("title", "Chart")

        values: list[float]
        if values_or_series and isinstance(values_or_series[0], dict):
            series_items = cast("list[dict[str, str | list[float]]]", values_or_series)
            first = series_items[0]
            values = cast("list[float]", first.get("values", []))
            title = str(first.get("name", title))
        else:
            values = cast("list[float]", values_or_series)
        result = self.execute(
            ops.OP_ADD_CHART,
            {
                "slide_index": slide_index,
                "chart_type": chart_type,
                "title": title,
                "categories": categories,
                "values": values,
                "x": x,
                "y": y,
                "w": w,
                "h": h,
            },
        )
        return int(cast("int", result.get("shape_id") or result.get("chart_id", 0)))

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
        """Add a combo (bar + line) chart to a slide.

        Args:
            slide_index: Zero-based slide index.
            categories: List of category labels.
            bar_series: List of bar series dicts with "name" and "values" keys.
            line_series: List of line series dicts with "name" and "values" keys.
            title: Chart title.
            bounds: (x, y, width, height) in EMU.

        Returns:
            Shape ID of the created chart.

        Example:
            chart_id = prs.add_combo_chart(
                0,
                ["Q1", "Q2", "Q3"],
                bar_series=[{"name": "Revenue", "values": [100, 200, 150]}],
                line_series=[{"name": "Growth %", "values": [10, 15, 12]}],
                title="Sales Overview",
                bounds=(Inches(1), Inches(1), Inches(8), Inches(5)),
            )
        """
        x, y, w, h = bounds
        result = self.execute(
            ops.OP_ADD_CHART,
            {
                "slide_index": slide_index,
                "chart_type": ChartType.COMBO,
                "title": title,
                "categories": categories,
                "values": [],
                "bar_series": [dict(s) for s in bar_series],
                "line_series": [dict(s) for s in line_series],
                "x": x,
                "y": y,
                "w": w,
                "h": h,
            },
        )
        return int(cast("int", result.get("shape_id") or result.get("chart_id", 0)))
