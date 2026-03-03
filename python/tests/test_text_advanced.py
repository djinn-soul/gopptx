"""Advanced text API integration tests."""

import pathlib

from gopptx import Presentation

EXPECTED_SLIDE_COUNT = 2
MIN_EXPECTED_SHAPES = 2


def _assert_positive_shape_id(value: object, label: str) -> None:
    if not isinstance(value, int):
        raise AssertionError(f"{label} must be int")
    if value <= 0:
        raise AssertionError(f"{label} must be positive")


def _assert_expected_texts(texts: list[str]) -> None:
    expected_text = (
        "Normal text, Bold text, Italic text, Underlined text, Colored and Highlighted."
    )
    if expected_text not in texts:
        raise AssertionError("expected rich text not found")
    if "Updated link" not in texts:
        raise AssertionError("expected updated link text not found")
    if "Frame Test" not in texts:
        raise AssertionError("expected frame text not found")
    if "Click me (shape layer)" not in texts:
        raise AssertionError("expected action shape text not found")


def test_shape_rich_text(tmp_path: pathlib.Path) -> None:
    """Test adding text runs with formatting and hyperlinks."""
    out_file = tmp_path / "rich_text.pptx"

    with Presentation.new("Rich Text Deck") as deck:
        slide = deck.add_slide("Rich Text Test")

        # Add a shape with multiple text runs
        shape_id = slide.add_shape(
            "rect",
            (100, 100, 400, 100),
            text="",  # Ignored since runs is provided
            runs=[
                {"text": "Normal text, "},
                {"text": "Bold text, ", "bold": True},
                {"text": "Italic text, ", "italic": True},
                {"text": "Underlined text, ", "underline": "sng"},
                {"text": "Colored ", "color": "FF0000"},
                {"text": "and Highlighted.", "highlight": "00FF00"},
            ],
        )

        _assert_positive_shape_id(shape_id, "shape_id")

        # Add another shape with a hyperlink
        link_shape_id = slide.add_shape(
            "rect",
            (100, 250, 400, 50),
            runs=[
                {
                    "text": "Click here to search",
                    "hyperlink": {
                        "address": "https://google.com",
                        "tooltip": "Search the web",
                    },
                }
            ],
        )

        # Also test update with runs
        slide.update_shape(
            link_shape_id,
            {
                "runs": [
                    {
                        "text": "Updated link",
                        "hyperlink": {
                            "address": "https://bing.com",
                            "tooltip": "Use Bing",
                        },
                    }
                ]
            },
        )

        # Test click_action directly on a shape
        action_shape = slide.add_shape(
            "rect",
            (100, 350, 400, 100),
            text="Click me (shape layer)",
            click_action={
                "address": "https://github.com",
                "tooltip": "Go to Github",
            },
        )
        _assert_positive_shape_id(action_shape, "action_shape")

        # Test Text Frame integration (margins, autofit, word wrap)
        frame_shape_id = slide.add_shape(
            "rect",
            (500, 100, 200, 200),
            text="Frame Test",
            text_frame={
                "margin_top": 91440,
                "margin_bottom": 91440,
                "margin_left": 182880,
                "margin_right": 182880,
                "word_wrap": False,
                "auto_fit": True,
                "vertical_align": "ctr",
            },
        )
        _assert_positive_shape_id(frame_shape_id, "frame_shape_id")

        deck.save(str(out_file))

    # Re-open and see if it loads successfully
    with Presentation(str(out_file)) as deck_in:
        slides = deck_in.slides
        if len(slides) != EXPECTED_SLIDE_COUNT:
            raise AssertionError("expected 2 slides")
        shapes = slides[1].list_shapes()
        if len(shapes) < MIN_EXPECTED_SHAPES:
            raise AssertionError("expected at least 2 shapes")

        # Our simple get text API will just concat the raw text inside
        texts = [s.get("Text", "") for s in shapes]
        _assert_expected_texts(texts)
