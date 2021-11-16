// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// RadioButtonStyle defines a radio button.
type RadioButtonStyle struct {
	Widget
	Clickable
	Key                string
	Output             *string
	Label              string
	CheckedStateIcon   *Icon
	UncheckedStateIcon *Icon
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, output *string, key string, label string, options ...Option) func(gtx C) D {
	r := RadioButtonStyle{
		Label:              label,
		Output:             output,
		Key:                key,
		CheckedStateIcon:   th.RadioChecked,
		UncheckedStateIcon: th.RadioUnchecked,
	}
	r.th = th
	for _, option := range options {
		option.apply(&r)
	}
	r.SetupTabs()
	return func(gtx C) D {
		isSelected := *r.Output == r.Key
		dims := r.layout(gtx, isSelected)
		gtx.Constraints.Min = dims.Size
		for r.Clicked() {
			if r.Output != nil {
				*r.Output = r.Key
			}
			if r.handler != nil {
				r.handler(true)
			}
		}
		return dims
	}
}

type RbOption func(style *RadioButtonStyle)

// Do is an optional parameter to set a callback when the button is clicked
func Do(f func()) RbOption {
	foo := func(b bool) { f() }
	return func(b *RadioButtonStyle) {
		b.handler = foo
	}
}

func (b RbOption) apply(cfg interface{}) {
	b(cfg.(*RadioButtonStyle))
}

func (r *RadioButtonStyle) layout(gtx C, checked bool) D {
	icon := r.CheckedStateIcon
	if !checked {
		icon = r.UncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					size := gtx.Px(r.th.TextSize.Scale(1.8))
					if r.Hovered() || r.Focused() {
						radius := float32(size) / 2
						paint.FillShape(gtx.Ops,
							MulAlpha(r.th.OnBackground, 70),
							clip.Circle{Center: f32.Point{X: radius, Y: radius}, Radius: radius}.Op(gtx.Ops))
					}
					return D{Size: image.Point{X: size, Y: size}}
				}),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
						size := gtx.Px(r.th.TextSize.Scale(1.8))
						gtx.Constraints.Min = image.Point{X: size}
						icon.Layout(gtx, ColDisabled(r.th.OnBackground, gtx.Queue == nil))
						return D{Size: image.Point{X: size, Y: size}}
					})
				}),
			)
		}),

		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: Zv, Right: r.th.TextSize, Bottom: Zv, Left: Zv}.Layout(gtx, func(gtx C) D {
				paint.ColorOp{Color: r.th.OnBackground}.Add(gtx.Ops)
				lbl := r.Label
				if lbl == "" {
					lbl = r.Key
				}
				return Label(r.th, lbl)(gtx) //  text.Start, 1.0
			})
		}),
	)
	gtx.Constraints.Min = dims.Size
	r.LayoutClickable(gtx)
	r.HandleClicks(gtx)
	r.HandleKeys(gtx)
	return dims
}
