"""ShadowFormat parity tests (upstream python-pptx issue #130)."""

import pathlib
import re
import zipfile

import pytest
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches

EFFECT_LST = re.compile(r"<a:effectLst\s*/>|<a:effectLst>.*?</a:effectLst>", re.DOTALL)


def _deck_with_shape(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "shadow.pptx"
    with Presentation.new("Shadow Test") as pres:
        pres.slides[0].add_shape("rect", (Inches(1), Inches(1), Inches(2), Inches(1)))
        pres.save(output_path)
    return output_path


def _rect(pres: Presentation):
    return next(shape for shape in pres.slides[0].shapes if shape.shape_type == "rect")


def _effect_list_xml(path: pathlib.Path) -> str | None:
    with zipfile.ZipFile(path) as archive:
        slide_xml = archive.read("ppt/slides/slide1.xml").decode("utf-8")
    match = EFFECT_LST.search(slide_xml)
    return match.group(0) if match else None


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    return _deck_with_shape(tmp_path)


def test_shadow_inherit_defaults_to_true(deck: pathlib.Path) -> None:
    """A shape with no effect list of its own inherits from the shape style."""
    pres = Presentation()
    pres.open(deck)
    try:
        assert _rect(pres).shadow.inherit is True
    finally:
        pres.close()


def test_shadow_inherit_false_writes_empty_effect_list(deck: pathlib.Path) -> None:
    """inherit=False suppresses the inherited shadow, the ask in issue #130."""
    pres = Presentation()
    pres.open(deck)
    _rect(pres).shadow.inherit = False
    pres.save(deck)
    pres.close()

    assert _effect_list_xml(deck) == "<a:effectLst/>"

    reopened = Presentation()
    reopened.open(deck)
    try:
        assert _rect(reopened).shadow.inherit is False
    finally:
        reopened.close()


def test_explicit_shadow_after_inherit_false(deck: pathlib.Path) -> None:
    """Explicit attributes replace a suppressed shadow instead of conflicting."""
    pres = Presentation()
    pres.open(deck)
    _rect(pres).shadow.inherit = False
    pres.save(deck)
    pres.close()

    pres = Presentation()
    pres.open(deck)
    shadow = _rect(pres).shadow
    shadow.color = "FF0000"
    shadow.blur_radius = 60000
    shadow.distance = 40000
    shadow.angle = 45.0
    pres.save(deck)
    pres.close()

    effects = _effect_list_xml(deck)
    assert effects is not None
    assert '<a:outerShdw blurRad="60000" dist="40000" dir="2700000">' in effects
    assert '<a:srgbClr val="FF0000"/>' in effects

    reopened = Presentation()
    reopened.open(deck)
    try:
        shadow = _rect(reopened).shadow
        assert shadow.color == "FF0000"
        assert shadow.blur_radius == 60000
        assert shadow.distance == 40000
        assert shadow.angle == pytest.approx(45.0)
        assert shadow.inherit is False
    finally:
        reopened.close()


def test_shadow_inherit_true_clears_explicit_effects(deck: pathlib.Path) -> None:
    """Setting inherit=True hands the shape back to its style."""
    pres = Presentation()
    pres.open(deck)
    _rect(pres).shadow.color = "00FF00"
    pres.save(deck)
    pres.close()
    assert _effect_list_xml(deck) is not None

    pres = Presentation()
    pres.open(deck)
    _rect(pres).shadow.inherit = True
    pres.save(deck)
    pres.close()

    assert _effect_list_xml(deck) is None

    reopened = Presentation()
    reopened.open(deck)
    try:
        assert _rect(reopened).shadow.inherit is True
    finally:
        reopened.close()
