package wid

import (
	"gioui.org/layout"
	"gioui.org/op"
	"image"
	"log"
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

func MakeFlex(axis layout.Axis, spacing layout.Spacing, widgets... layout.Widget) layout.Widget {
	node := node{}
	var ops op.Ops
	log.Println("MakeFlex()")
	gtx := layout.Context{Ops: &ops, Constraints: layout.Exact(image.Point{1000,50})}
	for _,w := range widgets {
		node.addChild(w)
		dims:=w(gtx)
		log.Println("Dims = %v", dims.Size)
	}
	return func(gtx C) D {
		var children []layout.FlexChild
		for i := 0; i < len(node.children); i++ {
			w := *node.children[i].w
			if axis == layout.Horizontal {
				children = append(children, layout.Flexed(1.0 , func(gtx C) D {return w(gtx)}))
			} else {
				children = append(children, layout.Rigid(func(gtx C) D {return w(gtx)}))
			}
		}
		return layout.Flex{Axis: axis, Alignment: layout.Middle, Spacing: spacing}.Layout(gtx, children...)
	}
}
