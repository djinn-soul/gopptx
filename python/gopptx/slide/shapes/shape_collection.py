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

    def __getitem__(self, index: int | str) -> ShapeProxy:
        """Return the shape at an index, or the first shape with that name.

        A string key raises KeyError when nothing matches, so a typo is not
        mistaken for an empty slide (Issue #309).
        """
        if isinstance(index, str):
            found = self.get_by_name(index)
            if found is None:
                msg = f"no shape named {index!r} on slide {self._slide.index}"
                raise KeyError(msg)
            return found
        ids = self._shape_ids()
        if index < 0:
            index += len(ids)
        if index < 0 or index >= len(ids):
            raise IndexError("shape index out of range")
        return self._slide.shape(ids[index])

    def get_by_name(self, name: str) -> ShapeProxy | None:
        """Return the first shape with this name, or None (Issue #309).

        Shape names are not unique in OOXML, so the first match in document
        order wins, matching how PowerPoint's selection pane resolves them.
        """
        for shape in self:
            if shape.name == name:
                return shape
        return None

    def names(self) -> list[str]:
        """Return every shape name on the slide, in document order (Issue #309)."""
        return [shape.name for shape in self]

    def __iter__(self) -> Iterator[ShapeProxy]:
        """Iterate shape proxies."""
        for shape_id in self._shape_ids():
            yield self._slide.shape(shape_id)

    def iter_leaf_shapes(self) -> Iterator[ShapeProxy]:
        """Yield non-group shapes recursively in document order (Issue #435)."""
        for shape in self:
            if shape.shape_type == "grpSp":
                yield from _iter_group_leaves(shape)
            else:
                yield shape

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


def _iter_group_leaves(group: ShapeProxy) -> Iterator[ShapeProxy]:
    for child_host in group.shapes:
        child = cast("ShapeProxy", child_host)
        if child.shape_type == "grpSp":
            yield from _iter_group_leaves(child)
        else:
            yield child
