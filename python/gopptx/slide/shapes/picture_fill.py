"""Shared payload builder for picture fills (Issue #234).

Both a shape fill and a table-cell fill go through the same ``set_picture_fill``
op, so the payload they send is built in one place.
"""

from __future__ import annotations

import base64
import os

_CROP_EDGES = ("left", "top", "right", "bottom")


def picture_fill_payload(
    image: str | bytes | os.PathLike[str],
    *,
    tile: bool,
    crop: tuple[float, float, float, float] | None,
) -> dict[str, object]:
    """Build the set_picture_fill payload shared by shape and table-cell fills."""
    payload: dict[str, object] = {"mode": "tile" if tile else "stretch"}
    if isinstance(image, bytes):
        payload["image_data"] = base64.b64encode(image).decode("ascii")
    else:
        payload["image_path"] = os.fspath(image)
    if crop is not None:
        if len(crop) != len(_CROP_EDGES):
            raise ValueError("crop must be a (left, top, right, bottom) tuple")
        payload["crop"] = dict(zip(_CROP_EDGES, (float(v) for v in crop), strict=True))
    return payload
