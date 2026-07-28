"""Presentation content mixin aggregator for gopptx library."""

from __future__ import annotations

from .charts import PresentationChartMixin
from .color_mixin import PresentationColorMixin
from .comments import PresentationCommentMixin
from .custom_xml import PresentationCustomXMLMixin
from .export import PresentationExportMixin
from .grayscale import PresentationGrayscaleMixin
from .headers_footers_mixin import PresentationHeaderFooterMixin
from .layout_theme_mixin import PresentationThemeMixin
from .notes import PresentationNotesMixin
from .shapes import PresentationShapeBatchMixin, PresentationShapeMixin
from .tables import PresentationTableMixin
from .tables.table_builders import PresentationTableBuilders
from .text.text_batch_mixin import PresentationTextBatchMixin
from .text.text_mixin import PresentationTextMixin
from .vba import PresentationVBAMixin


class PresentationContentMixin(
    PresentationHeaderFooterMixin,
    PresentationThemeMixin,
    PresentationColorMixin,
    PresentationShapeBatchMixin,
    PresentationTableMixin,
    PresentationTableBuilders,
    PresentationTextBatchMixin,
    PresentationTextMixin,
    PresentationShapeMixin,
    PresentationNotesMixin,
    PresentationChartMixin,
    PresentationCommentMixin,
    PresentationGrayscaleMixin,
    PresentationVBAMixin,
    PresentationCustomXMLMixin,
    PresentationExportMixin,
):
    """Mixin providing content manipulation methods for Presentation."""
