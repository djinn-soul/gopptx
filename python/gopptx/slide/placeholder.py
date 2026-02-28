"""Placeholder proxy class for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from .slide import Slide


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
    def placeholder_format(self) -> str:
        """The type of this placeholder."""
        return self._type

    @property
    def name(self) -> str:
        """The name of this placeholder."""
        return self._name

    def insert_text(self, text: str, **style_kwargs: Any) -> None:
        """Replace the placeholder with text.

        Args:
            text: The text to insert.
            **style_kwargs: Optional text style properties (size_pt, bold, italic, color, font).
        """
        # Normalize style keys (e.g. size -> size_pt, colour -> color)
        text_style = {}
        for k, v in style_kwargs.items():
            key = k
            if k == "size" or k == "font_size":
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
