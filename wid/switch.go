// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/op"

	"gioui.org/io/pointer"

	"gioui.org/io/semantic"

	"gioui.org/widget"

	"gioui.org/layout"
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
	// Calculate sizes
	s.trackWidth = unit.Dp(s.th.TextSize) * 1.5
	s.trackLength = s.trackWidth * 13 / 8
	s.btnOnSize = s.trackWidth * 3 / 4
	s.btnOffSize = s.trackWidth / 2
	s.trackColorOn = s.th.Bg(Primary)
	s.trackColorOff = s.th.Bg(SurfaceVariant)
	s.trackOutline = s.th.Fg(Outline)
	s.thumbColorOn = s.th.Fg(Primary)
	s.thumbColorOff = s.th.Fg(Outline)
	s.trackStroke = s.th.BorderThickness
	// Default padding. Can be changed with option Padds()
	s.padding = layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(5), Right: unit.Dp(5)}
	// The pointer to a variable receiving switch on/off state
	s.StatePtr = statePtr
	for _, option := range options {
		option.apply(s)
	}
	return func(gtx C) D {
		return s.padding.Layout(gtx, func(gtx C) D {
			return s.Layout(gtx)
		})
	}
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {

	if s.sw.Changed() {
		GuiLock.Lock()
		*s.StatePtr = s.sw.Value
		if s.onUserChange != nil {
			s.onUserChange()
		}
		GuiLock.Unlock()
	} else {
		GuiLock.RLock()
		s.sw.Value = *s.StatePtr
		GuiLock.RUnlock()
	}

	width := gtx.Dp(s.trackLength)
	height := gtx.Dp(s.trackWidth)
	offSize := gtx.Dp(s.btnOffSize)
	onSize := gtx.Dp(s.btnOnSize)
	stroke := float32(gtx.Dp(s.trackStroke))
	r := gtx.Dp(s.trackWidth / 4)
	trackRect := image.Rect(0, 0, width, height)
	if s.sw.Focused() && s.sw.Hovered() {
		s.hoverShadow = MulAlpha(s.th.Bg(Primary), 120)
	} else if s.sw.Focused() {
		s.hoverShadow = MulAlpha(s.th.Bg(Primary), 90)
	} else if s.sw.Hovered() {
		s.hoverShadow = MulAlpha(s.th.Bg(Primary), 60)
	} else {
		s.hoverShadow = MulAlpha(s.th.Bg(Primary), 0)
	}

	if s.sw.Value {
		// Draw track in on position, filled rounded
		paint.FillShape(gtx.Ops, s.trackColorOn, clip.UniformRRect(trackRect, height/2).Op(gtx.Ops))
		// Draw thumb,
		paint.FillShape(gtx.Ops, s.thumbColorOn,
			clip.Ellipse{image.Point{X: 3 * r, Y: r / 2}, image.Point{X: onSize + 3*r, Y: onSize + r/2}}.Op(gtx.Ops))
		// Draw hover/focus shade.
		paint.FillShape(gtx.Ops, s.hoverShadow,
			clip.Ellipse{Min: image.Point{X: 2 * r, Y: -r / 2}, Max: image.Point{X: 7 * r, Y: 9 * r / 2}}.Op(gtx.Ops))
		// TODO: Draw icon
	} else {
		// First draw track in OFF position, outlined
		paint.FillShape(gtx.Ops, s.trackColorOff, clip.UniformRRect(trackRect, height/2).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, s.trackOutline,
			clip.Stroke{Path: clip.UniformRRect(trackRect, height/2).Path(gtx.Ops), Width: stroke}.Op())
		// Draw thumb
		paint.FillShape(gtx.Ops, s.thumbColorOff,
			clip.Ellipse{image.Point{X: +r, Y: +r}, image.Point{X: offSize + r, Y: offSize + r}}.Op(gtx.Ops))
		// Draw hover/focus shade.
		paint.FillShape(gtx.Ops, s.hoverShadow,
			clip.Ellipse{Min: image.Point{X: -r / 2, Y: -r / 2}, Max: image.Point{X: 9 * r / 2, Y: 9 * r / 2}}.Op(gtx.Ops))
		// TODO: Draw icon
	}
	// Set up click area.
	defer op.Offset(image.Point{-10, -10}).Push(gtx.Ops).Pop()
	sz := image.Pt(width+20, height+20)
	clickRect := image.Rect(0, 0, width+20, height+20)
	defer clip.UniformRRect(clickRect, height/2).Push(gtx.Ops).Pop()
	s.sw.Layout(gtx, func(gtx C) D {
		if s.description != "" {
			semantic.DescriptionOp(s.description).Add(gtx.Ops)
		}
		semantic.Switch.Add(gtx.Ops)
		return layout.Dimensions{Size: sz}
	})
	pointer.CursorPointer.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Point{X: width, Y: height}}
}
