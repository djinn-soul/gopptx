"""Public API exports for gopptx library."""

from __future__ import annotations

from .api_errors import GopptxError
from .builder import PresentationBuilder
from .presentation.presentation import Presentation
from .slide.slide import Slide

__all__ = ["GopptxError", "Presentation", "PresentationBuilder", "Slide"]
