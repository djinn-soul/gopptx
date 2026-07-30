"""Table.style, Cell.row_idx/col_idx, and full-width merge collapse (#27, #849, #636)."""

import pathlib
import zipfile

import pytest
from gopptx.presentation.presentation import Presentation
from gopptx.presentation.tables.table_styles import TableStyle
from gopptx.schemas import Inches

MEDIUM_STYLE_2 = "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"


def _table(pres: Presentation):
    slide = pres.slides[0]
    shape_id = next(
        s.id for s in slide.shapes if s.shape_type in {"tbl", "graphicFrame"}
    )
    return slide.table(shape_id)


def _make_deck(path: pathlib.Path, rows: int, cols: int) -> pathlib.Path:
    data = [[f"r{r}c{c}" for c in range(cols)] for r in range(rows)]
    with Presentation.new("Table styling") as pres:
        pres.slides[0].add_table(
            rows,
            cols,
            (Inches(0.8), Inches(1.5), Inches(8.0), Inches(3.0)),
            data=data,
        )
        pres.save(path)
    return path


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    return _make_deck(tmp_path / "table.pptx", 2, 2)


def _slide_xml(deck_path: pathlib.Path) -> str:
    with zipfile.ZipFile(deck_path) as archive:
        return archive.read("ppt/slides/slide1.xml").decode("utf-8")


def test_style_round_trips_by_guid(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    table = _table(pres)
    table.style = MEDIUM_STYLE_2
    assert table.style == MEDIUM_STYLE_2
    pres.save(deck)
    pres.close()

    assert f"<a:tableStyleId>{MEDIUM_STYLE_2}</a:tableStyleId>" in _slide_xml(deck)


def test_style_accepts_named_constant(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    table = _table(pres)
    table.style = TableStyle.MEDIUM_STYLE_2
    assert table.style == TableStyle.MEDIUM_STYLE_2
    pres.close()


def test_style_rejects_unknown_name(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    with pytest.raises(ValueError, match="Unknown style name"):
        _table(pres).style = "NOT_A_STYLE"
    pres.close()


def test_cell_row_and_col_idx_aliases(deck: pathlib.Path) -> None:
    """python-pptx spells these row_idx/col_idx (#849)."""
    pres = Presentation()
    pres.open(deck)
    cell = _table(pres).cell(1, 0)
    assert (cell.row_idx, cell.col_idx) == (cell.row, cell.col) == (1, 0)
    pres.close()


def test_full_width_merge_removes_merged_rows(tmp_path: pathlib.Path) -> None:
    """A vMerge spanning every column renders unmerged in PowerPoint (#636)."""
    path = _make_deck(tmp_path / "single_col.pptx", 4, 1)
    pres = Presentation()
    pres.open(path)
    table = _table(pres)
    table[1:3, 0:1].merge()
    table.invalidate_cache()

    assert table.row_count == 3
    assert table.cell(1, 0).row_span == 1
    pres.save(path)
    pres.close()

    assert "vMerge" not in _slide_xml(path)


def test_partial_width_merge_keeps_span(tmp_path: pathlib.Path) -> None:
    """A merge that leaves a column outside the range renders fine, so it is kept."""
    path = _make_deck(tmp_path / "two_col.pptx", 3, 2)
    pres = Presentation()
    pres.open(path)
    table = _table(pres)
    table[0:2, 0:1].merge()
    table.invalidate_cache()

    assert table.row_count == 3
    assert table.cell(0, 0).row_span == 2
    pres.close()
