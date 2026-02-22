import os
import pathlib

import pytest
from gopptx import Presentation

project_root = pathlib.Path(
    os.path.join(pathlib.Path(__file__).parent, "../..")
).resolve()
input_deck = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")


def test_table_batch_mode() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("Table Batch Test")
        shape_id = slide.add_table(3, 3, 1000, 1000, 5000, 2000)

        # Initialize table OUTSIDE batch to fetch structure
        table = slide.table(shape_id)

        with prs.batch():
            # These should all be batched and NOT trigger read operations
            table[0, 0].text = "R0C0"
            table[0, 1].text = "R0C1"
            table[-1, -1].text = "Last"  # Test negative indexing

            # Read back from local cache should work inside batch
            assert table[0, 0].text == "R0C0"
            assert table[-1, -1].text == "Last"

        # Verify after batch
        assert table[0, 0].text == "R0C0"
        assert table[0, 1].text == "R0C1"
        assert table[2, 2].text == "Last"


def test_table_negative_indexing() -> None:
    if not pathlib.Path(input_deck).exists():
        pytest.skip("smoke sample missing")

    with Presentation(input_deck) as prs:
        slide = prs.add_slide("Table Indexing Test")
        shape_id = slide.add_table(3, 3, 1000, 1000, 5000, 2000)
        table = slide.table(shape_id)

        table[0, 0].text = "TopLeft"
        table[2, 2].text = "BottomRight"

        assert table[-3, -3].text == "TopLeft"
        assert table[-1, -1].text == "BottomRight"

        with pytest.raises(IndexError):
            _ = table[3, 0]
        with pytest.raises(IndexError):
            _ = table[-4, 0]
