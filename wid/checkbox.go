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
	r.role = Surface
	r.padding = th.OutsidePadding
	r.Font = &th.DefaultFont
	for _, option := range options {
		option.apply(&r)
	}
	return func(gtx C) D {
		return r.Layout(gtx)
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
	c.th = th
	c.role = Surface
	c.padding = th.OutsidePadding
	c.Font = &th.DefaultFont
	for _, option := range options {
		option.apply(c)
	}
	return func(gtx C) D {
		return c.Layout(gtx)
	}
}

// Layout updates the checkBox and displays it.
func (c *CheckBoxDef) Layout(gtx C) D {
	c.HandleEvents(gtx)
	for c.Clicked() {
		c.Checked = !c.Checked
		GuiLock.Lock()
		if c.BoolValue != nil {
			*c.BoolValue = c.Checked
		} else if c.StrValue != nil {
			*c.StrValue = c.Key
		}
		GuiLock.Unlock()
		if c.onUserChange != nil {
			c.onUserChange()
		}
	}
	if c.BoolValue != nil {
		GuiLock.RLock()
		c.Checked = *c.BoolValue
		GuiLock.RUnlock()
	} else if c.StrValue != nil {
		GuiLock.RLock()
		c.Checked = *c.StrValue == c.Key
		GuiLock.RUnlock()
	}
	semantic.DisabledOp(gtx.Queue == nil).Add(gtx.Ops)

	icon := c.uncheckedStateIcon
	if c.Checked {
		icon = c.checkedStateIcon
	}

	macro := op.Record(gtx.Ops)
	gtx.Constraints.Min.Y = 0
	gtx.Constraints.Min.X = 0
	labelDim := widget.Label{MaxLines: 1}.Layout(gtx, c.th.Shaper, *c.Font, c.TextSize, c.Label)
	drawLabel := macro.Stop()
	dx := labelDim.Size.Y / 6
	dy := gtx.Dp(c.padding.Top + 1)
	defer op.Offset(image.Pt(dx, dy)).Push(gtx.Ops).Pop()
	// The hover/focus shadow extends outside the checkbox by 25%
	b := image.Rectangle{Min: image.Pt(-dx, -dx), Max: image.Pt(labelDim.Size.Y+dx, labelDim.Size.Y+dx)}
	background := color.NRGBA{}
	if c.Focused() && c.Hovered() {
		background = MulAlpha(c.Fg(), 70)
	} else if c.Focused() {
		background = MulAlpha(c.Fg(), 45)
	} else if c.Hovered() {
		background = MulAlpha(c.Fg(), 35)
	}
	paint.FillShape(gtx.Ops, background, clip.Ellipse(b).Op(gtx.Ops))

	col := c.Fg()
	if gtx.Queue == nil {
		col = Disabled(col)
	}
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{X: labelDim.Size.Y}
	icon.Layout(cgtx, col)
	dims := layout.Dimensions{
		Size: image.Point{
			X: labelDim.Size.Y + gtx.Dp(c.padding.Left+c.padding.Right),
			Y: labelDim.Size.Y + gtx.Dp(c.padding.Top+c.padding.Bottom)},
	}
	of := op.Offset(image.Pt(labelDim.Size.Y+gtx.Dp(c.padding.Left), 0)).Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	drawLabel.Add(gtx.Ops)
	if labelDim.Size.X > 0 {
		dims.Size.X += labelDim.Size.X + gtx.Dp(c.padding.Right)
	}
	of.Pop()
	c.SetupEventHandlers(gtx, dims.Size)
	pointer.CursorPointer.Add(gtx.Ops)
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
