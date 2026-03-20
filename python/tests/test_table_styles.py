"""Tests for improved table styling system."""

import pytest

from gopptx import Presentation
from gopptx.presentation.tables.table_styles import TableStyle
from gopptx.schemas import Inches


def test_table_style_constants_exist():
    """Verify TableStyle constants are defined."""
    assert TableStyle.MEDIUM_STYLE_2 == "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"
    assert TableStyle.LIGHT_STYLE_1 == "{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}"
    assert TableStyle.MEDIUM_STYLE_1 == "{5940675A-B579-460E-94D1-54222C63F5DA}"


def test_table_style_get_all():
    """Test TableStyle.get_all() returns dict of all styles."""
    styles = TableStyle.get_all()
    assert isinstance(styles, dict)
    assert "MEDIUM_STYLE_2" in styles
    assert "LIGHT_STYLE_1" in styles
    assert len(styles) > 0


def test_apply_style_by_constant():
    """Test applying style using TableStyle constant."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )

        table = prs.slides[0].table(table_id)

        # Apply by constant
        table.apply_style(TableStyle.MEDIUM_STYLE_2)
        # Should not raise


def test_apply_style_by_name_string():
    """Test applying style using style name string."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )

        table = prs.slides[0].table(table_id)

        # Apply by name
        table.apply_style("MEDIUM_STYLE_2")
        # Should not raise


def test_apply_style_by_guid():
    """Test applying style using raw GUID (backward compat)."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )

        table = prs.slides[0].table(table_id)

        # Apply by raw GUID
        table.apply_style("{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}")
        # Should not raise


def test_apply_style_invalid_name():
    """Test that invalid style name raises ValueError."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )

        table = prs.slides[0].table(table_id)

        with pytest.raises(ValueError, match="Unknown style name"):
            table.apply_style("INVALID_STYLE_NAME")


def test_presentation_set_table_style_by_name():
    """Test Presentation.set_table_style with style name."""
    with Presentation.new("Test") as prs:
        table_id = prs.add_table(
            slide=0,
            rows=2,
            cols=2,
            bounds=(Inches(1), Inches(1), Inches(5), Inches(3)),
        )

        # Apply by name via Presentation
        prs.set_table_style(0, table_id, "MEDIUM_STYLE_1")
        # Should not raise


def test_presentation_get_all_style_names():
    """Test getting all available style names."""
    with Presentation.new("Test") as prs:
        names = prs.get_all_table_style_names()
        # Names should be a list (may be empty if no styles in default template)
        assert isinstance(names, list)


def test_presentation_get_style_by_name():
    """Test finding style by name."""
    with Presentation.new("Test") as prs:
        # This will work if the presentation has any styles
        # Result may be None if no styles match
        result = prs.get_table_style_by_name("Some Style Name")
        assert result is None or isinstance(result, str)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
