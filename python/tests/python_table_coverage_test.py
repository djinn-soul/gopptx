import os  # noqa: D100
import pathlib

import pytest
from gopptx import GopptxError, Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")  # noqa: PTH118


@pytest.fixture
def table_shape():  # noqa: ANN201, D103
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("Table Coverage Test")
        shape_id = slide.add_table(4, 4, bounds=(1000, 1000, 5000, 2000))
        yield prs, slide, shape_id


def test_cell_properties(table_shape) -> None:  # noqa: ANN001, D103
    _prs, slide, shape_id = table_shape
    table = slide.table(shape_id)
    cell = table[0, 0]

    # Initial states
    assert cell.is_merge_origin is False  # noqa: S101
    assert cell.is_spanned is False  # noqa: S101
    assert cell.row_span == 1  # noqa: S101
    assert cell.col_span == 1  # noqa: S101

    # Merge and check
    table[0:2, 0:2].merge()
    assert table[0, 0].is_merge_origin is True  # noqa: S101
    assert table[0, 0].row_span == 2  # noqa: PLR2004, S101
    assert table[0, 0].col_span == 2  # noqa: PLR2004, S101
    assert table[1, 1].is_spanned is True  # noqa: S101


def test_cell_split_and_batch_error(table_shape) -> None:  # noqa: ANN001, D103
    prs, slide, shape_id = table_shape
    table = slide.table(shape_id)

    table[0:2, 0:2].merge()
    assert table[0, 0].is_merge_origin is True  # noqa: S101

    # Split
    table[0, 0].split()
    assert table[0, 0].is_merge_origin is False  # noqa: S101

    # Batch error
    with prs.batch():
        with pytest.raises(GopptxError) as exc:
            table[0, 0].split()
        assert exc.value.code == "BATCH_STRUCTURAL_CHANGE_NOT_ALLOWED"  # noqa: S101


def test_cell_range_merge_batch_error(table_shape) -> None:  # noqa: ANN001, D103
    prs, slide, shape_id = table_shape
    table = slide.table(shape_id)

    with prs.batch():
        with pytest.raises(GopptxError) as exc:
            table[0:2, 0:2].merge()
        assert exc.value.code == "BATCH_STRUCTURAL_CHANGE_NOT_ALLOWED"  # noqa: S101


def test_table_iter_and_cell_alias(table_shape) -> None:  # noqa: ANN001, D103
    _prs, slide, shape_id = table_shape
    table = slide.table(shape_id)

    cells = list(table.iter_cells())
    assert len(cells) == 16  # noqa: PLR2004, S101
    assert isinstance(cells[0].row, int)  # noqa: S101

    # Alias
    c1 = table.cell(2, 2)
    c2 = table[2, 2]
    assert c1.row == c2.row and c1.col == c2.col  # noqa: S101


def test_table_style_flags(table_shape) -> None:  # noqa: ANN001, D103
    _prs, slide, shape_id = table_shape
    table = slide.table(shape_id)

    # Header row
    table.header_row_enabled = True
    assert table.header_row_enabled is True  # noqa: S101
    table.header_row_enabled = False
    assert table.header_row_enabled is False  # noqa: S101

    # Banded rows
    table.banded_rows_enabled = True
    assert table.banded_rows_enabled is True  # noqa: S101
    table.banded_rows_enabled = False
    assert table.banded_rows_enabled is False  # noqa: S101


def test_table_indexing_errors(table_shape) -> None:  # noqa: ANN001, D103
    _prs, slide, shape_id = table_shape
    table = slide.table(shape_id)

    # Wrong tuple length
    with pytest.raises(TypeError):
        _ = table[0]  # type: ignore

    # Step error
    with pytest.raises(ValueError):
        _ = table[0:4:2, 0]

    # Out of bounds negative
    with pytest.raises(IndexError):
        _ = table[-5, 0]


def test_cell_repr(table_shape) -> None:  # noqa: ANN001, D103
    _prs, slide, shape_id = table_shape
    table = slide.table(shape_id)
    cell = table[0, 0]
    cell.text = "Hello"
    r = repr(cell)
    assert "Cell" in r  # noqa: S101
    assert "Hello" in r  # noqa: S101
