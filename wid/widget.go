// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

// Widget is tha base structure for widgets. It contains variables that (almost) all widgets share
type Widget struct {
	th      *Theme
	hint    string
	padding layout.Inset
	width   unit.Value
	fgColor color.NRGBA
}

// WidgetIf is the interface functions for widgets, used by options to set parameters
type WidgetIf interface {
	setWidth(width float32)
	setHint(hint string)
	setPadding(padding layout.Inset)
	setColor(c color.NRGBA)
}

// WidgetOption is a type for optional parameters when creating widgets
type WidgetOption func(WidgetIf)

// Option is the interface for optional parameters
type Option interface {
	apply(cfg interface{})
}

// Apply will apply all optional parameters. This can only be used when the widget has no own options.
func (wid *Widget) Apply(options ...Option) {
	for _, option := range options {
		option.apply(wid)
	}
}

func (wid WidgetOption) apply(cfg interface{}) {
	cc := cfg.(WidgetIf)
	wid(cc)
}

func (wid *Widget) setWidth(width float32) {
	wid.width = unit.Dp(width)
}

func (wid *Widget) setHint(hint string) {
	wid.hint = hint
}

func (wid *Widget) setPadding(padding layout.Inset) {
	wid.padding = padding
}

func (wid *Widget) setColor(c color.NRGBA) {
	wid.fgColor = c
}

// Pad is used to set default widget paddings
func (wid *Widget) Pad(t, r, b, l float32) {
	wid.padding = layout.Inset{Top: unit.Dp(t), Bottom: unit.Dp(b), Left: unit.Dp(l), Right: unit.Dp(r)}
}

// W is the option parameter for setting widget width
func W(width float32) WidgetOption {
	return func(w WidgetIf) {
		w.setWidth(width)
	}
}

// Max is an option parameter to set the widget width to fill all avaiable space
func Max() WidgetOption {
	return func(w WidgetIf) {
		w.setWidth(10000)
	}
}

// Min is an option parameter to set the widget to its minimum width (i.e. length of text)
func Min() WidgetOption {
	return func(w WidgetIf) {
		w.setWidth(0)
	}
}

// Hint is an option parameter to set the widget hint (tooltip)
func Hint(hint string) WidgetOption {
	return func(w WidgetIf) {
		w.setHint(hint)
	}
}

// Color is an option parameter to set widget color
func Color(c color.NRGBA) WidgetOption {
	return func(w WidgetIf) {
		w.setColor(c)
	}
}

// Pad is an option parameter to set customized padding. Noe that 1,2,3 or 4 paddings can be specified.
// If 1 is supplide, it is used for left,right,top,bottom, all with the same padding
// If 2 is supplied, the first is used for top/bottom, and the second for left and riht padding
// If 4 is supplied, it is used for top, right, bottom, left in that sequence.
// All values are in Dp (float32 device independent pixels)
func Pad(pads ...float32) WidgetOption {
	return func(w WidgetIf) {
		switch len(pads) {
		case 0:
			w.setPadding(layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(4), Right: unit.Dp(4)})
		case 1:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[0]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[0])})
		case 2:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[1])})
		case 3:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[2])})
		case 4:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[2]), Left: unit.Dp(pads[3])})
		}
	}
}

// CalcMin will calculate the minimum size of widget
func CalcMin(gtx C, width unit.Value) image.Point {
	min := gtx.Constraints.Min
	if width.V <= 1.0 {
		min.X = gtx.Px(width.Scale(float32(gtx.Constraints.Max.X)))
	} else if min.X < gtx.Px(width) {
		min.X = gtx.Px(width)
	}
	if min.X > gtx.Constraints.Max.X {
		min.X = gtx.Constraints.Max.X
	}
	return min
}
