"""Image part proxy for picture shapes for python-pptx parity (Issue #1084)."""

from __future__ import annotations

import base64
import hashlib
import os
import pathlib
from typing import TYPE_CHECKING, BinaryIO, Protocol

from typing_extensions import override

from ...presentation.shapes.image_inspection import infer_image_format

if TYPE_CHECKING:
    from collections.abc import Mapping


class _ImagePresentationProtocol(Protocol):
    def get_image_metadata(
        self, slide_index: int, shape_id: int
    ) -> Mapping[str, object]:
        """Return image metadata."""
        ...

    def swap_image_by_rel_id(
        self, slide_index: int, rel_id: str, data: bytes, img_format: str
    ) -> None:
        """Replace the image behind a relationship ID."""
        ...


class _ImageSlideProtocol(Protocol):
    @property
    def presentation(self) -> _ImagePresentationProtocol:
        """Return the owning presentation."""
        ...

    @property
    def index(self) -> int:
        """Return the slide index."""
        ...


class _ImageShapeProtocol(Protocol):
    @property
    def slide(self) -> _ImageSlideProtocol:
        """Return the owning slide."""
        ...

    @property
    def id(self) -> int:
        """Return the shape ID."""
        ...


class ImagePartProxy:
    """Proxy representing the underlying image asset of a picture shape (Issue #1084)."""

    def __init__(self, shape: _ImageShapeProtocol) -> None:
        """Create an image part proxy for a picture shape."""
        super().__init__()
        self._shape = shape
        self._meta_cache: dict[str, object] | None = None

    def _get_meta(self) -> dict[str, object]:
        if self._meta_cache is None:
            slide = self._shape.slide
            pres = slide.presentation
            slide_index = slide.index
            meta = pres.get_image_metadata(slide_index, self._shape.id)
            self._meta_cache = dict(meta)
        return self._meta_cache

    @property
    def blob(self) -> bytes:
        """Return the raw image byte content."""
        meta = self._get_meta()
        data = meta["data"]
        if isinstance(data, str):
            return base64.b64decode(data, validate=True)
        if isinstance(data, (bytes, bytearray)):
            return bytes(data)
        raise TypeError(
            f"image data must be bytes or base64 text, got {type(data).__name__}"
        )

    @property
    def content_type(self) -> str:
        """Return the MIME content type (e.g. 'image/tiff', 'image/png')."""
        return str(self._get_meta()["content_type"])

    @property
    def ext(self) -> str:
        """Return the image format extension (e.g. 'tif', 'png', 'jpeg')."""
        return str(self._get_meta()["format"]).lower()

    @property
    def filename(self) -> str:
        """Return the image filename."""
        return f"image.{self.ext}"

    @property
    def sha1(self) -> str:
        """Return SHA-1 hex digest of the image blob."""
        # Non-security content fingerprint retained for python-pptx API parity.
        return hashlib.sha1(self.blob, usedforsecurity=False).hexdigest()  # nosemgrep

    def save(self, file_or_path: str | os.PathLike[str] | BinaryIO) -> None:
        """Save the image blob to a file path or open binary stream.

        Args:
            file_or_path: File path string, Path object, or binary stream with write().
        """
        blob = self.blob
        if not isinstance(file_or_path, (str, os.PathLike)):
            file_or_path.write(blob)
            return
        path = pathlib.Path(file_or_path)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(blob)

    def replace(self, image_path_or_bytes: str | bytes | os.PathLike[str]) -> None:
        """Swap the image asset behind this picture shape (Issue #1142).

        Args:
            image_path_or_bytes: Raw image bytes, or a path to an image file.

        Raises:
            ValueError: The shape's image relationship could not be resolved.
        """
        if isinstance(image_path_or_bytes, (bytes, bytearray)):
            data = bytes(image_path_or_bytes)
        else:
            data = pathlib.Path(os.fspath(image_path_or_bytes)).read_bytes()
        # Sniff the bytes rather than trust a file extension: the bridge keys
        # the content type off this string, and a mislabelled part is what
        # PowerPoint refuses to open.
        img_format = infer_image_format(data)

        rel_id = str(self._get_meta().get("rel_id", ""))
        if not rel_id:
            raise ValueError("shape has no resolvable image relationship to replace")

        slide = self._shape.slide
        slide.presentation.swap_image_by_rel_id(slide.index, rel_id, data, img_format)
        self._meta_cache = None

    @override
    def __repr__(self) -> str:
        """Return a concise image-part description."""
        return (
            f"<ImagePartProxy format={self.ext!r} "
            f"content_type={self.content_type!r} bytes={len(self.blob)}>"
        )
