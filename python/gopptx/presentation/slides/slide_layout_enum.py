"""Named slide layout type constants for PowerPoint slide layouts."""

from __future__ import annotations


class SlideLayoutType:
    """Named slide layout type constants for all supported PowerPoint layouts.

    These constants provide a discoverable, type-safe way to specify slide layouts.
    ONLY the SlideLayoutType constants are accepted - string values are not supported.

    Examples:
        # Using enum-style constants (required)
        from gopptx.presentation.slides import SlideLayoutType

        # Add slide with specific layout
        slide = prs.add_slide("Title", layout=SlideLayoutType.BLANK)

        # Update slide layout
        prs.update_slide(0, layout=SlideLayoutType.TITLE_AND_CONTENT)

        # Discover available layouts
        all_layouts = SlideLayoutType.get_all()
    """

    # Standard Layouts
    BLANK = "blank"
    """Blank slide with no placeholders."""

    TITLE_ONLY = "title_only"
    """Slide with title placeholder only."""

    TITLE_AND_CONTENT = "title_and_content"
    """Slide with title and content (bullets) placeholders."""

    CENTERED_TITLE = "centered_title"
    """Slide with centered title layout."""

    @staticmethod
    def get_all() -> dict[str, str]:
        """Get all named slide layout type constants as a dictionary.

        Returns:
            Dictionary mapping constant names to layout values.
            Only includes currently supported layouts.

        Examples:
            all_layouts = SlideLayoutType.get_all()
            # Returns: {'BLANK': 'blank', 'TITLE_ONLY': 'title_only', ...}

            for name, value in all_layouts.items():
                print(f"SlideLayoutType.{name}")
        """
        return {
            "BLANK": SlideLayoutType.BLANK,
            "TITLE_ONLY": SlideLayoutType.TITLE_ONLY,
            "TITLE_AND_CONTENT": SlideLayoutType.TITLE_AND_CONTENT,
            "CENTERED_TITLE": SlideLayoutType.CENTERED_TITLE,
        }

    @staticmethod
    def validate(layout: str | None) -> str:
        """Validate a slide layout and return its value.

        Only accepts actual SlideLayoutType constant values (e.g., "blank", "title_only").
        Named constant strings are NOT supported - use the SlideLayoutType constants directly
        (e.g., SlideLayoutType.BLANK instead of "BLANK").

        Args:
            layout: Slide layout constant value (from SlideLayoutType enum).

        Returns:
            The same layout string (if valid).

        Raises:
            ValueError: If layout is not a valid layout value.

        Examples:
            SlideLayoutType.validate("blank")        # -> "blank" (valid)
            SlideLayoutType.validate("title_only")   # -> "title_only" (valid)
            SlideLayoutType.validate("BLANK")        # -> ValueError (use SlideLayoutType.BLANK)
            SlideLayoutType.validate("invalid")      # -> ValueError
        """
        if not layout:
            raise ValueError("layout cannot be empty")

        # Only accept actual layout values (not named constants)
        valid_layouts = set(SlideLayoutType.get_all().values())
        if layout in valid_layouts:
            return layout

        # Not found - provide helpful error message
        valid_values = ", ".join(sorted(valid_layouts))
        raise ValueError(
            f"Invalid layout {layout!r}. Use SlideLayoutType constants like SlideLayoutType.BLANK, SlideLayoutType.TITLE_ONLY. Valid internal values: {valid_values}"
        )

    @staticmethod
    def get_by_name(name: str) -> str | None:
        """Find a slide layout by its constant name.

        Args:
            name: Constant name like "BLANK", "TITLE_ONLY", "TITLE_AND_CONTENT".

        Returns:
            The layout value, or None if not found.

        Examples:
            SlideLayoutType.get_by_name("BLANK")           # -> "blank"
            SlideLayoutType.get_by_name("TITLE_ONLY")      # -> "title_only"
            SlideLayoutType.get_by_name("NOT_FOUND")       # -> None
        """
        return SlideLayoutType.get_all().get(name)


__all__ = ["SlideLayoutType"]
