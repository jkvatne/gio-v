// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"math"
)

// Container makes a container with given background
func Container(th *Theme, role UIRole, rr unit.Dp, padding layout.Inset, margin layout.Inset, widgets ...Wid) Wid {
	return func(gtx C) D {
		size := 0
		cgtx := gtx
		cgtx.Constraints.Min.Y = 0
		cgtx.Constraints.Max.X -= Px(gtx, padding.Left+margin.Left+padding.Right+margin.Right)
		if cgtx.Constraints.Min.X > cgtx.Constraints.Max.X {
			cgtx.Constraints.Min.X = cgtx.Constraints.Max.X
		}
		calls := make([]op.CallOp, len(widgets))
		dims := make([]D, len(widgets))
		remaining := gtx.Constraints.Max.Y
		for i, child := range widgets {
			macro := op.Record(gtx.Ops)
			cgtx.Constraints.Max.Y = remaining
			dims[i] = child(cgtx)
			calls[i] = macro.Stop()
			size += dims[i].Size.Y
			remaining = Max(0, remaining-dims[i].Size.Y)
		}
		// Increase size by padding
		size += Px(gtx, padding.Top+padding.Bottom)
		// Offset by margin
		defer op.Offset(image.Pt(Px(gtx, margin.Left), Px(gtx, margin.Top))).Push(gtx.Ops).Pop()
		// Draw surface
		outline := image.Rect(0, 0, gtx.Constraints.Max.X-Px(gtx, margin.Right+margin.Left), size)
		defer clip.UniformRRect(outline, Px(gtx, rr)).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, th.Bg[role])

		// Now do the actual drawing, with offsets
		y := Px(gtx, padding.Top)
		for i := range widgets {
			trans := op.Offset(image.Pt(Px(gtx, padding.Left), int(math.Round(float64(y))))).Push(gtx.Ops)
			calls[i].Add(gtx.Ops)
			trans.Pop()
			y += dims[i].Size.Y
			if y >= gtx.Constraints.Max.Y {
				break
			}
		}
		y += Px(gtx, padding.Bottom)
		y += Px(gtx, margin.Bottom)
		sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, y))
		return D{Size: sz, Baseline: sz.Y}
	}
}
