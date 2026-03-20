"""Theme definition and management for presentations."""

from .theme import (
    ColorScheme,
    FontScheme,
    Theme,
    create_aurora_theme,
    create_forest_theme,
    create_ocean_theme,
    create_sunset_theme,
    get_theme,
    list_themes,
)

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
