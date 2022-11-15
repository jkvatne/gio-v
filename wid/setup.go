// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
)

// Col makes a column of widgets.
func Col(widgets ...layout.Widget) layout.Widget {
	var children []layout.FlexChild
	for i := 0; i < len(widgets); i++ {
		w := widgets[i]
		children = append(children, layout.Rigid(
			func(gtx C) D {
				return w(gtx)
			},
		))
	}
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start, Spacing: layout.SpaceEnd}.Layout(gtx, children...)
	}
}

// Pad adds a padding around a widget.
func Pad(padding layout.Inset, w func(gtx C) D) func(gtx C) D {
	return func(gtx C) D {
		return padding.Layout(gtx, w)
	}
}
