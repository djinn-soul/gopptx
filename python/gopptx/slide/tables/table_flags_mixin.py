"""Table flag properties and style application mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ...presentation.tables.table_styles import TableStyle

if TYPE_CHECKING:
    from .table import Table


class TableFlagsMixin:
    """Mixin providing flag properties and style application for Table."""

    def _update_flags(self, flags: dict[str, bool]) -> None:
        table = cast("Table", self)
        table.prs.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {
                "slide_index": table.slide_index,
                "shape_id": table.shape_id,
                "flags": flags,
            },
        )
        if not getattr(table.prs, "_batch_active", False):
            table.invalidate_cache()

    @property
    def first_row(self) -> bool:
        """Compatibility alias for ``header_row_enabled``."""
        return cast("Table", self).header_row_enabled

    @first_row.setter
    def first_row(self, value: bool) -> None:
        cast("Table", self).header_row_enabled = value

    @property
    def horz_banding(self) -> bool:
        """Compatibility alias for ``banded_rows_enabled``."""
        return cast("Table", self).banded_rows_enabled

    @horz_banding.setter
    def horz_banding(self, value: bool) -> None:
        cast("Table", self).banded_rows_enabled = value

    @property
    def first_col(self) -> bool:
        """Whether first-column emphasis is enabled."""
        table = cast("Table", self)
        ensure_cache = getattr(table, "_ensure_cache", None)
        cache = getattr(table, "_cache", None)
        if callable(ensure_cache):
            ensure_cache()
        cache = getattr(table, "_cache", cache)
        cache_dict: dict[str, object] = (
            cast("dict[str, object]", cache) if isinstance(cache, dict) else {}
        )
        return cache_dict.get("first_col", False) is True

    @first_col.setter
    def first_col(self, value: bool) -> None:
        self._update_flags({"first_col": value})

    @property
    def last_col(self) -> bool:
        """Whether last-column emphasis is enabled."""
        table = cast("Table", self)
        ensure_cache = getattr(table, "_ensure_cache", None)
        cache = getattr(table, "_cache", None)
        if callable(ensure_cache):
            ensure_cache()
        cache = getattr(table, "_cache", cache)
        cache_dict: dict[str, object] = (
            cast("dict[str, object]", cache) if isinstance(cache, dict) else {}
        )
        return cache_dict.get("last_col", False) is True

    @last_col.setter
    def last_col(self, value: bool) -> None:
        self._update_flags({"last_col": value})

    @property
    def last_row(self) -> bool:
        """Whether last-row emphasis is enabled."""
        table = cast("Table", self)
        ensure_cache = getattr(table, "_ensure_cache", None)
        cache = getattr(table, "_cache", None)
        if callable(ensure_cache):
            ensure_cache()
        cache = getattr(table, "_cache", cache)
        cache_dict: dict[str, object] = (
            cast("dict[str, object]", cache) if isinstance(cache, dict) else {}
        )
        return cache_dict.get("last_row", False) is True

    @last_row.setter
    def last_row(self, value: bool) -> None:
        self._update_flags({"last_row": value})

    @property
    def vert_banding(self) -> bool:
        """Whether alternating column banding is enabled."""
        table = cast("Table", self)
        ensure_cache = getattr(table, "_ensure_cache", None)
        cache = getattr(table, "_cache", None)
        if callable(ensure_cache):
            ensure_cache()
        cache = getattr(table, "_cache", cache)
        cache_dict: dict[str, object] = (
            cast("dict[str, object]", cache) if isinstance(cache, dict) else {}
        )
        return cache_dict.get("band_col", False) is True

    @vert_banding.setter
    def vert_banding(self, value: bool) -> None:
        self._update_flags({"band_col": value})

    def apply_style(self, style: str | int) -> None:
        """Apply a table style by name or GUID."""
        style_guid = style
        if isinstance(style, str) and not style.startswith("{"):
            styles = TableStyle.get_all()
            if style not in styles:
                available = ", ".join(sorted(styles.keys()))
                raise ValueError(
                    f"Unknown style name '{style}'. Available: {available}"
                )
            style_guid = styles[style]

        table = cast("Table", self)
        table.prs.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": table.slide_index,
                "shape_id": table.shape_id,
                "style_guid": style_guid,
            },
        )
        if not getattr(table.prs, "_batch_active", False):
            table.invalidate_cache()
