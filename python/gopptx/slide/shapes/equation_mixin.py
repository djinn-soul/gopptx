"""Slide mixin for native PowerPoint equations (Issue #126)."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops

if TYPE_CHECKING:
    from ..contracts import SlidePresentationProtocol


class SlideEquationMixin:
    """Mixin adding equation support to Slide objects."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

        def _invalidate_shape_cache_if_present(self) -> None: ...

        def _invalidate_text_state_cache_if_present(self) -> None: ...

    def add_equation(
        self,
        latex: str,
        left: float,
        top: float,
        width: float,
        height: float,
        *,
        font_size_pt: float = 0,
    ) -> int:
        r"""Add a native PowerPoint equation from LaTeX (Issue #126).

        Args:
            latex: A LaTeX fragment. The supported subset is literal text,
                Greek letters and common symbols, ``^`` and ``_`` scripts,
                ``\\frac{a}{b}``, ``\\sqrt{x}``, ``\\sqrt[n]{x}``, and
                ``\\sum`` / ``\\prod`` / ``\\int`` with optional limits.
                Anything outside it raises rather than being dropped.
            left: Left edge of the equation box in EMU.
            top: Top edge of the equation box in EMU.
            width: Width of the equation box in EMU.
            height: Height of the equation box in EMU.
            font_size_pt: Equation text size; 0 leaves it at 24pt.

        Returns:
            Shape ID of the new equation box.

        The formula is written as OMML, the markup PowerPoint's own equation
        editor uses, so it stays editable and scales with the slide instead of
        being pasted in as a picture.
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "latex": latex,
            "x": int(left),
            "y": int(top),
            "w": int(width),
            "h": int(height),
        }
        if font_size_pt:
            payload["font_size_pt"] = float(font_size_pt)
        result = self._presentation.execute(ops.OP_ADD_EQUATION, payload)
        self._invalidate_shape_cache_if_present()
        self._invalidate_text_state_cache_if_present()
        return int(cast("int", result.get("shape_id", 0)))
