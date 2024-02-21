// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/semantic"
	"gioui.org/unit"
	"image"
	"image/color"

	"gioui.org/op/clip"

	"gioui.org/op/paint"

	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/widget"
)

// CheckBoxDef defines a checkbox widget
type CheckBoxDef struct {
	Base
	Clickable
	TooltipDef
	Label              string
	StrValue           *string
	BoolValue          *bool
	Checked            bool
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
	Key                string
}

// RadioButton returns a RadioButton with a label. The key specifies the initial value for the output
func RadioButton(th *Theme, value *string, key string, label string, options ...Option) Wid {
	r := &CheckBoxDef{
		Base: Base{
			th:        th,
			FontScale: 1.0,
			role:      Surface,
			Font:      &th.DefaultFont,
			padding:   th.DefaultPadding,
		},
		TooltipDef:         Tooltip(th),
		Label:              label,
		StrValue:           value,
		checkedStateIcon:   th.RadioChecked,
		uncheckedStateIcon: th.RadioUnchecked,
		Key:                key,
	}
	for _, option := range options {
		option.apply(r)
	}
	return r.Layout
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, options ...Option) Wid {
	c := &CheckBoxDef{
		Base: Base{
			th:        th,
			FontScale: 1.0,
			role:      Surface,
			Font:      &th.DefaultFont,
			padding:   th.DefaultPadding,
		},
		TooltipDef:         Tooltip(th),
		Label:              label,
		checkedStateIcon:   th.CheckBoxChecked,
		uncheckedStateIcon: th.CheckBoxUnchecked,
	}
	for _, option := range options {
		option.apply(c)
	}
	return c.Layout
}

// Layout updates the checkBox and displays it.
func (c *CheckBoxDef) Layout(gtx C) D {
	c.HandleEvents(&c.Clickable, gtx)
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
	GuiLock.RLock()
	if c.BoolValue != nil {
		c.Checked = *c.BoolValue
	} else if c.StrValue != nil {
		c.Checked = *c.StrValue == c.Key
	}
	GuiLock.RUnlock()

	icon := c.uncheckedStateIcon
	if c.Checked {
		icon = c.checkedStateIcon
	}

	iconSize := Px(gtx, c.th.TextSize*unit.Sp(c.FontScale))
	macro := op.Record(gtx.Ops)
	gtx.Constraints.Min.Y = 0
	gtx.Constraints.Min.X = 0
	ctx := gtx
	ctx.Constraints.Max.X -= Px(gtx, c.padding.Right+c.padding.Left) + iconSize
	if ctx.Constraints.Max.X < 10 {
		ctx.Constraints.Max.X = 10
	}
	// Calculate color of text and checkbox
	fgColor := c.Fg()
	if !gtx.Enabled() {
		fgColor = Disabled(fgColor)
	}

	// Draw label into macro
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: fgColor}.Add(gtx.Ops)
	labelDim := widget.Label{MaxLines: 1}.Layout(ctx, c.th.Shaper, *c.Font, c.th.TextSize*unit.Sp(c.FontScale), c.Label, colMacro.Stop())
	drawLabel := macro.Stop()

	// Calculate hover/focus background color
	background := color.NRGBA{}
	if gtx.Focused(&c.Clickable) && c.Hovered() {
		background = MulAlpha(c.Fg(), 90)
	} else if gtx.Focused(&c.Clickable) {
		background = MulAlpha(c.Fg(), 85)
	} else if c.Hovered() {
		background = MulAlpha(c.Fg(), 65)
	}
	// The hover/focus shadow extends outside the checkbox by the padding size
	pt, pb, pl, pr := ScaleInset(gtx, c.padding)
	b := image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(iconSize*10/9+pl+pr, iconSize*10/9+pt+pb)}
	paint.FillShape(gtx.Ops, background, clip.Ellipse(b).Op(gtx.Ops))

	// Icon layout size will be equal to the min x constraint.
	cgtx := gtx
	cgtx.Constraints.Min = image.Point{X: iconSize}
	// Offset for drawing icon
	defer op.Offset(image.Pt(pl, pt+iconSize/9)).Push(gtx.Ops).Pop()
	// Now draw icon
	iconDim := icon.Layout(cgtx, fgColor)
	size := image.Pt(labelDim.Size.X+pl+pl+iconDim.Size.X, labelDim.Size.Y+pt)
	if c.Label != "" {
		size.Y += iconSize / 9
	}
	// Handle events within the calculated size. Must be called before label offset
	c.SetupEventHandlers(gtx, size)
	semantic.CheckBox.Add(gtx.Ops)
	semantic.SelectedOp(c.Checked).Add(gtx.Ops)
	// Extra offset for drawing label
	defer op.Offset(image.Pt(iconSize+iconSize/9, -iconSize/9)).Push(gtx.Ops).Pop()
	drawLabel.Add(gtx.Ops)
	pointer.CursorPointer.Add(gtx.Ops)
	_ = c.TooltipDef.Layout(gtx, c.hint, c.th)
	return D{Size: size}
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
