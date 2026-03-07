"""Chart data builder classes for python-pptx-style ergonomics."""
# ruff: noqa: D102

from __future__ import annotations

from dataclasses import dataclass, field
from typing import cast


@dataclass(slots=True)
class CategorySeries:
    """One category chart series."""

    name: str
    values: list[float]


@dataclass(slots=True)
class CategoryChartData:
    """Builder for category-based chart data."""

    categories: list[str] = field(default_factory=list)
    series: list[CategorySeries] = field(default_factory=list)

    def add_category(self, value: str) -> None:
        self.categories.append(value)

    def add_series(self, name: str, values: list[float]) -> None:
        self.series.append(CategorySeries(name=name, values=[float(v) for v in values]))

    def to_add_chart_args(
        self,
    ) -> tuple[list[str], list[float] | list[dict[str, object]]]:
        if not self.series:
            return self.categories, []
        if len(self.series) == 1:
            return self.categories, list(self.series[0].values)
        payload: list[dict[str, object]] = [
            {"name": item.name, "categories": self.categories, "values": item.values}
            for item in self.series
        ]
        return self.categories, payload

    def to_update_payload(self) -> dict[str, object]:
        return {
            "categories": self.categories,
            "series": [
                {
                    "name": item.name,
                    "categories": self.categories,
                    "values": item.values,
                }
                for item in self.series
            ],
        }


@dataclass(slots=True)
class XySeries:
    """One XY chart series."""

    name: str
    x_values: list[float]
    y_values: list[float]
    sizes: list[float] | None = None


@dataclass(slots=True)
class XyChartData:
    """Builder for XY/scatter/bubble chart data."""

    series: list[XySeries] = field(default_factory=list)

    def add_series(
        self,
        name: str,
        x_values: list[float],
        y_values: list[float],
        sizes: list[float] | None = None,
    ) -> None:
        self.series.append(
            XySeries(
                name=name,
                x_values=[float(v) for v in x_values],
                y_values=[float(v) for v in y_values],
                sizes=[float(v) for v in sizes] if sizes is not None else None,
            )
        )

    def to_update_payload(self) -> dict[str, object]:
        out_series: list[dict[str, object]] = []
        for item in self.series:
            entry: dict[str, object] = {
                "name": item.name,
                "x_values": item.x_values,
                "y_values": item.y_values,
            }
            if item.sizes is not None:
                entry["sizes"] = item.sizes
            out_series.append(entry)
        return {"series": out_series}

    def to_add_chart_args(
        self,
    ) -> tuple[list[str], list[float] | list[dict[str, object]]]:
        payload = self.to_update_payload()
        # add_chart legacy path expects category/value signature; pass series payload.
        series = cast("list[dict[str, object]]", payload.get("series", []))
        return [], series
