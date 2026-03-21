"""ShapeType enum compatibility checks for add_shape APIs."""

from __future__ import annotations

from gopptx import Presentation
from gopptx.constants import ConnectorType, ShapeType
from gopptx.schemas import Inches


def test_slide_add_shape_accepts_shape_type_enum() -> None:
    """Slide.add_shape should accept ShapeType enum values."""
    with Presentation.new("ShapeType enum") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_shape(
            ShapeType.RECTANGLE,
            (Inches(1.0), Inches(1.0), Inches(2.0), Inches(1.0)),
            text="Enum shape",
        )
        assert shape_id > 0


def test_slide_add_connector_accepts_connector_type_enum() -> None:
    """Slide.add_connector should accept ConnectorType enum values."""
    with Presentation.new("ConnectorType enum") as pres:
        slide = pres.slides[0]
        shape_id = slide.add_connector(
            ConnectorType.STRAIGHT,
            Inches(1.0),
            Inches(1.0),
            Inches(3.0),
            Inches(1.0),
        )
        assert shape_id > 0
