// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"image"
	"image/color"
)
var zv = unit.Value{}

type RadioButtonStyle struct {
	Widget
	Clickable
	Key                string
	Group              *Enum
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

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func RadioButton(th *Theme, group *Enum, key string, label string) func(gtx C) D {
	r := RadioButtonStyle{
		Label:              label,
		Group:              group,
		Key:                key,
		TextColor:          th.OnBackground,
		IconColor:          th.OnBackground,
		TextSize:           th.TextSize.Scale(1.0),
		Size:               th.TextSize.Scale(1.5),
		shaper:             th.Shaper,
		checkedStateIcon:   th.RadioChecked,
		uncheckedStateIcon: th.RadioUnchecked,
	}
	r.th =th
	r.SetupTabs()
	return func(gtx C) D {
		hovered, hovering := r.Group.Hovered()
		dims := r.layout(gtx, r.Group.Value == r.Key, hovering && hovered == r.Key)
		gtx.Constraints.Min = dims.Size
		r.Group.Layout(gtx, r.Key)
		return dims
	}
}

func (r *RadioButtonStyle) layout(gtx layout.Context, checked, hovered bool) layout.Dimensions {
	var icon *Icon
	if checked {
		icon = r.checkedStateIcon
	} else {
		icon = r.uncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Px(r.Size) * 5 / 4
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !hovered {
						return dims
					}
					background := MulAlpha(r.IconColor, 70)
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
						size := gtx.Px(r.Size)
						col := r.IconColor
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
			return layout.Inset{zv,r.th.TextSize,zv,zv}.Layout(gtx, func(gtx C) D {
				paint.ColorOp{Color: r.IconColor}.Add(gtx.Ops)
				return Label(r.th, r.Label, text.Start, 1.0)(gtx)
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	return dims
}
