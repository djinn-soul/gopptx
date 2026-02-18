package transitions

import (
	"testing"
)

func TestTransitionOptionsValidation(t *testing.T) {
	tests := []struct {
		name    string
		opts    TransitionOptions
		wantErr bool
	}{
		{
			name:    "valid fade",
			opts:    TransitionOptions{Type: TransitionFade},
			wantErr: false,
		},
		{
			name:    "invalid type",
			opts:    TransitionOptions{Type: TransitionType("invalid")},
			wantErr: true,
		},
		{
			name:    "invalid orientation",
			opts:    TransitionOptions{Type: TransitionFade, Orientation: TransitionOrientHorizontal},
			wantErr: true,
		},
		{
			name:    "valid orientation for split",
			opts:    TransitionOptions{Type: TransitionSplit, Orientation: TransitionOrientHorizontal},
			wantErr: false,
		},
		{
			name:    "invalid direction for push",
			opts:    TransitionOptions{Type: TransitionPush, Direction: TransitionDirection("invalid")},
			wantErr: true,
		},
		{
			name:    "valid direction for push",
			opts:    TransitionOptions{Type: TransitionPush, Direction: TransitionDirUp},
			wantErr: false,
		},
		{
			name:    "valid sound",
			opts:    TransitionOptions{Type: TransitionFade, Sound: &TransitionSound{RelID: "rId2"}},
			wantErr: false,
		},
		{
			name:    "invalid sound missing relID",
			opts:    TransitionOptions{Type: TransitionFade, Sound: &TransitionSound{RelID: ""}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.opts.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("TransitionOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransitionXML(t *testing.T) {
	tests := []struct {
		name string
		opts TransitionOptions
		want string
	}{
		{
			name: "fade",
			opts: TransitionOptions{Type: TransitionFade, DurationMS: 2000},
			want: `<p:transition><p:fade/></p:transition>`,
		},
		{
			name: "push right",
			opts: TransitionOptions{Type: TransitionPush, Direction: TransitionDirRight},
			want: `<p:transition><p:push dir="r"/></p:transition>`,
		},
		{
			name: "shape serializes as wheel",
			opts: TransitionOptions{Type: TransitionShape},
			want: `<p:transition><p:wheel/></p:transition>`,
		},
		{
			name: "split horizontal out",
			opts: TransitionOptions{
				Type:        TransitionSplit,
				Direction:   TransitionDirOut,
				Orientation: TransitionOrientHorizontal,
			},
			want: `<p:transition><p:split dir="out" orient="horz"/></p:transition>`,
		},
		{
			name: "fade with sound",
			opts: TransitionOptions{
				Type:  TransitionFade,
				Sound: &TransitionSound{RelID: "rId2", Name: "Applause.wav"},
			},
			want: `<p:transition><p:fade/><p:sndAc><p:stSnd><p:snd r:embed="rId2" name="Applause.wav"/></p:stSnd></p:sndAc></p:transition>`,
		},
		{
			name: "fade with looping sound",
			opts: TransitionOptions{
				Type:  TransitionFade,
				Sound: &TransitionSound{RelID: "rId3", Loop: true},
			},
			want: `<p:transition><p:fade/><p:sndAc><p:stSnd loop="1"><p:snd r:embed="rId3"/></p:stSnd></p:sndAc></p:transition>`,
		},
		{
			name: "cut with sound",
			opts: TransitionOptions{
				Type:  TransitionCut,
				Sound: &TransitionSound{RelID: "rId4"},
			},
			want: `<p:transition><p:cut/><p:sndAc><p:stSnd><p:snd r:embed="rId4"/></p:stSnd></p:sndAc></p:transition>`,
		},
		{
			name: "sound with special characters",
			opts: TransitionOptions{
				Type:  TransitionFade,
				Sound: &TransitionSound{RelID: "rId5", Name: `Applause "Special".wav`},
			},
			want: `<p:transition><p:fade/><p:sndAc><p:stSnd><p:snd r:embed="rId5" name="Applause &quot;Special&quot;.wav"/></p:stSnd></p:sndAc></p:transition>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.XML()
			if got != tt.want {
				t.Errorf("TransitionOptions.XML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransitionXMLMorph(t *testing.T) {
	tests := []struct {
		name string
		opts TransitionOptions
		want string
	}{
		{
			name: "default morph",
			opts: TransitionOptions{Type: TransitionMorph},
			want: `<mc:AlternateContent xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"><mc:Choice xmlns:p159="http://schemas.microsoft.com/office/powerpoint/2015/09/main" Requires="p159"><p:transition spd="slow" xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" p14:dur="2000"><p159:morph option="byObject"/><p:extLst mod="1"><p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}"><p15:morph xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2015/09/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns=""/></p:ext></p:extLst></p:transition></mc:Choice><mc:Fallback><p:transition spd="slow"><p:fade/><p:extLst mod="1"><p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}"><p15:morph xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2015/09/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns=""/></p:ext></p:extLst></p:transition></mc:Fallback></mc:AlternateContent>`,
		},
		{
			name: "morph word option",
			opts: TransitionOptions{Type: TransitionMorph, MorphOption: MorphOptionWord},
			want: `<mc:AlternateContent xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"><mc:Choice xmlns:p159="http://schemas.microsoft.com/office/powerpoint/2015/09/main" Requires="p159"><p:transition spd="slow" xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" p14:dur="2000"><p159:morph option="byWord"/><p:extLst mod="1"><p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}"><p15:morph xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2015/09/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns=""/></p:ext></p:extLst></p:transition></mc:Choice><mc:Fallback><p:transition spd="slow"><p:fade/><p:extLst mod="1"><p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}"><p15:morph xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2015/09/main" xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns=""/></p:ext></p:extLst></p:transition></mc:Fallback></mc:AlternateContent>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.XML()
			if got != tt.want {
				t.Errorf("TransitionOptions.XML() = %v, want %v", got, tt.want)
			}
		})
	}
}
