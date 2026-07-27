"""Image format and native-size inspection for picture insertion."""

from __future__ import annotations

import os
from io import BytesIO
from pathlib import Path
from typing import TYPE_CHECKING, cast

# Pillow is a runtime dependency. The import is guarded so a wheel installed
# without it fails with an actionable message instead of an opaque
# ModuleNotFoundError raised from an unrelated module.
try:
    from PIL import Image
except ImportError as exc:  # pragma: no cover - depends on install shape
    raise ImportError(
        "Pillow is required to inspect image sources: pip install pillow"
    ) from exc

if TYPE_CHECKING:
    from os import PathLike

    from PIL.ImageFile import ImageFile

_EMU_PER_INCH = 914_400
_DEFAULT_DPI = 72.0
_DPI_COMPONENTS = 2

# An XML prolog, a DOCTYPE and a comment can all precede the root element, so
# the window is wide enough to reach <svg but still bounded.
_SVG_SNIFF_BYTES = 1024


def is_svg_source(source: str | bytes | PathLike[str]) -> bool:
    """Check whether the image source is SVG vector graphics."""
    if isinstance(source, (str, os.PathLike)):
        return os.fspath(source).lower().endswith(".svg")
    # An XML prolog alone is not enough: any XML-ish payload would otherwise
    # take the SVG fallback instead of failing to decode. Require an <svg root.
    return b"<svg" in source.lstrip()[:_SVG_SNIFF_BYTES].lower()


def infer_image_format(source: bytes) -> str:
    """Return the decoded image format used by the Go bridge."""
    if is_svg_source(source):
        return "svg"
    with Image.open(BytesIO(source)) as image:
        return _normalized_format(image)


def resolve_picture_source(
    source: str | bytes | PathLike[str] | None,
    options: dict[str, object],
) -> str | bytes:
    """Resolve picture input using explicit source, path, then data."""
    if source is not None:
        resolved = os.fspath(source) if isinstance(source, os.PathLike) else source
        if not resolved:
            raise ValueError("picture source must not be empty")
        return resolved
    path = options.get("path")
    if isinstance(path, (str, os.PathLike)):
        resolved_path: str = os.fspath(path)
        if not resolved_path:
            raise ValueError("picture path must not be empty")
        return resolved_path
    data = options.get("data")
    if isinstance(data, bytes):
        if not data:
            raise ValueError("picture data must not be empty")
        return data
    raise ValueError("picture source is required")


def picture_bounds(
    source: str | bytes | PathLike[str],
    left: float,
    top: float,
    width: float,
    height: float,
) -> tuple[float, float, float, float]:
    """Resolve omitted picture dimensions from native pixels and DPI."""
    if width < 0 or height < 0:
        raise ValueError("picture width and height must not be negative")
    if width > 0 and height > 0:
        return (left, top, width, height)

    if is_svg_source(source):
        native_width, native_height = 2_540_000, 2_540_000
    else:
        if not isinstance(source, bytes) and not Path(source).exists():
            raise FileNotFoundError(f"picture source not found: {os.fspath(source)}")
        with Image.open(
            BytesIO(source) if isinstance(source, bytes) else source
        ) as image:
            native_width, native_height = _native_size_emu(image)

    if width == 0 and height == 0:
        width, height = float(native_width), float(native_height)
    elif width == 0:
        width = round(height * native_width / native_height)
    else:
        height = round(width * native_height / native_width)
    return (left, top, width, height)


def _normalized_format(image: ImageFile) -> str:
    image_format = image.format
    if not image_format:
        raise ValueError("image format could not be determined")
    return image_format.strip().lower()


def _dpi_pair(image: ImageFile) -> tuple[float, float]:
    """Return the image's horizontal and vertical DPI.

    ``image.info`` is an untyped mapping, and a malformed file can put anything
    under "dpi", so each component is checked before use and falls back to the
    72 DPI default rather than propagating an unusable value.
    """
    raw: object = image.info.get("dpi")
    if not isinstance(raw, tuple):
        return (_DEFAULT_DPI, _DEFAULT_DPI)

    values = cast("tuple[object, ...]", raw)
    if len(values) < _DPI_COMPONENTS:
        return (_DEFAULT_DPI, _DEFAULT_DPI)

    components = [
        _positive_dpi(float(value)) if isinstance(value, (int, float)) else _DEFAULT_DPI
        for value in values[:_DPI_COMPONENTS]
    ]
    return (components[0], components[1])


def _native_size_emu(image: ImageFile) -> tuple[int, int]:
    horizontal_dpi, vertical_dpi = _dpi_pair(image)
    width_px, height_px = image.size
    return (
        round(width_px * _EMU_PER_INCH / horizontal_dpi),
        round(height_px * _EMU_PER_INCH / vertical_dpi),
    )


def _positive_dpi(value: object) -> float:
    if isinstance(value, (int, float)) and value > 0:
        return float(value)
    return _DEFAULT_DPI
