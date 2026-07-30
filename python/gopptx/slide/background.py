"""Slide background proxy for python-pptx parity (Issue #1126)."""

from __future__ import annotations

import base64
from typing import Protocol, cast

# Element only constructs trusted namespace-tagged nodes; it never parses input.
from xml.etree.ElementTree import Element  # nosec B405

from defusedxml.ElementTree import fromstring


class _BackgroundSlideProtocol(Protocol):
    def set_background(self, bg_type: str, **kwargs: object) -> None:
        """Set the slide background."""
        ...

    def get_background_xml(self) -> str:
        """Return the slide's current p:bg XML."""
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
        self._background_element: Element | None = None
        self._common_slide_data: Element | None = None

    def _refresh_element(self) -> None:
        """Drop the parsed subtree, so the next read fetches the slide again."""
        self._background_element = None
        self._common_slide_data = None

    def _load_element(self) -> None:
        """Parse the trusted background subtree returned by the local bridge."""
        namespace = "http://schemas.openxmlformats.org/presentationml/2006/main"
        self._common_slide_data = Element(f"{{{namespace}}}cSld")
        self._background_element = self._parse_background(namespace)
        self._common_slide_data.append(self._background_element)

    def _parse_background(self, namespace: str) -> Element:
        background_xml = self._slide.get_background_xml()
        if not background_xml:
            return Element(f"{{{namespace}}}bg")
        wrapper = (
            '<root xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" '
            'xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" '
            'xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">'
            f"{background_xml}</root>"
        )
        root = fromstring(wrapper)
        if len(root) == 0:
            return Element(f"{{{namespace}}}bg")
        return root[0]

    @property
    def _element(self) -> Element:
        """Return the slide's <p:bg> background element (Issue #1126).

        The subtree is fetched on first read rather than in __init__, so that
        building the proxy — which ``Slide.background`` does on every access —
        costs nothing until someone looks at the XML.

        Returns:
            The <p:bg> XML Element for this slide.
        """
        if self._background_element is None:
            self._load_element()
        return cast("Element", self._background_element)

    @property
    def _cSld(self) -> Element:
        """Return the slide's <p:cSld> container element."""
        if self._common_slide_data is None:
            self._load_element()
        return cast("Element", self._common_slide_data)

    def set_solid(self, color: str) -> None:
        """Set solid background color (e.g. 'FF0000')."""
        self._slide.set_background("solid", color=color)
        self._refresh_element()

    def solid(self, color: str = "FFFFFF") -> None:
        """python-pptx alias for set_solid."""
        self.set_solid(color)

    def set_gradient(self, colors: list[str], angle: int = 0) -> None:
        """Set gradient background colors and angle."""
        self._slide.set_background("gradient", colors=colors, angle=angle)
        self._refresh_element()

    def set_picture(self, image_path_or_bytes: str | bytes | object) -> None:
        """Set full-slide picture background from file path, bytes, or Path (Issue #1095)."""
        if isinstance(image_path_or_bytes, bytes):
            b64_data = base64.b64encode(image_path_or_bytes).decode("ascii")
            self._slide.set_background("image", image_data=b64_data)
        else:
            self._slide.set_background("image", image_path=str(image_path_or_bytes))
        self._refresh_element()
