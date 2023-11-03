package wid

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Resize provides a draggable handle in between two widgets for resizing their area.
type Resize struct {
	// ratio defines how much space is available to the first widget.
	axis   layout.Axis
	th     *Theme
	ratio  float32
	Length float32
	drag   gesture.Drag
	pos    float32
	start  float32
}

// SplitHorizontal is used to layout two widgets with a vertical splitter between.
func SplitHorizontal(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	rs := Resize{th: th, ratio: ratio, axis: layout.Horizontal}
	return func(gtx C) D {
		return rs.Layout(gtx, w1, w2)
	}
}

// SplitVertical is used to layout two widgets with a vertical splitter between.
func SplitVertical(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	rs := Resize{th: th, ratio: ratio, axis: layout.Vertical}
	return func(gtx C) D {
		return rs.Layout(gtx, w1, w2)
	}
}

func (rs *Resize) get(r image.Point) float32 {
	if rs.axis == layout.Horizontal {
		return float32(r.X)
	}
	return float32(r.Y)
}

// Layout displays w1 and w2 with handle in between.
func (rs *Resize) Layout(gtx C, w1 layout.Widget, w2 layout.Widget) D {
	max := rs.get(gtx.Constraints.Max)
	if rs.pos != 0 {
		rs.ratio = rs.pos / max
	}
	// Clamp the handle position, leaving it always visible.
	rs.ratio = Clamp(rs.ratio, 0, 1)
	rs.pos = rs.ratio * max
	f := layout.Flex{
		Axis: rs.axis,
	}.Layout(gtx,
		layout.Flexed(rs.ratio, w1),
		layout.Rigid(func(gtx C) D {
			return D{Size: rs.drawSash(gtx)}
		}),
		layout.Flexed(1-rs.ratio, w2),
	)
	// Handle drag events
	for _, e := range rs.drag.Events(gtx.Metric, gtx, gesture.Axis(rs.axis)) {
		p := e.Position.X
		if rs.axis == layout.Vertical {
			p = e.Position.Y
		}
		if e.Type == pointer.Press {
			rs.start = p - rs.ratio*max
		} else {
			rs.pos = p - rs.start
		}
	}
	// Add drag gesture capture
	d := Px(gtx, rs.th.SashWidth)/2 + 1
	p := int(rs.get(f.Size) * rs.ratio)
	if rs.axis == layout.Vertical {
		defer clip.Rect(image.Rect(0, p-d, f.Size.X, p+d)).Push(gtx.Ops).Pop()
	} else {
		defer clip.Rect(image.Rect(p-d, 0, p+d, f.Size.X)).Push(gtx.Ops).Pop()
	}
	rs.drag.Add(gtx.Ops)
	// Setup cursor for sash
	if rs.axis == layout.Horizontal {
		pointer.CursorColResize.Add(gtx.Ops)
	} else {
		pointer.CursorRowResize.Add(gtx.Ops)
	}
	return f
}

func (rs *Resize) drawSash(gtx C) image.Point {
	dims := gtx.Constraints.Max
	if rs.axis == layout.Horizontal {
		dims.X = Px(gtx, rs.th.SashWidth)
	} else {
		dims.Y = Px(gtx, rs.th.SashWidth)
	}
	defer clip.Rect{Max: dims}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: rs.th.SashColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return dims
}
