"""Convenience slide methods and iteration support for presentations."""

from __future__ import annotations

from typing import TYPE_CHECKING, cast

from ... import ops

if TYPE_CHECKING:
    from collections.abc import Iterator

    from ...slide.slide import Slide


class PresentationSlidesExtrasMixin:
    """Mixin providing convenience slide helpers and iteration."""

    def add_title_slide(self, title: str) -> Slide:
        """Add a title slide to the presentation."""
        return self.add_slide(title, layout="title_only")

    def add_bullet_slide(self, title: str, bullets: list[str]) -> Slide:
        """Add a bullet slide to the presentation."""
        return self.add_slide(title, layout="title_and_content", bullets=bullets)

    def add_paragraph_slide(
        self,
        title: str,
        paragraph: str,
        *,
        bounds: tuple[float, float, float, float] | None = None,
        layout: str | None = None,
    ) -> Slide:
        """Add a slide with one paragraph textbox using sensible default bounds."""
        slide = self.add_slide(title, layout=layout)
        slide.add_paragraph(paragraph, bounds=bounds)
        return slide

    def add_slide_from_markdown(self, markdown: str, *, layout: str = "") -> int:
        """Append slides generated from a Markdown string.

        Returns:
            Zero-based index of the first slide added, or -1 when none produced.
        """
        payload: dict[str, object] = {"markdown": markdown}
        if layout:
            payload["layout"] = layout
        result = self.execute(ops.OP_MARKDOWN_TO_SLIDES, payload)
        self.invalidate_cache()
        return int(cast("int", result.get("first_index", -1)))

    def add_slide_from_url(self, url: str, *, layout: str = "") -> int:
        """Fetch a web page and append slides generated from its content.

        Returns:
            Zero-based index of the first slide added.
        """
        payload: dict[str, object] = {"url": url}
        result = self.execute(ops.OP_URL_FETCH_TO_SLIDES, payload)
        self.invalidate_cache()
        _ = layout  # layout not yet forwarded for URL fetch
        return int(cast("int", result.get("first_index", -1)))

    def __getitem__(self, index: int | slice) -> Slide | list[Slide]:
        """Get a slide by index or slice."""
        return self.slides[index]

    def __len__(self) -> int:
        """Return the number of slides."""
        return self.slide_count

    def __iter__(self) -> Iterator[Slide]:
        """Iterate over all slides."""
        return iter(self.slides)
