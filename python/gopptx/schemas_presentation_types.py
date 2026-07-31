"""Presentation/document/search/table typed schema definitions."""

from __future__ import annotations

from typing import TYPE_CHECKING

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict

if TYPE_CHECKING:
    from .schemas_shape_types import Shape


def emu(value: int) -> int:
    """English Metric Unit."""
    return int(value)


def inches(value: float) -> int:
    """Inches to EMUs."""
    return int(value * 914400)


def point(value: float) -> int:
    """Points to EMUs."""
    return int(value * 12700)


Emu = emu
Inches = inches
Point = point


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
    """Core document properties (matching python-pptx CoreProperties)."""

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
    identifier: str
    language: str
    lastPrinted: str
    version: str


class SlideMetadata(TypedDict, total=False):
    """Metadata for a slide."""

    Index: int
    SlideID: int
    RelationshipID: str
    PartName: str
    Title: str
    is_hidden: bool
    hidden: bool


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
    margin_left: NotRequired[int]
    margin_right: NotRequired[int]
    margin_top: NotRequired[int]
    margin_bottom: NotRequired[int]


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
    style_id: str
    cells: list[TableCellInfo]
