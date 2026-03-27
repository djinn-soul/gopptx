"""Shape/text/image typed schema definitions."""

from __future__ import annotations

try:
    from typing import NotRequired, TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import NotRequired, TypedDict


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
    tab_stops: list[int]
    alignment: str
    level: int
    bullet_style: str
    bullet_char: str
    bullet_color: str
    bullet_size_pct: int
    line_spacing_pct: int
    line_spacing_pts: int
    space_before_pts: int
    space_after_pts: int


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


class FillFormat(TypedDict, total=False):
    """Shape fill settings."""

    solid: str
    transparency: float
    background: bool
    gradient: GradientFill
    pattern: PatternFill


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


class SlideImageRef(TypedDict):
    """Reference to an image embedded in a slide."""

    Index: int
    RelID: str
    Target: str


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
