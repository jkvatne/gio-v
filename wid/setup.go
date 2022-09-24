// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
)

// node defines the widget tree of the form.
type node struct {
	w        *layout.Widget
	children []node
}

func (n *node) addChild(w layout.Widget) {
	n.children = append(n.children, node{w: &w})
}

func makeNode(widgets []layout.Widget) node {
	node := node{}
	for _, w := range widgets {
		node.addChild(w)
	}
	return node
}

// Col makes a column of widgets.
func Col(widgets ...layout.Widget) layout.Widget {
	var children []layout.FlexChild
	node := makeNode(widgets)
	for i := 0; i < len(widgets); i++ {
		wg := *node.children[i].w
		children = append(children, layout.Rigid(func(gtx C) D { return wg(gtx) }))
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
