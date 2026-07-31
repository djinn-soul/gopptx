"""Shapes-domain package for slide object-model helpers."""

from .freeform_builder import FreeformBuilder
from .shape_collection import ShapeCollection
from .shape_proxy import (
    ShapeFillProxy,
    ShapeLineProxy,
    ShapeProxy,
    ShapeShadowProxy,
)
from .shape_text_frame import ShapeTextFrame

__all__ = [
    "FreeformBuilder",
    "ShapeCollection",
    "ShapeFillProxy",
    "ShapeLineProxy",
    "ShapeProxy",
    "ShapeShadowProxy",
    "ShapeTextFrame",
]
