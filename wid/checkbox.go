// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxStyle struct {
	checkable
	CheckBox *widget.Bool
	handler  func(b bool)
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, State *bool, handler func(b bool)) func(gtx C) D {
	c := &CheckBoxStyle{
		handler:  handler,
		CheckBox: new(widget.Bool),
		checkable: checkable{
			Label:              label,
			TextColor:          th.Palette.OnBackground,
			IconColor:          th.Palette.OnBackground,
			TextSize:           th.TextSize * 14.0 / 16.0,
			Size:               unit.Dp(th.TextSize) * 1.5,
			shaper:             th.Shaper,
			checkedStateIcon:   th.CheckBoxChecked,
			uncheckedStateIcon: th.CheckBoxUnchecked,
		},
	}
	c.CheckBox.Value = *State
	return func(gtx C) D {
		dims := c.Layout(gtx)
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the checkBox and displays it.
func (c CheckBoxStyle) Layout(gtx layout.Context) layout.Dimensions {
	return c.CheckBox.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		semantic.CheckBox.Add(gtx.Ops)
		return c.layout(gtx, c.CheckBox.Value, c.CheckBox.Hovered() || c.CheckBox.Focused())
	})
}

/*
func (c *CheckBoxStyle) layout(gtx C) D {
	icon := c.uncheckedStateIcon
	if *c.Value {
		icon = c.checkedStateIcon
	}
	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					size := c.Size * 4 / 3
					dims := D{
						Size: image.Point{X: size, Y: size},
					}
					if !c.checkable.Hovered() && !c.Focused() {
						return dims
					}

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
*/
