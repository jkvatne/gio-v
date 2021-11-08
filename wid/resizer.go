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
	// Ratio defines how much space is available to the first widget.
	Axis   layout.Axis
	Theme  *Theme
	Ratio  float32
	Length float32
	drag   gesture.Drag
	start  int
}

// SplitHorizontal is used to layout two widgets with a vertical splitter between.
func SplitHorizontal(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	rs := Resize{Theme: th, Ratio: ratio}
	rs.Axis = layout.Horizontal
	return func(gtx C) D {
		return rs.Layout(gtx, th, w1, w2)
	}
}

// SplitVertical is used to layout two widgets with a vertical splitter between.
func SplitVertical(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	rs := Resize{Theme: th, Ratio: ratio}
	rs.Axis = layout.Vertical
	return func(gtx C) D {
		return rs.Layout(gtx, th, w1, w2)
	}
}

// Layout displays w1 and w2 with handle in between.
// The widgets w1 and w2 must be able to gracefully resize their minimum and maximum dimensions
// in order for the resize to be smooth.
func (rs *Resize) Layout(gtx C, th *Theme, w1 layout.Widget, w2 layout.Widget) D {
	// Compute the first widget's max width/height.
	if rs.Axis == layout.Horizontal {
		rs.Length = float32(gtx.Constraints.Max.X)
		rs.dragging(gtx)
	} else {
		rs.Length = float32(gtx.Constraints.Max.Y)
		rs.dragging(gtx)
	}
	return layout.Flex{
		Axis: rs.Axis,
	}.Layout(gtx,
		layout.Flexed(rs.Ratio, w1),
		layout.Rigid(func(gtx C) D {
			dims := rs.drawSash(gtx)
			rs.setCursor(gtx, dims)
			return D{Size: dims}
		}),
		layout.Flexed(1-rs.Ratio, w2),
	)
}

func (rs *Resize) setCursor(gtx C, dims image.Point) {
	rect := image.Rectangle{Max: dims}
	defer pointer.Rect(rect).Push(gtx.Ops).Pop()
	rs.drag.Add(gtx.Ops)
	if rs.Axis == layout.Horizontal {
		pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
	} else {
		pointer.CursorNameOp{Name: pointer.CursorRowResize}.Add(gtx.Ops)
	}
}

func (rs *Resize) dragging(gtx C) {
	var dp int
	pos := int(rs.Ratio * rs.Length)
	for _, e := range rs.drag.Events(gtx.Metric, gtx, gesture.Axis(rs.Axis)) {
		if rs.Axis == layout.Horizontal {
			dp = int(e.Position.X)
		} else {
			dp = int(e.Position.Y)
		}
		if e.Type == pointer.Drag {
			pos += dp - rs.start
		}
		if e.Type == pointer.Press {
			rs.start = dp
			return
		}
	}
	// Clamp the handle position, leaving it always visible.
	rs.Ratio = float32(pos) / rs.Length
	if rs.Ratio < 0.05 {
		rs.Ratio = 0.05
	} else if rs.Ratio > 0.95 {
		rs.Ratio = 0.95
	}
}

func (rs *Resize) drawSash(gtx C) image.Point {
	var sashSize, dims image.Point
	if rs.Axis == layout.Horizontal {
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
