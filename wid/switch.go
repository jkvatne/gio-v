// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/io/semantic"

	"gioui.org/widget"

	"gioui.org/io/pointer"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// SwitchDef is the parameters for a slider
type SwitchDef struct {
	sw       widget.Bool
	th       *Theme
	StatePtr *bool
	padding  layout.Inset
	handler  func(b bool)
}

// Switch returns a widget for a switch
func Switch(th *Theme, statePtr *bool, handler func(b bool)) func(gtx C) D {
	s := &SwitchDef{}
	s.th = th
	s.StatePtr = statePtr
	s.handler = handler
	s.padding = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(5), Right: unit.Dp(5)}
	return func(gtx C) D {
		// semantic.Switch.Add(gtx.Ops)
		// semantic.SelectedOp(*s.Value).Add(gtx.Ops)
		// semantic.DisabledOp(gtx.Queue == nil).Add(gtx.Ops)
		dims := s.padding.Layout(gtx, func(gtx C) D { return s.Layout(gtx) })
		if handler != nil {
			s.handler(s.sw.Value)
		}
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {

	// Calculate sizes
	trackWidth := gtx.Sp(s.th.TextSize * 2.3)
	trackHeight := gtx.Sp(s.th.TextSize * 1.2)
	thumbSize := gtx.Sp(s.th.TextSize * 1.3)
	trackOff := (thumbSize - trackHeight) / 2

	// Find colors
	trackColor := MulAlpha(s.th.Primary, 0x80)
	dotColor := s.th.Primary
	ofs := trackWidth - thumbSize
	if s.sw.Changed() {
		*s.StatePtr = s.sw.Value
	} else {
		s.sw.Value = *s.StatePtr
	}
	if !*s.StatePtr {
		trackColor = Gray(trackColor)
		dotColor = s.th.Background
		ofs = 0
	}
	if gtx.Queue == nil {
		dotColor = Disabled(dotColor)
		trackColor = Disabled(trackColor)
	}

	// Draw track.
	trackCorner := trackHeight / 2
	trackRect := image.Rect(0, 0, trackWidth, trackHeight)
	t := op.Offset(image.Point{Y: trackOff}).Push(gtx.Ops)
	cl := clip.UniformRRect(trackRect, trackCorner).Push(gtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()
	t.Pop()

	st := op.Offset(image.Point{X: ofs}).Push(gtx.Ops)
	// Draw hover.
	if s.sw.Hovered() || s.sw.Focused() {
		r := thumbSize / 4
		o := op.Offset(image.Point{X: -r / 2, Y: -r / 2}).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, MulAlpha(s.th.Primary, 88),
			clip.Ellipse{image.Point{}, image.Point{X: thumbSize + r, Y: thumbSize + r}}.Op(gtx.Ops))
		o.Pop()
	}
	// Draw thumb outline
	paint.FillShape(gtx.Ops, s.th.OnBackground,
		clip.Ellipse{image.Point{}, image.Point{X: thumbSize, Y: thumbSize}}.Op(gtx.Ops))
	// Draw thumb inside
	o := op.Offset(image.Point{X: 2, Y: 2}).Push(gtx.Ops)
	paint.FillShape(gtx.Ops, dotColor,
		clip.Ellipse{image.Point{}, image.Point{X: thumbSize - 4, Y: thumbSize - 4}}.Op(gtx.Ops))
	o.Pop()
	st.Pop()

	// Set up click area.
	clickSize := gtx.Dp(40)
	clickOff := image.Point{
		X: (thumbSize - clickSize) / 2,
		Y: (trackHeight-clickSize)/2 + trackOff,
	}
	defer op.Offset(clickOff).Push(gtx.Ops).Pop()
	sz := image.Pt(clickSize, clickSize)
	defer clip.Ellipse(image.Rectangle{Max: sz}).Push(gtx.Ops).Pop()
	s.sw.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Switch.Add(gtx.Ops)
		return layout.Dimensions{Size: sz}
	})

	dims := image.Point{X: trackWidth, Y: thumbSize}
	return layout.Dimensions{Size: dims}
}
