"""Chart mutation mixin for the Presentation API."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from .chart_mixin_updates import PresentationChartUpdatesMixin
from .chart_type_coerce import coerce_chart_type
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
        chart_type: ChartType,
        categories: Sequence[str] | CategoryChartData | XyChartData,
        values_or_series: Sequence[float] | Sequence[dict[str, object]] | None = None,
        *,
        title: str = "Chart",
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
    ) -> int:
        """Add a chart to a slide.

        Prefer ``slide.add_chart(...)`` on the Slide object, which takes the same
        arguments without the slide_index and so cannot address the wrong slide.

        Args:
            slide_index: Zero-based slide index. Note that Presentation.new()
                already creates slide 0, so the first slide you add is index 1.
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

            chart_id = prs.add_chart(
                1,
                ChartType.COLUMN,
                ["Q1", "Q2", "Q3"],
                [100, 200, 150],
                title="Sales",
                bounds=(Inches(1), Inches(1), Inches(4), Inches(3)),
            )
        """
        chart_type_value = coerce_chart_type(chart_type)

        if hasattr(categories, "to_add_chart_args"):
            categories, values_or_series = cast(
                "CategoryChartData | XyChartData", categories
            ).to_add_chart_args()
        if values_or_series is None:
            values_or_series = []
        if len(bounds) != self._BOUNDS_LEN:
            raise ValueError("bounds must be a tuple of (x, y, w, h)")
        x, y, w, h = bounds

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
                "chart_type": chart_type_value,
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
                A series may also carry ``"secondary_axis": True`` to draw that
                series against its own value axis on the right instead of
                sharing the bar series' scale, which is what keeps a growth
                percentage readable next to revenue. Series without the flag
                stay on the primary axis, so a mixed set is drawn as two line
                plots. ``"secondary_axis_title"`` labels that axis and enables
                it on its own; when no series is marked, every line series moves
                to it.
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
        line_payload = [dict(s) for s in line_series]
        secondary_title = next(
            (
                str(s["secondary_axis_title"])
                for s in line_payload
                if s.get("secondary_axis_title")
            ),
            "",
        )
        secondary_axis = bool(secondary_title) or any(
            bool(s.get("secondary_axis")) for s in line_payload
        )
        result = self.execute(
            ops.OP_ADD_CHART,
            {
                "slide_index": slide_index,
                "chart_type": ChartType.COMBO,
                "title": title,
                "categories": categories,
                "values": [],
                "bar_series": [dict(s) for s in bar_series],
                "line_series": line_payload,
                "x": x,
                "y": y,
                "w": w,
                "h": h,
                "secondary_axis": secondary_axis,
                "secondary_value_axis_title": secondary_title,
            },
        )
        return int(cast("int", result.get("shape_id") or result.get("chart_id", 0)))
