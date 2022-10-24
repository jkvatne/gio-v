package wid

import (
	"image"

	"gioui.org/io/semantic"

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
	widget.Clickable
	disabler   *bool
	disabled   bool
	Font       text.Font
	shaper     text.Shaper
	index      *int
	items      []string
	hovered    []bool
	Visible    bool
	wasVisible int
	list       layout.Widget
	Items      []layout.Widget
	icon       *widget.Icon
}

// DropDown returns an initiated struct with drop-dow box setup info
func DropDown(th *Theme, index *int, items []string, options ...Option) layout.Widget {
	b := DropDownStyle{}
	b.icon, _ = widget.NewIcon(icons.NavigationArrowDropDown)
	b.th = th
	b.Font = text.Font{Weight: text.Medium}
	b.shaper = th.Shaper
	b.index = index
	b.items = items
	for i := range items {
		b.Items = append(b.Items, b.option(th, i))
		b.hovered = append(b.hovered, false)
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

	oldVisible := b.Visible
	if !b.Focused() {
		b.Visible = false
	}
	for b.Clicked() {
		b.Visible = !b.Visible
	}
	if b.Visible {
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
		// paintBorder(gtx, listClipRect, b.th.OnBackground, b.th.BorderThickness, 0)
		call.Add(gtx.Ops)
		stack.Pop()
		call = macro.Stop()
		op.Defer(gtx.Ops, call)

	} else {
		b.setHovered()
	}
	pointer.CursorPointer.Add(gtx.Ops)
	return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.Switch.Add(gtx.Ops)
		return dims
	})
}

func (b *DropDownStyle) setHovered() {
	if *b.index >= len(b.hovered) {
		*b.index = len(b.hovered) - 1
	}
	if *b.index < 0 {
		*b.index = 0
	}
	for i := 0; i < len(b.hovered); i++ {
		b.hovered[i] = false
	}
	b.hovered[*b.index] = true
}

func (b *DropDownStyle) option(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		for _, e := range gtx.Events(&b.items[i]) {
			if e, ok := e.(pointer.Event); ok {
				switch e.Type {
				case pointer.Release:
					*b.index = i
					b.Visible = false
					b.wasVisible = 0
					b.hovered[i] = false
				case pointer.Enter:
					for j := 0; j < len(b.hovered); j++ {
						b.hovered[j] = false
					}
					b.hovered[i] = true
				case pointer.Leave:
					b.hovered[i] = false
				case pointer.Cancel:
					b.setHovered()
				}
			}
		}
		if b.hovered[i] {
			c := MulAlpha(b.fgColor, 48)
			if Luminance(b.bgColor) > 28 {
				c = MulAlpha(b.fgColor, 16)
			}
			paint.ColorOp{Color: c}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		paint.ColorOp{Color: b.fgColor}.Add(gtx.Ops)
		lblWidget := func(gtx C) D {
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, text.Font{}, th.TextSize, b.items[i])
		}
		dims := layout.Inset{Top: unit.Dp(2), Left: unit.Dp(th.TextSize * 0.4), Right: unit.Dp(0)}.Layout(gtx, lblWidget)
		defer clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops).Pop()
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
		if b.Focused() || b.Hovered() {
			Shadow(rr, gtx.Dp(b.th.Elevation)).Layout(gtx)
		}
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X - 2,
			Y: gtx.Constraints.Min.Y - 2,
		}}
		paint.FillShape(gtx.Ops, b.bgColor, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		clip.UniformRRect(outline, rr).Push(gtx.Ops).Pop()
		paintBorder(gtx, outline, b.fgColor, b.th.BorderThickness, rr)
		oldIndex := *b.index
		if *b.index > len(b.hovered) {
			*b.index = len(b.hovered) - 1
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
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, b.shaper, b.Font, b.th.TextSize, b.items[*b.index])
		})
	}
}

// LayoutIcon draws the icon
func (b *DropDownStyle) LayoutIcon() layout.Widget {
	return func(gtx C) D {
		size := gtx.Sp(b.th.TextSize * 1.5)
		gtx.Constraints = layout.Exact(image.Pt(size, size))
		return b.icon.Layout(gtx, b.fgColor)
	}
}
