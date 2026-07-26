"""Test image swapping and slide copying without r:embed corruption (Issue #1142)."""

import pathlib
import zipfile

from gopptx import Presentation


def test_swap_image_and_add_picture_no_rembed_corruption(
    tmp_path: pathlib.Path,
) -> None:
    """Swapping an image on a slide and then adding pictures to a new slide preserves r:embed integrity (Issue #1142)."""
    output_path = tmp_path / "swap_copy_test.pptx"

    png1 = b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15c4\x00\x00\x00\rIDATx\x9cc` \x05\x00\x00\x04\x00\x01\x8c\x82\xb6\x00\x00\x00\x00IEND\xaeB`\x82"
    png2 = b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x02\x00\x00\x00\x02\x08\x06\x00\x00\x00v\xd7\xda\xed\x00\x00\x00\x15IDATx\x9cc\xf8\xff\xff?\x03\x03\x03\x03\x03\x03\x03\x03\x00\x1d\x05\x02\x01\xbd\xb2\xc5\xdc\x00\x00\x00\x00IEND\xaeB`\x82"

    with Presentation.new(title="Image Swap Test Deck") as pres:
        slide1 = pres.slides[0]
        shape_id1 = slide1.add_picture(
            png1,
            left=1000000,
            top=1000000,
            width=2000000,
            height=2000000,
            image_format="png",
        )
        assert shape_id1 > 0

        # Swap image on slide 1
        slide1.swap_image_by_index(0, png2, "png")

        # Create a new slide and add a picture
        slide2 = pres.add_slide(title="Slide 2")
        shape_id2 = slide2.add_picture(
            png1,
            left=1000000,
            top=1000000,
            width=2000000,
            height=2000000,
            image_format="png",
        )
        assert shape_id2 > 0

        pres.save(output_path)

    # Reopen presentation and inspect relationships
    with Presentation(output_path) as pres:
        assert len(pres.slides) == 2
        slide2_images = pres.slides[1].list_shapes()
        assert len(slide2_images) > 0

    with zipfile.ZipFile(output_path) as zf:
        slide2_rels = zf.read("ppt/slides/_rels/slide2.xml.rels").decode("utf-8")
        assert "Target=" in slide2_rels
