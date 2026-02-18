from __future__ import annotations

from .api_presentation_base import PresentationBase
from .api_presentation_content import PresentationContentMixin
from .api_presentation_slides import PresentationSlidesMixin


class Presentation(PresentationContentMixin, PresentationSlidesMixin, PresentationBase):
    """High-level wrapper for a PowerPoint presentation handled by the Go engine."""

