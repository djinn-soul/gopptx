from collections.abc import Mapping

from .schemas import Hyperlink

class RunHyperlink:
    address: str | None
    action: str | None
    tooltip: str | None
    target_slide: int | None
    jump: str | None
    macro: str | None
    history: bool | None
    highlight_click: bool | None
    end_sound: bool | None
    def __init__(
        self,
        *,
        address: str | None = None,
        action: str | None = None,
        tooltip: str | None = None,
        target_slide: int | None = None,
        jump: str | None = None,
        macro: str | None = None,
        history: bool | None = None,
        highlight_click: bool | None = None,
        end_sound: bool | None = None,
    ) -> None: ...
    def to_payload(self) -> dict[str, object]: ...

class Run:
    text: str
    bold: bool | None
    italic: bool | None
    underline: str | None
    strikethrough: bool | None
    subscript: bool | None
    superscript: bool | None
    color: str | None
    highlight: str | None
    font: str | None
    size_pt: int | None
    code: bool | None
    all_caps: bool | None
    small_caps: bool | None
    @property
    def hyperlink(self) -> Mapping[str, object] | RunHyperlink | None: ...
    @hyperlink.setter
    def hyperlink(self, value: Mapping[str, object] | RunHyperlink | None) -> None: ...
    @property
    def hover_action(self) -> Mapping[str, object] | RunHyperlink | None: ...
    @hover_action.setter
    def hover_action(
        self, value: Mapping[str, object] | RunHyperlink | None
    ) -> None: ...
    def __init__(
        self,
        text: str = "",
        *,
        bold: bool | None = None,
        italic: bool | None = None,
        underline: str | None = None,
        strikethrough: bool | None = None,
        subscript: bool | None = None,
        superscript: bool | None = None,
        color: str | None = None,
        highlight: str | None = None,
        font: str | None = None,
        size_pt: int | None = None,
        code: bool | None = None,
        all_caps: bool | None = None,
        small_caps: bool | None = None,
        hyperlink: Hyperlink | RunHyperlink | None = None,
        hover_action: Hyperlink | RunHyperlink | None = None,
    ) -> None: ...
    def to_payload(self) -> dict[str, object]: ...

class TextFrameProps:
    margin_top: int | None
    margin_bottom: int | None
    margin_left: int | None
    margin_right: int | None
    word_wrap: bool | None
    auto_fit: bool | None
    auto_fit_type: str | None
    vertical_align: str | None
    orientation: str | None
    columns: int | None
    rotation: float | None
    def __init__(
        self,
        *,
        margin_top: int | None = None,
        margin_bottom: int | None = None,
        margin_left: int | None = None,
        margin_right: int | None = None,
        word_wrap: bool | None = None,
        auto_fit: bool | None = None,
        auto_fit_type: str | None = None,
        vertical_align: str | None = None,
        orientation: str | None = None,
        columns: int | None = None,
        rotation: float | None = None,
        vertical_anchor: str | None = None,
        auto_size: str | None = None,
        text_direction: str | None = None,
        column_count: int | None = None,
        text_rotation: float | None = None,
    ) -> None: ...
    def to_payload(self) -> dict[str, object]: ...

class ParagraphProps:
    indent: int | None
    hanging: int | None
    tab_stops: list[int] | None
    alignment: str | None
    bullet_style: str | None
    bullet_char: str | None
    bullet_color: str | None
    bullet_size_pct: int | None
    level: int | None
    line_spacing_pct: int | None
    line_spacing_pts: int | None
    space_before_pts: int | None
    space_after_pts: int | None
    def __init__(
        self,
        *,
        indent: int | None = None,
        hanging: int | None = None,
        tab_stops: list[int] | None = None,
        left_margin: int | None = None,
        hanging_indent: int | None = None,
        alignment: str | None = None,
        bullet_style: str | None = None,
        bullet_char: str | None = None,
        bullet_color: str | None = None,
        bullet_size_pct: int | None = None,
        level: int | None = None,
        line_spacing_pct: int | None = None,
        line_spacing_pts: int | None = None,
        space_before_pts: int | None = None,
        space_after_pts: int | None = None,
        tabs: list[int] | None = None,
    ) -> None: ...
    def to_payload(self) -> dict[str, object]: ...
