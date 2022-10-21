package wid

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// SeparatorStyle defines material rendering parameters for separator
type SeparatorStyle struct {
	Base
	thickness unit.Dp
	length    int
}

// Separator creates a material separator widget
func Separator(th *Theme, thickness unit.Dp, options ...Option) layout.Widget {
	s := SeparatorStyle{}
	s.thickness = thickness
	s.fgColor = th.OnBackground
	s.Apply(options...)

	return func(gtx C) D {
		dim := gtx.Constraints.Max
		dim.Y = gtx.Dp(s.thickness) + gtx.Dp(s.padding.Top) + gtx.Dp(s.padding.Bottom)
		size := image.Pt(dim.X-gtx.Dp(s.padding.Left)-gtx.Dp(s.padding.Right), gtx.Dp(s.thickness))
		if w := gtx.Dp(s.Base.width); w > size.X {
			size.X = w
		}
		defer op.Offset(image.Pt(gtx.Dp(s.padding.Left), gtx.Dp(s.padding.Top))).Push(gtx.Ops).Pop()
		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: s.fgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: dim}
	}
}
