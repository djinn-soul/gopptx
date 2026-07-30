"""find_and_replace across split runs, and its scope option."""

from __future__ import annotations

from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation

if TYPE_CHECKING:
    import pathlib

_BOX = (0.0, 0.0, 3000000.0, 500000.0)


def test_replaces_phrase_split_across_runs() -> None:
    """PowerPoint stores a typed phrase across runs; matching must still work."""
    with Presentation.new(title="Deck") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(*_BOX)
        slide.set_shape_runs(
            shape_id,
            [{"text": "Hello "}, {"text": "World", "bold": True}],
        )

        assert pres.find_and_replace("Hello World", "Hello Mars") == 1

        matched = [s for s in slide.list_shapes() if s["ID"] == shape_id]
        assert matched[0]["Text"] == "Hello Mars"


def test_scope_slides_leaves_notes_untouched() -> None:
    """The default scope is slides only, as it has always been."""
    with Presentation.new(title="Deck") as pres:
        slide = pres.slides[0]
        _ = slide.add_textbox(*_BOX, text="ACME on slide")
        slide.notes = "ACME in speaker notes"

        assert pres.find_and_replace("ACME", "NEWCO") == 1
        assert slide.notes == "ACME in speaker notes"


def test_scope_slides_and_notes_covers_notes() -> None:
    """slides+notes is what a deck-wide rename wants."""
    with Presentation.new(title="Deck") as pres:
        slide = pres.slides[0]
        _ = slide.add_textbox(*_BOX, text="ACME on slide")
        slide.notes = "ACME in speaker notes"

        assert pres.find_and_replace("ACME", "NEWCO", scope="slides+notes") == 2
        assert slide.notes == "NEWCO in speaker notes"


def test_unknown_scope_is_rejected() -> None:
    """A typo must fail rather than silently fall back to slides."""
    with (
        Presentation.new(title="Deck") as pres,
        pytest.raises(Exception, match="scope"),
    ):
        _ = pres.find_and_replace("a", "b", scope="bogus")


def test_replace_survives_a_save_round_trip(tmp_path: pathlib.Path) -> None:
    """The rewritten runs are valid enough to reopen."""
    output = tmp_path / "replaced.pptx"
    with Presentation.new(title="Deck") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_textbox(*_BOX)
        slide.set_shape_runs(shape_id, [{"text": "Q1 "}, {"text": "revenue"}])
        assert pres.find_and_replace("Q1 revenue", "Q2 revenue") == 1
        pres.save(output)

    with Presentation(output) as reopened:
        texts = [s["Text"] for s in reopened.slides[0].list_shapes()]
        assert "Q2 revenue" in texts
