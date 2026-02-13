import ctypes
import os
import sys

# Load the shared library
lib_path = "bindings/c/build/gopptx.dll" if sys.platform == "win32" else "bindings/c/build/libgopptx.so"
if sys.platform == "darwin":
    lib_path = "bindings/c/build/libgopptx.dylib"

if not os.path.exists(lib_path):
    print(f"Error: Library not found at {lib_path}. Run build script first.")
    sys.exit(1)

lib = ctypes.CDLL(lib_path)

# Define function signatures
lib.deck_open.argtypes = [ctypes.c_char_p]
lib.deck_open.restype = ctypes.c_void_p

lib.deck_add_slide.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
lib.deck_add_slide.restype = ctypes.c_int

lib.deck_save.argtypes = [ctypes.c_void_p, ctypes.c_char_p]
lib.deck_save.restype = ctypes.c_int

lib.deck_last_error.argtypes = [ctypes.c_void_p]
lib.deck_last_error.restype = ctypes.c_char_p

lib.deck_free_string.argtypes = [ctypes.c_char_p]
lib.deck_free_string.restype = None

lib.deck_close.argtypes = [ctypes.c_void_p]
lib.deck_close.restype = None

def get_last_error(h):
    err_ptr = lib.deck_last_error(h)
    if err_ptr:
        err_msg = ctypes.string_at(err_ptr).decode('utf-8')
        lib.deck_free_string(err_ptr)
        return err_msg
    return "Unknown error"

# Demo usage
deck_path = "examples/assets/01/01_basic_pptx.pptx".encode('utf-8')
out_path = "examples/python_modified.pptx".encode('utf-8')

print("Opening deck...")
h = lib.deck_open(deck_path)
if not h:
    print("Failed to open deck")
    sys.exit(1)

print(f"Deck opened with handle: {h}")

print("Adding slide...")
if lib.deck_add_slide(h, "Added via Python".encode('utf-8')) != 0:
    print(f"Error adding slide: {get_last_error(h)}")
else:
    print("Slide added successfully")

print("Saving deck...")
if lib.deck_save(h, out_path) != 0:
    print(f"Error saving deck: {get_last_error(h)}")
else:
    print(f"Deck saved to {out_path.decode('utf-8')}")

print("Closing deck...")
lib.deck_close(h)
print("Done!")
