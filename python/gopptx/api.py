from __future__ import annotations

import ctypes
import json
import os
import sys
import uuid
from typing import Any, Dict, Optional, cast

from . import ops
from .types import (
    Author,
    ChartDataUpdate,
    ChartSelector,
    Comment,
    CoreProperties,
    PresentationMetadata,
    Section,
    Shape,
    ShapeProps,
    ShapeSearchQuery,
    ShapeSearchResult,
    ShapeUpdate,
    SlideChartRef,
    SlideLayoutInfo,
    SlideMasterCloneResult,
    SlideMetadata,
)


class GopptxError(Exception):
    """Base exception for gopptx library errors."""

    def __init__(self, message: str, code: Optional[str] = None):
        super().__init__(message)
        self.code = code


class Slide:
    """Proxy object for a slide within a presentation."""

    def __init__(self, presentation: Presentation, metadata: SlideMetadata):
        self._presentation = presentation
        self._metadata = metadata

    @property
    def index(self) -> int:
        """The zero-based index of this slide."""
        # Try to find current index by ID to handle removals/moves
        for s in self._presentation.slides_metadata:
            if s["SlideID"] == self.slide_id:
                return int(s["Index"])
        return -1

    @property
    def slide_id(self) -> int:
        """The unique internal ID of this slide."""
        return self._metadata["SlideID"]

    @property
    def title(self) -> str:
        """The title text of this slide."""
        return self._metadata["Title"]

    @title.setter
    def title(self, value: str) -> None:
        """Updates the slide title."""
        self._presentation.set_slide_title(self.index, value)
        self._metadata["Title"] = value

    @property
    def notes(self) -> str:
        """The speaker notes for this slide."""
        return self._presentation.get_notes(self.index)

    @notes.setter
    def notes(self, value: str) -> None:
        """Updates the speaker notes."""
        self._presentation.set_notes(self.index, value)

    def add_shape(
        self,
        shape_type: str,
        x: float,
        y: float,
        w: float,
        h: float,
        text: Optional[str] = None,
        properties: Optional[ShapeProps] = None,
    ) -> int:
        """Adds a new shape to this slide."""
        return self._presentation.add_shape(self.index, shape_type, x, y, w, h, text=text, properties=properties)

    def add_image(self, path: str, x: float, y: float, w: float, h: float) -> int:
        """Adds an image to this slide from a local file path."""
        return self._presentation.add_image(self.index, path, x, y, w, h)

    def remove_shape(self, shape_id: int) -> None:
        """Removes a shape from this slide."""
        self._presentation.remove_shape(self.index, shape_id)

    def update_shape(self, shape_id: int, updates: ShapeUpdate) -> None:
        """Updates properties of an existing shape on this slide."""
        self._presentation.update_shape(self.index, shape_id, updates)

    def list_shapes(self) -> list[Shape]:
        """Returns a list of all shapes on this slide."""
        return self._presentation.list_shapes(self.index)

    def list_charts(self) -> list[SlideChartRef]:
        """Returns a list of all charts on this slide."""
        return self._presentation.list_slide_charts(self.index)

    def add_chart(
        self,
        chart_type: str,
        categories: list[str],
        values: list[float],
        title: str = "Chart",
        x: int = 0,
        y: int = 0,
        w: int = 0,
        h: int = 0,
    ) -> None:
        """Adds a chart to this slide."""
        self._presentation.add_chart(self.index, chart_type, categories, values, title=title, x=x, y=y, w=w, h=h)

    def update(self, title: Optional[str] = None, layout: Optional[str] = None, bullets: Optional[list[str]] = None) -> None:
        """Updates this slide's content."""
        self._presentation.update_slide(self.index, title=title, layout=layout, bullets=bullets)
        if title:
            self._metadata["Title"] = title

    def remove(self) -> None:
        """Removes this slide from the presentation."""
        self._presentation.remove_slide(self.index)

    def duplicate(self, insert_at: Optional[int] = None) -> Slide:
        """Duplicates this slide and returns the new Slide object."""
        new_idx = self._presentation.duplicate_slide(self.index, insert_at=insert_at)
        return self._presentation.slides[new_idx]

    def __repr__(self) -> str:
        return f"<Slide index={self.index} title='{self.title}'>"


class Presentation:
    """High-level wrapper for a PowerPoint presentation handled by the Go engine."""

    _lib = None

    def __init__(self, path: Optional[str] = None):
        self._load_library()
        self._handle: Optional[int] = None
        if path:
            self.open(path)

    @classmethod
    def _load_library(cls) -> None:
        if cls._lib:
            return

        pkg_dir = os.path.dirname(__file__)
        if sys.platform == "win32":
            lib_name = "gopptx.dll"
        elif sys.platform == "darwin":
            lib_name = "libgopptx.dylib"
        else:
            lib_name = "libgopptx.so"

        search_paths = [
            os.path.join(pkg_dir, lib_name),
            os.path.join(pkg_dir, "../../bindings/c/build", lib_name),
        ]

        lib_path = None
        for candidate in search_paths:
            if os.path.exists(candidate):
                lib_path = os.path.abspath(candidate)
                break

        if not lib_path:
            raise GopptxError(
                f"Could not find shared library {lib_name}. Please build it first."
            )

        cls._lib = ctypes.CDLL(lib_path)
        cls._lib.deck_open.argtypes = [ctypes.c_char_p]
        cls._lib.deck_open.restype = ctypes.c_void_p

        cls._lib.deck_new.argtypes = [ctypes.c_char_p]
        cls._lib.deck_new.restype = ctypes.c_void_p

        cls._lib.deck_execute_json.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
        cls._lib.deck_execute_json.restype = ctypes.c_void_p

        cls._lib.deck_save.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
        cls._lib.deck_save.restype = ctypes.c_int

        cls._lib.deck_last_error.argtypes = [ctypes.c_void_p]
        cls._lib.deck_last_error.restype = ctypes.c_void_p

        cls._lib.deck_global_error.argtypes = []
        cls._lib.deck_global_error.restype = ctypes.c_void_p

        cls._lib.deck_free_string.argtypes = [ctypes.c_void_p]
        cls._lib.deck_free_string.restype = None

        cls._lib.deck_close.argtypes = [ctypes.c_void_p]
        cls._lib.deck_close.restype = None

    @classmethod
    def new(cls, title: str) -> Presentation:
        """Create a new presentation with the given title."""
        pres = cls()
        handle = cls._lib.deck_new(title.encode("utf-8"))
        if not handle:
            err_ptr = cls._lib.deck_global_error()
            msg = ctypes.string_at(err_ptr).decode("utf-8") if err_ptr else "Unknown error"
            if err_ptr:
                cls._lib.deck_free_string(err_ptr)
            raise GopptxError(f"Failed to create new deck: {msg}")
        pres._handle = int(handle)
        return pres

    def execute(self, op: str, payload: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        """Execute a command operation against the Go engine."""
        if not self._handle:
            raise GopptxError("Presentation is not open.")

        envelope = {
            "api_version": 1,
            "request_id": str(uuid.uuid4()),
            "op": op,
            "payload": payload or {},
        }
        json_input = json.dumps(envelope).encode("utf-8")
        res_ptr = self._lib.deck_execute_json(self._handle, json_input)
        if not res_ptr:
            raise GopptxError("Received null response from Go engine")

        try:
            res_str = ctypes.string_at(res_ptr).decode("utf-8")
            response = json.loads(res_str)
            if not response.get("ok", False):
                err = response.get("error", {})
                raise GopptxError(
                    err.get("message", "Unknown engine error"),
                    code=err.get("code"),
                )
            result = response.get("result")
            if result is None:
                return {}
            if not isinstance(result, dict):
                raise GopptxError("Invalid response payload type")
            return cast(Dict[str, Any], result)
        finally:
            self._lib.deck_free_string(res_ptr)

    @property
    def slide_count(self) -> int:
        """Returns the total number of slides in the presentation."""
        result = self.execute(ops.OP_SLIDE_COUNT, {})
        return int(result.get("count", 0))

    @property
    def metadata(self) -> PresentationMetadata:
        """Returns presentation metadata as a dictionary with keys 'title', 'slide_count', and 'size'."""
        result = self.execute(ops.OP_GET_METADATA, {})
        return cast(PresentationMetadata, result)

    @property
    def slides(self) -> list[Slide]:
        """Returns a list of all slides in the presentation as Slide objects."""
        return [Slide(self, m) for m in self.slides_metadata]

    @property
    def slides_metadata(self) -> list[SlideMetadata]:
        """Internal: returns the raw slide metadata from the engine."""
        result = self.execute(ops.OP_LIST_SLIDES, {})
        return cast(list[SlideMetadata], result.get("slides", []))

    @property
    def sections(self) -> list[Section]:
        """Returns a list of all grouped sections in the presentation."""
        result = self.execute(ops.OP_GET_SECTIONS, {})
        return cast(list[Section], result.get("sections", []))

    def add_slide(self, title: str, layout: Optional[str] = None, bullets: Optional[list[str]] = None) -> Slide:
        """Adds a new slide and returns the new Slide object."""
        payload = {"title": title}
        if layout:
            payload["layout"] = layout
        if bullets:
            payload["bullets"] = bullets
        result = self.execute(ops.OP_ADD_SLIDE, payload)
        idx = int(result.get("index", -1))
        return self.slides[idx]

    def remove_slide(self, index: int) -> None:
        """Removes the slide at the specified zero-based index."""
        result = self.execute(ops.OP_REMOVE_SLIDE, {"index": index})
        return None

    def move_slide(self, from_index: int, to_index: int) -> None:
        """Moves a slide from one position to another."""
        result = self.execute(ops.OP_MOVE_SLIDE, {"from": from_index, "to": to_index})
        return None

    def duplicate_slide(self, index: int, insert_at: Optional[int] = None) -> int:
        """Duplicates a slide and returns the new slide index."""
        if insert_at is None:
            insert_at = index + 1
        result = self.execute(ops.OP_DUPLICATE_SLIDE, {"index": index, "insert_at": insert_at})
        return int(result.get("new_index", -1))

    def update_slide(
        self,
        index: int,
        title: Optional[str] = None,
        layout: Optional[str] = None,
        bullets: Optional[list[str]] = None,
    ) -> None:
        """Updates a slide's content."""
        payload = {"slide_index": index}
        if title is not None:
            payload["title"] = title
        if layout is not None:
            payload["layout"] = layout
        if bullets is not None:
            payload["bullets"] = bullets
        self.execute(ops.OP_UPDATE_SLIDE, payload)

    def add_chart(
        self,
        slide_index: int,
        chart_type: str,
        categories: list[str],
        values: list[float],
        title: str = "Chart",
        x: int = 0,
        y: int = 0,
        w: int = 0,
        h: int = 0,
    ) -> None:
        """Adds a chart to a slide."""
        self.execute(ops.OP_ADD_CHART, {
            "slide_index": slide_index,
            "chart_type": chart_type,
            "title": title,
            "categories": categories,
            "values": values,
            "x": x,
            "y": y,
            "w": w,
            "h": h,
        })

    def set_slide_title(self, index: int, title: str) -> None:
        """Updates the title of a specific slide."""
        self.execute(ops.OP_SET_SLIDE_TITLE, {"slide_index": index, "title": title})

    def merge_from_file(self, path: str) -> None:
        """Appends all slides from another PPTX file."""
        self.execute(ops.OP_MERGE_FROM_FILE, {"path": path})

    def find_and_replace(self, find_text: str, replace_text: str) -> int:
        """Finds and replaces text in all shapes across the entire presentation. Returns replacement count."""
        result = self.execute(ops.OP_FIND_AND_REPLACE, {"find": find_text, "replace": replace_text})
        return int(result.get("replacements", 0))

    def search_shapes(self, query: ShapeSearchQuery) -> list[ShapeSearchResult]:
        """Searches for shapes matching the given criteria across the presentation."""
        result = self.execute(ops.OP_SEARCH_SHAPES, cast(Dict[str, Any], query))
        return cast(list[ShapeSearchResult], result.get("results", []))

    def get_authors(self) -> list[Author]:
        """Returns a list of all comment authors in the presentation."""
        result = self.execute(ops.OP_GET_AUTHORS, {})
        return cast(list[Author], result.get("authors", []))

    def add_author(self, name: str, initials: str) -> int:
        """Registers a new comment author and returns their ID."""
        result = self.execute(ops.OP_ADD_AUTHOR, {"name": name, "initials": initials})
        return int(result.get("author_id", -1))

    def get_comments(self, slide_index: int) -> list[Comment]:
        """Returns all comments on the specified slide."""
        result = self.execute(ops.OP_GET_COMMENTS, {"slide_index": slide_index})
        return cast(list[Comment], result.get("comments", []))

    def add_comment(self, slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0) -> None:
        """Adds a comment to a slide."""
        self.execute(ops.OP_ADD_COMMENT, {
            "slide_index": slide_index,
            "author_id": author_id,
            "text": text,
            "x": x,
            "y": y,
        })

    def remove_comment(self, slide_index: int, author_id: int, author_index: int) -> None:
        """Removes a comment from a slide using author ID and their sequential comment index."""
        self.execute(ops.OP_REMOVE_COMMENT, {
            "slide_index": slide_index,
            "author_id": author_id,
            "author_index": author_index,
        })

    def add_section(self, name: str, slide_indices: list[int]) -> None:
        """Adds a new grouped section for slides."""
        self.execute(ops.OP_ADD_SECTION, {"name": name, "slide_indices": slide_indices})

    def remove_section(self, name: str) -> None:
        """Removes a section by name."""
        self.execute(ops.OP_REMOVE_SECTION, {"name": name})

    def rename_section(self, old_name: str, new_name: str) -> None:
        """Renames a section."""
        self.execute(ops.OP_RENAME_SECTION, {"old_name": old_name, "new_name": new_name})

    @property
    def core_properties(self) -> CoreProperties:
        """Returns the presentation's core properties (Dublin Core metadata)."""
        result = self.execute(ops.OP_GET_CORE_PROPERTIES, {})
        return cast(CoreProperties, result)

    @core_properties.setter
    def core_properties(self, props: CoreProperties) -> None:
        """Updates the presentation's core properties."""
        self.execute(ops.OP_SET_CORE_PROPERTIES, props)

    @property
    def title(self) -> str:
        """The presentation title from core properties."""
        return self.core_properties.get("title", "")

    @title.setter
    def title(self, value: str) -> None:
        """Updates the presentation title."""
        props = self.core_properties
        props["title"] = value
        self.core_properties = props

    def add_title_slide(self, title: str) -> Slide:
        """Convenience: adds a slide with a title only."""
        return self.add_slide(title, layout="title_only")

    def add_bullet_slide(self, title: str, bullets: list[str]) -> Slide:
        """Convenience: adds a slide with title and bullet points."""
        return self.add_slide(title, layout="title_and_content", bullets=bullets)

    def __getitem__(self, index: int) -> Slide:
        """Access a slide by its zero-based index."""
        return self.slides[index]

    def __len__(self) -> int:
        """Returns the number of slides."""
        return self.slide_count

    def __iter__(self):
        """Allows iterating over slides."""
        return iter(self.slides)

    def apply_theme(self, theme_name: str) -> None:
        """Applies a theme preset to the presentation."""
        self.execute(ops.OP_APPLY_THEME, {"theme_name": theme_name})

    def set_slide_size(self, width: int, height: int) -> None:
        """Sets the slide size in EMUs."""
        self.execute(ops.OP_SET_SLIDE_SIZE, {"width": width, "height": height})

    def list_slide_layouts(self) -> list[SlideLayoutInfo]:
        """Returns a list of all available slide layouts in the presentation."""
        result = self.execute(ops.OP_LIST_SLIDE_LAYOUTS, {})
        return cast(list[SlideLayoutInfo], result.get("layouts", []))

    def rebind_slide_layout(self, slide_index: int, layout_part: str) -> None:
        """Changes the layout of a specific slide."""
        self.execute(ops.OP_REBIND_SLIDE_LAYOUT, {"slide_index": slide_index, "layout_part": layout_part})

    def clone_layout_master_family(self, layout_part: str) -> SlideMasterCloneResult:
        """Clones a layout and its parent master/theme into the presentation."""
        result = self.execute(ops.OP_CLONE_LAYOUT_MASTER_FAMILY, {"layout_part": layout_part})
        return cast(SlideMasterCloneResult, result)

    def list_shapes(self, slide_index: int) -> list[Shape]:
        """Returns a list of all shapes on the specified slide."""
        result = self.execute(ops.OP_LIST_SHAPES, {"slide_index": slide_index})
        return cast(list[Shape], result.get("shapes", []))

    def add_shape(
        self,
        slide_index: int,
        shape_type: str,
        x: float,
        y: float,
        w: float,
        h: float,
        text: Optional[str] = None,
        properties: Optional[ShapeProps] = None,
    ) -> int:
        """Adds a new shape to a slide and returns its ID."""
        payload = {
            "slide_index": slide_index,
            "type": shape_type,
            "x": x,
            "y": y,
            "w": w,
            "h": h,
        }
        if text is not None:
            payload["text"] = text
        if properties is not None:
            payload["properties"] = properties
        result = self.execute(ops.OP_ADD_SHAPE, payload)
        return int(result.get("shape_id", -1))

    def add_image(self, slide_index: int, path: str, x: float, y: float, w: float, h: float) -> int:
        """Adds an image to a slide from a local file path and returns its ID."""
        result = self.execute(ops.OP_ADD_IMAGE, {
            "slide_index": slide_index,
            "path": path,
            "x": x,
            "y": y,
            "w": w,
            "h": h,
        })
        return int(result.get("shape_id", -1))

    def remove_shape(self, slide_index: int, shape_id: int) -> None:
        """Removes a shape from a slide."""
        self.execute(ops.OP_REMOVE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id})

    def update_shape(self, slide_index: int, shape_id: int, updates: ShapeUpdate) -> None:
        """Updates properties of an existing shape."""
        self.execute(ops.OP_UPDATE_SHAPE, {"slide_index": slide_index, "shape_id": shape_id, "updates": updates})

    def get_notes(self, slide_index: int) -> str:
        """Returns the speaker notes for a slide as plain text."""
        result = self.execute(ops.OP_GET_NOTES, {"slide_index": slide_index})
        return str(result.get("text", ""))

    def set_notes(self, slide_index: int, text: str) -> None:
        """Sets the speaker notes for a slide."""
        self.execute(ops.OP_SET_NOTES, {"slide_index": slide_index, "text": text})

    def list_slide_charts(self, slide_index: int) -> list[SlideChartRef]:
        """Returns a list of all charts on the specified slide."""
        result = self.execute(ops.OP_LIST_SLIDE_CHARTS, {"slide_index": slide_index})
        return cast(list[SlideChartRef], result.get("charts", []))

    def update_chart_data(self, slide_index: int, chart_selector: ChartSelector, data: ChartDataUpdate) -> None:
        """Updates the data of a chart on a slide."""
        self.execute(ops.OP_UPDATE_CHART_DATA, {"slide_index": slide_index, "chart_selector": chart_selector, "data": data})


    def _get_last_error(self) -> str:
        if not self._handle:
            return "No active session"
        err_ptr = self._lib.deck_last_error(self._handle)
        if err_ptr:
            err_msg = ctypes.string_at(err_ptr).decode("utf-8")
            self._lib.deck_free_string(err_ptr)
            return err_msg
        return "Unknown error"

    def open(self, path: str) -> None:
        """Open an existing PPTX file."""
        if self._handle:
            self.close()
        handle = self._lib.deck_open(path.encode("utf-8"))
        if not handle:
            err_ptr = self._lib.deck_global_error()
            msg = ctypes.string_at(err_ptr).decode("utf-8") if err_ptr else "Unknown error"
            if err_ptr:
                self._lib.deck_free_string(err_ptr)
            raise GopptxError(f"Failed to open deck: {msg}")
        self._handle = int(handle)

    def save(self, path: str) -> None:
        """Save the presentation to the specified path."""
        if not self._handle:
            raise GopptxError("Presentation is not open.")
        rc = self._lib.deck_save(self._handle, path.encode("utf-8"))
        if rc != 0:
            raise GopptxError(f"Failed to save deck: {self._get_last_error()}")

    def close(self) -> None:
        """Close the presentation and release resources."""
        if self._handle:
            self._lib.deck_close(self._handle)
            self._handle = None

    def __enter__(self) -> "Presentation":
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        self.close()

    def __del__(self) -> None:
        self.close()

    def __repr__(self) -> str:
        return f"<Presentation title='{self.title}' slides={self.slide_count}>"
