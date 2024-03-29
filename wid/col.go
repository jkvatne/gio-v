// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
	"gioui.org/op"
	"image"
	"math"
)

// Col makes a column of widgets. It is not scrollable, but
// weights are used to split the available area.
// Set weight to 0 for fixed height widgets, and 1 for flexible widgets (like lists)
func Col(weights []float32, widgets ...Wid) Wid {
	// Interpret the constant SpaceDistribute as many 1.0 weights
	if len(weights) == 1 && weights[0] == 1.0 {
		weights = make([]float32, len(widgets))
		for i := 0; i < len(widgets); i++ {
			weights[i] = 1.0
		}
	}
	// Make weigths[] as long as widgets[]
	for len(weights) < len(widgets) {
		weights = append(weights, 0.0)
	}
	return func(gtx C) D {
		size := 0
		cgtx := gtx
		cgtx.Constraints.Min.Y = 0
		calls := make([]op.CallOp, len(widgets))
		dims := make([]D, len(widgets))
		minY := gtx.Constraints.Min.Y
		// Lay out Rigid children. (with weight==0.0)
		var totalWeight float32
		remaining := gtx.Constraints.Max.Y
		cgtx.Constraints.Min.Y = 0
		for i, child := range widgets {
			totalWeight += weights[i]
			if weights[i] == 0 {
				macro := op.Record(gtx.Ops)
				cgtx.Constraints.Max.Y = remaining
				dims[i] = child(cgtx)
				calls[i] = macro.Stop()
				size += dims[i].Size.Y
				remaining = Max(0, remaining-dims[i].Size.Y)
			}
		}
		// fraction is the rounding error from a Flex weighting.
		var fraction float32
		flexTotal := remaining
		// Lay out Flexed children (with weight>0)
		for i, child := range widgets {
			// Note that flexed children are dropt if there are no remaining spaace
			if weights[i] > 0 && remaining > 0 {
				childSize := float32(flexTotal) * weights[i] / totalWeight
				flexSize := int(childSize + fraction + .5)
				fraction = childSize - float32(flexSize)
				flexSize = Min(flexSize, remaining)
				macro := op.Record(gtx.Ops)
				cgtx.Constraints = layout.Constraints{
					Min: image.Pt(gtx.Constraints.Min.X, 0),
					Max: image.Pt(gtx.Constraints.Max.X, flexSize)}
				// Layout flex rows
				dims[i] = child(cgtx)
				calls[i] = macro.Stop()
				size += dims[i].Size.Y
				remaining = Max(0, remaining-dims[i].Size.Y)
			}
		}
		space := Max(0, minY-size)
		maxX := gtx.Constraints.Min.X
		for i := range widgets {
			if c := dims[i].Size.X; c > maxX {
				maxX = c
			}
		}
		var y float32
		// Now do the actual drawing, with offsets
		for i := range widgets {
			trans := op.Offset(image.Pt(0, int(math.Round(float64(y))))).Push(gtx.Ops)
			calls[i].Add(gtx.Ops)
			trans.Pop()
			if weights[i] == 0 {
				y += float32(dims[i].Size.Y)
			} else {
				y += float32(dims[i].Size.Y) + float32(space)/float32(len(widgets))
				if y >= float32(gtx.Constraints.Max.Y) {
					break
				}
			}
		}
		sz := gtx.Constraints.Constrain(image.Pt(maxX, int(y)))
		return D{Size: sz, Baseline: sz.Y}
	}
}
