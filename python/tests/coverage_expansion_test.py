"""Coverage-expansion tests for runtime, batch, and helper behaviors."""

import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from gopptx import GopptxError, Presentation, ops
from gopptx.presentation import helpers
from gopptx.utils import normalize_table_index


class TestCoverageExpansion(unittest.TestCase):
    """Exercise lower-traffic API paths for regression coverage."""

    def setUp(self) -> None:
        """Generate a temporary PPTX fixture for each test."""
        self._tmpdir = tempfile.TemporaryDirectory()
        self.test_pptx = str(Path(self._tmpdir.name) / "simple.pptx")
        with Presentation.new("Simple Test Deck") as pres:
            pres.add_slide(
                "Slide One",
                layout="title_and_content",
                bullets=["Bullet A", "Bullet B"],
            )
            pres.add_slide("Slide Two", layout="title_only")
            pres.save(self.test_pptx)

    def tearDown(self) -> None:
        """Clean up the temporary directory."""
        self._tmpdir.cleanup()

    def test_runtime_unopened_errors(self) -> None:
        """Unopened presentation raises consistently for runtime operations."""
        pres = Presentation()
        with self.assertRaisesRegex(GopptxError, "Presentation is not open."):
            pres.execute("any_op")
        with self.assertRaisesRegex(GopptxError, "Presentation is not open."):
            pres.save("any_path")

        self.assertIsInstance(repr(pres), str)

    def test_batch_read_op_not_allowed(self) -> None:
        """Read operations are blocked while batch mode is active."""
        with (
            Presentation(self.test_pptx) as pres,
            pres.batch(),
            self.assertRaisesRegex(
                GopptxError,
                r"read operation 'slide_count' is not allowed inside batch",
            ),
        ):
            pres.execute(ops.OP_SLIDE_COUNT)

    def test_execute_batch_empty(self) -> None:
        """Empty batch submission returns empty result list."""
        with Presentation(self.test_pptx) as pres:
            self.assertEqual(pres.execute_batch([]), [])

    def test_nested_batch_not_allowed(self) -> None:
        """Nested begin_batch calls are rejected."""
        with Presentation(self.test_pptx) as pres:
            pres.begin_batch()
            with self.assertRaisesRegex(
                GopptxError, r"nested batch\(\) calls are not allowed"
            ):
                pres.begin_batch()
            pres.abort_batch()

    def test_abort_batch(self) -> None:
        """Aborting a batch restores pre-batch slide count."""
        with Presentation(self.test_pptx) as pres:
            # Get initial count
            initial_count = pres.slide_count
            pres.begin_batch()
            pres.add_slide("New Slide")
            pres.abort_batch()
            # Count should remain the same
            self.assertEqual(pres.slide_count, initial_count)

    def test_metadata_cache(self) -> None:
        """Metadata cache returns same object until explicit invalidation."""
        with Presentation(self.test_pptx) as pres:
            m1 = pres.metadata
            m2 = pres.metadata
            self.assertIs(m1, m2)
            pres.invalidate_cache()
            m3 = pres.metadata
            self.assertIsNot(m1, m3)

    def test_presentation_getitem(self) -> None:
        """Presentation indexing supports direct, negative, and slice access."""
        with Presentation(self.test_pptx) as pres:
            count = pres.slide_count
            # Simple indexing
            s0 = pres[0]
            self.assertEqual(s0.index, 0)

            # Negative indexing
            s_last = pres[-1]
            self.assertEqual(s_last.index, count - 1)

            with self.assertRaises(IndexError):
                pres[100]

            # Slicing
            slides = pres[0:1]
            self.assertEqual(len(slides), 1)

            with self.assertRaises(TypeError):
                pres["invalid"]

    def test_reopen_already_open(self) -> None:
        """Opening an already-open presentation rebinds cleanly."""
        with Presentation(self.test_pptx) as pres:
            pres.open(self.test_pptx)  # Should close previous and open new
            self.assertGreater(pres.slide_count, 0)

    def test_properties_and_protection(self) -> None:
        """Property and protection operations dispatch through the bridge."""
        with Presentation(self.test_pptx) as pres:
            # Core properties
            props = pres.get_core_properties()
            self.assertIn("title", props)

            pres.set_core_properties(props)

            # Title shortcut
            orig_title = pres.title
            pres.title = "New Title"
            self.assertEqual(pres.title, "New Title")
            pres.title = orig_title

            # Sections
            sections = pres.get_sections()
            self.assertIsInstance(sections, list)

            # Mark as final / password (smoke test for bridge dispatch)
            pres.set_mark_as_final(final=True)
            pres.set_mark_as_final(final=False)
            pres.set_modify_password("secret")

    def test_slide_management_expansion(self) -> None:
        """Slide add/update/move iteration paths remain stable."""
        with Presentation(self.test_pptx) as pres:
            initial_count = pres.slide_count
            # Add with layout
            s2 = pres.add_slide("Layout Slide", layout="title_and_content")
            self.assertEqual(s2.title, "Layout Slide")

            # Update with layout
            pres.update_slide(1, title="Updated Title", layout="blank")

            # Add title/bullet slides
            pres.add_title_slide("Title Only")
            pres.add_bullet_slide("Bullets", ["B1", "B2"])
            self.assertEqual(pres.slide_count, initial_count + 3)

            # Len, Iter
            self.assertEqual(len(pres), initial_count + 3)
            slide_list = list(iter(pres))
            self.assertEqual(len(slide_list), initial_count + 3)

    def test_batch_mode_placeholder_slide(self) -> None:
        """Batch-mode slide creation returns placeholder-safe slide object."""
        with Presentation(self.test_pptx) as pres, pres.batch():
            s = pres.add_slide("Batch Slide")
            # Cannot check index in batch mode because it triggers a read op
            self.assertEqual(s.title, "Batch Slide")

    def test_layout_theme_expansion(self) -> None:
        """Layout rebinding and theme application smoke paths run."""
        with Presentation(self.test_pptx) as pres:
            layouts = pres.list_slide_layouts()
            if layouts:
                name = layouts[0]["name"]
                # Rebind by name
                pres.rebind_slide_layout(0, name)

            # Apply theme
            from gopptx.presentation.theme import get_theme

            pres.apply_theme(get_theme("aurora"))  # Apply a built-in theme

    def test_remove_comment_errors(self) -> None:
        """Comment removal validates index and selector arguments."""
        with Presentation(self.test_pptx) as pres:
            with self.assertRaises(ValueError):
                pres.remove_comment(999)  # Unknown index in cache

            with self.assertRaises(TypeError):
                pres.remove_comment(0, author_id=1)  # Missing author_index

    def test_chart_expansion(self) -> None:
        """Chart add/update dispatch accepts series-style payloads."""
        with Presentation(self.test_pptx) as pres:
            # Invalid bounds
            with self.assertRaises(ValueError):
                pres.add_chart(0, "bar", ["A"], [10.0], bounds=(0, 0, 100))

            # Series-style input
            series = [{"name": "S1", "values": [10.0, 20.0]}]
            pres.add_chart(0, "bar", ["C1", "C2"], series)

            # For update_chart_data with dict selector, we need a valid chart RelID usually.
            charts = pres.list_slide_charts(0)
            if charts:
                rel_id = charts[0].get("RelID") or charts[0].get("rel_id")
                data = {
                    "categories": ["C1", "C2"],
                    "series": [
                        {
                            "name": "S1",
                            "values": [15.0, 25.0],
                            "categories": ["C1", "C2"],
                        }
                    ],
                }
                pres.update_chart_data(0, {"rel_id": rel_id}, data)

    def test_chart_expansion_convenience_methods(self) -> None:
        """Chart convenience update/replace methods dispatch correctly."""
        with Presentation(self.test_pptx) as pres:
            pres.add_chart(0, "bar", ["A", "B"], [10.0, 20.0])
            pres.update_chart_data_by_index(
                0,
                0,
                {
                    "categories": ["A", "B"],
                    "series": [{"values": [11.0, 22.0]}],
                },
            )
            pres.replace_chart_data_by_index(0, 0, ["X", "Y"], [1.0, 2.0])

            charts = pres.list_slide_charts(0)
            if charts:
                rel_id = charts[0].get("RelID") or charts[0].get("rel_id")
                if isinstance(rel_id, str) and rel_id:
                    pres.update_chart_data_by_rel_id(
                        0,
                        rel_id,
                        {
                            "categories": ["R1", "R2"],
                            "series": [{"values": [3.0, 4.0]}],
                        },
                    )
                    pres.replace_chart_data_by_rel_id(
                        0, rel_id, ["RR1", "RR2"], [5.0, 6.0]
                    )

    def test_normalize_table_index(self) -> None:
        """Table-index normalizer accepts integer-like values only."""
        self.assertEqual(normalize_table_index(1.0), 1)
        with self.assertRaises(ValueError):
            normalize_table_index(1.5)
        with self.assertRaises(ValueError):
            normalize_table_index("not_an_int")
        with self.assertRaises(ValueError):
            normalize_table_index(value=True)

    def test_json_fallback(self) -> None:
        """JSON helper falls back when orjson is unavailable."""
        # Test helpers.json_dumps and json_loads when orjson is missing
        with patch("gopptx.presentation.helpers._orjson", None):
            data = {"a": 1}
            dumped = helpers.json_dumps(data)
            self.assertEqual(dumped, b'{"a":1}')
            loaded = helpers.json_loads(dumped)
            self.assertEqual(loaded, data)


if __name__ == "__main__":
    unittest.main()
