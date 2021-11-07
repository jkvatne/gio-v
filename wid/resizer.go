package wid

import (
	"image"

	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

// Resize provides a draggable handle in between two widgets for resizing their area.
type Resize struct {
	// Axis defines how the widgets and the handle are laid out.
	Axis layout.Axis
	// Ratio defines how much space is available to the first widget.
	Ratio float32
	float float
	Theme *Theme
}

// Split is used to layout two widgets with a splitter between. Axis can be Horizontal or Vertical
func Split(th *Theme, Axis layout.Axis, w1 layout.Widget, w2 layout.Widget) func(gtx C) D {
	r := Resize{Axis: Axis}
	r.Theme = th
	r.Ratio = 0.5
	return func(gtx C) D {
		return r.Layout(gtx, th, w1, w2)
	}
}

func HorSep(gtx C) D {
	dim := gtx.Constraints.Max
	dim.Y = 12
	size := image.Pt(dim.X, 8)
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: RGB(0x777777)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: dim}
}

func VertSep(gtx C) D {
	dim := gtx.Constraints.Max
	dim.X = 12
	size := image.Pt(12, dim.Y)
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: RGB(0x777777)}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: dim}
}

// Layout displays w1 and w2 with handle in between.
//
// The widgets w1 and w2 must be able to gracefully resize their minimum and maximum dimensions
// in order for the resize to be smooth.
func (rs *Resize) Layout(gtx layout.Context, th *Theme, w1 layout.Widget, w2 layout.Widget) layout.Dimensions {
	var dims D
	// Compute the first widget's max width/height.
	m := op.Record(gtx.Ops)
	if rs.Axis == layout.Horizontal {
		rs.float.Length = gtx.Constraints.Max.X
		rs.float.Pos = int(rs.Ratio * float32(rs.float.Length))
		dims = rs.float.Layout(gtx, rs.Axis, VertSep)
	} else {
		rs.float.Length = gtx.Constraints.Max.Y
		rs.float.Pos = int(rs.Ratio * float32(rs.float.Length))
		dims = rs.float.Layout(gtx, rs.Axis, HorSep)
	}
	c := m.Stop()
	rs.Ratio = float32(rs.float.Pos) / float32(rs.float.Length)
	return layout.Flex{
		Axis: rs.Axis,
	}.Layout(gtx,
		layout.Flexed(rs.Ratio, w1),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			c.Add(gtx.Ops)
			return dims
		}),
		layout.Flexed(1-rs.Ratio, w2),
	)
}

type float struct {
	Length int // max constraint for the axis
	Pos    int // position in pixels of the handle
	drag   gesture.Drag
}

func (f *float) Layout(gtx layout.Context, axis layout.Axis, w layout.Widget) layout.Dimensions {
	gtx.Constraints.Min = image.Point{}
	dims := w(gtx)

	var de *pointer.Event
	for _, e := range f.drag.Events(gtx.Metric, gtx, gesture.Axis(axis)) {
		if e.Type == pointer.Drag {
			de = &e
		}
	}
	if de != nil {
		xy := de.Position.X
		if axis == layout.Vertical {
			xy = de.Position.Y
		}
		f.Pos += int(xy)
	}

	// Clamp the handle position, leaving it always visible.
	if f.Pos < 0 {
		f.Pos = 0
	} else if f.Pos > f.Length {
		f.Pos = f.Length
	}

	rect := image.Rectangle{Max: dims.Size}
	defer pointer.Rect(rect).Push(gtx.Ops).Pop()
	f.drag.Add(gtx.Ops)

	if axis == layout.Horizontal {
		pointer.CursorNameOp{Name: pointer.CursorColResize}.Add(gtx.Ops)
	} else {
		pointer.CursorNameOp{Name: pointer.CursorRowResize}.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: dims.Size}
}
