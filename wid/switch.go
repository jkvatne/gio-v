// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/op"

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

// DefaultSwitchDef returs a Switchdef with values from the theme
func DefaultSwitchDef(th *Theme) *SwitchDef {
	return &SwitchDef{
		Base: Base{
			th:     th,
			margin: th.DefaultMargin,
		},
		trackWidth:    unit.Dp(th.TextSize) * 1.5,
		trackLength:   unit.Dp(th.TextSize) * 2.4,
		btnOnSize:     unit.Dp(th.TextSize) * 1.15,
		btnOffSize:    unit.Dp(th.TextSize) * 0.8,
		trackColorOn:  th.Bg[Primary],
		trackColorOff: th.Bg[SurfaceVariant],
		trackOutline:  th.Fg[Outline],
		thumbColorOn:  th.Fg[Primary],
		thumbColorOff: th.Fg[Outline],
		trackStroke:   th.BorderThickness,
	}
}

// Switch returns a widget for a switch
func Switch(th *Theme, statePtr *bool, options ...Option) layout.Widget {
	// The pointer to a variable receiving switch on/off state
	s := DefaultSwitchDef(th)
	s.StatePtr = statePtr
	for _, option := range options {
		option.apply(s)
	}
	return s.Layout
}

// Layout updates the switch and displays it.
func (s *SwitchDef) Layout(gtx C) D {
	mt, mb, ml, mr := ScaleInset(gtx, s.margin)
	if *s.StatePtr != s.sw.Value {
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

	// Offset by margin
	defer op.Offset(image.Pt(ml, mt)).Push(gtx.Ops).Pop()

	width := Px(gtx, s.trackLength)
	height := Px(gtx, s.trackWidth)
	offSize := Px(gtx, s.btnOffSize)
	onSize := Px(gtx, s.btnOnSize)
	stroke := float32(Px(gtx, s.trackStroke))
	r := Px(gtx, s.trackWidth/4)
	trackRect := image.Rect(0, 0, width, height)
	if gtx.Focused(&s.sw) && s.sw.Hovered() {
		s.hoverShadow = MulAlpha(s.th.Bg[Primary], 120)
	} else if gtx.Focused(&s.sw) {
		s.hoverShadow = MulAlpha(s.th.Bg[Primary], 95)
	} else if s.sw.Hovered() {
		s.hoverShadow = MulAlpha(s.th.Bg[Primary], 75)
	} else {
		s.hoverShadow = MulAlpha(s.th.Bg[Primary], 0)
	}

	if s.sw.Value {
		// Draw track in on position, filled rounded
		paint.FillShape(gtx.Ops, s.trackColorOn, clip.UniformRRect(trackRect, height/2).Op(gtx.Ops))
		// Draw thumb,
		paint.FillShape(gtx.Ops, s.thumbColorOn,
			clip.Ellipse{Min: image.Point{X: 3 * r, Y: r / 2}, Max: image.Point{X: onSize + 3*r, Y: onSize + r/2}}.Op(gtx.Ops))
		// Draw hover/focus shade.
		paint.FillShape(
			gtx.Ops,
			s.hoverShadow,
			clip.Ellipse{Min: image.Point{X: 2 * r, Y: -r / 2}, Max: image.Point{X: 7 * r, Y: 9 * r / 2}}.Op(gtx.Ops),
		)
		// TODO: Draw icon
	} else {
		// First draw track in OFF position, outlined
		paint.FillShape(gtx.Ops, s.trackColorOff, clip.UniformRRect(trackRect, height/2).Op(gtx.Ops))
		paint.FillShape(gtx.Ops, s.trackOutline,
			clip.Stroke{Path: clip.UniformRRect(trackRect, height/2).Path(gtx.Ops), Width: stroke}.Op())
		// Draw thumb
		paint.FillShape(gtx.Ops, s.thumbColorOff,
			clip.Ellipse{Min: image.Point{X: +r, Y: +r}, Max: image.Point{X: offSize + r, Y: offSize + r}}.Op(gtx.Ops))
		// Draw hover/focus shade.
		paint.FillShape(gtx.Ops, s.hoverShadow,
			clip.Ellipse{Min: image.Point{X: -r / 2, Y: -r / 2}, Max: image.Point{X: 9 * r / 2, Y: 9 * r / 2}}.Op(gtx.Ops))
		// TODO: Draw icon
	}
	// Set up click area.
	defer op.Offset(image.Point{X: -10, Y: -10}).Push(gtx.Ops).Pop()
	sz := image.Pt(width+20, height+20)
	clickRect := image.Rect(0, 0, width+20, height+20)
	defer clip.UniformRRect(clickRect, height/2).Push(gtx.Ops).Pop()
	s.sw.Layout(gtx, func(gtx C) D {
		if s.hint != "" {
			// semantic.DescriptionOp(s.hint).Add(gtx.Ops)
		}
		// semantic.Switch.Add(gtx.Ops)
		return layout.Dimensions{Size: sz}
	})

	return layout.Dimensions{Size: image.Point{X: width + ml + mr, Y: height + mt + mb}}
}
