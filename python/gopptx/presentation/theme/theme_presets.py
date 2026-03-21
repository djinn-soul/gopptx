"""Built-in theme presets for gopptx."""

from __future__ import annotations

from .theme import ColorScheme, FontScheme, Theme


def create_aurora_theme() -> Theme:
    """Create the Aurora theme - cool blues and teals."""
    return Theme(
        name="Aurora",
        colors=ColorScheme(
            name="Aurora",
            dk1="0F172A",
            lt1="FFFFFF",
            dk2="1E293B",
            lt2="F7FAFC",
            accent1="2563EB",
            accent2="14B8A6",
            accent3="F59E0B",
            accent4="8B5CF6",
            accent5="EC4899",
            accent6="0EA5E9",
            hlink="2563EB",
        ),
        fonts=FontScheme(name="Aurora", major_font="Aptos Display", minor_font="Aptos"),
        metadata={"description": "Cool blues and teals", "author": "gopptx"},
    )


def create_ocean_theme() -> Theme:
    """Create the Ocean theme - deep blues and greens."""
    return Theme(
        name="Ocean",
        colors=ColorScheme(
            name="Ocean",
            dk1="0F172A",
            lt1="FFFFFF",
            dk2="164E63",
            lt2="E0F2FE",
            accent1="0369A1",
            accent2="059669",
            accent3="0891B2",
            accent4="0E7490",
            accent5="088397",
            accent6="06B6D4",
            hlink="0369A1",
        ),
        fonts=FontScheme(name="Ocean", major_font="Aptos Display", minor_font="Aptos"),
        metadata={"description": "Deep ocean colors", "author": "gopptx"},
    )


def create_sunset_theme() -> Theme:
    """Create the Sunset theme - warm oranges and reds."""
    return Theme(
        name="Sunset",
        colors=ColorScheme(
            name="Sunset",
            dk1="1F2937",
            lt1="FFFFFF",
            dk2="78350F",
            lt2="FEF3C7",
            accent1="F59E0B",
            accent2="EA580C",
            accent3="DC2626",
            accent4="CA8A04",
            accent5="B91C1C",
            accent6="F97316",
            hlink="EA580C",
        ),
        fonts=FontScheme(name="Sunset", major_font="Aptos Display", minor_font="Aptos"),
        metadata={"description": "Warm sunset colors", "author": "gopptx"},
    )


def create_forest_theme() -> Theme:
    """Create the Forest theme - greens and earth tones."""
    return Theme(
        name="Forest",
        colors=ColorScheme(
            name="Forest",
            dk1="14532D",
            lt1="FFFFFF",
            dk2="166534",
            lt2="DBEAFE",
            accent1="15803D",
            accent2="7C2D12",
            accent3="92400E",
            accent4="3F6212",
            accent5="84CC16",
            accent6="16A34A",
            hlink="15803D",
        ),
        fonts=FontScheme(name="Forest", major_font="Aptos Display", minor_font="Aptos"),
        metadata={"description": "Natural green palette", "author": "gopptx"},
    )


BUILT_IN_THEMES = {
    "aurora": create_aurora_theme,
    "ocean": create_ocean_theme,
    "sunset": create_sunset_theme,
    "forest": create_forest_theme,
}


def get_theme(name: str) -> Theme:
    """Get a built-in theme by name."""
    if name.lower() not in BUILT_IN_THEMES:
        available = ", ".join(BUILT_IN_THEMES.keys())
        raise ValueError(
            f"Unknown theme '{name}'. Available built-in themes: {available}"
        )
    return BUILT_IN_THEMES[name.lower()]()


def list_themes() -> list[str]:
    """List all available built-in theme names."""
    return list(BUILT_IN_THEMES.keys())
