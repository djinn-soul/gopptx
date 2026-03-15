package action

import "testing"

func TestValidateHyperlinkAction_FileAndProgramSchemeChecks(t *testing.T) {
	tests := []struct {
		name    string
		action  HyperlinkAction
		wantErr bool
	}{
		{
			name:    "valid file uri",
			action:  HyperlinkFile("file:///C:/docs/report.xlsx"),
			wantErr: false,
		},
		{
			name:    "valid relative file",
			action:  HyperlinkFile("docs/report.xlsx"),
			wantErr: false,
		},
		{
			name:    "valid windows drive path accepted cross platform",
			action:  HyperlinkFile(`D:\docs\report.xlsx`),
			wantErr: false,
		},
		{
			name:    "invalid non-file scheme",
			action:  HyperlinkFile("https://example.com/report.xlsx"),
			wantErr: true,
		},
		{
			name:    "escaped traversal rejected",
			action:  HyperlinkFile("file:///C:/safe/%2e%2e/secret.txt"),
			wantErr: true,
		},
		{
			name:    "restricted program path rejected",
			action:  HyperlinkProgram(`C:\Windows\System32\calc.exe`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHyperlinkAction(tt.action, "test")
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateHyperlinkAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
