"""Comment mixin for the Presentation API."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops
from ..helpers import PresentationMixinBase

if TYPE_CHECKING:
    from ...schemas import Author, Comment


class PresentationCommentMixin(PresentationMixinBase):
    """Mixin providing comment and author management methods."""

    def get_authors(self) -> list[Author]:
        """Get all comment authors in the presentation."""
        result = self.execute(ops.OP_GET_AUTHORS, {})
        return cast("list[Author]", result.get("authors", []))

    def add_author(self, name: str, initials: str) -> int:
        """Add a comment author to the presentation."""
        result = self.execute(ops.OP_ADD_AUTHOR, {"name": name, "initials": initials})
        author_id = result.get("author_id")
        if not isinstance(author_id, int):
            msg = "bridge response author_id must be an int"
            raise TypeError(msg)
        return author_id

    def get_comments(self, slide_index: int) -> list[Comment]:
        """Get all comments on a slide."""
        result = self.execute(ops.OP_GET_COMMENTS, {"slide_index": slide_index})
        raw_comments = result.get("comments")
        comments = cast(
            "list[Comment]", raw_comments if isinstance(raw_comments, list) else []
        )
        for item in cast("list[dict[str, object]]", comments):
            if "Index" in item and "index" not in item:
                item["index"] = item["Index"]
        return comments

    def add_comment(
        self, slide_index: int, author_id: int, text: str, x: int = 0, y: int = 0
    ) -> int:
        """Add a comment to a slide."""
        self.execute(
            ops.OP_ADD_COMMENT,
            {
                "slide_index": slide_index,
                "author_id": author_id,
                "text": text,
                "x": x,
                "y": y,
            },
        )
        comments = self.get_comments(slide_index)
        author_index = 0
        for item in reversed(cast("list[dict[str, object]]", comments)):
            c_author = item.get("AuthorID", item.get("author_id"))
            if c_author is not None and int(str(c_author)) == author_id:
                c_idx = item.get("Index", item.get("index", 0))
                author_index = int(str(c_idx)) if c_idx is not None else 0
                break
        self._comment_ref_cache[author_index] = (slide_index, author_id, author_index)
        return author_index

    def remove_comment(
        self,
        slide_index_or_index: int,
        author_id: int | None = None,
        author_index: int | None = None,
    ) -> None:
        """Remove a comment from a slide."""
        if author_id is None and author_index is None:
            ref = self._comment_ref_cache.get(slide_index_or_index)
            if ref is None:
                raise ValueError(
                    "unknown comment index; call remove_comment(slide_index, author_id, author_index)"
                )
            slide_index, author_id, author_index = ref
        else:
            slide_index = slide_index_or_index
            if author_id is None or author_index is None:
                raise TypeError(
                    "remove_comment requires either (comment_index) or (slide_index, author_id, author_index)"
                )
        self.execute(
            ops.OP_REMOVE_COMMENT,
            {
                "slide_index": slide_index,
                "author_id": author_id,
                "author_index": author_index,
            },
        )
