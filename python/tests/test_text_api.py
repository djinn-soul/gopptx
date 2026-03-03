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
        has_hello_text = any(shape["Text"].startswith("Hello") for shape in shapes)
        if not has_hello_text:
            raise AssertionError("expected to find shape text starting with 'Hello'")
