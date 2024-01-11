package wid

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"math"
	"time"
)

var (
	dialog          layout.Widget
	dialogStartTime time.Time
	startX          int
	startY          int
)

func Show(d layout.Widget) {
	dialog = d
	dialogStartTime = time.Now()
	startX = int(mouseX)
	startY = int(mouseY)
}

func Hide() {
	dialog = nil
	dialogStartTime = time.Time{}
}

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
		f := Min(1.0, float64(time.Since(dialogStartTime))/float64(time.Second/4))

		// Shade underlying form
		// Draw surface all over the underlying form with the transparent surface color
		outline := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
		defer clip.Rect(outline).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, WithAlpha(Black, uint8(f*200)))

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

		// Make animated expansion
		if f < 1.0 {
			// Draw the dialog surface with caclculated margins
			x := int(f*float64(ml) + (1-f)*float64(startX))
			y := int(f*float64(mt) + (1-f)*float64(startY))
			defer op.Offset(image.Pt(x, y)).Push(gtx.Ops).Pop()
			dx := gtx.Constraints.Max.X - ml - mr
			dy := size + pt + pb
			dx = int(f * float64(dx))
			dy = int(f * float64(dy))
			outline = image.Rect(0, 0, dx, dy)
			defer clip.UniformRRect(outline, Px(gtx, unit.Dp(20))).Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, th.Bg[role])
			sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, size+pb+pt))
			Invalidate()
			return D{Size: sz, Baseline: sz.Y}
		} else {
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
			sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, size+pb+pt))
			return D{Size: sz, Baseline: sz.Y}
		}
	}
}
