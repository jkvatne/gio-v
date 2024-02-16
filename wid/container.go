// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"math"
)

// Container makes a container widget with background color according to role/theme
func Container(th *Theme, role UIRole, rr unit.Dp, padding layout.Inset, margin layout.Inset, widgets ...Wid) Wid {
	tag := 0
	collapsed := false
	return func(gtx C) D {
		size := 0
		// Scale margins and paddings
		mt, mb, ml, mr := ScaleInset(gtx, margin)
		pt, pb, pl, pr := ScaleInset(gtx, padding)
		c := gtx
		c.Constraints.Min.Y = 0
		c.Constraints.Max.X -= pl + pr + pl + mr
		c.Constraints.Min.X = Min(c.Constraints.Min.X, c.Constraints.Max.X)
		calls := make([]op.CallOp, len(widgets))
		dims := make([]D, len(widgets))
		remaining := gtx.Constraints.Max.Y
		// Create macros and dimensions for all widgets in the container
		for i, child := range widgets {
			macro := op.Record(gtx.Ops)
			c.Constraints.Max.Y = remaining
			dims[i] = child(c)
			calls[i] = macro.Stop()
			size += dims[i].Size.Y
			remaining = Max(0, remaining-dims[i].Size.Y)
		}
		// Offset by margin
		defer op.Offset(image.Pt(ml, mt)).Push(gtx.Ops).Pop()

		if collapsed {
			d := widgets[1](c)
			size = d.Size.Y
			// Draw surface
			outline := image.Rect(0, 0, gtx.Constraints.Max.X-mr-ml, size+pt+pb)
			defer clip.UniformRRect(outline, Px(gtx, rr)).Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, th.Bg[role])
			widgets[1](c)
		} else {
			// Negative padding is used for flexible insets
			// The container will center the contents
			if free := gtx.Constraints.Max.Y - size; free > 0 && pt < 0 {
				pt = free / 2
				pb = free / 2
			}
			// Draw surface
			outline := image.Rect(0, 0, gtx.Constraints.Max.X-mr-ml, size+pt+pb)
			defer clip.UniformRRect(outline, Px(gtx, rr)).Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, th.Bg[role])
			// Now do the actual drawing of widgets in the container, with offsets
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
		}
		sz := gtx.Constraints.Constrain(image.Pt(gtx.Constraints.Max.X, size+mt+mb+pt+pb))
		defer op.Offset(image.Pt(sz.X-50-pr-mr, 0)).Push(gtx.Ops).Pop()
		c.Constraints.Max = image.Pt(50, 50)
		c.Constraints.Min = c.Constraints.Max
		if collapsed {
			dropDownIcon.Layout(c, th.Fg[role])
		} else {
			dropUpIcon.Layout(c, th.Fg[role])
		}
		// Handle events
		for _, ev := range gtx.Events(&tag) {
			if ev, ok := ev.(pointer.Event); ok {
				switch ev.Kind {
				case pointer.Release:
					collapsed = !collapsed
				}
			}
		}
		// Setup event handler
		defer clip.UniformRRect(image.Rect(0, 0, 50, 50), 0).Push(gtx.Ops).Pop()
		pointer.InputOp{
			Tag:   &tag,
			Kinds: pointer.Release,
		}.Add(gtx.Ops)

		return D{Size: sz, Baseline: sz.Y}
	}
}
