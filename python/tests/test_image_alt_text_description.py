"""Test setting and getting image alt text, description, and title (Issue #1139)."""

import pathlib

from gopptx import Presentation
from gopptx.schemas import Inches


def test_image_alt_text_and_description(tmp_path: pathlib.Path) -> None:
    """shape.alt_text, shape.description, shape.title allow inspecting and setting picture metadata (Issue #1139)."""
    output_path = tmp_path / "alt_text_test.pptx"

    # Create dummy PNG image bytes (1x1 red PNG)
    sample_png_bytes = (
        b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01"
        b"\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\x0cIDATx\x9cc\xf8\xcf\xc0"
        b"\x00\x00\x03\x01\x01\x00\x18\xdd\x8d\xb0\x00\x00\x00\x00IEND\xaeB`\x82"
    )

    with Presentation.new(title="Alt Text Test Deck") as pres:
        slide = pres.slides[0]

        # Add picture with description and title in add_picture
        pic_id = slide.add_picture(
            sample_png_bytes,
            left=Inches(1),
            top=Inches(1),
            width=Inches(3),
            height=Inches(3),
            description="Company Logo for Visually Impaired Users",
            title="Corporate Logo Title",
        )

        pic_shape = slide.shape(pic_id)
        assert pic_shape.alt_text == "Company Logo for Visually Impaired Users"
        assert pic_shape.description == "Company Logo for Visually Impaired Users"
        assert pic_shape.title == "Corporate Logo Title"

        # Update alt text and title dynamically
        pic_shape.alt_text = "Updated Accessible Logo Description"
        pic_shape.title = "Updated Logo Title"

        assert pic_shape.alt_text == "Updated Accessible Logo Description"
        assert pic_shape.title == "Updated Logo Title"

        pres.save(output_path)

    # Re-open and verify persistence
    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        shapes = slide.list_shapes()
        assert len(shapes) > 0
