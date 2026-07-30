"""Test picture shape .image property and ImagePartProxy (Issue #1084)."""

import base64
import pathlib

from gopptx import Presentation


def test_picture_shape_image_part(tmp_path: pathlib.Path) -> None:
    """Picture shapes expose .image with .blob, .ext, .content_type, .filename, and .sha1 (Issue #1084)."""
    output_path = tmp_path / "picture_image_test.pptx"

    # Valid 1x1 PNG image bytes via base64
    b64_str = (
        "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk"
        "+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
    )
    png_bytes = base64.b64decode(b64_str)

    with Presentation.new(title="Picture Image Part Test") as pres:
        slide = pres.slides[0]

        # Add picture using bytes
        pic_id = slide.add_picture(
            png_bytes, left=1000000, top=1000000, width=2000000, height=2000000
        )
        assert pic_id > 0

        pres.save(output_path)

    # Re-open presentation and inspect shape.image
    with Presentation(output_path) as pres:
        slide = pres.slides[0]
        # Find picture shape
        pic_shape = None
        for shape in slide.shapes:
            if shape.id == pic_id:
                pic_shape = shape
                break

        assert pic_shape is not None
        img_part = pic_shape.image
        assert img_part is not None
        assert img_part.blob == png_bytes
        assert img_part.ext in {"png", "image/png"}
        assert img_part.content_type in {"image/png", "png"}
        assert img_part.filename.endswith(".png") or img_part.filename.endswith("png")
        assert len(img_part.sha1) == 40

        # Test saving image to path
        saved_img_path = tmp_path / "extracted_picture.png"
        img_part.save(saved_img_path)
        assert saved_img_path.exists()
        assert saved_img_path.read_bytes() == png_bytes
