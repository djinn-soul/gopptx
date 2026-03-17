"""Fluent builder for adding shapes to slides."""

from __future__ import annotations

from typing_extensions import override

from gopptx.constants import (
    SHAPE_CLOUD,
    SHAPE_DIAMOND,
    SHAPE_ELLIPSE,
    SHAPE_HEART,
    SHAPE_HEXAGON,
    SHAPE_PARALLELOGRAM,
    SHAPE_PENTAGON,
    SHAPE_RECTANGLE,
    SHAPE_RIGHT_TRIANGLE,
    SHAPE_ROUNDED_RECTANGLE,
    SHAPE_STAR_5,
    SHAPE_STAR_6,
    SHAPE_TRIANGLE,
)


class ShapeBuilder:
    """Fluent builder for shape creation parameters.

    Construct a builder with a factory class method, chain formatting calls,
    then pass the builder to ``slide.shapes.add()``::

        from gopptx.shapes import ShapeBuilder

        proxy = slide.shapes.add(
            ShapeBuilder.rectangle(1.0, 1.0, 4.0, 2.0)
            .with_text("Hello")
            .with_fill("4472C4")
            .with_line("FFFFFF", width_pt=1.0)
        )

    All positions and sizes are in **inches**. Line/shadow sizes use EMU
    (1 point = 12700 EMU).
    """

    _EMU_PER_INCH: int = 914_400
    _EMU_PER_PT: float = 12_700.0

    def __init__(
        self,
        shape_type: str,
        x: float,
        y: float,
        w: float,
        h: float,
    ) -> None:
        """Create a builder for *shape_type* at position (*x*, *y*) with size (*w* x *h*) in inches."""
        super().__init__()
        self._shape_type = shape_type
        self._x = x
        self._y = y
        self._w = w
        self._h = h
        self._text: str | None = None
        self._properties: dict[str, object] = {}

    # ── Factory class methods ────────────────────────────────────────────────

    @classmethod
    def of(
        cls, shape_type: str, x: float, y: float, w: float, h: float
    ) -> ShapeBuilder:
        """Generic factory — use any OOXML preset geometry token."""
        return cls(shape_type, x, y, w, h)

    @classmethod
    def rectangle(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a rectangle builder."""
        return cls(SHAPE_RECTANGLE, x, y, w, h)

    @classmethod
    def rounded_rectangle(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a rounded-rectangle builder."""
        return cls(SHAPE_ROUNDED_RECTANGLE, x, y, w, h)

    @classmethod
    def ellipse(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create an ellipse (or circle) builder."""
        return cls(SHAPE_ELLIPSE, x, y, w, h)

    @classmethod
    def triangle(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create an isosceles-triangle builder."""
        return cls(SHAPE_TRIANGLE, x, y, w, h)

    @classmethod
    def right_triangle(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a right-triangle builder."""
        return cls(SHAPE_RIGHT_TRIANGLE, x, y, w, h)

    @classmethod
    def diamond(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a diamond builder."""
        return cls(SHAPE_DIAMOND, x, y, w, h)

    @classmethod
    def pentagon(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a pentagon builder."""
        return cls(SHAPE_PENTAGON, x, y, w, h)

    @classmethod
    def hexagon(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a hexagon builder."""
        return cls(SHAPE_HEXAGON, x, y, w, h)

    @classmethod
    def parallelogram(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a parallelogram builder."""
        return cls(SHAPE_PARALLELOGRAM, x, y, w, h)

    @classmethod
    def cloud(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a cloud builder."""
        return cls(SHAPE_CLOUD, x, y, w, h)

    @classmethod
    def heart(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a heart builder."""
        return cls(SHAPE_HEART, x, y, w, h)

    @classmethod
    def star5(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a 5-point star builder."""
        return cls(SHAPE_STAR_5, x, y, w, h)

    @classmethod
    def star6(cls, x: float, y: float, w: float, h: float) -> ShapeBuilder:
        """Create a 6-point star builder."""
        return cls(SHAPE_STAR_6, x, y, w, h)

    # ── Text ────────────────────────────────────────────────────────────────

    def with_text(self, text: str) -> ShapeBuilder:
        """Set the shape's text content."""
        self._text = text
        return self

    # ── Fill ─────────────────────────────────────────────────────────────────

    def with_fill(self, color: str) -> ShapeBuilder:
        """Apply a solid fill. *color* is a 6-digit hex string without ``#``."""
        self._properties["fill"] = {"solid": color}
        return self

    def with_no_fill(self) -> ShapeBuilder:
        """Make the shape background transparent."""
        self._properties["fill"] = {"background": True}
        return self

    # ── Line / border ────────────────────────────────────────────────────────

    def with_line(
        self,
        color: str,
        *,
        width_emu: int | None = None,
        dash_style: str | None = None,
    ) -> ShapeBuilder:
        """Style the shape border.

        Args:
            color: 6-digit hex color without ``#``.
            width_emu: Line width in EMU (1 pt = 12700 EMU).
            dash_style: OOXML dash token such as ``"dash"``, ``"dot"``,
                ``"dashDot"``, ``"lgDash"``, ``"lgDashDot"``, ``"solid"``.
        """
        line: dict[str, object] = {"color": color}
        if width_emu is not None:
            line["width_emu"] = width_emu
        if dash_style is not None:
            line["dash_style"] = dash_style
        self._properties["line"] = line
        return self

    def with_no_line(self) -> ShapeBuilder:
        """Remove the shape border (zero-width transparent line)."""
        self._properties["line"] = {"color": "FFFFFF", "width_emu": 0}
        return self

    # ── Shadow ───────────────────────────────────────────────────────────────

    def with_shadow(
        self,
        color: str = "000000",
        *,
        blur_emu: int | None = None,
        distance_emu: int | None = None,
        angle_deg: float = 45.0,
    ) -> ShapeBuilder:
        """Apply a drop shadow.

        Args:
            color: 6-digit hex shadow color.
            blur_emu: Blur radius in EMU (1 pt = 12700 EMU).
            distance_emu: Shadow offset in EMU.
            angle_deg: Shadow angle in degrees (default 45).
        """
        shadow: dict[str, object] = {"color": color, "angle_deg": angle_deg}
        if blur_emu is not None:
            shadow["blur_emu"] = blur_emu
        if distance_emu is not None:
            shadow["distance_emu"] = distance_emu
        self._properties["shadow"] = shadow
        return self

    # ── Geometry transforms ──────────────────────────────────────────────────

    def with_rotation(self, degrees: float) -> ShapeBuilder:
        """Rotate the shape by *degrees* clockwise."""
        self._properties["rotation"] = degrees
        return self

    def flip_horizontal(self) -> ShapeBuilder:
        """Mirror the shape horizontally."""
        self._properties["flip_h"] = True
        return self

    def flip_vertical(self) -> ShapeBuilder:
        """Mirror the shape vertically."""
        self._properties["flip_v"] = True
        return self

    # ── Internal helpers ─────────────────────────────────────────────────────

    @property
    def shape_type(self) -> str:
        """The OOXML preset geometry token."""
        return self._shape_type

    @property
    def bounds(self) -> tuple[float, float, float, float]:
        """(x, y, w, h) in EMU."""
        emu = self._EMU_PER_INCH
        return (
            self._x * emu,
            self._y * emu,
            self._w * emu,
            self._h * emu,
        )

    def to_kwargs(self) -> dict[str, object]:
        """Return the keyword arguments for ``slide.add_shape()``."""
        kwargs: dict[str, object] = {}
        if self._text is not None:
            kwargs["text"] = self._text
        if self._properties:
            kwargs["properties"] = dict(self._properties)
        return kwargs

    @override
    def __repr__(self) -> str:
        """Return a diagnostic repr."""
        return f"<ShapeBuilder type={self._shape_type!r} bounds=({self._x}, {self._y}, {self._w}, {self._h})>"
