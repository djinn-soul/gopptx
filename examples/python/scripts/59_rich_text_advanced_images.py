import os
import pathlib
import sys

from gopptx import Presentation
from gopptx.schemas import Inches

# Add project root to sys.path to find 'gopptx' package
project_root = pathlib.Path(__file__).parent.parent.parent.parent.resolve()
sys.path.append(os.path.join(project_root, "python"))


def generate_example():
    output_dir = project_root / "examples" / "output"
    output_dir.mkdir(exist_ok=True)
    out_file = output_dir / "59_rich_text_advanced_images.pptx"

    # Path to a sample image
    img_path = project_root / "examples" / "assets" / "test_image.png"
    if not img_path.exists():
        # Create a dummy image if not exists
        try:
            import PIL.Image
            import PIL.ImageDraw

            img = PIL.Image.new("RGB", (200, 200), color=(73, 109, 137))
            d = PIL.ImageDraw.Draw(img)
            d.text((10, 10), "Sample Image", fill=(255, 255, 0))
            img.save(img_path)
        except ImportError:
            # Fallback if PIL not available (though it usually is in this dev env)
            pathlib.Path(img_path).write_bytes(b"dummy data")

    print(f"Generating example: {out_file}")

    with Presentation.new("Advanced Text and Images") as deck:
        # Slide 1: Rich Text and Formatting
        slide1 = deck.add_slide("Rich Text & Formatting")

        slide1.add_shape(
            "rect",
            (Inches(0.5), Inches(1.5), Inches(8), Inches(2.5)),
            runs=[
                {
                    "text": "This is a single shape with multiple runs:\n",
                    "bold": True,
                    "size_pt": 24,
                },
                {"text": "Bold, ", "bold": True},
                {"text": "Italic, ", "italic": True},
                {"text": "Underlined, ", "underline": "sng"},
                {"text": "Strikethrough, ", "strikethrough": "sng"},
                {"text": "Red text, ", "color": "FF0000"},
                {"text": "Highlighted, ", "highlight": "FFFF00"},
                {"text": "Subscript", "subscript": True},
                {"text": " and ", "bold": False},
                {"text": "Superscript", "superscript": True},
                {"text": ".\n\n", "bold": False},
                {
                    "text": "Click here for Google",
                    "hyperlink": {
                        "address": "https://google.com",
                        "tooltip": "Search engine",
                    },
                },
            ],
        )

        # Slide 2: Advanced Image Operations
        slide2 = deck.add_slide("Advanced Image Operations")

        # Original Image
        slide2.add_shape(
            "rect",
            (Inches(0.5), Inches(1), Inches(3), Inches(0.5)),
            text="Original Image",
        )
        slide2.add_image(
            str(img_path), (Inches(0.5), Inches(1.5), Inches(3), Inches(3))
        )

        # Cropped, Rotated, and Flipped Image
        slide2.add_shape(
            "rect",
            (Inches(4.5), Inches(1), Inches(4), Inches(0.5)),
            text="Cropped, 45\u00b0 Rotated, H-Flipped",
        )
        slide2.add_image(
            str(img_path),
            (Inches(5), Inches(1.5), Inches(3), Inches(3)),
            crop={"left": 0.2, "right": 0.2, "top": 0.2, "bottom": 0.2},
            rotation=45.0,
            flip_h=True,
        )

        # Slide 3: Shape Layer Click Actions
        slide3 = deck.add_slide("Shape-Level Interactivity")

        slide3.add_shape(
            "rect",
            (Inches(1), Inches(2), Inches(4), Inches(1)),
            text="This entire shape is a link to Github",
            click_action={
                "address": "https://github.com/djinn-soul/gopptx",
                "tooltip": "Go to Repository",
            },
        )

        # Slide 4: Text Frame Settings
        slide4 = deck.add_slide("Text Frame Settings")

        slide4.add_shape(
            "rect",
            (Inches(1), Inches(1.5), Inches(3), Inches(3)),
            text="No wrap, large margins, center align",
            text_frame={
                "margin_top": Inches(0.5),
                "margin_left": Inches(0.5),
                "word_wrap": False,
                "vertical_align": "ctr",
            },
        )

        deck.save(str(out_file))

    print("Generation complete.")
    return str(out_file)


if __name__ == "__main__":
    generate_example()
