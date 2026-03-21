"""Convenience builders for creating tables with common patterns."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from ...presentation.presentation import Presentation


class PresentationTableBuilders:
    """Mixin providing convenience table builder methods."""

    if TYPE_CHECKING:

        def add_table(
            self,
            slide: int | None = None,
            rows: int | None = None,
            cols: int | None = None,
            bounds: tuple[int, int, int, int] | None = None,
            **kwargs: object,
        ) -> int:
            """Add a table (signature for type hints)."""
            ...

    def add_table_from_rows(
        self: Presentation,
        slide: int,
        rows: list[list[str]],
        bounds: tuple[int, int, int, int],
        *,
        first_row: bool = True,
        band_row: bool = True,
        **kwargs: object,
    ) -> int:
        """Create a table from a list of row data.

        This is the recommended way to add tables for documentation and examples.

        Args:
            slide: Zero-based slide index.
            rows: List of rows, each row is a list of strings.
            bounds: (x, y, cx, cy) table position and size in EMU.
            first_row: Treat first row as header (default True).
            band_row: Enable row alternating colors (default True).
            column_widths: Optional EMU widths per column.
            **kwargs: Additional flags (first_col, last_col, band_col, etc).

        Returns:
            Shape ID of the created table.

        Example:
            table_id = prs.add_table_from_rows(
                slide=0,
                rows=[
                    ["Item", "Qty", "Status"],
                    ["Widgets", "50", "Ready"],
                    ["Gadgets", "30", "Queued"],
                ],
                bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
                first_row=True,
                band_row=True,
            )
        """
        row_count = len(rows)
        col_count = max(len(row) for row in rows) if rows else 0

        return self.add_table(
            slide=slide,
            rows=row_count,
            cols=col_count,
            bounds=bounds,
            data=rows,
            first_row=first_row,
            band_row=band_row,
            **kwargs,
        )

    def add_table_from_dicts(
        self: Presentation,
        slide: int,
        rows: list[dict[str, str]],
        bounds: tuple[int, int, int, int] | None = None,
        *,
        column_names: list[str] | None = None,
        **kwargs: object,
    ) -> int:
        """Create a table from a list of dictionaries.

        Each dict represents one data row. Column names determine the order.

        Args:
            slide: Zero-based slide index.
            rows: List of dicts with consistent keys.
            column_names: Ordered list of keys to use as columns and header.
                         If None, uses sorted keys from first row.
            bounds: (x, y, cx, cy) table bounds. Required if add_table doesn't have defaults.
            first_row: Treat first row as header (default True).
            band_row: Enable row alternating colors (default True).
            **kwargs: Additional flags.

        Returns:
            Shape ID of created table.

        Example:
            table_id = prs.add_table_from_dicts(
                slide=0,
                rows=[
                    {"item": "Widgets", "qty": "50", "status": "Ready"},
                    {"item": "Gadgets", "qty": "30", "status": "Queued"},
                ],
                column_names=["item", "qty", "status"],
                bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
            )
        """
        if not rows:
            raise ValueError("rows list cannot be empty")

        # Determine column order
        if column_names is None:
            column_names = sorted(rows[0].keys())

        # Build 2D array with header row
        data: list[list[str]] = [column_names]  # Header
        data.extend([row_dict.get(col, "") for col in column_names] for row_dict in rows)

        return self.add_table_from_rows(
            slide=slide,
            rows=data,
            bounds=bounds,
            **kwargs,
        )


__all__ = ["PresentationTableBuilders"]
