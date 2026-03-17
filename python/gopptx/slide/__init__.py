"""Slide-domain modules."""

from .chart import (
    Chart,
    ChartCollection,
    ChartLegend,
    ChartPlot,
    ChartPlots,
    ChartTitle,
    DataLabels,
)
from .chart.data import CategoryChartData, CategorySeries, XyChartData, XySeries
from .placeholders.placeholder import (
    BodyPlaceholder,
    ChartPlaceholder,
    PicturePlaceholder,
    Placeholder,
    PlaceholderFormat,
    TablePlaceholder,
    TitlePlaceholder,
)
from .placeholders.placeholder_collection import PlaceholderCollection
from .shapes.freeform_builder import FreeformBuilder
from .shapes.shape_proxy import ShapeCollection, ShapeProxy
from .slide import Slide
from .tables.table import (
    Cell,
    CellRange,
    Table,
    TableColumn,
    TableColumns,
    TableRow,
    TableRows,
)
from .text.text_frame import TextFrameProps
from .text.text_model import (
    ShapeParagraphCollection,
    ShapeParagraphProxy,
    ShapeRunCollection,
    ShapeRunProxy,
    ShapeTextFrame,
)
from .text.text_paragraph import ParagraphProps
from .text.text_run import Run, RunHyperlink

__all__ = [
    "BodyPlaceholder",
    "CategoryChartData",
    "CategorySeries",
    "Cell",
    "CellRange",
    "Chart",
    "ChartCollection",
    "ChartLegend",
    "ChartPlaceholder",
    "ChartPlot",
    "ChartPlots",
    "ChartTitle",
    "DataLabels",
    "FreeformBuilder",
    "ParagraphProps",
    "PicturePlaceholder",
    "Placeholder",
    "PlaceholderCollection",
    "PlaceholderFormat",
    "Run",
    "RunHyperlink",
    "ShapeCollection",
    "ShapeParagraphCollection",
    "ShapeParagraphProxy",
    "ShapeProxy",
    "ShapeRunCollection",
    "ShapeRunProxy",
    "ShapeTextFrame",
    "Slide",
    "Table",
    "TableColumn",
    "TableColumns",
    "TablePlaceholder",
    "TableRow",
    "TableRows",
    "TextFrameProps",
    "TitlePlaceholder",
    "XyChartData",
    "XySeries",
]
