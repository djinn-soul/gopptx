"""Procedural row/column growth tests (upstream python-pptx issue #1016)."""

import pathlib

import pytest
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches


def _table(pres: Presentation):
    slide = pres.slides[0]
    shape_id = next(
        s.id for s in slide.shapes if s.shape_type in {"tbl", "graphicFrame"}
    )
    return slide.table(shape_id)


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "grow.pptx"
    with Presentation.new("Grow table") as pres:
        pres.slides[0].add_table(
            2,
            2,
            (Inches(0.6), Inches(1.4), Inches(8.6), Inches(1.2)),
            data=[["r0c0", "r0c1"], ["r1c0", "r1c1"]],
        )
        pres.save(output_path)
    return output_path


def test_rows_and_columns_grow_procedurally(deck: pathlib.Path) -> None:
    """The issue's script grew a table with rows.add() / columns.add()."""
    pres = Presentation()
    pres.open(deck)
    table = _table(pres)
    assert (table.row_count, table.col_count) == (2, 2)

    for _ in range(2):
        table.rows.add()
    for _ in range(2):
        table.columns.add(Inches(1.4))

    table = _table(pres)
    assert (table.row_count, table.col_count) == (4, 4)
    for row in range(table.row_count):
        for col in range(table.col_count):
            table.cell(row, col).text = f"r{row}c{col}"
    pres.save(deck)
    pres.close()

    reopened = Presentation()
    reopened.open(deck)
    try:
        table = _table(reopened)
        assert (table.row_count, table.col_count) == (4, 4)
        assert table.cell(0, 0).text == "r0c0"
        assert table.cell(3, 3).text == "r3c3"
    finally:
        reopened.close()


def test_columns_add_defaults_to_last_column_width(deck: pathlib.Path) -> None:
    """python-pptx 0.6 allowed a bare columns.add(); width is optional here too."""
    pres = Presentation()
    pres.open(deck)
    try:
        table = _table(pres)
        last_width = table.columns[table.col_count - 1].width

        table.columns.add()
        table = _table(pres)

        assert table.col_count == 3
        assert table.columns[2].width == last_width
    finally:
        pres.close()


def test_insert_and_remove_round_trip(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    try:
        table = _table(pres)
        table.columns.insert(1)
        assert _table(pres).col_count == 3
        table = _table(pres)
        table.columns.remove(1)
        assert _table(pres).col_count == 2

        table = _table(pres)
        table.rows.insert(1)
        assert _table(pres).row_count == 3
        table = _table(pres)
        table.rows.remove(1)
        assert _table(pres).row_count == 2
    finally:
        pres.close()
