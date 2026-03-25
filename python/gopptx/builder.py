"""Fluent builder for creating PowerPoint presentations from scratch."""

from __future__ import annotations

from dataclasses import dataclass, field

from .presentation.presentation import Presentation


@dataclass
class _SlideSpec:
    title: str
    bullets: list[str] = field(default_factory=list)
    layout: str = ""


class PresentationBuilder:
    """Fluent builder for creating presentations from scratch.

    Creates a new presentation, applies settings, and returns a live
    :class:`~gopptx.Presentation` object ready for further editing or saving.

    Example::

        from gopptx import PresentationBuilder

        prs = (
            PresentationBuilder("Q1 Review")
            .with_author("Alice")
            .with_theme("corporate")
            .add_title_slide("Q1 2025 Results")
            .add_bullet_slide("Highlights", ["Revenue +12%", "New markets: 3"])
            .add_bullet_slide("Next Steps", ["Hire 5 engineers", "Launch APAC"])
            .build()
        )
        prs.save_as("output.pptx")
    """

    def __init__(self, title: str) -> None:
        """Initialise the builder with the presentation title."""
        super().__init__()
        self._title = title
        self._author: str = ""
        self._subject: str = ""
        self._keywords: str = ""
        self._description: str = ""
        self._theme: str = ""
        self._slide_width: float | None = None
        self._slide_height: float | None = None
        self._modify_password: str = ""
        self._mark_as_final: bool = False
        self._slides: list[_SlideSpec] = []

    # ── Metadata ────────────────────────────────────────────────────────────

    def with_author(self, author: str) -> PresentationBuilder:
        """Set the author / creator metadata field."""
        self._author = author
        return self

    def with_subject(self, subject: str) -> PresentationBuilder:
        """Set the subject metadata field."""
        self._subject = subject
        return self

    def with_keywords(self, keywords: str) -> PresentationBuilder:
        """Set the keywords metadata field."""
        self._keywords = keywords
        return self

    def with_description(self, description: str) -> PresentationBuilder:
        """Set the description metadata field."""
        self._description = description
        return self

    # ── Appearance ──────────────────────────────────────────────────────────

    def with_theme(self, theme: str) -> PresentationBuilder:
        """Apply a named theme (e.g. ``"corporate"``, ``"modern"``, ``"dark"``)."""
        self._theme = theme
        return self

    def with_slide_size(
        self, width_inches: float, height_inches: float
    ) -> PresentationBuilder:
        """Override the default slide dimensions (in inches).

        Common sizes:
        - 16x9 widescreen: ``with_slide_size(13.33, 7.5)``
        - 4x3 standard:    ``with_slide_size(10.0, 7.5)``
        """
        self._slide_width = width_inches
        self._slide_height = height_inches
        return self

    # ── Security ────────────────────────────────────────────────────────────

    def with_modify_password(self, password: str) -> PresentationBuilder:
        """Require a password to modify the presentation."""
        self._modify_password = password
        return self

    def with_mark_as_final(self, *, final: bool = True) -> PresentationBuilder:
        """Mark the presentation as final (read-only intent flag)."""
        self._mark_as_final = final
        return self

    # ── Slides ───────────────────────────────────────────────────────────────

    def add_title_slide(self, title: str, layout: str = "") -> PresentationBuilder:
        """Add a title-only slide."""
        self._slides.append(_SlideSpec(title=title, layout=layout))
        return self

    def add_bullet_slide(
        self,
        title: str,
        bullets: list[str],
        layout: str = "",
    ) -> PresentationBuilder:
        """Add a slide with a title and bullet points."""
        self._slides.append(
            _SlideSpec(title=title, bullets=list(bullets), layout=layout)
        )
        return self

    # ── Build ────────────────────────────────────────────────────────────────

    def build(self) -> Presentation:
        """Create the presentation and return a live :class:`~gopptx.Presentation`.

        The returned object has a handle to the Go engine and can be further
        edited with any ``Presentation`` method before calling ``save_as()``.
        """
        prs = Presentation.new(self._title)
        self._apply_core_properties(prs)
        if self._theme:
            prs.apply_theme(self._theme)
        if self._slide_width is not None and self._slide_height is not None:
            emu_per_inch = 914400
            prs.set_slide_size(
                int(self._slide_width * emu_per_inch),
                int(self._slide_height * emu_per_inch),
            )
        for spec in self._slides:
            prs.add_slide(title=spec.title, bullets=spec.bullets, layout=spec.layout)
        if self._modify_password:
            prs.set_modify_password(self._modify_password)
        if self._mark_as_final:
            prs.set_mark_as_final(final=True)
        return prs

    def _apply_core_properties(self, prs: Presentation) -> None:
        """Push metadata fields onto the presentation if any were set."""
        props: dict[str, object] = {"title": self._title}
        if self._author:
            props["creator"] = self._author
        if self._subject:
            props["subject"] = self._subject
        if self._keywords:
            props["keywords"] = self._keywords
        if self._description:
            props["description"] = self._description
        if len(props) > 1:
            prs.set_core_properties(props)  # type: ignore[arg-type]

    def save(self, path: str) -> None:
        """Build the presentation and save it to *path* in one step."""
        prs = self.build()
        prs.save(path)

    def to_bytes(self) -> bytes:
        """Build the presentation and return it as a byte string."""
        prs = self.build()
        return prs.to_bytes()
