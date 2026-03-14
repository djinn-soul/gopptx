import os
import pathlib

import pytest
from gopptx import Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")  # noqa: PTH118
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")  # noqa: PTH118


def test_table_command_api() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("Table Command API Test")

        shape_id = slide.add_table(3, 3, bounds=(1000, 1000, 5000, 2000))
        table = slide.get_table(shape_id)

        assert table["row_count"] == 3
        assert table["col_count"] == 3

        slide.set_table_cell_text(shape_id, 0, 0, "Header 1")
        cell = slide.get_table_cell(shape_id, 0, 0)
        assert cell["text"] == "Header 1"

        slide.set_table_flags(shape_id, {"first_row": True, "band_row": False})
        table = slide.get_table(shape_id)
        assert table["first_row"] is True
        assert table["band_row"] is False

        slide.merge_table_cells(shape_id, (1, 1, 2, 2))
        origin = slide.get_table_cell(shape_id, 1, 1)
        assert origin["is_merge_origin"] is True
        assert origin["row_span"] == 2
        assert origin["col_span"] == 2

        spanned = slide.get_table_cell(shape_id, 2, 2)
        assert spanned["is_spanned"] is True

        slide.split_table_cell(shape_id, 1, 1)
        origin = slide.get_table_cell(shape_id, 1, 1)
        spanned = slide.get_table_cell(shape_id, 2, 2)
        assert origin["row_span"] == 1
        assert origin["col_span"] == 1
        assert spanned["is_spanned"] is False
