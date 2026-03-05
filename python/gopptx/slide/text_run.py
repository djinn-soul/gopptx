"""Run-level text facade helpers for python-pptx-style ergonomics."""

from __future__ import annotations

from collections.abc import Mapping
from typing import cast


class RunHyperlink:
    """Mutable hyperlink object used by `Run`."""

    __slots__ = (
        "action",
        "address",
        "end_sound",
        "highlight_click",
        "history",
        "jump",
        "macro",
        "target_slide",
        "tooltip",
    )

    def __init__(  # noqa: PLR0913, D107
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
    ) -> None:
        super().__init__()
        self.address = address
        self.action = action
        self.tooltip = tooltip
        self.target_slide = target_slide
        self.jump = jump
        self.macro = macro
        self.history = history
        self.highlight_click = highlight_click
        self.end_sound = end_sound

    @classmethod
    def from_payload(
        cls, payload: Mapping[str, object] | RunHyperlink | None
    ) -> RunHyperlink | None:
        """Build a hyperlink facade from a bridge payload."""
        if payload is None:
            return None
        if isinstance(payload, RunHyperlink):
            return payload
        return cls(
            address=_as_optional_string(payload.get("address")),
            action=_as_optional_string(payload.get("action")),
            tooltip=_as_optional_string(payload.get("tooltip")),
            target_slide=_as_optional_int(payload.get("target_slide")),
            jump=_as_optional_string(payload.get("jump")),
            macro=_as_optional_string(payload.get("macro")),
            history=_as_optional_bool(payload.get("history")),
            highlight_click=_as_optional_bool(payload.get("highlight_click")),
            end_sound=_as_optional_bool(payload.get("end_sound")),
        )

    def to_payload(self) -> dict[str, object]:
        """Convert this hyperlink facade to bridge payload format."""
        payload: dict[str, object] = {}
        if self.address is not None:
            payload["address"] = self.address
        if self.action is not None:
            payload["action"] = self.action
        if self.tooltip is not None:
            payload["tooltip"] = self.tooltip
        if self.target_slide is not None:
            payload["target_slide"] = self.target_slide
        if self.jump is not None:
            payload["jump"] = self.jump
        if self.macro is not None:
            payload["macro"] = self.macro
        if self.history is not None:
            payload["history"] = self.history
        if self.highlight_click is not None:
            payload["highlight_click"] = self.highlight_click
        if self.end_sound is not None:
            payload["end_sound"] = self.end_sound
        return payload


class Run:
    """Run-level facade supporting python-pptx-style `run.hyperlink.address`."""

    __slots__ = (
        "_hover_action",
        "_hyperlink",
        "all_caps",
        "bold",
        "code",
        "color",
        "font",
        "highlight",
        "italic",
        "size_pt",
        "small_caps",
        "strikethrough",
        "subscript",
        "superscript",
        "text",
        "underline",
    )

    def __init__(  # noqa: PLR0913, D107
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
        hyperlink: Mapping[str, object] | RunHyperlink | None = None,
        hover_action: Mapping[str, object] | RunHyperlink | None = None,
    ) -> None:
        super().__init__()
        self.text = text
        self.bold = bold
        self.italic = italic
        self.underline = underline
        self.strikethrough = strikethrough
        self.subscript = subscript
        self.superscript = superscript
        self.color = color
        self.highlight = highlight
        self.font = font
        self.size_pt = size_pt
        self.code = code
        self.all_caps = all_caps
        self.small_caps = small_caps
        self._hyperlink = RunHyperlink.from_payload(hyperlink)
        self._hover_action = RunHyperlink.from_payload(hover_action)

    @property
    def hyperlink(self) -> RunHyperlink | None:
        """Return optional run hyperlink."""
        return self._hyperlink

    @hyperlink.setter
    def hyperlink(self, value: RunHyperlink) -> None:
        self._hyperlink = RunHyperlink.from_payload(value)

    @property
    def hover_action(self) -> RunHyperlink | None:
        """Return the optional hover action."""
        return self._hover_action

    @hover_action.setter
    def hover_action(self, value: RunHyperlink | None) -> None:
        self._hover_action = RunHyperlink.from_payload(value)

    def to_payload(self) -> dict[str, object]:  # noqa: C901, PLR0912
        """Convert this run facade to bridge payload format."""
        payload: dict[str, object] = {"text": self.text}
        if self.bold is not None:
            payload["bold"] = self.bold
        if self.italic is not None:
            payload["italic"] = self.italic
        if self.underline is not None:
            payload["underline"] = self.underline
        if self.strikethrough is not None:
            payload["strikethrough"] = self.strikethrough
        if self.subscript is not None:
            payload["subscript"] = self.subscript
        if self.superscript is not None:
            payload["superscript"] = self.superscript
        if self.color is not None:
            payload["color"] = self.color
        if self.highlight is not None:
            payload["highlight"] = self.highlight
        if self.font is not None:
            payload["font"] = self.font
        if self.size_pt is not None:
            payload["size_pt"] = self.size_pt
        if self.code is not None:
            payload["code"] = self.code
        if self.all_caps is not None:
            payload["all_caps"] = self.all_caps
        if self.small_caps is not None:
            payload["small_caps"] = self.small_caps
        if self._hyperlink is not None:
            payload["hyperlink"] = self._hyperlink.to_payload()
        if self._hover_action is not None:
            payload["hover_action"] = self._hover_action.to_payload()
        return payload


def serialize_runs_for_payload(runs: object) -> object:
    """Serialize a run list containing dict runs and/or `Run` facades."""
    if not isinstance(runs, list):
        return runs
    serialized: list[object] = []
    for item in cast("list[object]", runs):
        if isinstance(item, Run):
            serialized.append(item.to_payload())
            continue
        if isinstance(item, Mapping):
            as_dict: dict[str, object] = dict(cast("Mapping[str, object]", item))
            hyperlink = as_dict.get("hyperlink")
            hover_action = as_dict.get("hover_action")
            if isinstance(hyperlink, RunHyperlink):
                as_dict["hyperlink"] = hyperlink.to_payload()
            if isinstance(hover_action, RunHyperlink):
                as_dict["hover_action"] = hover_action.to_payload()
            serialized.append(as_dict)
            continue
        serialized.append(item)
    return serialized


def _as_optional_string(value: object) -> str | None:
    if value is None:
        return None
    if isinstance(value, str):
        return value
    return str(value)


def _as_optional_int(value: object) -> int | None:
    if value is None:
        return None
    if isinstance(value, int):
        return value
    return None


def _as_optional_bool(value: object) -> bool | None:
    if value is None:
        return None
    if isinstance(value, bool):
        return value
    return None
