"""python-pptx-compatible picture insertion for slide proxies."""

from __future__ import annotations

import os
from typing import TYPE_CHECKING

from ...presentation.shapes.image_inspection import (
    picture_bounds,
    resolve_picture_source,
)


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
        image_file: str | bytes | os.PathLike[str] | None,
        left: float = 0,
        top: float = 0,
        width: float = 0,
        height: float = 0,
        **kwargs: object,
    ) -> int:
        """Add a picture with optional description, alt text, and title."""
        effective_source = resolve_picture_source(image_file, kwargs)
        bounds = picture_bounds(effective_source, left, top, width, height)
        options = {
            key: value for key, value in kwargs.items() if key not in {"path", "data"}
        }
        if isinstance(effective_source, bytes):
            return self.add_image(None, bounds, data=effective_source, **options)
        return self.add_image(os.fspath(effective_source), bounds, **options)
