"""Constants for gopptx library."""

from __future__ import annotations

import sys
from enum import Enum

from gopptx import shape_types as _shape_types

if sys.version_info >= (3, 11):
    from enum import StrEnum
else:

    class StrEnum(str, Enum):
        """Python 3.10 compatibility shim for ``enum.StrEnum``.

        ``Enum`` overrides ``__str__``/``__format__`` to render as
        ``ClassName.MEMBER``; ``enum.StrEnum`` restores the plain ``str``
        behaviour so members stringify to their value. Mirror that here so
        members sent over the FFI bridge carry the DrawingML token on 3.10.
        """

        __str__ = str.__str__
        __format__ = str.__format__  # type: ignore[assignment]


# Theme Presets
THEME_CORPORATE = "Corporate"
THEME_MODERN = "Modern"
THEME_VIBRANT = "Vibrant"
THEME_DARK = "Dark"
THEME_NATURE = "Nature"
THEME_TECH = "Tech"
THEME_CARBON = "Carbon"

# Shape presets live in gopptx.shape_types, generated from the Go constants so
# all 202 are exposed rather than the 40 this module used to hand-maintain.
# The names below are re-exported for the callers that already import them from
# here; new code should reach for gopptx.shape_types (or ShapeType) directly.
SHAPE_RECTANGLE = _shape_types.SHAPE_RECTANGLE
SHAPE_ROUNDED_RECTANGLE = _shape_types.SHAPE_ROUNDED_RECTANGLE
SHAPE_ELLIPSE = _shape_types.SHAPE_ELLIPSE
SHAPE_TRIANGLE = _shape_types.SHAPE_TRIANGLE
SHAPE_RIGHT_TRIANGLE = _shape_types.SHAPE_RIGHT_TRIANGLE
SHAPE_DIAMOND = _shape_types.SHAPE_DIAMOND
SHAPE_PENTAGON = _shape_types.SHAPE_PENTAGON
SHAPE_HEXAGON = _shape_types.SHAPE_HEXAGON
SHAPE_PARALLELOGRAM = _shape_types.SHAPE_PARALLELOGRAM
SHAPE_CLOUD = _shape_types.SHAPE_CLOUD
SHAPE_HEART = _shape_types.SHAPE_HEART
SHAPE_STAR_5 = _shape_types.SHAPE_STAR_5
SHAPE_STAR_6 = _shape_types.SHAPE_STAR_6
SHAPE_STAR_7 = _shape_types.SHAPE_STAR_7
SHAPE_STAR_8 = _shape_types.SHAPE_STAR_8
SHAPE_STAR_10 = _shape_types.SHAPE_STAR_10
SHAPE_STAR_12 = _shape_types.SHAPE_STAR_12
SHAPE_STAR_16 = _shape_types.SHAPE_STAR_16
SHAPE_STAR_24 = _shape_types.SHAPE_STAR_24
SHAPE_STAR_32 = _shape_types.SHAPE_STAR_32

# Flowchart Shapes
SHAPE_FLOWCHART_PROCESS = _shape_types.SHAPE_FLOWCHART_PROCESS
SHAPE_FLOWCHART_DECISION = _shape_types.SHAPE_FLOWCHART_DECISION
# The DrawingML token for the data shape is flowChartInputOutput; this module
# used to carry "flowChartData", which is not a member of ST_ShapeType and so
# drew a plain rectangle.
SHAPE_FLOWCHART_DATA = _shape_types.SHAPE_FLOWCHART_DATA
SHAPE_FLOWCHART_PREDEFINED_PROCESS = _shape_types.SHAPE_FLOWCHART_PREDEFINED_PROCESS
SHAPE_FLOWCHART_INTERNAL_STORAGE = _shape_types.SHAPE_FLOWCHART_INTERNAL_STORAGE
SHAPE_FLOWCHART_DOCUMENT = _shape_types.SHAPE_FLOWCHART_DOCUMENT
SHAPE_FLOWCHART_MULTIDOCUMENT = _shape_types.SHAPE_FLOWCHART_MULTIDOCUMENT
SHAPE_FLOWCHART_TERMINATOR = _shape_types.SHAPE_FLOWCHART_TERMINATOR
SHAPE_FLOWCHART_PREPARATION = _shape_types.SHAPE_FLOWCHART_PREPARATION
SHAPE_FLOWCHART_MANUAL_INPUT = _shape_types.SHAPE_FLOWCHART_MANUAL_INPUT
SHAPE_FLOWCHART_MANUAL_OPERATION = _shape_types.SHAPE_FLOWCHART_MANUAL_OPERATION

# Connectors
SHAPE_LINE = _shape_types.SHAPE_LINE
SHAPE_CURVED_CONNECTOR_2 = _shape_types.SHAPE_CURVED_CONNECTOR_2
SHAPE_CURVED_CONNECTOR_3 = _shape_types.SHAPE_CURVED_CONNECTOR_3
SHAPE_CURVED_CONNECTOR_4 = _shape_types.SHAPE_CURVED_CONNECTOR_4
SHAPE_CURVED_CONNECTOR_5 = _shape_types.SHAPE_CURVED_CONNECTOR_5
SHAPE_BENT_CONNECTOR_2 = _shape_types.SHAPE_BENT_CONNECTOR_2
SHAPE_BENT_CONNECTOR_3 = _shape_types.SHAPE_BENT_CONNECTOR_3
SHAPE_BENT_CONNECTOR_4 = _shape_types.SHAPE_BENT_CONNECTOR_4
SHAPE_BENT_CONNECTOR_5 = _shape_types.SHAPE_BENT_CONNECTOR_5


# ShapeType is generated from the Go constants, so it covers every preset the
# library can draw rather than the subset this module used to list by hand.
ShapeType = _shape_types.ShapeType


class ConnectorType(StrEnum):
    """Typed connector tokens for ``add_connector`` calls."""

    STRAIGHT = "straight"
    ELBOW = "elbow"
    CURVED = "curved"
    LINE = SHAPE_LINE
    CURVED_CONNECTOR_2 = SHAPE_CURVED_CONNECTOR_2
    CURVED_CONNECTOR_3 = SHAPE_CURVED_CONNECTOR_3
    CURVED_CONNECTOR_4 = SHAPE_CURVED_CONNECTOR_4
    CURVED_CONNECTOR_5 = SHAPE_CURVED_CONNECTOR_5
    BENT_CONNECTOR_2 = SHAPE_BENT_CONNECTOR_2
    BENT_CONNECTOR_3 = SHAPE_BENT_CONNECTOR_3
    BENT_CONNECTOR_4 = SHAPE_BENT_CONNECTOR_4
    BENT_CONNECTOR_5 = SHAPE_BENT_CONNECTOR_5


class LineDashStyle(StrEnum):
    """DrawingML line-dash tokens used by shape and chart line facades."""

    SOLID = "solid"
    DASH = "dash"
    DASH_DOT = "dashDot"
    LONG_DASH = "lgDash"
    LONG_DASH_DOT = "lgDashDot"
    LONG_DASH_DOT_DOT = "lgDashDotDot"
    ROUND_DOT = "sysDot"
    SQUARE_DOT = "sysDash"
    SYSTEM_DASH_DOT = "sysDashDot"
    SYSTEM_DASH_DOT_DOT = "sysDashDotDot"


MSO_LINE = LineDashStyle


# Placeholder Types
PLACEHOLDER_TITLE = "title"
PLACEHOLDER_BODY = "body"
PLACEHOLDER_CENTERED_TITLE = "ctrTitle"
PLACEHOLDER_SUBTITLE = "subTitle"
PLACEHOLDER_DATE_TIME = "dt"
PLACEHOLDER_SLIDE_NUMBER = "sldNum"
PLACEHOLDER_FOOTER = "ftr"
PLACEHOLDER_HEADER = "hdr"
PLACEHOLDER_OBJECT = "obj"
PLACEHOLDER_CHART = "chart"
PLACEHOLDER_TABLE = "tbl"
PLACEHOLDER_CLIP_ART = "clipArt"
PLACEHOLDER_DIAGRAM = "dgm"
PLACEHOLDER_MEDIA = "media"
PLACEHOLDER_SLIDE_IMAGE = "sldImg"
PLACEHOLDER_PICTURE = "pic"


class PlaceholderType(StrEnum):
    """Typed placeholder tokens for placeholder-aware APIs."""

    TITLE = PLACEHOLDER_TITLE
    BODY = PLACEHOLDER_BODY
    CENTERED_TITLE = PLACEHOLDER_CENTERED_TITLE
    SUBTITLE = PLACEHOLDER_SUBTITLE
    DATE_TIME = PLACEHOLDER_DATE_TIME
    SLIDE_NUMBER = PLACEHOLDER_SLIDE_NUMBER
    FOOTER = PLACEHOLDER_FOOTER
    HEADER = PLACEHOLDER_HEADER
    OBJECT = PLACEHOLDER_OBJECT
    CHART = PLACEHOLDER_CHART
    TABLE = PLACEHOLDER_TABLE
    CLIP_ART = PLACEHOLDER_CLIP_ART
    DIAGRAM = PLACEHOLDER_DIAGRAM
    MEDIA = PLACEHOLDER_MEDIA
    SLIDE_IMAGE = PLACEHOLDER_SLIDE_IMAGE
    PICTURE = PLACEHOLDER_PICTURE


# Alignment
ALIGN_LEFT = "l"
ALIGN_CENTER = "ctr"
ALIGN_RIGHT = "r"
ALIGN_JUSTIFY = "just"

# Vertical Alignment
VALIGN_TOP = "t"
VALIGN_CENTER = "ctr"
VALIGN_BOTTOM = "b"

# Underline Styles
UNDERLINE_NONE = "none"
UNDERLINE_SINGLE = "sng"
UNDERLINE_DOUBLE = "dbl"
UNDERLINE_DOTTED = "dotted"
UNDERLINE_DASH = "dash"


# Pattern Fill Types (Issue #1127)
class PatternType(StrEnum):
    """DrawingML preset pattern types for pattern fills."""

    PERCENT_5 = "pct5"
    PERCENT_10 = "pct10"
    PERCENT_20 = "pct20"
    PERCENT_25 = "pct25"
    PERCENT_30 = "pct30"
    PERCENT_40 = "pct40"
    PERCENT_50 = "pct50"
    PERCENT_60 = "pct60"
    PERCENT_70 = "pct70"
    PERCENT_75 = "pct75"
    PERCENT_80 = "pct80"
    PERCENT_90 = "pct90"
    CROSS = "cross"
    DIAGONAL_CROSS = "diagCross"


MSO_PATTERN = PatternType

# Slide Size Constants (EMU)
SIZE_4X3_WIDTH = 9144000
SIZE_4X3_HEIGHT = 6858000
SIZE_16X9_WIDTH = 12192000
SIZE_16X9_HEIGHT = 6858000
