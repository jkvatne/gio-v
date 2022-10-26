// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Checkable struct {
	Base
	Label              string
	TextSize           unit.Sp
	Size               unit.Dp
	checkedStateIcon   *widget.Icon
	uncheckedStateIcon *widget.Icon
}

func (c *Checkable) layout(gtx layout.Context, checked, hovered bool, focused bool) layout.Dimensions {
	var icon *widget.Icon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Dp(c.Size) * 6 / 7
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					b := image.Rectangle{Min: image.Pt(-size/4, -size/4), Max: image.Pt(size*5/4, size*5/4)}
					background := color.NRGBA{}
					if focused && hovered {
						background = MulAlpha(c.fgColor, 70)
					} else if focused {
						background = MulAlpha(c.fgColor, 45)
					} else if hovered {
						background = MulAlpha(c.fgColor, 35)
					}
					paint.FillShape(gtx.Ops, background, clip.Ellipse(b).Op(gtx.Ops))
					return dims
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(c.Size)
						col := c.fgColor
						if gtx.Queue == nil {
							col = Disabled(col)
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
			return layout.UniformInset(2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				paint.ColorOp{Color: c.fgColor}.Add(gtx.Ops)
				return widget.Label{}.Layout(gtx, c.th.Shaper, *c.Font, c.TextSize, c.Label)
			})
		}),
	)
	return dims
}
