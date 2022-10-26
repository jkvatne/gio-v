// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/op"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// LabelDef is the setup for a label.
type LabelDef struct {
	Base
	// Face defines the text style.
	Font text.Font
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	TextSize unit.Sp
	Stringer func() string
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
		d.TextSize = d.th.TextSize * 1.5
	}
}

// Small makes text 20% smaller.
func Small() LabelOption {
	return func(d *LabelDef) {
		d.TextSize = d.th.TextSize * 0.8
	}
}

// Size set the relative font size (1.0 gives normal text).
func Size(size float32) LabelOption {
	return func(d *LabelDef) {
		d.TextSize = d.th.TextSize * unit.Sp(size)
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
	tl := widget.Label{Alignment: l.Alignment, MaxLines: l.MaxLines}
	dims := tl.Layout(gtx, l.th.Shaper, l.Font, l.TextSize, l.Stringer())
	return dims
}

// Value returns a widget for a value given by stringer function
func Value(th *Theme, s func() string, options ...Option) func(gtx C) D {
	w := LabelDef{
		Stringer:  s,
		TextSize:  th.TextSize,
		Alignment: text.Start,
		Font:      text.Font{Weight: text.Medium, Style: text.Regular},
		MaxLines:  1,
	}
	w.padding = th.LabelPadding
	w.th = th
	w.fgColor = th.Fg(Canvas)
	w.bgColor = th.Bg(Canvas)
	for _, option := range options {
		option.apply(&w)
	}

	return func(gtx C) D {
		macro := op.Record(gtx.Ops)
		dims := w.padding.Layout(gtx, func(gtx C) D {
			return w.Layout(gtx)
		})
		call := macro.Stop()
		defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, w.bgColor)
		call.Add(gtx.Ops)
		return dims
	}
}

// Label returns a widget for a label showing a string
func Label(th *Theme, str string, options ...Option) func(gtx C) D {
	s := func() string { return str }
	return Value(th, s, options...)
}
