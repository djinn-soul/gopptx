"""Public API exports for gopptx library."""

from __future__ import annotations

from .api_errors import GopptxError
from .presentation.presentation import Presentation
from .slide.slide import Slide

__all__ = ["GopptxError", "Presentation", "Slide"]
