"""Integration tests for Presentation.apply_theme() method."""

import pytest
from gopptx import Presentation
from gopptx.presentation.theme import (
    ColorScheme,
    FontScheme,
    Theme,
    create_aurora_theme,
    create_forest_theme,
    create_ocean_theme,
    create_sunset_theme,
    get_theme,
)


class TestPresentationApplyTheme:
    """Tests for Presentation.apply_theme() integration."""

    def test_apply_builtin_theme_aurora(self):
        """Presentation can apply Aurora theme."""
        with Presentation.new("Test") as prs:
            theme = create_aurora_theme()
            prs.apply_theme(theme)
            # Should not raise

    def test_apply_builtin_theme_ocean(self):
        """Presentation can apply Ocean theme."""
        with Presentation.new("Test") as prs:
            theme = create_ocean_theme()
            prs.apply_theme(theme)
            # Should not raise

    def test_apply_builtin_theme_sunset(self):
        """Presentation can apply Sunset theme."""
        with Presentation.new("Test") as prs:
            theme = create_sunset_theme()
            prs.apply_theme(theme)
            # Should not raise

    def test_apply_builtin_theme_forest(self):
        """Presentation can apply Forest theme."""
        with Presentation.new("Test") as prs:
            theme = create_forest_theme()
            prs.apply_theme(theme)
            # Should not raise

    def test_apply_custom_theme(self):
        """Presentation can apply custom theme."""
        custom_theme = Theme(
            name="Custom",
            colors=ColorScheme(
                name="Custom",
                accent1="FF0000",
                accent2="00FF00",
            ),
            fonts=FontScheme(
                name="Custom",
                major_font="Arial",
                minor_font="Verdana",
            ),
        )

        with Presentation.new("Test") as prs:
            prs.apply_theme(custom_theme)
            # Should not raise

    def test_apply_theme_via_get_theme(self):
        """Presentation can apply theme via get_theme()."""
        with Presentation.new("Test") as prs:
            theme = get_theme("aurora")
            prs.apply_theme(theme)
            # Should not raise

    def test_apply_theme_with_slides(self):
        """Theme applies correctly when slides already exist."""
        with Presentation.new("Test") as prs:
            # Add slides before applying theme
            prs.add_slide("Slide 1")
            prs.add_slide("Slide 2")

            # Then apply theme
            theme = create_aurora_theme()
            prs.apply_theme(theme)

            # Verify presentation is still intact
            assert len(prs.slides) == 3  # Blank + 2 added

    def test_apply_theme_multiple_times(self):
        """Can apply different themes sequentially."""
        with Presentation.new("Test") as prs:
            aurora = create_aurora_theme()
            ocean = create_ocean_theme()

            # Apply first theme
            prs.apply_theme(aurora)

            # Then apply another
            prs.apply_theme(ocean)

            # Should have ocean colors now
            # (No direct access to theme colors, but should not raise)

    def test_apply_theme_and_save(self, tmp_path):
        """Themed presentation can be saved and reopened."""
        pptx_file = tmp_path / "themed.pptx"

        # Create with theme
        with Presentation.new("Themed") as prs:
            theme = create_aurora_theme()
            prs.apply_theme(theme)
            prs.add_slide("Content Slide")
            prs.save(str(pptx_file))

        # Verify file exists
        assert pptx_file.exists()

        # Reopen and verify
        prs2 = Presentation()
        prs2.open(str(pptx_file))
        try:
            assert len(prs2.slides) == 2
        finally:
            prs2.close()

    def test_apply_theme_before_content(self):
        """Can apply theme before adding content."""
        with Presentation.new("Test") as prs:
            theme = create_sunset_theme()
            prs.apply_theme(theme)

            # Add content after theme
            prs.add_slide("First Slide", layout="title_only")
            prs.add_slide("Second Slide", layout="title_and_content")

            assert len(prs.slides) == 3

    def test_apply_theme_after_content(self):
        """Can apply theme after adding content."""
        with Presentation.new("Test") as prs:
            # Add content first
            prs.add_slide("First Slide")
            prs.add_slide("Second Slide")

            # Then apply theme
            theme = create_forest_theme()
            prs.apply_theme(theme)

            assert len(prs.slides) == 3

    def test_theme_with_complete_metadata(self):
        """Theme with metadata applies correctly."""
        custom_theme = Theme(
            name="Branded",
            colors=ColorScheme(
                name="BrandColors",
                dk1="1a1a1a",
                accent1="FF6B6B",
            ),
            fonts=FontScheme(
                name="BrandFonts",
                major_font="Calibri Light",
                minor_font="Calibri",
            ),
            metadata={
                "description": "Company branded theme",
                "version": "1.0",
                "author": "Design Team",
            },
        )

        with Presentation.new("Branded") as prs:
            prs.apply_theme(custom_theme)

            # Metadata is preserved in theme object
            assert custom_theme.metadata["description"] == "Company branded theme"


class TestThemeApplicationFlow:
    """Tests for complete theme application workflows."""

    def test_create_presentation_with_theme(self, tmp_path):
        """Complete flow: create presentation with theme and save."""
        pptx_file = tmp_path / "output.pptx"

        with Presentation.new("Theme Demo") as prs:
            # Apply theme
            theme = get_theme("aurora")
            prs.apply_theme(theme)

            # Add themed content
            prs.slides[0].title = "Aurora Theme"
            prs.slides[0].body = "Cool blues and teals"

            prs.add_slide(
                "Benefits",
                layout="title_and_content",
                bullets=[
                    "Consistent visual identity",
                    "Professional appearance",
                    "Easy theme switching",
                ],
            )

            prs.add_slide("Colors", layout="title_only")

            # Save
            prs.save(str(pptx_file))

        # Verify saved file
        assert pptx_file.exists()
        assert pptx_file.stat().st_size > 0

        # Can reopen
        prs2 = Presentation()
        prs2.open(str(pptx_file))
        try:
            assert len(prs2.slides) == 3  # Initial blank + 2 added
        finally:
            prs2.close()

    def test_switch_theme_workflow(self, tmp_path):
        """Can create multiple versions with different themes."""
        themes_to_test = [
            ("aurora", create_aurora_theme()),
            ("ocean", create_ocean_theme()),
            ("sunset", create_sunset_theme()),
            ("forest", create_forest_theme()),
        ]

        for theme_name, theme in themes_to_test:
            pptx_file = tmp_path / f"{theme_name}.pptx"

            with Presentation.new(f"{theme.name} Demo") as prs:
                prs.apply_theme(theme)

                prs.slides[0].title = theme.name
                prs.slides[0].body = theme.metadata.get("description", "")

                prs.add_slide("Slide 2")

                prs.save(str(pptx_file))

            assert pptx_file.exists()

    def test_custom_branded_theme_workflow(self, tmp_path):
        """Create presentation with custom branded theme."""
        brand_theme = Theme(
            name="MyBrand",
            colors=ColorScheme(
                name="MyBrandColors",
                dk1="003366",
                lt1="FFFFFF",
                accent1="FF6B35",
                accent2="004E89",
                accent3="1B998B",
                accent4="2E8B57",
                accent5="B22222",
                accent6="FFD700",
            ),
            fonts=FontScheme(
                name="MyBrandFonts",
                major_font="Garamond",
                minor_font="Calibri",
            ),
            metadata={
                "brand": "MyBrand",
                "year": "2024",
            },
        )

        pptx_file = tmp_path / "branded.pptx"

        with Presentation.new("Brand Guidelines") as prs:
            prs.apply_theme(brand_theme)

            prs.slides[0].title = brand_theme.name
            prs.slides[0].body = "Official brand theme"

            for i in range(1, 4):
                prs.add_slide(f"Section {i}", layout="title_and_content")

            prs.save(str(pptx_file))

        assert pptx_file.exists()


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
