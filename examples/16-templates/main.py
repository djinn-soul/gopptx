"""Demonstrate using a data-driven template to generate an invoice presentation.

This example demonstrates:
- Building a presentation that acts as a template with placeholder tokens
- Populating template data (company, date, invoice ID, line items)
- Using add_table_from_rows() to render line item data
- Saving the rendered result as a final presentation
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation
from gopptx.presentation.slides import SlideLayoutType
from gopptx.schemas import Inches

_LINE_ITEMS = [
    {"num": "1", "name": "Widget A", "qty": "10", "price": "$1,200"},
    {"num": "2", "name": "Widget B", "qty": "4", "price": "$800"},
    {"num": "3", "name": "Consulting", "qty": "8h", "price": "$1,600"},
]

_CONTEXT = {
    "company": "Acme Corp",
    "date": "2026-03-01",
    "invoice_id": "INV-2026-0042",
    "total": "$3,600",
    "notes": "Thank you for your business!",
}


def _add_cover_slide(prs: Presentation) -> None:
    """Add the invoice cover slide with company details."""
    prs.add_bullet_slide(
        f"INVOICE \u2022 {_CONTEXT['company']}",
        [
            f"Date:        {_CONTEXT['date']}",
            f"Invoice #:   {_CONTEXT['invoice_id']}",
            f"Total:       {_CONTEXT['total']}",
        ],
    )


def _add_line_items_slide(prs: Presentation) -> None:
    """Add the line items table slide."""
    prs.add_slide("Line Items", layout=SlideLayoutType.TITLE_ONLY)
    idx = prs.slide_count - 1

    rows = [["#", "Item", "Qty", "Price"]]
    rows.extend(
        [item["num"], item["name"], item["qty"], item["price"]] for item in _LINE_ITEMS
    )

    prs.add_table_from_rows(
        idx,
        rows,
        bounds=(Inches(0.5), Inches(1.5), Inches(9), Inches(0.5 + 0.6 * len(rows))),
        first_row=True,
        band_row=True,
    )


def _add_notes_slide(prs: Presentation) -> None:
    """Add the client notes slide."""
    prs.add_bullet_slide(
        "Client Notes",
        [
            _CONTEXT["notes"],
            "Payment due within 30 days.",
            f"Account team: {_CONTEXT['company']} Finance",
        ],
    )


def main() -> None:
    """Create a rendered invoice presentation from template data."""
    output_dir = Path("examples/output")
    output_dir.mkdir(exist_ok=True)

    with Presentation.new("Invoice Template") as prs:
        _add_cover_slide(prs)
        _add_line_items_slide(prs)
        _add_notes_slide(prs)

        output_path = output_dir / "16-templates.pptx"
        prs.save(str(output_path))
        print(f"Saved: {output_path}")

    print("\n=== SUMMARY ===")
    print("Generated an invoice presentation from template data:")
    print(f"  Company: {_CONTEXT['company']}")
    print(f"  Invoice: {_CONTEXT['invoice_id']}")
    print(f"  Total:   {_CONTEXT['total']}")
    print(f"  Items:   {len(_LINE_ITEMS)} line items")


if __name__ == "__main__":
    main()
