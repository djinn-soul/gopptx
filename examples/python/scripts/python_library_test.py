import os
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../.."))
sys.path.append(os.path.join(project_root, "python"))

# Create directory if missing
output_dir = os.path.join(project_root, "examples/output")
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

try:
    print("Using high-level gopptx library...")
    input_path = os.path.join(project_root, "examples/assets/01/01_basic_pptx.pptx")
    with Presentation(input_path) as pres:
        print(f"Adding slides to {pres._handle}...")
        pres.add_slide("Slide 1: High-level API")
        pres.add_slide("Slide 2: Context Manager Support")

        out_path = os.path.join(output_dir, "python_library_output.pptx")
        pres.save(out_path)
        print(f"Saved to {out_path}")

    print("Success! The library is working as expected.")

except Exception as e:
    print(f"Error: {e}")
