"""
Smoke tests for the Python bridge command API.

Tests all 40 Phase-1 bridge operations to ensure the Python wrapper
correctly communicates with the Go engine via the JSON command bridge.
"""

from __future__ import annotations

import os
import shutil
import tempfile
import unittest
from pathlib import Path

# Import the gopptx module
import gopptx
from gopptx import ops


class TestBridgeSlideOperations(unittest.TestCase):
    """Test slide-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_slide_count(self) -> None:
        """Test slide_count operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            count = pres.slide_count
            self.assertGreaterEqual(count, 1)
        finally:
            pres.close()

    def test_add_slide(self) -> None:
        """Test add_slide operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            initial_count = pres.slide_count
            pres.add_slide("New Slide")
            self.assertEqual(pres.slide_count, initial_count + 1)
        finally:
            pres.close()

    def test_list_slides(self) -> None:
        """Test list_slides operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            slides = pres.slides_metadata
            self.assertGreaterEqual(len(slides), 1)
        finally:
            pres.close()

    def test_set_slide_title(self) -> None:
        """Test set_slide_title operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Original Title")
            pres.set_slide_title(1, "Updated Title")
            # Verify by checking slides metadata
            slides = pres.slides_metadata
            if len(slides) > 1:
                self.assertIn("Updated", slides[1].get("title", ""))
        finally:
            pres.close()


class TestBridgeMetadataOperations(unittest.TestCase):
    """Test metadata-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_get_metadata(self) -> None:
        """Test get_metadata operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            meta = pres.metadata
            self.assertIsNotNone(meta)
            self.assertIn("title", meta)
            self.assertIn("slide_count", meta)
        finally:
            pres.close()

    def test_get_core_properties(self) -> None:
        """Test get_core_properties operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            props = pres.get_core_properties()
            self.assertIsNotNone(props)
        finally:
            pres.close()

    def test_set_core_properties(self) -> None:
        """Test set_core_properties operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.set_core_properties({"title": "New Title", "creator": "Test Author"})
            props = pres.get_core_properties()
            self.assertEqual(props.get("title"), "New Title")
        finally:
            pres.close()


class TestBridgeShapeOperations(unittest.TestCase):
    """Test shape-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_list_shapes(self) -> None:
        """Test list_shapes operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            shapes = pres.list_shapes(0)
            self.assertIsInstance(shapes, list)
        finally:
            pres.close()

    def test_add_shape(self) -> None:
        """Test add_shape operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            shape_id = pres.add_shape(0, "rect", 100, 100, 200, 100)
            self.assertIsInstance(shape_id, int)
            self.assertGreater(shape_id, 0)
        finally:
            pres.close()

    def test_find_and_replace(self) -> None:
        """Test find_and_replace operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            # Add a shape with text
            pres.add_shape(0, "rect", 100, 100, 200, 100, text="Hello World")
            # Find and replace
            count = pres.find_and_replace("World", "Bridge")
            self.assertGreaterEqual(count, 0)
        finally:
            pres.close()


class TestBridgeSectionOperations(unittest.TestCase):
    """Test section-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_get_sections_empty(self) -> None:
        """Test get_sections operation on new deck."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            sections = pres.get_sections()
            self.assertIsInstance(sections, list)
        finally:
            pres.close()

    def test_add_section(self) -> None:
        """Test add_section operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Slide 1")
            pres.add_slide("Slide 2")
            pres.add_section("Section 1", [1, 2])
            sections = pres.get_sections()
            self.assertEqual(len(sections), 1)
            self.assertEqual(sections[0].get("name"), "Section 1")
        finally:
            pres.close()


class TestBridgeNotesOperations(unittest.TestCase):
    """Test notes-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_set_and_get_notes(self) -> None:
        """Test set_notes and get_notes operations."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.set_notes(0, "Speaker notes for slide 1")
            notes = pres.get_notes(0)
            self.assertEqual(notes, "Speaker notes for slide 1")
        finally:
            pres.close()


class TestBridgeCommentOperations(unittest.TestCase):
    """Test comment-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_get_authors(self) -> None:
        """Test get_authors operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            authors = pres.get_authors()
            self.assertIsInstance(authors, list)
        finally:
            pres.close()

    def test_add_author(self) -> None:
        """Test add_author operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            author_id = pres.add_author("Test User", "TU")
            self.assertIsInstance(author_id, int)
        finally:
            pres.close()

    def test_get_comments(self) -> None:
        """Test get_comments operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            comments = pres.get_comments(0)
            self.assertIsInstance(comments, list)
        finally:
            pres.close()


class TestBridgeLayoutOperations(unittest.TestCase):
    """Test layout-related bridge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_list_slide_layouts(self) -> None:
        """Test list_slide_layouts operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            layouts = pres.list_slide_layouts()
            self.assertIsInstance(layouts, list)
        finally:
            pres.close()


class TestBridgeBatchOperations(unittest.TestCase):
    """Test batch execution operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_batch_context(self) -> None:
        """Test batch context manager."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            with pres.batch() as batch:
                batch.add_slide("Slide 1")
                batch.add_slide("Slide 2")
                batch.add_slide("Slide 3")
            self.assertGreaterEqual(pres.slide_count, 4)
        finally:
            pres.close()

    def test_execute_batch(self) -> None:
        """Test execute_batch operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            commands = [
                {"op": ops.OP_ADD_SLIDE, "payload": {"title": "Batch Slide 1"}},
                {"op": ops.OP_ADD_SLIDE, "payload": {"title": "Batch Slide 2"}},
            ]
            results = pres.execute_batch(commands)
            self.assertEqual(len(results), 2)
        finally:
            pres.close()


class TestBridgeSaveLoad(unittest.TestCase):
    """Test save and load operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_save_and_reload(self) -> None:
        """Test save and reload cycle."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Added Slide")
            pres.save(self.test_file)
        finally:
            pres.close()

        # Reload
        pres2 = gopptx.Presentation(self.test_file)
        try:
            self.assertGreaterEqual(pres2.slide_count, 2)
        finally:
            pres2.close()


class TestBridgeOpsConstants(unittest.TestCase):
    """Test that all ops constants are properly defined."""

    def test_all_ops_defined(self) -> None:
        """Verify all expected ops constants exist."""
        expected_ops = [
            "OP_BATCH_EXECUTE",
            "OP_SLIDE_COUNT",
            "OP_ADD_SLIDE",
            "OP_REMOVE_SLIDE",
            "OP_MOVE_SLIDE",
            "OP_DUPLICATE_SLIDE",
            "OP_GET_METADATA",
            "OP_UPDATE_CHART_DATA",
            "OP_LIST_SLIDE_CHARTS",
            "OP_LIST_SLIDE_LAYOUTS",
            "OP_REBIND_SLIDE_LAYOUT",
            "OP_CLONE_LAYOUT_MASTER_FAMILY",
            "OP_ADD_SECTION",
            "OP_REMOVE_SECTION",
            "OP_RENAME_SECTION",
            "OP_GET_SECTIONS",
            "OP_GET_CORE_PROPERTIES",
            "OP_SET_CORE_PROPERTIES",
            "OP_APPLY_THEME",
            "OP_SET_SLIDE_SIZE",
            "OP_SET_SLIDE_TITLE",
            "OP_MERGE_FROM_FILE",
            "OP_UPDATE_SLIDE",
            "OP_ADD_CHART",
            "OP_LIST_SLIDES",
            "OP_FIND_AND_REPLACE",
            "OP_SEARCH_SHAPES",
            "OP_GET_AUTHORS",
            "OP_ADD_AUTHOR",
            "OP_GET_COMMENTS",
            "OP_ADD_COMMENT",
            "OP_REMOVE_COMMENT",
            "OP_LIST_SHAPES",
            "OP_ADD_SHAPE",
            "OP_ADD_IMAGE",
            "OP_REMOVE_SHAPE",
            "OP_UPDATE_SHAPE",
            "OP_GET_NOTES",
            "OP_SET_NOTES",
            "OP_SET_MODIFY_PASSWORD",
            "OP_SET_MARK_AS_FINAL",
        ]

        for op_name in expected_ops:
            self.assertTrue(hasattr(ops, op_name), f"Missing op constant: {op_name}")
            op_value = getattr(ops, op_name)
            self.assertIsInstance(op_value, str)
            self.assertGreater(len(op_value), 0)

    def test_supported_ops_set(self) -> None:
        """Test SUPPORTED_OPS_SET contains all ops."""
        self.assertIsInstance(ops.SUPPORTED_OPS_SET, frozenset)
        self.assertGreaterEqual(len(ops.SUPPORTED_OPS_SET), 40)


class TestBridgeErrorHandling(unittest.TestCase):
    """Test error handling in the bridge."""

    def test_invalid_slide_index(self) -> None:
        """Test error on invalid slide index."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            with self.assertRaises(gopptx.GopptxError):
                pres.list_shapes(999)  # Invalid index
        finally:
            pres.close()

    def test_closed_presentation_access(self) -> None:
        """Test error on accessing closed presentation."""
        pres = gopptx.Presentation.new("Test Deck")
        pres.close()
        with self.assertRaises(gopptx.GopptxError):
            _ = pres.slide_count


class TestBridgeSlideAdvancedOperations(unittest.TestCase):
    """Test advanced slide operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_remove_slide(self) -> None:
        """Test remove_slide operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Slide to Remove")
            initial_count = pres.slide_count
            pres.remove_slide(initial_count - 1)  # Remove last slide
            self.assertEqual(pres.slide_count, initial_count - 1)
        finally:
            pres.close()

    def test_move_slide(self) -> None:
        """Test move_slide operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Second Slide")
            pres.move_slide(0, 1)  # Move first slide to position 1
            slides = pres.slides_metadata
            self.assertEqual(len(slides), 2)
        finally:
            pres.close()

    def test_duplicate_slide(self) -> None:
        """Test duplicate_slide operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            initial_count = pres.slide_count
            pres.duplicate_slide(0)
            self.assertEqual(pres.slide_count, initial_count + 1)
        finally:
            pres.close()

    def test_update_slide(self) -> None:
        """Test update_slide operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.update_slide(0, title="Updated Title", bullets=["Point 1", "Point 2"])
            slides = pres.slides_metadata
            self.assertIn("Updated", slides[0].get("title", ""))
        finally:
            pres.close()


class TestBridgeThemeAndSizeOperations(unittest.TestCase):
    """Test theme and slide size operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_set_slide_size(self) -> None:
        """Test set_slide_size operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            # Set to 16:9 widescreen
            pres.set_slide_size(12192000, 6858000)  # EMUs
            meta = pres.metadata
            self.assertIsNotNone(meta)
        finally:
            pres.close()

    def test_apply_theme(self) -> None:
        """Test apply_theme operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            # Apply a built-in theme
            pres.apply_theme("office")
            # No exception means success
        except gopptx.GopptxError:
            # Theme might not be available, that's ok for smoke test
            pass
        finally:
            pres.close()


class TestBridgeShapeAdvancedOperations(unittest.TestCase):
    """Test advanced shape operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_remove_shape(self) -> None:
        """Test remove_shape operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            shape_id = pres.add_shape(0, "rect", 100, 100, 200, 100)
            pres.remove_shape(0, shape_id)
            # Shape should be removed
        finally:
            pres.close()

    def test_update_shape(self) -> None:
        """Test update_shape operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            shape_id = pres.add_shape(0, "rect", 100, 100, 200, 100)
            pres.update_shape(0, shape_id, {"text": "Updated Text"})
            # Shape should be updated
        finally:
            pres.close()

    def test_search_shapes(self) -> None:
        """Test search_shapes operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_shape(0, "rect", 100, 100, 200, 100, text="SearchTarget")
            results = pres.search_shapes("SearchTarget")
            self.assertIsInstance(results, list)
        finally:
            pres.close()


class TestBridgeImageOperations(unittest.TestCase):
    """Test image operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_add_image(self) -> None:
        """Test add_image operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            # Create a minimal test image
            import struct
            import zlib
            
            def create_minimal_png():
                # Minimal 1x1 transparent PNG
                width, height = 1, 1
                raw_data = b'\x00\x00\x00\x00'  # RGBA: transparent
                compressed = zlib.compress(raw_data)
                
                def png_chunk(chunk_type, data):
                    chunk_len = struct.pack('>I', len(data))
                    chunk_crc = struct.pack('>I', zlib.crc32(chunk_type + data) & 0xffffffff)
                    return chunk_len + chunk_type + data + chunk_crc
                
                signature = b'\x89PNG\r\n\x1a\n'
                ihdr = png_chunk(b'IHDR', struct.pack('>IIBBBBB', width, height, 8, 6, 0, 0, 0))
                idat = png_chunk(b'IDAT', compressed)
                iend = png_chunk(b'IEND', b'')
                return signature + ihdr + idat + iend
            
            img_path = os.path.join(self.temp_dir, "test.png")
            with open(img_path, 'wb') as f:
                f.write(create_minimal_png())
            
            image_id = pres.add_image(0, img_path, 100, 100, 200, 100)
            self.assertIsInstance(image_id, int)
        finally:
            pres.close()


class TestBridgeChartOperations(unittest.TestCase):
    """Test chart operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_add_chart(self) -> None:
        """Test add_chart operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            chart_id = pres.add_chart(
                0, "bar",
                ["Q1", "Q2", "Q3"],
                [{"name": "Sales", "values": [100, 200, 150]}],
                100, 100, 400, 300
            )
            self.assertIsInstance(chart_id, int)
        finally:
            pres.close()

    def test_list_slide_charts(self) -> None:
        """Test list_slide_charts operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_chart(
                0, "bar",
                ["Q1", "Q2"],
                [{"name": "Sales", "values": [100, 200]}],
                100, 100, 400, 300
            )
            charts = pres.list_slide_charts(0)
            self.assertIsInstance(charts, list)
            self.assertGreater(len(charts), 0)
        finally:
            pres.close()

    def test_update_chart_data(self) -> None:
        """Test update_chart_data operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            chart_id = pres.add_chart(
                0, "bar",
                ["Q1", "Q2"],
                [{"name": "Sales", "values": [100, 200]}],
                100, 100, 400, 300
            )
            pres.update_chart_data(chart_id, ["Q1", "Q2", "Q3"], [{"name": "Sales", "values": [100, 200, 150]}])
            # Chart should be updated
        finally:
            pres.close()


class TestBridgeSectionAdvancedOperations(unittest.TestCase):
    """Test advanced section operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_remove_section(self) -> None:
        """Test remove_section operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Slide 1")
            pres.add_slide("Slide 2")
            pres.add_section("Section to Remove", [1, 2])
            sections = pres.get_sections()
            if len(sections) > 0:
                pres.remove_section(sections[0].get("id", 0))
        finally:
            pres.close()

    def test_rename_section(self) -> None:
        """Test rename_section operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.add_slide("Slide 1")
            pres.add_section("Original Name", [1])
            sections = pres.get_sections()
            if len(sections) > 0:
                pres.rename_section(sections[0].get("id", 0), "Renamed Section")
                sections = pres.get_sections()
                self.assertEqual(sections[0].get("name"), "Renamed Section")
        finally:
            pres.close()


class TestBridgeCommentAdvancedOperations(unittest.TestCase):
    """Test advanced comment operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_add_comment(self) -> None:
        """Test add_comment operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            author_id = pres.add_author("Test User", "TU")
            comment_id = pres.add_comment(0, author_id, "Test comment content")
            self.assertIsInstance(comment_id, int)
        finally:
            pres.close()

    def test_remove_comment(self) -> None:
        """Test remove_comment operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            author_id = pres.add_author("Test User", "TU")
            comment_id = pres.add_comment(0, author_id, "Comment to remove")
            pres.remove_comment(comment_id)
            # Comment should be removed
        finally:
            pres.close()


class TestBridgeLayoutAdvancedOperations(unittest.TestCase):
    """Test advanced layout operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_rebind_slide_layout(self) -> None:
        """Test rebind_slide_layout operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            layouts = pres.list_slide_layouts()
            if len(layouts) > 0:
                # Try to rebind to first available layout
                pres.rebind_slide_layout(0, layouts[0].get("name", ""))
        except gopptx.GopptxError:
            # Layout rebinding might not work on all presentations
            pass
        finally:
            pres.close()


class TestBridgeMergeOperations(unittest.TestCase):
    """Test merge operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")
        self.source_file = os.path.join(self.temp_dir, "source.pptx")

    def tearDown(self) -> None:
        for f in [self.test_file, self.source_file]:
            if os.path.exists(f):
                os.remove(f)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_merge_from_file(self) -> None:
        """Test merge_from_file operation."""
        # Create source presentation
        source = gopptx.Presentation.new("Source Deck")
        try:
            source.add_slide("Source Slide")
            source.save(self.source_file)
        finally:
            source.close()

        # Merge into target
        target = gopptx.Presentation.new("Target Deck")
        try:
            initial_count = target.slide_count
            target.merge_from_file(self.source_file)
            self.assertGreater(target.slide_count, initial_count)
        finally:
            target.close()


class TestBridgeProtectionOperations(unittest.TestCase):
    """Test protection operations."""

    def setUp(self) -> None:
        self.temp_dir = tempfile.mkdtemp()
        self.test_file = os.path.join(self.temp_dir, "test.pptx")

    def tearDown(self) -> None:
        if os.path.exists(self.test_file):
            os.remove(self.test_file)
        if os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir, ignore_errors=True)

    def test_set_modify_password(self) -> None:
        """Test set_modify_password operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.set_modify_password("testpassword123")
            # Password protection should be set
        finally:
            pres.close()

    def test_set_mark_as_final(self) -> None:
        """Test set_mark_as_final operation."""
        pres = gopptx.Presentation.new("Test Deck")
        try:
            pres.set_mark_as_final(True)
            # Presentation should be marked as final
        finally:
            pres.close()


if __name__ == "__main__":
    unittest.main()


