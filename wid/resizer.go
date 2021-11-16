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
	start  int
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
	// Compute the first widget's max width/height.
	if rs.axis == layout.Horizontal {
		rs.dragging(gtx, float32(gtx.Constraints.Max.X))
	} else {
		rs.dragging(gtx, float32(gtx.Constraints.Max.Y))
	}
	return layout.Flex{
		Axis: rs.axis,
	}.Layout(gtx,
		layout.Flexed(rs.ratio, w1),
		layout.Rigid(func(gtx C) D {
			dims := rs.drawSash(gtx)
			rs.setCursor(gtx, dims)
			return D{Size: dims}
		}),
		layout.Flexed(1-rs.ratio, w2),
	)
}

func (rs *Resize) setCursor(gtx C, dims image.Point) {
	rect := image.Rectangle{Max: dims}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()
	rs.drag.Add(gtx.Ops)
	if rs.axis == layout.Horizontal {
		pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
	} else {
		pointer.CursorNameOp{Name: pointer.CursorRowResize}.Add(gtx.Ops)
	}
}

func clamp(v float32, lo float32, hi float32) float32 {
	if v < lo {
		return lo
	} else if v > hi {
		return 0.95
	}
	return v
}

func (rs *Resize) dragging(gtx C, length float32) {
	var dp int
	pos := int(rs.ratio * length)
	for _, e := range rs.drag.Events(gtx.Metric, gtx, gesture.Axis(rs.axis)) {
		if rs.axis == layout.Horizontal {
			dp = int(e.Position.X)
		} else {
			dp = int(e.Position.Y)
		}
		if e.Type == pointer.Drag {
			pos += dp - rs.start
		} else if e.Type == pointer.Press {
			rs.start = dp
			return
		}
	}
	// Clamp the handle position, leaving it always visible.
	rs.ratio = clamp(float32(pos)/length, 0.05, 0.95)
}

func (rs *Resize) drawSash(gtx C) image.Point {
	var sashSize, dims image.Point
	if rs.axis == layout.Horizontal {
		dims = gtx.Constraints.Max
		dims.X = gtx.Px(rs.Theme.SashWidth)
		sashSize = image.Pt(gtx.Px(rs.Theme.SashWidth), dims.Y)
	} else {
		dims = gtx.Constraints.Max
		dims.Y = gtx.Px(rs.Theme.SashWidth)
		sashSize = image.Pt(dims.X, gtx.Px(rs.Theme.SashWidth))
	}
	defer clip.Rect{Max: sashSize}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: rs.Theme.SashColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return dims
}
