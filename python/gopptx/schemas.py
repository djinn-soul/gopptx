"""Type definitions for gopptx library."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import override

from . import schemas_chart_layout as _schemas_chart_layout
from . import schemas_chart_series as _schemas_chart_series
from . import schemas_presentation_types as _schemas_presentation_types
from . import schemas_shape_style as _schemas_shape_style
from . import schemas_shape_types as _schemas_shape_types

if TYPE_CHECKING:
    from typing_extensions import Self

emu = _schemas_presentation_types.emu
inches = _schemas_presentation_types.inches
point = _schemas_presentation_types.point
Emu = _schemas_presentation_types.Emu
Inches = _schemas_presentation_types.Inches
Point = _schemas_presentation_types.Point

SlideSize = _schemas_presentation_types.SlideSize
PresentationMetadata = _schemas_presentation_types.PresentationMetadata
CoreProperties = _schemas_presentation_types.CoreProperties
SlideMetadata = _schemas_presentation_types.SlideMetadata
Section = _schemas_presentation_types.Section
ShapeSearchQuery = _schemas_presentation_types.ShapeSearchQuery
ShapeSearchResult = _schemas_presentation_types.ShapeSearchResult
Author = _schemas_presentation_types.Author
Comment = _schemas_presentation_types.Comment
BatchCommand = _schemas_presentation_types.BatchCommand
BatchErrorDetail = _schemas_presentation_types.BatchErrorDetail
BatchItemResult = _schemas_presentation_types.BatchItemResult
TableCellInfo = _schemas_presentation_types.TableCellInfo
TableInfo = _schemas_presentation_types.TableInfo

TextFrame = _schemas_shape_types.TextFrame
Paragraph = _schemas_shape_types.Paragraph
FillFormat = _schemas_shape_types.FillFormat
GradientStop = _schemas_shape_types.GradientStop
GradientFill = _schemas_shape_types.GradientFill
PatternFill = _schemas_shape_types.PatternFill
PictureFill = _schemas_shape_style.PictureFill
PictureFillCrop = _schemas_shape_style.PictureFillCrop
FreeformPoint = _schemas_shape_style.FreeformPoint
FreeformSegment = _schemas_shape_style.FreeformSegment
FreeformPath = _schemas_shape_style.FreeformPath
FreeformGeometry = _schemas_shape_style.FreeformGeometry
EffectiveColor = _schemas_shape_style.EffectiveColor
EffectiveString = _schemas_shape_style.EffectiveString
EffectiveFloat = _schemas_shape_style.EffectiveFloat
EffectiveBool = _schemas_shape_style.EffectiveBool
EffectivePosition = _schemas_shape_style.EffectivePosition
EffectiveShapeStyle = _schemas_shape_style.EffectiveShapeStyle
LineFormat = _schemas_shape_types.LineFormat
ShadowFormat = _schemas_shape_types.ShadowFormat
GlowFormat = _schemas_shape_types.GlowFormat
BlurFormat = _schemas_shape_types.BlurFormat
SoftEdgeFormat = _schemas_shape_types.SoftEdgeFormat
ReflectionFormat = _schemas_shape_types.ReflectionFormat
ShapeProps = _schemas_shape_types.ShapeProps
ImageMetadata = _schemas_shape_types.ImageMetadata
SlideImageRef = _schemas_shape_types.SlideImageRef
SlideMediaRef = _schemas_shape_types.SlideMediaRef
ImageCrop = _schemas_shape_types.ImageCrop
Hyperlink = _schemas_shape_types.Hyperlink
TextRun = _schemas_shape_types.TextRun
ShapeTextParagraph = _schemas_shape_types.ShapeTextParagraph
ShapeAdjustment = _schemas_shape_types.ShapeAdjustment
ShapeAdjustmentValue = _schemas_shape_types.ShapeAdjustmentValue
ShapeUpdate = _schemas_shape_types.ShapeUpdate
Shape = _schemas_shape_types.Shape
GrayscaleShapeRef = _schemas_shape_types.GrayscaleShapeRef
GrayscaleTextRef = _schemas_shape_types.GrayscaleTextRef
GrayscalePlaceholderRef = _schemas_shape_types.GrayscalePlaceholderRef
GrayscaleScope = _schemas_shape_types.GrayscaleScope

ChartDataUpdate = _schemas_chart_layout.ChartDataUpdate
ChartDataSource = _schemas_chart_layout.ChartDataSource
ChartAxisState = _schemas_chart_layout.ChartAxisState
ChartState = _schemas_chart_layout.ChartState
ChartFormatUpdate = _schemas_chart_layout.ChartFormatUpdate
ChartSelector = _schemas_chart_layout.ChartSelector
ChartSeriesData = _schemas_chart_layout.ChartSeriesData
ChartTrendlineSpec = _schemas_chart_series.ChartTrendlineSpec
ChartErrorBarSpec = _schemas_chart_series.ChartErrorBarSpec
ChartDataPointSpec = _schemas_chart_series.ChartDataPointSpec
ChartDataLabelPointSpec = _schemas_chart_series.ChartDataLabelPointSpec
ChartSeriesInvertSpec = _schemas_chart_series.ChartSeriesInvertSpec
ChartDataTableSpec = _schemas_chart_series.ChartDataTableSpec
ChartLineFormatSpec = _schemas_chart_series.ChartLineFormatSpec
ChartSeriesFormatSpec = _schemas_chart_series.ChartSeriesFormatSpec
ChartSeriesLinesSpec = _schemas_chart_series.ChartSeriesLinesSpec
PlaceholderInfo = _schemas_chart_layout.PlaceholderInfo
SlideChartRef = _schemas_chart_layout.SlideChartRef
SlideLayoutInfo = _schemas_chart_layout.SlideLayoutInfo
SlideMasterCloneResult = _schemas_chart_layout.SlideMasterCloneResult

_HEX_RGB_DIGITS = 6
_CHANNEL_MAX = 255


class RGBColor(tuple[int, int, int]):
    """An (r, g, b) triple that formats as the uppercase hex OOXML expects."""

    __slots__ = ()

    def __new__(cls, r: int, g: int, b: int) -> Self:
        """Build a color from three 0-255 channel values."""
        channels = (int(r), int(g), int(b))
        for value in channels:
            if not 0 <= value <= _CHANNEL_MAX:
                msg = f"RGB channel out of range: {channels}"
                raise ValueError(msg)
        return super().__new__(cls, channels)

    @classmethod
    def from_string(cls, hex_str: str) -> Self:
        """Parse ``RRGGBB`` or ``#RRGGBB``, raising on anything else."""
        clean = hex_str.lstrip("#")
        if len(clean) != _HEX_RGB_DIGITS:
            msg = f"expected a 6-digit hex color, got {hex_str!r}"
            raise ValueError(msg)
        try:
            return cls(int(clean[0:2], 16), int(clean[2:4], 16), int(clean[4:6], 16))
        except ValueError as exc:
            msg = f"expected a 6-digit hex color, got {hex_str!r}"
            raise ValueError(msg) from exc

    @override
    def __str__(self) -> str:
        """Return the color as uppercase ``RRGGBB``."""
        return f"{self[0]:02X}{self[1]:02X}{self[2]:02X}"
