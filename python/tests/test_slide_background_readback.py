"""Regression tests for reading the slide's real background XML."""

# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING

from gopptx.slide.background import SlideBackground

if TYPE_CHECKING:
    from xml.etree.ElementTree import Element


class _BackgroundSlide:
    def __init__(self) -> None:
        super().__init__()
        self.xml = (
            "<p:bg><p:bgPr><a:solidFill>"
            '<a:srgbClr val="112233"/>'
            "</a:solidFill></p:bgPr></p:bg>"
        )

    def get_background_xml(self) -> str:
        return self.xml

    def set_background(self, bg_type: str, **kwargs: object) -> None:
        assert bg_type == "solid"
        color = str(kwargs["color"])
        self.xml = (
            "<p:bg><p:bgPr><a:solidFill>"
            f'<a:srgbClr val="{color}"/>'
            "</a:solidFill></p:bgPr></p:bg>"
        )


def _find_color(element: Element) -> str:
    for child in element.iter():
        if child.tag.endswith("srgbClr"):
            return str(child.attrib["val"])
    raise AssertionError("background has no srgbClr")


def test_background_proxy_loads_and_refreshes_real_xml() -> None:
    slide = _BackgroundSlide()
    background = SlideBackground(slide)

    assert _find_color(background._element) == "112233"

    background.set_solid("AABBCC")

    assert _find_color(background._element) == "AABBCC"
    assert list(background._cSld) == [background._element]
