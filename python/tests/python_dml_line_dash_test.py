"""Parity tests for generic shape line dash-style controls."""

from __future__ import annotations

import pytest
from gopptx import Presentation


def test_shape_line_dash_rejects_unsupported_value() -> None:
    """Test shape line dash rejects unsupported value."""
    with Presentation.new("DML Line Dash Invalid") as prs:
        slide = prs.slides[0]
        with pytest.raises(Exception, match=r"line.dash_style"):
            slide.add_shape(
                "rect",
                (1000000, 1000000, 5000000, 1500000),
                properties={"line": {"dash_style": "zigzag"}},
            )
