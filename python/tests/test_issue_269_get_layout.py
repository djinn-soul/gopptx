"""Issue #269: SlideMaster.get_layout(slide_layout_id, default=None).

The ids come from the master's p:sldLayoutIdLst, so the tests check them against
the master XML rather than trusting whatever the facade reports.
"""

from __future__ import annotations

import re
import zipfile
from typing import TYPE_CHECKING

from gopptx import Presentation

if TYPE_CHECKING:
    from pathlib import Path

_SLD_LAYOUT_ID_RE = re.compile(r'<p:sldLayoutId[^>]*\bid="(\d+)"')


def test_layout_ids_match_the_master_xml(tmp_path: Path) -> None:
    """slide_layout_id reports the ids the master actually lists."""
    output_path = tmp_path / "layout_ids.pptx"

    with Presentation.new(title="Issue 269") as prs:
        facade_ids = [
            layout.slide_layout_id for layout in prs.slide_masters[0].slide_layouts
        ]
        prs.save(output_path)

    with zipfile.ZipFile(output_path) as zf:
        master_xml = zf.read("ppt/slideMasters/slideMaster1.xml").decode("utf-8")
    xml_ids = [int(value) for value in _SLD_LAYOUT_ID_RE.findall(master_xml)]

    assert xml_ids
    assert facade_ids == xml_ids


def test_get_layout_returns_the_matching_layout() -> None:
    """A known id resolves to the layout carrying it."""
    with Presentation.new(title="Issue 269 hit") as prs:
        master = prs.slide_masters[0]
        expected = master.slide_layouts[1]

        found = master.get_layout(expected.slide_layout_id)

        assert found is not None
        assert found.name == expected.name
        assert found.part_name == expected.part_name


def test_get_layout_returns_none_for_unknown_id() -> None:
    """An unknown id yields None rather than raising."""
    with Presentation.new(title="Issue 269 miss") as prs:
        assert prs.slide_masters[0].get_layout(999_999) is None


def test_get_layout_returns_the_supplied_default() -> None:
    """The default is returned unchanged when the id is unknown."""
    with Presentation.new(title="Issue 269 default") as prs:
        master = prs.slide_masters[0]
        sentinel = master.slide_layouts[0]

        assert master.get_layout(999_999, sentinel) is sentinel


def test_layout_ids_are_unique() -> None:
    """Two layouts never share an id, or get_layout would be ambiguous."""
    with Presentation.new(title="Issue 269 unique") as prs:
        ids = [layout.slide_layout_id for layout in prs.slide_masters[0].slide_layouts]

        assert len(ids) == len(set(ids))
        assert all(layout_id > 0 for layout_id in ids)
