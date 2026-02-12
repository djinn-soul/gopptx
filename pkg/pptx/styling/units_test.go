package styling

import (
	"testing"
)

func TestLengthConversions(t *testing.T) {
	tests := []struct {
		name    string
		length  Length
		wantIn  float64
		wantCm  float64
		wantPt  float64
		wantEmu int64
	}{
		{
			name:    "1 inch",
			length:  Inches(1),
			wantIn:  1.0,
			wantCm:  2.54,
			wantPt:  72.0,
			wantEmu: 914400,
		},
		{
			name:    "1 cm",
			length:  Centimeters(1),
			wantIn:  1.0 / 2.54,
			wantCm:  1.0,
			wantPt:  72.0 / 2.54,
			wantEmu: 360000,
		},
		{
			name:    "72 pt",
			length:  Points(72),
			wantIn:  1.0,
			wantCm:  2.54,
			wantPt:  72.0,
			wantEmu: 914400,
		},
		{
			name:    "10000 emu",
			length:  Emu(10000),
			wantIn:  10000.0 / 914400.0,
			wantCm:  10000.0 / 360000.0,
			wantPt:  10000.0 / 12700.0,
			wantEmu: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.length.Inches(); got != tt.wantIn {
				t.Errorf("Length.Inches() = %v, want %v", got, tt.wantIn)
			}
			if got := tt.length.Cm(); got != tt.wantCm {
				t.Errorf("Length.Cm() = %v, want %v", got, tt.wantCm)
			}
			if got := tt.length.Pt(); got != tt.wantPt {
				t.Errorf("Length.Pt() = %v, want %v", got, tt.wantPt)
			}
			if got := tt.length.Emu(); got != tt.wantEmu {
				t.Errorf("Length.Emu() = %v, want %v", got, tt.wantEmu)
			}
		})
	}
}

func TestLengthClamping(t *testing.T) {
	if got := Inches(10000); got != MaxEMU {
		t.Errorf("Inches(10000) = %v, want %v (MaxEMU)", got, MaxEMU)
	}
	if got := Inches(-10000); got != -MaxEMU {
		t.Errorf("Inches(-10000) = %v, want %v (-MaxEMU)", got, -MaxEMU)
	}
}

func TestFontSize(t *testing.T) {
	tests := []struct {
		pt   float64
		want int
	}{
		{12.0, 1200},
		{10.5, 1050},
		{0.0, 0},
	}
	for _, tt := range tests {
		if got := FontSize(tt.pt); got != tt.want {
			t.Errorf("FontSize(%v) = %v, want %v", tt.pt, got, tt.want)
		}
	}
}
