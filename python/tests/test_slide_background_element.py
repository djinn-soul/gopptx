"""Test slide.background._element returns <p:bg> element (Issue #1126)."""

from gopptx import Presentation
from gopptx.slide.background import SlideBackground


def test_slide_background_element_returns_bg_not_csld() -> None:
    """Slide.background._element must return the <p:bg> element, not <p:cSld> (Issue #1126)."""
    with Presentation.new(title="Background Element Test Deck") as pres:
        slide = pres.slides[0]

        bg_proxy = slide.background
        assert isinstance(bg_proxy, SlideBackground)

        bg_element = bg_proxy._element
        assert bg_element is not None

        # Local tag name must be "bg", NOT "cSld"
        tag_local = (
            bg_element.tag.split("}")[-1] if "}" in bg_element.tag else bg_element.tag
        )
        assert tag_local == "bg", (
            f"expected <p:bg> element, got <{tag_local}> (Issue #1126)"
        )

        # _cSld must return <p:cSld>
        csld_element = bg_proxy._cSld
        csld_local = (
            csld_element.tag.split("}")[-1]
            if "}" in csld_element.tag
            else csld_element.tag
        )
        assert csld_local == "cSld", f"expected <p:cSld> element, got <{csld_local}>"


def test_slide_background_fill_helpers() -> None:
    """Slide.background helper methods set solid/gradient fills."""
    with Presentation.new(title="Background Fill Test Deck") as pres:
        slide = pres.slides[0]
        slide.background.set_solid("3070B3")
        slide.background.set_gradient(["FF0000", "0000FF"], angle=90)
