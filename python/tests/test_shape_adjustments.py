"""Setting preset-geometry adjustments, the yellow handles (Issue #1017)."""

import pathlib

import pytest
from gopptx import Presentation
from gopptx.api_errors import GopptxError


def test_adjustment_is_settable_and_completes_the_guide_set(
    tmp_path: pathlib.Path,
) -> None:
    """A named handle is written, and the preset's other guide gets its default.

    PowerPoint refuses to open a file whose round2SameRect carries only adj1,
    so the whole guide set is always written.
    """
    output_path = tmp_path / "shape_adjustments.pptx"

    with Presentation.new(title="Adjustments") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_shape(
            "round2SameRect", (1000000, 1000000, 3000000, 2000000)
        )
        shape = slide.shape(shape_id)
        shape.style.set_adjustments([{"name": "adj1", "value": 0.25}])
        pres.save(output_path)

    with Presentation(output_path) as pres:
        adjustments = pres.slides[0].shape(shape_id).style.adjustments
        formulas = {item["Name"]: item["Formula"] for item in adjustments}
        assert formulas == {"adj1": "val 25000", "adj2": "val 0"}


def test_setting_one_adjustment_keeps_the_others(tmp_path: pathlib.Path) -> None:
    """Naming one handle leaves the value already on the shape alone."""
    output_path = tmp_path / "shape_adjustments_merge.pptx"

    with Presentation.new(title="Adjustments Merge") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_shape("round2SameRect", (0, 0, 2000000, 1000000))
        shape = slide.shape(shape_id)
        shape.style.set_adjustments([
            {"name": "adj1", "value": 0.2},
            {"name": "adj2", "value": 0.4},
        ])
        shape.style.set_adjustments([{"name": "adj2", "value": 0.9}])
        pres.save(output_path)

    with Presentation(output_path) as pres:
        adjustments = pres.slides[0].shape(shape_id).style.adjustments
        formulas = {item["Name"]: item["Formula"] for item in adjustments}
        assert formulas == {"adj1": "val 20000", "adj2": "val 90000"}


def test_preset_without_adjustments_is_refused(tmp_path: pathlib.Path) -> None:
    """A rect has no handles, and is refused rather than written and corrupted."""
    with Presentation.new(title="Adjustments Refused") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_shape("rect", (0, 0, 1000000, 1000000))
        with pytest.raises(GopptxError):
            slide.shape(shape_id).style.set_adjustments([
                {"name": "adj1", "value": 0.5}
            ])
