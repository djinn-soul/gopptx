from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from gopptx import ConnectorType, Presentation, ShapeType
from gopptx.presentation.export.export_mixin import PDFOptions
from gopptx.schemas import Inches


@dataclass(frozen=True)
class UsageEntry:
    id: str
    level: str
    title: str
    focus: str
    bullets: tuple[str, ...]


ENTRIES: tuple[UsageEntry, ...] = (
    UsageEntry(
        "S01",
        "Simple",
        "Hello World Deck",
        "Create first slide and save.",
        ("Create a new deck", "Add first slide", "Save PPTX output"),
    ),
    UsageEntry(
        "S02",
        "Simple",
        "Basic Text Frame",
        "Add controlled text regions.",
        ("Add title text frame", "Place body text blocks", "Save text-frame demo"),
    ),
    UsageEntry(
        "S03",
        "Simple",
        "Large Title Styling",
        "Highlight title readability.",
        ("Large heading text", "Strong visual hierarchy", "Readable slide structure"),
    ),
    UsageEntry(
        "S04",
        "Simple",
        "Shapes Showcase",
        "Use core shape primitives.",
        (
            "Rectangle content block",
            "Structured layout zones",
            "Quick slide composition",
        ),
    ),
    UsageEntry(
        "S05",
        "Simple",
        "Hyperlinks Basics",
        "Create link-ready content.",
        ("Title and context text", "Link placeholder callout", "Output for review"),
    ),
    UsageEntry(
        "S06",
        "Simple",
        "Speaker Notes",
        "Prepare presenter context.",
        ("Slide with speaker context", "Notes-ready structure", "Presentation save"),
    ),
    UsageEntry(
        "S07",
        "Simple",
        "Slide Properties",
        "Tune slide-level metadata.",
        ("Slide title setup", "Property-oriented content", "Consistent output"),
    ),
    UsageEntry(
        "S08",
        "Simple",
        "Background Fills",
        "Build visual slide base.",
        ("Background emphasis text", "Visual contrast block", "Ready-to-theme output"),
    ),
    UsageEntry(
        "I01",
        "Intermediate",
        "Radar Chart",
        "Show category vs value structure.",
        ("Radar categories", "Value interpretation", "Chart-ready narrative"),
    ),
    UsageEntry(
        "I02",
        "Intermediate",
        "Bubble Chart",
        "Represent multi-dimensional data.",
        ("Bubble chart context", "Axis interpretation", "Data story frame"),
    ),
    UsageEntry(
        "I03",
        "Intermediate",
        "Animations",
        "Stage progressive reveal intent.",
        ("Animation narrative", "Step-by-step points", "Playback-ready slide"),
    ),
    UsageEntry(
        "I04",
        "Intermediate",
        "Slide Duplication",
        "Model repeatable slide patterns.",
        ("Source slide concept", "Duplicate/edit workflow", "Versionable output"),
    ),
    UsageEntry(
        "I05",
        "Intermediate",
        "Image Stamping",
        "Place image-centric callouts.",
        ("Image placement zones", "Caption placeholders", "Visual asset pass"),
    ),
    UsageEntry(
        "I06",
        "Intermediate",
        "Presentation Properties",
        "Control document metadata.",
        ("Authoring metadata", "Title/subject flow", "Governance-friendly deck"),
    ),
    UsageEntry(
        "I07",
        "Intermediate",
        "Section Management",
        "Organize multi-part decks.",
        ("Section introduction", "Grouping strategy", "Navigation-friendly output"),
    ),
    UsageEntry(
        "I08",
        "Intermediate",
        "Commenting API",
        "Enable review workflows.",
        ("Reviewer context", "Comment target zone", "Collaborative output"),
    ),
    UsageEntry(
        "C01",
        "Complex",
        "Markdown + Mermaid Deck",
        "Generate decks from text specs.",
        ("Markdown ingestion", "Diagram rendering intent", "Automated slide synthesis"),
    ),
    UsageEntry(
        "C02",
        "Complex",
        "URL Fetch to Slides",
        "Convert web content into slides.",
        ("Web capture target", "Transformation notes", "Slide conversion output"),
    ),
    UsageEntry(
        "C03",
        "Complex",
        "Multi Master Deck",
        "Work with master/layout families.",
        ("Master family context", "Layout variance", "Template-aware output"),
    ),
    UsageEntry(
        "C04",
        "Complex",
        "Smart Merge Assets",
        "Merge decks and reuse content.",
        ("Source/target merge", "Asset reuse intent", "Conflict-safe output"),
    ),
    UsageEntry(
        "C05",
        "Complex",
        "Presentation Protection",
        "Apply protection-oriented settings.",
        ("Protection scenario", "Integrity expectations", "Secure distribution output"),
    ),
    UsageEntry(
        "C06",
        "Complex",
        "Advanced Hyperlinks",
        "Create internal/external navigation.",
        ("Cross-slide navigation", "External link intent", "Interactive deck output"),
    ),
    UsageEntry(
        "C07",
        "Complex",
        "Export & Distribution Pipeline",
        "Batch-export PPTX to PNG and PDF for reporting pipelines and CI/CD artifacts.",
        ("Build deck", "Export PDF", "Batch PNG publish"),
    ),
    UsageEntry(
        "C08",
        "Complex",
        "Rich Slide Composite",
        "Combine multiple advanced capabilities.",
        ("Composite layout intent", "Mixed content blocks", "Production-style output"),
    ),
)


ROOT = Path(__file__).resolve().parents[2]
PPTX_DIR = ROOT / "docs" / "assets" / "pptx" / "usage"
PDF_DIR = ROOT / "docs" / "assets" / "pdf" / "usage"
PNG_DIR = ROOT / "docs" / "assets" / "images" / "usage"
USAGE_DIR = ROOT / "docs" / "showcase" / "usages"


def entry_pptx_name(entry: UsageEntry) -> str:
    return f"{entry.id.lower()}-python.pptx"


def entry_png_name(entry: UsageEntry) -> str:
    return f"{entry.id.lower()}-python.png"


def entry_pdf_name(entry: UsageEntry) -> str:
    return f"{entry.id.lower()}-python.pdf"


def _build_c07_pipeline_slide(slide) -> None:
    cards = [
        (
            0.75,
            1.55,
            2.4,
            1.12,
            "Build PPTX\nquarterly_report.pptx",
            "D9EAF7",
            "4A86B8",
        ),
        (3.45, 1.55, 2.4, 1.12, "Export PDF\nquarterly_report.pdf", "E2F0D9", "5F8C5A"),
        (6.15, 1.55, 2.4, 1.12, "Batch PNG\nslide_*.png", "FCE4D6", "C27A2C"),
    ]
    for left, top, width, height, text, fill, line in cards:
        slide.add_shape(
            ShapeType.ROUNDED_RECTANGLE,
            (Inches(left), Inches(top), Inches(width), Inches(height)),
            text=text,
            properties={
                "fill": {"solid": fill},
                "line": {"color": line, "width_emu": 12700},
            },
        )

    connectors = [
        (3.17, 2.1, 3.38, 2.1),
        (5.87, 2.1, 6.08, 2.1),
    ]
    for begin_x, begin_y, end_x, end_y in connectors:
        slide.add_connector(
            ConnectorType.STRAIGHT,
            Inches(begin_x),
            Inches(begin_y),
            Inches(end_x),
            Inches(end_y),
            properties={"line": {"color": "64748B", "width_emu": 19050}},
        )

    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(0.95), Inches(3.25), Inches(7.55), Inches(1.0)),
        text="CI/CD friendly: generate one PPTX, archive a PDF, and publish all slide PNGs as build artifacts.",
        properties={
            "fill": {"solid": "F8FAFC"},
            "line": {"color": "CBD5E1", "width_emu": 12700},
        },
    )


def _build_c07_targets_slide(slide) -> None:
    targets = [
        (
            0.8,
            1.45,
            2.35,
            1.65,
            "PPTX source",
            "docs/assets/pptx/usage/c07-python.pptx",
            "DCE6F2",
            "5B9BD5",
        ),
        (
            3.25,
            1.45,
            2.35,
            1.65,
            "PDF export",
            "docs/assets/pdf/usage/c07-python.pdf",
            "E2F0D9",
            "70AD47",
        ),
        (
            5.7,
            1.45,
            2.35,
            1.65,
            "PNG batch",
            "docs/assets/images/usage/c07-python.png",
            "FCE4D6",
            "ED7D31",
        ),
    ]
    for left, top, width, height, title, path, fill, line in targets:
        slide.add_shape(
            ShapeType.ROUNDED_RECTANGLE,
            (Inches(left), Inches(top), Inches(width), Inches(height)),
            text=f"{title}\n\n{path}",
            properties={
                "fill": {"solid": fill},
                "line": {"color": line, "width_emu": 12700},
            },
        )

    slide.add_shape(
        ShapeType.ROUNDED_RECTANGLE,
        (Inches(0.85), Inches(3.55), Inches(7.45), Inches(0.95)),
        text="Reporting pipelines can publish the PDF for review and the PNGs for docs, while the PPTX remains the editable source artifact.",
        properties={
            "fill": {"solid": "FFFFFF"},
            "line": {"color": "CBD5E1", "width_emu": 12700},
        },
    )


def _build_c07_presentation(pres: Presentation) -> None:
    _build_c07_pipeline_slide(pres.slides[0])
    _build_c07_targets_slide(pres.add_slide("Distribution Targets"))


def generate_pptx(entry: UsageEntry) -> None:
    out_path = PPTX_DIR / entry_pptx_name(entry)
    out_path.parent.mkdir(parents=True, exist_ok=True)

    if entry.id == "C07":
        pdf_path = PDF_DIR / entry_pdf_name(entry)
        pdf_path.parent.mkdir(parents=True, exist_ok=True)
        with Presentation.new(f"{entry.id} {entry.title}") as pres:
            _build_c07_presentation(pres)
            pres.save(str(out_path))
            pres.save_as_pdf(str(pdf_path), options=PDFOptions(driver="native"))
        return

    with Presentation.new(f"{entry.id} {entry.title}") as pres:
        pres.add_bullet_slide(entry.title, list(entry.bullets))
        pres.save(str(out_path))


def go_code(entry: UsageEntry) -> str:
    if entry.id == "C07":
        return (
            "package main\n\n"
            "import (\n"
            '    "os"\n'
            '    "os/exec"\n'
            '    "path/filepath"\n\n'
            '    "github.com/djinn-soul/gopptx/pkg/pptx"\n'
            '    "github.com/djinn-soul/gopptx/pkg/pptx/export"\n'
            ")\n\n"
            "func buildPipelineSlide() pptx.SlideContent {\n"
            '    slide := pptx.NewSlide("Export & Distribution Pipeline")\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(0.75, 1.55, 2.4, 1.12).\n"
            '        WithText("Build PPTX\\nquarterly_report.pptx").\n'
            '        WithFill(pptx.NewShapeFill("D9EAF7")).\n'
            '        WithLine(pptx.NewShapeLine("4A86B8", pptx.Points(1))))\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(3.45, 1.55, 2.4, 1.12).\n"
            '        WithText("Export PDF\\nquarterly_report.pdf").\n'
            '        WithFill(pptx.NewShapeFill("E2F0D9")).\n'
            '        WithLine(pptx.NewShapeLine("5F8C5A", pptx.Points(1))))\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(6.15, 1.55, 2.4, 1.12).\n"
            '        WithText("Batch PNG\\nslide_*.png").\n'
            '        WithFill(pptx.NewShapeFill("FCE4D6")).\n'
            '        WithLine(pptx.NewShapeLine("C27A2C", pptx.Points(1))))\n'
            "    slide = slide.AddConnector(pptx.NewStraightConnector(3.17, 2.1, 3.38, 2.1).\n"
            '        WithLine(pptx.NewShapeLine("64748B", pptx.Points(1.5))))\n'
            "    slide = slide.AddConnector(pptx.NewStraightConnector(5.87, 2.1, 6.08, 2.1).\n"
            '        WithLine(pptx.NewShapeLine("64748B", pptx.Points(1.5))))\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(0.95, 3.25, 7.55, 1.0).\n"
            '        WithText("CI/CD friendly: generate one PPTX, archive a PDF, and publish all slide PNGs as build artifacts.").\n'
            '        WithFill(pptx.NewShapeFill("F8FAFC")).\n'
            '        WithLine(pptx.NewShapeLine("CBD5E1", pptx.Points(1))))\n'
            "    return slide\n"
            "}\n\n"
            "func buildTargetsSlide() pptx.SlideContent {\n"
            '    slide := pptx.NewSlide("Distribution Targets")\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(0.8, 1.45, 2.35, 1.65).\n"
            '        WithText("PPTX source\\ndocs/assets/pptx/usage/c07-python.pptx").\n'
            '        WithFill(pptx.NewShapeFill("DCE6F2")).\n'
            '        WithLine(pptx.NewShapeLine("5B9BD5", pptx.Points(1))))\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(3.25, 1.45, 2.35, 1.65).\n"
            '        WithText("PDF export\\ndocs/assets/pdf/usage/c07-python.pdf").\n'
            '        WithFill(pptx.NewShapeFill("E2F0D9")).\n'
            '        WithLine(pptx.NewShapeLine("70AD47", pptx.Points(1))))\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(5.7, 1.45, 2.35, 1.65).\n"
            '        WithText("PNG batch\\ndocs/assets/images/usage/c07-python.png").\n'
            '        WithFill(pptx.NewShapeFill("FCE4D6")).\n'
            '        WithLine(pptx.NewShapeLine("ED7D31", pptx.Points(1))))\n'
            "    slide = slide.AddShape(pptx.NewRoundedRectangle(0.85, 3.55, 7.45, 0.95).\n"
            '        WithText("Reporting pipelines can publish the PDF for review and the PNGs for docs, while the PPTX remains the editable source artifact.").\n'
            '        WithFill(pptx.NewShapeFill("FFFFFF")).\n'
            '        WithLine(pptx.NewShapeLine("CBD5E1", pptx.Points(1))))\n'
            "    return slide\n"
            "}\n\n"
            "func main() {\n"
            '    deckPath := filepath.Join("docs", "assets", "pptx", "usage", "c07-python.pptx")\n'
            '    pdfPath := filepath.Join("docs", "assets", "pdf", "usage", "c07-python.pdf")\n'
            '    pngDir, err := os.MkdirTemp("", "gopptx-c07-png-")\n'
            "    if err != nil { panic(err) }\n"
            "    defer os.RemoveAll(pngDir)\n\n"
            '    deck, err := pptx.CreateWithSlides("C07 Export & Distribution Pipeline", []pptx.SlideContent{buildPipelineSlide(), buildTargetsSlide()})\n'
            "    if err != nil { panic(err) }\n"
            "    if err := os.WriteFile(deckPath, deck, 0o600); err != nil { panic(err) }\n"
            "    if err := export.PDFFromFileWithOptions(deckPath, pdfPath, export.PDFOptions{Driver: export.PDFDriverNative}); err != nil { panic(err) }\n"
            '    if err := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", "scripts/tools/visual_regression/export_pptx_png.ps1", "-PptxPath", deckPath, "-OutDir", pngDir).Run(); err != nil { panic(err) }\n'
            "}\n"
        )
    bullet_lines = "\n".join([f'\tslide.AddBullet("{b}")' for b in entry.bullets])
    return (
        "package main\n\n"
        'import "github.com/djinn-soul/gopptx/pkg/gopptx"\n\n'
        "func main() {\n"
        f'\tpres := &gopptx.Presentation{{Title: "{entry.id} {entry.title}"}}\n'
        "\tslide := pres.AddSlide()\n"
        f'\tslide.Title = "{entry.title}"\n'
        f"{bullet_lines}\n"
        f'\t_ = pres.Save("{entry.id.lower()}-go.pptx")\n'
        "}\n"
    )


def python_code(entry: UsageEntry) -> str:
    if entry.id == "C07":
        return (
            "from pathlib import Path\n"
            "from gopptx import ConnectorType, Presentation, ShapeType\n"
            "from gopptx.presentation.export.export_mixin import PDFOptions\n"
            "from gopptx.schemas import Inches\n\n"
            "def build_pipeline_slide(slide):\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(0.75), Inches(1.55), Inches(2.4), Inches(1.12)),\n"
            '        text="Build PPTX\\nquarterly_report.pptx",\n'
            '        properties={"fill": {"solid": "D9EAF7"}, "line": {"color": "4A86B8", "width_emu": 12700}},\n'
            "    )\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(3.45), Inches(1.55), Inches(2.4), Inches(1.12)),\n"
            '        text="Export PDF\\nquarterly_report.pdf",\n'
            '        properties={"fill": {"solid": "E2F0D9"}, "line": {"color": "5F8C5A", "width_emu": 12700}},\n'
            "    )\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(6.15), Inches(1.55), Inches(2.4), Inches(1.12)),\n"
            '        text="Batch PNG\\nslide_*.png",\n'
            '        properties={"fill": {"solid": "FCE4D6"}, "line": {"color": "C27A2C", "width_emu": 12700}},\n'
            "    )\n"
            '    slide.add_connector(ConnectorType.STRAIGHT, Inches(3.17), Inches(2.1), Inches(3.38), Inches(2.1), properties={"line": {"color": "64748B", "width_emu": 19050}})\n'
            '    slide.add_connector(ConnectorType.STRAIGHT, Inches(5.87), Inches(2.1), Inches(6.08), Inches(2.1), properties={"line": {"color": "64748B", "width_emu": 19050}})\n'
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(0.95), Inches(3.25), Inches(7.55), Inches(1.0)),\n"
            '        text="CI/CD friendly: generate one PPTX, archive a PDF, and publish all slide PNGs as build artifacts.",\n'
            '        properties={"fill": {"solid": "F8FAFC"}, "line": {"color": "CBD5E1", "width_emu": 12700}},\n'
            "    )\n\n"
            "def build_targets_slide(slide):\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(0.8), Inches(1.45), Inches(2.35), Inches(1.65)),\n"
            '        text="PPTX source\\ndocs/assets/pptx/usage/c07-python.pptx",\n'
            '        properties={"fill": {"solid": "DCE6F2"}, "line": {"color": "5B9BD5", "width_emu": 12700}},\n'
            "    )\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(3.25), Inches(1.45), Inches(2.35), Inches(1.65)),\n"
            '        text="PDF export\\ndocs/assets/pdf/usage/c07-python.pdf",\n'
            '        properties={"fill": {"solid": "E2F0D9"}, "line": {"color": "70AD47", "width_emu": 12700}},\n'
            "    )\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(5.7), Inches(1.45), Inches(2.35), Inches(1.65)),\n"
            '        text="PNG batch\\ndocs/assets/images/usage/c07-python.png",\n'
            '        properties={"fill": {"solid": "FCE4D6"}, "line": {"color": "ED7D31", "width_emu": 12700}},\n'
            "    )\n"
            "    slide.add_shape(\n"
            "        ShapeType.ROUNDED_RECTANGLE,\n"
            "        (Inches(0.85), Inches(3.55), Inches(7.45), Inches(0.95)),\n"
            '        text="Reporting pipelines can publish the PDF for review and the PNGs for docs, while the PPTX remains the editable source artifact.",\n'
            '        properties={"fill": {"solid": "FFFFFF"}, "line": {"color": "CBD5E1", "width_emu": 12700}},\n'
            "    )\n\n"
            "repo = Path.cwd()\n"
            'pptx_path = repo / "docs" / "assets" / "pptx" / "usage" / "c07-python.pptx"\n'
            'pdf_path = repo / "docs" / "assets" / "pdf" / "usage" / "c07-python.pdf"\n'
            "pptx_path.parent.mkdir(parents=True, exist_ok=True)\n"
            "pdf_path.parent.mkdir(parents=True, exist_ok=True)\n"
            'with Presentation.new("C07 Export & Distribution Pipeline") as p:\n'
            "    build_pipeline_slide(p.slides[0])\n"
            '    build_targets_slide(p.add_slide("Distribution Targets"))\n'
            "    p.save(str(pptx_path))\n"
            '    p.save_as_pdf(str(pdf_path), options=PDFOptions(driver="native"))\n'
        )
    bullet_list = ",\n        ".join([f'"{b}"' for b in entry.bullets])
    return (
        "from gopptx import Presentation\n\n"
        f'with Presentation.new("{entry.id} {entry.title}") as p:\n'
        "    p.add_bullet_slide(\n"
        f'        "{entry.title}",\n'
        "        [\n"
        f"        {bullet_list}\n"
        "        ],\n"
        "    )\n"
        f'    p.save("docs/assets/pptx/usage/{entry.id.lower()}-python.pptx")\n'
    )


def render_level_page(level: str) -> str:
    level_entries = [e for e in ENTRIES if e.level == level]
    title = f"{level} Usages ({len(level_entries)})"
    intro = {
        "Simple": "Start here for first-time adoption and foundational slide automation.",
        "Intermediate": "Use these when you need richer visuals and editing workflows.",
        "Complex": "Use these for enterprise-grade decks, template families, and automation pipelines.",
    }[level]
    out: list[str] = [
        f"# {title}",
        "",
        intro,
        "",
        "Each usage is code-first and screenshot is generated from that Python code.",
    ]

    for entry in level_entries:
        downloads = [
            f"**Download PPTX:** [{entry_pptx_name(entry)}](../../assets/pptx/usage/{entry_pptx_name(entry)})"
        ]
        if entry.id == "C07":
            downloads.append(
                f"**Download PDF:** [{entry_pdf_name(entry)}](../../assets/pdf/usage/{entry_pdf_name(entry)})"
            )
        out.extend([
            "",
            f"## {entry.id} - {entry.title}",
            "",
            f"**Focus:** {entry.focus}",
            "",
            "**Go code**",
            "",
            "```go",
            go_code(entry).rstrip(),
            "```",
            "",
            "**Python code**",
            "",
            "```python",
            python_code(entry).rstrip(),
            "```",
            "",
            *downloads,
            "",
            "Screenshot generated from the Python code above using `export_pptx_png.ps1`.",
            "",
            f"![{entry.title}](../../assets/images/usage/{entry_png_name(entry)})",
        ])

    return "\n".join(out) + "\n"


def render_index_page() -> str:
    return (
        "\n".join([
            "# Usage Catalog (24)",
            "",
            "Organized from easiest to most advanced, with both Python and Go code in every entry.",
            "",
            "## Levels",
            "",
            "1. [Simple (8)](simple.md): quick wins and first workflows.",
            "2. [Intermediate (8)](intermediate.md): production-ready features.",
            "3. [Complex (8)](complex.md): advanced integrations and template-heavy pipelines.",
            "",
            "## Generation Contract",
            "",
            "- PPTX is generated from Python snippets in this catalog.",
            "- PNG screenshots are exported from those generated PPTX files.",
            "- Markdown, code snippets, PPTX names, and PNG names are produced from one source list.",
        ])
        + "\n"
    )


def main() -> None:
    PPTX_DIR.mkdir(parents=True, exist_ok=True)
    PDF_DIR.mkdir(parents=True, exist_ok=True)
    PNG_DIR.mkdir(parents=True, exist_ok=True)
    USAGE_DIR.mkdir(parents=True, exist_ok=True)

    for entry in ENTRIES:
        generate_pptx(entry)

    (USAGE_DIR / "index.md").write_text(render_index_page(), encoding="utf-8")
    (USAGE_DIR / "simple.md").write_text(render_level_page("Simple"), encoding="utf-8")
    (USAGE_DIR / "intermediate.md").write_text(
        render_level_page("Intermediate"), encoding="utf-8"
    )
    (USAGE_DIR / "complex.md").write_text(
        render_level_page("Complex"), encoding="utf-8"
    )

    print(f"Generated {len(ENTRIES)} PPTX files and usage markdown pages.")


if __name__ == "__main__":
    main()
