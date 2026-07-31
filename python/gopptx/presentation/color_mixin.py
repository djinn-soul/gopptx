"""Presentation-level colour resolution against the theme colour scheme."""

from __future__ import annotations

from typing import TYPE_CHECKING

from .. import ops
from ..color import ColorFormat
from .helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ..color import ThemeColor

# The 12 standard OOXML theme colour slots, as the bridge reports them.
_SCHEME_SLOTS = (
    "dk1",
    "lt1",
    "dk2",
    "lt2",
    "accent1",
    "accent2",
    "accent3",
    "accent4",
    "accent5",
    "accent6",
    "hlink",
    "fol_hlink",
)


class PresentationColorMixin(PresentationMixinBase):
    """Reads the theme colour scheme and resolves colours against it."""

    @property
    def theme_color_scheme(self) -> dict[str, str]:
        """Return the theme's colour slots as ``{slot: hex}``.

        Slots the theme leaves undefined are omitted. Colours defined with
        ``a:sysClr`` resolve to their recorded ``lastClr`` fallback.
        """
        result = self.execute(ops.OP_GET_THEME_COLOR_SCHEME, {})
        scheme: dict[str, str] = {}
        for slot in _SCHEME_SLOTS:
            value = result.get(slot)
            if isinstance(value, str) and value:
                scheme[slot] = value.upper()
        return scheme

    def resolve_color(self, color: ColorFormat) -> str | None:
        """Resolve a colour to hex RGB, looking up theme colours in this deck.

        Returns:
            The resolved hex string, or ``None`` when the colour is unset or
            its theme slot is undefined.
        """
        return color.to_rgb(self.theme_color_scheme)

    def theme_color_rgb(
        self, theme_color: ThemeColor | str, *, brightness: float = 0.0
    ) -> str | None:
        """Return one theme slot as hex RGB, with an optional brightness tweak.

        Accepts a ThemeColor or the bare slot name, e.g. ``"accent1"``.

        Raises:
            ValueError: If brightness is outside -1.0..1.0, or the slot name is
                not an OOXML scheme slot.
        """
        return self.resolve_color(
            ColorFormat.from_theme(theme_color, brightness=brightness)
        )
