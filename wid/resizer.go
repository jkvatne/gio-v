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
	Theme  *Theme
	ratio  float32
	Length float32
	drag   gesture.Drag
	pos    float32
	start  float32
}

// SplitHorizontal is used to layout two widgets with a vertical splitter between.
func SplitHorizontal(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	rs := Resize{Theme: th, ratio: ratio, axis: layout.Horizontal}
	return func(gtx C) D {
		return rs.Layout(gtx, w1, w2)
	}
}

// SplitVertical is used to layout two widgets with a vertical splitter between.
func SplitVertical(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	rs := Resize{Theme: th, ratio: ratio, axis: layout.Vertical}
	return func(gtx C) D {
		return rs.Layout(gtx, w1, w2)
	}
}

// Layout displays w1 and w2 with handle in between.
func (rs *Resize) Layout(gtx C, w1 layout.Widget, w2 layout.Widget) D {
	max := float32(gtx.Constraints.Max.Y)
	if rs.axis == layout.Horizontal {
		max = float32(gtx.Constraints.Max.X)
	}
	if rs.pos != 0 {
		rs.ratio = rs.pos / max
	}
	// Clamp the handle position, leaving it always visible.
	rs.ratio = Clamp(rs.ratio, 0.05, 0.95)
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
	d := gtx.Dp(rs.Theme.SashWidth)/2 + 1
	if rs.axis == layout.Vertical {
		p := int(rs.ratio * float32(f.Size.Y))
		defer clip.Rect(image.Rect(0, p-d, f.Size.X, p+d)).Push(gtx.Ops).Pop()
	} else {
		p := int(rs.ratio * float32(f.Size.X))
		defer clip.Rect(image.Rect(p-d, 0, p+d, f.Size.X)).Push(gtx.Ops).Pop()
	}
	rs.drag.Add(gtx.Ops)
	if rs.axis == layout.Horizontal {
		pointer.CursorColResize.Add(gtx.Ops)
	} else {
		pointer.CursorRowResize.Add(gtx.Ops)
	}
	return f
}

func (rs *Resize) drawSash(gtx C) image.Point {
	var sashSize, dims image.Point
	if rs.axis == layout.Horizontal {
		dims = gtx.Constraints.Max
		dims.X = gtx.Dp(rs.Theme.SashWidth)
		sashSize = image.Pt(gtx.Dp(rs.Theme.SashWidth), dims.Y)
	} else {
		dims = gtx.Constraints.Max
		dims.Y = gtx.Dp(rs.Theme.SashWidth)
		sashSize = image.Pt(dims.X, gtx.Dp(rs.Theme.SashWidth))
	}
	defer clip.Rect{Max: sashSize}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: rs.Theme.SashColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return dims
}
