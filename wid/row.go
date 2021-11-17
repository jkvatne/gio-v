// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/f32"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type rowDef struct {
	Clickable
}

// Row returns a widget grid row with selectable color.
func Row(th *Theme, selected *bool, weights []float32, widgets ...layout.Widget) layout.Widget {
	r := rowDef{}
	return func(gtx C) D {
		bgColor := th.Background
		if r.Hovered() {
			bgColor = Interpolate(th.Background, th.Primary, 0.05)
		} else if selected != nil && *selected {
			bgColor = Interpolate(th.Background, th.Primary, 0.1)
		}
		dims := make([]D, len(widgets))
		call := make([]op.CallOp, len(widgets))
		// Calculate widths
		fracSum := float32(0.0)
		fixSum := float32(0.0)
		for _, w := range weights {
			if w < 1.0 {
				fracSum += w
			} else {
				fixSum += float32(gtx.Px((th.TextSize).Scale(w)))
			}
		}
		scale := (float32(gtx.Constraints.Max.X) - fixSum) / float32(gtx.Px((th.TextSize).Scale(fracSum)))
		// Check child sizes
		for i, child := range widgets {
			c := gtx
			if weights != nil {
				if weights[i] < 1.0 {
					c.Constraints.Max.X = gtx.Px((th.TextSize).Scale(weights[i] * scale))
				} else {
					c.Constraints.Max.X = gtx.Px(th.TextSize.Scale(weights[i]))
				}
			} else {
				c.Constraints.Max.X = gtx.Constraints.Max.X / len(widgets)
			}
			c.Constraints.Min.X = c.Constraints.Max.X
			macro := op.Record(c.Ops)
			dims[i] = child(c)
			call[i] = macro.Stop()
		}
		macro := op.Record(gtx.Ops)
		pos := float32(0)
		// Generate all the rendering commands for the children,
		// translated to correct location
		for i := range widgets {
			trans := op.Offset(f32.Pt(pos, 0)).Push(gtx.Ops)
			call[i].Add(gtx.Ops)
			trans.Pop()
			pos += float32(dims[i].Size.X)
		}
		dim := D{Size: image.Pt(int(pos), dims[0].Size.Y)}
		drawAll := macro.Stop()
		// Draw background
		defer clip.Rect{Max: dim.Size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		gtx.Constraints.Min = dim.Size
		r.LayoutClickable(gtx)
		r.HandleClicks(gtx)
		r.HandleToggle(selected, nil)
		// Then play the macro to draw all the children
		drawAll.Add(gtx.Ops)
		return dim
	}
}
