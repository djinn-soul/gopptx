"""Chart type coercion shared by the presentation-level and slide-level APIs.

Hand-written companion to the generated :mod:`chart_types` module. Kept separate
so that ``go generate`` can overwrite the enum without touching this logic.
"""

from __future__ import annotations

import warnings

from .chart_types import ChartType

__all__ = ["coerce_chart_type"]

_REMOVAL_NOTE = "Bare chart_type strings will be rejected in a future release."


def coerce_chart_type(chart_type: ChartType | str) -> str:
    """Validate a chart type and return its wire value.

    A :class:`ChartType` member passes straight through. A bare string that
    happens to name a valid chart type is still accepted, but is deprecated:
    the enum is what makes the value discoverable and checkable, and bare
    strings will stop being accepted in a future release.

    Args:
        chart_type: A ChartType member, or its wire value such as "pie".

    Returns:
        The wire value, as a plain string.

    Raises:
        ValueError: If chart_type is empty or is not a known chart type.
    """
    if not isinstance(chart_type, ChartType):
        message = " ".join((
            f"chart_type should be a ChartType member, not {chart_type!r}.",
            _REMOVAL_NOTE,
        ))
        warnings.warn(message, DeprecationWarning, stacklevel=3)
    return ChartType.validate(chart_type)
