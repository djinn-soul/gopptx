"""Presentation content mixin aggregator for gopptx library."""

from __future__ import annotations

from .charts import PresentationChartMixin
from .comments import PresentationCommentMixin
from .notes_mixin import PresentationNotesMixin
from .shape_batch_mixin import PresentationShapeBatchMixin
from .shapes.shapes_tables import PresentationShapeMixin
from .table_mixin import PresentationTableMixin
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
