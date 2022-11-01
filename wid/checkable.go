// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/op"

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
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
}

func (c *Checkable) layout(gtx layout.Context, checked, hovered bool, focused bool) layout.Dimensions {
	var icon *Icon
	if checked {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Dp(c.Size)
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					// The hover/focus shadow extends outside the checkbox by 25%
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
					size := gtx.Dp(c.Size)
					col := c.fgColor
					if gtx.Queue == nil {
						col = Disabled(col)
					}
					gtx.Constraints.Min = image.Point{X: size}
					defer op.Offset(image.Pt(5, 2)).Push(gtx.Ops).Pop()
					icon.Layout(gtx, col)
					dims := layout.Dimensions{
						Size: image.Point{X: size + 5, Y: size + 2},
					}
					return dims
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(0).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				paint.ColorOp{Color: c.fgColor}.Add(gtx.Ops)
				return widget.Label{}.Layout(gtx, c.th.Shaper, *c.Font, c.TextSize, c.Label)
			})
		}),
	)
	return dims
}
