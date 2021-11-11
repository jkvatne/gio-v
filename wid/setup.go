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

// MakeList makes a horizontal or vertical list
func MakeList(th *Theme, dir layout.Axis, widgets ...layout.Widget) layout.Widget {
	node := makeNode(widgets)
	listStyle := ListStyle{
		list:           &layout.List{Axis: dir},
		ScrollbarStyle: MakeScrollbarStyle(th),
	}
	return func(gtx C) D {
		var ch []layout.Widget
		for i := 0; i < len(node.children); i++ {
			ch = append(ch, *node.children[i].w)
		}
		return listStyle.Layout(
			gtx,
			len(ch),
			func(gtx C, i int) D {
				return ch[i](gtx)
			},
		)
	}
}

func makeChildren(rigid bool, weights []float32, widgets ...layout.Widget) []layout.FlexChild {
	var children []layout.FlexChild
	node := makeNode(widgets)
	for i := 0; i < len(node.children); i++ {
		wg := *node.children[i].w
		w := float32(1.0)
		if len(weights) > i {
			w = weights[i]
		}
		if rigid {
			children = append(children, layout.Rigid(func(gtx C) D { return wg(gtx) }))
		} else {
			children = append(children, layout.Flexed(w, func(gtx C) D { return wg(gtx) }))
		}
	}
	return children
}

// Col makes a column of widgets.
func Col(widgets ...layout.Widget) layout.Widget {
	children := makeChildren(true, nil, widgets...)
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceEnd}.Layout(gtx, children...)
	}
}
