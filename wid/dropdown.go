package wid

import (
	"gioui.org/f32"
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

type DropDownDef struct {
	Clickable
	Widget
	shadow     ShadowStyle
	disabler   *bool
	Font       text.Font
	shaper     text.Shaper
	items      []string
	hovered    []bool
	Visible    bool
	wasVisible int
	list  layout.Widget
	Items []layout.Widget
	icon  *Icon
}

func DropDown(th *Theme, index int, items []string, options ...Option) func(gtx C) D {
	b := DropDownDef{}
	b.icon, _ = NewIcon(icons.NavigationArrowDropDown)
	b.SetupTabs()
	b.th = th
	b.Font = text.Font{Weight: text.Medium}
	b.shadow = Shadow(th.CornerRadius, th.Elevation)
	b.shaper = th.Shaper
	b.index = index
	b.items = items
	for i := range items {
		b.Items = append(b.Items, b.option(th, i))
		b.hovered = append(b.hovered, false)
	}
	b.list = MakeList(th, layout.Vertical, b.Items...)
	b.Pad(5,2,1,2)
	for _, option := range options {
		option.apply(&b)
	}
	return func(gtx C) D {
		dims := b.Layout(gtx)
		for b.Clicked() {
			b.Visible = !b.Visible
		}

		if b.wasVisible > 0 {
			gtx.Constraints.Min = image.Pt(dims.Size.X, dims.Size.Y)
			gtx.Constraints.Max = image.Pt(dims.Size.X, 9999)
			macro := op.Record(gtx.Ops)
			dims2 := b.list(gtx)
			r := f32.Rect(0, 0, float32(dims2.Size.X), float32(dims2.Size.Y))
			call := macro.Stop()
			macro = op.Record(gtx.Ops)
			op.Offset(f32.Pt(0, float32(dims.Size.Y))).Add(gtx.Ops)
			clip.UniformRRect(r, 0).Add(gtx.Ops)
			paint.Fill(gtx.Ops, b.th.Background)
			// Draw a border around all options
			PaintBorder(gtx, r, b.th.OnBackground, b.th.BorderThickness, unit.Value{})
			call.Add(gtx.Ops)
			call = macro.Stop()
			op.Defer(gtx.Ops, call)

		}
		ok := b.Visible
		if b.wasVisible > 0 {
			ok = b.Hovered()
			for _, x := range b.hovered {
				if x {
					ok = true
				}
			}
		}
		if ok && b.wasVisible < 10 {
			b.wasVisible++
			op.InvalidateOp{}.Add(gtx.Ops)
		}
		if !ok && b.wasVisible > 0 {
			b.wasVisible--
			op.InvalidateOp{}.Add(gtx.Ops)
		}
		if b.wasVisible == 0 {
			b.Visible = false
		}
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

func (b *DropDownDef) option(th *Theme, i int) func(gtx C) D {
	return func(gtx C) D {
		for _, e := range gtx.Events(&b.items[i]) {
			if e, ok := e.(pointer.Event); ok {
				switch e.Type {
				case pointer.Release:
					b.index = i
					b.Visible = false
					b.wasVisible = 0
					b.hovered[i] = false
				case pointer.Enter:
					b.hovered[i] = true
				case pointer.Leave:
					b.hovered[i] = false
				case pointer.Cancel:
					//b.hovered[i]=false
				}
			}
		}
		if b.hovered[i] {
			c := MulAlpha(b.th.OnBackground, 48)
			if Luminance(b.th.OnBackground) > 28 {
				c = MulAlpha(b.th.OnBackground, 16)
			}
			paint.ColorOp{Color: c}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		paint.ColorOp{Color: th.OnBackground}.Add(gtx.Ops)
		lblWidget := func(gtx C) D {
			return aLabel{Alignment: text.Start, MaxLines: 1}.Layout(gtx, th.Shaper, text.Font{}, th.TextSize, b.items[i])
		}
		dims := layout.Inset{Top: unit.Dp(2), Left: th.TextSize.Scale(0.4), Right: unit.Dp(2)}.Layout(gtx, lblWidget)
		pointer.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Add(gtx.Ops)
		pointer.InputOp{
			Tag:   &b.items[i],
			Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		return dims
	}
}

func (b *DropDownDef) Layout(gtx C) D {
	b.disabled = false
	if b.disabler != nil && *b.disabler || GlobalDisable {
		gtx = gtx.Disabled()
		b.disabled = true
	}
	min := CalcMin(gtx, b.width)
	return b.padding.Layout(gtx, func(gtx C) D {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(b.LayoutBackground()),
			layout.Stacked(
				func(gtx C) D {
					if min.X!=0 {
						gtx.Constraints.Min = min
						gtx.Constraints.Max.X = min.X
					}
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start}.Layout(
						gtx,
						layout.Flexed(1.0, b.LayoutLabel()),
						layout.Rigid(b.LayoutIcon()),
					)
				}),
			layout.Expanded(b.LayoutClickable),
			layout.Expanded(b.HandleClicks),
			layout.Expanded(b.HandleKeys),
		)
	})
}

func (b *DropDownDef) LayoutBackground() func(gtx C) D {
	return func(gtx C) D {
		if b.Focused() || b.Hovered() {
			Shadow(b.th.CornerRadius, b.th.Elevation).Layout(gtx)
		}
		rr := Pxr(gtx, b.th.CornerRadius)
		if rr > float32(gtx.Constraints.Min.Y)/2.0 {
			rr = float32(gtx.Constraints.Min.Y) / 2.0
		}
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		clip.UniformRRect(outline, rr).Add(gtx.Ops)
		PaintBorder(gtx, outline, b.th.OnBackground, b.th.BorderThickness, b.th.CornerRadius)
		return D{Size: gtx.Constraints.Min}
	}
}

func (b *DropDownDef) LayoutLabel() layout.Widget {
	return func(gtx C) D {
		if gtx.Px(b.width) > gtx.Constraints.Min.X {
			gtx.Constraints.Min.X = gtx.Px(b.width)
		}
		return b.th.LabelInset.Layout(gtx, func(gtx C) D {
			paint.ColorOp{Color: b.th.OnBackground}.Add(gtx.Ops)
			if b.index < 0 {
				b.index = 0
			}
			if b.index >= len(b.items) {
				b.index = len(b.items) - 1
			}
			return aLabel{Alignment: text.Start}.Layout(gtx, b.shaper, b.Font, b.th.TextSize, b.items[b.index])
		})
	}
}

func (b *DropDownDef) LayoutIcon() layout.Widget {
	return func(gtx C) D {
		size := gtx.Px(b.th.TextSize.Scale(1.5))
		gtx.Constraints = layout.Exact(image.Pt(size, size))
		return b.icon.Layout(gtx, b.th.OnBackground)
	}
}
