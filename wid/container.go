// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/event"
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
// func Container(th *Theme, role UIRole, rr unit.Dp, padding layout.Inset, margin layout.Inset, widgets ...Wid) Wid {
func Container(th *Theme, options ...any) Wid {
	tag := 0
	n := 0
	collapsable := false
	collapsed := false
	role := PrimaryContainer
	margin := th.DefaultPadding
	padding := th.DefaultMargin
	rr := th.DialogCorners
	var description Wid
	var widgets []Wid
	// Read in all options to change from default values to something else.
	for _, option := range options {
		if v, ok := option.(UIRole); ok {
			role = v
		} else if v, ok := option.(layout.Inset); ok {
			if n == 0 {
				padding = v
				n++
			} else {
				margin = v
			}
		} else if v, ok := option.(bool); ok {
			collapsable = v
		} else if v, ok := option.(string); ok {
			description = Label(th, v)
		} else if v, ok := option.(unit.Dp); ok {
			rr = v
		} else if v, ok := option.(int); ok {
			rr = unit.Dp(v)
		} else if v, ok := option.(Wid); ok {
			widgets = append(widgets, v)
		} else if v, ok := option.([]Wid); ok {
			widgets = append(widgets, v...)
		} else {
			panic("Unknown argument to Container()")
		}
	}

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

		if collapsable && collapsed {
			if len(widgets) > 1 {
				if description == nil {
					size := max(gtx.Sp(th.TextSize), gtx.Dp(rr*2))
					outline := image.Rect(0, 0, gtx.Constraints.Max.X-mr-ml, size+pt+pb)
					c := clip.UniformRRect(outline, Px(gtx, rr)).Push(gtx.Ops)
					paint.Fill(gtx.Ops, th.Bg[role])
					c.Pop()
				} else {
					macro := op.Record(gtx.Ops)
					size = description(gtx).Size.Y
					call := macro.Stop()
					// Draw surface
					outline := image.Rect(0, 0, gtx.Constraints.Max.X-mr-ml, size+pt+pb)
					c := clip.UniformRRect(outline, Px(gtx, rr)).Push(gtx.Ops)
					paint.Fill(gtx.Ops, th.Bg[role])
					trans := op.Offset(image.Pt(pl, pt)).Push(gtx.Ops)
					call.Add(gtx.Ops)
					trans.Pop()
					c.Pop()
				}
			}
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
		defer op.Offset(image.Pt(sz.X-50-pr-mr, gtx.Dp(rr/2-5))).Push(gtx.Ops).Pop()
		c.Constraints.Max = image.Pt(50, 50)
		c.Constraints.Min = c.Constraints.Max
		if collapsable {
			if collapsed {
				dropDownIcon.Layout(c, th.Fg[role])
			} else {
				dropUpIcon.Layout(c, th.Fg[role])
			}
		}
		// Setup event handler
		defer clip.UniformRRect(image.Rect(0, 0, 50, 50), 0).Push(gtx.Ops).Pop()
		event.Op(gtx.Ops, &tag)

		for {
			event, ok := gtx.Event(pointer.Filter{
				Target: &tag,
				Kinds:  pointer.Release,
			})
			if !ok {
				break
			}
			if _, ok := event.(pointer.Event); ok {
				if collapsable {
					collapsed = !collapsed
				}
			}
		}
		return D{Size: sz, Baseline: sz.Y}
	}
}
