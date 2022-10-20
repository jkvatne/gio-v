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
