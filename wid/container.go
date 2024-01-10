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

func ConfirmDialog(th *Theme, heading string, text string, lbl1 string, on1 func()) layout.Widget {
	return Container(th, TransparentSurface, 0, FlexInset, NoInset,
		Dialog(th, PrimaryContainer,
			Label(th, heading, Heading(), Middle()),
			Label(th, text, Middle()),
			Separator(th, 0, Pads(10)),
			Row(th, nil, SpaceRightAdjust,
				TextButton(th, lbl1, Do(on1)),
			),
		),
	)
}

func YesNoDialog(th *Theme, heading string, text string, lbl1, lbl2 string, on1, on2 func()) layout.Widget {
	return Dialog(th, PrimaryContainer,
		Label(th, heading, Heading(), Middle()),
		Label(th, text, Middle()),
		Separator(th, 0, Pads(10)),
		Row(th, nil, SpaceRightAdjust,
			TextButton(th, lbl1, Do(on1)),
			TextButton(th, lbl2, Do(on2)),
		),
	)
}
func Dialog(th *Theme, role UIRole, widgets ...Wid) Wid {
	return func(gtx C) D {
		pt := Px(gtx, th.DialogPadding.Top)
		pb := Px(gtx, th.DialogPadding.Bottom)
		pl := Px(gtx, th.DialogPadding.Left)
		pr := Px(gtx, th.DialogPadding.Right)
		// Shade underlying form
		// Draw surface all over the underlying form with the transparent surface color
		outline := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		defer clip.Rect(outline).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, th.Bg[TransparentSurface])
		// Calculate dialog constraints
		ctx := gtx
		ctx.Constraints.Min.Y = 0
		// Margins left and right for a constant maximum dialog size of 35 characters
		ml := Max(12, (ctx.Constraints.Max.X-pl-pr-gtx.Metric.Sp(th.TextSize)*20)/2)
		mr := Max(12, (ctx.Constraints.Max.X-pl-pr-gtx.Metric.Sp(th.TextSize)*20)/2)
		ctx.Constraints.Max.X = gtx.Constraints.Max.X - ml - mr - pl - pr
		calls := make([]op.CallOp, len(widgets))
		dims := make([]D, len(widgets))
		size := 0
		for i, child := range widgets {
			macro := op.Record(gtx.Ops)
			dims[i] = child(ctx)
			calls[i] = macro.Stop()
			size += dims[i].Size.Y
		}
		mt := (gtx.Constraints.Max.Y - size) / 2

		// Draw the dialog surface with caclculated margins
		defer op.Offset(image.Pt(ml, mt)).Push(gtx.Ops).Pop()
		outline = image.Rect(0, 0, gtx.Constraints.Max.X-ml-mr, size+pt+pb)
		defer clip.UniformRRect(outline, Px(gtx, unit.Dp(20))).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, th.Bg[role])
		// Now do the actual drawing of the widgets, with offsets
		y := 33
		for i := range widgets {
			trans := op.Offset(image.Pt(33, int(math.Round(float64(y))))).Push(gtx.Ops)
			calls[i].Add(gtx.Ops)
			trans.Pop()
			y += dims[i].Size.Y
		}
		sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, size+66))
		return D{Size: sz, Baseline: sz.Y}
	}
}

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
		pt := Px(gtx, padding.Top)
		pb := Px(gtx, padding.Bottom)
		pl := Px(gtx, padding.Left)
		ml := Px(gtx, margin.Left)
		mr := Px(gtx, margin.Right)
		mt := Px(gtx, margin.Top)
		mb := Px(gtx, margin.Bottom)
		free := gtx.Constraints.Max.Y - size
		if free > 0 {
			if pt < 0 {
				pt = free / 2
			}
			if pb < 0 {
				pb = free / 2
			}
		}
		// Offset by margin
		defer op.Offset(image.Pt(ml, mt)).Push(gtx.Ops).Pop()
		// Draw surface
		outline := image.Rect(0, 0, gtx.Constraints.Max.X-mr-ml, size+pt+pb)
		defer clip.UniformRRect(outline, Px(gtx, rr)).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, th.Bg[role])

		// Now do the actual drawing, with offsets
		y := pt
		for i := range widgets {
			trans := op.Offset(image.Pt(pl, int(math.Round(float64(y))))).Push(gtx.Ops)
			calls[i].Add(gtx.Ops)
			trans.Pop()
			y += dims[i].Size.Y
			if y >= gtx.Constraints.Max.Y {
				break
			}
		}
		sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, size+mt+mb+pt+pb))
		return D{Size: sz, Baseline: sz.Y}
	}
}
