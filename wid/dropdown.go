package wid

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

// DropDownStyle is the struct for dropdown lists.
type DropDownStyle struct {
	Base
	Clickable
	disabler    *bool
	disabled    bool
	items       []string
	itemHovered []bool
	listVisible bool
	list        layout.Widget
	Items       []layout.Widget
}

var icon *widget.Icon

// DropDown returns an initiated struct with drop-dow box setup info
func DropDown(th *Theme, index *int, items []string, options ...Option) layout.Widget {
	b := DropDownStyle{}
	b.th = th
	b.fgColor = th.Fg(Canvas)
	b.bgColor = th.Bg(Canvas)
	b.Font = &th.DefaultFont
	b.index = index
	b.items = items
	for i := range items {
		b.Items = append(b.Items, b.option(th, i))
		b.itemHovered = append(b.itemHovered, false)
	}
	b.list = List(th, Overlay, b.Items...)
	b.padding = th.LabelPadding
	for _, option := range options {
		option.apply(&b)
	}
	return b.Layout
}

// Layout adds padding to a dropdown box drawn with b.layout().
func (b *DropDownStyle) Layout(gtx C) D {
	return b.padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return b.layout(gtx)
	})
}

func (b *DropDownStyle) layout(gtx C) D {
	if b.width > 0 {
		gtx.Constraints.Min.X = gtx.Dp(b.width)
		gtx.Constraints.Max.X = gtx.Dp(b.width)
	}

	b.disabled = false
	if b.disabler != nil && *b.disabler {
		gtx = gtx.Disabled()
		b.disabled = true
	}
	min := CalcMin(gtx, b.width)
	dims := layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(b.LayoutBackground()),
		layout.Stacked(
			func(gtx C) D {
				if min.X > 0 {
					gtx.Constraints.Min = min
					gtx.Constraints.Max.X = min.X
				}
				if b.width > 0 {
					gtx.Constraints.Max.X = gtx.Dp(b.width)
					gtx.Constraints.Min.X = gtx.Dp(b.width)
				}
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(
					gtx,
					layout.Flexed(1.0, b.LayoutLabel()),
					layout.Rigid(b.LayoutIcon()),
				)
			},
		),
	)

	oldVisible := b.listVisible
	if !b.Focused() {
		b.listVisible = false
	}
	for b.Clicked() {
		b.listVisible = !b.listVisible
	}
	if b.listVisible {
		if !oldVisible {
			b.setHovered()
		}
		gtx.Constraints.Min = image.Pt(dims.Size.X, dims.Size.Y)
		gtx.Constraints.Max.Y = gtx.Constraints.Max.Y - dims.Size.Y - 5
		macro := op.Record(gtx.Ops)
		d := b.list(gtx)
		listClipRect := image.Rect(0, 0, gtx.Constraints.Min.X, d.Size.Y)
		call := macro.Stop()
		macro = op.Record(gtx.Ops)
		op.Offset(image.Pt(0, dims.Size.Y)).Add(gtx.Ops)
		stack := clip.UniformRRect(listClipRect, 0).Push(gtx.Ops)
		paint.Fill(gtx.Ops, b.bgColor)
		// Draw a border around all options
		paintBorder(gtx, listClipRect, b.th.Fg(Outline), b.th.BorderThickness, 0)
		call.Add(gtx.Ops)
		stack.Pop()
		call = macro.Stop()
		op.Defer(gtx.Ops, call)

	} else {
		b.setHovered()
	}
	pointer.CursorPointer.Add(gtx.Ops)
	b.SetupEventHandlers(gtx, dims.Size)
	return dims
}

func (b *DropDownStyle) setHovered() {
	if *b.index >= len(b.itemHovered) {
		*b.index = len(b.itemHovered) - 1
	}
	if *b.index < 0 {
		*b.index = 0
	}
	for i := 0; i < len(b.itemHovered); i++ {
		b.itemHovered[i] = false
	}
	b.itemHovered[*b.index] = true
}

func (b *DropDownStyle) option(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		for _, e := range gtx.Events(&b.items[i]) {
			if e, ok := e.(pointer.Event); ok {
				switch e.Type {
				case pointer.Release:
					*b.index = i
					b.listVisible = false
					b.itemHovered[i] = false
				case pointer.Enter:
					for j := 0; j < len(b.itemHovered); j++ {
						b.itemHovered[j] = false
					}
					b.itemHovered[i] = true
				case pointer.Leave:
					b.itemHovered[i] = false
				case pointer.Cancel:
					b.setHovered()
				}
			}
		}
		paint.ColorOp{Color: b.fgColor}.Add(gtx.Ops)
		lblWidget := func(gtx C) D {
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, *b.Font, th.TextSize, b.items[i])
		}
		dims := layout.Inset{Top: unit.Dp(4), Left: unit.Dp(th.TextSize * 0.4), Right: unit.Dp(0)}.Layout(gtx, lblWidget)
		defer clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops).Pop()
		if b.itemHovered[i] {
			c := MulAlpha(b.fgColor, 48)
			if Luminance(b.bgColor) > 28 {
				c = MulAlpha(b.fgColor, 160)
			}
			paint.ColorOp{Color: c}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		pointer.InputOp{
			Tag:   &b.items[i],
			Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		return dims
	}
}

// LayoutBackground draws the background.
func (b *DropDownStyle) LayoutBackground() func(gtx C) D {
	return func(gtx C) D {
		rr := rr(gtx, b.th.BorderCornerRadius, gtx.Constraints.Min.Y)
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X - 2,
			Y: gtx.Constraints.Min.Y - 2,
		}}
		if b.Focused() || b.Hovered() {
			DrawShadow(gtx, outline, gtx.Dp(b.th.BorderCornerRadius), 11)
		}
		paint.FillShape(gtx.Ops, b.bgColor, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		clip.UniformRRect(outline, rr).Push(gtx.Ops).Pop()
		paintBorder(gtx, outline, b.fgColor, b.th.BorderThickness, rr)
		oldIndex := *b.index
		b.HandleEvents(gtx)
		if *b.index > len(b.itemHovered) {
			*b.index = len(b.itemHovered) - 1
		}
		if *b.index != oldIndex {
			b.setHovered()

		}
		return D{Size: gtx.Constraints.Min}
	}
}

// LayoutLabel draws the label
func (b *DropDownStyle) LayoutLabel() layout.Widget {
	return func(gtx C) D {
		if gtx.Dp(b.width) > gtx.Constraints.Min.X {
			gtx.Constraints.Min.X = gtx.Dp(b.width)
		}
		// A little trick to bring the label closer to the arrow, and avoid a big gap.
		pad := b.th.DropDownPadding
		pad.Right = unit.Dp(-5)
		return pad.Layout(gtx, func(gtx C) D {
			paint.ColorOp{Color: b.fgColor}.Add(gtx.Ops)
			if *b.index < 0 {
				*b.index = 0
			}
			if *b.index >= len(b.items) {
				*b.index = len(b.items) - 1
			}
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, b.th.Shaper, *b.Font, b.th.TextSize, b.items[*b.index])
		})
	}
}

// LayoutIcon draws the icon
func (b *DropDownStyle) LayoutIcon() layout.Widget {
	return func(gtx C) D {
		size := gtx.Sp(b.th.TextSize * 2)
		gtx.Constraints = layout.Exact(image.Pt(size, size))
		return icon.Layout(gtx, b.fgColor)
	}
}

func init() {
	icon, _ = widget.NewIcon(icons.NavigationArrowDropDown)
}
