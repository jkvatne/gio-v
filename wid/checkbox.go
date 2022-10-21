// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/pointer"
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
			TextSize:           th.TextSize,
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
		return c.layout(gtx, c.CheckBox.Value, c.CheckBox.Hovered() || c.CheckBox.Focused())
	})
}
