// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/unit"

	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type rowDef struct {
	widget.Clickable
	padTop        unit.Sp
	padBtm        unit.Sp
	gridLineWidth unit.Dp
	gridColor     color.NRGBA
}

// SpaceClose is a shortcut for specifying that the row elements are placed close together, left to right
var SpaceClose []float32

// SpaceDistribute should disribute the widgets on a row evenly, with equal space for each
var SpaceDistribute []float32 = nil

// Calculate widths
func calcWidths(gtx C, textSize unit.Sp, weights []float32, widths []int) {
	w := make([]float32, len(widths))
	for i := 0; i < len(w); i++ {
		if weights == nil {
			// If weights is nil it means all widgets distributed equaly (like 1,1,1,1...)
			w[i] = 1.0
		} else if len(weights) == 0 {
			// If weights has no values, it means use native size (like 0,0,0,0...)
			w[i] = float32(widths[i])
		} else if i < len(weights) && weights[i] > 1.0 {
			// Weights > 1 is given in characters, do rescale to pixels
			w[i] = float32(gtx.Sp(textSize * unit.Sp(weights[i]) / 2))
		} else if i < len(weights) {
			w[i] = weights[i]
		}
		if w[i] == 0.0 {
			w[i] = float32(widths[i])
		}
	}
	fracSum := float32(0.0)
	fixSum := float32(0.0)
	for i, w := range w {
		if w <= 1.0 {
			fracSum += w
		} else {
			if widths[i] > 0 {
				fixSum += float32(widths[i])
			} else {
				fixSum += w
			}
		}
	}
	scale := float32(1.0)
	if fracSum > 0 {
		scale = (float32(gtx.Constraints.Min.X) - fixSum) / fracSum
	}
	for i := range w {
		if w[i] <= 1.0 {
			widths[i] = Max(1, int(w[i]*scale))
		} else {
			widths[i] = Min(int(w[i]), gtx.Constraints.Min.X)
		}
	}
}

// GridRow returns a widget grid row with a grid separating columns and rows
func GridRow(th *Theme, pbgColor *color.NRGBA, gridLineWidth unit.Dp, weights []float32, widgets ...layout.Widget) layout.Widget {
	r := rowDef{}
	bgColor := color.NRGBA{}
	if (pbgColor != nil) && (*pbgColor != color.NRGBA{}) {
		bgColor = *pbgColor
	}
	r.padTop = th.RowPadTop
	r.padBtm = th.RowPadBtm
	r.gridLineWidth = gridLineWidth
	r.gridColor = th.Fg(Outline)
	return func(gtx C) D {
		return r.rowLayout(gtx, th.TextSize, bgColor, weights, widgets...)
	}
}

// Row returns a widget grid row with selectable color.
func Row(th *Theme, pbgColor *color.NRGBA, weights []float32, widgets ...layout.Widget) layout.Widget {
	r := rowDef{}
	bgColor := color.NRGBA{}
	if (pbgColor != nil) && (*pbgColor != color.NRGBA{}) {
		bgColor = *pbgColor
	}
	r.padTop = th.RowPadTop
	r.padBtm = th.RowPadBtm
	return func(gtx C) D {
		return r.rowLayout(gtx, th.TextSize, bgColor, weights, widgets...)
	}
}

func (r *rowDef) rowLayout(gtx C, textSize unit.Sp, bgColor color.NRGBA, weights []float32, widgets ...layout.Widget) D {
	call := make([]op.CallOp, len(widgets))
	dim := make([]D, len(widgets))
	widths := make([]int, len(widgets))
	// Fill in size where width is given as zero
	for i, child := range widgets {
		if i < len(weights) && weights[i] == 0 || len(weights) == 0 {
			c := gtx
			macro := op.Record(c.Ops)
			c.Constraints.Max.X = inf
			c.Constraints.Min.X = 0
			dim := child(c)
			widths[i] = dim.Size.X
			_ = macro.Stop()
		}
	}
	calcWidths(gtx, textSize, weights, widths)
	// Check child sizes and make macros for each widget in a row
	yMax := 0
	c := gtx
	pos := make([]int, len(widgets)+1)
	// For each column in the row
	for i, child := range widgets {
		if len(widths) > i {
			c.Constraints.Max.X = widths[i]
			if widths[i] == 0 {
				c.Constraints.Max.X = inf
			}
			c.Constraints.Min.X = widths[i]
		} else {
			if widths[i] == 0 {
				c.Constraints.Max.X = inf
			}
			c.Constraints.Max.X = widths[i]
			c.Constraints.Min.X = 0
		}
		macro := op.Record(c.Ops)
		dim[i] = child(c)
		if widths[i] < dim[i].Size.X {
			pos[i+1] = pos[i] + dim[i].Size.X
		} else {
			pos[i+1] = pos[i] + widths[i]
		}
		call[i] = macro.Stop()
		if yMax < dim[i].Size.Y {
			yMax = dim[i].Size.Y
		}
	}
	macro := op.Record(gtx.Ops)
	// Generate all the rendering commands for the children,
	// translated to correct location.
	yMax += gtx.Sp(r.padBtm + r.padTop)
	for i := range widgets {
		trans := op.Offset(image.Pt(pos[i], 0)).Push(gtx.Ops)
		call[i].Add(gtx.Ops)
		// Draw a vertical separator
		if r.gridLineWidth > 0 {
			gw := gtx.Dp(r.gridLineWidth)
			outline := image.Rect(gw/2, gw/2, pos[i+1]-pos[i]-gw/2, yMax)
			paint.FillShape(gtx.Ops,
				Black,
				clip.Stroke{
					Path:  clip.Rect(outline).Path(),
					Width: float32(gw),
				}.Op(),
			)
		}
		trans.Pop()
	}
	// The row width is now the position after the last drawn widget + padBtm
	dims := D{Size: image.Pt(pos[len(widgets)], yMax)}
	drawAll := macro.Stop()
	// Draw background.
	defer clip.Rect{Max: image.Pt(dims.Size.X /*gtx.Constraints.Max.X*/, dims.Size.Y)}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	// Skip the top padding by offseting distance padTop
	defer op.Offset(image.Pt(0, gtx.Sp(r.padTop))).Push(gtx.Ops).Pop()
	// Then play the macro to draw all the children.
	drawAll.Add(gtx.Ops)
	return dims
}
