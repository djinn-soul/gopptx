"""Selective grayscale conversion mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import (
        GrayscalePlaceholderRef,
        GrayscaleScope,
        GrayscaleShapeRef,
        GrayscaleTextRef,
    )


class PresentationGrayscaleMixin(PresentationMixinBase):
    """Methods for converting selected presentation content to grayscale."""

    def convert_to_grayscale(
        self,
        *,
        slides: list[int] | None = None,
        shapes: list[GrayscaleShapeRef] | None = None,
        text: list[GrayscaleTextRef] | None = None,
        placeholders: list[GrayscalePlaceholderRef] | None = None,
        scope: GrayscaleScope | None = None,
    ) -> None:
        """Convert selected slides, shapes, placeholders, runs, images, and backgrounds to grayscale."""
        resolved = scope or {}
        payload: dict[str, object] = {
            "colors": resolved.get("colors", True),
            "images": resolved.get("images", True),
            "backgrounds": resolved.get("backgrounds", True),
        }
        if slides is not None:
            payload["slides"] = slides
        if shapes is not None:
            payload["shapes"] = shapes
        if text is not None:
            payload["text"] = text
        if placeholders is not None:
            payload["placeholders"] = placeholders
        self.execute(ops.OP_CONVERT_TO_GRAYSCALE, payload)
