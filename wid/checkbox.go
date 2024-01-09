// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/unit"
	"image"
	"image/color"

	"gioui.org/op/clip"

	"gioui.org/op/paint"

	"gioui.org/io/semantic"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxDef struct {
	Base
	Clickable
	Tooltip
	Label              string
	StrValue           *string
	BoolValue          *bool
	Checked            bool
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
	Key                string
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, value *string, key string, label string, options ...Option) func(gtx C) D {
	r := CheckBoxDef{
		Label:              label,
		StrValue:           value,
		checkedStateIcon:   th.RadioChecked,
		uncheckedStateIcon: th.RadioUnchecked,
		Key:                key,
	}
	r.th = th
	r.FontScale = 1.0
	r.role = Surface
	r.margin = th.DefaultMargin
	r.Font = &th.DefaultFont
	r.Tooltip = PlatformTooltip(th)
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
		checkedStateIcon:   th.CheckBoxChecked,
		uncheckedStateIcon: th.CheckBoxUnchecked,
	}
	c.th = th
	c.FontScale = 1.0
	c.role = Surface
	c.margin = th.DefaultMargin
	c.Font = &th.DefaultFont
	c.Tooltip = PlatformTooltip(th)
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
	semantic.EnabledOp(gtx.Queue == nil).Add(gtx.Ops)

	icon := c.uncheckedStateIcon
	if c.Checked {
		icon = c.checkedStateIcon
	}

	iconSize := unit.Sp(c.FontScale) * c.th.FontSp()
	macro := op.Record(gtx.Ops)
	gtx.Constraints.Min.Y = 0
	gtx.Constraints.Min.X = 0
	ctx := gtx
	ctx.Constraints.Max.X -= Px(gtx, c.margin.Right+unit.Dp(iconSize))
	if ctx.Constraints.Max.X < 0 {
		ctx.Constraints.Max.X = 0
	}
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: c.Fg()}.Add(gtx.Ops)
	labelDim := widget.Label{MaxLines: 1}.Layout(ctx, c.th.Shaper, *c.Font, iconSize, c.Label, colMacro.Stop())
	drawLabel := macro.Stop()
	dx := labelDim.Size.Y / 6
	dy := Px(gtx, c.margin.Top)
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
	iconDim := icon.Layout(cgtx, col)
	px := Px(gtx, c.margin.Left+c.margin.Right)
	py := Px(gtx, c.margin.Top+c.margin.Bottom)
	dims := layout.Dimensions{
		Size: image.Point{
			X: labelDim.Size.X + px + iconDim.Size.X,
			Y: labelDim.Size.Y + py,
		}}
	of := op.Offset(image.Pt(labelDim.Size.Y+Px(gtx, c.margin.Left), 0)).Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	drawLabel.Add(gtx.Ops)
	of.Pop()
	c.SetupEventHandlers(gtx, dims.Size)
	pointer.CursorPointer.Add(gtx.Ops)
	_ = c.Tooltip.Layout(gtx, c.hint, func(gtx C) D {
		return D{Size: dims.Size}
	})
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
