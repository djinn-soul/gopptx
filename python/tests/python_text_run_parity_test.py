"""Parity tests for run-level hyperlink facade objects."""

import zipfile
from pathlib import Path

from gopptx import Presentation, Run, RunHyperlink


def test_run_hyperlink_address_facade_round_trip(tmp_path: Path) -> None:
    """Ensure run/hyperlink facade objects serialize through bridge payloads."""
    out_path = tmp_path / "run_hyperlink_facade.pptx"

    with Presentation.new("Run Facade") as prs:
        slide = prs.slides[0]
        first = Run("OpenAI")
        first.hyperlink.address = "https://openai.com"
        first.hyperlink.tooltip = "OpenAI home"

        second = Run(" Next")
        second.hover_action = RunHyperlink(macro="HoverMacro")

        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 5000000, 1200000),
            runs=[first, second],
            text_frame={
                "margin_top": 90000,
                "margin_bottom": 90000,
                "margin_left": 180000,
                "margin_right": 180000,
                "word_wrap": True,
                "auto_fit_type": "shape",
                "vertical_align": "ctr",
            },
        )
        assert shape_id > 0

        jump_run = Run("Jump")
        jump_run.hyperlink.jump = "nextslide"
        jump_run.hover_action = RunHyperlink(macro="HoverMacro")
        slide.update_shape(shape_id, {"runs": [jump_run]})

        mixed_hyperlink = RunHyperlink(address="https://example.com")
        slide.add_shape(
            "rect",
            (1000000, 2500000, 5000000, 1200000),
            runs=[{"text": "mixed", "hyperlink": mixed_hyperlink}],
        )
        prs.save(out_path)

    with Presentation(out_path) as prs:
        texts = [shape["Text"] for shape in prs.slides[0].list_shapes()]
        assert "Jump" in texts
        assert "mixed" in texts

    with zipfile.ZipFile(out_path) as zf:
        slide_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
    assert "hlinkMouseOver" in slide_xml
    assert "ppaction://macro?name=HoverMacro" in slide_xml
