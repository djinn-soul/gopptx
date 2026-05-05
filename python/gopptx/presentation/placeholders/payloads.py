"""Placeholder content payload builders."""

from __future__ import annotations

import math
from typing import cast

from ...utils import is_four_number_bounds

_RESERVED_CHART_OPTION_KEYS = frozenset({
    "chart_type",
    "categories",
    "values",
    "x_values",
    "y_values",
    "sizes",
    "x",
    "y",
    "w",
    "h",
})


def build_placeholder_table_payload(
    table: object,
    table_rows: object,
    table_cols: object,
) -> dict[str, object] | None:
    """Build a normalized table payload from direct dict or row/col counts."""
    if isinstance(table, dict):
        return cast("dict[str, object]", table)
    normalized_rows = _normalize_table_dimension(table_rows, "table_rows")
    normalized_cols = _normalize_table_dimension(table_cols, "table_cols")
    if normalized_rows is not None or normalized_cols is not None:
        if normalized_rows is None or normalized_cols is None:
            raise ValueError(
                "table_rows and table_cols must both be provided as positive integers"
            )
        return {
            "rows": [
                ["" for _ in range(normalized_cols)] for _ in range(normalized_rows)
            ]
        }
    return None


def apply_series_chart_values(
    chart_payload: dict[str, object],
    chart_values: list[object],
) -> None:
    """Populate chart payload values from flat or series-style chart value inputs."""
    if chart_values and isinstance(chart_values[0], dict):
        if len(chart_values) > 1:
            raise ValueError(
                "placeholder charts currently support exactly one series; got multiple series"
            )
        first = cast("dict[str, object]", chart_values[0])
        if isinstance(first.get("values"), list):
            chart_payload["values"] = _validated_finite_number_list(
                first["values"], "values"
            )
        if isinstance(first.get("x_values"), list):
            chart_payload["x_values"] = _validated_finite_number_list(
                first["x_values"], "x_values"
            )
        if isinstance(first.get("y_values"), list):
            chart_payload["y_values"] = _validated_finite_number_list(
                first["y_values"], "y_values"
            )
        if isinstance(first.get("sizes"), list):
            chart_payload["sizes"] = _validated_finite_number_list(
                first["sizes"], "sizes"
            )
        if isinstance(first.get("categories"), list):
            chart_payload["categories"] = _validated_string_list(
                cast("list[object]", first["categories"]), "categories"
            )
        title = first.get("name")
        if isinstance(title, str):
            chart_payload["title"] = title
        return
    chart_payload["values"] = _validated_finite_number_list(chart_values, "values")


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
        chart_payload["categories"] = _validated_string_list(
            cast("list[object]", chart_categories), "chart_categories"
        )
    if isinstance(chart_values, list):
        apply_series_chart_values(chart_payload, cast("list[object]", chart_values))
    if isinstance(chart_options, dict):
        options = cast("dict[str, object]", chart_options)
        chart_payload.update(_validated_chart_options(options))
    if is_four_number_bounds(bounds):
        x, y, w, h = bounds
        chart_payload["x"] = int(x)
        chart_payload["y"] = int(y)
        chart_payload["w"] = int(w)
        chart_payload["h"] = int(h)
    return chart_payload


def _validated_chart_options(raw_options: object) -> dict[str, object]:
    options = cast("dict[object, object]", raw_options)
    blocked_keys = sorted(
        str(key)
        for key in options
        if isinstance(key, str) and key in _RESERVED_CHART_OPTION_KEYS
    )
    if blocked_keys:
        raise ValueError(
            "chart_options may not include reserved keys: " + ", ".join(blocked_keys)
        )
    return cast("dict[str, object]", options)


def _normalize_table_dimension(value: object, field: str) -> int | None:
    if value is None:
        return None
    if isinstance(value, bool) or not isinstance(value, int):
        raise ValueError(f"{field} must be a positive integer, got {value!r}")
    if value <= 0:
        raise ValueError(f"{field} must be a positive integer, got {value!r}")
    return value


def _validated_finite_number_list(values: object, field: str) -> list[float]:
    if not isinstance(values, list):
        raise ValueError(f"{field} must be an array of finite numbers")
    typed_values = cast("list[object]", values)
    out: list[float] = []
    for item in typed_values:
        if isinstance(item, bool) or not isinstance(item, int | float):
            raise ValueError(f"{field} must be an array of finite numbers")
        number = float(item)
        if not math.isfinite(number):
            raise ValueError(f"{field} must be an array of finite numbers")
        out.append(number)
    return out


def _validated_string_list(values: object, field: str) -> list[str]:
    if not isinstance(values, list):
        raise ValueError(f"{field} must be an array of strings")
    typed_values = cast("list[object]", values)
    out: list[str] = []
    for item in typed_values:
        if not isinstance(item, str):
            raise ValueError(f"{field} must be an array of strings")
        out.append(item)
    return out
