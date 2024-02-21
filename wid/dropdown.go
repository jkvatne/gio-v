package wid

import (
	"gioui.org/io/event"
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
	items        []string
	itemHovered  []bool
	outlineColor color.NRGBA
	listVisible  bool
	inList       bool
	list         Wid
	Items        []Wid
	label        string
	labelSize    float32
	above        bool
}

var (
	dropUpIcon   *Icon
	dropDownIcon *Icon
)

// DropDown returns an initiated struct with drop-dow box setup info
func DropDown(th *Theme, index *int, items []string, options ...Option) layout.Widget {
	b := &DropDownStyle{}
	b.th = th
	b.role = Canvas
	b.outlineColor = th.Fg[Outline]
	b.Font = &th.DefaultFont
	b.index = index
	b.items = items
	b.labelSize = th.LabelSplit
	b.borderWidth = b.th.BorderThickness
	b.ClickMovesFocus = true
	for i := range items {
		b.Items = append(b.Items, b.optionWidget(th, i))
		b.itemHovered = append(b.itemHovered, false)
	}
	b.list = List(th, Overlay, b.Items...)
	b.cornerRadius = th.BorderCornerRadius
	b.margin = th.DefaultMargin
	b.padding = th.DefaultPadding
	for _, option := range options {
		option.apply(b)
	}
	if b.label == "" {
		b.labelSize = 0
	}
	return b.Layout
}

func (d *DropDownStyle) setLabel(s string) {
	d.label = s
}

func (d *DropDownStyle) setLabelSize(w float32) {
	d.labelSize = w
}

func (d *DropDownStyle) Layout(gtx C) D {
	d.CheckDisabler(gtx)
	d.maxIndex = len(d.items)
	// Move to offset the external margin around both label and edit
	defer op.Offset(image.Pt(
		Px(gtx, d.margin.Left),
		Px(gtx, d.margin.Top))).Push(gtx.Ops).Pop()

	// If a width is given, and it is within constraints, limit size
	if w := Px(gtx, d.width); w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
	}
	// And reduce the size to make space for the margin+padding
	gtx.Constraints.Min.X -= Px(gtx, d.padding.Left+d.padding.Right+d.margin.Left+d.margin.Right)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X

	d.HandleEvents(gtx)

	GuiLock.RLock()
	idx := *d.index
	GuiLock.RUnlock()

	// Add outside label to the left of the dropdown box
	if d.label != "" {
		o := op.Offset(image.Pt(0, Px(gtx, d.padding.Top))).Push(gtx.Ops)
		paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
		oldMaxX := gtx.Constraints.Max.X
		ofs := int(float32(oldMaxX) * d.labelSize)
		gtx.Constraints.Max.X = ofs - Px(gtx, d.padding.Left)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: d.Fg()}.Add(gtx.Ops)
		ll := widget.Label{Alignment: text.End, MaxLines: 1}
		ll.Layout(gtx, d.th.Shaper, *d.Font, d.th.TextSize, d.label, colMacro.Stop())
		o.Pop()
		gtx.Constraints.Max.X = oldMaxX - ofs
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()
	}

	// Draw text with top/left padding offset
	textMacro := op.Record(gtx.Ops)
	o := op.Offset(image.Pt(Px(gtx, d.padding.Left), Px(gtx, d.padding.Top))).Push(gtx.Ops)
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
		gtx.Constraints.Max.X+Px(gtx, d.padding.Left+d.padding.Right),
		dims.Size.Y+Px(gtx, d.padding.Bottom+d.padding.Top))}

	// Draw border. Need to undo previous top padding offset first
	r := Px(gtx, d.cornerRadius)
	if r > border.Max.Y/2 {
		r = border.Max.Y / 2
	}
	if d.borderWidth > 0 {
		w := float32(Px(gtx, d.borderWidth))
		if gtx.Focused(&d.Clickable) {
			paintBorder(gtx, border, d.outlineColor, w*2, r)
		} else if d.Hovered() {
			paintBorder(gtx, border, d.outlineColor, w*3/2, r)
		} else {
			paintBorder(gtx, border, d.Fg(), w, r)
		}
	}
	drawTextMacro.Add(gtx.Ops)

	// Draw icon using foreground color
	iconSize := image.Pt(border.Max.Y, border.Max.Y)
	o = op.Offset(image.Pt(border.Max.X-iconSize.X, 0)).Push(gtx.Ops)
	c := gtx
	c.Constraints.Max = iconSize
	c.Constraints.Min = iconSize
	dropDownIcon.Layout(c, d.Fg())
	o.Pop()

	oldVisible := d.listVisible
	if d.listVisible && !gtx.Focused(&d.Clickable) {
		d.listVisible = false
	}
	for d.Clicked() {
		d.listVisible = !d.listVisible
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
			d.above = WinY-mouseY < o.Size.Y+dims.Size.Y
			d.setHovered(idx)
		}

		dy := dims.Size.Y + Px(gtx, d.padding.Top) + Px(gtx, d.padding.Bottom)
		if d.above {
			dy = -o.Size.Y
		}

		for {
			event, ok := gtx.Event(pointer.Filter{
				Target: d,
				Kinds:  pointer.Enter | pointer.Leave,
			})
			if !ok {
				break
			}
			ev, ok := event.(pointer.Event)
			if !ok {
				continue
			}
			switch ev.Kind {
			case pointer.Enter:
				d.listVisible = true
			case pointer.Leave:
				d.listVisible = false
			default:
			}
		}

		dropdownMacro := op.Record(gtx.Ops)
		r := op.Offset(image.Pt(0, dy)).Push(gtx.Ops)

		// Fill background and draw list
		bw := Px(gtx, unit.Dp(1.5))
		cl := clip.Rect{Max: image.Pt(border.Max.X, o.Size.Y+bw)}.Push(gtx.Ops)
		paint.Fill(gtx.Ops, d.th.Bg[Canvas])
		theListMacro.Add(gtx.Ops)
		cl.Pop()
		// Draw frame
		paintBorder(gtx, image.Rect(0, 0, listClipRect.Max.X, listClipRect.Max.Y), d.outlineColor, float32(bw), 0)

		// Handle mouse enter/leave into list area, inluding original value area
		cr := listClipRect
		cr.Min.Y = -dy
		clr := clip.Rect(cr).Push(gtx.Ops)
		pass := pointer.PassOp{}.Push(gtx.Ops)
		event.Op(gtx.Ops, d)
		pass.Pop()
		clr.Pop()

		// Draw a border around all options
		w := float32(Px(gtx, d.borderWidth))
		paintBorder(gtx, listClipRect, d.th.Fg[Outline], w, 0)
		r.Pop()
		// Save and defer execution
		dropDownListCall := dropdownMacro.Stop()
		op.Defer(gtx.Ops, dropDownListCall)
	} else {
		d.setHovered(idx)
	}

	sz := image.Pt(gtx.Constraints.Max.X, border.Max.Y-border.Min.Y+Px(gtx, d.margin.Bottom+d.margin.Top))
	d.SetupEventHandlers(gtx, dims.Size)
	pointer.CursorPointer.Add(gtx.Ops)

	return D{Size: sz}
}

func (d *DropDownStyle) setHovered(h int) {
	for i := 0; i < len(d.itemHovered); i++ {
		d.itemHovered[i] = false
	}
	d.itemHovered[h] = true
}

func (d *DropDownStyle) optionWidget(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		for {
			event, ok := gtx.Event(pointer.Filter{
				Target: &d.itemHovered[i],
				Kinds:  pointer.Release | pointer.Enter | pointer.Leave,
			})
			if !ok {
				break
			}
			ev, ok := event.(pointer.Event)
			if !ok {
				continue
			}
			switch ev.Kind {
			case pointer.Release:
				GuiLock.Lock()
				*d.index = i
				GuiLock.Unlock()
				d.listVisible = false
				d.itemHovered[i] = false
				// Force redraw when item is clicked
				Invalidate()
			case pointer.Enter:
				for j := 0; j < len(d.itemHovered); j++ {
					d.itemHovered[j] = false
				}
				d.itemHovered[i] = true
			case pointer.Leave:
				d.itemHovered[i] = false
			default:
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

		event.Op(gtx.Ops, &d.itemHovered[i])
		return dims
	}
}

func init() {
	dropDownIcon, _ = NewIcon(icons.NavigationArrowDropDown)
	dropUpIcon, _ = NewIcon(icons.NavigationArrowDropUp)
}
