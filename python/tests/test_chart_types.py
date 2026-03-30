"""Tests for ChartType enum and chart type validation."""

import unittest

from gopptx.presentation.charts import ChartType

EXPECTED_CHART_TYPES = {
    "COLUMN": "bar",
    "BAR": "bar",
    "BAR_HORIZONTAL": "barHorizontal",
    "BAR_STACKED": "barStacked",
    "BAR_STACKED_100": "barStacked100",
    "LINE": "line",
    "LINE_MARKERS": "lineMarkers",
    "LINE_STACKED": "lineStacked",
    "SCATTER": "scatter",
    "AREA": "area",
    "AREA_STACKED": "areaStacked",
    "AREA_STACKED_100": "areaStacked100",
    "PIE": "pie",
    "DOUGHNUT": "doughnut",
    "BUBBLE": "bubble",
    "RADAR": "radar",
    "RADAR_FILLED": "radarFilled",
    "STOCK_HLC": "stockHLC",
    "STOCK_OHLC": "stockOHLC",
    "COMBO": "combo",
}


class TestChartTypeConstants(unittest.TestCase):
    """Test ChartType named constants."""

    def test_constant_values_match_expected_map(self) -> None:
        """Every expected constant maps to the expected raw chart value."""
        for name, expected_value in EXPECTED_CHART_TYPES.items():
            actual = getattr(ChartType, name)
            self.assertEqual(actual, expected_value)

    def test_bar_alias(self) -> None:
        """BAR remains an alias for COLUMN."""
        self.assertEqual(ChartType.BAR, ChartType.COLUMN)
        self.assertEqual(ChartType.BAR, "bar")


class TestChartTypeGetAll(unittest.TestCase):
    """Test ChartType.get_all() method."""

    def test_get_all_returns_dict(self) -> None:
        """get_all returns a dictionary."""
        result = ChartType.get_all()
        self.assertIsInstance(result, dict)

    def test_get_all_contains_full_chart_set(self) -> None:
        """get_all contains all exposed chart constants."""
        result = ChartType.get_all()
        self.assertEqual(result, EXPECTED_CHART_TYPES)

    def test_get_all_includes_aliases(self) -> None:
        """Both COLUMN and BAR exist and resolve to the same value."""
        result = ChartType.get_all()
        self.assertIn("COLUMN", result)
        self.assertIn("BAR", result)
        self.assertEqual(result["COLUMN"], result["BAR"])

    def test_get_all_unique_value_count_matches_go_surface(self) -> None:
        """Go surface has 19 distinct chart type values (BAR aliases COLUMN)."""
        result = ChartType.get_all()
        self.assertEqual(len(set(result.values())), 19)


class TestChartTypeValidate(unittest.TestCase):
    """Test ChartType.validate() method."""

    def test_validate_accepts_all_supported_values(self) -> None:
        """validate accepts each raw value exposed by ChartType constants."""
        for value in set(EXPECTED_CHART_TYPES.values()):
            self.assertEqual(ChartType.validate(value), value)

    def test_validate_rejects_constant_name_strings(self) -> None:
        """validate must reject constant names like 'COLUMN'."""
        for name in EXPECTED_CHART_TYPES:
            with self.assertRaises(ValueError):
                ChartType.validate(name)

    def test_validate_invalid_raises_error(self) -> None:
        """validate with invalid type raises ValueError."""
        with self.assertRaises(ValueError):
            ChartType.validate("INVALID_CHART_TYPE")

    def test_validate_invalid_error_message(self) -> None:
        """Error message includes 'Invalid chart_type' marker."""
        with self.assertRaises(ValueError) as context:
            ChartType.validate("INVALID")
        self.assertIn("Invalid chart_type", str(context.exception))

    def test_validate_empty_string_raises_error(self) -> None:
        """validate with empty string raises error."""
        with self.assertRaises(ValueError):
            ChartType.validate("")

    def test_validate_none_raises_error(self) -> None:
        """validate with None raises error."""
        with self.assertRaises(ValueError):
            ChartType.validate(None)  # type: ignore[arg-type]

    def test_validate_case_sensitive(self) -> None:
        """Lowercase constant names are rejected if not raw values."""
        with self.assertRaises(ValueError):
            ChartType.validate("column")


class TestChartTypeGetByName(unittest.TestCase):
    """Test ChartType.get_by_name() method."""

    def test_get_by_name_for_all_constants(self) -> None:
        """get_by_name resolves every known constant name."""
        for name, value in EXPECTED_CHART_TYPES.items():
            self.assertEqual(ChartType.get_by_name(name), value)

    def test_get_by_name_not_found(self) -> None:
        """get_by_name returns None when constant is unknown."""
        self.assertIsNone(ChartType.get_by_name("NOT_FOUND"))

    def test_get_by_name_case_sensitive(self) -> None:
        """get_by_name is case-sensitive."""
        self.assertIsNone(ChartType.get_by_name("line"))


class TestChartTypeIntegration(unittest.TestCase):
    """Integration tests for ChartType usage."""

    def test_validate_accepts_constant_values(self) -> None:
        """validate accepts values passed from constants."""
        for name in EXPECTED_CHART_TYPES:
            value = ChartType.get_by_name(name)
            self.assertIsNotNone(value)
            self.assertEqual(ChartType.validate(value), value)

    def test_combo_round_trip(self) -> None:
        """Combo chart constant resolves and validates."""
        combo_value = ChartType.get_by_name("COMBO")
        self.assertEqual(combo_value, "combo")
        self.assertEqual(ChartType.validate(combo_value), "combo")


if __name__ == "__main__":
    unittest.main()
