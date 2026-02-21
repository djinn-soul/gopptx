import os
import pathlib
import sys

from gopptx import Presentation

# Add project root to sys.path to find 'gopptx' package
project_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../.."))
sys.path.append(os.path.join(project_root, "python"))


# Create output directory
output_dir = os.path.join(project_root, "examples/output")
pathlib.Path(output_dir).mkdir(exist_ok=True, parents=True)

output_path = os.path.join(output_dir, "python_sections_test.pptx")

try:
    print("Testing Sections property and dynamic indexing...")
    with Presentation.new("Sections Test") as pres:
        print(f"Presentation: {pres}")

        # 1. Create a few slides
        s1 = pres.add_slide("Slide 1")
        s2 = pres.add_slide("Slide 2")
        s3 = pres.add_slide("Slide 3")

        # 2. Add sections
        print("Adding sections...")
        pres.add_section("Introduction", [0, 1])
        pres.add_section("Main Content", [2, 3])

        # 3. Access sections
        sections = pres.sections
        print(f"Found {len(sections)} sections:")
        for sec in sections:
            print(f"  - {sec['Name']} (Slides: {sec['SlideIDs']})")

        # 4. Verify dynamic indexing
        print(f"Slide 3 original index: {s3.index}")
        pres.remove_slide(0)
        print(f"After removing first slide, Slide 3 index: {s3.index}")

        if s3.index == 2:
            print("Dynamic indexing works!")
        else:
            print(f"Dynamic indexing FAILED (expected 2, got {s3.index})")

        pres.save(output_path)
        print(f"Saved to {output_path}")

    print("Success! Sections and dynamic indexing verified.")

except Exception as e:
    print(f"Error during verification: {e}")
    import traceback

    traceback.print_exc()
