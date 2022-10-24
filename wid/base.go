// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
)

type UiState uint8

const (
	DisabledState UiState = iota
	EnabledState
	HoveredState
	FocusedState
	PressedState
)

// Base is tha base structure for widgets. It contains variables that (almost) all widgets share
type Base struct {
	th           *Theme
	state        UiState
	hint         string
	padding      layout.Inset
	onChange     func()
	disabled     bool
	disabler     *bool
	width        unit.Dp
	role         UIRole
	cornerRadius unit.Dp
	fgColor      color.NRGBA
	bgColor      color.NRGBA
}

// BaseIf is the interface functions for widgets, used by options to set parameters
type BaseIf interface {
	setWidth(width float32)
	setHint(hint string)
	setPadding(padding layout.Inset)
	setRole(role UIRole)
	setBgColor(c color.NRGBA)
	setFgColor(c color.NRGBA)
	setHandler(h func())
	getTheme() *Theme
}

// BaseOption is a type for optional parameters when creating widgets
type BaseOption func(BaseIf)

// Option is the interface for optional parameters
type Option interface {
	apply(cfg interface{})
}

// Apply will apply all optional parameters. This can only be used when the widget has no own options.
func (wid *Base) Apply(options ...Option) {
	for _, option := range options {
		option.apply(wid)
	}
}

func (wid BaseOption) apply(cfg interface{}) {
	cc := cfg.(BaseIf)
	wid(cc)
}

func (wid *Base) getTheme() *Theme {
	return wid.th
}

func (wid *Base) setWidth(width float32) {
	wid.width = unit.Dp(width)
}

func (wid *Base) setHint(hint string) {
	wid.hint = hint
}

func (wid *Base) setRole(role UIRole) {
	wid.role = role
}

func (wid *Base) setPadding(padding layout.Inset) {
	wid.padding = padding
}

func (wid *Base) setFgColor(c color.NRGBA) {
	wid.fgColor = c
}

func (wid *Base) setBgColor(c color.NRGBA) {
	wid.bgColor = c
}

func (wid *Base) setHandler(h func()) {
	wid.onChange = h
}

// Pad is used to set default widget paddings
func (wid *Base) Pad(t, r, b, l float32) {
	wid.padding = layout.Inset{Top: unit.Dp(t), Bottom: unit.Dp(b), Left: unit.Dp(l), Right: unit.Dp(r)}
}

// Do is an optional parameter to set a callback when widget state changes
func Do(f func()) BaseOption {
	return func(w BaseIf) {
		w.setHandler(f)
	}
}

// W is the option parameter for setting widget width
func W(width float32) BaseOption {
	return func(w BaseIf) {
		w.setWidth(width)
	}
}

// Max is an option parameter to set the widget width to fill all available space
func Max() BaseOption {
	return func(w BaseIf) {
		w.setWidth(10000)
	}
}

// Min is an option parameter to set the widget to its minimum width (i.e. length of text)
func Min() BaseOption {
	return func(w BaseIf) {
		w.setWidth(0)
	}
}

// Hint is an option parameter to set the widget hint (tooltip)
func Hint(hint string) BaseOption {
	return func(w BaseIf) {
		w.setHint(hint)
	}
}

// Fg is an option parameter to set widget foreground color
func Fg(c color.NRGBA) BaseOption {
	return func(w BaseIf) {
		w.setFgColor(c)
	}
}

// Bg is an option parameter to set widget background color
func Bg(c color.NRGBA) BaseOption {
	return func(w BaseIf) {
		w.setBgColor(c)
	}
}

func Role(r UIRole) BaseOption {
	return func(w BaseIf) {
		w.setBgColor(w.getTheme().Bg(r))
		w.setFgColor(w.getTheme().Fg(r))
		w.setRole(r)
	}
}

// Pads is an option parameter to set customized padding. Noe that 1,2,3 or 4 paddings can be specified.
// If 1 is supplied, it is used for left,right,top,bottom, all with the same padding
// If 2 is supplied, the first is used for top/bottom, and the second for left and right padding
// If 4 is supplied, it is used for top, right, bottom, left in that sequence.
// All values are in Dp (float32 device independent pixels)
func Pads(pads ...float32) BaseOption {
	return func(w BaseIf) {
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
func CalcMin(gtx C, width unit.Dp) image.Point {
	min := gtx.Constraints.Min
	if width <= 1.0 {
		min.X = gtx.Dp(width * unit.Dp(gtx.Constraints.Max.X))
	} else if width != 0 {
		min.X = gtx.Dp(width)
	}
	if min.X > gtx.Constraints.Max.X {
		min.X = gtx.Constraints.Max.X
	}
	return min
}
