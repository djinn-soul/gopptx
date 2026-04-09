"""Slide mixin providing SmartArt, animation, and transition methods."""

from __future__ import annotations

from typing import TYPE_CHECKING

from ... import ops
from ...presentation.helpers import get_required_int

if TYPE_CHECKING:
    from ..contracts import SlidePresentationProtocol

# EMU conversions: 914400 EMU = 1 inch
_INCHES_TO_EMU = 914400


def _optional_payload_str(value: object) -> str:
    """Return empty string for None, string form otherwise."""
    if value is None:
        return ""
    return str(value)


class SlideSmartArtAnimMixin:
    """Mixin adding add_smartart, add_animation, and set_transition to Slide."""

    if TYPE_CHECKING:
        _presentation: SlidePresentationProtocol  # pyright: ignore[reportUninitializedInstanceVariable]

        @property
        def index(self) -> int:
            """Slide index."""
            ...

        def _invalidate_shape_cache_if_present(self) -> None: ...

    def add_smartart(
        self,
        layout: str,
        items: list[str] | None = None,
        bounds: tuple[float, float, float, float] | None = None,
        *,
        nodes: list[dict[str, object]] | None = None,
    ) -> int:
        """Add a SmartArt diagram to the slide.

        Args:
            layout: SmartArt layout URI.  Use constants from
                ``gopptx.smartart`` (e.g. ``SMARTART_BASIC_LIST``) or pass a
                raw OOXML URN.
            items: Flat list of text items (for list/process/cycle layouts).
                Ignored when ``nodes`` is provided.
            bounds: ``(left, top, width, height)`` in inches.
            nodes: Nested node tree for hierarchy layouts (org chart, hierarchy).
                Each node is a dict with ``"text"`` (str) and optional
                ``"children"`` (list of the same shape).

        Returns:
            The shape ID of the inserted graphic frame.

        Examples::

            from gopptx.smartart import SMARTART_BASIC_PROCESS, SMARTART_ORG_CHART

            # Flat items for process/list layouts
            shape_id = slide.add_smartart(
                SMARTART_BASIC_PROCESS,
                ["Plan", "Execute", "Review"],
                (1.0, 2.0, 8.0, 4.0),
            )

            # Nested nodes for org chart / hierarchy
            shape_id = slide.add_smartart(
                SMARTART_ORG_CHART,
                bounds=(1.0, 2.0, 8.0, 4.0),
                nodes=[
                    {"text": "CEO", "children": [
                        {"text": "VP Engineering", "children": [
                            {"text": "Engineer"},
                        ]},
                        {"text": "VP Sales"},
                    ]},
                ],
            )
        """
        if bounds is None:
            bounds = (1.0, 2.0, 8.0, 4.0)
        left, top, width, height = bounds
        x = int(left * _INCHES_TO_EMU)
        y = int(top * _INCHES_TO_EMU)
        cx = int(width * _INCHES_TO_EMU)
        cy = int(height * _INCHES_TO_EMU)

        payload: dict[str, object] = {
            "slide_index": self.index,
            "layout": layout,
            "x": x,
            "y": y,
            "cx": cx,
            "cy": cy,
        }
        if nodes is not None:
            payload["nodes"] = nodes
        else:
            payload["items"] = items or []

        result = self._presentation.execute(ops.OP_ADD_SMART_ART, payload)
        self._invalidate_shape_cache_if_present()
        return get_required_int(result, "shape_id")

    def update_smartart(
        self,
        shape_id: int,
        items: list[str],
    ) -> None:
        """Replace text items in an existing SmartArt diagram.

        Args:
            shape_id: The shape ID of the SmartArt graphic frame.
            items: New list of text items.  Excess items are dropped;
                extra slots are cleared.

        Example::

            slide.update_smartart(shape_id, ["Step 1", "Step 2", "Step 3"])
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "items": items,
        }
        self._presentation.execute(ops.OP_UPDATE_SMART_ART, payload)
        self._invalidate_shape_cache_if_present()

    def delete_smartart(self, shape_id: int) -> None:
        """Delete a SmartArt diagram by shape ID."""
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
        }
        self._presentation.execute(ops.OP_DELETE_SMART_ART, payload)
        self._invalidate_shape_cache_if_present()

    def change_smartart_layout(self, shape_id: int, layout: str) -> None:
        """Change the layout URI of an existing SmartArt diagram."""
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "layout": layout,
        }
        self._presentation.execute(ops.OP_CHANGE_SMART_ART_LAYOUT, payload)
        self._invalidate_shape_cache_if_present()

    def set_smartart_style(
        self,
        shape_id: int,
        *,
        quick_style: str | None = None,
        color_style: str | None = None,
    ) -> None:
        """Set SmartArt quick style and/or color style URIs."""
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
        }
        if quick_style is not None:
            payload["quick_style"] = quick_style
        if color_style is not None:
            payload["color_style"] = color_style
        self._presentation.execute(ops.OP_SET_SMART_ART_STYLE, payload)
        self._invalidate_shape_cache_if_present()

    def set_smartart_nodes(self, shape_id: int, items: list[str]) -> None:
        """Replace SmartArt node text using a flat items list."""
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "items": items,
        }
        self._presentation.execute(ops.OP_SET_SMART_ART_NODES, payload)
        self._invalidate_shape_cache_if_present()

    def add_animation(
        self,
        shape_id: int,
        effect: str,
        *,
        trigger: str = "onClick",
        duration_ms: int = 500,
        delay_ms: int = 0,
    ) -> None:
        """Add an animation effect to a shape on this slide.

        Args:
            shape_id: Numeric shape ID as returned by add_shape/add_textbox etc.
            effect: Effect token.  Use constants from ``gopptx.animations``
                (e.g. ``ANIMATION_ENTRANCE_FADE``).
            trigger: When the animation starts.  One of ``"onClick"``,
                ``"withPrev"``, or ``"afterPrev"``.
            duration_ms: Duration in milliseconds (default 500).
            delay_ms: Delay before the animation starts, in milliseconds.

        Example::

            from gopptx.animations import ANIMATION_ENTRANCE_FADE
            slide.add_animation(shape_id, ANIMATION_ENTRANCE_FADE, duration_ms=800)
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "shape_id": shape_id,
            "effect": effect,
            "trigger": trigger,
            "duration_ms": duration_ms,
            "delay_ms": delay_ms,
        }
        self._presentation.execute(ops.OP_ADD_ANIMATION, payload)

    def set_transition(
        self,
        transition_type: str,
        *,
        duration_ms: int = 0,
        advance_ms: int | None = None,
        disable_advance_on_click: bool = False,
    ) -> None:
        """Set the slide transition.

        Args:
            transition_type: Transition token.  Use constants from
                ``gopptx.transitions`` (e.g. ``TRANSITION_FADE``).
            duration_ms: Transition duration in milliseconds (0 = default).
            advance_ms: Auto-advance after this many milliseconds.  ``None``
                means click-advance only (the default).
            disable_advance_on_click: When ``True``, disable click-to-advance
                for this slide transition (writes ``advClick="0"``).

        Example::

            from gopptx.transitions import TRANSITION_PUSH
            slide.set_transition(TRANSITION_PUSH, duration_ms=600)
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "transition_type": transition_type,
            "duration_ms": duration_ms,
            "advance_ms": advance_ms if advance_ms is not None else -1,
            "disable_advance_on_click": disable_advance_on_click,
        }
        self._presentation.execute(ops.OP_SET_SLIDE_TRANSITION, payload)

    def set_background(self, bg_type: str, **kwargs: object) -> None:
        """Set the slide background.

        Args:
            bg_type: Background type: ``"solid"``, ``"gradient"``,
                ``"image"``, or ``"theme"``.
            **kwargs: Optional background settings:
                ``color`` - Hex RGB color for solid backgrounds, e.g. ``"FF0000"``.
                ``colors`` - List of hex RGB colors for gradient backgrounds.
                ``angle`` - Gradient angle in degrees (0-360).
                ``image_path`` - Local file path for an image background.
                ``image_data`` - Base64-encoded image data for an image background.
                ``color_ref`` - Theme color token for theme backgrounds, e.g. ``"accent1"``.

        Example::

            slide.set_background("solid", color="3070B3")
            slide.set_background("gradient", colors=["FF0000", "0000FF"], angle=90)
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "type": bg_type,
            "color": _optional_payload_str(kwargs.get("color")),
            "colors": list(kwargs.get("colors") or []),  # type: ignore[arg-type]
            "angle": int(kwargs.get("angle") or 0),  # type: ignore[arg-type]
            "image_path": _optional_payload_str(kwargs.get("image_path")),
            "image_data": _optional_payload_str(kwargs.get("image_data")),
            "color_ref": _optional_payload_str(kwargs.get("color_ref")),
        }
        self._presentation.execute(ops.OP_SET_SLIDE_BACKGROUND, payload)

    def set_header_footer(
        self,
        *,
        footer: str = "",
        show_footer: bool = False,
        show_slide_num: bool = False,
        show_date_time: bool = False,
        date_time_text: str = "",
    ) -> None:
        """Configure the header/footer overlay for this slide.

        Args:
            footer: Footer text to display.
            show_footer: Whether to show the footer.
            show_slide_num: Whether to show the slide number.
            show_date_time: Whether to show the date/time.
            date_time_text: Fixed date/time string (empty = auto).

        Example::

            slide.set_header_footer(footer="Confidential", show_footer=True,
                                    show_slide_num=True)
        """
        payload: dict[str, object] = {
            "slide_index": self.index,
            "footer": footer,
            "show_footer": show_footer,
            "show_slide_num": show_slide_num,
            "show_date_time": show_date_time,
            "date_time_text": date_time_text,
        }
        self._presentation.execute(ops.OP_SET_SLIDE_HEADER_FOOTER, payload)
