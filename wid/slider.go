// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"
)

type SliderStyle struct {
	Widget
	Clickable
	Float    *Float
	Min, Max float32
	Color    color.NRGBA
	FingerSize unit.Value
}

// Slider is for selecting a value in a range.
func Slider(th *Theme, min, max float32) layout.Widget {
	s := SliderStyle{
		Min:        min,
		Max:        max,
		Float:      &Float{},
		Color:      th.OnBackground,
		FingerSize: th.FingerSize,
	}
	s.SetupTabs()
	return func(gtx C) D {
		return s.Layout(gtx)
	}
}

func (s *SliderStyle) Layout(gtx layout.Context) layout.Dimensions {
	thumbRadius := gtx.Px(unit.Dp(18))
	trackWidth := gtx.Px(unit.Dp(14))

	axis := s.Float.Axis
	// Keep a minimum length so that the track is always visible.
	minLength := thumbRadius + 3*thumbRadius + thumbRadius
	// Try to expand to finger size, but only if the constraints
	// allow for it.
	touchSizePx := min(gtx.Px(s.FingerSize), axis.Convert(gtx.Constraints.Max).Y)
	sizeMain := max(axis.Convert(gtx.Constraints.Min).X, minLength)
	sizeCross := max(2*thumbRadius, touchSizePx)
	size := axis.Convert(image.Pt(sizeMain, sizeCross))

	o := axis.Convert(image.Pt(thumbRadius, 0))
	op.Offset(layout.FPt(o)).Add(gtx.Ops)
	gtx.Constraints.Min = axis.Convert(image.Pt(sizeMain-2*thumbRadius, sizeCross))

	size = gtx.Constraints.Min
	s.Float.length = float32(s.Float.Axis.Convert(size).X)

	var de *pointer.Event
	for _, e := range s.Float.drag.Events(gtx.Metric, gtx, gesture.Axis(s.Float.Axis)) {
		if e.Type == pointer.Press || e.Type == pointer.Drag {
			de = &e
		}
	}

	value := s.Float.Value
	if de != nil {
		xy := de.Position.X
		if s.Float.Axis == layout.Vertical {
			xy = de.Position.Y
		}
		s.Float.pos = xy / s.Float.length
		value = s.Min + (s.Max-s.Min)*s.Float.pos
	} else if s.Min != s.Max {
		s.Float.pos = (value - s.Min) / (s.Max - s.Min)
	}
	// Unconditionally call setValue in case min, max, or value changed.
	s.setValue(value, s.Min, s.Max)

	if s.Float.pos < 0 {
		s.Float.pos = 0
	} else if s.Float.pos > 1 {
		s.Float.pos = 1
	}

	margin := s.Float.Axis.Convert(image.Pt(thumbRadius, 0))
	rect := image.Rectangle{
		Min: margin.Mul(-1),
		Max: size.Add(margin),
	}
	pointer.Rect(rect).Add(gtx.Ops)
	s.Float.drag.Add(gtx.Ops)


	gtx.Constraints.Min = gtx.Constraints.Min.Add(axis.Convert(image.Pt(0, sizeCross)))
	thumbPos := thumbRadius + int(s.Pos())

	color := s.Color
	if gtx.Queue == nil {
		color = Disabled(color)
	}

	// Draw track before thumb.
	st := op.Save(gtx.Ops)
	track := image.Rectangle{
		Min: axis.Convert(image.Pt(thumbRadius, sizeCross/2-trackWidth/2)),
		Max: axis.Convert(image.Pt(thumbPos, sizeCross/2+trackWidth/2)),
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, color)
	st.Load()

	// Draw track after thumb.
	st = op.Save(gtx.Ops)
	track = image.Rectangle{
		Min: axis.Convert(image.Pt(thumbPos, axis.Convert(track.Min).Y)),
		Max: axis.Convert(image.Pt(sizeMain-thumbRadius, axis.Convert(track.Max).Y)),
	}
	clip.Rect(track).Add(gtx.Ops)
	paint.Fill(gtx.Ops, MulAlpha(color, 96))
	st.Load()

	// Draw thumb.
	pt := axis.Convert(image.Pt(thumbPos, sizeCross/2))
	if s.Hovered() {
		color = Hovered(color)
	}
	paint.FillShape(gtx.Ops, color,
		clip.Circle{
			Center: f32.Point{X: float32(pt.X), Y: float32(pt.Y)},
			Radius: float32(thumbRadius),
		}.Op(gtx.Ops))

	s.LayoutClickable(gtx)
	s.HandleKeys(gtx)
	s.HandleClicks(gtx)

	return layout.Dimensions{Size: size}
}


// Float is for selecting a value in a range.
type Float struct {
	Value float32
	Axis  layout.Axis
	drag    gesture.Drag
	pos     float32 // position normalized to [0, 1]
	length  float32
	changed bool
}

// Dragging returns whether the value is being interacted with.
func (s *SliderStyle) Dragging() bool { return s.Float.drag.Dragging() }


func (s *SliderStyle) setValue(value, min, max float32) {
	if min > max {
		min, max = max, min
	}
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	if s.Float.Value != value {
		s.Float.Value = value
		s.Float.changed = true
	}
}

// Pos reports the selected position.
func (s *SliderStyle) Pos() float32 {
	return s.Float.pos * s.Float.length
}

// Changed reports whether the value has changed since
// the last call to Changed.
func (s *SliderStyle) Changed() bool {
	changed := s.Float.changed
	s.Float.changed = false
	return changed
}
