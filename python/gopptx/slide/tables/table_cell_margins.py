"""Margin properties shared by live table-cell proxies."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ._protocols import TableWriteProto


class _CellMarginHost(Protocol):
    _table: TableWriteProto
    row: int
    col: int

    def _margin(self, name: str, default: int) -> int: ...

    def _set_margin(self, name: str, value: int) -> None: ...


class CellMarginMixin:
    """Expose OOXML table-cell margins in EMU."""

    def _margin(self: _CellMarginHost, name: str, default: int) -> int:
        info = self._table.get_cell_info(self.row, self.col)
        return cast("int", info.get(name, default))

    def _set_margin(self: _CellMarginHost, name: str, value: int) -> None:
        self._table.update_cell(self.row, self.col, {name: int(value)})

    @property
    def margin_left(self: _CellMarginHost) -> int:
        """Return left cell margin in EMU."""
        return self._margin("margin_left", 91440)

    @margin_left.setter
    def margin_left(self: _CellMarginHost, value: int) -> None:
        self._set_margin("margin_left", value)

    @property
    def margin_right(self: _CellMarginHost) -> int:
        """Return right cell margin in EMU."""
        return self._margin("margin_right", 91440)

    @margin_right.setter
    def margin_right(self: _CellMarginHost, value: int) -> None:
        self._set_margin("margin_right", value)

    @property
    def margin_top(self: _CellMarginHost) -> int:
        """Return top cell margin in EMU."""
        return self._margin("margin_top", 45720)

    @margin_top.setter
    def margin_top(self: _CellMarginHost, value: int) -> None:
        self._set_margin("margin_top", value)

    @property
    def margin_bottom(self: _CellMarginHost) -> int:
        """Return bottom cell margin in EMU."""
        return self._margin("margin_bottom", 45720)

    @margin_bottom.setter
    def margin_bottom(self: _CellMarginHost, value: int) -> None:
        self._set_margin("margin_bottom", value)
