package elements

import (
	"testing"

	"github.com/djinn-soul/gopptx/pkg/pptx/smartart"
)

func TestValidateSlideSmartArt(t *testing.T) {
	tests := []struct {
		name    string
		slide   SlideContent
		wantErr bool
	}{
		{
			name: "Valid SmartArt",
			slide: SlideContent{
				Title:  "Valid Slide",
				SmartArtDiagrams: []smartart.SmartArt{
					smartart.NewSmartArt(smartart.BasicBlockList).
						AddNode(smartart.NewNode("Node 1")),
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid SmartArt - No Nodes",
			slide: SlideContent{
				Title:  "Invalid Slide",
				SmartArtDiagrams: []smartart.SmartArt{
					smartart.NewSmartArt(smartart.BasicBlockList),
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid SmartArt - Empty Node Text",
			slide: SlideContent{
				Title:  "Invalid Slide",
				SmartArtDiagrams: []smartart.SmartArt{
					smartart.NewSmartArt(smartart.BasicBlockList).
						AddNode(smartart.NewNode("")),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.slide.Validate(1); (err != nil) != tt.wantErr {
				t.Errorf("SlideContent.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
