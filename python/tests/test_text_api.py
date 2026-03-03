import os

from gopptx import Presentation
from gopptx.schemas import TextFrame, TextRun


def test_add_shape_text_frame_controls(tmp_path):
    output_path = os.path.join(tmp_path, "text_frame.pptx")

    with Presentation.new(title="Text APIs") as prs:
        slide = prs.slides[0]

        # Test word_wrap, auto_fit_type and margins
        text_frame_opts: TextFrame = {
            "margin_top": 100000,
            "margin_bottom": 100000,
            "margin_left": 200000,
            "margin_right": 200000,
            "word_wrap": True,
            "auto_fit_type": "shape",
        }

        runs: list[TextRun] = [
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
        assert shape_id > 0

        prs.save(output_path)

    # Reload to verify
    with Presentation(output_path) as prs:
        shapes = prs.slides[0].list_shapes()
        found = False
        for shape in shapes:
            if shape["Text"].startswith("Hello"):
                found = True
        assert found
