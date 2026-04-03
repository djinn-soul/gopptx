"""Named chart type constants for PowerPoint chart operations."""

from __future__ import annotations


class ChartType:
    """Named chart type constants for all supported PowerPoint chart types.

    These constants provide a discoverable, type-safe way to specify chart types.
    ONLY the ChartType constants are accepted - string values are not supported.

    Examples:
        # Using enum-style constants (required)
        from gopptx.presentation.charts import ChartType

        chart_id = slide.add_chart(
            ChartType.COLUMN,
            ["Q1", "Q2", "Q3"],
            [100, 200, 150],
            title="Sales by Quarter",
        )

        # Other chart types
        chart_id = slide.add_chart(
            ChartType.PIE,
            ["Product A", "Product B", "Product C"],
            [25.0, 35.0, 40.0],
            title="Sales Mix",
        )

        # Discover available types
        all_types = ChartType.get_all()  # Returns dict of all constants
    """

    # =========================================================================
    # Supported Chart Types
    # =========================================================================
    COLUMN = "bar"
    """Column/Vertical bar chart - displays data as vertical bars."""

    BAR = "bar"
    """Alias for COLUMN - vertical bar chart."""

    BAR_HORIZONTAL = "barHorizontal"
    """Horizontal bar chart."""

    BAR_STACKED = "barStacked"
    """Stacked bar chart."""

    BAR_STACKED_100 = "barStacked100"
    """100% stacked bar chart."""

    LINE = "line"
    """Line chart - connects data points with lines."""

    LINE_MARKERS = "lineMarkers"
    """Line chart with markers."""

    LINE_STACKED = "lineStacked"
    """Stacked line chart."""

    SCATTER = "scatter"
    """Scatter/XY chart."""

    AREA = "area"
    """Area chart."""

    AREA_STACKED = "areaStacked"
    """Stacked area chart."""

    AREA_STACKED_100 = "areaStacked100"
    """100% stacked area chart."""

    PIE = "pie"
    """Pie chart - shows data as slices of a circle."""

    DOUGHNUT = "doughnut"
    """Doughnut chart."""

    BUBBLE = "bubble"
    """Bubble chart."""

    RADAR = "radar"
    """Radar chart."""

    RADAR_FILLED = "radarFilled"
    """Filled radar chart."""

    STOCK_HLC = "stockHLC"
    """Stock high-low-close chart."""

    STOCK_OHLC = "stockOHLC"
    """Stock open-high-low-close chart."""

    COMBO = "combo"
    """Combo chart - mixes bar and line series on one category axis."""

    @staticmethod
    def get_all() -> dict[str, str]:
        """Get all named chart type constants as a dictionary.

        Returns:
            Dictionary mapping constant names to chart type values.
            Only includes currently supported chart types.

        Examples:
            all_types = ChartType.get_all()
            # Returns a dictionary of all available named chart constants.

            for name, value in all_types.items():
                print(f"ChartType.{name}")
        """
        return {
            "COLUMN": ChartType.COLUMN,
            "BAR": ChartType.BAR,
            "BAR_HORIZONTAL": ChartType.BAR_HORIZONTAL,
            "BAR_STACKED": ChartType.BAR_STACKED,
            "BAR_STACKED_100": ChartType.BAR_STACKED_100,
            "LINE": ChartType.LINE,
            "LINE_MARKERS": ChartType.LINE_MARKERS,
            "LINE_STACKED": ChartType.LINE_STACKED,
            "SCATTER": ChartType.SCATTER,
            "AREA": ChartType.AREA,
            "AREA_STACKED": ChartType.AREA_STACKED,
            "AREA_STACKED_100": ChartType.AREA_STACKED_100,
            "PIE": ChartType.PIE,
            "DOUGHNUT": ChartType.DOUGHNUT,
            "BUBBLE": ChartType.BUBBLE,
            "RADAR": ChartType.RADAR,
            "RADAR_FILLED": ChartType.RADAR_FILLED,
            "STOCK_HLC": ChartType.STOCK_HLC,
            "STOCK_OHLC": ChartType.STOCK_OHLC,
            "COMBO": ChartType.COMBO,
        }

    @staticmethod
    def validate(chart_type: str | None) -> str:
        """Validate a chart type and return its value.

        Only accepts actual ChartType constant values (e.g., "bar", "line", "pie").
        Named constant strings are NOT supported - use the ChartType constants directly
        (e.g., ChartType.COLUMN instead of "COLUMN").

        Args:
            chart_type: Chart type constant value (from ChartType enum).

        Returns:
            The same chart type string (if valid).

        Raises:
            ValueError: If chart_type is not a valid chart type value.

        Examples:
            ChartType.validate(ChartType.COLUMN)   # -> "bar"
            ChartType.validate(ChartType.SCATTER)  # -> "scatter"
            ChartType.validate("COLUMN")           # -> ValueError
            ChartType.validate("invalid")          # -> ValueError
        """
        if not chart_type:
            raise ValueError("chart_type cannot be empty")

        # Only accept actual chart type values (not named constants)
        valid_types = set(ChartType.get_all().values())
        if chart_type in valid_types:
            return chart_type

        # Not found - provide helpful error message
        valid_values = ", ".join(sorted(valid_types))
        raise ValueError(
            f"Invalid chart_type {chart_type!r}. Use ChartType constants like ChartType.COLUMN, ChartType.SCATTER, ChartType.PIE. Valid internal values: {valid_values}"
        )

    @staticmethod
    def get_by_name(name: str) -> str | None:
        """Find a chart type by its constant name.

        Args:
            name: Constant name like "COLUMN", "PIE", "LINE_MARKERS".

        Returns:
            The chart type value, or None if not found.

        Examples:
            ChartType.get_by_name("COLUMN")      # -> "bar"
            ChartType.get_by_name("PIE")          # -> "pie"
            ChartType.get_by_name("NOT_FOUND")    # -> None
        """
        return ChartType.get_all().get(name)


__all__ = ["ChartType"]
