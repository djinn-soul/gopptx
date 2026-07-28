"""Shape-media operations for the presentation facade."""

from __future__ import annotations

import base64
from pathlib import Path
from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import get_required_int
from .image_inspection import infer_image_format, picture_bounds, resolve_picture_source
from .image_options import reject_unknown_image_options
from .shape_media_av_mixin import PresentationShapeAVMixin
from .shape_payload_mixin import PresentationShapePayloadMixin

if TYPE_CHECKING:
    from os import PathLike

    from ...schemas import ImageMetadata, SlideImageRef, SlideMediaRef


class PresentationShapeMediaMixin(
    PresentationShapeAVMixin, PresentationShapePayloadMixin
):
    """Methods that add and inspect image/video/audio/OLE shapes."""

    def add_image(
        self,
        slide_index: int,
        source: str | bytes | None = None,
        bounds: tuple[float, float, float, float] = (0, 0, 0, 0),
        **kwargs: object,
    ) -> int:
        """Add an image to a slide and return the created shape ID.

        Raises:
            TypeError: If a keyword option is not recognized.
        """
        reject_unknown_image_options(kwargs)
        payload = self._init_bounds_payload(slide_index, bounds)
        self._resolve_image_source(payload, source, kwargs)
        self._resolve_image_options(payload, kwargs)

        result = self.execute(ops.OP_ADD_IMAGE, payload)
        return get_required_int(result, "shape_id")

    def add_picture(
        self,
        slide_index: int,
        source: str | bytes | None = None,
        left: float = 0,
        top: float = 0,
        width: float = 0,
        height: float = 0,
        **kwargs: object,
    ) -> int:
        """Add a picture shape to a slide (python-pptx compatible API).

        Supports optional ``description``, ``alt_text``, and ``title`` parameters.
        """
        effective_source = resolve_picture_source(source, kwargs)
        bounds = picture_bounds(effective_source, left, top, width, height)
        self._validate_picture_metadata(kwargs)
        options = {
            key: value for key, value in kwargs.items() if key not in {"path", "data"}
        }
        return self.add_image(slide_index, effective_source, bounds=bounds, **options)

    @staticmethod
    def _validate_picture_metadata(kwargs: dict[str, object]) -> None:
        """Validate optional picture metadata supplied through keyword options."""
        for key in ("description", "alt_text", "title"):
            value = kwargs.get(key)
            if value is not None and not isinstance(value, str):
                raise TypeError(f"{key} must be a string")

    def _resolve_image_source(
        self,
        payload: dict[str, object],
        source: str | bytes | None,
        kwargs: dict[str, object],
    ) -> None:
        path = kwargs.get("path")
        data = kwargs.get("data")
        if source:
            self._set_source_payload(payload, source)
            if isinstance(source, bytes):
                self._set_image_format(payload, source, kwargs)
        elif isinstance(path, str):
            self._set_source_payload(payload, path)
        elif isinstance(data, bytes):
            self._set_source_payload(payload, data)
            self._set_image_format(payload, data, kwargs)

    @staticmethod
    def _set_image_format(
        payload: dict[str, object],
        data: bytes,
        kwargs: dict[str, object],
    ) -> None:
        fmt = kwargs.get("image_format") or kwargs.get("img_format")
        payload["format"] = fmt if isinstance(fmt, str) else infer_image_format(data)

    @staticmethod
    def _resolve_image_options(
        payload: dict[str, object], kwargs: dict[str, object]
    ) -> None:
        options: dict[str, object] = {}
        crop = kwargs.get("crop")
        if isinstance(crop, dict):
            options["crop"] = cast("dict[str, object]", crop)

        rot = kwargs.get("rotation")
        if isinstance(rot, (int, float)):
            options["rotation"] = rot

        for key in ("flip_h", "flip_v"):
            val = kwargs.get(key)
            if isinstance(val, bool):
                options[key] = val

        descr = kwargs.get("description") or kwargs.get("alt_text")
        if isinstance(descr, str):
            options["description"] = descr
            options["alt_text"] = descr

        title = kwargs.get("title")
        if isinstance(title, str):
            options["title"] = title

        if options:
            payload["options"] = options

    def get_image_metadata(self, slide_index: int, shape_id: int) -> ImageMetadata:
        """Get dimensions and format metadata for an image shape."""
        result = self.execute(
            ops.OP_GET_IMAGE_METADATA,
            {"slide_index": slide_index, "shape_id": shape_id},
        )
        return cast("ImageMetadata", result)

    def list_slide_images(self, slide_index: int) -> list[SlideImageRef]:
        """List all images embedded in a slide.

        Args:
            slide_index: Zero-based index of the slide.

        Returns:
            List of SlideImageRef dicts with keys: index, rel_id, target.
        """
        result = self.execute(ops.OP_LIST_SLIDE_IMAGES, {"slide_index": slide_index})
        return cast("list[SlideImageRef]", result.get("images", []))

    def list_slide_media(self, slide_index: int) -> list[SlideMediaRef]:
        """List every media relationship on a slide: images, sounds and movies.

        Images already had ``list_slide_images``; an embedded movie or sound
        could only be found by walking relationships by hand (issue #1049).

        Args:
            slide_index: Zero-based index of the slide.

        Returns:
            List of SlideMediaRef dicts: index, rel_id, kind, target, part_path,
            content_type, size_bytes and external.
        """
        result = self.execute(ops.OP_LIST_SLIDE_MEDIA, {"slide_index": slide_index})
        return cast("list[SlideMediaRef]", result.get("media", []))

    def extract_media(self, part_path: str) -> bytes:
        """Return the bytes of one media part.

        Args:
            part_path: Package part path, as reported by list_slide_media
                (for example ``ppt/media/media1.mp4``).
        """
        result = self.execute(ops.OP_EXTRACT_MEDIA, {"part_path": part_path})
        return base64.b64decode(cast("str", result.get("data", "")))

    def save_media(self, part_path: str, destination: str | PathLike[str]) -> int:
        """Write one media part out to a file, and return the bytes written."""
        data = self.extract_media(part_path)
        Path(destination).write_bytes(data)
        return len(data)

    def swap_image_by_index(
        self,
        slide_index: int,
        image_index: int,
        data: bytes,
        img_format: str,
    ) -> None:
        """Replace an image at a given position within a slide.

        Args:
            slide_index: Zero-based slide index.
            image_index: Zero-based position of the image within the slide's
                image list (as returned by list_slide_images).
            data: Raw image bytes.
            img_format: Image format string (e.g. 'png', 'jpeg').
        """
        self.execute(
            ops.OP_SWAP_IMAGE_BY_INDEX,
            {
                "slide_index": slide_index,
                "image_index": image_index,
                "data": base64.b64encode(data).decode(),
                "format": img_format,
            },
        )

    def swap_image_by_rel_id(
        self,
        slide_index: int,
        rel_id: str,
        data: bytes,
        img_format: str,
    ) -> None:
        """Replace an image identified by its relationship ID.

        Args:
            slide_index: Zero-based slide index.
            rel_id: Relationship ID of the image to replace (e.g. 'rId3').
            data: Raw image bytes.
            img_format: Image format string (e.g. 'png', 'jpeg').
        """
        self.execute(
            ops.OP_SWAP_IMAGE_BY_REL_ID,
            {
                "slide_index": slide_index,
                "rel_id": rel_id,
                "data": base64.b64encode(data).decode(),
                "format": img_format,
            },
        )
