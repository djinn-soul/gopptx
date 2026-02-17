from __future__ import annotations

from typing import List, Optional

try:
    from typing import TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import TypedDict


class SlideSize(TypedDict):
    width: int
    height: int


class PresentationMetadata(TypedDict):
    title: str
    slide_count: int
    size: SlideSize


class CoreProperties(TypedDict, total=False):
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
    name: str


class ShapeUpdate(TypedDict, total=False):
    text: str
    x: int
    y: int
    w: int
    h: int


class Shape(TypedDict):
    ID: int
    Name: str
    Type: str
    Text: str
    X: int
    Y: int
    W: int
    H: int


class ChartSelector(TypedDict, total=False):
    index: int
    rel_id: str


class ChartSeriesData(TypedDict, total=False):
    name: str
    categories: List[str]
    values: List[float]
    x_values: List[float]
    y_values: List[float]
    sizes: List[float]


class ChartDataUpdate(TypedDict, total=False):
    categories: List[str]
    series: List[ChartSeriesData]


class SlideChartRef(TypedDict):
    Index: int
    RelID: str
    ChartPart: str


class SlideLayoutInfo(TypedDict):
    Part: str
    Name: str
    MasterPart: str


class SlideMasterCloneResult(TypedDict):
    MasterPart: str
    ThemePart: str
    LayoutMap: dict[str, str]


class SlideMetadata(TypedDict):
    Index: int
    SlideID: int
    RelationshipID: str
    PartName: str
    Title: str


class Section(TypedDict):
    Name: str
    GUID: str
    SlideIDs: List[int]


class ShapeSearchQuery(TypedDict, total=False):
    name_contains: str
    type_equals: str
    text_contains: str
    case_sensitive: bool


class ShapeSearchResult(TypedDict):
    SlideIndex: int
    Shape: Shape


class Author(TypedDict):
    ID: int
    Name: str
    Initials: str
    LastIndex: int


class Comment(TypedDict):
    AuthorID: int
    Text: str
    Created: str
    X: int
    Y: int
    Index: int
