// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/widget"
)

// RadioButtonStyle defines a radio button.
type RadioButtonStyle struct {
	// Widget
	checkable
	Key     string
	handler func(s string)
	Group   *widget.Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, group *widget.Enum, key string, label string, options ...Option) func(gtx C) D {
	r := RadioButtonStyle{
		checkable: checkable{
			Label:              label,
			TextColor:          th.OnSurface,
			IconColor:          th.OnBackground,
			TextSize:           th.TextSize * 14.0 / 16.0,
			Size:               25,
			shaper:             th.Shaper,
			checkedStateIcon:   th.RadioChecked,
			uncheckedStateIcon: th.RadioUnchecked,
		},
		Key:   key,
		Group: group,
	}
	for _, option := range options {
		option.apply(&r)
	}
	return func(gtx C) D {
		return r.Layout(gtx)
	}
}

type RbOption func(style *RadioButtonStyle)

// Do is an optional parameter to set a callback when the button is clicked
func Do(f func(s string)) RbOption {
	return func(b *RadioButtonStyle) {
		b.handler = func(s string) { f(b.Group.Value) }
	}
}

func (b RbOption) apply(cfg interface{}) {
	b(cfg.(*RadioButtonStyle))
}

// Layout updates enum and displays the radio button.
func (r RadioButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	hovered, hovering := r.Group.Hovered()
	focus, focused := r.Group.Focused()
	return r.Group.Layout(gtx, r.Key, func(gtx layout.Context) layout.Dimensions {
		semantic.RadioButton.Add(gtx.Ops)
		highlight := hovering && hovered == r.Key || focused && focus == r.Key
		if r.Group.Changed() {
			r.handler(r.Group.Value)
		}
		return r.layout(gtx, r.Group.Value == r.Key, highlight)
	})
}

/*
func (r *RadioButtonStyle) layout(gtx C, checked bool) D {
	icon := r.CheckedStateIcon
	if !checked {
		icon = r.UncheckedStateIcon
	}

	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					size := gtx.Sp(r.th.TextSize * 1.8)
					if r.Hovered() || r.Focused() {
						paint.FillShape(gtx.Ops,
							MulAlpha(r.th.OnBackground, 70),
							clip.Ellipse{image.Point{}, image.Pt(size, size)}.Op(gtx.Ops))
					}
					return D{Size: image.Point{X: size, Y: size}}
				}),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(1)).Layout(gtx, func(gtx C) D {
						size := gtx.Sp(r.th.TextSize * 1.3)
						gtx.Constraints.Min = image.Point{X: size}
						icon.Layout(gtx, ColDisabled(r.th.OnBackground, gtx.Queue == nil))
						return D{Size: image.Point{X: size, Y: size}}
					})
				}),
			)
		}),

		layout.Rigid(func(gtx C) D {
			return layout.Inset{}.Layout(gtx, func(gtx C) D {
				paint.ColorOp{Color: r.th.OnBackground}.Add(gtx.Ops)
				lbl := r.Label
				if lbl == "" {
					lbl = r.Key
				}
				paint.ColorOp{Color: r.th.OnBackground}.Add(gtx.Ops)
				tl := widget.Label{Alignment: text.Start, MaxLines: 1}
				return tl.Layout(gtx, r.th.Shaper, text.Font{Weight: text.Medium, Style: text.Regular}, r.th.TextSize, lbl)
			})
		}),
	)
	gtx.Constraints.Min = dims.Size
	// r.LayoutClickable(gtx)
	// r.HandleClicks(gtx)
	// r.HandleKeys(gtx)
	return dims
}
*/
