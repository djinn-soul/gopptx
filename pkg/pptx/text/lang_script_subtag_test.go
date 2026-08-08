package text

import "testing"

// A tag that names its script means it: the primary language's usual script is
// only a default for tags that leave it unsaid.
func TestExplicitScriptSubtagOutranksTheLanguageDefault(t *testing.T) {
	cases := map[string]string{
		"ku-Latn":    "Latn",
		"ku-Arab":    ScriptCodeArabic,
		"ku":         ScriptCodeArabic,
		"sr-Cyrl":    "Cyrl",
		"zh-Hans":    ScriptCodeHanSimp,
		"zh-Hant":    ScriptCodeHanTrad,
		"zh-CN":      ScriptCodeHanSimp,
		"zh-TW":      ScriptCodeHanTrad,
		"ku-latn-iq": "Latn",
		"KU-LATN":    "Latn",
	}
	for tag, want := range cases {
		if got := LanguageToScript(tag); got != want {
			t.Errorf("LanguageToScript(%q) = %q, want %q", tag, got, want)
		}
	}
}

// A two-letter second subtag is a region and must not be mistaken for a script.
func TestRegionSubtagIsNotReadAsAScript(t *testing.T) {
	if got := LanguageToScript("ar-EG"); got != ScriptCodeArabic {
		t.Errorf("LanguageToScript(ar-EG) = %q, want %q", got, ScriptCodeArabic)
	}
	if got := LanguageToScript("en-US"); got != "" {
		t.Errorf("LanguageToScript(en-US) = %q, want no script", got)
	}
}

func TestDirectionFollowsAnExplicitScript(t *testing.T) {
	rtl := []string{"ku", "ku-Arab", "ar", "he", "fa-Arab"}
	ltr := []string{"ku-Latn", "en", "sr-Cyrl", "az-Latn"}

	for _, tag := range rtl {
		if !IsRTLLanguage(tag) {
			t.Errorf("IsRTLLanguage(%q) = false, want true", tag)
		}
	}
	for _, tag := range ltr {
		if IsRTLLanguage(tag) {
			t.Errorf("IsRTLLanguage(%q) = true, want false", tag)
		}
	}
}
