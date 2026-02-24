import pathlib
import sys
import traceback

from gopptx import Presentation

project_root = (pathlib.Path(__file__).parent / ".." / ".." / "..").resolve()
sys.path.append(str(project_root / "python"))
output_dir = project_root / "examples" / "output"
output_dir.mkdir(exist_ok=True, parents=True)
input_path = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"
output_path = output_dir / "python_comprehensive_output.pptx"
if not input_path.exists():
    print(f"Input file {input_path} not found. Skipping test.")
    sys.exit(0)
try:
    print(f"Opening presentation: {input_path}")
    with Presentation(str(input_path)) as pres:
        count = pres.slide_count
        print(f"Original slide count: {count}")
        meta = pres.metadata
        print(f"Metadata: {meta}")
        idx = pres.add_slide("New Python Slide")
        print(f"Added slide at index: {idx.index}")
        print("Adding shape...")
        shape_id = pres.add_shape(
            idx.index,
            "rect",
            bounds=(1000000, 1000000, 2000000, 1000000),
            text="Hello from Python!",
        )
        print(f"Added shape with ID: {shape_id}")
        print("Listing shapes...")
        shapes = pres.list_shapes(idx.index)
        for s in shapes:
            print(f"  Shape: {s}")
        print("Updating shape...")
        pres.update_shape(idx.index, shape_id, {"text": "Updated Text", "x": 1500000})
        print("Setting notes...")
        pres.set_notes(idx.index, "These are speaker notes added via Python API.")
        notes = pres.get_notes(idx.index)
        print(f"Retrieved notes: {notes}")
        print("Listing layouts...")
        layouts = pres.list_slide_layouts()
        if layouts:
            print(f"  First layout: {layouts[0]}")
        print("Updating core properties...")
        props = pres.core_properties
        props["title"] = "Comprehensive Python Test"
        props["creator"] = "GoPPTX Python API"
        pres.core_properties = props
        pres.save(str(output_path))
        print(f"Saved to {output_path}")
    print("Success! Python API expansion verified.")
except (OSError, RuntimeError, ValueError) as e:
    print(f"Error during verification: {e}")
    traceback.print_exc()
