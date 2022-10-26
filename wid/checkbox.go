// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxDef struct {
	Checkable
	CheckBox *widget.Bool
	handler  func(b bool)
	State    *bool
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, options ...Option) func(gtx C) D {
	c := &CheckBoxDef{
		CheckBox: new(widget.Bool),
		Checkable: Checkable{
			Label:              label,
			TextSize:           th.TextSize,
			Size:               unit.Dp(th.TextSize) * 1.5,
			checkedStateIcon:   th.CheckBoxChecked,
			uncheckedStateIcon: th.CheckBoxUnchecked,
		},
	}
	c.Font = &th.DefaultFont
	c.th = th
	c.fgColor = th.Fg(Surface)
	for _, option := range options {
		option.apply(c)
	}
	return func(gtx C) D {
		dims := c.Layout(gtx)
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the checkBox and displays it.
func (c CheckBoxDef) Layout(gtx layout.Context) layout.Dimensions {
	return c.CheckBox.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return c.layout(gtx, c.CheckBox.Value, c.CheckBox.Hovered() || c.CheckBox.Focused())
	})
}

// CheckboxOption is options specific to Checkboxes
type CheckboxOption func(w *CheckBoxDef)

// Bool is an option parameter to set the variable updated
func Bool(b *bool) CheckboxOption {
	return func(c *CheckBoxDef) {
		c.State = b
	}
}

func (e CheckboxOption) apply(cfg interface{}) {
	e(cfg.(*CheckBoxDef))
}
