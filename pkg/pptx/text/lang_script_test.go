package text

import "testing"

func TestLanguageToScript(t *testing.T) {
	cases := []struct {
		lang string
		want string
	}{
		{"ja", "Jpan"},
		{"ja-JP", "Jpan"},
		{"ko", "Hang"},
		{"zh_CN", "Hans"},
		{"zh-SG", "Hans"},
		{"zh-TW", "Hant"},
		{"zh", "Hant"},
		{"ar-SA", "Arab"},
		{"he", "Hebr"},
		{"th", "Thai"},
		{"hi", "Deva"},
		{"ka", "Geor"},
		{"en-US", ""},
		{"fr", ""},
		{"", ""},
		{"   ", ""},
		{"zz", ""},
	}
	for _, tc := range cases {
		if got := LanguageToScript(tc.lang); got != tc.want {
			t.Errorf("LanguageToScript(%q) = %q, want %q", tc.lang, got, tc.want)
		}
	}
}

func TestScriptKindForLanguage(t *testing.T) {
	cases := []struct {
		lang string
		want ScriptKind
	}{
		{"en-US", ScriptLatin},
		{"vi", ScriptComplex},
		{"ja", ScriptEastAsian},
		{"zh-CN", ScriptEastAsian},
		{"ko-KR", ScriptEastAsian},
		{"ar", ScriptComplex},
		{"th", ScriptComplex},
		{"", ScriptLatin},
	}
	for _, tc := range cases {
		if got := ScriptKindForLanguage(tc.lang); got != tc.want {
			t.Errorf("ScriptKindForLanguage(%q) = %v, want %v", tc.lang, got, tc.want)
		}
	}
}
