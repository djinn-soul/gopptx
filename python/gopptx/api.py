"""Public API exports for gopptx library."""

from __future__ import annotations

from .api_errors import GopptxError
from .api_presentation import Presentation
from .api_slide import Slide

__all__ = ["GopptxError", "Presentation", "Slide"]
