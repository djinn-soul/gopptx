"""Slide-domain mixins for presentation APIs."""

from .properties_mixin import PresentationPropertiesMixin
from .sections_mixin import PresentationSectionMixin
from .slides_mixin import PresentationSlidesMixin

__all__ = [
    "PresentationPropertiesMixin",
    "PresentationSectionMixin",
    "PresentationSlidesMixin",
]
