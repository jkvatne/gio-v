// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image"
	"image/color"

	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/op"
)

type rowDef struct {
	widget.Clickable
	th            *Theme
	padTop        unit.Dp
	padBtm        unit.Dp
	gridLineWidth unit.Dp
	gridColor     color.NRGBA
}

// SpaceClose is a shortcut for specifying that the row elements are placed close together, left to right
var SpaceClose []float32

// SpaceDistribute should disribute the widgets on a row evenly, with equal space for each
var SpaceDistribute = []float32{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0}
var SpaceRightAdjust = []float32{-1.0}
var SpaceCenter = []float32{-2.0}
var FlexInset = layout.Inset{-1, -1, -1, -1}
var NoInset = layout.Inset{}

// calcWidths will calculate widths
func calcWidths(gtx C, textSize unit.Sp, weights []float32, count int) (widths []int) {
	w := make([]float32, count)
	widths = make([]int, count)
	if len(weights) >= 1 && weights[0] == -1 {
		widths[0] = -1
		return widths
	}
	for i := 0; i < len(w); i++ {
		if len(weights) == 0 {
			// If weights is nil, place all widgets as close as possible
			w[i] = 0.0
		} else if len(weights) == 1 && weights[0] == 1.0 {
			// If weights are {1.0} then distribute equaly (like equaly1,1,1,1...)
			w[i] = 1.0
		} else if i < len(weights) && weights[i] > 1.0 {
			// Weights > 1 is given in characters, do rescale to pixels
			w[i] = float32(Px(gtx, textSize*unit.Sp(weights[i])/2))
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
			widths[i] = Max(0, int(w[i]*scale))
		} else {
			widths[i] = Min(int(w[i]), gtx.Constraints.Min.X)
		}
	}
	return widths
}

// GridRow returns a widget grid row with a grid separating columns and rows
func GridRow(th *Theme, option ...interface{}) layout.Widget {
	return Row(th, option...)
}

// Row returns a widget grid row with selectable color.
// func Row(th *Theme, pbgColor *color.NRGBA, weights []float32, widgets ...layout.Widget) layout.Widget {
func Row(th *Theme, option ...interface{}) layout.Widget {
	r := rowDef{
		th:     th,
		padTop: th.RowPadTop,
		padBtm: th.RowPadBtm,
	}
	bgColor := color.NRGBA{}
	var weights []float32
	var widgets []layout.Widget
	i := 0
	for ; i < len(option); i++ {
		if v, ok := option[i].(*color.NRGBA); ok {
			bgColor = *v
		} else if v, ok := option[i].([]float32); ok {
			weights = v
		} else if v, ok := option[i].(unit.Dp); ok {
			r.gridLineWidth = v
		} else if v, ok := option[i].(layout.Widget); ok {
			widgets = append(widgets, v)
		}
	}
	return func(gtx C) D {
		return r.rowLayout(gtx, th.TextSize, bgColor, weights, widgets...)
	}
}

func (r *rowDef) rowLayout(gtx C, textSize unit.Sp, bgColor color.NRGBA, weights []float32, widgets ...layout.Widget) D {
	call := make([]op.CallOp, len(widgets))
	dim := make([]D, len(widgets))
	w := make([]float32, len(widgets))
	copy(w, weights)
	if len(weights) == 1 && w[0] == -2 {
		w[0] = 0
	}
	// Calculate fixed-width columns (with weight=0)
	for i, child := range widgets {
		if i < len(w) && w[i] == 0 {
			macro := op.Record(gtx.Ops)
			c := gtx
			c.Constraints.Min.X = 0
			dim[i] = child(c)
			// Back calculate equivalent width in char widths. Add 1 to avoid rounding errors.
			w[i] = 1 + 2*float32(dim[i].Size.X)/float32(Px(gtx, textSize))
			call[i] = macro.Stop()
		}
	}
	widths := calcWidths(gtx, textSize, w, len(widgets))
	// Check child sizes and make macros for each widget in a row
	yMax := 0
	totSize := 0
	c := gtx
	pos := make([]int, len(widgets)+1)
	// For each column in the row, make macros to draw the widget
	for i, child := range widgets {
		if len(widths) >= 1 && widths[0] <= -1.0 {
			c.Constraints.Min.X = 0
		} else if len(widths) > i {
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
		totSize = pos[i+1]
		call[i] = macro.Stop()
		if yMax < dim[i].Size.Y {
			yMax = dim[i].Size.Y
		}
	}
	if len(widths) >= 1 && widths[0] == -1.0 {
		delta := Max(gtx.Constraints.Max.X-totSize, 0)
		for i := 0; i < len(pos); i++ {
			pos[i] += delta
		}
	} else if len(weights) == 1 && weights[0] == -2.0 {
		delta := Max(gtx.Constraints.Max.X-totSize, 0) / 2
		for i := 0; i < len(pos); i++ {
			pos[i] += delta
		}
	}
	macro := op.Record(gtx.Ops)
	// Generate all the rendering commands for the children,
	// translated to correct location.
	yMax += Px(gtx, r.padBtm+r.padTop)
	for i := range widgets {
		trans := op.Offset(image.Pt(pos[i], 0)).Push(gtx.Ops)
		call[i].Add(gtx.Ops)
		// Draw a vertical separator
		if r.gridLineWidth > 0 {
			gw := Px(gtx, r.gridLineWidth)
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
	defer op.Offset(image.Pt(0, Px(gtx, r.padTop))).Push(gtx.Ops).Pop()
	// Then play the macro to draw all the children.
	drawAll.Add(gtx.Ops)
	return dims
}
