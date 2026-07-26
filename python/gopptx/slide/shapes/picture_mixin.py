"""python-pptx-compatible picture insertion for slide proxies."""

from __future__ import annotations

from typing import TYPE_CHECKING


class SlidePictureMixin:
    """Add pictures while preserving shape-cache behavior."""

    if TYPE_CHECKING:

        @property
        def index(self) -> int:
            """Return the zero-based slide index."""
            ...

        def add_image(
            self,
            path: str | None,
            bounds: tuple[float, float, float, float],
            **kwargs: object,
        ) -> int:
            """Add an image using a source path and bounds."""
            ...

    def add_picture(
        self,
        image_file: str | bytes | None,
        left: float = 0,
        top: float = 0,
        width: float = 0,
        height: float = 0,
        **kwargs: object,
    ) -> int:
        """Add a picture with optional description, alt text, and title."""
        bounds = (left, top, width, height)
        if isinstance(image_file, str):
            return self.add_image(image_file, bounds, **kwargs)
        if isinstance(image_file, bytes):
            return self.add_image(None, bounds, data=image_file, **kwargs)
        if image_file is not None:
            return self.add_image(str(image_file), bounds, **kwargs)
        return self.add_image(None, bounds, **kwargs)
