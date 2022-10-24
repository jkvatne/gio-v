// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

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
	Base
	sw            widget.Bool
	StatePtr      *bool
	trackColorOn  color.NRGBA
	trackColorOff color.NRGBA
	trackOutline  color.NRGBA
	thumbColorOn  color.NRGBA
	thumbColorOff color.NRGBA
	hoverShadow   color.NRGBA
	trackWidth    unit.Dp
	trackStroke   unit.Dp
	trackLength   unit.Dp
	btnOnSize     unit.Dp
	btnOffSize    unit.Dp
}

// Switch returns a widget for a switch
func Switch(th *Theme, statePtr *bool, options ...Option) func(gtx C) D {
	s := &SwitchDef{}
	s.th = th
	s.StatePtr = statePtr
	s.padding = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(5), Right: unit.Dp(5)}
	// Calculate sizes
	s.trackWidth = unit.Dp(s.th.TextSize) * 1.8
	s.trackLength = s.trackWidth * 13 / 8
	s.btnOnSize = s.trackWidth * 3 / 4
	s.btnOffSize = s.trackWidth / 2
	s.trackColorOn = s.th.Bg(Primary)
	s.trackColorOff = s.th.Bg(SurfaceVariant)
	s.trackOutline = s.th.Fg(Outline)
	s.thumbColorOn = s.th.Fg(Primary)
	s.thumbColorOff = s.th.Fg(Outline)
	s.trackStroke = s.th.BorderThickness
	s.hoverShadow = MulAlpha(Gray(s.fgColor), 88)
	return func(gtx C) D {
		semantic.Switch.Add(gtx.Ops)
		dims := s.padding.Layout(gtx, func(gtx C) D { return s.Layout(gtx) })
		if s.onChange != nil {
			s.onChange()
		}
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {

	if s.sw.Changed() {
		*s.StatePtr = s.sw.Value
	} else {
		s.sw.Value = *s.StatePtr
	}

	length := gtx.Dp(s.trackLength)
	width := gtx.Dp(s.trackWidth)
	offSize := gtx.Dp(s.btnOffSize)
	onSize := gtx.Dp(s.btnOnSize)
	stroke := float32(gtx.Dp(s.trackStroke))
	r := gtx.Dp(s.trackWidth / 4)
	trackRect := image.Rect(0, 0, length, width)
	if *s.StatePtr {
		// Draw track in on position, filled rounded
		paint.FillShape(gtx.Ops, s.trackColorOn, clip.UniformRRect(trackRect, width/2).Op(gtx.Ops))
		// Draw thumb,
		paint.FillShape(gtx.Ops, s.thumbColorOn,
			clip.Ellipse{image.Point{X: 3 * r, Y: r / 2}, image.Point{X: onSize + 3*r, Y: onSize + r/2}}.Op(gtx.Ops))

		if s.sw.Hovered() || s.sw.Focused() {
			// Hover is a transparent big circle over the thumb
			paint.FillShape(gtx.Ops, s.hoverShadow,
				clip.Ellipse{Min: image.Point{X: -r, Y: -r}, Max: image.Point{X: offSize + r, Y: offSize + r}}.Op(gtx.Ops))
		}
		// TODO: Draw icon
	} else if !*s.StatePtr {
		// First draw track in OFF position, outlined
		paint.FillShape(gtx.Ops, s.trackColorOff, clip.UniformRRect(trackRect, width/2).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, s.trackOutline,
			clip.Stroke{Path: clip.UniformRRect(trackRect, width/2).Path(gtx.Ops), Width: stroke}.Op())
		// Draw thumb
		paint.FillShape(gtx.Ops, s.thumbColorOff,
			clip.Ellipse{image.Point{X: +r, Y: +r}, image.Point{X: offSize + r, Y: offSize + r}}.Op(gtx.Ops))
		// Draw hover.
		if s.sw.Hovered() || s.sw.Focused() {
			// st := op.Offset(image.Point{X: -r, Y: -r}).Push(gtx.Ops)
			paint.FillShape(gtx.Ops, MulAlpha(s.bgColor, 88),
				clip.Ellipse{image.Point{X: -r, Y: -r}, image.Point{X: +r, Y: offSize + r}}.Op(gtx.Ops))
		}

		// TODO: Draw icon
	}
	// Set up click area.
	clickSize := gtx.Dp(s.th.FingerSize)
	clickOff := image.Point{
		X: (width - clickSize) / 2,
		Y: (width - clickSize) / 2,
	}
	defer op.Offset(clickOff).Push(gtx.Ops).Pop()
	sz := image.Pt(clickSize, clickSize)
	defer clip.Ellipse(image.Rectangle{Max: sz}).Push(gtx.Ops).Pop()
	s.sw.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Switch.Add(gtx.Ops)
		return layout.Dimensions{Size: sz}
	})

	return layout.Dimensions{Size: image.Point{X: length, Y: width}}
}
