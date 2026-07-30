"""Gridline styling for a chart axis (issue #984)."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .line_format import build_line_format

if TYPE_CHECKING:
    from ...schemas import ChartAxisState, ChartLineFormatSpec


class ChartAxisGridlineMixin:
    """Adds colour, width and dash controls for an axis' gridlines.

    The visibility flags alone leave PowerPoint drawing its default grey line;
    the style lives in the gridline's own ``c:spPr``.
    """

    def _payload(self) -> ChartAxisState:
        """Return the current axis state payload from the concrete proxy."""
        raise NotImplementedError

    def _apply_axis_format(self, name: str, value: object) -> None:
        """Apply one axis-scoped formatting key, prefixed by the axis name."""
        raise NotImplementedError

    @property
    def major_gridlines_format(self) -> ChartLineFormatSpec | None:
        """Return the explicit style of the major gridlines, if any."""
        return self._gridline_format("major_gridline_format")

    @property
    def minor_gridlines_format(self) -> ChartLineFormatSpec | None:
        """Return the explicit style of the minor gridlines, if any."""
        return self._gridline_format("minor_gridline_format")

    def _gridline_format(self, key: str) -> ChartLineFormatSpec | None:
        raw: object = self._payload().get(key)
        if not isinstance(raw, dict):
            return None
        return cast("ChartLineFormatSpec", raw)

    def format_major_gridlines(
        self,
        *,
        color: str | None = None,
        width_emu: int | None = None,
        dash: str | None = None,
        none: bool | None = None,
    ) -> None:
        """Style the major gridlines, drawing them if the axis has none."""
        self._format_gridlines(
            "major", color=color, width_emu=width_emu, dash=dash, none=none
        )

    def format_minor_gridlines(
        self,
        *,
        color: str | None = None,
        width_emu: int | None = None,
        dash: str | None = None,
        none: bool | None = None,
    ) -> None:
        """Style the minor gridlines, drawing them if the axis has none."""
        self._format_gridlines(
            "minor", color=color, width_emu=width_emu, dash=dash, none=none
        )

    def _format_gridlines(
        self,
        which: str,
        *,
        color: str | None,
        width_emu: int | None,
        dash: str | None,
        none: bool | None,
    ) -> None:
        spec = build_line_format(
            color=color,
            width_emu=width_emu,
            dash=dash,
            none=none,
            name=f"{which} gridlines",
        )
        self._apply_axis_format(f"{which}_gridline_format", spec)
