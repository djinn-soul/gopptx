"""Tests for table row/column dimension accessors."""
# ruff: noqa: D103,PLR2004

from __future__ import annotations

from gopptx import Presentation


def test_table_row_and_column_dimension_accessors() -> None:
    with Presentation.new("Table Dimensions") as prs:
        slide = prs.add_slide("Table")
        shape_id = slide.add_table(2, 2, bounds=(1000, 1000, 5000, 2000))
        table = slide.table(shape_id)

        table.rows[0].height = 250000
        table.columns[1].width = 1800000

        table_after = slide.table(shape_id)
        assert table_after.rows[0].height == 250000  # noqa: S101
        assert table_after.columns[1].width == 1800000  # noqa: S101
