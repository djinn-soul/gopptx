"""Public API exports for gopptx library."""

from __future__ import annotations

from .api_errors import GopptxError
from .builder import PresentationBuilder
from .presentation.presentation import Presentation
from .slide.chart.data import CategoryChartData, CategorySeries, XySeries
from .slide.shapes.freeform_builder import FreeformBuilder
from .slide.slide import Slide
from .slide.tables.table import Cell, CellRange, Table

__all__ = [
    "CategoryChartData",
    "CategorySeries",
    "Cell",
    "CellRange",
    "FreeformBuilder",
    "GopptxError",
    "Presentation",
    "PresentationBuilder",
    "Slide",
    "Table",
    "XySeries",
]
