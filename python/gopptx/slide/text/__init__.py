"""Text-domain package for slide object-model helpers."""

from .text_frame import TextFrameProps, serialize_text_frame_for_payload
from .text_model import (
    ShapeParagraphCollection,
    ShapeParagraphProxy,
    ShapeRunCollection,
    ShapeRunProxy,
    ShapeTextFrame,
)
from .text_paragraph import ParagraphProps, serialize_paragraph_for_payload
from .text_run import Run, RunHyperlink, serialize_runs_for_payload

__all__ = [
    "ParagraphProps",
    "Run",
    "RunHyperlink",
    "ShapeParagraphCollection",
    "ShapeParagraphProxy",
    "ShapeRunCollection",
    "ShapeRunProxy",
    "ShapeTextFrame",
    "TextFrameProps",
    "serialize_paragraph_for_payload",
    "serialize_runs_for_payload",
    "serialize_text_frame_for_payload",
]
