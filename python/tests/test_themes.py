"""Tests for theme system (ColorScheme, FontScheme, Theme classes)."""

import pytest

from gopptx.presentation.theme import (
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


class TestColorScheme:
    """Tests for ColorScheme class."""

    def test_default_colors(self):
        """ColorScheme has sensible defaults."""
        scheme = ColorScheme(name="Test")
        assert scheme.name == "Test"
        assert scheme.dk1 == "000000"  # Black
        assert scheme.lt1 == "FFFFFF"  # White
        assert scheme.dk2 == "1F2937"
        assert scheme.lt2 == "F9FAFB"
        assert scheme.accent1 == "2563EB"
        assert scheme.accent2 == "059669"

    def test_custom_colors(self):
        """ColorScheme allows custom color values."""
        scheme = ColorScheme(
            name="Custom",
            dk1="111111",
            accent1="FF0000",
            accent2="00FF00",
        )
        assert scheme.dk1 == "111111"
        assert scheme.accent1 == "FF0000"
        assert scheme.accent2 == "00FF00"

    def test_to_dict_includes_all_required(self):
        """to_dict() includes all core colors."""
        scheme = ColorScheme(name="Test")
        result = scheme.to_dict()
        assert result["dk1"] == "000000"
        assert result["lt1"] == "FFFFFF"
        assert result["dk2"] == "1F2937"
        assert result["lt2"] == "F9FAFB"
        assert result["accent1"] == "2563EB"
        assert result["accent2"] == "059669"
        assert result["accent3"] == "DC2626"
        assert result["accent4"] == "7C3AED"
        assert result["accent5"] == "DB2777"
        assert result["accent6"] == "0891B2"

    def test_to_dict_includes_optional_hyperlinks(self):
        """to_dict() includes hyperlink colors when provided."""
        scheme = ColorScheme(
            name="Test",
            hlink="0066CC",
            fol_hlink="CC0066",
        )
        result = scheme.to_dict()
        assert result["hlink"] == "0066CC"
        assert result["fol_hlink"] == "CC0066"

    def test_to_dict_excludes_none_hyperlinks(self):
        """to_dict() excludes None hyperlink colors."""
        scheme = ColorScheme(name="Test")
        result = scheme.to_dict()
        assert "hlink" not in result
        assert "fol_hlink" not in result


class TestFontScheme:
    """Tests for FontScheme class."""

    def test_default_fonts(self):
        """FontScheme has sensible defaults."""
        scheme = FontScheme(name="Test")
        assert scheme.name == "Test"
        assert scheme.major_font == "Aptos Display"
        assert scheme.minor_font == "Aptos"

    def test_custom_fonts(self):
        """FontScheme allows custom font values."""
        scheme = FontScheme(
            name="Custom",
            major_font="Calibri Light",
            minor_font="Calibri",
        )
        assert scheme.major_font == "Calibri Light"
        assert scheme.minor_font == "Calibri"

    def test_to_dict(self):
        """to_dict() returns font names."""
        scheme = FontScheme(
            name="Test",
            major_font="Arial",
            minor_font="Verdana",
        )
        result = scheme.to_dict()
        assert result["major_font"] == "Arial"
        assert result["minor_font"] == "Verdana"


class TestTheme:
    """Tests for Theme class."""

    def test_create_theme(self):
        """Theme combines color and font schemes."""
        colors = ColorScheme(name="TestColors")
        fonts = FontScheme(name="TestFonts")
        theme = Theme(name="TestTheme", colors=colors, fonts=fonts)

        assert theme.name == "TestTheme"
        assert theme.colors is colors
        assert theme.fonts is fonts
        assert theme.metadata == {}

    def test_theme_with_metadata(self):
        """Theme can store optional metadata."""
        colors = ColorScheme(name="Colors")
        fonts = FontScheme(name="Fonts")
        metadata = {"description": "Custom theme", "author": "TestAuthor"}
        theme = Theme(
            name="Test",
            colors=colors,
            fonts=fonts,
            metadata=metadata,
        )

        assert theme.metadata == metadata

    def test_to_dict(self):
        """to_dict() returns complete theme configuration."""
        colors = ColorScheme(
            name="Colors",
            accent1="FF0000",
        )
        fonts = FontScheme(
            name="Fonts",
            major_font="Arial",
        )
        theme = Theme(
            name="Test",
            colors=colors,
            fonts=fonts,
            metadata={"author": "Test"},
        )

        result = theme.to_dict()
        assert result["name"] == "Test"
        assert result["metadata"]["author"] == "Test"
        assert isinstance(result["colors"], dict)
        assert isinstance(result["fonts"], dict)
        assert result["colors"]["accent1"] == "FF0000"
        assert result["fonts"]["major_font"] == "Arial"


class TestBuiltInThemes:
    """Tests for built-in theme factories."""

    def test_aurora_theme(self):
        """Aurora theme has correct colors and fonts."""
        theme = create_aurora_theme()
        assert theme.name == "Aurora"
        assert theme.colors.name == "Aurora"
        assert theme.fonts.name == "Aurora"
        # Aurora colors are cool blues and teals
        assert theme.colors.accent1 == "2563EB"  # Blue
        assert theme.colors.accent2 == "14B8A6"  # Teal
        assert theme.fonts.major_font == "Aptos Display"
        assert "description" in theme.metadata

    def test_ocean_theme(self):
        """Ocean theme has deep ocean colors."""
        theme = create_ocean_theme()
        assert theme.name == "Ocean"
        assert theme.colors.name == "Ocean"
        # Ocean colors are deep blues and greens
        assert theme.colors.accent1 == "0369A1"  # Sky blue
        assert theme.colors.accent2 == "059669"  # Emerald
        assert "Deep ocean" in theme.metadata.get("description", "")

    def test_sunset_theme(self):
        """Sunset theme has warm colors."""
        theme = create_sunset_theme()
        assert theme.name == "Sunset"
        # Sunset colors are warm oranges and reds
        assert theme.colors.accent1 == "F59E0B"  # Amber
        assert theme.colors.accent2 == "EA580C"  # Orange
        assert theme.colors.accent3 == "DC2626"  # Red
        assert "Warm sunset" in theme.metadata.get("description", "")

    def test_forest_theme(self):
        """Forest theme has natural green palette."""
        theme = create_forest_theme()
        assert theme.name == "Forest"
        # Forest colors are greens and earth tones
        assert theme.colors.accent1 == "15803D"  # Green
        assert theme.colors.accent2 == "7C2D12"  # Brown
        assert "Natural green" in theme.metadata.get("description", "")


class TestThemeRegistry:
    """Tests for theme registry functions."""

    def test_list_themes(self):
        """list_themes() returns available theme names."""
        themes = list_themes()
        assert isinstance(themes, list)
        assert "aurora" in themes
        assert "ocean" in themes
        assert "sunset" in themes
        assert "forest" in themes

    def test_get_theme_aurora(self):
        """get_theme('aurora') returns Aurora theme."""
        theme = get_theme("aurora")
        assert theme.name == "Aurora"
        assert theme.colors.accent1 == "2563EB"

    def test_get_theme_ocean(self):
        """get_theme('ocean') returns Ocean theme."""
        theme = get_theme("ocean")
        assert theme.name == "Ocean"

    def test_get_theme_sunset(self):
        """get_theme('sunset') returns Sunset theme."""
        theme = get_theme("sunset")
        assert theme.name == "Sunset"

    def test_get_theme_forest(self):
        """get_theme('forest') returns Forest theme."""
        theme = get_theme("forest")
        assert theme.name == "Forest"

    def test_get_theme_case_insensitive(self):
        """get_theme() is case-insensitive."""
        theme_lower = get_theme("aurora")
        theme_upper = get_theme("AURORA")
        theme_mixed = get_theme("AuRoRa")

        assert theme_lower.name == theme_upper.name == theme_mixed.name

    def test_get_theme_invalid_name(self):
        """get_theme() raises ValueError for unknown theme."""
        with pytest.raises(ValueError) as exc_info:
            get_theme("nonexistent")
        assert "Unknown theme" in str(exc_info.value)
        assert "aurora" in str(exc_info.value)


class TestThemeImmutability:
    """Tests for theme immutability patterns."""

    def test_builtin_themes_independent(self):
        """Each call to builtin factory creates independent instance."""
        aurora1 = create_aurora_theme()
        aurora2 = create_aurora_theme()

        # Same values
        assert aurora1.name == aurora2.name
        assert aurora1.colors.accent1 == aurora2.colors.accent1

        # But different instances
        assert aurora1 is not aurora2
        assert aurora1.colors is not aurora2.colors
        assert aurora1.fonts is not aurora2.fonts

    def test_custom_theme_modification(self):
        """Custom themes can be modified."""
        theme = Theme(
            name="Test",
            colors=ColorScheme(name="Colors", accent1="FF0000"),
            fonts=FontScheme(name="Fonts"),
        )

        # Original value
        assert theme.colors.accent1 == "FF0000"

        # Modify
        theme.colors.accent1 = "00FF00"
        assert theme.colors.accent1 == "00FF00"


class TestThemeColorPalettes:
    """Tests for color palette completeness."""

    def test_all_themes_have_required_colors(self):
        """All built-in themes define all required colors."""
        required_colors = [
            "dk1",
            "lt1",
            "dk2",
            "lt2",
            "accent1",
            "accent2",
            "accent3",
            "accent4",
            "accent5",
            "accent6",
        ]

        for theme_factory in [
            create_aurora_theme,
            create_ocean_theme,
            create_sunset_theme,
            create_forest_theme,
        ]:
            theme = theme_factory()
            color_dict = theme.colors.to_dict()

            for color_name in required_colors:
                assert color_name in color_dict, (
                    f"{theme.name} missing color: {color_name}"
                )
                assert (
                    color_dict[color_name]
                ), f"{theme.name} {color_name} is empty"

    def test_all_themes_have_fonts(self):
        """All built-in themes define fonts."""
        for theme_factory in [
            create_aurora_theme,
            create_ocean_theme,
            create_sunset_theme,
            create_forest_theme,
        ]:
            theme = theme_factory()
            assert theme.fonts.major_font, f"{theme.name} missing major_font"
            assert theme.fonts.minor_font, f"{theme.name} missing minor_font"

    def test_color_values_are_hex(self):
        """All color values are valid hex strings."""
        for theme_factory in [
            create_aurora_theme,
            create_ocean_theme,
            create_sunset_theme,
            create_forest_theme,
        ]:
            theme = theme_factory()
            color_dict = theme.colors.to_dict()

            for color_name, color_value in color_dict.items():
                # Should be 6-char hex (no # prefix)
                assert isinstance(color_value, str), (
                    f"{theme.name} {color_name} is not string"
                )
                assert len(color_value) == 6, (
                    f"{theme.name} {color_name} not 6-char hex: {color_value}"
                )
                assert all(
                    c in "0123456789ABCDEFabcdef" for c in color_value
                ), (
                    f"{theme.name} {color_name} not valid hex: {color_value}"
                )


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
