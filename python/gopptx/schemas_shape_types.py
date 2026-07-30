"""Shape/text/image typed schema definitions."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .schemas_shape_style import FreeformGeometry, PictureFill
    from .schemas_text_types import Paragraph, TextFrame


class GradientStop(TypedDict, total=False):
    """Gradient stop settings."""

    position_pct: float
    color: str
    transparency: float


class GradientFill(TypedDict, total=False):
    """Linear gradient settings."""

    angle_deg: float
    stops: list[GradientStop]


class PatternFill(TypedDict, total=False):
    """Pattern fill settings."""

    preset: str
    fg_color: str
    bg_color: str


class FillFormat(TypedDict, total=False):
    """Shape fill settings."""

    solid: str
    transparency: float
    background: bool
    gradient: GradientFill
    pattern: PatternFill
    picture: PictureFill


class LineFormat(TypedDict, total=False):
    """Shape line settings."""

    color: str
    width_emu: int
    dash_style: str
    start_arrow: str
    start_arrow_width: str
    start_arrow_length: str
    end_arrow: str
    end_arrow_width: str
    end_arrow_length: str


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


class ImageMetadata(TypedDict):
    """Basic image properties returned by the bridge."""

    width: int
    height: int
    format: str
    hash: NotRequired[str]
    content_type: NotRequired[str]
    # Where the bytes live, so a caller can address the same image again.
    rel_id: NotRequired[str]
    part_path: NotRequired[str]


class SlideImageRef(TypedDict):
    """Reference to an image embedded in a slide."""

    Index: int
    RelID: str
    Target: str


class SlideMediaRef(TypedDict, total=False):
    """Reference to one media relationship on a slide: image, audio or video."""

    index: int
    rel_id: str
    kind: str
    target: str
    part_path: str
    content_type: str
    size_bytes: int
    external: bool


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


class ShapeTextParagraph(TypedDict, total=False):
    """One shape-text paragraph with independent runs and properties."""

    runs: list[TextRun]
    paragraph: Paragraph


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


class ShapeAdjustment(TypedDict, total=False):
    """One adjustment as read back from a shape's preset geometry."""

    Name: str
    Formula: str


class ShapeAdjustmentValue(TypedDict, total=False):
    """One preset-geometry adjustment: a yellow handle in PowerPoint's UI.

    ``value`` is a fraction (0.5 is the halfway point); ``formula`` overrides it
    with a raw OOXML guide expression.
    """

    name: str
    value: float
    formula: str


class ShapeUpdate(TypedDict, total=False):
    """Shape update parameters."""

    text: str
    runs: list[TextRun]
    paragraphs: list[ShapeTextParagraph]
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
    transparent_color: str
    rotation: float
    flip_h: bool
    flip_v: bool
    x: int
    y: int
    w: int
    h: int
    description: str
    alt_text: str
    title: str


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
    PlaceholderIndex: int | None
    PlaceholderType: str
    fill: FillFormat
    line: LineFormat
    shadow: ShadowFormat
    flip_h: bool
    flip_v: bool
    freeform: FreeformGeometry
    Shapes: list[Shape]
