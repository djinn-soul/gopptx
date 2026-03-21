"""Table style constants and utilities for gopptx.

Instead of using raw GUIDs like "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}",
use named constants for common table styles.
"""

from __future__ import annotations


class TableStyle:
    """Named table style constants for common PowerPoint table styles."""

    # Light Styles
    LIGHT_STYLE_1 = "{B9AC3A68-259E-4EED-9050-4AE35E7F2B2D}"
    LIGHT_STYLE_2 = "{B4BFB6E8-F1D8-4F8B-9FA3-A8B7E5D8A1C2}"
    LIGHT_STYLE_3 = "{C1D9E8A0-5F7A-4B2E-8D6F-9E2A5C8B1D3F}"

    # Medium Styles
    MEDIUM_STYLE_1 = "{5940675A-B579-460E-94D1-54222C63F5DA}"
    MEDIUM_STYLE_1_ACCENT_1 = "{5940675A-B579-460E-94D1-54222C63F5DA}"
    MEDIUM_STYLE_2 = "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"
    MEDIUM_STYLE_2_ACCENT_1 = "{5C22544A-7EE6-4342-B048-85BDC9FD1C3A}"
    MEDIUM_STYLE_3 = "{3D296D3D-A3F5-4A1C-9F2E-B8D1C5A7E9F2}"
    MEDIUM_STYLE_4 = "{A1C4D3E8-7F2B-4A5E-9D1F-C6E8A2B5F7D1}"

    # Dark Styles
    DARK_STYLE_1 = "{6C7A8F3D-2E5C-4B7A-9F1D-E8C2A5B7D9E1}"
    DARK_STYLE_2 = "{8F5C3A1E-7D2F-4A8C-B1E5-D7A9C3E8F2B6}"

    # Themed Styles (Accent colors)
    THEMED_STYLE_1_ACCENT_1 = "{F472CDA9-FDAA-48DA-9F9A-0FF3F4F8F037}"
    THEMED_STYLE_1_ACCENT_2 = "{D4C8D5E9-3A7F-4B2E-8D6F-9E2A5C8B1D3F}"
    THEMED_STYLE_1_ACCENT_3 = "{C1D9E8A0-5F7A-4B2E-8D6F-9E2A5C8B1D3F}"
    THEMED_STYLE_1_ACCENT_4 = "{B4BFB6E8-F1D8-4F8B-9FA3-A8B7E5D8A1C2}"
    THEMED_STYLE_1_ACCENT_5 = "{A1C4D3E8-7F2B-4A5E-9D1F-C6E8A2B5F7D1}"
    THEMED_STYLE_1_ACCENT_6 = "{8F5C3A1E-7D2F-4A8C-B1E5-D7A9C3E8F2B6}"

    # Default style (no style applied)
    DEFAULT = ""

    @staticmethod
    def get_all() -> dict[str, str]:
        """Return all named styles as a dict.

        Returns:
            Dictionary mapping style names to GUIDs.

        Example:
            >>> styles = TableStyle.get_all()
            >>> for name, guid in styles.items():
            ...     print(f"{name}: {guid}")
        """
        return {
            "LIGHT_STYLE_1": TableStyle.LIGHT_STYLE_1,
            "LIGHT_STYLE_2": TableStyle.LIGHT_STYLE_2,
            "LIGHT_STYLE_3": TableStyle.LIGHT_STYLE_3,
            "MEDIUM_STYLE_1": TableStyle.MEDIUM_STYLE_1,
            "MEDIUM_STYLE_1_ACCENT_1": TableStyle.MEDIUM_STYLE_1_ACCENT_1,
            "MEDIUM_STYLE_2": TableStyle.MEDIUM_STYLE_2,
            "MEDIUM_STYLE_2_ACCENT_1": TableStyle.MEDIUM_STYLE_2_ACCENT_1,
            "MEDIUM_STYLE_3": TableStyle.MEDIUM_STYLE_3,
            "MEDIUM_STYLE_4": TableStyle.MEDIUM_STYLE_4,
            "DARK_STYLE_1": TableStyle.DARK_STYLE_1,
            "DARK_STYLE_2": TableStyle.DARK_STYLE_2,
            "THEMED_STYLE_1_ACCENT_1": TableStyle.THEMED_STYLE_1_ACCENT_1,
            "THEMED_STYLE_1_ACCENT_2": TableStyle.THEMED_STYLE_1_ACCENT_2,
            "THEMED_STYLE_1_ACCENT_3": TableStyle.THEMED_STYLE_1_ACCENT_3,
            "THEMED_STYLE_1_ACCENT_4": TableStyle.THEMED_STYLE_1_ACCENT_4,
            "THEMED_STYLE_1_ACCENT_5": TableStyle.THEMED_STYLE_1_ACCENT_5,
            "THEMED_STYLE_1_ACCENT_6": TableStyle.THEMED_STYLE_1_ACCENT_6,
            "DEFAULT": TableStyle.DEFAULT,
        }


__all__ = ["TableStyle"]
