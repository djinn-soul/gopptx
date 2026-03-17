"""Placeholder content payload builders."""

from __future__ import annotations

from typing import cast

from ...utils import is_four_number_bounds


def build_placeholder_table_payload(
    table: object,
    table_rows: object,
    table_cols: object,
) -> dict[str, object] | None:
    """Build a normalized table payload from direct dict or row/col counts."""
    if isinstance(table, dict):
        return cast("dict[str, object]", table)
    if isinstance(table_rows, int) and isinstance(table_cols, int):
        return {"rows": [["" for _ in range(table_cols)] for _ in range(table_rows)]}
    return None


def apply_series_chart_values(
    chart_payload: dict[str, object],
    chart_values: list[object],
) -> None:
    """Populate chart payload values from flat or series-style chart value inputs."""
    if chart_values and isinstance(chart_values[0], dict):
        first = cast("dict[str, object]", chart_values[0])
        if isinstance(first.get("values"), list):
            chart_payload["values"] = cast("list[float]", first["values"])
        if isinstance(first.get("x_values"), list):
            chart_payload["x_values"] = cast("list[float]", first["x_values"])
        if isinstance(first.get("y_values"), list):
            chart_payload["y_values"] = cast("list[float]", first["y_values"])
        if isinstance(first.get("sizes"), list):
            chart_payload["sizes"] = cast("list[float]", first["sizes"])
        if isinstance(first.get("categories"), list):
            chart_payload["categories"] = cast("list[str]", first["categories"])
        title = first.get("name")
        if isinstance(title, str):
            chart_payload["title"] = title
        return
    chart_payload["values"] = cast("list[float]", chart_values)


def build_placeholder_chart_payload(
    kwargs: dict[str, object],
    bounds: object,
) -> dict[str, object] | None:
    """Build a normalized chart payload from placeholder insertion kwargs."""
    chart = kwargs.get("chart")
    chart_type = kwargs.get("chart_type")
    chart_categories = kwargs.get("chart_categories")
    chart_values = kwargs.get("chart_values")
    chart_options = kwargs.get("chart_options")
    if isinstance(chart, dict):
        return cast("dict[str, object]", chart)
    if not isinstance(chart_type, str):
        return None

    chart_payload: dict[str, object] = {"chart_type": chart_type}
    if isinstance(chart_categories, list):
        chart_payload["categories"] = chart_categories
    if isinstance(chart_values, list):
        apply_series_chart_values(chart_payload, cast("list[object]", chart_values))
    if isinstance(chart_options, dict):
        chart_payload.update(cast("dict[str, object]", chart_options))
    if is_four_number_bounds(bounds):
        x, y, w, h = bounds
        chart_payload["x"] = int(x)
        chart_payload["y"] = int(y)
        chart_payload["w"] = int(w)
        chart_payload["h"] = int(h)
    return chart_payload
