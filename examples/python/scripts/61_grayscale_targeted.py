"""Generate a grayscale demo deck with isolated targeting examples.

Outputs:
  - examples/output/61_grayscale_targeted_source.pptx
  - examples/output/61_grayscale_targeted_result.pptx
"""

from __future__ import annotations

from pathlib import Path

from gopptx import GrayscaleScope, PlaceholderType, Presentation, Slide


def find_shape_by_placeholder(
    slide: Slide, placeholder_type: str | PlaceholderType
) -> int:
    for shape in slide.list_shapes():
        if shape.get("PlaceholderType") == placeholder_type:
            return int(shape["ID"])
    raise ValueError(
        f"placeholder type {placeholder_type!r} not found on slide {slide.index}"
    )


def print_inventory(prs: Presentation) -> None:
    print("Slide inventory:")
    for slide in prs.slides:
        print(f"slide_index={slide.index}, title={slide.title!r}")
        print("  placeholders:")
        for placeholder in slide.list_placeholders():
            print(
                "   ",
                placeholder.get("index"),
                placeholder.get("type"),
                placeholder.get("name"),
            )
        print("  shapes:")
        for shape in slide.list_shapes():
            print(
                "   ",
                shape.get("ID"),
                shape.get("Name"),
                shape.get("Text"),
                "ph_type=",
                shape.get("PlaceholderType"),
                "ph_index=",
                shape.get("PlaceholderIndex"),
            )


def main() -> None:
    root = Path(__file__).resolve().parents[3]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)

    source_path = output_dir / "61_grayscale_targeted_source.pptx"
    result_path = output_dir / "61_grayscale_targeted_result.pptx"

    with Presentation.new("Grayscale Targeted Example") as prs:
        # Slide 0: title placeholder lookup by enum-style placeholder type.
        title_slide = prs.slides[0]
        title_slide.title = "Title Placeholder Target"
        title_slide.set_background("solid", color="F6F2EA")

        # Slide 1: footer placeholder lookup by placeholder type.
        footer_slide = prs.add_slide("Footer Placeholder Target")
        footer_slide.set_background("solid", color="F6F2EA")
        footer_slide.add_shape(
            "rect",
            (914400, 1524000, 4572000, 914400),
            text="Only the footer should grayscale here",
            properties={"fill": {"solid": "FFF7ED"}, "line": {"color": "9A3412"}},
        )
        footer_slide.set_header_footer(
            footer="Confidential Footer",
            show_footer=True,
            show_slide_num=True,
        )

        # Slide 2: specific text run by shape_id/run_indices.
        text_slide = prs.add_slide("Text Run Target")
        text_slide.set_background("solid", color="EFEFEF")
        text_shape_id = text_slide.add_shape(
            "rect",
            (914400, 914400, 4572000, 1143000),
            text="placeholder",
        )
        text_slide.set_shape_runs(
            text_shape_id,
            [
                {"text": "Run 0 turns grayscale", "color": "FF0000"},
                {"text": " | Run 1 stays blue", "color": "0000FF"},
            ],
        )

        # Slide 3: specific shape by shape_id.
        shape_slide = prs.add_slide("Shape Target")
        shape_slide.set_background("solid", color="EFEFEF")
        accent_shape_id = shape_slide.add_shape(
            "ellipse",
            (914400, 2286000, 1828800, 914400),
            text="Accent shape",
            properties={"fill": {"solid": "FF6A00"}, "line": {"color": "1A2B8F"}},
        )
        shape_slide.add_shape(
            "rect",
            (3200400, 2286000, 1828800, 914400),
            text="Control shape (untouched)",
            properties={"fill": {"solid": "00A676"}, "line": {"color": "004D40"}},
        )

        # Slide 4: slide background only.
        background_slide = prs.add_slide("Background Target")
        background_slide.set_background("solid", color="3070B3")
        background_slide.add_shape(
            "rect",
            (914400, 1524000, 4572000, 1143000),
            text="Only background should grayscale on this slide",
            properties={"fill": {"solid": "FFFFFF"}, "line": {"color": "1F2937"}},
        )

        # Slide 5: embedded image only.
        image_slide = prs.add_slide("Image Target")
        image_slide.set_background("solid", color="EFEFEF")
        image_shape_id = image_slide.add_image(
            str(root / "examples" / "assets" / "test_image.png"),
            (1828800, 1371600, 2743200, 2057400),
        )
        image_slide.add_shape(
            "rect",
            (914400, 914400, 1828800, 914400),
            text="Color shape should remain",
            properties={"fill": {"solid": "8B5CF6"}, "line": {"color": "4C1D95"}},
        )

        title_shape_id = find_shape_by_placeholder(title_slide, PlaceholderType.TITLE)
        footer_shape_id = find_shape_by_placeholder(
            footer_slide, PlaceholderType.FOOTER
        )
        title_slide.set_shape_runs(
            title_shape_id,
            [
                {"text": "Title placeholder", "color": "D62828"},
                {"text": " stays discoverable", "color": "1D4ED8"},
            ],
        )
        footer_slide.set_shape_runs(
            footer_shape_id,
            [
                {"text": "Confidential Footer", "color": "C2410C"},
            ],
        )

        prs.save(str(source_path))
        print_inventory(prs)

        prs.convert_to_grayscale(
            placeholders=[
                {"slide_index": title_slide.index, "type": PlaceholderType.TITLE}
            ],
            scope=GrayscaleScope(colors=True, images=False, backgrounds=False),
        )
        prs.convert_to_grayscale(
            placeholders=[
                {"slide_index": footer_slide.index, "type": PlaceholderType.FOOTER}
            ],
            scope=GrayscaleScope(colors=True, images=False, backgrounds=False),
        )
        prs.convert_to_grayscale(
            text=[
                {
                    "slide_index": text_slide.index,
                    "shape_id": text_shape_id,
                    "run_indices": [0],
                }
            ],
            scope=GrayscaleScope(colors=True, images=False, backgrounds=False),
        )
        prs.convert_to_grayscale(
            shapes=[{"slide_index": shape_slide.index, "shape_id": accent_shape_id}],
            scope=GrayscaleScope(colors=True, images=False, backgrounds=False),
        )
        prs.convert_to_grayscale(
            slides=[background_slide.index],
            scope=GrayscaleScope(colors=False, images=False, backgrounds=True),
        )
        prs.convert_to_grayscale(
            shapes=[{"slide_index": image_slide.index, "shape_id": image_shape_id}],
            scope=GrayscaleScope(colors=False, images=True, backgrounds=False),
        )
        prs.save(str(result_path))

    print(f"Created source deck: {source_path}")
    print(f"Created grayscale deck: {result_path}")


if __name__ == "__main__":
    main()
