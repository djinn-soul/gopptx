"""Format proxy classes for shape fill, line, and shadow."""
# pyright: reportPrivateUsage=false, reportMissingSuperCall=false

from __future__ import annotations

from typing import TYPE_CHECKING, Protocol, cast

from ...schemas import RGBColor

if TYPE_CHECKING:
    from ...schemas import (
        FillFormat,
        LineFormat,
        PictureFill,
        ShadowFormat,
        Shape,
        ShapeUpdate,
    )


class _ShapeProto(Protocol):
    """Structural protocol for the shape object used by format proxies."""

    def shape_record(self) -> Shape: ...
    def apply_update(self, patch: ShapeUpdate) -> None: ...


class _FillColorProxy:
    """Proxy for fill fore_color."""

    def __init__(self, fill: _ShapeFillProxy) -> None:
        self._fill = fill

    @property
    def rgb(self) -> RGBColor | None:
        hex_str = self._fill.solid_color
        if not hex_str:
            return None
        return RGBColor.from_string(hex_str.lstrip("#"))

    @rgb.setter
    def rgb(self, color: RGBColor | None) -> None:
        self._fill.solid_color = None if color is None else f"#{color}"


class _LineFillProxy:
    """Proxy for line fill."""

    def __init__(self, line: _ShapeLineProxy) -> None:
        self._line = line

    def solid(self) -> None:
        """Give the line a solid fill, defaulting to black when it has no color."""
        if not self._line.color:
            self._line.color = "#000000"

    @property
    def fore_color(self) -> _LineColorProxy:
        return _LineColorProxy(self._line)


class _LineColorProxy:
    """Proxy for line color."""

    def __init__(self, line: _ShapeLineProxy) -> None:
        self._line = line

    @property
    def rgb(self) -> RGBColor | None:
        hex_str = self._line.color
        if not hex_str:
            return None
        return RGBColor.from_string(hex_str.lstrip("#"))

    @rgb.setter
    def rgb(self, color: RGBColor | None) -> None:
        self._line.color = "" if color is None else f"#{color}"


class _ShapeFillProxy:
    """Live fill proxy."""

    def __init__(self, shape: _ShapeProto) -> None:
        self._shape = shape

    def solid(self) -> None:
        """Set fill to solid."""
        if not self.solid_color:
            self.solid_color = "#000000"

    @property
    def fore_color(self) -> _FillColorProxy:
        """Return fill color proxy."""
        return _FillColorProxy(self)

    @property
    def type(self) -> str:
        """Return fill type string."""
        if self.solid_color:
            return "solid"
        if self.picture:
            return "picture"
        return "none"

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

    @property
    def picture(self) -> PictureFill | None:
        """Return the image used as this shape's fill, if any.

        Reports the relationship id, the resolved package part, the tile/stretch
        mode and any source-rectangle crop. Read-only: setting a picture fill
        goes through the picture APIs.
        """
        raw = cast("object", self._payload().get("picture"))
        if not isinstance(raw, dict):
            return None
        return cast("PictureFill", raw)

    def background(self) -> None:
        self._apply(cast("FillFormat", {"background": True}))


class _ShapeLineProxy:
    """Live line proxy."""

    def __init__(self, shape: _ShapeProto) -> None:
        self._shape = shape

    @property
    def fill(self) -> _LineFillProxy:
        """Return line fill proxy."""
        return _LineFillProxy(self)

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

    @property
    def start_arrow(self) -> str | None:
        value = self._payload().get("start_arrow")
        return str(value) if isinstance(value, str) else None

    @start_arrow.setter
    def start_arrow(self, value: str) -> None:
        self._apply({"start_arrow": value})

    @property
    def start_arrow_width(self) -> str | None:
        value = self._payload().get("start_arrow_width")
        return str(value) if isinstance(value, str) else None

    @start_arrow_width.setter
    def start_arrow_width(self, value: str) -> None:
        self._apply({"start_arrow_width": value})

    @property
    def start_arrow_length(self) -> str | None:
        value = self._payload().get("start_arrow_length")
        return str(value) if isinstance(value, str) else None

    @start_arrow_length.setter
    def start_arrow_length(self, value: str) -> None:
        self._apply({"start_arrow_length": value})

    @property
    def end_arrow(self) -> str | None:
        value = self._payload().get("end_arrow")
        return str(value) if isinstance(value, str) else None

    @end_arrow.setter
    def end_arrow(self, value: str) -> None:
        self._apply({"end_arrow": value})

    @property
    def end_arrow_width(self) -> str | None:
        value = self._payload().get("end_arrow_width")
        return str(value) if isinstance(value, str) else None

    @end_arrow_width.setter
    def end_arrow_width(self, value: str) -> None:
        self._apply({"end_arrow_width": value})

    @property
    def end_arrow_length(self) -> str | None:
        value = self._payload().get("end_arrow_length")
        return str(value) if isinstance(value, str) else None

    @end_arrow_length.setter
    def end_arrow_length(self, value: str) -> None:
        self._apply({"end_arrow_length": value})


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
        # An explicit effect already overrides the shape style, so carrying a
        # stale ``inherit`` alongside it would only describe the same state twice.
        payload.pop("inherit", None)
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
    def inherit(self) -> bool:
        """Return whether the shape takes its effects from the shape style.

        True when the shape carries no effect list of its own, matching
        python-pptx's ``ShadowFormat.inherit`` (upstream issue #130).
        """
        payload = self._payload()
        value = payload.get("inherit")
        if isinstance(value, bool):
            return value
        return not any(
            key in payload for key in ("color", "blur_emu", "distance_emu", "angle_deg")
        )

    @inherit.setter
    def inherit(self, value: bool) -> None:
        # Sent alone: inherit describes the whole effect list, so any previously
        # set explicit attribute must not travel with it.
        self._shape.apply_update(cast("ShapeUpdate", {"shadow": {"inherit": value}}))


ShapeFillProxy = _ShapeFillProxy
ShapeLineProxy = _ShapeLineProxy
ShapeShadowProxy = _ShapeShadowProxy
