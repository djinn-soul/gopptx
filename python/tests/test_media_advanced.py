import pathlib

from gopptx import Presentation
from gopptx.schemas import Inches


def test_image_advanced(tmp_path):
    # Create a dummy image
    img_path = tmp_path / "test.png"
    with pathlib.Path(img_path).open("wb") as f:
        # Minimal valid PNG
        f.write(
            b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n2\xb4\x00\x00\x00\x00IEND\xaeB`\x82"
        )

    output_path = tmp_path / "output.pptx"
    with Presentation.new("Media Test") as pres:
        slide = pres.slides[0]

        # Test 1: Add image with crop and rotation
        shape_id = slide.add_image(
            img_path,
            (Inches(1), Inches(1), Inches(2), Inches(2)),
            crop={"left": 0.1, "right": 0.1, "top": 0.1, "bottom": 0.1},
            rotation=45.0,
            flip_h=True,
        )
        assert shape_id > 0

        # Test 2: Get metadata
        meta = pres.get_image_metadata(0, shape_id)
        assert meta["width"] == 1
        assert meta["height"] == 1
        assert meta["format"] == "png"

        # Test 3: Add from bytes
        data = pathlib.Path(img_path).read_bytes()

        shape_id2 = slide.add_image(
            None, (Inches(4), Inches(1), Inches(2), Inches(2)), data=data, format="png"
        )
        assert shape_id2 > shape_id

        pres.save(output_path)


def test_media_full(tmp_path):
    # Create dummy media files
    video_path = tmp_path / "test.mp4"
    pathlib.Path(video_path).write_bytes(b"fake video data")

    poster_path = tmp_path / "poster.png"
    with pathlib.Path(poster_path).open("wb") as f:
        # Minimal valid PNG
        f.write(
            b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n2\xb4\x00\x00\x00\x00IEND\xaeB`\x82"
        )

    excel_path = tmp_path / "test.xlsx"
    pathlib.Path(excel_path).write_bytes(b"fake excel data")

    output_path = tmp_path / "media_output.pptx"
    with Presentation.new("Media Integration") as pres:
        slide = pres.slides[0]

        # Test 1: Add Video
        video_id = slide.add_video(
            video_path,
            (Inches(1), Inches(1), Inches(4), Inches(3)),
            poster_frame=poster_path,
            mime_type="video/mp4",
        )
        assert video_id > 0

        # Test 2: Add OLE Object
        ole_id = slide.add_ole_object(
            excel_path,
            (Inches(1), Inches(5), Inches(2), Inches(2)),
            prog_id="Excel.Sheet.12",
            icon=poster_path,
        )
        assert ole_id > video_id

        pres.save(output_path)

    assert output_path.exists()
