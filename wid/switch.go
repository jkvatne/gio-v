// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// SwitchDef is the parameters for a slider
type SwitchDef struct {
	Clickable
	th      *Theme
	Value   *bool
	changed bool
	padding layout.Inset
}

// Switch returns a widget for a switch
func Switch(th *Theme, State *bool, handler func(b bool)) func(gtx C) D {
	s := &SwitchDef{}
	s.th = th
	s.SetupTabs()
	s.Value = State
	s.handler = handler
	s.padding = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(5), Right: unit.Dp(5)}
	return func(gtx C) D {
		dims := s.padding.Layout(gtx, func(gtx C) D { return s.Layout(gtx) })
		if handler != nil {
			s.HandleToggle(s.Value, &s.changed)
		} else {
			s.HandleToggle(s.Value, &s.changed)
		}
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {

	// Calculate sizes
	trackWidth := gtx.Px(s.th.TextSize.Scale(2.1))
	trackHeight := gtx.Px(s.th.TextSize.Scale(0.8))
	thumbSize := gtx.Px(s.th.TextSize.Scale(1.05))
	trackOff := float32(thumbSize-trackHeight) * 0.5
	thumbRadius := float32(thumbSize) / 2

	// Find colors
	trackColor := MulAlpha(s.th.Primary, 0x80)
	dotColor := s.th.Primary
	if !*s.Value {
		trackColor = Gray(trackColor)
		dotColor = s.th.Background
	}
	if gtx.Queue == nil {
		dotColor = Disabled(dotColor)
		trackColor = Disabled(trackColor)
	}

	// Draw track.
	trackCorner := float32(trackHeight) / 2
	trackRect := f32.Rectangle{Max: f32.Point{
		X: float32(trackWidth),
		Y: float32(trackHeight),
	}}
	t := op.Offset(f32.Point{Y: trackOff}).Push(gtx.Ops)
	cl := clip.UniformRRect(trackRect, trackCorner).Push(gtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()
	t.Pop()

	// Compute thumb offset based on switch on/off state
	ofs := float32(0)
	if *s.Value {
		ofs = float32(trackWidth - thumbSize)
	}
	st := op.Offset(f32.Point{X: ofs}).Push(gtx.Ops)

	// Draw hover/focused circle
	hoverRadius := float32(2.4 * thumbRadius)
	if s.Hovered() || s.Focused() {
		paint.FillShape(gtx.Ops, MulAlpha(s.th.Primary, 88),
			clip.Ellipse{f32.Point{}, f32.Point{X: hoverRadius, Y: hoverRadius}}.Op(gtx.Ops))
	}

	// Draw thumb shadow, a translucent disc slightly larger than the thumb itself.
	for i := 6; i > 0; i-- {
		s := op.Offset(f32.Point{Y: float32(i) * 0.4}).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, color.NRGBA{A: alpha[i]},
			clip.Ellipse{f32.Point{}, f32.Point{X: hoverRadius, Y: hoverRadius}}.Op(gtx.Ops))
		s.Pop()
	}
	// Draw thumb.
	paint.FillShape(gtx.Ops, dotColor,
		clip.Ellipse{f32.Point{}, f32.Point{X: hoverRadius, Y: hoverRadius}}.Op(gtx.Ops))

	st.Pop()
	// Set area for click and hover
	gtx.Constraints.Min = image.Pt(trackWidth, thumbSize)
	// Handle clicks and keyboard
	s.LayoutClickable(gtx)
	s.HandleClicks(gtx)
	s.HandleKeys(gtx)
	return D{Size: image.Point{X: trackWidth, Y: thumbSize}}
}
