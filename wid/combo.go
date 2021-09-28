package wid

import (
	"gio-v/f32color"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

type ComboDef struct {
	Clickable
	th           *Theme
	shadow       ShadowStyle
	disabler     *bool
	Font         text.Font
	shaper       text.Shaper
	Width        unit.Value
	items        []string
	hovered      []bool
	Visible      bool
	list         layout.Widget
	options      []layout.Widget
	icon         *Icon
}

func (b *ComboDef) option(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		b.hovered[i] = false
		for _, e := range gtx.Events(&b.items[i]) {
			if e, ok := e.(pointer.Event); ok {
				switch e.Type {
				case pointer.Release:
					b.index = i
					b.Visible = false
				}
			}
		}
		if b.index == i {
			b.hovered[i]=true
		}
		if b.hovered[i] {
			paint.ColorOp{Color: f32color.MulAlpha(b.th.OnBackground, 16)}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		paint.ColorOp{Color: th.OnBackground}.Add(gtx.Ops)
		dims := aLabel{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, text.Font{}, th.TextSize.Scale(0.8), b.items[i])
		pointer.Rect(image.Rect(0,0,dims.Size.X, dims.Size.Y)).Add(gtx.Ops)
		pointer.InputOp{
			Tag:   &b.items[i],
			Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		return dims
	}
}

func Combo(th *Theme, index int, items []string) func(gtx C) D {
	b := ComboDef{}
	b.icon, _ = NewIcon(icons.NavigationArrowDropDown)
	b.Width = unit.Dp(0)
	b.SetupTabs()
	b.th = th
	b.Font = text.Font{Weight: text.Medium}
	b.shadow = Shadow(th.CornerRadius, th.Elevation)
	b.shaper = th.Shaper
	b.index = index
	b.items = items
	for i, _ := range items {
		b.options = append(b.options, b.option(th, i))
		b.hovered = append(b.hovered, false)
	}
	b.list = MakeList(th, layout.Vertical, b.options...)

	return func(gtx C) D {
		dims := b.Layout(gtx)
		for b.Clicked() {
			b.Visible = !b.Visible
		}
		if b.Visible {
			gtx.Constraints.Min = image.Pt(dims.Size.X, dims.Size.Y)
			gtx.Constraints.Max = image.Pt(dims.Size.X, 9999)
			macro := op.Record(gtx.Ops)
			dims2 := b.list(gtx)
			r := f32.Rect(0, 0, float32(dims2.Size.X), float32(dims2.Size.Y))
			call := macro.Stop()
			macro = op.Record(gtx.Ops)
			op.Offset(f32.Pt(0, float32(dims.Size.Y))).Add(gtx.Ops)
			clip.UniformRRect(r, 0).Add(gtx.Ops)
			paint.Fill(gtx.Ops, b.th.Palette.Background)
			// Draw a border around all options
			paintBorder(gtx, r, b.th.Palette.OnBackground, b.th.BorderThickness, unit.Value{})
			for _, ev := range gtx.Events(&b.eventKey) {
				switch ke := ev.(type) {
				case key.FocusEvent:
					b.focused = ke.Focus
				case key.Event:
					if !b.focused || ke.State != key.Press {
						break
					}
					if ke.Name == key.NameUpArrow {
						b.index--
					}
					if ke.Name == key.NameDownArrow {
						b.index++
					}
				}
			}
			key.InputOp{Tag: &b.eventKey, Hint: key.HintAny}.Add(gtx.Ops)
			call.Add(gtx.Ops)
			call = macro.Stop()
			op.Defer(gtx.Ops, call)

		}
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

func (b *ComboDef) Layout(gtx layout.Context) layout.Dimensions {
	b.disabled = false
	if b.disabler != nil && *b.disabler || GlobalDisable {
		gtx = gtx.Disabled()
		b.disabled = true
	}
	min := gtx.Constraints.Min
	if b.Width.V <= 1.0 {
		min.X = gtx.Px(b.Width.Scale(float32(gtx.Constraints.Max.X)))
	} else if min.X < gtx.Px(b.Width) {
		min.X = gtx.Px(b.Width)
	}
	if min.X > gtx.Constraints.Max.X {
		min.X = gtx.Constraints.Max.X
	}
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(b.LayoutBackground()),
		layout.Stacked(
			func(gtx C) D {
				gtx.Constraints.Min = min
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
					gtx,
					layout.Rigid(b.LayoutLabel()),
					layout.Rigid(b.LayoutIcon()),
				)
			}),
		layout.Expanded(b.LayoutClickable),
	)
}

func (b *ComboDef) LayoutBackground() func(gtx C) D {
	return func(gtx C) D {
		if b.Focused() || b.Hovered() {
			b.shadow.Layout(gtx)
		}
		rr := gtx.Pxr(b.th.CornerRadius)
		if rr > float32(gtx.Constraints.Min.Y)/2.0 {
			rr = float32(gtx.Constraints.Min.Y) / 2.0
		}
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		clip.UniformRRect(outline, rr).Add(gtx.Ops)
		paintBorder(gtx, outline, b.th.Palette.Primary, b.th.BorderThickness, b.th.CornerRadius)
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
}

func (b *ComboDef) LayoutLabel() layout.Widget {
	return func(gtx C) D {
		return b.th.LabelInset.Layout(gtx, func(gtx C) D {
			paint.ColorOp{Color: b.th.Palette.Primary}.Add(gtx.Ops)
			if b.index<0 {b.index=0}
			if b.index>=len(b.items) {b.index=len(b.items)-1}
			return aLabel{Alignment: text.Start}.Layout(gtx, b.shaper, b.Font, b.th.TextSize, b.items[b.index])
		})
	}
}

func (b *ComboDef) LayoutIcon() layout.Widget {
	return func(gtx C) D {
		size := gtx.Px(b.th.TextSize.Scale(1.5))
		gtx.Constraints = layout.Exact(image.Pt(size, size))
		return b.icon.Layout(gtx, b.th.Palette.OnBackground)
	}
}
