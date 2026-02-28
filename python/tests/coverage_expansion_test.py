import os
import unittest
from unittest.mock import patch, MagicMock
import pytest
from gopptx import Presentation, GopptxError, ops
from gopptx.utils import normalize_table_index
from gopptx.presentation import helpers

class TestCoverageExpansion(unittest.TestCase):
    def setUp(self):
        self.test_pptx = os.path.join("testdata", "ppt_rs", "simple.pptx")

    def test_runtime_unopened_errors(self):
        pres = Presentation()
        with self.assertRaisesRegex(GopptxError, "Presentation is not open."):
            pres.execute("any_op")
        with self.assertRaisesRegex(GopptxError, "Presentation is not open."):
            pres.save("any_path")
        
        # Test repr on unopened - it currently raises GopptxError which we should handle or expect
        with self.assertRaises(GopptxError):
            repr(pres)

    def test_batch_read_op_not_allowed(self):
        with Presentation(self.test_pptx) as pres:
            with pres.batch():
                with self.assertRaisesRegex(GopptxError, r"read operation 'slide_count' is not allowed inside batch"):
                    pres.execute(ops.OP_SLIDE_COUNT)

    def test_execute_batch_empty(self):
        with Presentation(self.test_pptx) as pres:
            self.assertEqual(pres.execute_batch([]), [])

    def test_nested_batch_not_allowed(self):
        with Presentation(self.test_pptx) as pres:
            pres.begin_batch()
            with self.assertRaisesRegex(GopptxError, r"nested batch\(\) calls are not allowed"):
                pres.begin_batch()
            pres.abort_batch()

    def test_abort_batch(self):
        with Presentation(self.test_pptx) as pres:
            # Get initial count
            initial_count = pres.slide_count
            pres.begin_batch()
            pres.add_slide("New Slide")
            pres.abort_batch()
            # Count should remain the same
            self.assertEqual(pres.slide_count, initial_count)

    def test_metadata_cache(self):
        with Presentation(self.test_pptx) as pres:
            m1 = pres.metadata
            m2 = pres.metadata
            self.assertIs(m1, m2)
            pres.invalidate_cache()
            m3 = pres.metadata
            self.assertIsNot(m1, m3)

    def test_presentation_getitem(self):
        with Presentation(self.test_pptx) as pres:
            count = pres.slide_count
            # Simple indexing
            s0 = pres[0]
            self.assertEqual(s0.index, 0)
            
            # Negative indexing
            s_last = pres[-1]
            self.assertEqual(s_last.index, count - 1)
            
            with self.assertRaises(IndexError):
                _ = pres[100]
            
            # Slicing
            slides = pres[0:1]
            self.assertEqual(len(slides), 1)
            
            with self.assertRaises(TypeError):
                _ = pres["invalid"]

    def test_reopen_already_open(self):
        with Presentation(self.test_pptx) as pres:
            pres.open(self.test_pptx) # Should close previous and open new
            self.assertGreater(pres.slide_count, 0)

    def test_properties_and_protection(self):
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

    def test_slide_management_expansion(self):
        with Presentation(self.test_pptx) as pres:
            initial_count = pres.slide_count
            # Add with layout
            s2 = pres.add_slide("Layout Slide", layout="Title Slide")
            self.assertEqual(s2.title, "Layout Slide")
            
            # Update with layout
            pres.update_slide(1, title="Updated Title", layout="Blank")
            
            # Add title/bullet slides
            pres.add_title_slide("Title Only")
            pres.add_bullet_slide("Bullets", ["B1", "B2"])
            self.assertEqual(pres.slide_count, initial_count + 3)
            
            # Len, Iter
            self.assertEqual(len(pres), initial_count + 3)
            slide_list = list(iter(pres))
            self.assertEqual(len(slide_list), initial_count + 3)

    def test_batch_mode_placeholder_slide(self):
        with Presentation(self.test_pptx) as pres:
            with pres.batch():
                s = pres.add_slide("Batch Slide")
                # Cannot check index in batch mode because it triggers a read op
                self.assertEqual(s.title, "Batch Slide")

    def test_layout_theme_expansion(self):
        with Presentation(self.test_pptx) as pres:
            layouts = pres.list_slide_layouts()
            if layouts:
                name = layouts[0]["name"]
                # Rebind by name
                pres.rebind_slide_layout(0, name)
            
            # Apply theme
            pres.apply_theme("Office") # Tests the "office" -> "Corporate" mapping
            # "Default" might not be supported by engine, skip it or use another known one
            
    def test_remove_comment_errors(self):
        with Presentation(self.test_pptx) as pres:
            with self.assertRaises(ValueError):
                pres.remove_comment(999) # Unknown index in cache
            
            with self.assertRaises(TypeError):
                pres.remove_comment(0, author_id=1) # Missing author_index

    def test_chart_expansion(self):
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
                    "series": [{"name": "S1", "values": [15.0, 25.0], "categories": ["C1", "C2"]}]
                }
                pres.update_chart_data(0, {"rel_id": rel_id}, data)

    def test_normalize_table_index(self):
        self.assertEqual(normalize_table_index(1.0), 1)
        with self.assertRaises(ValueError):
            normalize_table_index(1.5)
        with self.assertRaises(ValueError):
            normalize_table_index("not_an_int")
        with self.assertRaises(ValueError):
            normalize_table_index(True)

    def test_json_fallback(self):
        # Test helpers.json_dumps and json_loads when orjson is missing
        with patch("gopptx.presentation.helpers._orjson", None):
            data = {"a": 1}
            dumped = helpers.json_dumps(data)
            self.assertEqual(dumped, b'{"a":1}')
            loaded = helpers.json_loads(dumped)
            self.assertEqual(loaded, data)

if __name__ == "__main__":
    unittest.main()
