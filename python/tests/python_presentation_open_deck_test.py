"""Presentation.open_deck opens an existing deck like the constructor does."""

from gopptx.presentation.presentation import Presentation


def test_open_deck_opens_saved_file(tmp_path):
    path = str(tmp_path / "open_deck.pptx")
    with Presentation.new("Open Deck") as pres:
        pres.slides[0].title = "Opened Through open_deck"
        pres.save(path)

    with Presentation.open_deck(path) as reopened:
        assert len(reopened.slides) == 1
        title_ph = reopened.slides[0].placeholders.title
        assert title_ph is not None
        assert title_ph.text == "Opened Through open_deck"


def test_open_deck_matches_constructor(tmp_path):
    path = str(tmp_path / "open_deck_parity.pptx")
    with Presentation.new("Parity") as pres:
        pres.save(path)

    with Presentation.open_deck(path) as from_helper, Presentation(path) as from_ctor:
        assert len(from_helper.slides) == len(from_ctor.slides)
