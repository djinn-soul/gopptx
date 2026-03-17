"""Operation constants shared by gopptx Python runtime."""

from __future__ import annotations

OP_ADD_ANIMATION = "add_animation"
OP_ADD_AUDIO = "add_audio"
OP_ADD_AUTHOR = "add_author"
OP_ADD_CHART = "add_chart"
OP_ADD_COMMENT = "add_comment"
OP_ADD_CONNECTOR = "add_connector"
OP_ADD_CONNECTORS = "add_connectors"
OP_ADD_CUSTOM_XML = "add_custom_xml"
OP_ADD_GROUP_SHAPE = "add_group_shape"
OP_ADD_IMAGE = "add_image"
OP_ADD_MERMAID_SHAPE = "add_mermaid_shape"
OP_ADD_OLE_OBJECT = "add_ole_object"
OP_ADD_SECTION = "add_section"
OP_ADD_SHAPE = "add_shape"
OP_ADD_SLIDE = "add_slide"
OP_ADD_SLIDE_LAYOUT = "add_slide_layout"
OP_ADD_SLIDE_MASTER = "add_slide_master"
OP_ADD_SMART_ART = "add_smartart"
OP_ADD_TABLE = "add_table"
OP_ADD_TEXTBOX = "add_textbox"
OP_ADD_TEXTBOXES = "add_textboxes"
OP_ADD_VBA = "add_vba"
OP_ADD_VIDEO = "add_video"
OP_APPEND_SHAPE_RUN = "append_shape_run"
OP_APPLY_THEME = "apply_theme"
OP_BATCH_EXECUTE = "batch_execute"
OP_BUILD_FREEFORM = "build_freeform"
OP_CLONE_LAYOUT_MASTER_FAMILY = "clone_layout_master_family"
OP_DEFINE_TABLE_STYLE = "define_table_style"
OP_DUPLICATE_SLIDE = "duplicate_slide"
OP_DUPLICATE_SLIDE_AFTER = "duplicate_slide_after"
OP_EXPORT_HTML = "export_html"
OP_EXPORT_PDF = "export_pdf"
OP_FIND_AND_REPLACE = "find_and_replace"
OP_GET_AUTHORS = "get_authors"
OP_GET_CHART_STATE = "get_chart_state"
OP_GET_COMMENTS = "get_comments"
OP_GET_CORE_PROPERTIES = "get_core_properties"
OP_GET_HANDOUT_MASTER = "get_handout_master"
OP_GET_IMAGE_METADATA = "get_image_metadata"
OP_GET_LAYOUT_PLACEHOLDERS = "get_layout_placeholders"
OP_GET_LAYOUT_SHAPES = "get_layout_shapes"
OP_GET_MASTER_PLACEHOLDERS = "get_master_placeholders"
OP_GET_MASTER_SHAPES = "get_master_shapes"
OP_GET_METADATA = "get_metadata"
OP_GET_NOTES = "get_notes"
OP_GET_SECTIONS = "get_sections"
OP_GET_SHAPE_RUNS = "get_shape_runs"
OP_GET_SHAPE_TEXT_STATE = "get_shape_text_state"
OP_GET_SLIDE_LAYOUT_REF = "get_slide_layout_ref"
OP_GET_SLIDE_TEXT_STATES = "get_slide_text_states"
OP_GET_TABLE = "get_table"
OP_GET_THEME_INVENTORY = "get_theme_inventory"
OP_GROUP_SHAPES = "group_shapes"
OP_HAS_DIGITAL_SIGNATURE = "has_digital_signature"
OP_LIST_CUSTOM_XML = "list_custom_xml"
OP_LIST_MASTER_LAYOUTS = "list_master_layouts"
OP_LIST_NOTES_PLACEHOLDERS = "list_notes_placeholders"
OP_LIST_NOTES_SHAPES = "list_notes_shapes"
OP_LIST_PLACEHOLDERS = "list_placeholders"
OP_LIST_SHAPES = "list_shapes"
OP_LIST_SLIDES = "list_slides"
OP_LIST_SLIDE_CHARTS = "list_slide_charts"
OP_LIST_SLIDE_IMAGES = "list_slide_images"
OP_LIST_SLIDE_LAYOUTS = "list_slide_layouts"
OP_LIST_SLIDE_MASTERS = "list_slide_masters"
OP_LIST_TABLE_STYLES = "list_table_styles"
OP_MARKDOWN_TO_SLIDES = "markdown_to_slides"
OP_MERGE_FROM_EDITOR = "merge_from_editor"
OP_MERGE_FROM_FILE = "merge_from_file"
OP_MERGE_TABLE_CELLS = "merge_table_cells"
OP_MOVE_SHAPE_TO_BACK = "move_shape_to_back"
OP_MOVE_SHAPE_TO_FRONT = "move_shape_to_front"
OP_MOVE_SHAPE_TO_INDEX = "move_shape_to_index"
OP_MOVE_SLIDE = "move_slide"
OP_NOTES_SLIDE_EXISTS = "notes_slide_exists"
OP_REBIND_SLIDE_LAYOUT = "rebind_slide_layout"
OP_REMOVE_COMMENT = "remove_comment"
OP_REMOVE_CUSTOM_XML = "remove_custom_xml"
OP_REMOVE_SECTION = "remove_section"
OP_REMOVE_SHAPE = "remove_shape"
OP_REMOVE_SLIDE = "remove_slide"
OP_REMOVE_SLIDE_LAYOUT = "remove_slide_layout"
OP_REMOVE_SLIDE_MASTER = "remove_slide_master"
OP_RENAME_SECTION = "rename_section"
OP_REPAIR = "repair"
OP_RESERVE_SHAPE_I_DS = "reserve_shape_ids"
OP_RESERVE_SHAPE_IDS = OP_RESERVE_SHAPE_I_DS
OP_SEARCH_SHAPES = "search_shapes"
OP_SET_CORE_PROPERTIES = "set_core_properties"
OP_SET_GLOBAL_THEME_PRESET = "set_global_theme_preset"
OP_SET_MARK_AS_FINAL = "set_mark_as_final"
OP_SET_MODIFY_PASSWORD = "set_modify_password"
OP_SET_NOTES = "set_notes"
OP_SET_NOTES_SHAPE_PROPS = "set_notes_shape_props"
OP_SET_NOTES_SHAPE_TEXT = "set_notes_shape_text"
OP_SET_PLACEHOLDER_CONTENT = "set_placeholder_content"
OP_SET_SHAPE_RUNS = "set_shape_runs"
OP_SET_SLIDE_BACKGROUND = "set_slide_background"
OP_SET_SLIDE_HEADER_FOOTER = "set_slide_header_footer"
OP_SET_SLIDE_SHAPE_RUNS = "set_slide_shape_runs"
OP_SET_SLIDE_SIZE = "set_slide_size"
OP_SET_SLIDE_TITLE = "set_slide_title"
OP_SET_SLIDE_TRANSITION = "set_slide_transition"
OP_SET_TABLE_COLUMN_WIDTH = "set_table_column_width"
OP_SET_TABLE_ROW_HEIGHT = "set_table_row_height"
OP_SET_TABLE_STYLE = "set_table_style"
OP_SET_THEME_COLOR_SCHEME = "set_theme_color_scheme"
OP_SET_THEME_FONT_SCHEME = "set_theme_font_scheme"
OP_SLIDE_COUNT = "slide_count"
OP_SPLIT_TABLE_CELL = "split_table_cell"
OP_SWAP_IMAGE_BY_INDEX = "swap_image_by_index"
OP_SWAP_IMAGE_BY_REL_ID = "swap_image_by_rel_id"
OP_UNGROUP_SHAPES = "ungroup_shapes"
OP_UPDATE_CHART_DATA = "update_chart_data"
OP_UPDATE_CHART_DATA_BATCH = "update_chart_data_batch"
OP_UPDATE_CHART_FORMATTING = "update_chart_formatting"
OP_UPDATE_DECK_RUN_TEXTS = "update_deck_run_texts"
OP_UPDATE_HANDOUT_MASTER = "update_handout_master"
OP_UPDATE_NOTES_MASTER = "update_notes_master"
OP_UPDATE_SHAPE = "update_shape"
OP_UPDATE_SHAPE_RUN_TEXT = "update_shape_run_text"
OP_UPDATE_SLIDE = "update_slide"
OP_UPDATE_SLIDE_RUN_TEXTS = "update_slide_run_texts"
OP_UPDATE_SMART_ART = "update_smartart"
OP_UPDATE_TABLE_CELL = "update_table_cell"
OP_UPDATE_TABLE_FLAGS = "update_table_flags"
OP_URL_FETCH_TO_SLIDES = "url_fetch_to_slides"
OP_VALIDATE = "validate"

SUPPORTED_OPS = (
    OP_ADD_ANIMATION,
    OP_ADD_AUDIO,
    OP_ADD_AUTHOR,
    OP_ADD_CHART,
    OP_ADD_COMMENT,
    OP_ADD_CONNECTOR,
    OP_ADD_CONNECTORS,
    OP_ADD_CUSTOM_XML,
    OP_ADD_GROUP_SHAPE,
    OP_ADD_IMAGE,
    OP_ADD_MERMAID_SHAPE,
    OP_ADD_OLE_OBJECT,
    OP_ADD_SECTION,
    OP_ADD_SHAPE,
    OP_ADD_SLIDE,
    OP_ADD_SLIDE_LAYOUT,
    OP_ADD_SLIDE_MASTER,
    OP_ADD_SMART_ART,
    OP_ADD_TABLE,
    OP_ADD_TEXTBOX,
    OP_ADD_TEXTBOXES,
    OP_ADD_VBA,
    OP_ADD_VIDEO,
    OP_APPEND_SHAPE_RUN,
    OP_APPLY_THEME,
    OP_BATCH_EXECUTE,
    OP_BUILD_FREEFORM,
    OP_CLONE_LAYOUT_MASTER_FAMILY,
    OP_DEFINE_TABLE_STYLE,
    OP_DUPLICATE_SLIDE,
    OP_DUPLICATE_SLIDE_AFTER,
    OP_EXPORT_HTML,
    OP_EXPORT_PDF,
    OP_FIND_AND_REPLACE,
    OP_GET_AUTHORS,
    OP_GET_CHART_STATE,
    OP_GET_COMMENTS,
    OP_GET_CORE_PROPERTIES,
    OP_GET_HANDOUT_MASTER,
    OP_GET_IMAGE_METADATA,
    OP_GET_LAYOUT_PLACEHOLDERS,
    OP_GET_LAYOUT_SHAPES,
    OP_GET_MASTER_PLACEHOLDERS,
    OP_GET_MASTER_SHAPES,
    OP_GET_METADATA,
    OP_GET_NOTES,
    OP_GET_SECTIONS,
    OP_GET_SHAPE_RUNS,
    OP_GET_SHAPE_TEXT_STATE,
    OP_GET_SLIDE_LAYOUT_REF,
    OP_GET_SLIDE_TEXT_STATES,
    OP_GET_TABLE,
    OP_GET_THEME_INVENTORY,
    OP_GROUP_SHAPES,
    OP_HAS_DIGITAL_SIGNATURE,
    OP_LIST_CUSTOM_XML,
    OP_LIST_MASTER_LAYOUTS,
    OP_LIST_NOTES_PLACEHOLDERS,
    OP_LIST_NOTES_SHAPES,
    OP_LIST_PLACEHOLDERS,
    OP_LIST_SHAPES,
    OP_LIST_SLIDES,
    OP_LIST_SLIDE_CHARTS,
    OP_LIST_SLIDE_IMAGES,
    OP_LIST_SLIDE_LAYOUTS,
    OP_LIST_SLIDE_MASTERS,
    OP_LIST_TABLE_STYLES,
    OP_MARKDOWN_TO_SLIDES,
    OP_MERGE_FROM_EDITOR,
    OP_MERGE_FROM_FILE,
    OP_MERGE_TABLE_CELLS,
    OP_MOVE_SHAPE_TO_BACK,
    OP_MOVE_SHAPE_TO_FRONT,
    OP_MOVE_SHAPE_TO_INDEX,
    OP_MOVE_SLIDE,
    OP_NOTES_SLIDE_EXISTS,
    OP_REBIND_SLIDE_LAYOUT,
    OP_REMOVE_COMMENT,
    OP_REMOVE_CUSTOM_XML,
    OP_REMOVE_SECTION,
    OP_REMOVE_SHAPE,
    OP_REMOVE_SLIDE,
    OP_REMOVE_SLIDE_LAYOUT,
    OP_REMOVE_SLIDE_MASTER,
    OP_RENAME_SECTION,
    OP_REPAIR,
    OP_RESERVE_SHAPE_I_DS,
    OP_SEARCH_SHAPES,
    OP_SET_CORE_PROPERTIES,
    OP_SET_GLOBAL_THEME_PRESET,
    OP_SET_MARK_AS_FINAL,
    OP_SET_MODIFY_PASSWORD,
    OP_SET_NOTES,
    OP_SET_NOTES_SHAPE_PROPS,
    OP_SET_NOTES_SHAPE_TEXT,
    OP_SET_PLACEHOLDER_CONTENT,
    OP_SET_SHAPE_RUNS,
    OP_SET_SLIDE_BACKGROUND,
    OP_SET_SLIDE_HEADER_FOOTER,
    OP_SET_SLIDE_SHAPE_RUNS,
    OP_SET_SLIDE_SIZE,
    OP_SET_SLIDE_TITLE,
    OP_SET_SLIDE_TRANSITION,
    OP_SET_TABLE_COLUMN_WIDTH,
    OP_SET_TABLE_ROW_HEIGHT,
    OP_SET_TABLE_STYLE,
    OP_SET_THEME_COLOR_SCHEME,
    OP_SET_THEME_FONT_SCHEME,
    OP_SLIDE_COUNT,
    OP_SPLIT_TABLE_CELL,
    OP_SWAP_IMAGE_BY_INDEX,
    OP_SWAP_IMAGE_BY_REL_ID,
    OP_UNGROUP_SHAPES,
    OP_UPDATE_CHART_DATA,
    OP_UPDATE_CHART_DATA_BATCH,
    OP_UPDATE_CHART_FORMATTING,
    OP_UPDATE_DECK_RUN_TEXTS,
    OP_UPDATE_HANDOUT_MASTER,
    OP_UPDATE_NOTES_MASTER,
    OP_UPDATE_SHAPE,
    OP_UPDATE_SHAPE_RUN_TEXT,
    OP_UPDATE_SLIDE,
    OP_UPDATE_SLIDE_RUN_TEXTS,
    OP_UPDATE_SMART_ART,
    OP_UPDATE_TABLE_CELL,
    OP_UPDATE_TABLE_FLAGS,
    OP_URL_FETCH_TO_SLIDES,
    OP_VALIDATE,
)

SUPPORTED_OPS_SET = frozenset(SUPPORTED_OPS)
