import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = (pathlib.Path(__file__).parent / ".." / ".." / "..").resolve()
sys.path.append(str(project_root / "python"))

# Create directory if missing
output_dir = project_root / "examples" / "output"
output_dir.mkdir(exist_ok=True, parents=True)

try:
    print("Using high-level gopptx library...")
    input_path = project_root / "examples" / "assets" / "01" / "01_basic_pptx.pptx"
    with Presentation(str(input_path)) as pres:
        print("Adding slides to presentation...")
        pres.add_slide("Slide 1: High-level API")
        pres.add_slide("Slide 2: Context Manager Support")

        out_path = output_dir / "python_library_output.pptx"
        pres.save(str(out_path))
        print(f"Saved to {out_path}")

    print("Success! The library is working as expected.")

except (OSError, RuntimeError, ValueError) as e:
    print(f"Error: {e}")
