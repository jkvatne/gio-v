package wid

import (
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
	thickness    unit.Value
	paddingAbove unit.Value
	paddingBelow unit.Value
}

// Separator creates a material separator widget
func Separator(th *Theme, thickness unit.Value, paddingAbove unit.Value, paddingBelow unit.Value) layout.Widget {
	s := SeparatorStyle{
		thickness:    thickness,
		paddingAbove: paddingAbove,
		paddingBelow: paddingBelow,
	}
	return func(gtx C) D {
		dim := gtx.Constraints.Max
		dim.Y = gtx.Px(s.thickness) + gtx.Px(s.paddingAbove) + gtx.Px(s.paddingBelow)
		op.Offset(f32.Pt(0, float32(gtx.Px(paddingAbove)))).Add(gtx.Ops)
		size := dim
		size.Y = gtx.Px(s.thickness)
		clip.Rect{Max: size}.Add(gtx.Ops)
		paint.ColorOp{Color: th.OnBackground}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: dim}
	}
}
