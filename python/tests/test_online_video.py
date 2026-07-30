"""Online video linking tests (upstream python-pptx issue #1034)."""

import pathlib
import zipfile

import pytest
from gopptx.api_errors import GopptxError
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches

URL = "https://example.com/videos/demo.mp4"
BOUNDS = (Inches(2.2), Inches(1.6), Inches(5.6), Inches(3.2))


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "video.pptx"
    with Presentation.new("Online video") as pres:
        pres.save(output_path)
    return output_path


def test_online_video_uses_external_relationship(deck: pathlib.Path) -> None:
    """The media is linked, not embedded: that is what the issue asks for."""
    pres = Presentation()
    pres.open(deck)
    shape_id = pres.add_online_video(0, URL, BOUNDS, alt_text="Demo")
    assert shape_id > 0
    pres.save(deck)
    pres.close()

    with zipfile.ZipFile(deck) as archive:
        rels = archive.read("ppt/slides/_rels/slide1.xml.rels").decode()
        slide_xml = archive.read("ppt/slides/slide1.xml").decode()
        media = [n for n in archive.namelist() if n.startswith("ppt/media/")]

    assert 'TargetMode="External"' in rels
    assert URL in rels
    assert "<a:videoFile" in slide_xml
    # Only the poster image is embedded; the video itself is not.
    assert all(not name.endswith(".mp4") for name in media)


def test_online_video_default_poster_is_visible(deck: pathlib.Path) -> None:
    """A video renders as its poster, so the default one cannot be blank."""
    pres = Presentation()
    pres.open(deck)
    pres.add_online_video(0, URL, BOUNDS)
    pres.save(deck)
    pres.close()

    with zipfile.ZipFile(deck) as archive:
        posters = [n for n in archive.namelist() if n.startswith("ppt/media/image")]
        assert posters
        poster = archive.read(posters[0])

    # PNG IHDR carries the dimensions at a fixed offset; a 1x1 poster is
    # invisible on the slide.
    width = int.from_bytes(poster[16:20], "big")
    height = int.from_bytes(poster[20:24], "big")
    assert width > 1 and height > 1, f"default poster is {width}x{height}"


@pytest.mark.parametrize(
    "bad_url",
    ["", "   ", "file:///C:/videos/demo.mp4", "javascript:alert(1)", "https://"],
)
def test_online_video_rejects_unusable_urls(deck: pathlib.Path, bad_url: str) -> None:
    pres = Presentation()
    pres.open(deck)
    try:
        with pytest.raises(GopptxError):
            pres.add_online_video(0, bad_url, BOUNDS)
    finally:
        pres.close()


def test_online_video_reuses_relationship_for_same_url(deck: pathlib.Path) -> None:
    pres = Presentation()
    pres.open(deck)
    try:
        pres.add_online_video(0, URL, BOUNDS)
        pres.add_online_video(0, URL, BOUNDS)
        pres.save(deck)
    finally:
        pres.close()

    with zipfile.ZipFile(deck) as archive:
        rels = archive.read("ppt/slides/_rels/slide1.xml.rels").decode()
    assert rels.count(URL) == 1
