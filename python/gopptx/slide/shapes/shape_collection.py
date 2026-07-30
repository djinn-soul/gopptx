"""python-pptx-style slide shapes collection facade."""
# pyright: reportPrivateUsage=false

from __future__ import annotations

from typing import TYPE_CHECKING, cast

if TYPE_CHECKING:
    import os
    from collections.abc import Iterator

    from ...constants import ShapeType
    from ...schemas import ShapeProps
    from ...shapes import ShapeBuilder
    from .shape_proxy import ShapeProxy, _ShapeProxySlideProto


class ShapeCollection:
    """python-pptx-style slide shapes collection."""

    def __init__(self, slide: _ShapeProxySlideProto) -> None:  # pyright: ignore[reportMissingSuperCall]
        """Create a shape collection for a slide."""
        self._slide = slide

    def _shape_ids(self) -> list[int]:
        out: list[int] = []
        for shape in self._slide.list_shapes():
            shape_id = shape.get("ID", shape.get("id"))
            if shape_id is None:
                continue
            out.append(int(str(shape_id)))
        return out

    def __len__(self) -> int:
        """Return shape count."""
        return len(self._shape_ids())

    def __getitem__(self, index: int) -> ShapeProxy:
        """Return shape proxy at index."""
        ids = self._shape_ids()
        if index < 0:
            index += len(ids)
        if index < 0 or index >= len(ids):
            raise IndexError("shape index out of range")
        return self._slide.shape(ids[index])

    def __iter__(self) -> Iterator[ShapeProxy]:
        """Iterate shape proxies."""
        for shape_id in self._shape_ids():
            yield self._slide.shape(shape_id)

    def add(self, builder: ShapeBuilder) -> ShapeProxy:
        """Add a shape from a ShapeBuilder and return its proxy."""
        shape_id = self._slide.add_shape(
            builder.shape_type,
            builder.bounds,
            **cast("dict[str, ShapeProps]", builder.to_kwargs()),
        )
        return self._slide.shape(shape_id)

    def add_textbox(
        self,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        text: str = "",
        **kwargs: str | ShapeProps,
    ) -> ShapeProxy:
        """Add a textbox to slide and return shape proxy."""
        shape_id = self._slide.add_textbox(
            left, top, width, height, text=text, **kwargs
        )
        return self._slide.shape(shape_id)

    def add_shape(
        self,
        shape_type: ShapeType | object,
        left: float,
        top: float,
        width: float,
        height: float,
        **kwargs: str | ShapeProps,
    ) -> ShapeProxy:
        """Add an autoshape to slide and return shape proxy."""
        bounds = (left, top, width, height)
        shape_id = self._slide.add_shape(
            cast("ShapeType", shape_type), bounds, **kwargs
        )
        return self._slide.shape(shape_id)

    def add_picture(
        self,
        image_path_or_file: str | bytes | os.PathLike[str] | None,
        left: float,
        top: float,
        width: float,
        height: float,
        **kwargs: object,
    ) -> ShapeProxy:
        """Add a picture to slide and return shape proxy."""
        shape_id = self._slide.add_picture(
            image_path_or_file, left, top, width, height, **kwargs
        )
        return self._slide.shape(shape_id)

    def remove(self, shape: object) -> None:
        """Remove a shape proxy or shape ID from the collection (Issue #41)."""
        shape_id = getattr(shape, "id", shape)
        if isinstance(shape_id, int):
            remove_func = getattr(self._slide, "remove_shape", None)
            if callable(remove_func):
                remove_func(shape_id)
            else:
                self._slide.presentation.remove_shape(self._slide.index, shape_id)
        else:
            raise TypeError("shape must be a ShapeProxy or integer shape ID")

    def __delitem__(self, index: int) -> None:
        """Delete shape at index (Issue #41)."""
        proxy = self[index]
        proxy.delete()

    def clear(self) -> int:
        """Remove every shape on the slide and return how many were removed.

        Ids are snapshotted first: removing a shape renumbers nothing, but it
        does invalidate the slide's shape cache, so iterating live would skip
        entries (Issue #96).
        """
        shape_ids = self._shape_ids()
        for shape_id in shape_ids:
            self.remove(shape_id)
        return len(shape_ids)
