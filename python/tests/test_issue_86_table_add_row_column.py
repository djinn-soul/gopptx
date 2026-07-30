"""Issue #86: Table.add_row() and Table.add_column().

python-pptx spells these on the table itself; gopptx already had
rows.add()/columns.add(). Asserts the grid really grew in the saved a:tbl, not
just that the proxy reports a larger count.
"""

from __future__ import annotations

import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path


def _slide_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/slides/slide1.xml").decode("utf-8")


def test_add_row_and_add_column_grow_the_grid(tmp_path: Path) -> None:
    """Both helpers extend the table and the new cell is addressable."""
    output_path = tmp_path / "grown.pptx"

    with Presentation.new(title="Issue 86") as prs:
        slide = prs.slides[0]
        table_id = slide.add_table(
            2, 2, bounds=(Inches(1), Inches(1), Inches(6), Inches(2))
        )
        table = slide.table(table_id)

        table.add_row()
        table.add_column(Inches(1.5))

        assert table.row_count == 3
        assert table.col_count == 3

        table.cell(2, 2).text = "GROWN_CELL"
        prs.save(output_path)

    xml = _slide_xml(output_path)
    assert xml.count("<a:tr ") == 3
    assert xml.count("<a:gridCol") == 3
    assert "GROWN_CELL" in xml


def test_add_row_returns_the_new_row(tmp_path: Path) -> None:
    """The returned proxy points at the row that was just appended."""
    with Presentation.new(title="Issue 86 row") as prs:
        slide = prs.slides[0]
        table_id = slide.add_table(
            2, 2, bounds=(Inches(1), Inches(1), Inches(6), Inches(2))
        )
        table = slide.table(table_id)

        row = table.add_row()

        assert table.row_count == 3
        assert row is table.rows[2] or row.index == table.rows[2].index


def test_add_column_defaults_to_the_last_column_width(tmp_path: Path) -> None:
    """A bare add_column() copies the trailing column width, as in python-pptx."""
    with Presentation.new(title="Issue 86 width") as prs:
        slide = prs.slides[0]
        table_id = slide.add_table(
            2, 2, bounds=(Inches(1), Inches(1), Inches(6), Inches(2))
        )
        table = slide.table(table_id)
        last_width = table.columns[table.col_count - 1].width

        table.add_column()

        assert table.col_count == 3
        assert table.columns[2].width == last_width


def test_add_row_column_are_aliases_of_the_collections(tmp_path: Path) -> None:
    """Table.add_row/add_column and rows.add/columns.add reach the same grid."""
    with Presentation.new(title="Issue 86 alias") as prs:
        slide = prs.slides[0]
        table_id = slide.add_table(
            1, 1, bounds=(Inches(1), Inches(1), Inches(4), Inches(1))
        )
        table = slide.table(table_id)

        table.add_row()
        table.rows.add()
        table.add_column()
        table.columns.add()

        assert table.row_count == 3
        assert table.col_count == 3
