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
	items           []string
	itemHovered     []bool
	outlineColor    color.NRGBA
	listVisible     bool
	inList          bool
	list            Wid
	Items           []Wid
	label           string
	labelSize       unit.Sp
	above           bool
	borderThickness unit.Dp
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
	b.borderThickness = b.th.BorderThickness
	for i := range items {
		b.Items = append(b.Items, b.option(th, i))
		b.itemHovered = append(b.itemHovered, false)
	}
	b.list = List(th, Overlay, b.Items...)
	b.cornerRadius = th.BorderCornerRadius
	b.padding = th.OutsidePadding
	for _, option := range options {
		option.apply(&b)
	}
	if b.label == "" {
		b.labelSize = 0
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
	b.CheckDisable(gtx)

	// Move to offset the external padding around both label and edit
	defer op.Offset(image.Pt(
		gtx.Dp(b.padding.Left),
		gtx.Dp(b.padding.Top))).Push(gtx.Ops).Pop()

	// If a width is given, and it is within constraints, limit size
	if w := gtx.Dp(b.width); w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
	}
	// And reduce the size to make space for the padding
	gtx.Constraints.Min.X -= gtx.Dp(b.padding.Left + b.padding.Right + b.th.InsidePadding.Left + b.th.InsidePadding.Right)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X

	b.HandleEvents(gtx)

	// Check for index range, because tha HandleEvents() function does not know the limits.
	idx := b.GetIndex(len(b.items))

	// Add outside label to the left of the dropdown box
	if b.label != "" {
		o := op.Offset(image.Pt(0, gtx.Dp(b.th.InsidePadding.Top))).Push(gtx.Ops)
		paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
		ctx := gtx
		ctx.Constraints.Max.X = gtx.Sp(b.labelSize)
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
		_ = widget.Label{Alignment: text.End, MaxLines: 1}.Layout(ctx, b.th.Shaper, *b.Font, b.th.TextSize, b.label, colMacro.Stop())
		o.Pop()
		ofs := gtx.Sp(b.labelSize) + gtx.Dp(b.th.InsidePadding.Left)
		// Move space used by label
		defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()
		gtx.Constraints.Max.X -= ofs
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
	}

	// Draw text with top/left padding offset
	textMacro := op.Record(gtx.Ops)
	o := op.Offset(image.Pt(gtx.Dp(b.th.InsidePadding.Left), gtx.Dp(b.th.InsidePadding.Top))).Push(gtx.Ops)
	paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
	tl := widget.Label{Alignment: text.Start, MaxLines: 1}
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
	dims := tl.Layout(gtx, b.th.Shaper, *b.Font, b.th.TextSize, b.items[*b.index], colMacro.Stop())
	o.Pop()
	drawTextMacro := textMacro.Stop()

	// Calculate widget size based on text size and padding, using all available x space
	dims.Size.X = gtx.Constraints.Max.X

	border := image.Rectangle{Max: image.Pt(
		gtx.Constraints.Max.X+gtx.Dp(b.th.InsidePadding.Left+b.th.InsidePadding.Right),
		dims.Size.Y+gtx.Dp(b.th.InsidePadding.Bottom+b.th.InsidePadding.Top))}

	// Draw border. Need to undo previous top padding offset first
	r := gtx.Dp(b.cornerRadius)
	if r > border.Max.Y/2 {
		r = border.Max.Y / 2
	}
	if b.borderThickness > 0 {
		if b.Focused() {
			paintBorder(gtx, border, b.outlineColor, b.borderThickness*2, r)
		} else if b.Hovered() {
			paintBorder(gtx, border, b.outlineColor, b.borderThickness*3/2, r)
		} else {
			paintBorder(gtx, border, b.Fg(), b.th.BorderThickness, r)
		}
	}
	drawTextMacro.Add(gtx.Ops)

	// Draw icon using foreground color
	iconSize := image.Pt(border.Max.Y, border.Max.Y)
	o = op.Offset(image.Pt(border.Max.X-iconSize.X, 0)).Push(gtx.Ops)
	c := gtx
	c.Constraints.Max = iconSize
	c.Constraints.Min = iconSize
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
	if !ok && !b.inList {
		b.listVisible = false
	}
	if b.listVisible {
		gtx.Constraints.Min = image.Pt(border.Max.X, dims.Size.Y)
		// Limit list length to 8 times the gross size of the dropdow
		gtx.Constraints.Max.Y = dims.Size.Y * 8
		gtx.Constraints.Max.X = gtx.Constraints.Min.X

		listMacro := op.Record(gtx.Ops)
		d := b.list(gtx)
		listClipRect := image.Rect(0, 0, border.Max.X, d.Size.Y)
		theListMacro := listMacro.Stop()

		if !oldVisible {
			b.above = WinY-CurrentY < d.Size.Y+dims.Size.Y
			b.setHovered(idx)
		}

		dy := dims.Size.Y + gtx.Dp(b.padding.Top) + gtx.Dp(b.padding.Bottom)
		if b.above {
			dy = -d.Size.Y
		}
		op.Offset(image.Pt(0, dy)).Add(gtx.Ops)

		for _, e := range gtx.Events(&b.role) {
			if e, ok := e.(pointer.Event); ok {
				switch e.Type {
				case pointer.Enter:
					b.inList = true
				case pointer.Leave:
					b.inList = false
				}
			}
		}

		dropdownMacro := op.Record(gtx.Ops)

		// Fill background and draw list
		cl := clip.Rect{Max: listClipRect.Max}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, b.th.Bg(Canvas))
		theListMacro.Add(gtx.Ops)
		cl.Pop()

		// Handle mouse enter/leave into list
		cl = clip.Rect(listClipRect).Push(gtx.Ops)
		pass := pointer.PassOp{}.Push(gtx.Ops)
		pointer.InputOp{
			Tag:   &b.role,
			Types: pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		cl.Pop()
		pass.Pop()

		// Draw a border around all options
		paintBorder(gtx, listClipRect, b.th.Fg(Outline), b.th.BorderThickness, 0)
		// Save and defer execution
		dropDownListCall := dropdownMacro.Stop()
		op.Defer(gtx.Ops, dropDownListCall)

	} else {
		b.setHovered(idx)
	}
	pointer.CursorPointer.Add(gtx.Ops)
	b.SetupEventHandlers(gtx, dims.Size)

	return D{Size: image.Pt(
		gtx.Constraints.Max.X,
		border.Max.Y-border.Min.Y+gtx.Dp(b.padding.Bottom+b.padding.Top))}
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
				}
			}
		}
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
		lblWidget := func(gtx C) D {
			m := op.Record(gtx.Ops)
			paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
			colMacro := m.Stop()
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, *b.Font, th.TextSize, b.items[i], colMacro)
		}
		dims := layout.Inset{Top: unit.Dp(4), Left: unit.Dp(th.TextSize * 0.4), Right: unit.Dp(0)}.Layout(gtx, lblWidget)
		defer clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops).Pop()
		c := color.NRGBA{}
		if *b.index == i {
			c = MulAlpha(b.Fg(), 64)
		} else if b.itemHovered[i] {
			c = MulAlpha(b.Fg(), 24)
		}
		paint.ColorOp{Color: c}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		pointer.InputOp{
			Tag:   &b.items[i],
			Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		return dims
	}
}

func init() {
	icon, _ = NewIcon(icons.NavigationArrowDropDown)
}
