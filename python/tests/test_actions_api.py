"""Action API fixture tests for click/hover jumps and validation."""

import zipfile
from pathlib import Path
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation

if TYPE_CHECKING:
    from gopptx.schemas import Hyperlink, TextRun

EXPECTED_SLIDE_COUNT = 3


def test_action_api_slide_jumps_and_macros(tmp_path: Path) -> None:
    """Action payloads persist expected hyperlinks and relationship wiring."""
    output_path = tmp_path / "actions_test.pptx"

    with Presentation.new(title="Actions API Test") as prs:
        # Create 3 slides to jump between
        slide1 = prs.slides[0]
        prs.add_slide("Slide 2")
        slide3 = prs.add_slide("Slide 3")

        # 1. Shape with a slide jump (to Slide 3)
        click_jump: Hyperlink = {
            "target_slide": 2,  # 0-indexed slide 3
            "tooltip": "Jump to Slide 3",
        }
        slide1.add_shape(
            "rect",
            (1000000, 1000000, 2000000, 1000000),
            text="Jump to Slide 3",
            click_action=click_jump,
        )

        # 2. Shape with a relative jump
        click_next: Hyperlink = {
            "jump": "nextslide",
        }
        slide1.add_shape(
            "rect",
            (4000000, 1000000, 2000000, 1000000),
            text="Next Slide",
            click_action=click_next,
        )

        # 3. Shape with a Hover Action returning a macro
        hover_macro: Hyperlink = {
            "macro": "MyCustomMacro",
            "tooltip": "Hover to run macro",
        }
        slide1.add_shape(
            "rect",
            (1000000, 3000000, 2000000, 1000000),
            text="Hover Macro",
            hover_action=hover_macro,
        )

        # 4. Text Run with a jump to First Slide
        run_hyperlink: Hyperlink = {"jump": "firstslide"}
        runs: list[TextRun] = [{"text": "Go to Beginning", "hyperlink": run_hyperlink}]
        slide3.add_shape("rect", (1000000, 1000000, 4000000, 1000000), runs=runs)

        prs.save(output_path)

    # Verify reopening doesn't crash from invalid relationships
    with Presentation(output_path) as prs:
        if prs.slide_count != EXPECTED_SLIDE_COUNT:
            raise AssertionError("expected 3 slides after save/reopen")
        # Ensure slide 1 has shapes
        shapes = prs.slides[0].list_shapes()
        expected_shape_count = 4
        if len(shapes) != expected_shape_count:
            raise AssertionError("expected four shapes on first slide")

    with zipfile.ZipFile(output_path) as zf:
        slide1_xml = zf.read("ppt/slides/slide1.xml").decode("utf-8")
        slide1_rels = zf.read("ppt/slides/_rels/slide1.xml.rels").decode("utf-8")
        slide3_xml = zf.read("ppt/slides/slide3.xml").decode("utf-8")

    if "hlinkMouseOver" not in slide1_xml:
        raise AssertionError("expected hover action link token")
    if "ppaction://macro?name=MyCustomMacro" not in slide1_xml:
        raise AssertionError("expected macro action token")
    if "ppaction://hlinksldjump" not in slide1_xml:
        raise AssertionError("expected slide jump action token")
    if "ppaction://hlinkshowjump?jump=nextslide" not in slide1_xml:
        raise AssertionError("expected relative jump token")
    if "relationships/slide" not in slide1_rels:
        raise AssertionError("expected slide relationship entry")
    if 'Target="slide3.xml"' not in slide1_rels:
        raise AssertionError("expected relationship target to slide3.xml")
    if "ppaction://hlinkshowjump?jump=firstslide" not in slide3_xml:
        raise AssertionError("expected firstslide jump token on slide 3")


def test_action_api_rejects_invalid_action_combinations() -> None:
    """Invalid action combinations fail fast with descriptive errors."""
    with Presentation.new(title="Actions Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match="mutually exclusive"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 2000000, 1000000),
                text="Invalid",
                click_action={
                    "address": "https://example.com",
                    "jump": "nextslide",
                },
            )

        with pytest.raises(Exception, match="unsupported jump target"):
            slide.add_shape(
                "rect",
                (1000000, 2000000, 2000000, 1000000),
                text="Invalid jump",
                click_action={
                    "jump": "homeslide",
                },
            )
