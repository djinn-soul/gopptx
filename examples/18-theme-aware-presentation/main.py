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


def add_theme_slide(prs: Presentation, title: str, points: list[str]) -> object:
    """Add a themed slide with title and bullet points."""
    return prs.add_slide(
        title,
        layout="title_and_content",
        bullets=points,
    )


def create_custom_corporate_theme() -> Theme:
    """Create a custom corporate theme with brand colors."""
    return Theme(
        name="Corporate Blue",
        colors=ColorScheme(
            name="Corporate Blue",
            dk1="003366",
            lt1="FFFFFF",
            dk2="004499",
            lt2="E6F0FF",
            accent1="0066CC",
            accent2="FF6600",
            accent3="009933",
            accent4="CC0000",
            accent5="9933CC",
            accent6="00CC99",
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


def _build_aurora_demo() -> None:
    """Create and save Aurora theme presentation."""
    with Presentation.new("Theme-Aware Presentation") as prs:
        aurora_theme = get_theme("aurora")
        prs.apply_theme(aurora_theme)
        print(f"   Applied: {aurora_theme.name}")
        print(f"   Colors: Primary={aurora_theme.colors.accent1}, Secondary={aurora_theme.colors.accent2}")
        print(f"   Fonts: Major={aurora_theme.fonts.major_font}, Minor={aurora_theme.fonts.minor_font}")
        prs.slides[0].title = "Theme-Aware Presentation"
        prs.slides[0].body = "Demonstrating unified visual identity with themes"
        add_theme_slide(prs, "What are Themes?", ["Unified color palette across all slides", "Consistent font families and sizes", "Professional, cohesive visual identity", "Easy to apply - one line of code!"])
        add_theme_slide(prs, "Aurora Theme", ["Dark navy text on light backgrounds", "Blue and teal accent colors", "Modern, professional feel", "Great for tech and corporate presentations"])
        add_theme_slide(prs, "Theme Components", ["Color Scheme: dk1, lt1, dk2, lt2, accent1-6, hyperlinks", "Font Scheme: major (heading) and minor (body) fonts", "Metadata: description, author, custom fields", "All components work together for cohesion"])
        prs.save("examples/output/18-theme-aware-presentation.pptx")
        print("\n   Saved: examples/output/18-theme-aware-presentation.pptx")


def _build_single_theme(name: str, slide_title: str, body: str, points: list[str], out: str) -> None:
    """Create and save a single built-in theme demo presentation."""
    with Presentation.new(f"{name.capitalize()} Theme Demo") as prs:
        theme = get_theme(name)
        prs.apply_theme(theme)
        prs.slides[0].title = f"{name.capitalize()} Theme"
        prs.slides[0].body = body
        add_theme_slide(prs, slide_title, points)
        prs.save(out)
        print(f"   Saved: {out}")


def _build_corporate_demo() -> None:
    """Create and save corporate custom theme presentation."""
    with Presentation.new("Custom Corporate Theme") as prs:
        corporate_theme = create_custom_corporate_theme()
        prs.apply_theme(corporate_theme)
        prs.slides[0].title = "Custom Corporate Theme"
        prs.slides[0].body = "Building your own branded themes"
        add_theme_slide(prs, "Creating Custom Themes", ["Define custom color schemes with your brand colors", "Choose fonts that match your brand identity", "Reuse themes across presentations", "Maintain consistency across all presentations"])
        add_theme_slide(prs, "Custom Theme Features", [f"Primary Color: {corporate_theme.colors.accent1}", f"Secondary Color: {corporate_theme.colors.accent2}", f"Major Font: {corporate_theme.fonts.major_font}", "Metadata support for documentation"])
        prs.save("examples/output/18-theme-custom-corporate.pptx")
        print("   Saved: examples/output/18-theme-custom-corporate.pptx")


def main() -> None:
    """Create presentations demonstrating theme system."""
    print("\n" + "=" * 70)
    print("THEME-AWARE PRESENTATION DEMO")
    print("=" * 70)

    print("\n1. Available Built-in Themes:")
    for theme_name in list_themes():
        print(f"   - {theme_name.upper()}")

    print("\n2. Creating presentation with Aurora theme...")
    _build_aurora_demo()

    print("\n3. Creating presentation with Ocean theme...")
    _build_single_theme("ocean", "Professional & Calm", "Deep blues and greens", ["Deep ocean blues for depth", "Teal and emerald accents", "Serene, professional aesthetic", "Perfect for corporate and finance"], "examples/output/18-theme-ocean.pptx")

    print("\n4. Creating presentation with Sunset theme...")
    _build_single_theme("sunset", "Energetic & Vibrant", "Warm oranges and reds", ["Warm sunset colors", "Orange, red, and gold accents", "Energetic, vibrant feel", "Great for creative and startup presentations"], "examples/output/18-theme-sunset.pptx")

    print("\n5. Creating presentation with Forest theme...")
    _build_single_theme("forest", "Calm & Organic", "Natural greens and earth tones", ["Natural green palette", "Earth tones and sage accents", "Calm, organic aesthetic", "Ideal for sustainability and environmental topics"], "examples/output/18-theme-forest.pptx")

    print("\n6. Creating presentation with Custom Corporate theme...")
    _build_corporate_demo()

    print("\n" + "=" * 70)
    print("SUMMARY")
    print("=" * 70)
    print("\nTheme System Benefits:")
    print("  [+] Consistent visual identity across all slides")
    print("  [+] Professional, polished appearance")
    print("  [+] Easy to switch themes (one line of code)")
    print("  [+] Built-in themes for immediate use")
    print("  [+] Custom themes for brand consistency")
    print("=" * 70 + "\n")


if __name__ == "__main__":
    main()
