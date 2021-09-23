// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gio-v/f32color"
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

type SwitchDef struct {
	Clickable
	Color struct {
		Enabled  color.NRGBA
		Disabled color.NRGBA
		Track    color.NRGBA
	}
	size unit.Value
	Value bool
	changed bool
}

func Switch(th *Theme, initialState bool, handler func(b bool)) func(gtx C) D {
	s := &SwitchDef{}
	s.SetupTabs()
	s.Color.Enabled = th.Palette.Primary
	s.Color.Disabled = th.Palette.Background
	s.Color.Track = f32color.MulAlpha(th.Palette.Primary, 0x88)
	s.size = th.TextSize
	s.Value = initialState
	s.handler = handler
	return func(gtx C) D {
		dims := s.Layout(gtx)
		if handler!=nil {
			s.HandleToggle(&s.Value, &s.changed)
		} else {
			s.HandleToggle(&s.Value, &s.changed)
		}
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {
	trackWidth := gtx.Px(s.size.Scale(2.2))
	trackHeight := gtx.Px(s.size.Scale(1.2))
	thumbSize := gtx.Px(s.size.Scale(1.0))
	trackOff := float32(thumbSize-trackHeight) * .4

	// Draw track.
	stack := op.Save(gtx.Ops)
	trackCorner := float32(trackHeight) / 2
	trackRect := f32.Rectangle{Max: f32.Point{
		X: float32(trackWidth),
		Y: float32(trackHeight),
	}}
	col := s.Color.Disabled
	if s.Value {
		col = s.Color.Enabled
	}
	if gtx.Queue == nil {
		col = f32color.Disabled(col)
	}
	trackColor := s.Color.Track
	op.Offset(f32.Point{Y: trackOff}).Add(gtx.Ops)
	clip.UniformRRect(trackRect, trackCorner).Add(gtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Load()

	// Compute thumb offset and color.
	stack = op.Save(gtx.Ops)
	if s.Value {
		off := trackWidth - thumbSize
		op.Offset(f32.Point{X: float32(off)}).Add(gtx.Ops)
	}

	thumbRadius := float32(thumbSize) / 2

	// Draw hover.
	if s.Hovered() || s.Focused() {
		r := 1.7 * thumbRadius
		background := f32color.MulAlpha(s.Color.Enabled, 70)
		paint.FillShape(gtx.Ops, background,
			clip.Circle{
				Center: f32.Point{X: thumbRadius, Y: thumbRadius},
				Radius: r,
			}.Op(gtx.Ops))
	}

	// Draw thumb shadow, a translucent disc slightly larger than the thumb itself.
	// Center shadow horizontally and slightly adjust its Y.
	paint.FillShape(gtx.Ops, ARGB(0x55000000),
		clip.Circle{
			Center: f32.Point{X: thumbRadius, Y: thumbRadius + .25},
			Radius: thumbRadius + 1,
		}.Op(gtx.Ops))

	// Draw thumb.
	paint.FillShape(gtx.Ops, col,
		clip.Circle{
			Center: f32.Point{X: thumbRadius, Y: thumbRadius},
			Radius: thumbRadius,
		}.Op(gtx.Ops))
	stack.Load()

	gtx.Constraints.Min = image.Pt(trackWidth, trackHeight)
	s.LayoutClickable(gtx)
	dims := image.Point{X: trackWidth, Y: thumbSize}
	return layout.Dimensions{Size: dims}
}
