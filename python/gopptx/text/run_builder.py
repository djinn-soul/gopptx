"""Fluent RunBuilder for composing rich text runs."""

from __future__ import annotations

from typing import TYPE_CHECKING

from typing_extensions import override

if TYPE_CHECKING:
    from ..schemas_shape_types import TextRun


class RunBuilder:
    """Fluent builder that produces a :class:`~gopptx.schemas_shape_types.TextRun` dict.

    Example::

        from gopptx.text import RunBuilder

        proxy.text_frame.set_runs([
            RunBuilder("Hello ").bold().color("FF0000").size_pt(24),
            RunBuilder("World").italic().font("Calibri"),
        ])
    """

    def __init__(self, text: str = "") -> None:
        """Create a builder with initial text content."""
        super().__init__()
        self._payload: dict[str, object] = {"text": text}

    # ------------------------------------------------------------------
    # Text content
    # ------------------------------------------------------------------

    def text(self, value: str) -> RunBuilder:
        """Set the text of the run."""
        self._payload["text"] = value
        return self

    # ------------------------------------------------------------------
    # Styling  (boolean flags accept keyword-only ``value`` to satisfy FBT rules)
    # ------------------------------------------------------------------

    def bold(self, *, value: bool = True) -> RunBuilder:
        """Apply bold formatting."""
        self._payload["bold"] = value
        return self

    def italic(self, *, value: bool = True) -> RunBuilder:
        """Apply italic formatting."""
        self._payload["italic"] = value
        return self

    def underline(self, style: str = "sng") -> RunBuilder:
        """Apply underline formatting.

        Args:
            style: Underline style token e.g. ``"sng"``, ``"dbl"``, ``"dotted"``.
        """
        self._payload["underline"] = style
        return self

    def strikethrough(self, style: str = "sng") -> RunBuilder:
        """Apply strikethrough formatting."""
        self._payload["strikethrough"] = style
        return self

    def subscript(self, *, value: bool = True) -> RunBuilder:
        """Apply subscript formatting."""
        self._payload["subscript"] = value
        return self

    def superscript(self, *, value: bool = True) -> RunBuilder:
        """Apply superscript formatting."""
        self._payload["superscript"] = value
        return self

    def color(self, hex_color: str) -> RunBuilder:
        """Set text colour as a hex RGB string (e.g. ``"FF0000"`` for red)."""
        self._payload["color"] = hex_color
        return self

    def highlight(self, hex_color: str) -> RunBuilder:
        """Set highlight colour as a hex RGB string."""
        self._payload["highlight"] = hex_color
        return self

    def font(self, name: str) -> RunBuilder:
        """Set the font face name."""
        self._payload["font"] = name
        return self

    def size_pt(self, points: int) -> RunBuilder:
        """Set the font size in points."""
        self._payload["size_pt"] = points
        return self

    def code(self, *, value: bool = True) -> RunBuilder:
        """Apply monospace / code formatting."""
        self._payload["code"] = value
        return self

    def all_caps(self, *, value: bool = True) -> RunBuilder:
        """Apply all-caps formatting."""
        self._payload["all_caps"] = value
        return self

    def small_caps(self, *, value: bool = True) -> RunBuilder:
        """Apply small-caps formatting."""
        self._payload["small_caps"] = value
        return self

    # ------------------------------------------------------------------
    # Hyperlink / hover action
    # ------------------------------------------------------------------

    def hyperlink(self, address: str, *, tooltip: str = "") -> RunBuilder:
        """Attach a URL hyperlink to the run.

        Args:
            address: The URL or slide target.
            tooltip: Optional hover tooltip text.
        """
        link: dict[str, object] = {"address": address}
        if tooltip:
            link["tooltip"] = tooltip
        self._payload["hyperlink"] = link
        return self

    def hover_action(self, address: str) -> RunBuilder:
        """Attach a hover-action hyperlink to the run."""
        self._payload["hover_action"] = {"address": address}
        return self

    # ------------------------------------------------------------------
    # Build
    # ------------------------------------------------------------------

    def build(self) -> TextRun:
        """Return the accumulated :class:`~gopptx.schemas_shape_types.TextRun` dict."""
        from ..schemas_shape_types import TextRun  # noqa: PLC0415

        return TextRun(**dict(self._payload))  # type: ignore[misc]

    def to_payload(self) -> dict[str, object]:
        """Return the raw dict payload (used internally by ``set_runs``)."""
        return dict(self._payload)

    @override
    def __repr__(self) -> str:
        """Return a developer-friendly representation."""
        text = self._payload.get("text", "")
        return f"RunBuilder({text!r})"
