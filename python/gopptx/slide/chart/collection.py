"""Chart collection facade."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .model import Chart

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ._protocols import ChartSlideProto


class ChartCollection:
    """Collection facade for slide chart proxies."""

    def __init__(self, slide: ChartSlideProto) -> None:
        """Bind a chart collection to a slide proxy."""
        super().__init__()
        self._slide = slide

    def _items(self) -> list[Chart]:
        """Build chart proxies from live slide chart references."""
        refs = self._slide.list_charts()
        out: list[Chart] = []
        for item in refs:
            ref = item
            index = int(ref.get("Index", ref.get("index", 0)))
            rel_id = str(ref.get("RelID", ref.get("rel_id", "")))
            chart_part = str(ref.get("ChartPart", ref.get("chart_part", "")))
            out.append(Chart(self._slide, index, rel_id, chart_part))
        out.sort(key=lambda chart: chart.index)
        return out

    def __len__(self) -> int:
        """Return number of charts on this slide."""
        return len(self._items())

    def __getitem__(self, index: int) -> Chart:
        """Return one chart by index."""
        items = self._items()
        if index < 0:
            index += len(items)
        if index < 0 or index >= len(items):
            raise IndexError("chart index out of range")
        return items[index]

    def __iter__(self) -> Iterator[Chart]:
        """Iterate charts in index order."""
        return iter(self._items())
