from .api import GopptxError, Presentation
from .ops import (
    OP_ADD_SLIDE,
    OP_DUPLICATE_SLIDE,
    OP_GET_METADATA,
    OP_MOVE_SLIDE,
    OP_REMOVE_SLIDE,
    OP_SLIDE_COUNT,
    SUPPORTED_OPS,
)
from .types import PresentationMetadata, SlideSize

__all__ = [
    "Presentation",
    "GopptxError",
    "PresentationMetadata",
    "SlideSize",
    "OP_SLIDE_COUNT",
    "OP_ADD_SLIDE",
    "OP_REMOVE_SLIDE",
    "OP_MOVE_SLIDE",
    "OP_DUPLICATE_SLIDE",
    "OP_GET_METADATA",
    "SUPPORTED_OPS",
]
