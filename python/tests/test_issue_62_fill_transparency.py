"""Issue #62: read and write the transparency (alpha) of a solid fill.

OOXML stores opacity, not transparency: alpha is (1 - transparency) in
thousandths of a percent, so 40% transparent is a:alpha val="60000".
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

import pytest
from gopptx import Presentation
from gopptx.schemas import Inches

if TYPE_CHECKING:
    from pathlib import Path

    from gopptx.slide.shapes.shape_proxy import ShapeProxy

_ALPHA_RE = re.compile(r'<a:alpha val="(\d+)"')


def _slide_xml(pptx_path: Path) -> str:
    with zipfile.ZipFile(pptx_path) as zf:
        return zf.read("ppt/slides/slide1.xml").decode("utf-8")


def _filled_shape(prs: Presentation, transparency: float | None) -> ShapeProxy:
    """Add a blue rectangle, optionally with a transparency, and return it."""
    slide = prs.slides[0]
    shape_id = slide.add_shape(
        "rectangle", bounds=(Inches(1), Inches(1), Inches(3), Inches(2))
    )
    shape = slide.shape(shape_id)
    shape.fill.solid_color = "#3498DB"
    if transparency is not None:
        shape.fill.transparency = transparency
    return shape


@pytest.mark.parametrize(
    ("transparency", "expected_alpha"),
    [(0.0, 100000), (0.25, 75000), (0.4, 60000), (1.0, 0)],
)
def test_transparency_maps_to_alpha(
    tmp_path: Path, transparency: float, expected_alpha: int
) -> None:
    """transparency=t writes alpha=(1-t), in thousandths of a percent."""
    output_path = tmp_path / f"alpha_{expected_alpha}.pptx"

    with Presentation.new(title="Issue 62") as prs:
        shape = _filled_shape(prs, transparency)

        assert shape.fill.transparency == pytest.approx(transparency)
        prs.save(output_path)

    found = _ALPHA_RE.search(_slide_xml(output_path))

    assert found is not None, _slide_xml(output_path)
    assert int(found.group(1)) == expected_alpha


def test_transparency_is_absent_until_set(tmp_path: Path) -> None:
    """A plain solid fill writes no alpha element at all."""
    output_path = tmp_path / "no_alpha.pptx"

    with Presentation.new(title="Issue 62 none") as prs:
        shape = _filled_shape(prs, None)

        assert shape.fill.transparency is None
        prs.save(output_path)

    assert "<a:alpha" not in _slide_xml(output_path)


def test_transparency_survives_reload(tmp_path: Path) -> None:
    """The value reads back from a saved deck."""
    output_path = tmp_path / "alpha_reload.pptx"

    with Presentation.new(title="Issue 62 reload") as prs:
        _filled_shape(prs, 0.4)
        prs.save(output_path)

    with Presentation(str(output_path)) as reloaded:
        shape = reloaded.slides[0].shapes[-1]

        assert shape.fill.transparency == pytest.approx(0.4)


@pytest.mark.parametrize("bad", [-0.1, 1.1])
def test_out_of_range_transparency_is_rejected(bad: float) -> None:
    """Values outside 0.0..1.0 raise rather than writing a bogus alpha."""
    with Presentation.new(title="Issue 62 range") as prs:
        shape = _filled_shape(prs, None)

        with pytest.raises(ValueError, match=r"between 0\.0 and 1\.0"):
            shape.fill.transparency = bad


def test_transparency_requires_a_solid_fill() -> None:
    """Alpha only means something alongside a colour."""
    with Presentation.new(title="Issue 62 nofill") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_shape(
            "rectangle", bounds=(Inches(1), Inches(1), Inches(2), Inches(1))
        )

        with pytest.raises(ValueError, match="requires a solid fill"):
            slide.shape(shape_id).fill.transparency = 0.5
