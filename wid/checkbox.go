// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"image"
	"image/color"
)

type CheckBoxDef struct {
	Widget
	Clickable
	Label              string
	TextColor          color.NRGBA
	Font               text.Font
	TextSize           unit.Value
	IconColor          color.NRGBA
	Size               unit.Value
	shaper             text.Shaper
	checkedStateIcon   *Icon
	uncheckedStateIcon *Icon
	Value              bool
	changed            bool
}

// Checkbox returns a widget that can be checked, with label, initial state and handler function
func Checkbox(th *Theme, label string, initialState bool, handler func(b bool)) func(gtx C) D {
	c := &CheckBoxDef{
		Label:              label,
		Value:              initialState,
		TextColor:          th.OnBackground,
		IconColor:          th.OnBackground,
		TextSize:           th.TextSize.Scale(1.0),
		Size:               th.TextSize.Scale(1.5),
		shaper:             th.Shaper,
		checkedStateIcon:   th.CheckBoxChecked,
		uncheckedStateIcon: th.CheckBoxUnchecked,
	}
	c.handler = handler
	c.SetupTabs()
	return func(gtx C) D {
		dims :=c.Layout(gtx)
		c.HandleToggle(&c.Value, &c.changed)
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

func (c *CheckBoxDef) Layout(gtx C) D {
	var icon *Icon
	if c.Value {
		icon = c.checkedStateIcon
	} else {
		icon = c.uncheckedStateIcon
	}
	dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					size := gtx.Px(c.Size) * 4 / 3
					dims := layout.Dimensions{
						Size: image.Point{X: size, Y: size},
					}
					if !c.Hovered() && !c.Focused() {
						return dims
					}

					background := MulAlpha(c.IconColor, 70)

					radius := float32(size) / 2
					paint.FillShape(gtx.Ops, background,
						clip.Circle{
							Center: f32.Point{X: radius, Y: radius},
							Radius: radius,
						}.Op(gtx.Ops))

					return dims
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						size := gtx.Px(c.Size)
						col := c.IconColor
						if gtx.Queue == nil {
							col = Disabled(col)
						}
						gtx.Constraints.Min = image.Point{X: size}
						icon.Layout(gtx, col)
						return layout.Dimensions{
							Size: image.Point{X: size, Y: size},
						}
					})
				}),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				paint.ColorOp{Color: c.TextColor}.Add(gtx.Ops)
				tl := aLabel{Alignment: text.Middle, MaxLines: 1}
				return tl.Layout(gtx, c.shaper, c.Font, c.TextSize, c.Label)
			})
		}),
	)
	stack := op.Save(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: dims.Size}).Add(gtx.Ops)
	gtx.Constraints.Min = dims.Size
	c.LayoutClickable(gtx)
	c.LayoutClickable(gtx)
	c.HandleClicks(gtx)
	c.HandleKeys(gtx)
	stack.Load()
	return dims
}

type Enum struct {
	Value    string
	hovered  string
	hovering bool
	changed bool
	clicks []gesture.Click
	values []string
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Changed reports whether Value has changed by user interaction since the last
// call to Changed.
func (e *Enum) Changed() bool {
	changed := e.changed
	e.changed = false
	return changed
}

// Hovered returns the key that is highlighted, or false if none are.
func (e *Enum) Hovered() (string, bool) {
	return e.hovered, e.hovering
}

// Layout adds the event handler for key.
func (e *Enum) Layout(gtx layout.Context, key string) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)

	if index(e.values, key) == -1 {
		e.values = append(e.values, key)
		e.clicks = append(e.clicks, gesture.Click{})
		e.clicks[len(e.clicks)-1].Add(gtx.Ops)
	} else {
		idx := index(e.values, key)
		clk := &e.clicks[idx]
		for _, ev := range clk.Events(gtx) {
			switch ev.Type {
			case gesture.TypeClick:
				if new := e.values[idx]; new != e.Value {
					e.Value = new
					e.changed = true
				}
			}
		}
		if e.hovering && e.hovered == key {
			e.hovering = false
		}
		if clk.Hovered() {
			e.hovered = key
			e.hovering = true
		}
		clk.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: gtx.Constraints.Min}
}
