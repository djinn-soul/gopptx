"""Minimal freeform builder compatible with python-pptx style usage."""

from __future__ import annotations

from typing import Protocol

MIN_POINTS_FOR_FREEFORM = 2


class _FreeformCommitter(Protocol):
    def commit_freeform(
        self,
        slide_index: int,
        points: list[tuple[float, float]],
        *,
        close: bool,
        options: dict[str, object] | None = None,
    ) -> int: ...


class FreeformBuilder:
    """Line-segment freeform builder for a single slide."""

    def __init__(
        self,
        committer: _FreeformCommitter,
        slide_index: int,
        *,
        start_x: float = 0,
        start_y: float = 0,
        scale: tuple[float, float] | float = 1.0,
    ) -> None:
        """Initialize builder state for one slide."""
        self._committer = committer
        self._slide_index = slide_index
        if isinstance(scale, tuple):
            self._scale_x, self._scale_y = scale
        else:
            self._scale_x = scale
            self._scale_y = scale
        self._points: list[tuple[float, float]] = [(start_x, start_y)]

    def add_line_to(self, x: float, y: float) -> FreeformBuilder:
        """Append a line segment endpoint."""
        self._points.append((x, y))
        return self

    def add_line_segments(self, points: list[tuple[float, float]]) -> FreeformBuilder:
        """Append multiple line segment endpoints."""
        self._points.extend(points)
        return self

    def move_to(self, x: float, y: float) -> FreeformBuilder:
        """Set the initial point before any segments are added."""
        if len(self._points) > 1:
            raise ValueError("move_to() is only allowed before line segments are added")
        self._points[0] = (x, y)
        return self

    def convert_to_shape(self, *, close: bool = False, **kwargs: object) -> int:
        """Create the freeform shape and return its shape ID."""
        if len(self._points) < MIN_POINTS_FOR_FREEFORM:
            raise ValueError(
                "freeform requires at least one line segment before convert_to_shape()"
            )
        scaled_points = [
            (x * self._scale_x, y * self._scale_y) for x, y in self._points
        ]
        return self._committer.commit_freeform(
            self._slide_index,
            scaled_points,
            close=close,
            options=kwargs or None,
        )
