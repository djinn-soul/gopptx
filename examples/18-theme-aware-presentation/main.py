"""Example showing theme-aware presentation creation.

This example demonstrates:
- Applying built-in themes (Aurora, Ocean, Sunset, Forest)
- Creating custom themes with color and font schemes
- Using themes to maintain consistent visual identity across slides
- Theme discovery with list_themes() and get_theme()
"""

from gopptx import Presentation
from gopptx.presentation.theme import (
    ColorScheme,
    FontScheme,
    Theme,
    get_theme,
    list_themes,
)
from gopptx.schemas import Inches


def add_theme_slide(prs, title, points):
    """Add a themed slide with title and bullet points."""
    slide = prs.add_slide(
        title,
        layout="title_and_content",
        bullets=points,
    )
    return slide


def create_custom_corporate_theme():
    """Create a custom corporate theme with brand colors."""
    return Theme(
        name="Corporate Blue",
        colors=ColorScheme(
            name="Corporate Blue",
            dk1="003366",  # Dark navy
            lt1="FFFFFF",  # White
            dk2="004499",  # Navy
            lt2="E6F0FF",  # Light blue
            accent1="0066CC",  # Primary blue
            accent2="FF6600",  # Orange accent
            accent3="009933",  # Green accent
            accent4="CC0000",  # Red accent
            accent5="9933CC",  # Purple accent
            accent6="00CC99",  # Teal accent
            hlink="0066CC",
        ),
        fonts=FontScheme(
            name="Corporate Blue",
            major_font="Calibri Light",
            minor_font="Calibri",
        ),
        metadata={
            "description": "Professional corporate theme",
            "brand": "ACME Corp",
            "author": "gopptx",
        },
    )


def main():
    """Create presentations demonstrating theme system."""

    # ===== Part 1: Built-in Themes =====
    print("\n" + "=" * 70)
    print("THEME-AWARE PRESENTATION DEMO")
    print("=" * 70)

    print("\n1. Available Built-in Themes:")
    for theme_name in list_themes():
        print(f"   - {theme_name.upper()}")

    # ===== Part 2: Aurora Theme (default in this example) =====
    print("\n2. Creating presentation with Aurora theme...")
    with Presentation.new("Theme-Aware Presentation") as prs:
        aurora_theme = get_theme("aurora")
        prs.apply_theme(aurora_theme)

        print(f"   Applied: {aurora_theme.name}")
        print(
            f"   Colors: Primary={aurora_theme.colors.accent1}, "
            f"Secondary={aurora_theme.colors.accent2}"
        )
        print(
            f"   Fonts: Major={aurora_theme.fonts.major_font}, "
            f"Minor={aurora_theme.fonts.minor_font}"
        )

        # Title slide
        prs.slides[0].title = "Theme-Aware Presentation"
        prs.slides[0].body = "Demonstrating unified visual identity with themes"

        # Content slides
        add_theme_slide(
            prs,
            "What are Themes?",
            [
                "Unified color palette across all slides",
                "Consistent font families and sizes",
                "Professional, cohesive visual identity",
                "Easy to apply - one line of code!",
            ],
        )

        add_theme_slide(
            prs,
            "Aurora Theme",
            [
                "Dark navy text on light backgrounds",
                "Blue and teal accent colors",
                "Modern, professional feel",
                "Great for tech and corporate presentations",
            ],
        )

        add_theme_slide(
            prs,
            "Theme Components",
            [
                "Color Scheme: dk1, lt1, dk2, lt2, accent1-6, hyperlinks",
                "Font Scheme: major (heading) and minor (body) fonts",
                "Metadata: description, author, custom fields",
                "All components work together for cohesion",
            ],
        )

        # Save
        prs.save("examples/output/18-theme-aware-presentation.pptx")
        print("\n   Saved: examples/output/18-theme-aware-presentation.pptx")

    # ===== Part 3: Ocean Theme =====
    print("\n3. Creating presentation with Ocean theme...")
    with Presentation.new("Ocean Theme Demo") as prs:
        ocean_theme = get_theme("ocean")
        prs.apply_theme(ocean_theme)

        prs.slides[0].title = "Ocean Theme"
        prs.slides[0].body = "Deep blues and greens"

        add_theme_slide(
            prs,
            "Professional & Calm",
            [
                "Deep ocean blues for depth",
                "Teal and emerald accents",
                "Serene, professional aesthetic",
                "Perfect for corporate and finance",
            ],
        )

        prs.save("examples/output/18-theme-ocean.pptx")
        print("   Saved: examples/output/18-theme-ocean.pptx")

    # ===== Part 4: Sunset Theme =====
    print("\n4. Creating presentation with Sunset theme...")
    with Presentation.new("Sunset Theme Demo") as prs:
        sunset_theme = get_theme("sunset")
        prs.apply_theme(sunset_theme)

        prs.slides[0].title = "Sunset Theme"
        prs.slides[0].body = "Warm oranges and reds"

        add_theme_slide(
            prs,
            "Energetic & Vibrant",
            [
                "Warm sunset colors",
                "Orange, red, and gold accents",
                "Energetic, vibrant feel",
                "Great for creative and startup presentations",
            ],
        )

        prs.save("examples/output/18-theme-sunset.pptx")
        print("   Saved: examples/output/18-theme-sunset.pptx")

    # ===== Part 5: Forest Theme =====
    print("\n5. Creating presentation with Forest theme...")
    with Presentation.new("Forest Theme Demo") as prs:
        forest_theme = get_theme("forest")
        prs.apply_theme(forest_theme)

        prs.slides[0].title = "Forest Theme"
        prs.slides[0].body = "Natural greens and earth tones"

        add_theme_slide(
            prs,
            "Calm & Organic",
            [
                "Natural green palette",
                "Earth tones and sage accents",
                "Calm, organic aesthetic",
                "Ideal for sustainability and environmental topics",
            ],
        )

        prs.save("examples/output/18-theme-forest.pptx")
        print("   Saved: examples/output/18-theme-forest.pptx")

    # ===== Part 6: Custom Corporate Theme =====
    print("\n6. Creating presentation with Custom Corporate theme...")
    with Presentation.new("Custom Corporate Theme") as prs:
        corporate_theme = create_custom_corporate_theme()
        prs.apply_theme(corporate_theme)

        prs.slides[0].title = "Custom Corporate Theme"
        prs.slides[0].body = "Building your own branded themes"

        add_theme_slide(
            prs,
            "Creating Custom Themes",
            [
                "Define custom color schemes with your brand colors",
                "Choose fonts that match your brand identity",
                "Reuse themes across presentations",
                "Maintain consistency across all presentations",
            ],
        )

        add_theme_slide(
            prs,
            "Custom Theme Features",
            [
                f"Primary Color: {corporate_theme.colors.accent1}",
                f"Secondary Color: {corporate_theme.colors.accent2}",
                f"Major Font: {corporate_theme.fonts.major_font}",
                "Metadata support for documentation",
            ],
        )

        prs.save("examples/output/18-theme-custom-corporate.pptx")
        print("   Saved: examples/output/18-theme-custom-corporate.pptx")

    # ===== Summary =====
    print("\n" + "=" * 70)
    print("SUMMARY")
    print("=" * 70)
    print("\nTheme System Benefits:")
    print("  [+] Consistent visual identity across all slides")
    print("  [+] Professional, polished appearance")
    print("  [+] Easy to switch themes (one line of code)")
    print("  [+] Built-in themes for immediate use")
    print("  [+] Custom themes for brand consistency")
    print("\nHow to Use:")
    print("  1. Import: from gopptx.presentation.theme import get_theme, Theme")
    print("  2. Get theme: theme = get_theme('aurora')")
    print("  3. Apply: prs.apply_theme(theme)")
    print("  4. Add content: All slides automatically use theme colors/fonts")
    print("\nCreate Custom Theme:")
    print("  - Define ColorScheme with your brand colors")
    print("  - Define FontScheme with your brand fonts")
    print("  - Combine into Theme instance")
    print("  - Apply like any built-in theme")
    print("=" * 70 + "\n")


if __name__ == "__main__":
    main()
