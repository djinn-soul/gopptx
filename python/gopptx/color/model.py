"""Colour model: RGB, theme colours, brightness, and RGB resolution."""

from __future__ import annotations

from enum import Enum
from typing import TYPE_CHECKING

from typing_extensions import override

if TYPE_CHECKING:
    from collections.abc import Mapping

_HEX_DIGITS = 6
_CHANNEL_MAX = 255
_BRIGHTNESS_MIN = -1.0
_BRIGHTNESS_MAX = 1.0


class ColorType(str, Enum):
    """How a colour is specified."""

    RGB = "rgb"
    """An explicit hex RGB value."""

    SCHEME = "scheme"
    """A reference to a slot in the presentation's theme colour scheme."""

    NONE = "none"
    """No colour: the shape inherits it or is left unfilled."""


class ThemeColor(str, Enum):
    """Theme colour slots, named as OOXML writes them in ``a:schemeClr``.

    ``TEXT_1``/``BACKGROUND_1`` are the aliases PowerPoint's UI uses for the
    ``dk1``/``lt1`` slots, and they resolve to the same scheme entries.
    """

    DARK_1 = "dk1"
    LIGHT_1 = "lt1"
    DARK_2 = "dk2"
    LIGHT_2 = "lt2"
    TEXT_1 = "tx1"
    BACKGROUND_1 = "bg1"
    TEXT_2 = "tx2"
    BACKGROUND_2 = "bg2"
    ACCENT_1 = "accent1"
    ACCENT_2 = "accent2"
    ACCENT_3 = "accent3"
    ACCENT_4 = "accent4"
    ACCENT_5 = "accent5"
    ACCENT_6 = "accent6"
    HYPERLINK = "hlink"
    FOLLOWED_HYPERLINK = "folHlink"


# tx1/bg1/tx2/bg2 are aliases; the scheme itself only stores dk/lt slots.
_SCHEME_SLOT_ALIASES: Mapping[str, str] = {
    "tx1": "dk1",
    "bg1": "lt1",
    "tx2": "dk2",
    "bg2": "lt2",
    "folHlink": "fol_hlink",
}


class ColorFormat:
    """A colour value: explicit RGB or a theme slot, plus a brightness tweak.

    ``brightness`` runs from -1.0 (black) through 0.0 (the colour itself) to
    1.0 (white), matching how PowerPoint presents "Lighter"/"Darker" variants.
    """

    __slots__ = ("_brightness", "_rgb", "_theme_color", "_type")

    def __init__(
        self,
        color_type: ColorType = ColorType.NONE,
        *,
        rgb: str | None = None,
        theme_color: ThemeColor | None = None,
        brightness: float = 0.0,
    ) -> None:
        """Build a colour of the given type.

        Raises:
            ValueError: If the value does not match the type, the hex string is
                malformed, or brightness is outside -1.0..1.0.
        """
        super().__init__()
        if not _BRIGHTNESS_MIN <= brightness <= _BRIGHTNESS_MAX:
            raise ValueError("brightness must be between -1.0 and 1.0")
        if color_type is ColorType.RGB:
            if rgb is None:
                raise ValueError("an RGB colour requires rgb")
            rgb = normalize_hex_color(rgb)
        elif color_type is ColorType.SCHEME:
            if theme_color is None:
                raise ValueError("a scheme colour requires theme_color")
        self._type = color_type
        self._rgb = rgb
        self._theme_color = theme_color
        self._brightness = float(brightness)

    @classmethod
    def from_rgb(cls, rgb: str, *, brightness: float = 0.0) -> ColorFormat:
        """Build an explicit RGB colour from a hex string."""
        return cls(ColorType.RGB, rgb=rgb, brightness=brightness)

    @classmethod
    def from_theme(
        cls, theme_color: ThemeColor, *, brightness: float = 0.0
    ) -> ColorFormat:
        """Build a colour that references a theme slot."""
        return cls(ColorType.SCHEME, theme_color=theme_color, brightness=brightness)

    @property
    def type(self) -> ColorType:
        """Return how this colour is specified."""
        return self._type

    @property
    def rgb(self) -> str | None:
        """Return the explicit hex RGB value, or ``None`` for other types."""
        return self._rgb

    @property
    def theme_color(self) -> ThemeColor | None:
        """Return the referenced theme slot, or ``None`` for other types."""
        return self._theme_color

    @property
    def brightness(self) -> float:
        """Return the brightness adjustment, from -1.0 to 1.0."""
        return self._brightness

    def with_brightness(self, brightness: float) -> ColorFormat:
        """Return a copy of this colour with a different brightness."""
        return ColorFormat(
            self._type,
            rgb=self._rgb,
            theme_color=self._theme_color,
            brightness=brightness,
        )

    def to_rgb(self, scheme: Mapping[str, str] | None = None) -> str | None:
        """Resolve this colour to a hex RGB string.

        Args:
            scheme: The presentation's theme colour scheme, needed only to
                resolve a theme colour. ``Presentation.theme_color_scheme``
                returns it in the expected shape.

        Returns:
            The resolved hex string, or ``None`` when there is no colour or the
            scheme does not define the referenced slot.
        """
        base = self._base_rgb(scheme)
        if base is None:
            return None
        return apply_brightness(base, self._brightness)

    def _base_rgb(self, scheme: Mapping[str, str] | None) -> str | None:
        if self._type is ColorType.RGB:
            return self._rgb
        if self._type is not ColorType.SCHEME or self._theme_color is None:
            return None
        if scheme is None:
            return None
        slot = self._theme_color.value
        value = scheme.get(slot) or scheme.get(_SCHEME_SLOT_ALIASES.get(slot, slot))
        if not isinstance(value, str) or not value:
            return None
        return normalize_hex_color(value)

    @override
    def __eq__(self, other: object) -> bool:
        """Compare type, value, and brightness."""
        if not isinstance(other, ColorFormat):
            return NotImplemented
        return (
            self._type is other._type
            and self._rgb == other._rgb
            and self._theme_color is other._theme_color
            and self._brightness == other._brightness
        )

    @override
    def __hash__(self) -> int:
        """Hash on the same fields used for equality."""
        return hash((self._type, self._rgb, self._theme_color, self._brightness))

    @override
    def __repr__(self) -> str:
        """Return a developer-friendly representation."""
        if self._type is ColorType.RGB:
            body = f"rgb={self._rgb!r}"
        elif self._type is ColorType.SCHEME and self._theme_color is not None:
            body = f"theme_color={self._theme_color.name}"
        else:
            body = "none"
        return f"ColorFormat({body}, brightness={self._brightness})"


def normalize_hex_color(value: str) -> str:
    """Return an uppercase 6-digit hex string, accepting a leading ``#``.

    Raises:
        ValueError: If the value is not a 6-digit hex colour.
    """
    text = value.strip().removeprefix("#")
    if len(text) != _HEX_DIGITS:
        raise ValueError(f"expected a 6-digit hex colour, got {value!r}")
    try:
        _ = int(text, 16)
    except ValueError as exc:
        raise ValueError(f"expected a 6-digit hex colour, got {value!r}") from exc
    return text.upper()


def apply_brightness(rgb: str, brightness: float) -> str:
    """Lighten or darken a hex colour.

    Positive brightness moves each channel toward white and negative toward
    black, which is how PowerPoint renders its lighter/darker theme variants.
    """
    if not brightness:
        return normalize_hex_color(rgb)
    base = normalize_hex_color(rgb)
    channels = [int(base[index : index + 2], 16) for index in (0, 2, 4)]
    adjusted: list[int] = []
    for channel in channels:
        if brightness > 0:
            value = channel + (_CHANNEL_MAX - channel) * brightness
        else:
            value = channel * (1.0 + brightness)
        adjusted.append(max(0, min(_CHANNEL_MAX, round(value))))
    return "".join(f"{value:02X}" for value in adjusted)
