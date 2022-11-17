// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"fmt"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/exp/constraints"
	"image"
	"strconv"
)

// LabelDef is the setup for a label.
type LabelDef struct {
	Base
	// Face defines the text style.
	Font text.Font
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	TextSize unit.Sp
	Stringer func(dp int) string
	dp       int
}

// LabelOption is options specific to Edits.
type LabelOption func(w *LabelDef)

// Bold is an option parameter to set the widget hint (tooltip).
func Bold() LabelOption {
	return func(d *LabelDef) {
		d.Font.Weight = text.Bold
	}
}

// Weight sets the font weight.
func Weight(weight text.Weight) LabelOption {
	return func(d *LabelDef) {
		d.Font.Weight = weight
	}
}

// Middle will align text in the middle.
func Middle() LabelOption {
	return func(d *LabelDef) {
		d.Alignment = text.Middle
	}
}

// End will align text to the end.
func Right() LabelOption {
	return func(d *LabelDef) {
		d.Alignment = text.End
	}
}

func Dp(n int) LabelOption {
	return func(d *LabelDef) {
		d.dp = n
	}
}

func (e LabelOption) apply(cfg interface{}) {
	e(cfg.(*LabelDef))
}

// Layout will draw the label
func (l LabelDef) Layout(gtx C) D {
	paint.ColorOp{Color: l.Fg()}.Add(gtx.Ops)
	tl := widget.Label{Alignment: l.Alignment, MaxLines: l.MaxLines}
	c := gtx
	if l.MaxLines == 1 {
		// This is a hack to avoid splitting the line when only one line is allowed
		c.Constraints.Max.X = inf
	}
	GuiLock.RLock()
	str := l.Stringer(l.dp)
	GuiLock.RUnlock()
	dims := tl.Layout(c, l.th.Shaper, l.Font, l.TextSize*unit.Sp(l.FontSize), str)
	// NB: Use Min.X instead of Max.X in order to fill screen width. Max.X is very large to allow scrolling wide widgets.
	dims.Size.X = Min(gtx.Constraints.Max.X, dims.Size.X)
	if dims.Size.Y > 100 {
		dims.Size.Y++
	}
	return dims
}

// Value returns a widget for a value given by stringer function
func StringerValue(th *Theme, s func(dp int) string, options ...Option) func(gtx C) D {
	w := LabelDef{
		Stringer:  s,
		TextSize:  th.TextSize,
		Alignment: text.Start,
		Font:      text.Font{Weight: text.Medium, Style: text.Regular},
		MaxLines:  0,
	}
	w.padding = th.LabelPadding
	w.th = th
	// Default to Canvas role (typically black for LightMode and white for DarkMode
	w.role = Canvas
	w.FontSize = 1.0
	for _, option := range options {
		option.apply(&w)
	}

	return func(gtx C) D {
		macro := op.Record(gtx.Ops)
		// NB: Use Min.X instead of Max.X in order to fill screen width. Max.X is very large to allow scrolling wide widgets.
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
		dim := w.padding.Layout(gtx, func(gtx C) D {
			return w.Layout(gtx)
		})
		call := macro.Stop()
		defer clip.Rect(image.Rectangle{Max: dim.Size}).Push(gtx.Ops).Pop()
		if w.bgColor != nil {
			paint.Fill(gtx.Ops, w.Bg())
		}
		call.Add(gtx.Ops)
		return dim
	}
}

type Value interface {
	int | float64 | float32 | string | *int | *float64 | *float32 | *string
}

// Label returns a widget for a label showing a string
func Label[V Value](th *Theme, v V, options ...Option) func(gtx C) D {
	if x, ok := any(v).(int); ok {
		s := func(dp int) string { return fmt.Sprintf("%d", x) }
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*int); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			return fmt.Sprintf("%d", *x)
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(float64); ok {
		s := func(dp int) string { return strconv.FormatFloat(float64(x), 'f', dp, 64) }
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*float64); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			return strconv.FormatFloat(float64(*x), 'f', dp, 64)
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(float32); ok {
		s := func(dp int) string { return strconv.FormatFloat(float64(x), 'f', dp, 32) }
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*float32); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			return strconv.FormatFloat(float64(*x), 'f', dp, 32)
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(string); ok {
		s := func(dp int) string { return x }
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*string); ok {
		GuiLock.RLock()
		defer GuiLock.RUnlock()
		s := func(dp int) string { return *x }
		return StringerValue(th, s, options...)
	}
	s := func(dp int) string { return fmt.Sprintf("%v", v) }
	return StringerValue(th, s, options...)
}

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Ordered](x, y T) T {
	if x >= y {
		return x
	}
	return y
}
