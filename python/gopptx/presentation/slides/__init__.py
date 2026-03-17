"""Slide-domain mixins for presentation APIs."""

from .layout_theme_mixin import PresentationLayoutMixin, PresentationThemeMixin
from .master import SlideLayout, SlideLayouts, SlideMaster, SlideMasters
from .properties_mixin import PresentationPropertiesMixin
from .sections_mixin import PresentationSectionMixin
from .slide_lookup_mixin import PresentationSlideLookupMixin
from .slide_proxy_mixin import PresentationSlideProxyMixin
from .slides_mixin import PresentationSlidesMixin

__all__ = [
    "PresentationLayoutMixin",
    "PresentationPropertiesMixin",
    "PresentationSectionMixin",
    "PresentationSlideLookupMixin",
    "PresentationSlideProxyMixin",
    "PresentationSlidesMixin",
    "PresentationThemeMixin",
    "SlideLayout",
    "SlideLayouts",
    "SlideMaster",
    "SlideMasters",
]
