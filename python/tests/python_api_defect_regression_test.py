"""Regressions for four API defects found while verifying the documentation.

Each of these shipped: a constant no method accepted, a return type that varied
by how the deck was opened, types reachable only through private module paths,
and read results whose key casing depended on which handler answered.
"""

from __future__ import annotations

import warnings

import gopptx
import pytest
from gopptx import Presentation

THEME_CONSTANTS = [
    gopptx.THEME_CORPORATE,
    gopptx.THEME_MODERN,
    gopptx.THEME_VIBRANT,
    gopptx.THEME_DARK,
    gopptx.THEME_NATURE,
    gopptx.THEME_TECH,
    gopptx.THEME_CARBON,
]

OFFICE_PRESETS = [
    "office",
    "office2013",
    "facet",
    "integral",
    "ion",
    "retrospect",
    "slice",
    "wisp",
]


@pytest.mark.parametrize("name", THEME_CONSTANTS)
def test_exported_theme_constants_are_accepted(name: str) -> None:
    """Every exported THEME_* constant names a theme apply_theme resolves.

    These constants were in __all__ but rejected by the only method that takes a
    theme name, because the string path reached the Office-preset resolver only.
    """
    with Presentation.new("Theme") as pres:
        pres.apply_theme(name)


@pytest.mark.parametrize("name", OFFICE_PRESETS)
def test_office_preset_names_still_accepted(name: str) -> None:
    with Presentation.new("Theme") as pres:
        pres.apply_theme(name)


def test_unknown_theme_name_is_rejected_and_lists_the_alternatives() -> None:
    with Presentation.new("Theme") as pres, pytest.raises(gopptx.GopptxError) as excinfo:
        pres.apply_theme("NotATheme")
    assert "accepted names are" in str(excinfo.value)


def test_validate_returns_a_list_and_no_false_orphans(tmp_path) -> None:
    """validate() is a list on both paths, and edits do not fake orphan slides.

    The manifest parts are rebuilt at save time, so validating the raw part set
    reported every slide added since opening as an orphan.
    """
    out = tmp_path / "deck.pptx"
    with Presentation.new("Validate") as pres:
        pres.add_slide("A")
        pres.add_slide("B")
        issues = pres.validate()
        assert isinstance(issues, list)
        assert issues == []
        pres.save(str(out))

    with Presentation(str(out)) as pres:
        issues = pres.validate()
        assert isinstance(issues, list)
        assert issues == []


@pytest.mark.parametrize(
    "name",
    ["ChartType", "PDFOptions", "HTMLOptions", "SlideLayoutType", "Inches", "Point", "Emu"],
)
def test_types_are_reachable_from_the_package_root(name: str) -> None:
    assert hasattr(gopptx, name)
    assert name in gopptx.__all__


def test_read_results_carry_both_key_spellings() -> None:
    """Responses expose snake_case aliases without dropping the Go keys."""
    with Presentation.new("Keys") as pres:
        slide = pres.add_slide("S")
        pres.add_textbox(
            slide.index,
            gopptx.Inches(1),
            gopptx.Inches(1),
            gopptx.Inches(3),
            gopptx.Inches(1),
            text="x",
        )

        shape = pres.list_shapes(slide.index)[0]
        assert shape["id"] == shape["ID"]
        assert shape["alt_text"] == shape["AltText"]

        pres.add_section("Sec", [slide.index])
        section = pres.get_sections()[0]
        # SlideIDs pluralises an acronym — the snake_case helper used to turn it
        # into "slide_i_ds".
        assert "slide_ids" in section
        assert section["slide_ids"] == section["SlideIDs"]


def test_chart_type_enum_does_not_warn() -> None:
    """Passing a ChartType member is the non-deprecated path."""
    with warnings.catch_warnings():
        warnings.simplefilter("error", DeprecationWarning)
        with Presentation.new("Chart") as pres:
            slide = pres.add_slide("S")
            pres.add_chart(
                slide.index,
                gopptx.ChartType.COLUMN,
                ["Q1", "Q2"],
                [1.0, 2.0],
                bounds=(
                    gopptx.Inches(1),
                    gopptx.Inches(1),
                    gopptx.Inches(4),
                    gopptx.Inches(3),
                ),
            )
