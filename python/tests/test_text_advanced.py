from gopptx import Presentation


def test_shape_rich_text(tmp_path):
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

        assert isinstance(shape_id, int)
        assert shape_id > 0

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
        assert isinstance(action_shape, int)
        assert action_shape > 0

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
        assert isinstance(frame_shape_id, int)
        assert frame_shape_id > 0

        deck.save(str(out_file))

    # Re-open and see if it loads successfully
    with Presentation(str(out_file)) as deck_in:
        slides = deck_in.slides
        assert len(slides) == 2
        shapes = slides[1].list_shapes()
        print("DUMP SHAPES:", shapes)
        assert len(shapes) >= 2

        # Our simple get text API will just concat the raw text inside
        texts = [s.get("Text", "") for s in shapes]
        assert (
            "Normal text, Bold text, Italic text, Underlined text, Colored and Highlighted."
            in texts
        )
        assert "Updated link" in texts
        assert "Frame Test" in texts
        assert "Click me (shape layer)" in texts
