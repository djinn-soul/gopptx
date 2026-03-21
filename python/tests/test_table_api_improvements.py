"""Tests for improved table API with data loading and convenience builders."""

import pytest
from gopptx import Presentation
from gopptx.schemas import Inches


def test_add_table_new_named_parameters():
    """Test add_table with new named parameter style."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=3,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )
        assert table_id > 0

        # Verify table exists
        table = prs.slides[0].table(table_id)
        assert table.row_count == 3
        assert table.col_count == 2


def test_add_table_with_data():
    """Test add_table with data parameter."""
    with Presentation.new("Test") as prs:
        data = [
            ["Item", "Qty"],
            ["Widgets", "50"],
            ["Gadgets", "30"],
        ]
        table_id = prs.add_table(
            slide=0,
            rows=3,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
            data=data,
        )

        # Verify cells are populated
        table = prs.slides[0].table(table_id)
        assert table[0, 0].text == "Item"
        assert table[0, 1].text == "Qty"
        assert table[1, 0].text == "Widgets"
        assert table[1, 1].text == "50"
        assert table[2, 0].text == "Gadgets"
        assert table[2, 1].text == "30"


def test_add_table_with_flags():
    """Test add_table with flag parameters."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
            first_row=True,
            band_row=True,
        )

        table = prs.slides[0].table(table_id)
        assert table.header_row_enabled is True
        assert table.banded_rows_enabled is True


def test_add_table_with_all_flags():
    """Test add_table with all flag parameters."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
            first_row=True,
            first_col=True,
            last_row=True,
            last_col=True,
            band_row=True,
            band_col=True,
        )

        table = prs.slides[0].table(table_id)
        assert table.header_row_enabled is True
        assert table.first_col is True
        assert table.last_row is True
        assert table.last_col is True
        assert table.banded_rows_enabled is True
        assert table.vert_banding is True


def test_add_table_backward_compat():
    """Test that old positional signature still works."""
    with Presentation.new("Test") as prs:
        # Old style call (should still work)
        x, y, cx, cy = Inches(1), Inches(1), Inches(5), Inches(3)
        table_id = prs.add_table(0, 2, 2, (x, y, cx, cy))
        assert table_id > 0

        table = prs.slides[0].table(table_id)
        assert table.row_count == 2
        assert table.col_count == 2


def test_table_set_data():
    """Test Table.set_data() method."""
    with Presentation.new("Test") as prs:
        # Create a table with initial data
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
            data=[["A", "B"], ["C", "D"]],
        )

        table = prs.slides[0].table(table_id)

        # Replace data
        table.set_data([["X", "Y"], ["Z", "W"]])

        # Verify updated
        assert table[0, 0].text == "X"
        assert table[0, 1].text == "Y"
        assert table[1, 0].text == "Z"
        assert table[1, 1].text == "W"


def test_table_set_data_dimension_mismatch():
    """Test set_data raises on dimension mismatch."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )
        table = prs.slides[0].table(table_id)

        # Wrong number of rows
        with pytest.raises(ValueError, match="Row count mismatch"):
            table.set_data([["A", "B"]])

        # Wrong number of columns
        with pytest.raises(ValueError, match="column count mismatch"):
            table.set_data([["A"], ["B"]])


def test_add_table_from_rows():
    """Test add_table_from_rows convenience method."""
    with Presentation.new("Test") as prs:
        data = [
            ["Item", "Qty", "Status"],
            ["Widgets", "50", "Ready"],
            ["Gadgets", "30", "Queued"],
        ]
        table_id = prs.add_table_from_rows(
            slide=0,
            rows=data,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
            first_row=True,
            band_row=True,
        )

        assert table_id > 0
        table = prs.slides[0].table(table_id)
        assert table.row_count == 3
        assert table.col_count == 3
        assert table[0, 0].text == "Item"
        assert table[1, 0].text == "Widgets"
        assert table.header_row_enabled is True
        assert table.banded_rows_enabled is True


def test_add_table_from_dicts():
    """Test add_table_from_dicts convenience method."""
    with Presentation.new("Test") as prs:
        rows = [
            {"item": "Widgets", "qty": "50", "status": "Ready"},
            {"item": "Gadgets", "qty": "30", "status": "Queued"},
        ]
        table_id = prs.add_table_from_dicts(
            slide=0,
            rows=rows,
            column_names=["item", "qty", "status"],
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
            first_row=True,
            band_row=True,
        )

        assert table_id > 0
        table = prs.slides[0].table(table_id)
        assert table.row_count == 3  # header + 2 data rows
        assert table.col_count == 3

        # Header row
        assert table[0, 0].text == "item"
        assert table[0, 1].text == "qty"
        assert table[0, 2].text == "status"

        # Data rows
        assert table[1, 0].text == "Widgets"
        assert table[1, 1].text == "50"
        assert table[2, 0].text == "Gadgets"


def test_add_table_from_dicts_auto_column_names():
    """Test add_table_from_dicts with automatic column name detection."""
    with Presentation.new("Test") as prs:
        rows = [
            {"a": "1", "b": "2", "c": "3"},
            {"a": "4", "b": "5", "c": "6"},
        ]
        table_id = prs.add_table_from_dicts(
            slide=0,
            rows=rows,
            bounds=(Inches(0.8), Inches(1.5), Inches(8), Inches(2.5)),
        )

        table = prs.slides[0].table(table_id)
        # Should have auto-sorted columns: a, b, c
        assert table[0, 0].text == "a"
        assert table[0, 1].text == "b"
        assert table[0, 2].text == "c"


def test_slide_add_table_new_style():
    """Test Slide.add_table with new data parameter."""
    with Presentation.new("Test") as prs:
        slide = prs.slides[0]
        table_id = slide.add_table(
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
            data=[["A", "B"], ["C", "D"]],
            first_row=True,
        )

        table = slide.table(table_id)
        assert table[0, 0].text == "A"
        assert table[1, 1].text == "D"
        assert table.header_row_enabled is True


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
