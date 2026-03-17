"""Presentation content mixin aggregator for gopptx library."""

from __future__ import annotations

from .charts import PresentationChartMixin
from .comments import PresentationCommentMixin
from .custom_xml import PresentationCustomXMLMixin
from .export import PresentationExportMixin
from .notes import PresentationNotesMixin
from .shapes import PresentationShapeBatchMixin, PresentationShapeMixin
from .tables import PresentationTableMixin
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
    PresentationCustomXMLMixin,
    PresentationExportMixin,
):
    """Mixin providing content manipulation methods for Presentation."""
