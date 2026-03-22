from gopptx.slide.text.text_frame import TextFrameProps
import pytest

def test_text_frame_props_legacy_alias_override():
    """Verify that legacy aliases override new fields, even if new fields are invalid."""
    # This would have failed before the fix due to early normalization of 'auto_fit_type'
    props = TextFrameProps(auto_fit_type="invalid_type", auto_size="shape_to_fit_text")
    assert props.auto_fit_type == "shape"

    props = TextFrameProps(vertical_align="invalid_align", vertical_anchor="middle")
    assert props.vertical_align == "ctr"

    props = TextFrameProps(orientation="invalid_orient", text_direction="vertical")
    assert props.orientation == "vert"

def test_text_frame_props_normalization():
    """Verify that values are still normalized."""
    props = TextFrameProps(auto_fit_type="spautofit")
    assert props.auto_fit_type == "shape"
    
    props = TextFrameProps(vertical_align="center")
    assert props.vertical_align == "ctr"

if __name__ == "__main__":
    pytest.main([__file__])
