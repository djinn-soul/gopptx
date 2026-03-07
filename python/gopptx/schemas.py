"""Type definitions for gopptx library."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict

from . import schemas_chart_layout as _schemas_chart_layout

ChartDataUpdate = _schemas_chart_layout.ChartDataUpdate
ChartFormatUpdate = _schemas_chart_layout.ChartFormatUpdate
ChartSelector = _schemas_chart_layout.ChartSelector
ChartSeriesData = _schemas_chart_layout.ChartSeriesData
PlaceholderInfo = _schemas_chart_layout.PlaceholderInfo
SlideChartRef = _schemas_chart_layout.SlideChartRef
SlideLayoutInfo = _schemas_chart_layout.SlideLayoutInfo
SlideMasterCloneResult = _schemas_chart_layout.SlideMasterCloneResult


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


RGBColor = str  # Hex string like 'FF0000'


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
    creator: str  # python-pptx: author
    keywords: str
    description: str  # python-pptx: comments
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


class TextFrame(TypedDict, total=False):
    """Text frame settings."""

    margin_top: int
    margin_bottom: int
    margin_left: int
    margin_right: int
    word_wrap: bool
    auto_fit: bool
    auto_fit_type: str
    vertical_align: str
    orientation: str
    columns: int
    rotation: float


class Paragraph(TypedDict, total=False):
    """Paragraph settings."""

    indent: int
    hanging: int


class FillFormat(TypedDict, total=False):
    """Shape fill settings."""

    solid: str
    background: bool
    gradient: GradientFill
    pattern: PatternFill


class GradientStop(TypedDict, total=False):
    """Gradient stop settings."""

    position_pct: float
    color: str


class GradientFill(TypedDict, total=False):
    """Linear gradient settings."""

    angle_deg: float
    stops: list[GradientStop]


class PatternFill(TypedDict, total=False):
    """Pattern fill settings."""

    preset: str
    fg_color: str
    bg_color: str


class LineFormat(TypedDict, total=False):
    """Shape line settings."""

    color: str
    width_emu: int
    dash_style: str


class ShadowFormat(TypedDict, total=False):
    """Shape shadow settings."""

    inherit: bool
    color: str
    blur_emu: int
    distance_emu: int
    angle_deg: float


class GlowFormat(TypedDict, total=False):
    """Shape glow settings."""

    color: str
    radius_emu: int


class BlurFormat(TypedDict, total=False):
    """Shape blur settings."""

    radius_emu: int


class SoftEdgeFormat(TypedDict, total=False):
    """Shape soft-edge settings."""

    radius_emu: int


class ReflectionFormat(TypedDict, total=False):
    """Shape reflection settings."""

    blur_emu: int
    distance_emu: int


class ShapeProps(TypedDict, total=False):
    """Shape properties."""

    name: str
    text_frame: TextFrame
    paragraph: Paragraph
    fill: FillFormat
    line: LineFormat
    shadow: ShadowFormat
    glow: GlowFormat
    blur: BlurFormat
    soft_edge: SoftEdgeFormat
    reflection: ReflectionFormat
    hover_action: Hyperlink
    crop: ImageCrop
    rotation: float
    flip_h: bool
    flip_v: bool


class ImageMetadata(TypedDict):
    """Basic image properties returned by the bridge."""

    width: int
    height: int
    format: str
    hash: NotRequired[str]


class ImageCrop(TypedDict, total=False):
    """Cropping offsets (0.0 to 1.0)."""

    left: float
    right: float
    top: float
    bottom: float


class Hyperlink(TypedDict, total=False):
    """Hyperlink properties."""

    address: str
    action: str
    tooltip: str
    target_slide: int
    jump: str
    macro: str
    history: bool
    highlight_click: bool
    end_sound: bool


class TextRun(TypedDict, total=False):
    """Text run properties."""

    text: str
    bold: bool
    italic: bool
    underline: str
    strikethrough: str
    subscript: bool
    superscript: bool
    color: str
    highlight: str
    font: str
    size_pt: int
    code: bool
    all_caps: bool
    small_caps: bool
    hyperlink: Hyperlink
    hover_action: Hyperlink


class ShapeUpdate(TypedDict, total=False):
    """Shape update parameters."""

    text: str
    runs: list[TextRun]
    text_frame: TextFrame
    paragraph: Paragraph
    fill: FillFormat
    line: LineFormat
    shadow: ShadowFormat
    glow: GlowFormat
    blur: BlurFormat
    soft_edge: SoftEdgeFormat
    reflection: ReflectionFormat
    click_action: Hyperlink
    crop: ImageCrop
    rotation: float
    flip_h: bool
    flip_v: bool
    x: int
    y: int
    w: int
    h: int


class Shape(TypedDict, total=False):
    """Shape information."""

    ID: int
    Name: str
    Type: str
    Text: str
    X: int
    Y: int
    W: int
    H: int
    fill: FillFormat
    line: LineFormat
    shadow: ShadowFormat


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
