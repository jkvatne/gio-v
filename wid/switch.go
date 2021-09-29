// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
)

type SwitchDef struct {
	Clickable
	th *Theme
	size    unit.Value
	Value   bool
	changed bool
	padding layout.Inset
}

func Switch(th *Theme, initialState bool, handler func(b bool)) func(gtx C) D {
	s := &SwitchDef{}
	s.th = th
	s.SetupTabs()
	s.size = th.TextSize
	s.Value = initialState
	s.handler = handler
	s.padding = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(5), Right: unit.Dp(5)}
	return func(gtx C) D {
		//dims := s.Layout(gtx)
		dims := s.padding.Layout(gtx, func(gtx C) D {return s.Layout(gtx)})
		if handler != nil {
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
	trackHeight := gtx.Px(s.size.Scale(1.0))
	thumbSize := gtx.Px(s.size.Scale(1.25))
	trackOff := float32(thumbSize-trackHeight) * .4

	// Draw track.
	trackCorner := float32(trackHeight) / 2
	trackRect := f32.Rectangle{Max: f32.Point{
		X: float32(trackWidth),
		Y: float32(trackHeight),
	}}
	trackColor :=  MulAlpha(s.th.Primary, 0x80)
	dotColor := s.th.Primary
	if !s.Value {
		trackColor = Gray(trackColor)
		dotColor = s.th.Background
	}
	stack := op.Save(gtx.Ops)
	op.Offset(f32.Point{Y: trackOff}).Add(gtx.Ops)
	clip.UniformRRect(trackRect, trackCorner).Add(gtx.Ops)
	paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	stack.Load()

	if gtx.Queue == nil {
		dotColor = Disabled(dotColor)
		trackColor = Disabled(trackColor)
	}


	// Compute thumb offset and color.
	stack = op.Save(gtx.Ops)
	if s.Value {
		off := trackWidth - thumbSize
		op.Offset(f32.Point{X: float32(off)}).Add(gtx.Ops)
	}

	thumbRadius := float32(thumbSize) / 2

	// Draw hover.
	if s.Hovered() || s.Focused() {
		r := 1.4 * thumbRadius
		paint.FillShape(gtx.Ops,  MulAlpha(s.th.Primary, 88),
			clip.Circle{
				Center: f32.Point{X: thumbRadius, Y: thumbRadius},
				Radius: r,
			}.Op(gtx.Ops))
	}

	// Draw thumb shadow, a translucent disc slightly larger than the thumb itself.
	// Center shadow horizontally and slightly adjust its Y.
	paint.FillShape(gtx.Ops, trackColor,
		clip.Circle{
			Center: f32.Point{X: thumbRadius, Y: thumbRadius + 0.05},
			Radius: thumbRadius + 1.5,
		}.Op(gtx.Ops))

	// Draw thumb.
	paint.FillShape(gtx.Ops, dotColor,
		clip.Circle{
			Center: f32.Point{X: thumbRadius, Y: thumbRadius},
			Radius: thumbRadius-1,
		}.Op(gtx.Ops))
	stack.Load()

	gtx.Constraints.Min = image.Pt(trackWidth, trackHeight)
	s.LayoutClickable(gtx)
	s.HandleClicks(gtx)
	s.HandleKeys(gtx)
	dims := image.Point{X: trackWidth, Y: thumbSize}
	return layout.Dimensions{Size: dims}
}
