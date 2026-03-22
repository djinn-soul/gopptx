"""Presentation-level header/footer control mixin."""

from __future__ import annotations

from .. import ops
from .helpers import PresentationMixinBase


class PresentationHeaderFooterMixin(PresentationMixinBase):
    """Mixin providing presentation-level header/footer control."""

    def __init__(self, *args: object, **kwargs: object) -> None:
        """Initialize header/footer defaults."""
        super().__init__(*args, **kwargs)
        self._header_footer_defaults: dict[str, object] = {
            "footer": "",
            "show_footer": False,
            "show_slide_num": False,
            "show_date_time": False,
            "date_time_text": "",
        }

    def set_header_footer(
        self,
        footer: str = "",
        *,
        show_footer: bool = False,
        show_slide_num: bool = False,
        show_date_time: bool = False,
        date_time_text: str = "",
    ) -> None:
        """Set header/footer for ALL slides in the presentation.

        This is a presentation-wide setting that applies to existing slides
        and becomes the default for all subsequently added slides.
        Individual slides can override using slide.set_header_footer().

        Args:
            footer: Footer text to display.
            show_footer: Whether to show the footer.
            show_slide_num: Whether to show the slide number.
            show_date_time: Whether to show the date/time.
            date_time_text: Fixed date/time string (empty = auto).

        Example:
            # Apply footer to ALL slides (existing and future)
            prs.set_header_footer(
                footer="Confidential",
                show_footer=True,
                show_slide_num=True,
                show_date_time=True,
            )

            # Now add slides - they all have the footer
            prs.add_slide("Slide 1")  # Has footer
            prs.add_slide("Slide 2")  # Has footer

            # Override on specific slide
            prs.slides[0].set_header_footer(
                footer="Different Footer",
                show_footer=True,
                show_slide_num=False,  # Turn off slide numbers on this one
            )
        """
        # Store as defaults for future slides
        self._header_footer_defaults = {
            "footer": footer,
            "show_footer": show_footer,
            "show_slide_num": show_slide_num,
            "show_date_time": show_date_time,
            "date_time_text": date_time_text,
        }

        # Apply to all existing slides
        for slide_index in range(len(self.slides)):
            payload: dict[str, object] = {
                "slide_index": slide_index,
                "footer": footer,
                "show_footer": show_footer,
                "show_slide_num": show_slide_num,
                "show_date_time": show_date_time,
                "date_time_text": date_time_text,
            }
            self.execute(ops.OP_SET_SLIDE_HEADER_FOOTER, payload)

    def get_header_footer(self, _slide_index: int) -> dict[str, object]:
        """Get header/footer configuration for a specific slide.

        Args:
            slide_index: Zero-based slide index.

        Returns:
            Dict with footer, show_footer, show_slide_num, show_date_time keys.
        """
        return dict(self._header_footer_defaults)


__all__ = ["PresentationHeaderFooterMixin"]
