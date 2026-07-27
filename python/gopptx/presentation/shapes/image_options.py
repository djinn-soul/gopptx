"""Keyword-option validation for image insertion.

add_image takes **kwargs, so a misspelled option used to be dropped without a
word and the picture appeared without the rotation, crop or alt text the caller
asked for. Validating the names here makes that a loud error instead.
"""

from __future__ import annotations

__all__ = ["IMAGE_OPTION_KEYS", "reject_unknown_image_options"]

# Every keyword the image mixin's source, format and option resolvers read.
IMAGE_OPTION_KEYS = frozenset({
    "path",
    "data",
    "image_format",
    "img_format",
    "crop",
    "rotation",
    "flip_h",
    "flip_v",
    "description",
    "alt_text",
    "title",
})


def reject_unknown_image_options(kwargs: dict[str, object]) -> None:
    """Raise on an image keyword option the mixin does not understand.

    Args:
        kwargs: Keyword options handed to add_image.

    Raises:
        TypeError: If any option name is not recognized.
    """
    unknown = sorted(set(kwargs) - IMAGE_OPTION_KEYS)
    if unknown:
        raise TypeError(
            " ".join((
                f"unexpected keyword argument(s) {', '.join(unknown)};",
                f"supported: {', '.join(sorted(IMAGE_OPTION_KEYS))}",
            ))
        )
