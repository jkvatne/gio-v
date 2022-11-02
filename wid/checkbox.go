// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/op/clip"

	"gioui.org/op/paint"

	"gioui.org/io/semantic"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxDef struct {
	Base
	Clickable
	Label              string
	StrValue           *string
	BoolValue          *bool
	Checked            bool
	TextSize           unit.Sp
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
	Key                string
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, value *string, key string, label string, options ...Option) func(gtx C) D {
	r := CheckBoxDef{
		Label:              label,
		StrValue:           value,
		TextSize:           th.TextSize,
		checkedStateIcon:   th.RadioChecked,
		uncheckedStateIcon: th.RadioUnchecked,
		Key:                key,
	}
	r.th = th
	r.fgColor = th.Fg(Surface)
	r.padding = th.LabelPadding // layout.Inset{d, d, d, d}
	r.Font = &th.DefaultFont
	for _, option := range options {
		option.apply(&r)
	}
	return func(gtx C) D {
		r.HandleEvents(gtx)
		dims := r.Layout(gtx)
		r.SetupEventHandlers(gtx, dims.Size)
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, options ...Option) func(gtx C) D {
	c := &CheckBoxDef{
		Label:              label,
		TextSize:           th.TextSize,
		checkedStateIcon:   th.CheckBoxChecked,
		uncheckedStateIcon: th.CheckBoxUnchecked,
	}
	c.Font = &th.DefaultFont
	c.th = th
	c.fgColor = th.Fg(Surface)
	c.padding = th.LabelPadding
	for _, option := range options {
		option.apply(c)
	}
	return func(gtx C) D {
		c.HandleEvents(gtx)
		dims := c.Layout(gtx)
		c.SetupEventHandlers(gtx, dims.Size)
		pointer.CursorPointer.Add(gtx.Ops)
		return dims
	}
}

// Layout updates the checkBox and displays it.
func (c *CheckBoxDef) Layout(gtx layout.Context) layout.Dimensions {
	for c.Clicked() {
		c.Checked = !c.Checked
		if c.BoolValue != nil {
			*c.BoolValue = c.Checked
		}
		if c.StrValue != nil {
			*c.StrValue = c.Key
		}
		if c.onUserChange != nil {
			c.onUserChange()
		}
	}
	if c.StrValue != nil {
		c.Checked = *c.StrValue == c.Key
	}
	semantic.DisabledOp(gtx.Queue == nil).Add(gtx.Ops)

	icon := c.uncheckedStateIcon
	if c.Checked {
		icon = c.checkedStateIcon
	}

	macro := op.Record(gtx.Ops)
	gtx.Constraints.Min.Y = 0
	labelDim := widget.Label{}.Layout(gtx, c.th.Shaper, *c.Font, c.TextSize, c.Label)
	drawLabel := macro.Stop()
	dx := labelDim.Size.Y / 6
	dy := gtx.Dp(c.padding.Top + 1)
	defer op.Offset(image.Pt(dx, dy)).Push(gtx.Ops).Pop()
	// The hover/focus shadow extends outside the checkbox by 25%
	b := image.Rectangle{Min: image.Pt(-dx, -dx), Max: image.Pt(dx, dx)}
	background := color.NRGBA{}
	if c.Focused() && c.Hovered() {
		background = MulAlpha(c.fgColor, 70)
	} else if c.Focused() {
		background = MulAlpha(c.fgColor, 45)
	} else if c.Hovered() {
		background = MulAlpha(c.fgColor, 35)
	}
	paint.FillShape(gtx.Ops, background, clip.Ellipse(b).Op(gtx.Ops))

	col := c.fgColor
	if gtx.Queue == nil {
		col = Disabled(col)
	}
	gtx.Constraints.Min = image.Point{X: labelDim.Size.Y}
	icon.Layout(gtx, col)
	dims := layout.Dimensions{
		Size: image.Point{X: labelDim.Size.Y + 5, Y: labelDim.Size.Y + gtx.Dp(c.padding.Top+c.padding.Bottom)},
	}
	defer op.Offset(image.Pt(dims.Size.X, 0)).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c.fgColor}.Add(gtx.Ops)
	drawLabel.Add(gtx.Ops)
	return dims
}

// CheckboxOption is options specific to Checkboxes
type CheckboxOption func(w *CheckBoxDef)

// Bool is an option parameter to set the variable updated
func Bool(b *bool) CheckboxOption {
	return func(c *CheckBoxDef) {
		c.BoolValue = b
	}
}

func (e CheckboxOption) apply(cfg interface{}) {
	e(cfg.(*CheckBoxDef))
}
