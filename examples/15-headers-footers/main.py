"""Example showing presentation-level and slide-level header/footer control."""

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def _build_demo_slides(prs: Presentation) -> None:
    """Create slides and apply footer settings."""
    prs.set_header_footer(
        footer="Confidential",
        show_footer=True,
        show_slide_num=True,
        show_date_time=True,
    )
    prs.add_slide("Slide 1", layout="title_only")
    prs.add_textbox(
        0,
        Inches(1),
        Inches(2),
        Inches(8),
        Inches(1),
        text="All slides have the same footer",
    )
    prs.add_slide("Slide 2", layout="title_only")
    prs.add_textbox(
        1,
        Inches(1),
        Inches(2),
        Inches(8),
        Inches(1),
        text="This slide also has the footer",
    )
    prs.add_slide("Slide 3", layout="title_only")
    prs.add_textbox(
        2, Inches(1), Inches(2), Inches(8), Inches(1), text="And this one too!"
    )
    prs.slides[1].set_header_footer(
        footer="Internal Use Only",
        show_footer=True,
        show_slide_num=False,
        show_date_time=True,
    )
    prs.slides[2].set_header_footer(
        footer="Do Not Distribute",
        show_footer=True,
        show_slide_num=True,
        show_date_time=False,
    )


def _print_footer_summary() -> None:
    """Print usage summary for header/footer API."""
    print("""
PRESENTATION-LEVEL (prs.set_header_footer):
  [+] Apply header/footer to ALL slides at once
  [+] Efficient - single method call instead of per-slide
  [+] Slides can override individual settings

PER-SLIDE LEVEL (slide.set_header_footer):
  [+] Override presentation-level settings
  [+] Each slide can have completely different footer
  [+] Mix and match: some slides default, others custom

USE CASES:
  1. All slides same footer --> Use prs.set_header_footer()
  2. Most slides same, few different --> Set presentation-level, override specific
  3. Every slide different --> Use slide.set_header_footer() on each
""")


def main() -> None:
    """Demonstrate presentation-level and per-slide header/footer control."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Headers and Footers Demo") as prs:
        _build_demo_slides(prs)

        print("\n" + "=" * 70)
        print("COMPARISON: What Each Slide Gets")
        print("=" * 70)
        print()
        print("Slide # | Footer Text          | Slide #  | Date/Time")
        print("-" * 70)
        print("Slide 1 | Confidential         | YES (1)  | YES     ")
        print("Slide 2 | Internal Use Only    | NO       | YES     ")
        print("Slide 3 | Do Not Distribute    | YES (3)  | NO      ")
        print()

        output_path = output_dir / "i03-python-headers-footers.pptx"
        prs.save(str(output_path))
        print(f"\nPresentation saved to {output_path}")

    print("\n" + "=" * 70)
    print("SUMMARY")
    print("=" * 70)
    _print_footer_summary()


if __name__ == "__main__":
    main()
