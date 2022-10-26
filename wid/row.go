// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/unit"

	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type rowDef struct {
	widget.Clickable
}

var SpaceClose = []float32{}
var SpaceDistribute []float32 = nil

// Calculate widths
func calcWidths(gtx C, textSize unit.Sp, weights []float32, widths []int) {
	if weights == nil {
		weights = []float32{}
		for range widths {
			weights = append(weights, 1.0)
		}
	} else if len(weights) == 0 {
		for range widths {
			weights = append(weights, 0)
		}
	}
	fracSum := float32(0.0)
	fixSum := float32(0.0)
	for _, w := range weights {
		if w <= 1.0 {
			fracSum += w
		} else {
			fixSum += w
		}
	}
	fixWidth := gtx.Sp(textSize * unit.Sp(fixSum) / 2)
	scale := float32(1.0)
	if fracSum > 0 {
		scale = float32(gtx.Constraints.Max.X-fixWidth) / fracSum
	}
	for i := range widths {
		if i < len(widths) && i < len(weights) {
			if weights != nil {
				if weights[i] <= 1.0 {
					widths[i] = int(weights[i] * scale)
				} else {
					widths[i] = gtx.Sp(textSize * unit.Sp(weights[i]) / 2)
				}
			} else {
				widths[i] = gtx.Constraints.Max.X / len(weights)
			}
		}
	}
}

// Row returns a widget grid row with selectable color.
func Row(th *Theme, selected *bool, weights []float32, widgets ...layout.Widget) layout.Widget {
	r := rowDef{}
	dims := make([]D, len(widgets))
	call := make([]op.CallOp, len(widgets))
	widths := make([]int, len(widgets))
	return func(gtx C) D {
		bgColor := th.Bg(Canvas)
		if r.Hovered() {
			// TODO bgColor = Interpolate(th.Bg(Primary), th.Fg(Primary), 0.05)
		} else if selected != nil && *selected {
			// TODO bgColor = Interpolate(th.Bg(Canvas), th.Fg(Primary), 0.1)
		}
		calcWidths(gtx, th.TextSize, weights, widths)
		// Check child sizes and make macros for each widget in a row
		yMax := 0
		c := gtx
		pos := make([]int, len(widgets)+1)
		for i, child := range widgets {
			if len(widths) > i {
				c.Constraints.Max.X = widths[i]
				if widths[i] == 0 {
					c.Constraints.Max.X = inf
				}
				c.Constraints.Min.X = widths[i]
			} else {
				c.Constraints.Max.X = inf
				c.Constraints.Min.X = 0
			}
			pos[i+1] = pos[i] + widths[i]
			macro := op.Record(c.Ops)
			dims[i] = child(c)
			call[i] = macro.Stop()
			if yMax < dims[i].Size.Y {
				yMax = dims[i].Size.Y
			}
		}
		macro := op.Record(gtx.Ops)
		// Generate all the rendering commands for the children,
		// translated to correct location.
		for i := range widgets {
			trans := op.Offset(image.Pt(pos[i], 0)).Push(gtx.Ops)
			call[i].Add(gtx.Ops)
			trans.Pop()
		}
		// The row width is now the position after the last drawn widget.
		dim := D{Size: image.Pt(pos[len(widgets)], yMax)}
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
