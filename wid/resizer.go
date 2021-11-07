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
	// Axis defines how the widgets and the handle are laid out.
	Axis layout.Axis
	// Ratio defines how much space is available to the first widget.
	Ratio  float32
	Theme  *Theme
	Length int // max constraint for the axis
	Pos    int // position in pixels of the handle
	drag   gesture.Drag
}

// Split is used to layout two widgets with a splitter between. Axis can be Horizontal or Vertical
func SplitHorizontal(th *Theme, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	r := Resize{Axis: layout.Horizontal}
	r.Theme = th
	r.Ratio = 0.5
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
	c, dims := rs.lo(gtx)
	rs.Ratio = float32(rs.Pos) / float32(rs.Length)
	return layout.Flex{
		Axis: rs.Axis,
	}.Layout(gtx,
		layout.Flexed(rs.Ratio, w1),
		layout.Rigid(func(gtx C) D {
			c.Add(gtx.Ops)
			return D{Size: dims}
		}),
		layout.Flexed(1-rs.Ratio, w2),
	)
}

func (rs *Resize) lo(gtx C) (op.CallOp, image.Point) {
	m := op.Record(gtx.Ops)
	rs.Length = gtx.Constraints.Max.X
	rs.Pos = int(rs.Ratio * float32(rs.Length))
	gtx.Constraints.Min = image.Point{}
	dims := gtx.Constraints.Max
	dims.X = gtx.Px(rs.Theme.SashWidth) + 2*gtx.Px(rs.Theme.SashPadding)
	size := image.Pt(gtx.Px(rs.Theme.SashWidth), dims.Y)
	m1 := op.Offset(f32.Pt(float32(gtx.Px(rs.Theme.SashPadding)), 0)).Push(gtx.Ops)
	m2 := clip.Rect{Max: size}.Push(gtx.Ops)
	paint.ColorOp{Color: rs.Theme.SashColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	m2.Pop()
	m1.Pop()

	var de *pointer.Event
	for _, e := range rs.drag.Events(gtx.Metric, gtx, gesture.Axis(rs.Axis)) {
		if e.Type == pointer.Drag {
			de = &e
		}
	}
	if de != nil {
		xy := de.Position.X
		if rs.Axis == layout.Vertical {
			xy = de.Position.Y
		}
		rs.Pos += int(xy)
	}

	// Clamp the handle position, leaving it always visible.
	if rs.Pos < 0 {
		rs.Pos = 0
	} else if rs.Pos > rs.Length {
		rs.Pos = rs.Length
	}

	rect := image.Rectangle{Max: dims}
	m3 := pointer.Rect(rect).Push(gtx.Ops)
	rs.drag.Add(gtx.Ops)
	if rs.Axis == layout.Horizontal {
		pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
	} else {
		pointer.CursorNameOp{Name: pointer.CursorRowResize}.Add(gtx.Ops)
	}
	m3.Pop()
	return m.Stop(), dims
}
