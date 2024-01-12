// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/io/key"

	"gioui.org/io/semantic"

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
	Base
	focused  bool
	hovered  bool
	axis     layout.Axis
	drag     gesture.Drag
	pos      float32 // position normalized to [0, 1]
	length   float32
	min, max float32
	Value    *float32
	keyTag   struct{}
}

// Slider is for selecting a value in a range.
func Slider(th *Theme, value *float32, minV, maxV float32, options ...Option) layout.Widget {
	s := SliderStyle{
		min:   minV,
		max:   maxV,
		Value: value,
	}
	s.th = th
	s.width = unit.Dp(99999)
	s.Apply(options...)

	return func(gtx C) D {
		s.handleKeys(gtx)
		m := op.Record(gtx.Ops)
		dims := s.Layout(gtx)
		c := m.Stop()
		defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()
		disabled := gtx.Queue == nil
		keys := key.Set("")
		if !disabled {
			keys = "(Ctrl)-[→,↓,←,↑,0,9,↓]"
			if !s.focused {
				keys = ""
			}
			key.InputOp{Tag: &s.keyTag, Keys: keys}.Add(gtx.Ops)
		} else {
			s.focused = false
		}
		c.Add(gtx.Ops)
		return dims
	}
}

func (s *SliderStyle) handleKeys(gtx C) {
	for _, ev := range gtx.Events(&s.keyTag) {
		switch ke := ev.(type) {
		case key.FocusEvent:
			s.focused = ke.Focus
		case key.Event:
			if ke.State == key.Press {
				d := float32(0.01)
				if ke.Modifiers.Contain(key.ModCtrl) {
					d = 0.1
				}
				switch ke.Name {
				case "0":
					s.pos = 0
				case "9":
					s.pos = 1.0
				case key.NameUpArrow, key.NameLeftArrow:
					s.pos -= d
				case key.NameDownArrow, key.NameRightArrow:
					s.pos += d
				}
				s.setValue()
			}
		}
	}
}

// Layout will draw the slider
func (s *SliderStyle) Layout(gtx C) D {
	w := Px(gtx, s.width)
	if w < gtx.Constraints.Min.X {
		gtx.Constraints.Min.X = w
	}
	thumbRadius := Px(gtx, s.th.TextSize*0.5)
	trackWidth := thumbRadius

	// Keep a minimum length so that the track is always visible.
	minLength := thumbRadius + 3*thumbRadius + thumbRadius
	// Try to expand to finger size, but only if the constraints allow for it.
	touchSizePx := Min(Px(gtx, s.th.FingerSize), s.axis.Convert(gtx.Constraints.Max).Y)
	sizeMain := Max(s.axis.Convert(gtx.Constraints.Min).X, minLength)
	sizeCross := Max(2*thumbRadius, touchSizePx)

	o := s.axis.Convert(image.Pt(thumbRadius, 0))
	op.Offset(o).Add(gtx.Ops)
	gtx.Constraints.Min = s.axis.Convert(image.Pt(sizeMain-2*thumbRadius, sizeCross))

	disabled := gtx.Queue == nil
	semantic.EnabledOp(disabled).Add(gtx.Ops)
	semantic.Switch.Add(gtx.Ops)

	size := gtx.Constraints.Min
	s.length = float32(s.axis.Convert(size).X)

	var de *pointer.Event
	for _, e := range s.drag.Update(gtx.Metric, gtx, gesture.Axis(s.axis)) {
		switch e.Kind {
		case pointer.Press, pointer.Drag:
			key.FocusOp{Tag: &s.keyTag}.Add(gtx.Ops)
			de = &e
		case pointer.Leave, pointer.Cancel:
			s.hovered = false
		case pointer.Enter:
			s.hovered = true
		default:
		}
	}
	if de != nil {
		xy := de.Position.X
		if s.axis == layout.Vertical {
			xy = de.Position.Y
		}
		s.pos = xy / (float32(thumbRadius) + s.length)
	}

	s.setValue()

	margin := s.axis.Convert(image.Pt(thumbRadius, 0))
	rect := image.Rectangle{
		Min: margin.Mul(-1),
		Max: size.Add(margin),
	}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()
	s.drag.Add(gtx.Ops)

	gtx.Constraints.Min = gtx.Constraints.Min.Add(s.axis.Convert(image.Pt(0, sizeCross)))
	thumbPos := thumbRadius + int(s.pos*(float32(sizeMain-thumbRadius*5)))

	color := WithAlpha(s.th.Fg[Canvas], 175)
	if gtx.Queue == nil {
		color = Disabled(color)
	}

	// Draw track before thumb.
	track := image.Rectangle{
		Min: s.axis.Convert(image.Pt(thumbRadius, sizeCross/2-trackWidth/2)),
		Max: s.axis.Convert(image.Pt(thumbPos, sizeCross/2+trackWidth/2)),
	}
	paint.FillShape(gtx.Ops, color, clip.RRect{
		Rect: image.Rect(track.Min.X, track.Min.Y, track.Max.X, track.Max.Y),
		SW:   5, NW: 5, NE: 5, SE: 5,
	}.Op(gtx.Ops))

	// Draw track after thumb.
	track = image.Rectangle{
		Min: s.axis.Convert(image.Pt(thumbPos, s.axis.Convert(track.Min).Y)),
		Max: s.axis.Convert(image.Pt(sizeMain-2*thumbRadius, s.axis.Convert(track.Max).Y)),
	}
	paint.FillShape(gtx.Ops, WithAlpha(color, 80), clip.RRect{
		Rect: image.Rect(track.Min.X, track.Min.Y, track.Max.X, track.Max.Y),
		SW:   5, NW: 5, NE: 5, SE: 5,
	}.Op(gtx.Ops))

	// Draw thumb.
	pt := s.axis.Convert(image.Pt(thumbPos, sizeCross/2))
	if s.hovered || s.focused {
		r := int(float32(thumbRadius) * 1.35)
		ul := image.Pt(pt.X-r, pt.Y-r)
		lr := image.Pt(pt.X+r, pt.Y+r)
		paint.FillShape(gtx.Ops, MulAlpha(s.th.Fg[Canvas], 88), clip.Ellipse{Min: ul, Max: lr}.Op(gtx.Ops))
	}
	r := thumbRadius
	ul := image.Pt(pt.X-r, pt.Y-r)
	lr := image.Pt(pt.X+r, pt.Y+r)
	paint.FillShape(gtx.Ops, s.th.Fg[Canvas], clip.Ellipse{Min: ul, Max: lr}.Op(gtx.Ops))

	return layout.Dimensions{Size: size}
}

func (s *SliderStyle) setValue() {
	if s.pos < 0 {
		s.pos = 0
	}
	if s.pos > 1.0 {
		s.pos = 1.0
	}
	GuiLock.Lock()
	*s.Value = s.pos*(s.max-s.min) + s.min
	GuiLock.Unlock()
}
