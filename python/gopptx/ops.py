"""Operation constants shared by gopptx Python runtime."""

from __future__ import annotations

OP_ADD_AUTHOR = "add_author"
OP_ADD_CHART = "add_chart"
OP_ADD_COMMENT = "add_comment"
OP_ADD_CUSTOM_XML = "add_custom_xml"
OP_ADD_IMAGE = "add_image"
OP_ADD_SECTION = "add_section"
OP_ADD_SHAPE = "add_shape"
OP_ADD_SLIDE = "add_slide"
OP_ADD_TABLE = "add_table"
OP_ADD_VBA = "add_vba"
OP_APPLY_THEME = "apply_theme"
OP_BATCH_EXECUTE = "batch_execute"
OP_CLONE_LAYOUT_MASTER_FAMILY = "clone_layout_master_family"
OP_DUPLICATE_SLIDE = "duplicate_slide"
OP_FIND_AND_REPLACE = "find_and_replace"
OP_GET_AUTHORS = "get_authors"
OP_GET_COMMENTS = "get_comments"
OP_GET_CORE_PROPERTIES = "get_core_properties"
OP_GET_METADATA = "get_metadata"
OP_GET_NOTES = "get_notes"
OP_GET_SECTIONS = "get_sections"
OP_GET_TABLE = "get_table"
OP_LIST_CUSTOM_XML = "list_custom_xml"
OP_LIST_MASTER_LAYOUTS = "list_master_layouts"
OP_LIST_SHAPES = "list_shapes"
OP_LIST_SLIDES = "list_slides"
OP_LIST_SLIDE_CHARTS = "list_slide_charts"
OP_LIST_SLIDE_LAYOUTS = "list_slide_layouts"
OP_LIST_SLIDE_MASTERS = "list_slide_masters"
OP_MERGE_FROM_FILE = "merge_from_file"
OP_MERGE_TABLE_CELLS = "merge_table_cells"
OP_MOVE_SHAPE_TO_BACK = "move_shape_to_back"
OP_MOVE_SHAPE_TO_FRONT = "move_shape_to_front"
OP_MOVE_SLIDE = "move_slide"
OP_REBIND_SLIDE_LAYOUT = "rebind_slide_layout"
OP_REMOVE_COMMENT = "remove_comment"
OP_REMOVE_CUSTOM_XML = "remove_custom_xml"
OP_REMOVE_SECTION = "remove_section"
OP_REMOVE_SHAPE = "remove_shape"
OP_REMOVE_SLIDE = "remove_slide"
OP_RENAME_SECTION = "rename_section"
OP_SEARCH_SHAPES = "search_shapes"
OP_SET_CORE_PROPERTIES = "set_core_properties"
OP_SET_MARK_AS_FINAL = "set_mark_as_final"
OP_SET_MODIFY_PASSWORD = "set_modify_password"
OP_SET_NOTES = "set_notes"
OP_SET_SLIDE_SIZE = "set_slide_size"
OP_SET_SLIDE_TITLE = "set_slide_title"
OP_SLIDE_COUNT = "slide_count"
OP_SPLIT_TABLE_CELL = "split_table_cell"
OP_UPDATE_CHART_DATA = "update_chart_data"
OP_UPDATE_SHAPE = "update_shape"
OP_UPDATE_SLIDE = "update_slide"
OP_UPDATE_TABLE_CELL = "update_table_cell"
OP_UPDATE_TABLE_FLAGS = "update_table_flags"

SUPPORTED_OPS = (
    OP_ADD_AUTHOR,
    OP_ADD_CHART,
    OP_ADD_COMMENT,
    OP_ADD_CUSTOM_XML,
    OP_ADD_IMAGE,
    OP_ADD_SECTION,
    OP_ADD_SHAPE,
    OP_ADD_SLIDE,
    OP_ADD_TABLE,
    OP_ADD_VBA,
    OP_APPLY_THEME,
    OP_BATCH_EXECUTE,
    OP_CLONE_LAYOUT_MASTER_FAMILY,
    OP_DUPLICATE_SLIDE,
    OP_FIND_AND_REPLACE,
    OP_GET_AUTHORS,
    OP_GET_COMMENTS,
    OP_GET_CORE_PROPERTIES,
    OP_GET_METADATA,
    OP_GET_NOTES,
    OP_GET_SECTIONS,
    OP_GET_TABLE,
    OP_LIST_CUSTOM_XML,
    OP_LIST_MASTER_LAYOUTS,
    OP_LIST_SHAPES,
    OP_LIST_SLIDES,
    OP_LIST_SLIDE_CHARTS,
    OP_LIST_SLIDE_LAYOUTS,
    OP_LIST_SLIDE_MASTERS,
    OP_MERGE_FROM_FILE,
    OP_MERGE_TABLE_CELLS,
    OP_MOVE_SHAPE_TO_BACK,
    OP_MOVE_SHAPE_TO_FRONT,
    OP_MOVE_SLIDE,
    OP_REBIND_SLIDE_LAYOUT,
    OP_REMOVE_COMMENT,
    OP_REMOVE_CUSTOM_XML,
    OP_REMOVE_SECTION,
    OP_REMOVE_SHAPE,
    OP_REMOVE_SLIDE,
    OP_RENAME_SECTION,
    OP_SEARCH_SHAPES,
    OP_SET_CORE_PROPERTIES,
    OP_SET_MARK_AS_FINAL,
    OP_SET_MODIFY_PASSWORD,
    OP_SET_NOTES,
    OP_SET_SLIDE_SIZE,
    OP_SET_SLIDE_TITLE,
    OP_SLIDE_COUNT,
    OP_SPLIT_TABLE_CELL,
    OP_UPDATE_CHART_DATA,
    OP_UPDATE_SHAPE,
    OP_UPDATE_SLIDE,
    OP_UPDATE_TABLE_CELL,
    OP_UPDATE_TABLE_FLAGS,
)

SUPPORTED_OPS_SET = frozenset(SUPPORTED_OPS)
