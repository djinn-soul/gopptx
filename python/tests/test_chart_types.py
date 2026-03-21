"""Tests for ChartType enum and chart type validation."""

import unittest

from gopptx.presentation.charts import ChartType


class TestChartTypeConstants(unittest.TestCase):
    """Test ChartType named constants."""

    def test_column_chart_type(self) -> None:
        """Test COLUMN constant."""
        self.assertEqual(ChartType.COLUMN, "bar")

    def test_bar_alias(self) -> None:
        """Test BAR is alias for COLUMN."""
        self.assertEqual(ChartType.BAR, ChartType.COLUMN)
        self.assertEqual(ChartType.BAR, "bar")

    def test_line_chart_type(self) -> None:
        """Test LINE constant."""
        self.assertEqual(ChartType.LINE, "line")

    def test_pie_chart_type(self) -> None:
        """Test PIE constant."""
        self.assertEqual(ChartType.PIE, "pie")


class TestChartTypeGetAll(unittest.TestCase):
    """Test ChartType.get_all() method."""

    def test_get_all_returns_dict(self) -> None:
        """Test get_all returns a dictionary."""
        result = ChartType.get_all()
        self.assertIsInstance(result, dict)

    def test_get_all_has_minimum_entries(self) -> None:
        """Test get_all has at least 3 entries (currently supported types)."""
        result = ChartType.get_all()
        self.assertGreaterEqual(len(result), 3)

    def test_get_all_contains_column(self) -> None:
        """Test get_all contains COLUMN."""
        result = ChartType.get_all()
        self.assertIn("COLUMN", result)
        self.assertEqual(result["COLUMN"], "bar")

    def test_get_all_contains_line(self) -> None:
        """Test get_all contains LINE."""
        result = ChartType.get_all()
        self.assertIn("LINE", result)
        self.assertEqual(result["LINE"], "line")

    def test_get_all_contains_pie(self) -> None:
        """Test get_all contains PIE."""
        result = ChartType.get_all()
        self.assertIn("PIE", result)
        self.assertEqual(result["PIE"], "pie")

    def test_get_all_values_are_strings(self) -> None:
        """Test all values in get_all are strings."""
        result = ChartType.get_all()
        for name, value in result.items():
            self.assertIsInstance(name, str)
            self.assertIsInstance(value, str)

    def test_get_all_includes_aliases(self) -> None:
        """Test get_all includes alias entries."""
        result = ChartType.get_all()
        # Both COLUMN and BAR should exist
        self.assertIn("COLUMN", result)
        self.assertIn("BAR", result)
        # They should have the same value
        self.assertEqual(result["COLUMN"], result["BAR"])


class TestChartTypeValidate(unittest.TestCase):
    """Test ChartType.validate() method."""

    def test_validate_constant_name_rejects(self) -> None:
        """Test validate rejects constant name strings."""
        # "COLUMN" is not accepted - must use the actual chart type value
        with self.assertRaises(ValueError):
            ChartType.validate("COLUMN")

    def test_validate_raw_chart_type(self) -> None:
        """Test validate with raw chart type value."""
        result = ChartType.validate("bar")
        self.assertEqual(result, "bar")

    def test_validate_line_constant_name_rejects(self) -> None:
        """Test validate rejects LINE constant name."""
        # "LINE" is a constant name, not a chart type value
        with self.assertRaises(ValueError):
            ChartType.validate("LINE")

    def test_validate_line_raw(self) -> None:
        """Test validate accepts raw line type."""
        result = ChartType.validate("line")
        self.assertEqual(result, "line")

    def test_validate_pie_constant_name_rejects(self) -> None:
        """Test validate rejects PIE constant name."""
        # "PIE" is a constant name, not a chart type value
        with self.assertRaises(ValueError):
            ChartType.validate("PIE")

    def test_validate_pie_raw(self) -> None:
        """Test validate accepts raw pie type."""
        result = ChartType.validate("pie")
        self.assertEqual(result, "pie")

    def test_validate_invalid_raises_error(self) -> None:
        """Test validate with invalid type raises ValueError."""
        with self.assertRaises(ValueError):
            ChartType.validate("INVALID_CHART_TYPE")

    def test_validate_invalid_error_message(self) -> None:
        """Test validate error message includes valid options."""
        with self.assertRaises(ValueError) as context:
            ChartType.validate("INVALID")
        error_msg = str(context.exception)
        self.assertIn("Invalid chart_type", error_msg)

    def test_validate_empty_string_raises_error(self) -> None:
        """Test validate with empty string raises error."""
        with self.assertRaises(ValueError):
            ChartType.validate("")

    def test_validate_none_raises_error(self) -> None:
        """Test validate with None raises error."""
        with self.assertRaises(ValueError):
            ChartType.validate(None)  # type: ignore

    def test_validate_rejects_constant_names(self) -> None:
        """Test validate rejects constant names (uppercase)."""
        # "COLUMN" is a constant name, not a valid chart type value
        with self.assertRaises(ValueError):
            ChartType.validate("COLUMN")

    def test_validate_case_sensitive(self) -> None:
        """Test validate is case sensitive."""
        # "column" (lowercase) is not a valid chart type
        with self.assertRaises(ValueError):
            ChartType.validate("column")

    def test_validate_supported_types(self) -> None:
        """Test validate with currently supported chart type values."""
        test_cases = [
            ("bar", "bar"),
            ("line", "line"),
            ("pie", "pie"),
        ]
        for input_val, expected in test_cases:
            result = ChartType.validate(input_val)
            self.assertEqual(result, expected)


class TestChartTypeGetByName(unittest.TestCase):
    """Test ChartType.get_by_name() method."""

    def test_get_by_name_column(self) -> None:
        """Test get_by_name for COLUMN."""
        result = ChartType.get_by_name("COLUMN")
        self.assertEqual(result, "bar")

    def test_get_by_name_line(self) -> None:
        """Test get_by_name for LINE."""
        result = ChartType.get_by_name("LINE")
        self.assertEqual(result, "line")

    def test_get_by_name_pie(self) -> None:
        """Test get_by_name for PIE."""
        result = ChartType.get_by_name("PIE")
        self.assertEqual(result, "pie")

    def test_get_by_name_not_found(self) -> None:
        """Test get_by_name returns None when not found."""
        result = ChartType.get_by_name("NOT_FOUND")
        self.assertIsNone(result)

    def test_get_by_name_case_sensitive(self) -> None:
        """Test get_by_name is case sensitive."""
        # "line" should not match "LINE"
        result = ChartType.get_by_name("line")
        self.assertIsNone(result)

    def test_get_by_name_all_constants(self) -> None:
        """Test get_by_name for all defined constants."""
        all_types = ChartType.get_all()
        for name in all_types:
            result = ChartType.get_by_name(name)
            self.assertEqual(result, all_types[name])


class TestChartTypeIntegration(unittest.TestCase):
    """Integration tests for ChartType usage."""

    def test_constant_equals_raw_string(self) -> None:
        """Test constant value equals corresponding raw string."""
        self.assertEqual(ChartType.COLUMN, "bar")
        self.assertEqual(ChartType.PIE, "pie")
        self.assertEqual(ChartType.LINE, "line")

    def test_validate_accepts_constant_value(self) -> None:
        """Test validate accepts constant value (raw string)."""
        # ChartType.COLUMN is "bar"
        result = ChartType.validate(ChartType.COLUMN)
        self.assertEqual(result, "bar")

    def test_validate_rejects_constant_name_string(self) -> None:
        """Test validate rejects constant name strings."""
        # "COLUMN" is a constant name, not accepted by validate
        with self.assertRaises(ValueError):
            ChartType.validate("COLUMN")

    def test_round_trip_get_by_name_and_validate_value(self) -> None:
        """Test round-trip: constant name -> get_by_name -> validate value."""
        for name in ["COLUMN", "LINE", "PIE", "BAR"]:
            # Get the raw value from the constant name
            value = ChartType.get_by_name(name)
            if value is not None:  # Only test if the constant exists
                # Validate the raw value (not the constant name)
                validated = ChartType.validate(value)
                self.assertEqual(validated, value)


if __name__ == "__main__":
    unittest.main()
