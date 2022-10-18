// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type checkable struct {
	Label              string
	TextColor          color.NRGBA
	Font               text.Font
	TextSize           unit.Sp
	IconColor          color.NRGBA
	Size               unit.Dp
	shaper             text.Shaper
	checkedStateIcon   *widget.Icon
	uncheckedStateIcon *widget.Icon
}

func (c *checkable) layout(gtx layout.Context, checked, hovered bool) layout.Dimensions {
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
					size := gtx.Dp(c.Size) * 4 / 3
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !hovered {
						return dims
					}

					background := MulAlpha(c.IconColor, 70)

					b := image.Rectangle{Max: image.Pt(size, size)}
					paint.FillShape(gtx.Ops, background, clip.Ellipse(b).Op(gtx.Ops))

					return dims
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(2).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						size := gtx.Dp(c.Size)
						col := c.IconColor
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
				paint.ColorOp{Color: c.TextColor}.Add(gtx.Ops)
				return widget.Label{}.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
			})
		}),
	)
	return dims
}