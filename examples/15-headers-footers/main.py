"""Example showing presentation-level and slide-level header/footer control."""

from pathlib import Path

from gopptx import Presentation
from gopptx.schemas import Inches


def main():
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Headers and Footers Demo") as prs:
        # =======================================================================
        # Example 1: Presentation-level (ALL slides get same footer)
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 1: Presentation-Level Headers/Footers")
        print("=" * 70)
        print("\nApplying to ALL slides:")
        print("  - Footer: 'Confidential'")
        print("  - Slide numbers: ON")
        print("  - Date/Time: ON")

        # Set header/footer for ALL slides at once
        prs.set_header_footer(
            footer="Confidential",
            show_footer=True,
            show_slide_num=True,
            show_date_time=True,
        )

        # Add slides - they all inherit the presentation-level footer
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
            2,
            Inches(1),
            Inches(2),
            Inches(8),
            Inches(1),
            text="And this one too!",
        )

        print("\nResult: Slides 1, 2, 3 all have 'Confidential' footer + slide #")

        # =======================================================================
        # Example 2: Per-Slide Override (change specific slide)
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 2: Per-Slide Override")
        print("=" * 70)

        # Override on Slide 2 - different footer, no slide numbers
        print("\nOverriding Slide 2:")
        print("  - Footer: 'Internal Use Only'")
        print("  - Slide numbers: OFF (override)")
        print("  - Date/Time: ON (inherited)")

        prs.slides[1].set_header_footer(
            footer="Internal Use Only",
            show_footer=True,
            show_slide_num=False,  # Turn off slide numbers on this one
            show_date_time=True,
        )

        print("\nResult:")
        print("  - Slide 1: 'Confidential' + slide #1 + date")
        print("  - Slide 2: 'Internal Use Only' + date (NO slide #)")
        print("  - Slide 3: 'Confidential' + slide #3 + date")

        # =======================================================================
        # Example 3: Per-Slide Custom (completely different)
        # =======================================================================
        print("\n" + "=" * 70)
        print("EXAMPLE 3: Complete Per-Slide Customization")
        print("=" * 70)

        # Slide 3 - totally different setup
        print("\nCustomizing Slide 3:")
        print("  - Footer: 'Do Not Distribute'")
        print("  - Slide numbers: ON")
        print("  - Date/Time: OFF (override)")

        prs.slides[2].set_header_footer(
            footer="Do Not Distribute",
            show_footer=True,
            show_slide_num=True,
            show_date_time=False,  # Turn off date/time on this one
        )

        print("\nResult:")
        print("  - Slide 1: 'Confidential' + slide #1 + date")
        print("  - Slide 2: 'Internal Use Only' + date (NO slide #)")
        print("  - Slide 3: 'Do Not Distribute' + slide #3 (NO date)")

        # =======================================================================
        # Example 4: Comparison Table
        # =======================================================================
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

        # Save presentation
        output_path = output_dir / "i03-python-headers-footers.pptx"
        prs.save(str(output_path))
        print(f"\nPresentation saved to {output_path}")

        # =======================================================================
        # Print Summary
        # =======================================================================
        print("\n" + "=" * 70)
        print("SUMMARY")
        print("=" * 70)
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


if __name__ == "__main__":
    main()
