package wid

import (
	"gioui.org/layout"
)

// node defines the widget tree of the form.
type node struct {
	w *layout.Widget
	children []*node
}

func (n *node) addChild(w layout.Widget) {
	n.children = append(n.children, &node{w:&w})
}

func MakeList(th *Theme, dir layout.Axis, widgets... layout.Widget) layout.Widget {
	node := node{}
	for _,w := range widgets {
		node.addChild(w)
	}
	listStyle := ListStyle{
		list:           &layout.List{Axis: dir},
		ScrollbarStyle: MakeScrollbarStyle(th),
	}
	return func(gtx C) D {
		var ch []layout.Widget
		for i:=0; i<len(node.children); i++ {
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

func MakeFlex(widgets... layout.Widget) layout.Widget {
	node := node{}
	for _,w := range widgets {
		node.addChild(w)
	}
	return func(gtx C) D {
		var ch []layout.FlexChild
		for i := 0; i < len(node.children); i++ {
			w := *node.children[i].w
			ch = append(ch, layout.Flexed(0.2, func(gtx C) D {
				return w(gtx)
			}))
		}
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx, ch...)
	}
}
