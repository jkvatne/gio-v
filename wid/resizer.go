package wid

import (
	"image"

	"gioui.org/op"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Resize provides a draggable handle in between two widgets for resizing their area.
type Resize struct {
	// Ratio defines how much space is available to the first widget.
	Ratio  float32
	Theme  *Theme
	Length int // max constraint for the axis
	pos    int // position in pixels of the handle
	drag   gesture.Drag
	start  int
}

// SplitHorizontal is used to layout two widgets with a vertical splitter between.
func SplitHorizontal(th *Theme, ratio float32, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	r := Resize{Theme: th, Ratio: ratio}
	return func(gtx C) D {
		return r.Layout(gtx, th, w1, w2)
	}
}

// Layout displays w1 and w2 with handle in between.
//
// The widgets w1 and w2 must be able to gracefully resize their minimum and maximum dimensions
// in order for the resize to be smooth.
func (rs *Resize) Layout(gtx C, th *Theme, w1 layout.Widget, w2 layout.Widget) D {
	// Compute the first widget's max width/height.
	rs.Length = gtx.Constraints.Max.X
	rs.pos = int(rs.Ratio * float32(rs.Length))
	rs.dragging(gtx, 0, gtx.Constraints.Max.X)
	rs.Ratio = float32(rs.pos) / float32(rs.Length)
	c, dims := rs.lo(gtx)
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(rs.Ratio, w1),
		layout.Rigid(func(gtx C) D {
			c.Add(gtx.Ops)
			return D{Size: dims}
		}),
		layout.Flexed(1-rs.Ratio, w2),
	)
}

func (rs *Resize) dragging(gtx C, lo int, hi int) {
	for _, e := range rs.drag.Events(gtx.Metric, gtx, gesture.Axis(layout.Horizontal)) {
		if e.Type == pointer.Drag {
			rs.pos += int(e.Position.X) - rs.start
		}
		if e.Type == pointer.Press {
			rs.start = int(e.Position.X)
		}
	}
	// Clamp the handle position, leaving it always visible.
	if rs.pos < lo {
		rs.pos = lo
	} else if rs.pos > hi {
		rs.pos = hi
	}
}

func (rs *Resize) drawSash(gtx C) image.Point {
	dims := gtx.Constraints.Max
	dims.X = gtx.Px(rs.Theme.SashWidth) + 2*gtx.Px(rs.Theme.SashPadding)
	size := image.Pt(gtx.Px(rs.Theme.SashWidth), dims.Y)
	m1 := op.Offset(f32.Pt(float32(gtx.Px(rs.Theme.SashPadding)), 0)).Push(gtx.Ops)
	m2 := clip.Rect{Max: size}.Push(gtx.Ops)
	paint.ColorOp{Color: rs.Theme.SashColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	m2.Pop()
	m1.Pop()
	return dims
}

func (rs *Resize) lo(gtx C) (op.CallOp, image.Point) {
	m := op.Record(gtx.Ops)
	gtx.Constraints.Min = image.Point{}
	dims := rs.drawSash(gtx)
	rect := image.Rectangle{Max: dims}
	m3 := pointer.Rect(rect).Push(gtx.Ops)
	rs.drag.Add(gtx.Ops)
	pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
	m3.Pop()
	return m.Stop(), dims
}
