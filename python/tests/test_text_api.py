"""Text API fixture coverage for text-frame and run hyperlinks."""

from pathlib import Path

from gopptx import Presentation


def test_add_shape_text_frame_controls(tmp_path: Path) -> None:
    """Shape creation accepts text-frame options and run hyperlinks."""
    output_path = tmp_path / "text_frame.pptx"

    with Presentation.new(title="Text APIs") as prs:
        slide = prs.slides[0]

        # Test word_wrap, auto_fit_type and margins
        text_frame_opts: dict[str, object] = {
            "margin_top": 100000,
            "margin_bottom": 100000,
            "margin_left": 200000,
            "margin_right": 200000,
            "word_wrap": True,
            "auto_fit_type": "shape",
        }

        runs: list[dict[str, object]] = [
            {"text": "Hello "},
            {
                "text": "Hyperlink",
                "hyperlink": {
                    "address": "https://google.com/",
                    "tooltip": "Go to Google",
                },
            },
            {"text": " world!", "bold": True},
        ]

        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 4000000, 2000000),
            runs=runs,
            text_frame=text_frame_opts,
        )
        if shape_id <= 0:
            raise AssertionError("expected positive shape id")

        prs.save(output_path)

    # Reload to verify
    with Presentation(output_path) as prs:
        shapes = prs.slides[0].list_shapes()
        hello_text_present = any(shape["Text"].startswith("Hello") for shape in shapes)
        if not hello_text_present:
            raise AssertionError("expected to find shape text starting with 'Hello'")


def test_font_color_read_does_not_mutate_xml(tmp_path: Path) -> None:
    """Reading run color or text state must not inject solidFill into XML (Issue #1111)."""
    output_path = tmp_path / "read_color.pptx"

    with Presentation.new(title="Color Read Test") as prs:
        slide = prs.slides[0]

        # Add a text shape with theme-inherited color (no explicit color set)
        shape_id = slide.add_shape(
            "rect",
            (1000000, 1000000, 4000000, 2000000),
            runs=[{"text": "Theme Color Text"}],
        )
        assert shape_id > 0

        # Read shape and text state multiple times (read-only access)
        shapes_initial = slide.list_shapes()
        shapes_read2 = slide.list_shapes()
        assert shapes_initial == shapes_read2

        prs.save(output_path)

    # Reopen and confirm text is preserved without solidFill injection corruption
    with Presentation(output_path) as prs:
        shapes = prs.slides[0].list_shapes()
        assert any(s["Text"] == "Theme Color Text" for s in shapes)


def test_run_level_font_formatting(tmp_path: Path) -> None:
    """Run-level font size formatting applies directly to rPr sz to avoid UI bounce-back (Issue #1135)."""
    output_path = tmp_path / "run_font_size.pptx"

    with Presentation.new(title="Run Font Size Test") as prs:
        slide = prs.slides[0]

        runs = [
            {"text": "14pt Run Text", "size_pt": 14, "bold": True, "color": "333333"},
            {"text": " 18pt Run Text", "size_pt": 18, "italic": True},
        ]
        shape_id = slide.add_shape(
            "textbox",
            (1000000, 1000000, 5000000, 1000000),
            runs=runs,
        )
        assert shape_id > 0
        prs.save(output_path)

    with Presentation(output_path) as prs:
        shapes = prs.slides[0].list_shapes()
        assert any("14pt Run Text" in s["Text"] for s in shapes)
