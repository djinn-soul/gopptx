"""python-pptx-compatible picture insertion for slide proxies."""

from __future__ import annotations

import os
from typing import TYPE_CHECKING

from ...presentation.shapes.image_inspection import (
    picture_bounds,
    resolve_picture_source,
)

if TYPE_CHECKING:
    from ..contracts import SlidePresentationProtocol


class SlidePictureMixin:
    """Add pictures while preserving shape-cache behavior."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

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
        if (
            isinstance(effective_source, (str, os.PathLike))
            and str(effective_source).lower().endswith(".svg")
        ) or (
            isinstance(effective_source, bytes) and b"<svg" in effective_source.lower()
        ):
            options["format"] = "svg"
            options["is_svg"] = True

        if isinstance(effective_source, bytes):
            return self.add_image(None, bounds, data=effective_source, **options)
        return self.add_image(os.fspath(effective_source), bounds, **options)

    def swap_image_by_index(
        self,
        image_index: int,
        data: bytes,
        img_format: str,
    ) -> None:
        """Swap an image on this slide by position index."""
        self._presentation.swap_image_by_index(
            self.index, image_index, data, img_format
        )

    def swap_image_by_rel_id(
        self,
        rel_id: str,
        data: bytes,
        img_format: str,
    ) -> None:
        """Swap an image on this slide by relationship ID."""
        self._presentation.swap_image_by_rel_id(self.index, rel_id, data, img_format)
