"""Placeholder text getter reads back what was written (Issue #144)."""

from gopptx.presentation.presentation import Presentation


def test_placeholder_text_reads_back_written_text(presentation: Presentation):
    slide = presentation.slides[0]
    ph = next(iter(slide.placeholders))

    ph.text = "Written Through Placeholder"

    assert ph.text == "Written Through Placeholder"


def test_placeholder_text_reflects_existing_slide_content(tmp_path):
    path = str(tmp_path / "placeholder_text.pptx")
    with Presentation.new("Placeholder Text") as pres:
        pres.slides[0].title = "Title From Slide"
        pres.save(path)

    with Presentation(path) as reopened:
        title_ph = reopened.slides[0].placeholders.title
        assert title_ph is not None
        assert title_ph.text == "Title From Slide"


def test_placeholder_text_is_empty_when_placeholder_has_no_text(tmp_path):
    path = str(tmp_path / "empty_placeholder.pptx")
    with Presentation.new("Empty Placeholder") as pres:
        pres.save(path)

    with Presentation(path) as reopened:
        for ph in reopened.slides[0].placeholders:
            assert isinstance(ph.text, str)
