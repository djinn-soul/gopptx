"""Shape-domain package for presentation content operations."""

from .shape_media_mixin import PresentationShapeMediaMixin
from .shape_payload_mixin import PresentationShapePayloadMixin
from .shape_text_runs_mixin import PresentationShapeTextRunMixin
from .shapes_tables import (
    PresentationNotesMixin,
    PresentationShapeMixin,
    PresentationTableMixin,
    PresentationTextMixin,
)

__all__ = [
    "PresentationNotesMixin",
    "PresentationShapeMediaMixin",
    "PresentationShapeMixin",
    "PresentationShapePayloadMixin",
    "PresentationShapeTextRunMixin",
    "PresentationTableMixin",
    "PresentationTextMixin",
]
