"""Advanced media API integration tests."""

import pathlib

from gopptx import Presentation
from gopptx.presentation.shapes.shape_media_mixin import PresentationShapeMediaMixin
from gopptx.schemas import Inches


def test_image_advanced(tmp_path: pathlib.Path) -> None:
    """Add images with options and validate metadata."""
    # Create a dummy image
    img_path = tmp_path / "test.png"
    img_path.write_bytes(
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
        if shape_id <= 0:
            raise AssertionError("image shape id must be positive")

        # Test 2: Get metadata
        meta = pres.get_image_metadata(0, shape_id)
        if meta["width"] != 1:
            raise AssertionError("expected image width 1")
        if meta["height"] != 1:
            raise AssertionError("expected image height 1")
        if meta["format"] != "png":
            raise AssertionError("expected png format")

        # Test 3: Add from bytes
        data = pathlib.Path(img_path).read_bytes()

        shape_id2 = slide.add_image(
            None,
            (Inches(4), Inches(1), Inches(2), Inches(2)),
            data=data,
            image_format="png",
        )
        if shape_id2 <= shape_id:
            raise AssertionError("expected second image id to increase")

        pres.save(output_path)


def test_media_full(tmp_path: pathlib.Path) -> None:
    """Add video and OLE payloads and save output deck."""
    # Create dummy media files
    video_path = tmp_path / "test.mp4"
    pathlib.Path(video_path).write_bytes(b"fake video data")

    poster_path = tmp_path / "poster.png"
    poster_path.write_bytes(
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
        if video_id <= 0:
            raise AssertionError("video id must be positive")

        # Test 2: Add OLE Object
        ole_id = slide.add_ole_object(
            excel_path,
            (Inches(1), Inches(5), Inches(2), Inches(2)),
            prog_id="Excel.Sheet.12",
            icon=poster_path,
        )
        if ole_id <= video_id:
            raise AssertionError("ole id must be greater than video id")

        pres.save(output_path)

    if not output_path.exists():
        raise AssertionError("expected output deck to exist")


def test_add_audio_uses_icon_payload_keys() -> None:
    """Audio icon payload uses backend-supported icon_path/icon_data keys."""

    class _FakeShapeMedia(PresentationShapeMediaMixin):
        def __init__(self) -> None:
            self.last_payload: dict[str, object] = {}

        def execute(self, _op: str, payload: dict[str, object]) -> dict[str, int]:
            self.last_payload = payload
            return {"shape_id": 1}

    fake = _FakeShapeMedia()
    _ = fake.add_audio(
        0,
        b"audio",
        (0, 0, 1, 1),
        icon=b"icon-bytes",
    )
    assert "icon_data" in fake.last_payload
    assert "poster_data" not in fake.last_payload
