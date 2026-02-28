package export

import "testing"

func TestParsePDFDriver(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PDFDriver
		wantErr bool
	}{
		{name: "default empty", input: "", want: PDFDriverAuto},
		{name: "auto", input: "auto", want: PDFDriverAuto},
		{name: "native", input: "native", want: PDFDriverNative},
		{name: "libreoffice", input: "libreoffice", want: PDFDriverLibreOffice},
		{name: "powerpoint", input: "powerpoint", want: PDFDriverPowerPoint},
		{name: "mixed case", input: "NaTiVe", want: PDFDriverNative},
		{name: "invalid", input: "chrome", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePDFDriver(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParsePDFDriver(%q) expected error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePDFDriver(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ParsePDFDriver(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
