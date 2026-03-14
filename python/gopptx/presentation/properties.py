"""Presentation properties/protection mixin."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from .. import ops
from .helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ..schemas import CoreProperties


class PresentationPropertiesMixin(PresentationMixinBase):
    """Mixin providing document properties and protection methods."""

    @property
    def core_properties(self) -> CoreProperties:
        """Get the core properties of the presentation."""
        return cast("CoreProperties", self.execute(ops.OP_GET_CORE_PROPERTIES, {}))

    def get_core_properties(self) -> CoreProperties:
        """Get the core properties of the presentation."""
        return self.core_properties

    @core_properties.setter
    def core_properties(self, props: CoreProperties) -> None:
        self.execute(ops.OP_SET_CORE_PROPERTIES, cast("dict[str, object]", props))

    def set_core_properties(self, props: CoreProperties) -> None:
        """Set the core properties of the presentation."""
        self.core_properties = props

    @property
    def title(self) -> str:
        """The title of the presentation."""
        return self.core_properties.get("title", "")

    @title.setter
    def title(self, value: str) -> None:
        props = self.core_properties
        props["title"] = value
        self.core_properties = props

    def set_modify_password(self, password: str) -> None:
        """Set the modify password for the presentation."""
        self.execute(ops.OP_SET_MODIFY_PASSWORD, {"password": password})

    def set_mark_as_final(self, *, final: bool = True) -> None:
        """Mark the presentation as final."""
        self.execute(ops.OP_SET_MARK_AS_FINAL, {"final": final})

    @property
    def author(self) -> str:
        """The author/creator of the presentation (python-pptx: author)."""
        return self.core_properties.get("creator", "")

    @author.setter
    def author(self, value: str) -> None:
        props = self.core_properties
        props["creator"] = value
        self.core_properties = props

    @property
    def comments(self) -> str:
        """The comments/description of the presentation (python-pptx: comments)."""
        return self.core_properties.get("description", "")

    @comments.setter
    def comments(self, value: str) -> None:
        props = self.core_properties
        props["description"] = value
        self.core_properties = props

    @property
    def identifier(self) -> str:
        """The identifier of the presentation."""
        return self.core_properties.get("identifier", "")

    @identifier.setter
    def identifier(self, value: str) -> None:
        props = self.core_properties
        props["identifier"] = value
        self.core_properties = props

    @property
    def language(self) -> str:
        """The language of the presentation."""
        return self.core_properties.get("language", "")

    @language.setter
    def language(self, value: str) -> None:
        props = self.core_properties
        props["language"] = value
        self.core_properties = props

    @property
    def last_printed(self) -> str:
        """The last printed date of the presentation."""
        return self.core_properties.get("lastPrinted", "")

    @last_printed.setter
    def last_printed(self, value: str) -> None:
        props = self.core_properties
        props["lastPrinted"] = value
        self.core_properties = props

    @property
    def version(self) -> str:
        """The version of the presentation."""
        return self.core_properties.get("version", "")

    @version.setter
    def version(self, value: str) -> None:
        props = self.core_properties
        props["version"] = value
        self.core_properties = props
