package wid

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type NodeType int
const (
	ListNode NodeType = iota
	FlexNode
)

// node defines the widget tree of the form.
type node struct {
	nodeType NodeType
	inset float32
	w *layout.Widget
	children []*node
}

func (n *node) addChild(w layout.Widget) {
	n.children = append(n.children, &node{nodeType: 0, w:&w})
}

func MakeList(th *Theme, dir layout.Axis, widgets... layout.Widget) layout.Widget {
	node := node{nodeType: ListNode, inset: 8.0}
	for _,w := range widgets {
		node.addChild(w)
	}
	listStyle := ListStyle{
		list:           &layout.List{Axis: dir},
		ScrollbarStyle: MakeScrollbarStyle(th),
	}
	return func(gtx C) D {return drawList(th, node, listStyle)(gtx)}
}


func drawList(th *Theme, n node,  listStyle ListStyle) func(gtx C) D {
	var ch []layout.Widget
	for i:=0; i<len(n.children); i++ {
		ch = append(ch, *n.children[i].w)
	}
	return func(gtx C) D {
		return listStyle.Layout(
			gtx,
			len(ch),
			func(gtx C, i int) D {
				return layout.UniformInset(unit.Dp(n.inset)).Layout(gtx, ch[i])
			},
		)
	}
}

func MakeFlex(widgets... layout.Widget) layout.Widget {
	node := node{nodeType: FlexNode, inset: 8.0}
	for _,w := range widgets {
		node.addChild(w)
	}
	return func(gtx C) D {return drawFlex(node)(gtx)}
}

func drawFlex(n node) func(gtx C) D {
	var ch []layout.FlexChild
	in := layout.UniformInset(unit.Dp(n.inset))
	for i := 0; i < len(n.children); i++ {
		w := *n.children[i].w
		ch = append(ch, layout.Rigid(func(gtx C) D {
			return in.Layout(gtx, w)
		}))
	}
	return func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx, ch...)
	}
}
