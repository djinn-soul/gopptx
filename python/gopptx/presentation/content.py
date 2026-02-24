"""Presentation content mixin aggregator for gopptx library."""

from __future__ import annotations

from .comments_charts import (
    PresentationChartMixin,
    PresentationCommentMixin,
)
from .shapes_tables import (
    PresentationNotesMixin,
    PresentationShapeMixin,
    PresentationTableMixin,
    PresentationTextMixin,
)


class PresentationContentMixin(
    PresentationTableMixin,
    PresentationTextMixin,
    PresentationShapeMixin,
    PresentationNotesMixin,
    PresentationChartMixin,
    PresentationCommentMixin,
):
    """Mixin providing content manipulation methods for Presentation."""
