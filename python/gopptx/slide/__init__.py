"""Slide-domain modules."""
# ruff: noqa: I001,RUF022

from .freeform_builder import FreeformBuilder
from .placeholder import (
    BodyPlaceholder,
    ChartPlaceholder,
    PicturePlaceholder,
    Placeholder,
    PlaceholderFormat,
    TablePlaceholder,
    TitlePlaceholder,
)
from .placeholder_collection import PlaceholderCollection
from .slide import Slide
from .table import (
    Cell,
    CellRange,
    Table,
    TableColumn,
    TableColumns,
    TableRow,
    TableRows,
)
from .shape_proxy import ShapeCollection, ShapeProxy
from .chart_model import (
    Chart,
    ChartCollection,
    ChartLegend,
    ChartPlot,
    ChartPlots,
    ChartTitle,
    DataLabels,
)
from .text_model import (
    ShapeParagraphCollection,
    ShapeParagraphProxy,
    ShapeRunCollection,
    ShapeRunProxy,
    ShapeTextFrame,
)
from .chart_data import CategoryChartData, CategorySeries, XyChartData, XySeries
from .text_frame import TextFrameProps
from .text_paragraph import ParagraphProps
from .text_run import Run, RunHyperlink

__all__ = [
    "BodyPlaceholder",
    "Cell",
    "CellRange",
    "ChartPlaceholder",
    "CategoryChartData",
    "CategorySeries",
    "Chart",
    "ChartCollection",
    "ChartLegend",
    "ChartPlot",
    "ChartPlots",
    "ChartTitle",
    "FreeformBuilder",
    "DataLabels",
    "ParagraphProps",
    "PicturePlaceholder",
    "Placeholder",
    "PlaceholderCollection",
    "PlaceholderFormat",
    "Run",
    "RunHyperlink",
    "Slide",
    "ShapeCollection",
    "ShapeProxy",
    "ShapeTextFrame",
    "ShapeParagraphCollection",
    "ShapeParagraphProxy",
    "ShapeRunCollection",
    "ShapeRunProxy",
    "Table",
    "TableColumn",
    "TableColumns",
    "TableRow",
    "TableRows",
    "TablePlaceholder",
    "TextFrameProps",
    "TitlePlaceholder",
    "XyChartData",
    "XySeries",
]
