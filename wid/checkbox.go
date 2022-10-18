// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxStyle struct {
	Label              string
	Color              color.NRGBA
	Font               text.Font
	TextSize           unit.Sp
	IconColor          color.NRGBA
	Size               unit.Sp
	shaper             text.Shaper
	checkedStateIcon   *widget.Icon
	uncheckedStateIcon *widget.Icon
	CheckBox           *widget.Bool
	Value              *bool
	handler            func(b bool)
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, State *bool, handler func(b bool)) func(gtx C) D {
	c := &CheckBoxStyle{
		Value:              State,
		handler:            handler,
		Label:              label,
		Color:              th.Palette.OnBackground,
		IconColor:          th.Palette.Background,
		TextSize:           th.TextSize * 14.0 / 16.0,
		Size:               th.TextSize * 1.5,
		shaper:             th.Shaper,
		checkedStateIcon:   th.CheckBoxChecked,
		uncheckedStateIcon: th.CheckBoxUnchecked,
	}
	c.Size = th.TextSize * 1.5
	c.handler = handler
	return func(gtx C) D {
		dims := c.layout(gtx)
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

func (c *CheckBoxStyle) layout(gtx C) D {
	icon := c.uncheckedStateIcon
	if *c.Value {
		icon = c.checkedStateIcon
	}
	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					size := gtx.Sp(c.Size) * 4 / 3
					dims := D{
						Size: image.Point{X: size, Y: size},
					}
					// if !c.Hovered() && !c.Focused() {
					//	return dims
					// }

					background := MulAlpha(c.IconColor, 70)
					paint.FillShape(gtx.Ops, background, clip.Ellipse{image.Point{}, image.Pt(size, size)}.Op(gtx.Ops))

					return dims
				}),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
						size := gtx.Sp(c.Size)
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
				paint.ColorOp{Color: c.IconColor}.Add(gtx.Ops)
				tl := widget.Label{Alignment: text.Start, MaxLines: 1}
				return tl.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
			})
		}),
	)
	gtx.Constraints.Min = dims.Size
	// c.LayoutClickable(gtx)
	// c.HandleClicks(gtx)
	// c.HandleKeys(gtx)
	return dims
}
