// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type LabelDef struct {
	// Face defines the text style.
	Font text.Font
	// Color is the text color.
	Color color.NRGBA
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Text     string
	TextSize unit.Value

	shaper text.Shaper
}

func CreateLabelDef(th *Theme, text string, align text.Alignment, relSize float32) LabelDef {
	return LabelDef{
		Text:     text,
		Color:    th.Palette.OnBackground,
		TextSize: th.TextSize.Scale(relSize),
		shaper:   th.Shaper,
		Alignment: align,
	}
}

func Label(th *Theme, text string, align text.Alignment, relSize float32) func(gtx C) D {
	if relSize == 0.0 {
		relSize = 1.0
	}
	return func(gtx C) D {
		return LabelDef{
			Text:     text,
			Color:    th.Palette.OnBackground,
			TextSize: th.TextSize.Scale(relSize),
			shaper:   th.Shaper,
			Alignment: align,
		}.Layout(gtx)
	}
}

func (l LabelDef) Layout(gtx layout.Context) layout.Dimensions {
	paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
		tl := aLabel{Alignment: l.Alignment, MaxLines: l.MaxLines}
	return tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text)
}
