"""Data-table proxy for the grid drawn under a chart's plot area."""

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ...schemas import ChartFormatUpdate, ChartState

_FONT_SIZE_MIN_PT = 1
_FONT_SIZE_MAX_PT = 400


class _DataTableChartProtocol(Protocol):
    """Minimal chart surface used by data-table helpers."""

    def snapshot(self) -> ChartState:
        """Return the current chart state."""
        ...

    def apply_format(self, fmt: ChartFormatUpdate) -> None:
        """Apply a chart format update."""
        ...


class ChartDataTable:
    """Live proxy for the ``c:dTable`` under a chart's plot area."""

    def __init__(self, chart: _DataTableChartProtocol) -> None:
        """Initialize a data-table proxy bound to a chart."""
        super().__init__()
        self._chart = chart

    def _payload(self) -> dict[str, object]:
        # The snapshot is a bridge payload, so the declared type is a promise
        # the runtime value has to be checked against.
        raw: object = self._chart.snapshot().get("data_table")
        if not isinstance(raw, dict):
            return {}
        return dict(cast("dict[str, object]", raw))

    def _flag(self, key: str) -> bool | None:
        value = self._payload().get(key)
        return value if isinstance(value, bool) else None

    def _apply(self, field: str, value: object) -> None:
        if not self.visible:
            raise ValueError("the chart has no data table; set visible first")
        self._chart.apply_format(
            cast("ChartFormatUpdate", {"data_table": {"show": True, field: value}})
        )

    @property
    def visible(self) -> bool:
        """Return whether the chart draws a data table."""
        return bool(self._payload())

    @visible.setter
    def visible(self, value: bool) -> None:
        self._chart.apply_format(
            cast("ChartFormatUpdate", {"data_table": {"show": bool(value)}})
        )

    @property
    def show_horizontal_border(self) -> bool | None:
        """Return whether horizontal cell borders are drawn."""
        return self._flag("show_horizontal_border")

    @show_horizontal_border.setter
    def show_horizontal_border(self, value: bool) -> None:
        self._apply("show_horizontal_border", bool(value))

    @property
    def show_vertical_border(self) -> bool | None:
        """Return whether vertical cell borders are drawn."""
        return self._flag("show_vertical_border")

    @show_vertical_border.setter
    def show_vertical_border(self, value: bool) -> None:
        self._apply("show_vertical_border", bool(value))

    @property
    def show_outline(self) -> bool | None:
        """Return whether the table outline is drawn."""
        return self._flag("show_outline")

    @show_outline.setter
    def show_outline(self, value: bool) -> None:
        self._apply("show_outline", bool(value))

    @property
    def show_legend_keys(self) -> bool | None:
        """Return whether legend keys are shown beside the series names."""
        return self._flag("show_keys")

    @show_legend_keys.setter
    def show_legend_keys(self, value: bool) -> None:
        self._apply("show_keys", bool(value))

    @property
    def font_size_pt(self) -> int | None:
        """Return the data-table font size in points."""
        value = self._payload().get("font_size_pt")
        # bool is a subclass of int, so it would otherwise become 1/0.
        if isinstance(value, bool) or not isinstance(value, int):
            return None
        return value

    @font_size_pt.setter
    def font_size_pt(self, value: int) -> None:
        self._apply("font_size_pt", _validated_font_size(value))

    def show(
        self,
        *,
        horizontal_border: bool | None = None,
        vertical_border: bool | None = None,
        outline: bool | None = None,
        legend_keys: bool | None = None,
        font_size_pt: int | None = None,
    ) -> None:
        """Turn the data table on, optionally setting its options at once.

        Raises:
            ValueError: If ``font_size_pt`` is outside the 1-400 point range.
        """
        spec: dict[str, object] = {"show": True}
        for key, value in (
            ("show_horizontal_border", horizontal_border),
            ("show_vertical_border", vertical_border),
            ("show_outline", outline),
            ("show_keys", legend_keys),
        ):
            if value is not None:
                spec[key] = bool(value)
        if font_size_pt is not None:
            spec["font_size_pt"] = _validated_font_size(font_size_pt)
        self._chart.apply_format(cast("ChartFormatUpdate", {"data_table": spec}))

    def hide(self) -> None:
        """Remove the data table from the chart."""
        self.visible = False


def _validated_font_size(value: int) -> int:
    """Return the font size in points, rejecting anything out of range.

    Raises:
        ValueError: If the size is outside the 1-400 point range.
    """
    if not _FONT_SIZE_MIN_PT <= value <= _FONT_SIZE_MAX_PT:
        message = (
            f"font_size_pt must be between {_FONT_SIZE_MIN_PT} and {_FONT_SIZE_MAX_PT}"
        )
        raise ValueError(message)
    return int(value)
