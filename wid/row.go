// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
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
	children := makeChildren(false, weights, widgets...)
	return func(gtx C) D {
		bgColor := th.Background
		if selected != nil && *selected {
			bgColor = Interpolate(th.Background, th.Primary, 0.1)
		} else if r.Hovered() {
			bgColor = Interpolate(th.Background, th.Primary, 0.05)
		}
		macro := op.Record(gtx.Ops)
		d := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
		call := macro.Stop()
		// Draw background
		defer clip.Rect{Max: d.Size}.Push(gtx.Ops).Pop()
		//if bgColor != th.Background {
		paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		//}
		// Then play the macro to draw children
		gtx.Constraints.Min = d.Size
		r.LayoutClickable(gtx)
		r.HandleClicks(gtx)
		r.HandleToggle(selected, nil)
		call.Add(gtx.Ops)
		return d
	}
}
