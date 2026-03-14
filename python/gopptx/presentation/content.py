"""Presentation content mixin aggregator for gopptx library."""

from __future__ import annotations

from .comments_charts import (
    PresentationChartMixin,
    PresentationCommentMixin,
)
from .shape_batch_mixin import PresentationShapeBatchMixin
from .shapes.shapes_tables import (
    PresentationNotesMixin,
    PresentationShapeMixin,
    PresentationTableMixin,
)
from .text.text_batch_mixin import PresentationTextBatchMixin
from .text.text_mixin import PresentationTextMixin
from .vba import PresentationVBAMixin


class PresentationContentMixin(
    PresentationShapeBatchMixin,
    PresentationTableMixin,
    PresentationTextBatchMixin,
    PresentationTextMixin,
    PresentationShapeMixin,
    PresentationNotesMixin,
    PresentationChartMixin,
    PresentationCommentMixin,
    PresentationVBAMixin,
):
    """Mixin providing content manipulation methods for Presentation."""
