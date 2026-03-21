"""Table flag properties and style application mixin."""

from __future__ import annotations

from ... import ops


class _TableFlagsMixin:
    """Mixin providing flag properties and style application for Table."""

    def _update_flags(self, flags: dict[str, bool]) -> None:
        self.prs.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {
                "slide_index": self.slide_index,
                "shape_id": self.shape_id,
                "flags": flags,
            },
        )
        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()

    @property
    def first_row(self) -> bool:
        """Compatibility alias for ``header_row_enabled``."""
        return self.header_row_enabled

    @first_row.setter
    def first_row(self, value: bool) -> None:
        self.header_row_enabled = value

    @property
    def horz_banding(self) -> bool:
        """Compatibility alias for ``banded_rows_enabled``."""
        return self.banded_rows_enabled

    @horz_banding.setter
    def horz_banding(self, value: bool) -> None:
        self.banded_rows_enabled = value

    @property
    def first_col(self) -> bool:
        """Whether first-column emphasis is enabled."""
        self._ensure_cache()
        return self._cache.get("first_col", False) is True if self._cache else False

    @first_col.setter
    def first_col(self, value: bool) -> None:
        self._update_flags({"first_col": value})

    @property
    def last_col(self) -> bool:
        """Whether last-column emphasis is enabled."""
        self._ensure_cache()
        return self._cache.get("last_col", False) is True if self._cache else False

    @last_col.setter
    def last_col(self, value: bool) -> None:
        self._update_flags({"last_col": value})

    @property
    def last_row(self) -> bool:
        """Whether last-row emphasis is enabled."""
        self._ensure_cache()
        return self._cache.get("last_row", False) is True if self._cache else False

    @last_row.setter
    def last_row(self, value: bool) -> None:
        self._update_flags({"last_row": value})

    @property
    def vert_banding(self) -> bool:
        """Whether alternating column banding is enabled."""
        self._ensure_cache()
        return self._cache.get("band_col", False) is True if self._cache else False

    @vert_banding.setter
    def vert_banding(self, value: bool) -> None:
        self._update_flags({"band_col": value})

    def apply_style(self, style: str | int) -> None:
        """Apply a table style by name or GUID."""
        from ...presentation.tables.table_styles import TableStyle

        style_guid = style
        if isinstance(style, str) and not style.startswith("{"):
            styles = TableStyle.get_all()
            if style not in styles:
                available = ", ".join(sorted(styles.keys()))
                raise ValueError(
                    f"Unknown style name '{style}'. Available: {available}"
                )
            style_guid = styles[style]

        self.prs.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": self.slide_index,
                "shape_id": self.shape_id,
                "style_guid": style_guid,
            },
        )
        if not getattr(self.prs, "_batch_active", False):
            self.invalidate_cache()
