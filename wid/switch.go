// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

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
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {

	// Calculate sizes
	trackWidth := gtx.Sp(s.th.TextSize * 2.1)
	trackHeight := gtx.Sp(s.th.TextSize * 0.8)
	thumbSize := gtx.Sp(s.th.TextSize * 1.05)
	trackOff := (thumbSize - trackHeight) / 2
	thumbRadius := thumbSize / 2

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
	trackCorner := trackHeight / 2
	trackRect := image.Rect(0, 0, trackWidth, trackHeight)
	t := op.Offset(image.Point{Y: trackOff}).Push(gtx.Ops)
	cl := clip.UniformRRect(trackRect, trackCorner).Push(gtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()
	t.Pop()

	// Compute thumb offset based on switch on/off state
	ofs := 0
	if *s.Value {
		ofs = trackWidth - thumbSize
	}
	st := op.Offset(image.Point{X: ofs}).Push(gtx.Ops)

	// Draw hover/focused circle
	hoverRadius := 2 * thumbRadius
	if s.Hovered() || s.Focused() {
		paint.FillShape(gtx.Ops, MulAlpha(s.th.Primary, 88),
			clip.Ellipse{image.Point{}, image.Point{X: hoverRadius, Y: hoverRadius}}.Op(gtx.Ops))
	}

	// Draw thumb shadow, a translucent disc slightly larger than the thumb itself.
	for i := 6; i > 0; i-- {
		s := op.Offset(image.Point{Y: i / 2}).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, color.NRGBA{A: alpha[i]},
			clip.Ellipse{image.Point{}, image.Point{X: hoverRadius, Y: hoverRadius}}.Op(gtx.Ops))
		s.Pop()
	}
	// Draw thumb.
	paint.FillShape(gtx.Ops, dotColor,
		clip.Ellipse{image.Point{}, image.Point{X: hoverRadius, Y: hoverRadius}}.Op(gtx.Ops))

	st.Pop()
	// Set area for click and hover
	gtx.Constraints.Min = image.Pt(trackWidth, thumbSize)
	// Handle clicks and keyboard
	s.LayoutClickable(gtx)
	s.HandleClicks(gtx)
	s.HandleKeys(gtx)
	return D{Size: image.Point{X: trackWidth, Y: thumbSize}}
}
