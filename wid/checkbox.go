// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxDef struct {
	Widget
	Clickable
	Label              string
	TextColor          color.NRGBA
	Font               text.Font
	TextSize           unit.Value
	IconColor          color.NRGBA
	Size               unit.Value
	shaper             text.Shaper
	CheckedStateIcon   *Icon
	UncheckedStateIcon *Icon
	Value              *bool
	changed            bool
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, State *bool, handler func(b bool)) func(gtx C) D {
	c := &CheckBoxDef{
		Label:              label,
		Value:              State,
		TextColor:          th.OnBackground,
		IconColor:          th.OnBackground,
		TextSize:           th.TextSize.Scale(1.0),
		Size:               th.TextSize.Scale(1.5),
		shaper:             th.Shaper,
		CheckedStateIcon:   th.CheckBoxChecked,
		UncheckedStateIcon: th.CheckBoxUnchecked,
	}
	c.handler = handler
	c.SetupTabs()
	return func(gtx C) D {
		dims := c.layout(gtx)
		c.HandleToggle(c.Value, &c.changed)
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

func (c *CheckBoxDef) layout(gtx C) D {
	icon := c.UncheckedStateIcon
	if *c.Value {
		icon = c.CheckedStateIcon
	}
	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					size := gtx.Px(c.Size) * 4 / 3
					dims := D{
						Size: image.Point{X: size, Y: size},
					}
					if !c.Hovered() && !c.Focused() {
						return dims
					}

					background := MulAlpha(c.IconColor, 70)

					radius := float32(size) / 2
					paint.FillShape(gtx.Ops, background,
						clip.Circle{
							Center: f32.Point{X: radius, Y: radius},
							Radius: radius,
						}.Op(gtx.Ops))

					return dims
				}),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
						size := gtx.Px(c.Size)
						col := c.IconColor
						if gtx.Queue == nil {
							col = Disabled(col)
						}
						gtx.Constraints.Min = image.Point{X: size}
						if gtx.Constraints.Min.X > gtx.Constraints.Max.X {
							gtx.Constraints.Max.X = gtx.Constraints.Min.X
						}
						icon.Layout(gtx, col)
						return D{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: unit.Dp(2), Right: unit.Dp(0), Bottom: unit.Dp(2), Left: unit.Dp(0)}.Layout(gtx, func(gtx C) D {
				paint.ColorOp{Color: c.TextColor}.Add(gtx.Ops)
				tl := aLabel{Alignment: text.Start, MaxLines: 1}
				return tl.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
			})
		}),
	)
	gtx.Constraints.Min = dims.Size
	c.LayoutClickable(gtx)
	c.HandleClicks(gtx)
	c.HandleKeys(gtx)
	return dims
}
