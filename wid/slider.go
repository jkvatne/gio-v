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
//	Widget
//	Clickable
	float *Float
	Min, Max float32
	Color    color.NRGBA
	FingerSize unit.Value
}

// Slider is for selecting a value in a range.
func Slider(th *Theme, min, max float32) layout.Widget {
	s := SliderStyle{
		Min:        min,
		Max:        max,
		float: &Float{},
		Color:      th.OnBackground,
		FingerSize: th.FingerSize,
	}
	//s.SetupTabs()
	return func(gtx C) D {
		return s.Layout(gtx)
	}
}

func (s SliderStyle) Layout(gtx layout.Context) layout.Dimensions {
	thumbRadius := gtx.Px(unit.Dp(18))
	trackWidth := gtx.Px(unit.Dp(14))

	axis := s.float.Axis
	// Keep a minimum length so that the track is always visible.
	minLength := thumbRadius + 3*thumbRadius + thumbRadius
	// Try to expand to finger size, but only if the constraints
	// allow for it.
	touchSizePx := min(gtx.Px(s.FingerSize), axis.Convert(gtx.Constraints.Max).Y)
	sizeMain := max(axis.Convert(gtx.Constraints.Min).X, minLength)
	sizeCross := max(2*thumbRadius, touchSizePx)
	size := axis.Convert(image.Pt(sizeMain, sizeCross))

	st := op.Save(gtx.Ops)
	o := axis.Convert(image.Pt(thumbRadius, 0))
	op.Offset(layout.FPt(o)).Add(gtx.Ops)
	gtx.Constraints.Min = axis.Convert(image.Pt(sizeMain-2*thumbRadius, sizeCross))
	s.float.Layout(gtx, thumbRadius, s.Min, s.Max)
	gtx.Constraints.Min = gtx.Constraints.Min.Add(axis.Convert(image.Pt(0, sizeCross)))
	thumbPos := thumbRadius + int(s.float.Pos())
	st.Load()

	color := s.Color
	if gtx.Queue == nil {
		color = Disabled(color)
	}

	// Draw track before thumb.
	st = op.Save(gtx.Ops)
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
	paint.FillShape(gtx.Ops, color,
		clip.Circle{
			Center: f32.Point{X: float32(pt.X), Y: float32(pt.Y)},
			Radius: float32(thumbRadius),
		}.Op(gtx.Ops))

	//s.LayoutClickable(gtx)
	//s.HandleClicks(gtx)
	//s.HandleKeys(gtx)

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
func (f *Float) Dragging() bool { return f.drag.Dragging() }

// Layout updates the value according to drag events along the f's main axis.
//
// The range of f is set by the minimum constraints main axis value.
func (f *Float) Layout(gtx layout.Context, pointerMargin int, min, max float32) layout.Dimensions {
	size := gtx.Constraints.Min
	f.length = float32(f.Axis.Convert(size).X)

	margin := f.Axis.Convert(image.Pt(pointerMargin, 0))
	rect := image.Rectangle{
		Min: margin.Mul(-1),
		Max: size.Add(margin),
	}
	pointer.Rect(rect).Add(gtx.Ops)
	f.drag.Add(gtx.Ops)

	var de *pointer.Event
	for _, e := range f.drag.Events(gtx.Metric, gtx, gesture.Axis(f.Axis)) {
		if e.Type == pointer.Press || e.Type == pointer.Drag {
			de = &e
		}
	}

	value := f.Value
	if de != nil {
		xy := de.Position.X
		if f.Axis == layout.Vertical {
			xy = de.Position.Y
		}
		f.pos = xy / f.length
		value = min + (max-min)*f.pos
	} else if min != max {
		f.pos = (value - min) / (max - min)
	}
	// Unconditionally call setValue in case min, max, or value changed.
	f.setValue(value, min, max)

	if f.pos < 0 {
		f.pos = 0
	} else if f.pos > 1 {
		f.pos = 1
	}

	//defer op.Save(gtx.Ops).Load()
	//margin := f.Axis.Convert(image.Pt(pointerMargin, 0))
	//rect := image.Rectangle{
	//	Min: margin.Mul(-1),
	//	Max: size.Add(margin),
	//}
//	pointer.Rect(rect).Add(gtx.Ops)
//	f.drag.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func (f *Float) setValue(value, min, max float32) {
	if min > max {
		min, max = max, min
	}
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	if f.Value != value {
		f.Value = value
		f.changed = true
	}
}

// Pos reports the selected position.
func (f *Float) Pos() float32 {
	return f.pos * f.length
}

// Changed reports whether the value has changed since
// the last call to Changed.
func (f *Float) Changed() bool {
	changed := f.changed
	f.changed = false
	return changed
}
