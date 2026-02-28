"""Type definitions for gopptx library."""

from __future__ import annotations

try:
    from typing import TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import TypedDict


class SlideSize(TypedDict):
    """Slide dimensions in EMUs."""

    width: int
    height: int


class PresentationMetadata(TypedDict):
    """Presentation metadata."""

    title: str
    slide_count: int
    size: SlideSize


class CoreProperties(TypedDict, total=False):
    """Core document properties."""

    title: str
    subject: str
    creator: str
    keywords: str
    description: str
    lastModifiedBy: str
    revision: str
    created: str
    modified: str
    category: str
    contentStatus: str


class ShapeProps(TypedDict, total=False):
    """Shape properties."""

    name: str


class ShapeUpdate(TypedDict, total=False):
    """Shape update parameters."""

    text: str
    x: int
    y: int
    w: int
    h: int


class Shape(TypedDict):
    """Shape information."""

    ID: int
    Name: str
    Type: str
    Text: str
    X: int
    Y: int
    W: int
    H: int


class ChartSelector(TypedDict, total=False):
    """Chart selector for identifying charts."""

    index: int
    rel_id: str


class ChartSeriesData(TypedDict, total=False):
    """Chart series data for updates."""

    name: str
    categories: list[str]
    values: list[float]
    x_values: list[float]
    y_values: list[float]
    sizes: list[float]


class ChartDataUpdate(TypedDict, total=False):
    """Chart data update payload."""

    categories: list[str]
    series: list[ChartSeriesData]


class SlideChartRef(TypedDict):
    """Reference to a chart on a slide."""

    Index: int
    RelID: str
    ChartPart: str


class SlideLayoutInfo(TypedDict):
    """Information about a slide layout."""

    Part: str
    Name: str
    MasterPart: str


class SlideMasterCloneResult(TypedDict):
    """Result of cloning a slide master."""

    MasterPart: str
    ThemePart: str
    LayoutMap: dict[str, str]


class SlideMetadata(TypedDict):
    """Metadata for a slide."""

    Index: int
    SlideID: int
    RelationshipID: str
    PartName: str
    Title: str


class Section(TypedDict):
    """Section in a presentation."""

    Name: str
    GUID: str
    SlideIDs: list[int]


class ShapeSearchQuery(TypedDict, total=False):
    """Query parameters for searching shapes."""

    name_contains: str
    type_equals: str
    text_contains: str
    case_sensitive: bool


class ShapeSearchResult(TypedDict):
    """Result of a shape search."""

    SlideIndex: int
    Shape: Shape


class Author(TypedDict):
    """Author information for comments."""

    ID: int
    Name: str
    Initials: str
    LastIndex: int


class Comment(TypedDict):
    """Comment on a slide."""

    AuthorID: int
    Text: str
    Created: str
    X: int
    Y: int
    Index: int


class BatchCommand(TypedDict, total=False):
    """Command for batch operations."""

    op: str
    payload: dict[str, object]
    request_id: str


class BatchErrorDetail(TypedDict, total=False):
    """Error details for a failed batch item."""

    code: str
    message: str
    details: dict[str, object]


class BatchItemResult(TypedDict, total=False):
    """Result of a single batch item."""

    ok: bool
    op: str
    request_id: str
    result: dict[str, object]
    error: BatchErrorDetail


class TableCellInfo(TypedDict):
    """Information about a table cell."""

    row: int
    col: int
    row_span: int
    col_span: int
    v_merge: bool
    h_merge: bool
    is_merge_origin: bool
    is_spanned: bool
    text: str


class TableInfo(TypedDict):
    """Information about a table."""

    row_count: int
    col_count: int
    first_row: bool
    first_col: bool
    last_row: bool
    last_col: bool
    band_row: bool
    band_col: bool
    cells: list[TableCellInfo]

