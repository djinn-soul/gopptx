"""Theme application mixin for presentations."""

from __future__ import annotations

from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from .theme.theme import Theme


class PresentationThemeMixin:
    """Mixin providing theme application methods for Presentation."""

    def apply_theme(self, theme: Theme) -> None:
        """Apply a theme to the presentation.

        Applies both the color scheme and font scheme from the theme to all
        slides in the presentation.

        Args:
            theme: Theme instance with color and font schemes.

        Examples:
            from gopptx import Presentation
            from gopptx.presentation.theme import get_theme

            with Presentation.new("My Presentation") as prs:
                # Apply a built-in theme
                aurora_theme = get_theme("aurora")
                prs.apply_theme(aurora_theme)

                # Add slides that will use the theme colors and fonts
                prs.add_slide("Title Slide")
                prs.add_slide("Content Slide")

                prs.save("output.pptx")

            # Or create a custom theme
            from gopptx.presentation.theme import Theme, ColorScheme, FontScheme

            custom_theme = Theme(
                name="Custom",
                colors=ColorScheme(
                    name="Custom",
                    dk1="1a1a1a",
                    lt1="ffffff",
                    accent1="ff6b6b",
                    accent2="4ecdc4",
                ),
                fonts=FontScheme(
                    name="Custom",
                    major_font="Calibri Light",
                    minor_font="Calibri",
                ),
            )
            prs.apply_theme(custom_theme)
        """
        # Apply color scheme
        color_dict = theme.colors.to_dict()
        self.set_theme_color_scheme(**color_dict)

        # Apply font scheme
        font_dict = theme.fonts.to_dict()
        self.set_theme_font_scheme(
            font_dict["major_font"],
            font_dict["minor_font"],
        )
