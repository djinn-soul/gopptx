package transitions

import (
	"fmt"
	"strings"
)

func (o TransitionOptions) XML() string {
	if o.Type == TransitionNone || (o.Type == TransitionCut && o.Sound == nil) {
		return ""
	}
	if o.Type == TransitionMorph {
		return o.morphXML()
	}

	var b strings.Builder
	b.WriteString(`<p:transition`)
	if o.DisableAdvanceOnClick {
		b.WriteString(` advClick="0"`)
	}
	if o.AdvanceAfterMS > 0 {
		fmt.Fprintf(&b, ` advTm="%d"`, o.AdvanceAfterMS)
	}
	if o.ThruBlk {
		b.WriteString(` thruBlk="1"`)
	}
	b.WriteString(`>`)
	b.WriteString(`<p:`)
	b.WriteString(o.Type.transitionElementName())

	if o.Direction != "" {
		fmt.Fprintf(&b, ` dir="%s"`, o.Direction)
	}
	if o.Orientation != "" {
		fmt.Fprintf(&b, ` orient="%s"`, o.Orientation)
	}
	if o.SpokeCount > 0 {
		fmt.Fprintf(&b, ` spokes="%d"`, o.SpokeCount)
	}
	if o.ThruBlk {
		b.WriteString(` thruBlk="1"`)
	}
	b.WriteString(`/>`)

	b.WriteString(transitionSoundXML(o.Sound))

	b.WriteString(`</p:transition>`)
	return b.String()
}

func (o TransitionOptions) morphXML() string {
	var b strings.Builder
	speed := "slow"
	durationMS := o.DurationMS
	if durationMS == 0 {
		durationMS = 2000
	}
	soundXML := transitionSoundXML(o.Sound)
	choiceOption := morphChoiceOption(o.MorphOption)
	const morphExt = `<p:extLst mod="1"><p:ext uri="{AE3914FA-7E93-4B9E-9A96-D1E12CAF14E6}">` +
		`<p15:morph xmlns:p15="http://schemas.microsoft.com/office/powerpoint/2015/09/main" ` +
		`xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns=""/>` +
		`</p:ext></p:extLst>`

	b.WriteString(`<mc:AlternateContent xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006">`)
	b.WriteString(
		`<mc:Choice xmlns:p159="http://schemas.microsoft.com/office/powerpoint/2015/09/main" Requires="p159">`,
	)
	fmt.Fprintf(
		&b,
		`<p:transition spd="%s" xmlns:p14="http://schemas.microsoft.com/office/powerpoint/2010/main" p14:dur="%d"`,
		speed,
		durationMS,
	)
	if o.DisableAdvanceOnClick {
		b.WriteString(` advClick="0"`)
	}
	if o.AdvanceAfterMS > 0 {
		fmt.Fprintf(&b, ` advTm="%d"`, o.AdvanceAfterMS)
	}
	b.WriteString(`>`)
	fmt.Fprintf(&b, `<p159:morph option="%s"/>`, choiceOption)
	b.WriteString(morphExt)
	b.WriteString(soundXML)
	b.WriteString(`</p:transition></mc:Choice>`)
	fmt.Fprintf(&b, `<mc:Fallback><p:transition spd="%s"`, speed)
	if o.DisableAdvanceOnClick {
		b.WriteString(` advClick="0"`)
	}
	if o.AdvanceAfterMS > 0 {
		fmt.Fprintf(&b, ` advTm="%d"`, o.AdvanceAfterMS)
	}
	b.WriteString(`><p:fade/>`)
	b.WriteString(morphExt)
	b.WriteString(soundXML)
	b.WriteString(`</p:transition></mc:Fallback></mc:AlternateContent>`)
	return b.String()
}

func morphChoiceOption(option MorphOption) string {
	switch option {
	case MorphOptionWord:
		return "byWord"
	case MorphOptionCharacter:
		return "byChar"
	default:
		return "byObject"
	}
}

func transitionSoundXML(sound *TransitionSound) string {
	if sound == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<p:sndAc><p:stSnd`)
	if sound.Loop {
		b.WriteString(` loop="1"`)
	}
	b.WriteString(`><p:snd r:embed="` + escape(sound.RelID) + `"`)
	if sound.Name != "" {
		b.WriteString(` name="` + escape(sound.Name) + `"`)
	}
	b.WriteString(`/>`)
	b.WriteString(`</p:stSnd></p:sndAc>`)
	return b.String()
}

func escape(value string) string {
	if !strings.ContainsAny(value, `&<>"'`) {
		return value
	}

	var b strings.Builder
	b.Grow(len(value))
	for _, r := range value {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (t TransitionType) XML() string {
	switch t {
	case TransitionNone, TransitionCut:
		// TransitionCut is the default and requires no XML unless options (like sound) are set.
		return ""
	case TransitionPush:
		return `<p:transition><p:push dir="r"/></p:transition>`
	case TransitionWipe:
		return `<p:transition><p:wipe dir="r"/></p:transition>`
	case TransitionSplit:
		return `<p:transition><p:split dir="out" orient="horz"/></p:transition>`
	case TransitionZoom:
		return `<p:transition><p:zoom dir="in"/></p:transition>`
	case TransitionFade:
		return `<p:transition><p:fade/></p:transition>`
	case TransitionReveal:
		return `<p:transition><p:reveal dir="r"/></p:transition>`
	case TransitionCover:
		return `<p:transition><p:cover dir="r"/></p:transition>`
	default:
		return fmt.Sprintf(`<p:transition><p:%s/></p:transition>`, t.transitionElementName())
	}
}

func (t TransitionType) transitionElementName() string {
	if t == TransitionShape {
		return string(TransitionClock)
	}
	return string(t)
}
