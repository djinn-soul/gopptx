"""Presentation table mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ...api_errors import GopptxError
from ...utils import normalize_table_index
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import TableCellInfo, TableInfo


class PresentationTableMixin(PresentationMixinBase):
    """Mixin providing table creation and manipulation methods."""

    def add_table(  # noqa: C901, PLR0913
        self,
        slide: int | None = None,
        slide_index: int | None = None,
        rows: int | None = None,
        cols: int | None = None,
        bounds: tuple[int, int, int, int] | None = None,
        x: int | None = None,
        y: int | None = None,
        cx: int | None = None,
        cy: int | None = None,
        data: list[list[str]] | None = None,
        first_row: bool = False,
        first_col: bool = False,
        last_row: bool = False,
        last_col: bool = False,
        band_row: bool = False,
        band_col: bool = False,
        column_widths: list[int] | None = None,
    ) -> int:
        """Add a table shape to a slide and return its shape ID.

        Args:
            slide: Zero-based slide index (preferred over slide_index).
            slide_index: Zero-based slide index (deprecated, use slide).
            rows: Number of rows.
            cols: Number of columns.
            bounds: Tuple of (x, y, cx, cy) in EMU.
            x: X position in EMU (used if bounds not provided).
            y: Y position in EMU (used if bounds not provided).
            cx: Width in EMU (used if bounds not provided).
            cy: Height in EMU (used if bounds not provided).
            data: Optional 2D array of cell text [[row0...], [row1...], ...].
            first_row: Enable first-row header formatting.
            first_col: Enable first-column emphasis.
            last_row: Enable last-row emphasis.
            last_col: Enable last-column emphasis.
            band_row: Enable alternating row colors.
            band_col: Enable alternating column colors.
            column_widths: Optional list of column widths in EMU.

        Returns:
            Shape ID of the created table.

        Examples:
            # New style - recommended with data and named params
            table_id = prs.add_table(
                slide=0,
                rows=3, cols=2,
                bounds=(Inches(1), Inches(2), Inches(8), Inches(4)),
                data=[["Item", "Qty"], ["Widgets", "50"], ["Gadgets", "30"]],
                first_row=True,
                band_row=True,
            )

            # Old style - still works (backward compat)
            table_id = prs.add_table(0, 3, 2, (x, y, cx, cy))
        """
        # Detect old positional calling pattern: add_table(slide_idx, rows, cols, bounds)
        # In that case: slide=slide_idx, slide_index=rows, rows=cols, cols=bounds
        if (
            slide is not None
            and isinstance(slide, int)
            and slide_index is not None
            and isinstance(slide_index, int)
            and rows is not None
            and isinstance(rows, int)
            and cols is not None
            and isinstance(cols, tuple)
            and len(cast("tuple", cols)) == 4
        ):
            # Old positional pattern detected: remap parameters
            slide_idx = slide
            rows_count = slide_index
            cols_count = rows
            bounds = cast("tuple", cols)
            rows = rows_count
            cols = cols_count
            # Validate
            if rows is None or cols is None:
                raise ValueError("'rows' and 'cols' parameters are required")
        else:
            # New named or mixed pattern
            slide_idx = slide if slide is not None else slide_index
            if slide_idx is None:
                raise ValueError("Either 'slide' or 'slide_index' must be provided")

            # Validate rows/cols
            if rows is None:
                raise ValueError("'rows' parameter is required")
            if cols is None:
                raise ValueError("'cols' parameter is required")

        # Parse bounds
        if bounds is not None:
            x, y, cx, cy = bounds
        elif x is None or y is None or cx is None or cy is None:
            raise ValueError("Either 'bounds' tuple or (x, y, cx, cy) must be provided")

        # 1. Create empty table
        result = self.execute(
            ops.OP_ADD_TABLE,
            {
                "slide_index": slide_idx,
                "rows": rows,
                "cols": cols,
                "x": x,
                "y": y,
                "cx": cx,
                "cy": cy,
            },
        )
        shape_id = int(cast("int", result.get("shape_id", 0)))

        # 2. Populate cells if data provided
        if data is not None:
            for row_idx, row in enumerate(data):
                for col_idx, text in enumerate(row):
                    self.set_table_cell_text(
                        slide_idx, shape_id, row_idx, col_idx, text
                    )

        # 3. Set column widths if provided
        if column_widths is not None:
            for col_idx, width in enumerate(column_widths):
                self.set_table_column_width(slide_idx, shape_id, col_idx, width)

        # 4. Set flags if any are true
        flags: dict[str, bool] = {
            "first_row": first_row,
            "first_col": first_col,
            "last_row": last_row,
            "last_col": last_col,
            "band_row": band_row,
            "band_col": band_col,
        }
        if any(flags.values()):
            self.set_table_flags(slide_idx, shape_id, flags)

        return shape_id

    def get_table(self, slide_index: int, shape_id: int) -> TableInfo:
        """Return serialized table information for a table shape."""
        result = self.execute(
            ops.OP_GET_TABLE,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("TableInfo", cast("dict[str, object]", result.get("table", {})))

    def set_table_style(self, slide_index: int, shape_id: int, style: str) -> None:
        """Apply a table style by name or GUID.

        Args:
            style: Style name (e.g., "MEDIUM_STYLE_2") or GUID string.
                   Use TableStyle constants for easy style selection.

        Example:
            from gopptx.presentation.tables import TableStyle

            prs.set_table_style(0, table_id, TableStyle.MEDIUM_STYLE_2)
            prs.set_table_style(0, table_id, "MEDIUM_STYLE_2")
        """
        from .table_styles import TableStyle

        # Handle style name lookup
        style_guid = style
        if isinstance(style, str) and not style.startswith("{"):
            # It's a style name, look it up
            styles = TableStyle.get_all()
            if style not in styles:
                available = ", ".join(sorted(styles.keys()))
                raise ValueError(
                    f"Unknown style name '{style}'. Available: {available}"
                )
            style_guid = styles[style]

        self.execute(
            ops.OP_SET_TABLE_STYLE,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "style_guid": style_guid,
            },
        )

    def define_table_style(self, name: str, style_id: str | None = None) -> str:
        """Define a custom table style and return its resolved style ID."""
        payload: dict[str, object] = {"name": name}
        if style_id is not None:
            payload["style_id"] = style_id
        result = self.execute(ops.OP_DEFINE_TABLE_STYLE, payload)
        return str(result.get("style_id", ""))

    def list_table_styles(self) -> list[dict[str, str]]:
        """List available table styles visible to the presentation.

        Returns:
            List of style dicts with 'name' and 'guid' keys.

        Example:
            styles = prs.list_table_styles()
            for style in styles:
                print(f"{style['name']}: {style['guid']}")
        """
        result = self.execute(ops.OP_LIST_TABLE_STYLES, {})
        return cast("list[dict[str, str]]", result.get("styles", []))

    def get_table_style_by_name(self, name: str) -> str | None:
        """Find a presentation table style GUID by name.

        Args:
            name: Style name to search for.

        Returns:
            Style GUID if found, None otherwise.

        Example:
            guid = prs.get_table_style_by_name("Medium Style 2 - Accent 1")
            if guid:
                prs.set_table_style(0, table_id, guid)
        """
        styles = self.list_table_styles()
        for style in styles:
            if style.get("name", "").lower() == name.lower():
                return style.get("guid")
        return None

    def get_all_table_style_names(self) -> list[str]:
        """Get all available table style names in the presentation.

        Returns:
            List of style names.

        Example:
            names = prs.get_all_table_style_names()
            print("Available styles:", names)
        """
        styles = self.list_table_styles()
        return [style.get("name", "") for style in styles if "name" in style]

    def set_table_flags(
        self,
        slide_index: int,
        shape_id: int,
        flags: dict[str, bool],
    ) -> None:
        """Set table display flags such as header-row or banded options."""
        self.execute(
            ops.OP_UPDATE_TABLE_FLAGS,
            {"slide_index": slide_index, "shape_id": shape_id, "flags": flags},
        )

    def set_table_cell_text(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        col: int,
        text: str,
    ) -> None:
        """Update the text value for one table cell."""
        self.execute(
            ops.OP_UPDATE_TABLE_CELL,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row": row,
                "col": col,
                "updates": {"text": text},
            },
        )

    def get_table_cell(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        col: int,
    ) -> TableCellInfo:
        """Return one table cell payload by zero-based row and column."""
        table = self.get_table(slide_index, shape_id)
        cells = cast("list[dict[str, object]]", table.get("cells", []))
        cell_map: dict[tuple[int, int], dict[str, object]] = {}
        for cell in cells:
            try:
                row_idx = normalize_table_index(cell["row"])
                col_idx = normalize_table_index(cell["col"])
            except (KeyError, ValueError):
                continue
            cell_map[row_idx, col_idx] = cell
        found = cell_map.get((row, col))
        if found is not None:
            return cast("TableCellInfo", found)
        raise GopptxError(f"table cell [{row},{col}] not found", code="OP_FAILED")

    def merge_table_cells(
        self,
        slide_index: int,
        shape_id: int,
        cell_range: tuple[int, int, int, int],
    ) -> None:
        """Merge a rectangular range of table cells."""
        row1, col1, row2, col2 = cell_range
        self.execute(
            ops.OP_MERGE_TABLE_CELLS,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row1": row1,
                "col1": col1,
                "row2": row2,
                "col2": col2,
            },
        )

    def split_table_cell(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        col: int,
    ) -> None:
        """Split a merged table cell back into its original cells."""
        self.execute(
            ops.OP_SPLIT_TABLE_CELL,
            {"slide_index": slide_index, "shape_id": shape_id, "row": row, "col": col},
        )

    def set_table_row_height(
        self,
        slide_index: int,
        shape_id: int,
        row: int,
        height: int,
    ) -> None:
        """Set the height of a specific table row."""
        self.execute(
            ops.OP_SET_TABLE_ROW_HEIGHT,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "row": row,
                "height": height,
            },
        )

    def set_table_column_width(
        self,
        slide_index: int,
        shape_id: int,
        col: int,
        width: int,
    ) -> None:
        """Set the width of a specific table column."""
        self.execute(
            ops.OP_SET_TABLE_COLUMN_WIDTH,
            {
                "slide_index": slide_index,
                "shape_id": shape_id,
                "col": col,
                "width": width,
            },
        )
