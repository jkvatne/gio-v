// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
)

// Col makes a column of widgets.
func Col(weights []float32, widgets ...layout.Widget) layout.Widget {
	return func(gtx C) D {
		return Flex{Axis: layout.Vertical, Alignment: layout.Start, Spacing: SpaceEnd}.Layout(gtx, weights, widgets...)
	}
}

// Flex lays out child elements along an axis,
// according to alignment and weights.
type Flex struct {
	// Axis is the main axis, either Horizontal or Vertical.
	Axis layout.Axis
	// Spacing controls the distribution of space left after
	// layout.
	Spacing Spacing
	// Alignment is the alignment in the cross axis.
	Alignment layout.Alignment
	// WeightSum is the sum of weights used for the weighted
	// size of Flexed children. If WeightSum is zero, the sum
	// of all Flexed weights is used.
	WeightSum float32
}

// Spacing determine the spacing mode for a Flex.
type Spacing uint8

const (
	// SpaceEnd leaves space at the end.
	SpaceEnd Spacing = iota
	// SpaceStart leaves space at the start.
	SpaceStart
	// SpaceSides shares space between the start and end.
	SpaceSides
	// SpaceAround distributes space evenly between children,
	// with half as much space at the start and end.
	SpaceAround
	// SpaceBetween distributes space evenly between children,
	// leaving no space at the start and end.
	SpaceBetween
	// SpaceEvenly distributes space evenly between children and
	// at the start and end.
	SpaceEvenly
)

// mainConstraint returns the min and max main constraints for axis a.
func mainConstraint(a layout.Axis, cs layout.Constraints) (int, int) {
	if a == layout.Horizontal {
		return cs.Min.X, cs.Max.X
	}
	return cs.Min.Y, cs.Max.Y
}

func crossConstraint(a layout.Axis, cs layout.Constraints) (int, int) {
	if a == layout.Horizontal {
		return cs.Min.Y, cs.Max.Y
	}
	return cs.Min.X, cs.Max.X
}

// constraints returns the constraints for axis a.
func constrain(a layout.Axis, mainMin, mainMax, crossMin, crossMax int) layout.Constraints {
	if a == layout.Horizontal {
		return layout.Constraints{Min: image.Pt(mainMin, crossMin), Max: image.Pt(mainMax, crossMax)}
	}
	return layout.Constraints{Min: image.Pt(crossMin, mainMin), Max: image.Pt(crossMax, mainMax)}
}

// Layout a list of children. The position of the children are
// determined by the specified order, but Rigid children are laid out
// before Flexed children.
func (f Flex) Layout(gtx layout.Context, weights []float32, children ...layout.Widget) layout.Dimensions {
	size := 0
	cs := gtx.Constraints
	mainMin, mainMax := mainConstraint(f.Axis, cs)
	crossMin, crossMax := crossConstraint(f.Axis, cs)
	remaining := mainMax
	var totalWeight float32
	cgtx := gtx
	calls := make([]op.CallOp, len(children))
	dims := make([]layout.Dimensions, len(children))
	// Lay out Rigid children. (with weight==0.0)
	for i, child := range children {
		if i < len(weights) && weights[i] > 0 {
			totalWeight += weights[i]
			continue
		}
		macro := op.Record(gtx.Ops)
		cgtx.Constraints = constrain(f.Axis, 0, remaining, crossMin, crossMax)
		dim := child(cgtx)
		c := macro.Stop()
		sz := f.Axis.Convert(dim.Size).X
		size += sz
		remaining -= sz
		if remaining < 0 {
			remaining = 0
		}
		calls[i] = c
		dims[i] = dim
	}
	if w := f.WeightSum; w != 0 {
		totalWeight = w
	}
	// fraction is the rounding error from a Flex weighting.
	var fraction float32
	flexTotal := remaining
	// Lay out Flexed children (with weight>0)
	for i, child := range children {
		if len(weights) <= i || weights[i] == 0 {
			continue
		}
		var flexSize int
		if remaining > 0 && totalWeight > 0 {
			// Apply weight and add any leftover fraction from a
			// previous Flexed.
			childSize := float32(flexTotal) * weights[i] / totalWeight
			flexSize = int(childSize + fraction + .5)
			fraction = childSize - float32(flexSize)
			if flexSize > remaining {
				flexSize = remaining
			}
		}
		macro := op.Record(gtx.Ops)
		cgtx.Constraints = constrain(f.Axis, flexSize, flexSize, crossMin, crossMax)
		dim := child(cgtx)
		c := macro.Stop()
		sz := f.Axis.Convert(dim.Size).X
		size += sz
		remaining -= sz
		if remaining < 0 {
			remaining = 0
		}
		calls[i] = c
		dims[i] = dim
	}
	maxCross := crossMin
	var maxBaseline int
	for i := range children {
		if c := f.Axis.Convert(dims[i].Size).Y; c > maxCross {
			maxCross = c
		}
		if b := dims[i].Size.Y - dims[i].Baseline; b > maxBaseline {
			maxBaseline = b
		}
	}
	var space int
	if mainMin > size {
		space = mainMin - size
	}
	var mainSize int
	switch f.Spacing {
	case SpaceSides:
		mainSize += space / 2
	case SpaceStart:
		mainSize += space
	case SpaceEvenly:
		mainSize += space / (1 + len(children))
	case SpaceAround:
		if len(children) > 0 {
			mainSize += space / (len(children) * 2)
		}
	}
	for i := range children {
		dims := dims[i]
		b := dims.Size.Y - dims.Baseline
		var cross int
		switch f.Alignment {
		case layout.End:
			cross = maxCross - f.Axis.Convert(dims.Size).Y
		case layout.Middle:
			cross = (maxCross - f.Axis.Convert(dims.Size).Y) / 2
		case layout.Baseline:
			if f.Axis == layout.Horizontal {
				cross = maxBaseline - b
			}
		}
		pt := f.Axis.Convert(image.Pt(mainSize, cross))
		trans := op.Offset(pt).Push(gtx.Ops)
		calls[i].Add(gtx.Ops)
		trans.Pop()
		mainSize += f.Axis.Convert(dims.Size).X
		if i < len(children)-1 {
			switch f.Spacing {
			case SpaceEvenly:
				mainSize += space / (1 + len(children))
			case SpaceAround:
				if len(children) > 0 {
					mainSize += space / len(children)
				}
			case SpaceBetween:
				if len(children) > 1 {
					mainSize += space / (len(children) - 1)
				}
			}
		}
	}
	switch f.Spacing {
	case SpaceSides:
		mainSize += space / 2
	case SpaceEnd:
		mainSize += space
	case SpaceEvenly:
		mainSize += space / (1 + len(children))
	case SpaceAround:
		if len(children) > 0 {
			mainSize += space / (len(children) * 2)
		}
	}
	sz := f.Axis.Convert(image.Pt(mainSize, maxCross))
	sz = cs.Constrain(sz)
	return layout.Dimensions{Size: sz, Baseline: sz.Y - maxBaseline}
}
