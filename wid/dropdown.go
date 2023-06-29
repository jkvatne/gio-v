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

func (d *DropDownStyle) setLabel(s string) {
	d.label = s
}

func (d *DropDownStyle) setBorder(w unit.Dp) {
	d.borderThickness = w
}

// Layout adds padding to a dropdown box drawn with b.layout().
func (d *DropDownStyle) Layout(gtx C) D {
	d.CheckDisable(gtx)

	// Move to offset the external padding around both label and edit
	defer op.Offset(image.Pt(
		gtx.Dp(d.padding.Left),
		gtx.Dp(d.padding.Top))).Push(gtx.Ops).Pop()

	// If a width is given, and it is within constraints, limit size
	if w := gtx.Dp(d.width); w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
	}
	// And reduce the size to make space for the padding
	gtx.Constraints.Min.X -= gtx.Dp(d.padding.Left + d.padding.Right + d.th.InsidePadding.Left + d.th.InsidePadding.Right)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X

	d.HandleEvents(gtx)

	// Check for index range, because tha HandleEvents() function does not know the limits.
	idx := d.GetIndex(len(d.items))

	// Add outside label to the left of the dropdown box
	if d.label != "" {
		o := op.Offset(image.Pt(0, gtx.Dp(d.th.InsidePadding.Top))).Push(gtx.Ops)
		paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
		ctx := gtx
		ctx.Constraints.Max.X = gtx.Sp(d.labelSize)
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
		_ = widget.Label{Alignment: text.End, MaxLines: 1}.Layout(ctx, d.th.Shaper, *d.Font, d.th.TextSize, d.label, colMacro.Stop())
		o.Pop()
		ofs := gtx.Sp(d.labelSize) + gtx.Dp(d.th.InsidePadding.Left)
		// Move space used by label
		defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()
		gtx.Constraints.Max.X -= ofs
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
	}

	// Draw text with top/left padding offset
	textMacro := op.Record(gtx.Ops)
	o := op.Offset(image.Pt(gtx.Dp(d.th.InsidePadding.Left), gtx.Dp(d.th.InsidePadding.Top))).Push(gtx.Ops)
	paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
	tl := widget.Label{Alignment: text.Start, MaxLines: 1}
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
	dims := tl.Layout(gtx, d.th.Shaper, *d.Font, d.th.TextSize, d.items[*d.index], colMacro.Stop())
	o.Pop()
	drawTextMacro := textMacro.Stop()

	// Calculate widget size based on text size and padding, using all available x space
	dims.Size.X = gtx.Constraints.Max.X

	border := image.Rectangle{Max: image.Pt(
		gtx.Constraints.Max.X+gtx.Dp(d.th.InsidePadding.Left+d.th.InsidePadding.Right),
		dims.Size.Y+gtx.Dp(d.th.InsidePadding.Bottom+d.th.InsidePadding.Top))}

	// Draw border. Need to undo previous top padding offset first
	r := gtx.Dp(d.cornerRadius)
	if r > border.Max.Y/2 {
		r = border.Max.Y / 2
	}
	if d.borderThickness > 0 {
		if d.Focused() {
			paintBorder(gtx, border, d.outlineColor, d.borderThickness*2, r)
		} else if d.Hovered() {
			paintBorder(gtx, border, d.outlineColor, d.borderThickness*3/2, r)
		} else {
			paintBorder(gtx, border, d.Fg(), d.th.BorderThickness, r)
		}
	}
	drawTextMacro.Add(gtx.Ops)

	// Draw icon using foreground color
	iconSize := image.Pt(border.Max.Y, border.Max.Y)
	o = op.Offset(image.Pt(border.Max.X-iconSize.X, 0)).Push(gtx.Ops)
	c := gtx
	c.Constraints.Max = iconSize
	c.Constraints.Min = iconSize
	icon.Layout(c, d.Fg())
	o.Pop()

	oldVisible := d.listVisible
	if !d.Focused() {
		d.listVisible = false
	}
	for d.Clicked() {
		d.listVisible = !d.listVisible
	}
	ok := false
	for i := 0; i < len(d.items); i++ {
		ok = ok || d.itemHovered[i]
	}
	if !ok && !d.inList {
		d.listVisible = false
	}
	if d.listVisible {
		gtx.Constraints.Min = image.Pt(border.Max.X, dims.Size.Y)
		// Limit list length to 8 times the gross size of the dropdown
		gtx.Constraints.Max.Y = dims.Size.Y * 8
		gtx.Constraints.Max.X = gtx.Constraints.Min.X

		listMacro := op.Record(gtx.Ops)
		o := d.list(gtx)
		listClipRect := image.Rect(0, 0, border.Max.X, o.Size.Y)
		theListMacro := listMacro.Stop()

		if !oldVisible {
			d.above = WinY-CurrentY < o.Size.Y+dims.Size.Y
			d.setHovered(idx)
		}

		dy := dims.Size.Y + gtx.Dp(d.padding.Top) + gtx.Dp(d.padding.Bottom)
		if d.above {
			dy = -o.Size.Y
		}
		op.Offset(image.Pt(0, dy)).Add(gtx.Ops)

		for _, ev := range gtx.Events(&d.role) {
			if ev, ok := ev.(pointer.Event); ok {
				switch ev.Type {
				case pointer.Enter:
					d.inList = true
				case pointer.Leave:
					d.inList = false
				}
			}
		}

		dropdownMacro := op.Record(gtx.Ops)

		// Fill background and draw list
		cl := clip.Rect{Max: listClipRect.Max}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, d.th.Bg(Canvas))
		theListMacro.Add(gtx.Ops)
		cl.Pop()

		// Handle mouse enter/leave into list
		cl = clip.Rect(listClipRect).Push(gtx.Ops)
		pass := pointer.PassOp{}.Push(gtx.Ops)
		pointer.InputOp{
			Tag:   &d.role,
			Types: pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		cl.Pop()
		pass.Pop()

		// Draw a border around all options
		paintBorder(gtx, listClipRect, d.th.Fg(Outline), d.th.BorderThickness, 0)
		// Save and defer execution
		dropDownListCall := dropdownMacro.Stop()
		op.Defer(gtx.Ops, dropDownListCall)

	} else {
		d.setHovered(idx)
	}
	pointer.CursorPointer.Add(gtx.Ops)
	d.SetupEventHandlers(gtx, dims.Size)

	return D{Size: image.Pt(
		gtx.Constraints.Max.X,
		border.Max.Y-border.Min.Y+gtx.Dp(d.padding.Bottom+d.padding.Top))}
}

func (d *DropDownStyle) setHovered(h int) {
	for i := 0; i < len(d.itemHovered); i++ {
		d.itemHovered[i] = false
	}
	d.itemHovered[h] = true
}

func (d *DropDownStyle) option(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		for _, ev := range gtx.Events(&d.items[i]) {
			if ev, ok := ev.(pointer.Event); ok {
				switch ev.Type {
				case pointer.Release:
					GuiLock.Lock()
					*d.index = i
					GuiLock.Unlock()
					d.listVisible = false
					d.itemHovered[i] = false
				case pointer.Enter:
					for j := 0; j < len(d.itemHovered); j++ {
						d.itemHovered[j] = false
					}
					d.itemHovered[i] = true
				case pointer.Leave:
					d.itemHovered[i] = false
				case pointer.Cancel:
				}
			}
		}
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
		lblWidget := func(gtx C) D {
			m := op.Record(gtx.Ops)
			paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
			colMacro := m.Stop()
			return widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, *d.Font, th.TextSize, d.items[i], colMacro)
		}
		dims := layout.Inset{Top: unit.Dp(4), Left: unit.Dp(th.TextSize * 0.4), Right: unit.Dp(0)}.Layout(gtx, lblWidget)
		defer clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Push(gtx.Ops).Pop()
		c := color.NRGBA{}
		if *d.index == i {
			c = MulAlpha(d.Fg(), 64)
		} else if d.itemHovered[i] {
			c = MulAlpha(d.Fg(), 24)
		}
		paint.ColorOp{Color: c}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		pointer.InputOp{
			Tag:   &d.items[i],
			Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		return dims
	}
}

func init() {
	icon, _ = NewIcon(icons.NavigationArrowDropDown)
}
