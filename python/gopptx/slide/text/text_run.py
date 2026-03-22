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

    def __init__(self, **kwargs: object) -> None:
        """Initialize with optional hyperlink attributes."""
        super().__init__()
        self.address = _as_optional_string(kwargs.get("address"))
        self.action = _as_optional_string(kwargs.get("action"))
        self.tooltip = _as_optional_string(kwargs.get("tooltip"))
        self.target_slide = _as_optional_int(kwargs.get("target_slide"))
        self.jump = _as_optional_string(kwargs.get("jump"))
        self.macro = _as_optional_string(kwargs.get("macro"))
        self.history = _as_optional_bool(kwargs.get("history"))
        self.highlight_click = _as_optional_bool(kwargs.get("highlight_click"))
        self.end_sound = _as_optional_bool(kwargs.get("end_sound"))

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

    @property
    def is_empty(self) -> bool:
        """Return True if this hyperlink has no attributes set."""
        return not any(getattr(self, attr) is not None for attr in self.__slots__)


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

    def __init__(self, text: str = "", **kwargs: object) -> None:
        """Initialize with optional run formatting attributes."""
        super().__init__()
        self.text = text
        self.bold = _as_optional_bool(kwargs.get("bold"))
        self.italic = _as_optional_bool(kwargs.get("italic"))
        self.underline = _as_optional_string(kwargs.get("underline"))
        self.strikethrough = _as_optional_bool(kwargs.get("strikethrough"))
        self.subscript = _as_optional_bool(kwargs.get("subscript"))
        self.superscript = _as_optional_bool(kwargs.get("superscript"))
        self.color = _as_optional_string(kwargs.get("color"))
        self.highlight = _as_optional_string(kwargs.get("highlight"))
        self.font = _as_optional_string(kwargs.get("font"))
        self.size_pt = _as_optional_int(kwargs.get("size_pt"))
        self.code = _as_optional_bool(kwargs.get("code"))
        self.all_caps = _as_optional_bool(kwargs.get("all_caps"))
        self.small_caps = _as_optional_bool(kwargs.get("small_caps"))
        self._hyperlink = RunHyperlink.from_payload(
            cast("Mapping[str, object] | RunHyperlink | None", kwargs.get("hyperlink"))
        )
        self._hover_action = RunHyperlink.from_payload(
            cast("Mapping[str, object] | RunHyperlink | None", kwargs.get("hover_action"))
        )

    @property
    def hyperlink(self) -> RunHyperlink:
        """Return run hyperlink, lazily initializing an empty facade."""
        if self._hyperlink is None:
            self._hyperlink = RunHyperlink()
        return self._hyperlink

    @hyperlink.setter
    def hyperlink(self, value: Mapping[str, object] | RunHyperlink | None) -> None:  # type: ignore[reportPropertyTypeMismatch]
        self._hyperlink = RunHyperlink.from_payload(value)

    @property
    def hover_action(self) -> RunHyperlink | None:
        """Return the optional hover action."""
        return self._hover_action

    @hover_action.setter
    def hover_action(self, value: Mapping[str, object] | RunHyperlink | None) -> None:  # type: ignore[reportPropertyTypeMismatch]
        self._hover_action = RunHyperlink.from_payload(value)

    def to_payload(self) -> dict[str, object]:
        """Convert this run facade to bridge payload format."""
        payload: dict[str, object] = {"text": self.text}
        for key in (
            "bold", "italic", "underline", "strikethrough", "subscript",
            "superscript", "color", "highlight", "font", "size_pt",
            "code", "all_caps", "small_caps",
        ):
            val = getattr(self, key)
            if val is not None:
                payload[key] = val
        if self._hyperlink is not None and not self._hyperlink.is_empty:
            payload["hyperlink"] = self._hyperlink.to_payload()
        if self._hover_action is not None and not self._hover_action.is_empty:
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
