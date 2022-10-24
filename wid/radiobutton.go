// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

// RadioButtonStyle defines a radio button.
type RadioButtonStyle struct {
	Base
	checkable
	Key   string
	Group *widget.Enum
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, group *widget.Enum, key string, label string, options ...Option) func(gtx C) D {
	r := RadioButtonStyle{
		checkable: checkable{
			Label:              label,
			TextSize:           th.TextSize,
			Size:               unit.Dp(th.TextSize) * 1.5,
			shaper:             th.Shaper,
			checkedStateIcon:   th.RadioChecked,
			uncheckedStateIcon: th.RadioUnchecked,
		},
		Key:   key,
		Group: group,
	}
	r.fgColor = th.Fg(Surface)
	for _, option := range options {
		option.apply(&r)
	}
	return func(gtx C) D {
		return r.Layout(gtx)
	}
}

// Layout updates enum and displays the radio button.
func (r RadioButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	hovered, hovering := r.Group.Hovered()
	focus, focused := r.Group.Focused()
	return r.Group.Layout(gtx, r.Key, func(gtx layout.Context) layout.Dimensions {
		semantic.RadioButton.Add(gtx.Ops)
		highlight := hovering && hovered == r.Key || focused && focus == r.Key
		if r.Group.Changed() {
			if r.onChange != nil {
				r.onChange()
			}
		}
		return r.layout(gtx, r.Group.Value == r.Key, highlight)
	})
}
