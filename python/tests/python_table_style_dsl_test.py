from __future__ import annotations

from gopptx import Presentation


def test_define_and_list_table_styles() -> None:
    with Presentation.new("Table Style DSL") as prs:
        style_id = prs.define_table_style("My Programmatic Style")
        assert style_id.startswith("{")
        styles = prs.list_table_styles()
        assert any(s["style_id"] == style_id for s in styles)


def test_apply_defined_table_style() -> None:
    with Presentation.new("Table Style Apply") as prs:
        slide = prs.add_slide("Table")
        shape_id = slide.add_table(2, 2, (0, 0, 2000000, 1000000))
        style_id = prs.define_table_style("Apply Style")
        slide.set_table_style(shape_id, style_id)
