package wid

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// SeparatorStyle defines material rendering parameters for separator
type SeparatorStyle struct {
	Widget
	thickness unit.Value
	length    int
}

// Separator creates a material separator widget
func Separator(th *Theme, thickness unit.Value, options ...Option) layout.Widget {
	s := SeparatorStyle{}
	s.thickness = thickness
	s.fgColor = th.OnBackground
	s.Apply(options...)

	return func(gtx C) D {
		dim := gtx.Constraints.Max
		dim.Y = gtx.Px(s.thickness) + gtx.Px(s.padding.Top) + gtx.Px(s.padding.Bottom)
		size := image.Pt(dim.X-gtx.Px(s.padding.Left)-gtx.Px(s.padding.Right), gtx.Px(s.thickness))
		if w := gtx.Px(s.Widget.width); w > size.X {
			size.X = w
		}
		defer op.Offset(f32.Pt(float32(gtx.Px(s.padding.Left)), float32(gtx.Px(s.padding.Top)))).Push(gtx.Ops).Pop()
		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: s.fgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: dim}
	}
}
