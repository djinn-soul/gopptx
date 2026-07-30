"""Slide background proxy for python-pptx parity (Issue #1126)."""

from __future__ import annotations

import base64
from typing import Protocol, cast

# Element only constructs trusted namespace-tagged nodes; it never parses input.
from xml.etree.ElementTree import Element  # nosec B405


class _BackgroundSlideProtocol(Protocol):
    def set_background(self, bg_type: str, **kwargs: object) -> None:
        """Set the slide background."""
        ...


class SlideBackground:
    """Proxy object representing a slide's background (Issue #1126).

    Fixes python-pptx issue #1126: ``slide.background._element`` returns the
    <p:bg> background element itself rather than <p:cSld>, preventing accidental
    deletion of <p:spTree> when copying or manipulating slide elements.
    """

    def __init__(self, slide: object) -> None:
        """Create a background proxy for a slide."""
        super().__init__()
        self._slide = cast("_BackgroundSlideProtocol", slide)
        namespace = "http://schemas.openxmlformats.org/presentationml/2006/main"
        self._common_slide_data = Element(f"{{{namespace}}}cSld")
        self._background_element = Element(f"{{{namespace}}}bg")
        self._common_slide_data.append(self._background_element)

    @property
    def _element(self) -> Element:
        """Return the slide's <p:bg> background element (Issue #1126).

        Returns:
            The <p:bg> XML Element for this slide.
        """
        return self._background_element

    @property
    def _cSld(self) -> Element:
        """Return the slide's <p:cSld> container element."""
        return self._common_slide_data

    def set_solid(self, color: str) -> None:
        """Set solid background color (e.g. 'FF0000')."""
        self._slide.set_background("solid", color=color)

    def solid(self, color: str = "FFFFFF") -> None:
        """python-pptx alias for set_solid."""
        self.set_solid(color)

    def set_gradient(self, colors: list[str], angle: int = 0) -> None:
        """Set gradient background colors and angle."""
        self._slide.set_background("gradient", colors=colors, angle=angle)

    def set_picture(self, image_path_or_bytes: str | bytes | object) -> None:
        """Set full-slide picture background from file path, bytes, or Path (Issue #1095)."""
        if isinstance(image_path_or_bytes, bytes):
            b64_data = base64.b64encode(image_path_or_bytes).decode("ascii")
            self._slide.set_background("image", image_data=b64_data)
        else:
            self._slide.set_background("image", image_path=str(image_path_or_bytes))
