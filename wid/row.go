// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type rowDef struct {
	widget.Clickable
}

// Row returns a widget grid row with selectable color.
func Row(th *Theme, selected *bool, weights []float32, widgets ...layout.Widget) layout.Widget {
	r := rowDef{}
	if weights == nil {
		weights = []float32{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	}
	dims := make([]D, len(widgets))
	call := make([]op.CallOp, len(widgets))
	widths := make([]int, len(widgets))
	return func(gtx C) D {
		bgColor := th.Background
		if r.Hovered() {
			bgColor = Interpolate(th.Background, th.Primary, 0.05)
		} else if selected != nil && *selected {
			bgColor = Interpolate(th.Background, th.Primary, 0.1)
		}
		calcWidths(gtx, th.TextSize, weights[:len(widgets)], widths)
		// Check child sizes and make macros for each widget in a row
		yMax := 0
		c := gtx
		for i, child := range widgets {
			if len(widths) > i {
				c.Constraints.Max.X = widths[i]
				c.Constraints.Min.X = c.Constraints.Max.X
			} else {
				c.Constraints.Max.X = inf
				c.Constraints.Min.X = 0
			}
			macro := op.Record(c.Ops)
			dims[i] = child(c)
			call[i] = macro.Stop()
			if yMax < dims[i].Size.Y {
				yMax = dims[i].Size.Y
			}
		}
		macro := op.Record(gtx.Ops)
		pos := 0
		// Generate all the rendering commands for the children,
		// translated to correct location.
		for i := range widgets {
			trans := op.Offset(image.Pt(pos, 0)).Push(gtx.Ops)
			call[i].Add(gtx.Ops)
			trans.Pop()
			pos += dims[i].Size.X
		}
		// The row width is now the position after the last drawn widget.
		dim := D{Size: image.Pt(int(pos), yMax)}
		drawAll := macro.Stop()
		// Draw background.
		defer clip.Rect{Max: dim.Size}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		gtx.Constraints.Min = dim.Size
		// r.LayoutClickable(gtx)
		// r.HandleClicks(gtx)
		// r.HandleToggle(selected, nil)
		// Then play the macro to draw all the children.
		drawAll.Add(gtx.Ops)
		return dim
	}
}
