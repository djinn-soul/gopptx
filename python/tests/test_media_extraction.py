"""Listing and extracting non-image media (Issue #1049)."""

import base64
import pathlib

from gopptx import Presentation

PNG_BYTES = base64.b64decode(
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk"
    "+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
)
MOVIE_BYTES = b"fake-mp4-bytes-for-test"


def test_video_is_listed_and_extractable(tmp_path: pathlib.Path) -> None:
    """An embedded movie is listed as media and can be written out as a file."""
    output_path = tmp_path / "media_extraction.pptx"

    with Presentation.new(title="Media Extraction") as pres:
        slide = pres.slides[0]
        slide.add_video(
            MOVIE_BYTES,
            (1000000, 1000000, 4000000, 3000000),
            poster_frame=PNG_BYTES,
            mime_type="video/mp4",
        )
        pres.save(output_path)

    with Presentation(output_path) as pres:
        media = pres.list_slide_media(0)
        videos = [ref for ref in media if ref["kind"] == "video"]
        assert videos, f"expected a video in {media}"

        video = videos[0]
        assert video["part_path"].startswith("ppt/media/")
        assert video["content_type"] == "video/mp4"
        assert video["size_bytes"] == len(MOVIE_BYTES)

        assert pres.extract_media(video["part_path"]) == MOVIE_BYTES

        destination = tmp_path / "extracted.mp4"
        written = pres.save_media(video["part_path"], destination)
        assert written == len(MOVIE_BYTES)
        assert destination.read_bytes() == MOVIE_BYTES

        # The poster frame is an image, and is listed alongside the movie.
        assert any(ref["kind"] == "image" for ref in media)


def test_slide_without_media_lists_nothing(tmp_path: pathlib.Path) -> None:
    """A slide with no media reports an empty list rather than failing."""
    output_path = tmp_path / "media_none.pptx"

    with Presentation.new(title="No Media") as pres:
        pres.save(output_path)

    with Presentation(output_path) as pres:
        assert pres.list_slide_media(0) == []
