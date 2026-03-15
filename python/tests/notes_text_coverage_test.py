import pytest
from gopptx import Presentation

def test_notes_text_model_properties():
    with Presentation.new("Notes Coverage") as prs:
        slide = prs.add_slide("S1")
        # Accessing notes property triggers creation of notes slide if it doesn't exist
        slide.notes = "Initial notes"
        notes_slide = slide.notes_slide
        assert notes_slide is not None
        
        tf = notes_slide.notes_text_frame
        assert tf is not None
        
        # Paragraph properties
        para = tf.paragraphs[0]
        para.font_size = 14
        para.bold = True
        para.italic = True
        para.color = "FF0000"
        para.alignment = "ctr"
        
        assert para.font_size == 14
        assert para.bold is True
        assert para.italic is True
        assert para.color == "FF0000"
        assert para.alignment == "ctr"
        
        # Run properties
        run = para.runs[0]
        run.font_size = 12
        run.bold = False
        run.italic = False
        run.color = "0000FF"
        run.text = "Updated run"
        
        assert run.font_size == 12
        assert run.bold is False
        assert run.italic is False
        assert run.color == "0000FF"
        assert run.text == "Updated run"
        
        # Alignment aliases
        para.align_left()
        assert para.alignment == "l"
        para.align_center()
        assert para.alignment == "ctr"
        para.align_right()
        assert para.alignment == "r"
        para.align_justify()
        assert para.alignment == "just"

def test_notes_text_frame_add_paragraph():
    with Presentation.new("Notes Coverage") as prs:
        slide = prs.add_slide("S1")
        slide.notes = "P1" # Initialize
        notes_slide = slide.notes_slide
        assert notes_slide is not None
        tf = notes_slide.notes_text_frame
        assert tf is not None
        
        p2 = tf.add_paragraph("P2")
        assert len(tf.paragraphs) == 2
        assert p2.text == "P2"
        
        tf.clear()
        assert len(tf.paragraphs) == 0
        assert slide.notes == ""
