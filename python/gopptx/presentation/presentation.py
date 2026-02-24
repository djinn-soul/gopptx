"""Main Presentation class for gopptx library."""

from __future__ import annotations

from .base import PresentationBase
from .content import PresentationContentMixin
from .slides import PresentationSlidesMixin


class Presentation(PresentationContentMixin, PresentationSlidesMixin, PresentationBase):
    """High-level wrapper for a PowerPoint presentation handled by the Go engine."""
