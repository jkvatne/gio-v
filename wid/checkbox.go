// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gio-v/f32color"
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"image"
	"image/color"
)

type CheckBoxDef struct {
	Clickable
	Label              string
	TextColor          color.NRGBA
	Font               text.Font
	TextSize           unit.Value
	IconColor          color.NRGBA
	Size               unit.Value
	shaper             text.Shaper
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
	Value              bool
	changed            bool
}

func (c *CheckBoxDef) Layout(gtx C) D {
	var icon *Icon
	if c.Value {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}
	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Px(c.Size) * 4 / 3
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !c.Hovered() && !c.Focused() {
						return dims
					}

					background := f32color.MulAlpha(c.IconColor, 70)

					radius := float32(size) / 2
					paint.FillShape(gtx.Ops, background,
						clip.Circle{
							Center: f32.Point{X: radius, Y: radius},
							Radius: radius,
						}.Op(gtx.Ops))

					return dims
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						size := gtx.Px(c.Size)
						col := c.IconColor
						if gtx.Queue == nil {
							col = f32color.Disabled(col)
						}
						gtx.Constraints.Min = image.Point{X: size}
						icon.Layout(gtx, col)
						return layout.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				paint.ColorOp{Color: c.TextColor}.Add(gtx.Ops)
				tl := aLabel{Alignment: text.Middle, MaxLines: 1}
				return tl.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
			})
		}),
	)
	stack := op.Save(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	gtx.Constraints.Min = dims.Size
	c.LayoutClickable(gtx)
	stack.Load()
	return dims
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, initialState bool, handler func(b bool)) func(gtx C) D {
	s := &CheckBoxDef{
		Label:              label,
		Value:              initialState,
		TextColor:          th.Palette.OnBackground,
		IconColor:          th.Palette.OnBackground,
		TextSize:           th.TextSize.Scale(1.0),
		Size:               th.TextSize.Scale(1.5),
		shaper:             th.Shaper,
		checkedStateIcon:   th.Icon.CheckBoxChecked,
		uncheckedStateIcon: th.Icon.CheckBoxUnchecked,
	}
	s.handler = handler
	s.SetupTabs()
	return func(gtx C) D {
		dims := s.Layout(gtx)
		s.HandleToggle(&s.Value, &s.changed)
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}
