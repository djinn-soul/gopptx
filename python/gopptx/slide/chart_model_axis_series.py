"""Axis and series helpers for live chart proxies."""
# ruff: noqa: D101,D102,D105,D107,SLF001,TC003
# pyright: reportMissingSuperCall=false, reportPrivateUsage=false

from __future__ import annotations

from collections.abc import Iterator
from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    from ..schemas import ChartAxisState, ChartFormatUpdate
    from .chart_model import Chart


class ChartAxis:
    def __init__(self, chart: Chart, *, axis_name: str) -> None:
        self._chart = chart
        self._axis_name = axis_name

    def _payload(self) -> ChartAxisState:
        snapshot = self._chart._snapshot()
        key = "category_axis" if self._axis_name == "category" else "value_axis"
        raw = snapshot.get(key, {})
        return cast("ChartAxisState", raw if isinstance(raw, dict) else {})

    @property
    def present(self) -> bool:
        return bool(self._payload().get("present", False))

    @property
    def tick_label_position(self) -> str | None:
        value = self._payload().get("tick_label_pos")
        return str(value) if isinstance(value, str) else None

    @tick_label_position.setter
    def tick_label_position(self, value: str) -> None:
        key = (
            "category_axis_tick_label_pos"
            if self._axis_name == "category"
            else "value_axis_tick_label_pos"
        )
        self._chart._apply_format(cast("ChartFormatUpdate", {key: value}))

    @property
    def has_major_gridlines(self) -> bool:
        return bool(self._payload().get("has_major_gridline", False))


class ChartSeries:
    def __init__(self, payload: dict[str, object]) -> None:
        self._payload = payload

    @property
    def name(self) -> str | None:
        value = self._payload.get("name")
        return str(value) if isinstance(value, str) else None

    @property
    def values(self) -> list[float]:
        raw = self._payload.get("values")
        if not isinstance(raw, list):
            return []
        values_raw = cast("list[object]", raw)
        return [float(item) for item in values_raw if isinstance(item, int | float)]


class ChartSeriesCollection:
    def __init__(self, payload: list[dict[str, object]]) -> None:
        self._payload = payload

    def __len__(self) -> int:
        return len(self._payload)

    def __getitem__(self, index: int) -> ChartSeries:
        if index < 0:
            index += len(self._payload)
        if index < 0 or index >= len(self._payload):
            raise IndexError("series index out of range")
        return ChartSeries(self._payload[index])

    def __iter__(self) -> Iterator[ChartSeries]:
        for item in self._payload:
            yield ChartSeries(item)
