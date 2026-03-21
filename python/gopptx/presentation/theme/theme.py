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


# ============================================================================
# Built-in Themes
# ============================================================================


def create_aurora_theme() -> Theme:
    """Create the Aurora theme - cool blues and teals.

    Features:
    - Dark navy text on light backgrounds
    - Blue and teal accent colors
    - Modern, professional feel

    Returns:
        Aurora theme instance.
    """
    return Theme(
        name="Aurora",
        colors=ColorScheme(
            name="Aurora",
            dk1="0F172A",  # Dark navy
            lt1="FFFFFF",  # White
            dk2="1E293B",  # Slate
            lt2="F7FAFC",  # Light gray
            accent1="2563EB",  # Blue
            accent2="14B8A6",  # Teal
            accent3="F59E0B",  # Amber
            accent4="8B5CF6",  # Purple
            accent5="EC4899",  # Pink
            accent6="0EA5E9",  # Cyan
            hlink="2563EB",
        ),
        fonts=FontScheme(
            name="Aurora",
            major_font="Aptos Display",
            minor_font="Aptos",
        ),
        metadata={"description": "Cool blues and teals", "author": "gopptx"},
    )


def create_ocean_theme() -> Theme:
    """Create the Ocean theme - deep blues and greens.

    Features:
    - Deep ocean blues
    - Teal and emerald accents
    - Professional, calm aesthetic

    Returns:
        Ocean theme instance.
    """
    return Theme(
        name="Ocean",
        colors=ColorScheme(
            name="Ocean",
            dk1="0F172A",  # Deep navy
            lt1="FFFFFF",  # White
            dk2="164E63",  # Ocean blue
            lt2="E0F2FE",  # Light sky
            accent1="0369A1",  # Sky blue
            accent2="059669",  # Emerald
            accent3="0891B2",  # Cyan
            accent4="0E7490",  # Teal
            accent5="088397",  # Darker teal
            accent6="06B6D4",  # Cyan
            hlink="0369A1",
        ),
        fonts=FontScheme(
            name="Ocean",
            major_font="Aptos Display",
            minor_font="Aptos",
        ),
        metadata={"description": "Deep ocean colors", "author": "gopptx"},
    )


def create_sunset_theme() -> Theme:
    """Create the Sunset theme - warm oranges and reds.

    Features:
    - Warm sunset colors
    - Orange, red, and gold accents
    - Energetic, vibrant feel

    Returns:
        Sunset theme instance.
    """
    return Theme(
        name="Sunset",
        colors=ColorScheme(
            name="Sunset",
            dk1="1F2937",  # Dark gray
            lt1="FFFFFF",  # White
            dk2="78350F",  # Dark brown
            lt2="FEF3C7",  # Pale yellow
            accent1="F59E0B",  # Amber
            accent2="EA580C",  # Orange
            accent3="DC2626",  # Red
            accent4="CA8A04",  # Yellow
            accent5="B91C1C",  # Dark red
            accent6="F97316",  # Bright orange
            hlink="EA580C",
        ),
        fonts=FontScheme(
            name="Sunset",
            major_font="Aptos Display",
            minor_font="Aptos",
        ),
        metadata={"description": "Warm sunset colors", "author": "gopptx"},
    )


def create_forest_theme() -> Theme:
    """Create the Forest theme - greens and earth tones.

    Features:
    - Natural green palette
    - Earth tones and sage accents
    - Calm, organic aesthetic

    Returns:
        Forest theme instance.
    """
    return Theme(
        name="Forest",
        colors=ColorScheme(
            name="Forest",
            dk1="14532D",  # Dark forest
            lt1="FFFFFF",  # White
            dk2="166534",  # Forest green
            lt2="DBEAFE",  # Light mint
            accent1="15803D",  # Green
            accent2="7C2D12",  # Brown
            accent3="92400E",  # Tan
            accent4="3F6212",  # Sage
            accent5="84CC16",  # Lime
            accent6="16A34A",  # Bright green
            hlink="15803D",
        ),
        fonts=FontScheme(
            name="Forest",
            major_font="Aptos Display",
            minor_font="Aptos",
        ),
        metadata={"description": "Natural green palette", "author": "gopptx"},
    )


# Theme registry
BUILT_IN_THEMES = {
    "aurora": create_aurora_theme,
    "ocean": create_ocean_theme,
    "sunset": create_sunset_theme,
    "forest": create_forest_theme,
}


def get_theme(name: str) -> Theme:
    """Get a built-in theme by name.

    Args:
        name: Theme name ('aurora', 'ocean', 'sunset', 'forest').

    Returns:
        Theme instance.

    Raises:
        ValueError: If theme name is not found.

    Examples:
        theme = get_theme("aurora")
        prs.apply_theme(theme)
    """
    if name.lower() not in BUILT_IN_THEMES:
        available = ", ".join(BUILT_IN_THEMES.keys())
        raise ValueError(
            f"Unknown theme '{name}'. "
            f"Available built-in themes: {available}"
        )
    return BUILT_IN_THEMES[name.lower()]()


def list_themes() -> list[str]:
    """List all available built-in theme names.

    Returns:
        List of theme names.

    Examples:
        themes = list_themes()  # ['aurora', 'ocean', 'sunset', 'forest']
    """
    return list(BUILT_IN_THEMES.keys())


__all__ = [
    "ColorScheme",
    "FontScheme",
    "Theme",
    "create_aurora_theme",
    "create_ocean_theme",
    "create_sunset_theme",
    "create_forest_theme",
    "get_theme",
    "list_themes",
]
