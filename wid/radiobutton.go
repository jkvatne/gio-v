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

type RadioButtonStyle struct {
	Widget
	Clickable
	Key                string
	Output             *string
	Label              string
	TextColor          color.NRGBA
	TextSize           unit.Value
	IconColor          color.NRGBA
	Size               unit.Value
	CheckedStateIcon   *Icon
	UncheckedStateIcon *Icon
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, output *string, key string, label string) func(gtx C) D {
	r := RadioButtonStyle{
		Label:              label,
		Output:             output,
		Key:                key,
		TextColor:          th.OnBackground,
		IconColor:          th.OnBackground,
		TextSize:           th.TextSize.Scale(1.0),
		Size:               th.TextSize.Scale(1.5),
		CheckedStateIcon:   th.RadioChecked,
		UncheckedStateIcon: th.RadioUnchecked,
	}
	r.th = th
	r.SetupTabs()
	return func(gtx C) D {
		isSelected := *r.Output == r.Key
		dims := r.layout(gtx, isSelected, r.Hovered())
		gtx.Constraints.Min = dims.Size
		for r.Clicked() {
			if r.Output != nil {
				*r.Output = r.Key
			}
		}
		return dims
	}
}

func (r *RadioButtonStyle) layout(gtx layout.Context, checked, hovered bool) layout.Dimensions {
	icon := r.CheckedStateIcon
	if !checked {
		icon = r.UncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Px(r.Size) * 5 / 4
					dims := layout.Dimensions{Size: image.Point{X: size, Y: size}}
					if r.Hovered() {
						radius := float32(size) / 2
						paint.FillShape(gtx.Ops,
							MulAlpha(r.IconColor, 70),
							clip.Circle{Center: f32.Point{X: radius, Y: radius},Radius: radius,
						}.Op(gtx.Ops))
					}
					return dims
				}),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min = image.Point{X: gtx.Px(r.Size)}
						icon.Layout(gtx, ColDisabled(r.IconColor, gtx.Queue==nil))
						return layout.Dimensions{
							Size: image.Point{X: gtx.Px(r.Size), Y: gtx.Px(r.Size)},
						}
					})
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{zv, r.th.TextSize, zv, zv}.Layout(gtx, func(gtx C) D {
				paint.ColorOp{Color: r.IconColor}.Add(gtx.Ops)
				return Label(r.th, r.Label, text.Start, 1.0)(gtx)
			})
		}),
	)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	gtx.Constraints.Min = dims.Size
	r.LayoutClickable(gtx)
	r.HandleClicks(gtx)
	r.HandleKeys(gtx)
	return dims
}