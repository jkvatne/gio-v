// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Col makes a column of widgets. It is not scrollable, but
// weights are used to split the available area.
// Set weigth to 0 for fixed height widgets, and 1 for flexible widgets (like lists)
func Col(weights []float32, widgets ...Wid) Wid {
	offsets := make([]int, len(widgets))
	return func(gtx C) D {
		size := 0
		startY := CurrentY
		var totalWeight float32
		cgtx := gtx
		cgtx.Constraints.Min.Y = 0
		calls := make([]op.CallOp, len(widgets))
		dims := make([]D, len(widgets))
		remaining := gtx.Constraints.Max.Y
		// Lay out Rigid children. (with weight==0.0)
		for i, child := range widgets {
			if i < len(weights) && weights[i] > 0 {
				totalWeight += weights[i]
			} else {
				CurrentY = offsets[i]
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
			if len(weights) <= i || weights[i] == 0 {
				continue
			}
			CurrentY = offsets[i]
			var flexSize int
			if remaining > 0 && totalWeight > 0 {
				childSize := float32(flexTotal) * weights[i] / totalWeight
				flexSize = int(childSize + fraction + .5)
				fraction = childSize - float32(flexSize)
				flexSize = Min(flexSize, remaining)
			}
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
		maxX := gtx.Constraints.Min.X
		for i := range widgets {
			if c := dims[i].Size.X; c > maxX {
				maxX = c
			}
		}
		var mainSize int
		CurrentY = startY
		// Now do the actual drawing, with offsets
		for i := range widgets {
			dims := dims[i]
			offsets[i] = CurrentY
			trans := op.Offset(image.Pt(0, mainSize)).Push(gtx.Ops)
			calls[i].Add(gtx.Ops)
			trans.Pop()
			CurrentY += dims.Size.Y
			mainSize += dims.Size.Y
			if mainSize >= gtx.Constraints.Max.Y {
				break
			}
		}
		// mainSize += Max(0, gtx.Constraints.Min.Y-size)
		sz := gtx.Constraints.Constrain(image.Pt(maxX, mainSize))
		return D{Size: sz, Baseline: sz.Y}
	}
}
