package wid

import (
	"image"
	"image/color"

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
	disabler        *bool
	items           []string
	itemHovered     []bool
	outlineColor    color.NRGBA
	listVisible     bool
	list            layout.Widget
	Items           []layout.Widget
	label           string
	labelSize       unit.Sp
	above           bool
	borderThickness unit.Dp
	insidePadding   layout.Inset
}

var icon *Icon

// DropDown returns an initiated struct with drop-dow box setup info
func DropDown(th *Theme, index *int, items []string, options ...Option) layout.Widget {
	b := DropDownStyle{}
	b.th = th
	b.role = Canvas
	b.outlineColor = th.Fg(Outline)
	b.Font = &th.DefaultFont
	b.index = index
	b.items = items
	b.labelSize = th.TextSize * 8
	b.insidePadding = th.LabelPadding
	b.borderThickness = b.th.BorderThickness
	for i := range items {
		b.Items = append(b.Items, b.option(th, i))
		b.itemHovered = append(b.itemHovered, false)
	}
	b.list = List(th, Overlay, b.Items...)
	b.cornerRadius = th.BorderCornerRadius
	b.padding = th.DropDownPadding
	for _, option := range options {
		option.apply(&b)
	}
	if b.label == "" {
		b.labelSize = 0
	}
	if b.borderThickness == 0 {
		b.insidePadding = layout.Inset{}
	}
	return b.Layout
}

func (e *DropDownStyle) setLabel(s string) {
	e.label = s
}

func (e *DropDownStyle) setBorder(w unit.Dp) {
	e.borderThickness = w
}

// Layout adds padding to a dropdown box drawn with b.layout().
func (b *DropDownStyle) Layout(gtx C) D {
	return b.layout(gtx)
}

func (b *DropDownStyle) layout(gtx C) D {

	b.HandleEvents(gtx)
	// Check for index range, because tha HandleEvents() function does not know the limits.
	GuiLock.RLock()
	idx := *b.index
	GuiLock.RUnlock()
	if idx < 0 {
		idx = 0
		GuiLock.Lock()
		*b.index = idx
		GuiLock.Unlock()
	}
	if idx >= len(b.items) {
		idx = len(b.items) - 1
		GuiLock.Lock()
		*b.index = idx
		GuiLock.Unlock()
	}

	if b.disabler != nil && *b.disabler {
		gtx = gtx.Disabled()
	}

	// Use all awailable x space, unless a width is given, and it is within constraints.
	w := gtx.Dp(b.width)
	if w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
		gtx.Constraints.Max.X = w
	}

	// Add outside padding
	defer op.Offset(image.Pt(gtx.Dp(b.padding.Left), gtx.Dp(b.padding.Top))).Push(gtx.Ops).Pop()
	o := op.Offset(image.Pt(gtx.Dp(b.insidePadding.Left), gtx.Dp(b.insidePadding.Top))).Push(gtx.Ops)
	paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
	ctx := gtx
	ctx.Constraints.Max.X = gtx.Sp(b.labelSize) - gtx.Dp(b.padding.Left)
	_ = widget.Label{Alignment: text.End, MaxLines: 1}.Layout(ctx, b.th.Shaper, *b.Font, b.th.TextSize, b.label)
	o.Pop()
	ofs := gtx.Sp(b.labelSize) + gtx.Dp(b.padding.Left)
	defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()

	// Draw text with top/left padding offset
	o = op.Offset(image.Pt(gtx.Dp(b.insidePadding.Left), gtx.Dp(b.insidePadding.Top))).Push(gtx.Ops)
	paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
	dims := widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, b.th.Shaper, *b.Font, b.th.TextSize, b.items[*b.index])
	o.Pop()

	// Calculate widget size based on text size and padding, using all available x space
	dims.Size.X = gtx.Constraints.Max.X
	dims.Size.Y = dims.Size.Y + gtx.Dp(b.insidePadding.Top+b.insidePadding.Bottom+b.padding.Top+b.padding.Bottom)

	border := image.Rectangle{Max: image.Pt(
		dims.Size.X-gtx.Dp(b.padding.Left+b.padding.Right)-ofs,
		dims.Size.Y-gtx.Dp(b.padding.Top+b.padding.Bottom))}

	// Draw border. Need to undo previous top padding offset first
	r := gtx.Dp(b.cornerRadius)
	if b.borderThickness > 0 {
		if b.Focused() {
			paintBorder(gtx, border, b.outlineColor, b.borderThickness*2, r)
		} else if b.Hovered() {
			paintBorder(gtx, border, b.outlineColor, b.borderThickness*3/2, r)
		} else {
			paintBorder(gtx, border, b.Fg(), b.th.BorderThickness, r)
		}
	}

	// Draw icon using forground color
	o = op.Offset(image.Pt(border.Max.X-border.Max.Y, 0)).Push(gtx.Ops)
	iconSize := image.Pt(border.Max.Y, border.Max.Y)
	c := gtx
	c.Constraints.Max = iconSize
	icon.Layout(c, b.Fg())
	o.Pop()

	oldVisible := b.listVisible
	if !b.Focused() {
		b.listVisible = false
	}
	for b.Clicked() {
		b.listVisible = !b.listVisible
	}
	ok := false
	for i := 0; i < len(b.items); i++ {
		ok = ok || b.itemHovered[i]
	}
	if !ok {
		b.listVisible = false
	}
	if b.listVisible {
		gtx.Constraints.Min = image.Pt(dims.Size.X, dims.Size.Y)
		gtx.Constraints.Max.Y = gtx.Constraints.Max.Y - dims.Size.Y - 5

		macro := op.Record(gtx.Ops)
		d := b.list(gtx)
		listClipRect := image.Rect(0, 0, gtx.Constraints.Min.X, d.Size.Y)
		call := macro.Stop()

		if !oldVisible {
			b.setHovered(idx)
			b.above = int(mouseY) > (winY - d.Size.Y)
		}

		macro = op.Record(gtx.Ops)
		dy := dims.Size.Y
		if b.above {
			dy = -d.Size.Y
		}
		op.Offset(image.Pt(0, dy)).Add(gtx.Ops)
		stack := clip.UniformRRect(listClipRect, 0).Push(gtx.Ops)
		paint.Fill(gtx.Ops, b.Bg())
		// Draw a border around all options
		paintBorder(gtx, listClipRect, b.th.Fg(Outline), b.th.BorderThicknessActive, 0)
		call.Add(gtx.Ops)
		stack.Pop()
		call = macro.Stop()
		op.Defer(gtx.Ops, call)

	} else {
		b.setHovered(idx)
	}
	pointer.CursorPointer.Add(gtx.Ops)
	b.SetupEventHandlers(gtx, dims.Size)
	return dims
}

func (b *DropDownStyle) setHovered(h int) {

	for i := 0; i < len(b.itemHovered); i++ {
		b.itemHovered[i] = false
	}
	b.itemHovered[h] = true
}

func (b *DropDownStyle) option(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		for _, e := range gtx.Events(&b.items[i]) {
			if e, ok := e.(pointer.Event); ok {
				switch e.Type {
				case pointer.Release:
					GuiLock.Lock()
					*b.index = i
					GuiLock.Unlock()
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
					b.setHovered(i)
				}
			}
		}
		paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
		lblWidget := func(gtx C) D {
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, *b.Font, th.TextSize, b.items[i])
		}
		dims := layout.Inset{Top: unit.Dp(4), Left: unit.Dp(th.TextSize * 0.4), Right: unit.Dp(0)}.Layout(gtx, lblWidget)
		defer clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops).Pop()
		if b.itemHovered[i] {
			c := MulAlpha(b.Fg(), 48)
			if Luminance(b.Bg()) > 28 {
				c = MulAlpha(b.Fg(), 160)
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

// LayoutIcon draws the icon
func (b *DropDownStyle) LayoutIcon() layout.Widget {
	return func(gtx C) D {
		size := gtx.Sp(b.th.TextSize * 1)
		gtx.Constraints = layout.Exact(image.Pt(size, size))
		return icon.Layout(gtx, b.Fg())
	}
}

func init() {
	icon, _ = NewIcon(icons.NavigationArrowDropDown)
}
