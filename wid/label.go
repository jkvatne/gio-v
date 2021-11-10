// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

// LabelDef is the setup for a label.
type LabelDef struct {
	Widget
	// Face defines the text style.
	Font text.Font
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Text     string
	TextSize unit.Value
	padding  layout.Inset
	shaper   text.Shaper
}

// Label  returns a widget for a label.
func Label(th *Theme, str string, options ...Option) func(gtx C) D {
	w := LabelDef{
		Text:      str,
		TextSize:  th.TextSize,
		shaper:    th.Shaper,
		Alignment: text.Start,
		Font:      text.Font{Weight: text.Medium, Style: text.Regular},
		padding:   th.LabelPadding,
	}
	w.th = th
	w.fgColor = th.OnBackground
	for _, option := range options {
		option.apply(&w)
	}
	return func(gtx C) D {
		return w.padding.Layout(gtx, func(gtx C) D {
			return w.Layout(gtx)
		})
	}
}

// LabelOption is options specific to Edits.
type LabelOption func(w *LabelDef)

// Bold is an option parameter to set the widget hint (tooltip).
func Bold() LabelOption {
	return func(d *LabelDef) {
		d.Font.Weight = text.Bold
	}
}

// Weight sets the font weight.
func Weight(weight text.Weight) LabelOption {
	return func(d *LabelDef) {
		d.Font.Weight = weight
	}
}

// Large makes text 50% larger.
func Large() LabelOption {
	return func(d *LabelDef) {
		d.TextSize = d.th.TextSize.Scale(1.5)
	}
}

// Small makes text 20% smaller.
func Small() LabelOption {
	return func(d *LabelDef) {
		d.TextSize = d.th.TextSize.Scale(0.8)
	}
}

// Size set the relative font size (1.0 gives normal text).
func Size(size float32) LabelOption {
	return func(d *LabelDef) {
		d.TextSize = d.th.TextSize.Scale(size)
	}
}

// Middle will align text in the middle.
func Middle() LabelOption {
	return func(d *LabelDef) {
		d.Alignment = text.Middle
	}
}

// End will align text to the end.
func End() LabelOption {
	return func(d *LabelDef) {
		d.Alignment = text.End
	}
}

func (e LabelOption) apply(cfg interface{}) {
	e(cfg.(*LabelDef))
}

// Layout will draw the label
func (l LabelDef) Layout(gtx C) D {
	paint.ColorOp{Color: l.fgColor}.Add(gtx.Ops)
	tl := aLabel{Alignment: l.Alignment, MaxLines: l.MaxLines}
	return tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text)
}
