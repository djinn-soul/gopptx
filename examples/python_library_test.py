from gopptx import Presentation
import os

# Create directory if missing
os.makedirs("examples", exist_ok=True)

try:
    print("Using high-level gopptx library...")
    with Presentation("examples/assets/01/01_basic_pptx.pptx") as pres:
        print(f"Adding slides to {pres._handle}...")
        pres.add_slide("Slide 1: High-level API")
        pres.add_slide("Slide 2: Context Manager Support")
        
        out_path = "examples/python_library_output.pptx"
        pres.save(out_path)
        print(f"Saved to {out_path}")

    print("Success! The library is working as expected.")

except Exception as e:
    print(f"Error: {e}")
