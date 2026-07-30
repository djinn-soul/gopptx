"""Test PatternType and MSO_PATTERN enum values (Issue #1127)."""

from gopptx import MSO_PATTERN, PatternType


def test_mso_pattern_percent_40_fixed() -> None:
    """MSO_PATTERN.PERCENT_40 must equal 'pct40' without typos (Issue #1127)."""
    assert hasattr(MSO_PATTERN, "PERCENT_40")
    assert MSO_PATTERN.PERCENT_40 == "pct40"
    assert PatternType.PERCENT_40 == "pct40"
    assert MSO_PATTERN.PERCENT_5 == "pct5"
    assert MSO_PATTERN.DIAGONAL_CROSS == "diagCross"
