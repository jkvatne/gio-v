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
)

type ComboDef struct {
	Clickable
	th           *Theme
	shadow       ShadowStyle
	disabler     *bool
	Font         text.Font
	TextSize     unit.Value
	CornerRadius unit.Value
	LabelInset   layout.Inset
	BorderWidth  unit.Value
	shaper       text.Shaper
	Width        unit.Value
	items      []string
	index int
	Visible bool
	list  layout.Widget
}



func Combo(th *Theme, index int, items []string) func(gtx C) D {
	s := th.TextSize.Scale(0.6)
	t := th.TextSize.Scale(0.4)
	c := th.TextSize.Scale(0.2)
	b := ComboDef{}
	b.Width = unit.Dp(300)
	b.SetupTabs()
	b.th = th
	b.TextSize =th.TextSize
	b.Font = text.Font{Weight: text.Bold}
	b.shadow = Shadow(c,c)
	b.CornerRadius = c
	b.BorderWidth = th.TextSize.Scale(0.2)
	b.shaper = th.Shaper
	b.LabelInset = layout.Inset{Top: t, Bottom: t, Left: s, Right: s}
	b.index = index
	b.items = items
	b.list = MakeList(
		th, layout.Vertical,
		Label(th, "Option1", text.Start, 1.0),
		Label(th, "Option2", text.Start, 1.0),
		Label(th, "Option3", text.Start, 1.0),
	)

	return func(gtx C) D {
		dims := b.Layout(gtx)
		b.HandleClick()
		if b.Visible {
			pos := f32.Pt(0, float32(gtx.Constraints.Min.Y))
			macro := op.Record(gtx.Ops)
			dims = b.list(gtx)
			call := macro.Stop()
			outline := f32.Rectangle{Max: f32.Point{
				X: float32(gtx.Constraints.Min.X),
				Y: float32(gtx.Constraints.Min.Y),
			}}
			clip.UniformRRect(outline, 0).Add(gtx.Ops)
			paint.Fill(gtx.Ops, b.th.Palette.Background)
			paintBorder(gtx, outline, b.th.Palette.OnBackground, b.BorderWidth.V, 0)
			macro = op.Record(gtx.Ops)
			op.Offset(pos).Add(gtx.Ops)
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
	b.Visible = true
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(b.LayoutBackground()),
		layout.Stacked(
			func(gtx C) D {
				gtx.Constraints.Min = min
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
					gtx,
					layout.Rigid(b.LayoutLabel()),
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
		rr := float32(gtx.Px(b.CornerRadius))
		if rr > float32(gtx.Constraints.Min.Y)/2.0 {
			rr = float32(gtx.Constraints.Min.Y) / 2.0
		}
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		clip.UniformRRect(outline, rr).Add(gtx.Ops)
		paintBorder(gtx, outline, b.th.Palette.Primary, b.BorderWidth.V, b.CornerRadius.V)
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
}

func (b *ComboDef) LayoutLabel() layout.Widget {
	return func(gtx C) D {
		return b.LabelInset.Layout(gtx, func(gtx C) D {
			paint.ColorOp{Color: b.th.Palette.Primary}.Add(gtx.Ops)
			return aLabel{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.items[b.index])
		})
	}
}
