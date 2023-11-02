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
}

// Space will create an open space, like separator but without any line drawn
func Space(size unit.Dp) layout.Widget {
	s := SeparatorStyle{}
	s.thickness = size
	return func(gtx C) D {
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, gtx.Dp(s.thickness))}
	}
}

// Separator creates a material separator widget
func Separator(th *Theme, thickness unit.Dp, options ...Option) layout.Widget {
	s := SeparatorStyle{}
	s.thickness = thickness
	s.role = Canvas
	s.Apply(options...)
	s.th = th
	return func(gtx C) D {
		dim := gtx.Constraints.Min
		dim.Y = th.Px(gtx, s.thickness+s.padding.Top+s.padding.Bottom)
		x := dim.X - th.Px(gtx, s.padding.Left+s.padding.Right)
		y := th.Px(gtx, s.thickness)
		size := image.Pt(x, y)
		if w := th.Px(gtx, s.Base.width); w > size.X {
			size.X = w
		}
		defer op.Offset(image.Pt(th.Px(gtx, s.padding.Left), th.Px(gtx, s.padding.Top))).Push(gtx.Ops).Pop()
		defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: s.Fg()}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: dim}
	}
}
