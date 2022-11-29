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
	return layout.Flex{
		Axis: rs.axis,
	}.Layout(gtx,
		layout.Flexed(rs.ratio, w1),
		layout.Rigid(func(gtx C) D {
			max := float32(gtx.Constraints.Max.Y)
			if rs.axis == layout.Horizontal {
				max = float32(gtx.Constraints.Max.X)
			}
			var dp int
			pos := int(rs.ratio * max)
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
				}
			}
			// Clamp the handle position, leaving it always visible.
			rs.ratio = Clamp(float32(pos)/max, 0.05, 0.95)
			// Draw the sash
			dims := rs.drawSash(gtx)
			// Setup drag to catch events within clip rect
			defer clip.Rect(image.Rectangle{Max: dims}).Push(gtx.Ops).Pop()
			rs.drag.Add(gtx.Ops)
			// Setup cursor when within sash
			if rs.axis == layout.Horizontal {
				pointer.CursorColResize.Add(gtx.Ops)
			} else {
				pointer.CursorRowResize.Add(gtx.Ops)
			}
			return D{Size: dims}
		}),
		layout.Flexed(1-rs.ratio, w2),
	)
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
