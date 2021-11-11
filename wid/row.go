// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type rowDef struct {
	Clickable
}

// MakeRow returns a widget grid row with selectable color
func MakeRow(th *Theme, axis layout.Axis, selected *bool, weight []float32, widgets ...layout.Widget) layout.Widget {
	var ops op.Ops
	var y int
	node := makeNode(widgets)
	gtx := layout.Context{Ops: &ops, Constraints: layout.Constraints{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 3000, Y: 600}}}
	for _, w := range widgets {
		d := w(gtx).Size
		if d.Y > y {
			y = d.Y
		}
	}
	r := rowDef{}

	return func(gtx C) D {
		var children []layout.FlexChild
		for i := 0; i < len(node.children); i++ {
			wg := *node.children[i].w
			w := float32(1.0)
			if len(weight) > i {
				w = weight[i]
			}
			children = append(children, layout.Flexed(w, func(gtx C) D { return wg(gtx) }))
		}
		bgColor := th.Background
		if *selected {
			bgColor = Interpolate(th.Background, th.Primary, 0.1)
		} else if r.Hovered() {
			bgColor = Interpolate(th.Background, th.Primary, 0.05)
		}
		macro := op.Record(gtx.Ops)
		d := layout.Flex{Axis: axis, Alignment: layout.Middle}.Layout(gtx, children...)
		call := macro.Stop()
		defer clip.Rect{Max: d.Size}.Push(gtx.Ops).Pop()
		// Draw background
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		// Then play the macro to draw children
		gtx.Constraints.Min = d.Size
		r.LayoutClickable(gtx)
		r.HandleClicks(gtx)
		r.HandleToggle(selected, nil)
		call.Add(gtx.Ops)
		return d
	}
}
