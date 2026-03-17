"""Type definitions for gopptx library."""

from __future__ import annotations

from . import schemas_chart_layout as _schemas_chart_layout
from . import schemas_presentation_types as _schemas_presentation_types
from . import schemas_shape_types as _schemas_shape_types

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
LineFormat = _schemas_shape_types.LineFormat
ShadowFormat = _schemas_shape_types.ShadowFormat
GlowFormat = _schemas_shape_types.GlowFormat
BlurFormat = _schemas_shape_types.BlurFormat
SoftEdgeFormat = _schemas_shape_types.SoftEdgeFormat
ReflectionFormat = _schemas_shape_types.ReflectionFormat
ShapeProps = _schemas_shape_types.ShapeProps
ImageMetadata = _schemas_shape_types.ImageMetadata
SlideImageRef = _schemas_shape_types.SlideImageRef
ImageCrop = _schemas_shape_types.ImageCrop
Hyperlink = _schemas_shape_types.Hyperlink
TextRun = _schemas_shape_types.TextRun
ShapeUpdate = _schemas_shape_types.ShapeUpdate
Shape = _schemas_shape_types.Shape

ChartDataUpdate = _schemas_chart_layout.ChartDataUpdate
ChartAxisState = _schemas_chart_layout.ChartAxisState
ChartState = _schemas_chart_layout.ChartState
ChartFormatUpdate = _schemas_chart_layout.ChartFormatUpdate
ChartSelector = _schemas_chart_layout.ChartSelector
ChartSeriesData = _schemas_chart_layout.ChartSeriesData
PlaceholderInfo = _schemas_chart_layout.PlaceholderInfo
SlideChartRef = _schemas_chart_layout.SlideChartRef
SlideLayoutInfo = _schemas_chart_layout.SlideLayoutInfo
SlideMasterCloneResult = _schemas_chart_layout.SlideMasterCloneResult

RGBColor = str  # Hex string like 'FF0000'
