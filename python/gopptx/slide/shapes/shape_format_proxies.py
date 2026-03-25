"""Format proxy classes for shape fill, line, and shadow."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

if TYPE_CHECKING:
    from ...schemas import (
        FillFormat,
        LineFormat,
        ShadowFormat,
        Shape,
        ShapeUpdate,
    )


class _ShapeProto(Protocol):
    """Structural protocol for the shape object used by format proxies."""

    def shape_record(self) -> Shape: ...
    def apply_update(self, patch: ShapeUpdate) -> None: ...


class _ShapeFillProxy:
    """Live fill proxy."""

    def __init__(self, shape: _ShapeProto) -> None:
        self._shape = shape

    def _payload(self) -> FillFormat:
        record = self._shape.shape_record()
        raw = cast("object", record.get("fill", record.get("Fill", {})))
        return cast("FillFormat", raw if raw is not None else {})

    def _apply(self, payload: FillFormat) -> None:
        self._shape.apply_update(cast("ShapeUpdate", {"fill": payload}))

    @property
    def solid_color(self) -> str | None:
        payload = self._payload()
        value = payload.get("solid")
        return str(value) if isinstance(value, str) else None

    @solid_color.setter
    def solid_color(self, value: str | None) -> None:
        if value is None:
            self._apply(cast("FillFormat", {"background": True}))
            return
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.pop("background", None)
        payload["solid"] = value
        self._apply(cast("FillFormat", payload))

    @property
    def transparency(self) -> float | None:
        value = self._payload().get("transparency")
        return float(value) if isinstance(value, int | float) else None

    @transparency.setter
    def transparency(self, value: float | None) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        if value is None:
            payload.pop("transparency", None)
            self._apply(cast("FillFormat", payload))
            return
        if value < 0.0 or value > 1.0:
            raise ValueError("fill.transparency must be between 0.0 and 1.0")
        if not isinstance(payload.get("solid"), str):
            raise ValueError("fill.transparency requires a solid fill color")
        payload["transparency"] = float(value)
        self._apply(cast("FillFormat", payload))

    def background(self) -> None:
        self._apply(cast("FillFormat", {"background": True}))


class _ShapeLineProxy:
    """Live line proxy."""

    def __init__(self, shape: _ShapeProto) -> None:
        self._shape = shape

    def _payload(self) -> LineFormat:
        record = self._shape.shape_record()
        raw = cast("object", record.get("line", record.get("Line", {})))
        return cast("LineFormat", raw if raw is not None else {})

    def _apply(self, patch: dict[str, object]) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.update(patch)
        self._shape.apply_update(cast("ShapeUpdate", {"line": payload}))

    @property
    def color(self) -> str | None:
        value = self._payload().get("color")
        return str(value) if isinstance(value, str) else None

    @color.setter
    def color(self, value: str) -> None:
        self._apply({"color": value})

    @property
    def width(self) -> int | None:
        value = self._payload().get("width_emu")
        return int(value) if isinstance(value, int) else None

    @width.setter
    def width(self, value: int) -> None:
        self._apply({"width_emu": value})

    @property
    def dash_style(self) -> str | None:
        value = self._payload().get("dash_style")
        return str(value) if isinstance(value, str) else None

    @dash_style.setter
    def dash_style(self, value: str) -> None:
        self._apply({"dash_style": value})


class _ShapeShadowProxy:
    """Live shadow proxy."""

    def __init__(self, shape: _ShapeProto) -> None:
        self._shape = shape

    def _payload(self) -> ShadowFormat:
        record = self._shape.shape_record()
        raw = cast("object", record.get("shadow", record.get("Shadow", {})))
        return cast("ShadowFormat", raw if raw is not None else {})

    def _apply(self, patch: dict[str, object]) -> None:
        payload = dict(cast("dict[str, object]", self._payload()))
        payload.update(patch)
        self._shape.apply_update(cast("ShapeUpdate", {"shadow": payload}))

    @property
    def color(self) -> str | None:
        value = self._payload().get("color")
        return str(value) if isinstance(value, str) else None

    @color.setter
    def color(self, value: str) -> None:
        self._apply({"color": value})

    @property
    def blur_radius(self) -> int | None:
        value = self._payload().get("blur_emu")
        return int(value) if isinstance(value, int) else None

    @blur_radius.setter
    def blur_radius(self, value: int) -> None:
        self._apply({"blur_emu": value})

    @property
    def distance(self) -> int | None:
        value = self._payload().get("distance_emu")
        return int(value) if isinstance(value, int) else None

    @distance.setter
    def distance(self, value: int) -> None:
        self._apply({"distance_emu": value})

    @property
    def angle(self) -> float | None:
        value = self._payload().get("angle_deg")
        return float(value) if isinstance(value, int | float) else None

    @angle.setter
    def angle(self, value: float) -> None:
        self._apply({"angle_deg": value})

    @property
    def inherit(self) -> bool | None:
        value = self._payload().get("inherit")
        return bool(value) if isinstance(value, bool) else None

    @inherit.setter
    def inherit(self, value: bool) -> None:
        self._apply({"inherit": value})


ShapeFillProxy = _ShapeFillProxy
ShapeLineProxy = _ShapeLineProxy
ShapeShadowProxy = _ShapeShadowProxy
