"""Custom XML example in Python.

Adds and inspects custom XML parts, then writes a round-trip deck.
"""

from __future__ import annotations

from pathlib import Path

from gopptx import Presentation


def main() -> None:
    root = Path(__file__).resolve().parents[2]
    output_dir = root / "examples" / "output"
    output_dir.mkdir(parents=True, exist_ok=True)
    output_path = output_dir / "27_custom_xml.pptx"
    roundtrip_path = output_dir / "27_custom_xml_roundtrip.pptx"

    with Presentation.new("Custom XML (Python)") as prs:
        prs.update_slide(0, title="Custom XML", bullets=["Metadata inside custom XML"])
        part_a = prs.add_custom_xml(
            "<CompanyData><Name>Acme</Name><ID>12345</ID></CompanyData>",
            root_element="CompanyData",
            namespace="http://schemas.example.com/company",
            properties={"Name": "Acme Corp", "ID": "12345"},
        )
        part_b = prs.add_custom_xml("<RawSettings><Theme>Dark</Theme></RawSettings>")
        parts = prs.list_custom_xml()
        prs.save(str(output_path))

    with Presentation(str(output_path)) as reopened:
        reopened_parts = reopened.list_custom_xml()
        if reopened_parts:
            reopened.remove_custom_xml(len(reopened_parts) - 1)
        reopened.save(str(roundtrip_path))

    print(f"Inserted part indexes: {part_a}, {part_b}")
    print(f"Custom XML part count before save: {len(parts)}")
    print(f"Created: {output_path}")
    print(f"Created: {roundtrip_path}")


if __name__ == "__main__":
    main()
