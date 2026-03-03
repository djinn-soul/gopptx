"""Placeholder proxy class for gopptx library."""

from __future__ import annotations

from collections import UserString
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .slide import Slide


class PlaceholderFormat(UserString):
    """String-compatible placeholder format with python-pptx-like attributes."""

    def __init__(self, value: str, idx: int) -> None:
        """Initialize the placeholder format payload."""
        super().__init__(value)
        self.idx_value = idx

    @property
    def type(self) -> str:
        """Return the placeholder type token."""
        return str(self.data)

    @property
    def idx(self) -> int:
        """Return the placeholder index."""
        return self.idx_value


class Placeholder:
    """Proxy object for a placeholder within a slide."""

    def __init__(self, slide: Slide, index: int, ph_type: str, name: str) -> None:
        """Initialize the placeholder proxy.

        Args:
            slide: The Parent Slide proxy object.
            index: The zero-based index of the placeholder.
            ph_type: The placeholder type (e.g., 'body', 'title', 'pic').
            name: The human-readable name of the placeholder.
        """
        self._slide = slide
        self._index = index
        self._type = ph_type
        self._name = name

    @property
    def idx(self) -> int:
        """The index of this placeholder."""
        return self._index

    @property
    def placeholder_format(self) -> PlaceholderFormat:
        """The placeholder format object (string-compatible)."""
        return PlaceholderFormat(self._type, self._index)

    @property
    def name(self) -> str:
        """The name of this placeholder."""
        return self._name

    def insert_text(self, text: str, **style_kwargs: object) -> None:
        """Replace the placeholder with text.

        Args:
            text: The text to insert.
            **style_kwargs: Optional text style properties (size_pt, bold, italic, color, font).
        """
        # Normalize style keys (e.g. size -> size_pt, colour -> color)
        text_style = {}
        for k, v in style_kwargs.items():
            key = k
            if k in {"size", "font_size"}:
                key = "size_pt"
            elif k == "font_name":
                key = "font"
            elif k == "colour":  # Handle British spelling
                key = "color"
            text_style[key] = v

        self._slide.set_placeholder_content(
            self.idx, self._type, text=text, text_style=text_style
        )

    def insert_picture(
        self,
        image_path: str,
        bounds: tuple[float, float, float, float] | None = None,
    ) -> None:
        """Replace the placeholder with a picture.

        Args:
            image_path: Path to the image file.
            bounds: Optional (x, y, width, height) in points relative to the placeholder.
        """
        self._slide.set_placeholder_content(
            self.idx, self._type, image_path=image_path, bounds=bounds
        )

    def __repr__(self) -> str:
        """Return a string representation of this placeholder."""
        return f"<Placeholder idx={self.idx} type='{self.placeholder_format}' name='{self.name}'>"


class TitlePlaceholder(Placeholder):
    """Placeholder subtype for title-like placeholders."""


class BodyPlaceholder(Placeholder):
    """Placeholder subtype for body/content placeholders."""


class PicturePlaceholder(Placeholder):
    """Placeholder subtype for picture placeholders."""


class ChartPlaceholder(Placeholder):
    """Placeholder subtype for chart placeholders."""


class TablePlaceholder(Placeholder):
    """Placeholder subtype for table placeholders."""


_PLACEHOLDER_TYPE_TO_CLASS: dict[str, type[Placeholder]] = {
    "title": TitlePlaceholder,
    "ctrTitle": TitlePlaceholder,
    "body": BodyPlaceholder,
    "obj": BodyPlaceholder,
    "pic": PicturePlaceholder,
    "chart": ChartPlaceholder,
    "tbl": TablePlaceholder,
}


def create_placeholder(
    slide: Slide, index: int, ph_type: str, name: str
) -> Placeholder:
    """Create a placeholder proxy using the most-specific subtype mapping."""
    cls = _PLACEHOLDER_TYPE_TO_CLASS.get(ph_type, Placeholder)
    return cls(slide, index, ph_type, name)
