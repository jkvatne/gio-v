// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

// LabelDef is the setup for a label
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
	padding  layout.Inset
	shaper   text.Shaper
}

// CreateLabelDef will make a LabelDef
func CreateLabelDef(th *Theme, text string, align text.Alignment, relSize float32) LabelDef {
	return LabelDef{
		Text:      text,
		Color:     th.OnBackground,
		TextSize:  th.TextSize.Scale(relSize),
		shaper:    th.Shaper,
		Alignment: align,
	}
}

// Label  returns a widget for a label
func Label(th *Theme, text string, align text.Alignment, relSize float32) func(gtx C) D {
	lbl := LabelDef{
		Text:      text,
		Color:     th.OnBackground,
		TextSize:  th.TextSize.Scale(relSize),
		shaper:    th.Shaper,
		Alignment: align,
	}
	lbl.padding = th.LabelPadding
	return func(gtx C) D {
		return lbl.padding.Layout(gtx, func(gtx C) D {
			return lbl.Layout(gtx)
		})
	}
}

// Layout will draw the label
func (l LabelDef) Layout(gtx C) D {
	paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
	tl := aLabel{Alignment: l.Alignment, MaxLines: l.MaxLines}
	return tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text)
}
