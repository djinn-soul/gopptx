"""Advanced media API integration tests."""

import pathlib
from io import BytesIO

from gopptx import Presentation
from gopptx.presentation.shapes.shape_media_mixin import PresentationShapeMediaMixin
from gopptx.schemas import Inches
from PIL import Image


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


def test_add_picture_with_description(tmp_path: pathlib.Path) -> None:
    """add_picture creates image shape with description, alt_text, and title options."""
    img_path = tmp_path / "chart_pic.png"
    img_path.write_bytes(
        b"\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n2\xb4\x00\x00\x00\x00IEND\xaeB`\x82"
    )

    output_path = tmp_path / "picture_output.pptx"
    with Presentation.new("Picture Test") as pres:
        slide = pres.slides[0]
        pic_id = slide.add_picture(
            img_path,
            left=Inches(1),
            top=Inches(1),
            width=Inches(3),
            height=Inches(2),
            description="Sample Alt Text Description",
            title="Sample Title",
        )
        assert pic_id > 0
        pres.save(output_path)

    assert output_path.exists()


def _png_bytes(width: int = 144, height: int = 72) -> bytes:
    buffer = BytesIO()
    Image.new("RGB", (width, height), "blue").save(buffer, format="PNG")
    return buffer.getvalue()


def test_presentation_add_picture_infers_byte_format_and_native_size() -> None:
    """Presentation-level picture insertion inspects byte sources."""

    class _FakeShapeMedia(PresentationShapeMediaMixin):
        def __init__(self) -> None:
            self.last_payload: dict[str, object] = {}

        def execute(self, _op: str, payload: dict[str, object]) -> dict[str, int]:
            self.last_payload = payload
            return {"shape_id": 1}

    fake = _FakeShapeMedia()
    _ = fake.add_picture(0, _png_bytes(), left=100, top=200)
    assert fake.last_payload["format"] == "png"
    assert fake.last_payload["w"] == Inches(2)
    assert fake.last_payload["h"] == Inches(1)


def test_presentation_add_picture_accepts_data_only_source() -> None:
    """Presentation-level insertion resolves the documented data option."""

    class _FakeShapeMedia(PresentationShapeMediaMixin):
        def __init__(self) -> None:
            self.last_payload: dict[str, object] = {}

        def execute(self, _op: str, payload: dict[str, object]) -> dict[str, int]:
            self.last_payload = payload
            return {"shape_id": 1}

    fake = _FakeShapeMedia()
    _ = fake.add_picture(0, None, data=_png_bytes())
    assert fake.last_payload["format"] == "png"
    assert fake.last_payload["w"] == Inches(2)
    assert fake.last_payload["h"] == Inches(1)


def test_slide_add_picture_accepts_data_only_source(tmp_path: pathlib.Path) -> None:
    """Slide-level insertion forwards a data-only picture to the bridge."""
    output_path = tmp_path / "data_only_picture.pptx"
    with Presentation.new("Data-only Picture") as pres:
        slide = pres.slides[0]
        picture_id = slide.add_picture(None, data=_png_bytes())
        metadata = slide.get_image_metadata(picture_id)
        assert metadata["format"] == "png"
        pres.save(output_path)

    assert output_path.exists()


def test_slide_add_picture_preserves_aspect_ratio_for_omitted_dimension(
    tmp_path: pathlib.Path,
) -> None:
    """Slide-level insertion derives height and embeds inferred PNG bytes."""
    output_path = tmp_path / "native_picture_size.pptx"
    with Presentation.new("Native Picture Size") as pres:
        slide = pres.slides[0]
        picture_id = slide.add_picture(
            _png_bytes(),
            left=Inches(1),
            top=Inches(1),
            width=Inches(1),
        )
        shape = next(item for item in slide.list_shapes() if item["ID"] == picture_id)
        assert shape["W"] == Inches(1)
        assert shape["H"] == Inches(0.5)
        metadata = slide.get_image_metadata(picture_id)
        assert metadata["format"] == "png"
        pres.save(output_path)

    assert output_path.exists()
