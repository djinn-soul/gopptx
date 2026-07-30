"""Chart data subset tests (upstream python-pptx issue #1043)."""

import io
import pathlib
import zipfile

import pytest
from gopptx.presentation.charts.chart_types import ChartType
from gopptx.presentation.presentation import Presentation
from gopptx.schemas import Inches

CATEGORIES = ["Q1", "Q2", "Q3", "Q4"]


@pytest.fixture
def deck(tmp_path: pathlib.Path) -> pathlib.Path:
    output_path = tmp_path / "chart.pptx"
    with Presentation.new("Chart") as pres:
        pres.slides[0].add_chart(
            ChartType.BAR,
            CATEGORIES,
            [{"name": "Revenue", "values": [12.0, 18.0, 15.0, 22.0]}],
            title="Chart",
            bounds=(Inches(1.2), Inches(1.3), Inches(7.0), Inches(3.6)),
        )
        pres.save(output_path)
    return output_path


def _embedded_sheet(deck_path: pathlib.Path) -> str:
    with zipfile.ZipFile(deck_path) as archive:
        chart_rels = archive.read("ppt/charts/_rels/chart1.xml.rels").decode()
        target = chart_rels.split('Target="../embeddings/')[1].split('"')[0]
        book_bytes = archive.read(f"ppt/embeddings/{target}")
    with zipfile.ZipFile(io.BytesIO(book_bytes)) as book:
        return book.read("xl/worksheets/sheet1.xml").decode()


def test_hidden_series_stays_in_workbook_only(deck: pathlib.Path) -> None:
    """The data behind the chart may be a superset of what it plots."""
    pres = Presentation()
    pres.open(deck)
    pres.update_chart_data(
        0,
        {"index": 0},
        {
            "categories": CATEGORIES,
            "series": [
                {"name": "Revenue", "values": [12.0, 18.0, 15.0, 22.0]},
                {"name": "Cost", "values": [8.0, 11.0, 9.0, 13.0], "hidden": True},
                {
                    "name": "Forecast",
                    "values": [14.0, 20.0, 17.0, 25.0],
                    "hidden": True,
                },
            ],
        },
    )
    pres.save(deck)
    pres.close()

    with zipfile.ZipFile(deck) as archive:
        chart_xml = archive.read("ppt/charts/chart1.xml").decode()
    assert chart_xml.count("<c:ser>") == 1
    assert "Forecast" not in chart_xml

    sheet = _embedded_sheet(deck)
    for header in ("Revenue", "Cost", "Forecast"):
        assert header in sheet, f"{header} should be in the embedded workbook"


def test_chart_update_declares_new_embedding_content_type(deck: pathlib.Path) -> None:
    """A chart data update writes a new embedding, which needs a content type.

    Without one the package is invalid OPC and PowerPoint refuses to open it
    with 0x80CB8002.
    """
    pres = Presentation()
    pres.open(deck)
    pres.update_chart_data(
        0,
        {"index": 0},
        {
            "categories": CATEGORIES,
            "series": [{"name": "Revenue", "values": [1.0, 2.0, 3.0, 4.0]}],
        },
    )
    pres.save(deck)
    pres.close()

    with zipfile.ZipFile(deck) as archive:
        content_types = archive.read("[Content_Types].xml").decode()
        embeddings = [n for n in archive.namelist() if n.startswith("ppt/embeddings/")]

    assert embeddings
    declares_default = 'Extension="xlsx"' in content_types
    for part in embeddings:
        assert declares_default or f'PartName="/{part}"' in content_types, (
            f"{part} has no content type"
        )
