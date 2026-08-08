package text

import "testing"

// Persian, Urdu, Pashto and Kurdish fell through the table to Latin, so their
// text was written to <a:latin> instead of <a:cs>.
func TestArabicScriptLanguagesAreComplex(t *testing.T) {
	for _, tag := range []string{"ar", "fa", "ur", "ps", "ku", "ckb", "fa-IR", "ur_PK"} {
		if got := LanguageToScript(tag); got != ScriptCodeArabic {
			t.Errorf("LanguageToScript(%q) = %q, want %q", tag, got, ScriptCodeArabic)
		}
		if got := ScriptKindForLanguage(tag); got != ScriptComplex {
			t.Errorf("ScriptKindForLanguage(%q) = %v, want ScriptComplex", tag, got)
		}
	}
}

func TestRTLLanguageDirectionAndFont(t *testing.T) {
	for _, tag := range []string{"ar", "he", "fa", "ur", "ps", "ku", "yi", "dv"} {
		if !IsRTLLanguage(tag) {
			t.Errorf("IsRTLLanguage(%q) = false, want true", tag)
		}
		if RTLDefaultFont(tag) == "" {
			t.Errorf("RTLDefaultFont(%q) is empty", tag)
		}
	}
	for _, tag := range []string{"en", "fr", "ja", "hi", ""} {
		if IsRTLLanguage(tag) {
			t.Errorf("IsRTLLanguage(%q) = true, want false", tag)
		}
		if got := RTLDefaultFont(tag); got != "" {
			t.Errorf("RTLDefaultFont(%q) = %q, want empty", tag, got)
		}
	}
}
