import pytest
from gopptx.slide.notes.notes_text_model import NotesTextFrame

class DummyNotesShape:
    def __init__(self, text=""):
        self.text = text

def test_notes_text_model_plain_traversal():
    shape = DummyNotesShape("P1\nP2\nP3")
    tf = NotesTextFrame(shape)
    
    assert len(tf.paragraphs) == 3
    assert tf.paragraphs[0].text == "P1"
    assert tf.paragraphs[1].text == "P2"
    assert tf.paragraphs[2].text == "P3"
    
    # Update paragraph
    tf.paragraphs[1].text = "Updated P2"
    assert shape.text == "P1\nUpdated P2\nP3"
    
    # Run access
    run = tf.paragraphs[0].runs[0]
    assert run.text == "P1"
    run.text = "R1"
    assert shape.text == "R1\nUpdated P2\nP3"
    
    # Add run (append to single run)
    tf.paragraphs[0].add_run(" Extra")
    assert shape.text == "R1 Extra\nUpdated P2\nP3"
    
    # Add paragraph
    tf.add_paragraph("P4")
    assert len(tf.paragraphs) == 4
    assert tf.paragraphs[3].text == "P4"
    
    # Clear
    tf.clear()
    assert tf.text == ""
    assert len(tf.paragraphs) == 1 # paragraph_texts returns [""] for empty text
    assert tf.paragraphs[0].text == ""

def test_notes_text_model_empty_init():
    shape = DummyNotesShape("")
    tf = NotesTextFrame(shape)
    tf.add_paragraph("First")
    assert tf.text == "First"
    tf.add_paragraph("Second")
    assert tf.text == "First\nSecond"

def test_notes_text_model_index_errors():
    shape = DummyNotesShape("P1")
    tf = NotesTextFrame(shape)
    with pytest.raises(IndexError):
        tf.paragraphs[10]
    with pytest.raises(IndexError):
        tf.paragraphs[0].runs[5]
    
    # NotesRun and NotesParagraph internal check for out of range if text changed externally
    run = tf.paragraphs[0].runs[0]
    shape.text = "" 
    # Now paragraphs is [""] length 1, run paragraph_index 0 is OK.
    # To trigger error we'd need paragraph_index > 0.
    para = NotesTextFrame(DummyNotesShape("P1")).paragraphs[0]
    shape.text = ""
    # para._index is 0, paragraphs length is 1. Still OK.
