"""Regressions for the review findings on the upstream-issue batch."""

import pathlib
import re
import zipfile

import pytest
from gopptx import Presentation
from gopptx.api_errors import GopptxError

_TEMPLATE = (
    pathlib.Path(__file__).resolve().parents[2]
    / "examples"
    / "assets"
    / "37"
    / "162301-moneybox-template-16x9.pptx"
)


def _part(path: pathlib.Path, name: str) -> str:
    with zipfile.ZipFile(path) as archive:
        return archive.read(name).decode("utf-8")


def _placeholder_types(xml: str) -> list[str]:
    found: list[str] = []
    for tag in re.findall(r"<p:ph\b[^>]*>", xml):
        match = re.search(r'type="([^"]+)"', tag)
        ph_type = match.group(1) if match else "body"
        if ph_type not in {"dt", "ftr", "sldNum"}:
            found.append(ph_type)
    return found


@pytest.mark.skipif(not _TEMPLATE.is_file(), reason="template asset not present")
def test_add_slide_takes_placeholders_from_layout_object(
    tmp_path: pathlib.Path,
) -> None:
    """A SlideLayout object decides the slide's placeholders, not just its binding.

    Adding a default slide and retargeting its layout relationship afterwards
    left the default title and body placeholders in place whatever the layout
    said, so a picture or extra content placeholder never appeared.
    """
    output = tmp_path / "layout_object.pptx"

    prs = Presentation()
    prs.open(str(_TEMPLATE))
    with prs:
        chosen = None
        for layout in prs.slide_layouts:
            wanted = _placeholder_types(_part(_TEMPLATE, layout.part))
            if len(wanted) >= 2:
                chosen = (layout, wanted)
                break
        assert chosen is not None, "template has no layout with placeholders"
        layout, wanted = chosen

        prs.add_slide(title="From Layout", bullets=["one"], layout=layout)
        prs.save(str(output))

    with zipfile.ZipFile(output) as archive:
        slides = sorted(
            n
            for n in archive.namelist()
            if n.startswith("ppt/slides/slide") and n.endswith(".xml")
        )
    slide_xml = _part(output, slides[-1])

    assert _placeholder_types(slide_xml) == wanted
    assert "From Layout" in slide_xml


def test_import_layout_inside_batch_is_refused(tmp_path: pathlib.Path) -> None:
    """A foreign layout cannot be imported from inside a batch.

    The import has to finish before the queued add can name the part it made,
    and a queued op returns nothing to name it with. Refusing beats the
    ``layout_part`` KeyError the queued path used to raise.
    """
    source_path = tmp_path / "source.pptx"
    with Presentation.new(title="source") as source:
        source.save(str(source_path))

    source = Presentation()
    source.open(str(source_path))
    with source, Presentation.new(title="dest") as dest:
        foreign = source.slide_layouts[1]
        with (
            dest.batch(),
            pytest.raises(GopptxError, match="not allowed inside a batch"),
        ):
            dest.add_slide(title="queued", layout=foreign)


def test_copied_slide_keeps_its_own_layout(tmp_path: pathlib.Path) -> None:
    """A copied slide brings its layout family instead of binding to a local part.

    Without the import the slide bound to whatever destination part shared the
    name, changing how it looks, or to nothing at all.
    """
    source_path = tmp_path / "copy_source.pptx"
    with Presentation.new(title="source") as source:
        source.add_slide(title="Source Slide", bullets=["one"])
        source.save(str(source_path))

    output = tmp_path / "copied.pptx"
    with Presentation.new(title="dest") as dest:
        dest.copy_slides_from(str(source_path), [0])
        dest.save(str(output))

    with zipfile.ZipFile(output) as archive:
        names = archive.namelist()
    masters = [n for n in names if n.startswith("ppt/slideMasters/slideMaster")]
    assert len(masters) >= 2, f"the copied slide's master was not imported: {masters}"

    rels = _part(output, "ppt/slides/_rels/slide2.xml.rels")
    assert "slideLayout" in rels

    # The copied slide must have its own presentation relationship: an import
    # that allocated rel ids mid-merge stole the one the slide needed and left a
    # package PowerPoint refused to open.
    presentation_rels = _part(output, "ppt/_rels/presentation.xml.rels")
    assert "slides/slide2.xml" in presentation_rels


def test_fit_text_keeps_paragraphs(tmp_path: pathlib.Path) -> None:
    """Shrinking the text must not merge the shape's paragraphs into one."""
    font = pathlib.Path("C:/Windows/Fonts/arial.ttf")
    if not font.is_file():
        pytest.skip("no TrueType font available to measure with")

    output = tmp_path / "fit_text.pptx"
    with Presentation.new(title="fit") as prs:
        slide = prs.slides[0]
        shape_id = slide.add_textbox(914400, 914400, 3657600, 1828800, text="alpha")
        frame = slide.shape(shape_id).text_frame
        frame.text = "alpha"
        frame.add_paragraph("beta")
        frame.add_paragraph("gamma")

        before = len(frame.paragraphs)
        frame.fit_text(font_file=str(font), max_size=28)
        prs.save(str(output))

    slide_xml = _part(output, "ppt/slides/slide1.xml")
    assert "<a:t>beta</a:t>" in slide_xml
    assert "<a:t>gamma</a:t>" in slide_xml
    assert before == 3


def test_equation_font_size_reaches_the_math_runs(tmp_path: pathlib.Path) -> None:
    """The requested size has to be on the equation, not on a trailing endParaRPr."""
    output = tmp_path / "equation.pptx"
    with Presentation.new(title="equation") as prs:
        prs.slides[0].add_equation(
            r"E = mc^2",
            left=914400,
            top=914400,
            width=5486400,
            height=1828800,
            font_size_pt=40,
        )
        prs.save(str(output))

    slide_xml = _part(output, "ppt/slides/slide1.xml")
    assert '<m:r><a:rPr lang="en-US" sz="4000"' in slide_xml
