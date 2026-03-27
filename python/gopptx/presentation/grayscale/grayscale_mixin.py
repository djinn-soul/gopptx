"""Selective grayscale conversion mixin."""

from __future__ import annotations

from ... import ops
from ..helpers import PresentationMixinBase


class PresentationGrayscaleMixin(PresentationMixinBase):
    """Methods for converting selected presentation content to grayscale."""

    def convert_to_grayscale(
        self,
        *,
        slides: list[int] | None = None,
        shapes: list[dict[str, int]] | None = None,
        text: list[dict[str, object]] | None = None,
        placeholders: list[dict[str, object]] | None = None,
        colors: bool = True,
        images: bool = True,
        backgrounds: bool = True,
    ) -> None:
        """Convert selected slides, shapes, placeholders, runs, images, and backgrounds to grayscale."""
        payload: dict[str, object] = {
            "colors": colors,
            "images": images,
            "backgrounds": backgrounds,
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
