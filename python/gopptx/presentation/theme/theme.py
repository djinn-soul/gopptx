"""Theme definition and management for presentations."""

from __future__ import annotations

from dataclasses import dataclass, field


@dataclass
class ColorScheme:
    """Color scheme for a presentation theme.

    Defines all colors used throughout the theme, including dark/light variants
    and accent colors.
    """

    name: str
    """Name of the color scheme (e.g., 'Aurora')."""

    # Core colors
    dk1: str = "000000"
    """Dark color 1 (typically dark gray or black for text)."""

    lt1: str = "FFFFFF"
    """Light color 1 (typically white or off-white)."""

    dk2: str = "1F2937"
    """Dark color 2 (darker shade for secondary elements)."""

    lt2: str = "F9FAFB"
    """Light color 2 (light shade for backgrounds)."""

    # Accent colors (up to 6)
    accent1: str = "2563EB"
    """Primary accent color."""

    accent2: str = "059669"
    """Secondary accent color."""

    accent3: str = "DC2626"
    """Tertiary accent color."""

    accent4: str = "7C3AED"
    """Quaternary accent color."""

    accent5: str = "DB2777"
    """Quinary accent color."""

    accent6: str = "0891B2"
    """Senary accent color."""

    # Special colors
    hlink: str | None = None
    """Hyperlink color (defaults to accent1)."""

    fol_hlink: str | None = None
    """Followed hyperlink color (defaults to accent4)."""

    def to_dict(self) -> dict[str, str]:
        """Convert to dictionary for bridge operations.

        Returns:
            Dictionary with all color values.
        """
        result = {
            "dk1": self.dk1,
            "lt1": self.lt1,
            "dk2": self.dk2,
            "lt2": self.lt2,
            "accent1": self.accent1,
            "accent2": self.accent2,
            "accent3": self.accent3,
            "accent4": self.accent4,
            "accent5": self.accent5,
            "accent6": self.accent6,
        }
        if self.hlink:
            result["hlink"] = self.hlink
        if self.fol_hlink:
            result["fol_hlink"] = self.fol_hlink
        return result


@dataclass
class FontScheme:
    """Font scheme for a presentation theme.

    Defines the major (heading) and minor (body) fonts used throughout.
    """

    name: str
    """Name of the font scheme (e.g., 'Aurora')."""

    major_font: str = "Aptos Display"
    """Font for major/heading text."""

    minor_font: str = "Aptos"
    """Font for minor/body text."""

    def to_dict(self) -> dict[str, str]:
        """Convert to dictionary for bridge operations.

        Returns:
            Dictionary with font names.
        """
        return {
            "major_font": self.major_font,
            "minor_font": self.minor_font,
        }


@dataclass
class Theme:
    """Complete presentation theme definition.

    A theme bundles color and font schemes together, providing a cohesive
    visual identity for presentations. Apply a theme to instantly set
    colors, fonts, and typography across all slides.

    Examples:
        # Define a custom theme
        aurora_theme = Theme(
            name="Aurora",
            colors=ColorScheme(
                name="Aurora",
                dk1="0F172A",
                lt1="FFFFFF",
                accent1="2563EB",
                accent2="14B8A6",
                accent3="F59E0B",
            ),
            fonts=FontScheme(
                name="Aurora",
                major_font="Aptos Display",
                minor_font="Aptos",
            ),
        )

        # Apply theme to presentation
        with Presentation.new("My Deck") as prs:
            prs.apply_theme(aurora_theme)
    """

    name: str
    """Name of the theme (e.g., 'Aurora', 'Ocean', 'Sunset')."""

    colors: ColorScheme
    """Color scheme for this theme."""

    fonts: FontScheme
    """Font scheme for this theme."""

    metadata: dict[str, str] = field(default_factory=dict)
    """Optional metadata (description, author, etc.)."""

    def to_dict(self) -> dict[str, object]:
        """Convert to dictionary for bridge operations.

        Returns:
            Dictionary with theme configuration.
        """
        return {
            "name": self.name,
            "colors": self.colors.to_dict(),
            "fonts": self.fonts.to_dict(),
            "metadata": self.metadata,
        }


__all__ = [
    "ColorScheme",
    "FontScheme",
    "Theme",
]
