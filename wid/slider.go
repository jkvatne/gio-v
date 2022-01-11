// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// SliderStyle is the parameters for a slider
type SliderStyle struct {
	Widget
	Clickable
	axis     layout.Axis
	drag     gesture.Drag
	pos      float32 // position normalized to [0, 1]
	length   float32
	min, max float32
	Value    *float32
}

// Slider is for selecting a value in a range.
func Slider(th *Theme, value *float32, minV, maxV float32, options ...Option) *SliderStyle {
	s := SliderStyle{
		min: minV,
		max: maxV,
	}
	s.Value = value
	s.th = th
	s.SetupTabs()
	s.width = unit.Dp(99999)
	s.Apply(options...)
	return &s
}

// Layout will draw the slider
func (s *SliderStyle) Layout(gtx C) D {
	gtx.Constraints.Min = CalcMin(gtx, s.width)
	thumbRadius := gtx.Px(s.th.TextSize.Scale(0.5))
	trackWidth := gtx.Px(s.th.TextSize.Scale(0.5))

	// Keep a minimum length so that the track is always visible.
	minLength := thumbRadius + 3*thumbRadius + thumbRadius
	// Try to expand to finger size, but only if the constraints
	// allow for it.
	touchSizePx := min(gtx.Px(s.th.FingerSize), s.axis.Convert(gtx.Constraints.Max).Y)
	sizeMain := max(s.axis.Convert(gtx.Constraints.Min).X, minLength)
	sizeCross := max(2*thumbRadius, touchSizePx)
	size := s.axis.Convert(image.Pt(sizeMain, sizeCross))

	o := s.axis.Convert(image.Pt(thumbRadius, 0))
	op.Offset(layout.FPt(o)).Add(gtx.Ops)
	gtx.Constraints.Min = s.axis.Convert(image.Pt(sizeMain-2*thumbRadius, sizeCross))

	size = gtx.Constraints.Min
	s.length = float32(s.axis.Convert(size).X)

	var de *pointer.Event
	for _, e := range s.drag.Events(gtx.Metric, gtx, gesture.Axis(s.axis)) {
		if e.Type == pointer.Press || e.Type == pointer.Drag {
			de = &e
		}
	}
	if s.HandleKeys(gtx) {
		if s.index != nil {
			s.pos = float32(*s.index) / 100.0
		}
		*s.Value = s.min + (s.max-s.min)*s.pos
	}
	if de != nil {
		xy := de.Position.X
		if s.axis == layout.Vertical {
			xy = de.Position.Y
		}
		s.pos = (xy - float32(thumbRadius)) / s.length
		*s.Value = s.min + (s.max-s.min)*s.pos
	} else if s.min != s.max {
		s.pos = (*s.Value - s.min) / (s.max - s.min)
	}
	if s.index != nil {
		*s.index = int(s.pos*100 + 0.5)
	}
	// Unconditionally call setValue in case min, max, or value changed.
	s.setValue(*s.Value, s.min, s.max)
	s.pos = clamp(s.pos, 0, 1)

	margin := s.axis.Convert(image.Pt(thumbRadius, 0))
	rect := image.Rectangle{
		Min: margin.Mul(-1),
		Max: size.Add(margin),
	}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()
	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops).Pop()
	s.drag.Add(gtx.Ops)

	gtx.Constraints.Min = gtx.Constraints.Min.Add(s.axis.Convert(image.Pt(0, sizeCross)))
	thumbPos := thumbRadius + int(s.pos*s.length)

	color := WithAlpha(s.th.OnBackground, 175)
	if gtx.Queue == nil {
		color = Disabled(color)
	}

	// Draw track before thumb.
	track := image.Rectangle{
		Min: s.axis.Convert(image.Pt(thumbRadius, sizeCross/2-trackWidth/2)),
		Max: s.axis.Convert(image.Pt(thumbPos, sizeCross/2+trackWidth/2)),
	}
	paint.FillShape(gtx.Ops, color, clip.RRect{
		Rect: f32.Rect(float32(track.Min.X), float32(track.Min.Y), float32(track.Max.X), float32(track.Max.Y)),
		SW:   5, NW: 5, NE: 5, SE: 5,
	}.Op(gtx.Ops))

	// Draw track after thumb.
	track = image.Rectangle{
		Min: s.axis.Convert(image.Pt(thumbPos, s.axis.Convert(track.Min).Y)),
		Max: s.axis.Convert(image.Pt(sizeMain-thumbRadius, s.axis.Convert(track.Max).Y)),
	}
	paint.FillShape(gtx.Ops, WithAlpha(color, 80), clip.RRect{
		Rect: f32.Rect(float32(track.Min.X), float32(track.Min.Y), float32(track.Max.X), float32(track.Max.Y)),
		SW:   5, NW: 5, NE: 5, SE: 5,
	}.Op(gtx.Ops))

	// Draw thumb.
	pt := s.axis.Convert(image.Pt(thumbPos, sizeCross/2))
	if s.Hovered() || s.Focused() {
		r := float32(thumbRadius) * 1.35
		ul := f32.Pt(float32(pt.X)-r, float32(pt.Y)-r)
		lr := f32.Pt(float32(pt.X)+r, float32(pt.Y)+r)
		paint.FillShape(gtx.Ops, MulAlpha(s.th.OnBackground, 88), clip.Ellipse{Min: ul, Max: lr}.Op(gtx.Ops))
	} else {
		color = s.th.OnBackground
	}
	r := thumbRadius
	ul := f32.Pt(float32(pt.X-r), float32(pt.Y-r))
	lr := f32.Pt(float32(pt.X+r), float32(pt.Y+r))
	paint.FillShape(gtx.Ops, s.th.OnBackground, clip.Ellipse{ul, lr}.Op(gtx.Ops))

	s.LayoutClickable(gtx)

	s.HandleClicks(gtx)

	return layout.Dimensions{Size: size}
}

func (s *SliderStyle) setValue(value, min, max float32) {
	if min > max {
		min, max = max, min
	}
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	if *s.Value != value {
		*s.Value = value
	}
	if s.Value != nil {
		*s.Value = value
	}
}
