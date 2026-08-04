"""Language-tag to script mapping (Issue #172).

PowerPoint renders a run from one of three typeface slots: ``<a:latin>`` for
Latin text, ``<a:ea>`` for East Asian text and ``<a:cs>`` for complex scripts.
Knowing a run's script is what lets a caller — or gopptx itself — put the font
in the slot the text is actually drawn from, instead of watching a Japanese or
Arabic run fall back to the theme font.

The mapping mirrors the Go implementation in ``pkg/pptx/text/lang_script.go``.
"""

from __future__ import annotations

# Primary language subtag -> ISO 15924 script code. Chinese is resolved
# separately because it depends on the region subtag.
_LANG_SCRIPT: dict[str, str] = {
    "ja": "Jpan",
    "ko": "Hang",
    "ar": "Arab",
    "he": "Hebr",
    "yi": "Hebr",
    "iw": "Hebr",
    "th": "Thai",
    "ti": "Ethi",
    "bwo": "Ethi",
    "eth": "Ethi",
    "kxh": "Ethi",
    "mdy": "Ethi",
    "bn": "Beng",
    "as": "Beng",
    "mni": "Beng",
    "gu": "Gujr",
    "km": "Khmr",
    "kn": "Knda",
    "pa": "Guru",
    "iu": "Cans",
    "chr": "Cher",
    "bo": "Tibt",
    "dv": "Thaa",
    "hi": "Deva",
    "ks": "Deva",
    "kok": "Deva",
    "mr": "Deva",
    "ne": "Deva",
    "sa": "Deva",
    "sd": "Deva",
    "te": "Telu",
    "ta": "Taml",
    "syr": "Syrc",
    "or": "Orya",
    "ml": "Mlym",
    "lo": "Laoo",
    "si": "Sinh",
    "vi": "Viet",
    "lha": "Viet",
    "nut": "Viet",
    "ka": "Geor",
}

# Scripts PowerPoint draws from the <a:ea> typeface; the rest of the non-Latin
# table is drawn from <a:cs>.
_EAST_ASIAN_SCRIPTS = frozenset({"Jpan", "Hang", "Hans", "Hant"})

SCRIPT_LATIN = "latin"
SCRIPT_EAST_ASIAN = "east_asian"
SCRIPT_COMPLEX = "complex"


def lang_to_script(lang_tag: str) -> str | None:
    """Return the ISO 15924 script code for a BCP-47 language tag.

    Returns ``None`` for Latin-script or unrecognised tags. Both ``-`` and ``_``
    are accepted as the subtag separator.

    >>> lang_to_script("ja")
    'Jpan'
    >>> lang_to_script("zh_CN")
    'Hans'
    >>> lang_to_script("zh-TW")
    'Hant'
    >>> lang_to_script("en-US") is None
    True
    """
    tag = (lang_tag or "").strip()
    if not tag:
        return None
    normalized = tag.replace("_", "-")
    lang, _, region = normalized.partition("-")
    lang = lang.lower()
    region = region.lower()

    if lang == "zh":
        # Mainland China and Singapore use simplified characters; everywhere
        # else that writes Chinese uses traditional.
        if region in {"cn", "sg"} or region.startswith("hans"):
            return "Hans"
        return "Hant"
    return _LANG_SCRIPT.get(lang)


def script_kind_for_language(lang_tag: str) -> str:
    """Return which typeface slot a language tag's text is drawn from.

    One of :data:`SCRIPT_LATIN`, :data:`SCRIPT_EAST_ASIAN` or
    :data:`SCRIPT_COMPLEX`, matching ``<a:latin>``, ``<a:ea>`` and ``<a:cs>``.
    """
    script = lang_to_script(lang_tag)
    if script is None:
        return SCRIPT_LATIN
    if script in _EAST_ASIAN_SCRIPTS:
        return SCRIPT_EAST_ASIAN
    return SCRIPT_COMPLEX
