from gopptx import Presentation
import os
import sys

# Ensure smoke sample exists
input_deck = "examples/assets/01/01_basic_pptx.pptx"
if not os.path.exists(input_deck):
    print("Error: smoke sample missing. Run smoke tests first.")
    exit(1)

try:
    print("Testing Command-based JSON API...")
    with Presentation(input_deck) as pres:
        print(f"Handle: {pres._handle}")
        # 1. Slide Count
        print("Getting slide count...")
        initial_count = pres.slide_count
        print(f"Initial slide count: {initial_count}")
        
        # 2. Metadata
        print("Getting metadata...")
        meta = pres.metadata
        print(f"Metadata Title: {meta['title']}")
        print(f"Metadata Size: {meta['size']['width']}x{meta['size']['height']}")

        # 3. Duplicate Slide
        print("Duplicating slide 0 to index 1...")
        new_idx = pres.duplicate_slide(0, 1)
        print(f"New slide index: {new_idx}")
        print(f"New slide count: {pres.slide_count}")

        # 4. Move Slide
        print("Moving slide 0 to the end...")
        pres.move_slide(0, pres.slide_count - 1)
        
        # 5. Remove Slide
        print("Removing slide index 1...")
        pres.remove_slide(1)
        print(f"Final slide count: {pres.slide_count}")

        # 6. Save
        out_path = "examples/python_management_output.pptx"
        print(f"Saving to {out_path}...")
        pres.save(out_path)
        print("Saved.")

    print("\nSuccess! JSON Command API is working for slide management.")

except Exception as e:
    print(f"\nError during testing: {e}")
    if hasattr(e, "code"):
        print(f"Error Code: {e.code}")
    sys.exit(1)
