"""Freeform geometry, picture-fill and effective-style schema definitions."""

from __future__ import annotations

try:
    from typing import TypedDict
except ImportError:  # pragma: no cover
    from typing_extensions import TypedDict


class PictureFillCrop(TypedDict):
    """Source-rectangle insets of a picture fill, as fractions of the image."""

    left: float
    top: float
    right: float
    bottom: float


class PictureFill(TypedDict, total=False):
    """Read-only view of an image used as a shape fill."""

    rel_id: str
    image_part: str
    external: str
    mode: str
    crop: PictureFillCrop


class FreeformPoint(TypedDict, total=False):
    """One point of a freeform path.

    ``x_expr``/``y_expr`` carry the raw attribute when the coordinate is a guide
    formula rather than a number, in which case ``x``/``y`` stay 0.
    """

    x: int
    y: int
    x_expr: str
    y_expr: str


class FreeformSegment(TypedDict, total=False):
    """One drawing command of a freeform path."""

    type: str
    points: list[FreeformPoint]
    width_radius: int
    height_radius: int
    start_angle_deg: float
    swing_angle_deg: float


class FreeformPath(TypedDict, total=False):
    """One path of a custom geometry, in its own coordinate space."""

    w: int
    h: int
    fill: str
    stroke: bool
    segments: list[FreeformSegment]


class FreeformGeometry(TypedDict, total=False):
    """Custom geometry (``<a:custGeom>``) read back from a shape."""

    paths: list[FreeformPath]


class EffectiveColor(TypedDict, total=False):
    """A colour resolved through the inheritance chain."""

    rgb: str
    scheme_slot: str
    source: str


class EffectiveString(TypedDict, total=False):
    """A string value resolved through the inheritance chain."""

    value: str
    source: str


class EffectiveFloat(TypedDict, total=False):
    """A numeric value resolved through the inheritance chain."""

    value: float
    source: str


class EffectiveBool(TypedDict, total=False):
    """A flag resolved through the inheritance chain."""

    value: bool
    source: str


class EffectivePosition(TypedDict, total=False):
    """A position and extent in EMU, resolved through the inheritance chain."""

    x: int
    y: int
    w: int
    h: int
    source: str


class EffectiveShapeStyle(TypedDict, total=False):
    """How a shape actually looks after inheritance.

    Each value records the ``source`` it was resolved from: ``shape``,
    ``layout``, ``master`` or ``theme``. A missing key means nothing in the
    chain defined that value.
    """

    fill_color: EffectiveColor
    font_color: EffectiveColor
    font_typeface: EffectiveString
    font_size_pt: EffectiveFloat
    bold: EffectiveBool
    italic: EffectiveBool
    position: EffectivePosition
    layout_part: str
    master_part: str
