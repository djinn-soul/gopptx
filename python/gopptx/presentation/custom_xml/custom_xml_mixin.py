"""Custom XML mixin for gopptx library."""

from __future__ import annotations

from typing import cast

from ... import ops
from ..helpers import PresentationMixinBase


class PresentationCustomXMLMixin(PresentationMixinBase):
    """Mixin providing custom XML part manipulation methods."""

    def add_custom_xml(
        self,
        content: str,
        root_element: str | None = None,
        namespace: str | None = None,
        properties: dict[str, str] | None = None,
    ) -> int:
        """Embed a custom XML part in the presentation.

        Args:
            content: Raw XML string to embed.
            root_element: Optional root element name for the XML part.
            namespace: Optional XML namespace URI.
            properties: Optional key/value metadata properties.

        Returns:
            Index of the newly added custom XML part.
        """
        payload: dict[str, object] = {"content": content}
        if root_element is not None:
            payload["root_element"] = root_element
        if namespace is not None:
            payload["namespace"] = namespace
        if properties is not None:
            payload["properties"] = properties
        result = self.execute(ops.OP_ADD_CUSTOM_XML, payload)
        return int(cast("int", result.get("index", 0)))

    def list_custom_xml(self) -> list[dict[str, object]]:
        """Return all custom XML parts embedded in the presentation.

        Each entry contains: content, root_element, namespace, properties.
        """
        result = self.execute(ops.OP_LIST_CUSTOM_XML, {})
        items = result.get("custom_xml", [])
        return cast("list[dict[str, object]]", items) if isinstance(items, list) else []

    def remove_custom_xml(self, index: int) -> None:
        """Remove a custom XML part by its index.

        Args:
            index: Zero-based index of the custom XML part to remove.
        """
        self.execute(ops.OP_REMOVE_CUSTOM_XML, {"index": index})
