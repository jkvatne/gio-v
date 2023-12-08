// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"fmt"
	"gioui.org/font"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"image"
	"math"
	"strconv"
)

// LabelDef is the setup for a label.
type LabelDef struct {
	Base
	// Face defines the text style.
	Font font.Font
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Stringer func(dp int) string
}

// LabelOption is options specific to Edits.
type LabelOption func(w *LabelDef)

// Bold is an option parameter to set the widget hint (tooltip).
func Bold() LabelOption {
	return func(d *LabelDef) {
		d.Font.Weight = font.Bold
	}
}

// Weight sets the font weight.
func Weight(weight font.Weight) LabelOption {
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

// Right will align text to the end.
func Right() LabelOption {
	return func(d *LabelDef) {
		d.Alignment = text.End
	}
}

func (e LabelOption) apply(cfg interface{}) {
	e(cfg.(*LabelDef))
}

// StringerValue returns a widget for a value given by stringer function
func StringerValue(th *Theme, s func(dp int) string, options ...Option) func(gtx C) D {
	w := LabelDef{
		Base: Base{
			th:        th,
			role:      Surface,
			padding:   th.OutsidePadding,
			FontScale: 1.0,
		},
		Stringer:  s,
		Alignment: text.Start,
		Font:      th.DefaultFont,
		MaxLines:  0,
	}

	// Apply options after initialization of LabelDef
	for _, option := range options {
		option.apply(&w)
	}

	return func(gtx C) D {
		if w.MaxLines == 1 {
			// This is a hack to avoid splitting the line when only one line is allowed
			gtx.Constraints.Max.X = inf
		}
		GuiLock.RLock()
		str := w.Stringer(w.Dp)
		GuiLock.RUnlock()
		macro := op.Record(gtx.Ops)
		o := op.Offset(image.Pt(Px(gtx, w.padding.Left), Px(gtx, w.padding.Top))).Push(gtx.Ops)
		tl := widget.Label{Alignment: w.Alignment, MaxLines: w.MaxLines}
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: w.Fg()}.Add(gtx.Ops)
		dims := tl.Layout(gtx, w.th.Shaper, w.Font, unit.Sp(w.FontScale)*w.th.FontSp(), str, colMacro.Stop())
		gtx.Constraints.Min.X -= Px(gtx, w.padding.Left+w.padding.Right)
		dims.Size.X = Min(gtx.Constraints.Max.X, dims.Size.X)
		o.Pop()
		// dims.Size.X += Px(gtx, w.padding.Left+w.padding.Right)
		dims.Size.Y += Px(gtx, w.padding.Bottom)
		call := macro.Stop()
		// Color background into the calculated size
		defer clip.Rect(image.Rectangle{Max: dims.Size}).Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, w.Bg())
		// Then do the text painting (apply the "call" macro)
		call.Add(gtx.Ops)
		return dims
	}
}

type Value interface {
	int | float64 | float32 | string | *int | *float64 | *float32 | *string
}

// Label returns a widget for a label showing a string
func Label[V Value](th *Theme, v V, options ...Option) func(gtx C) D {
	if x, ok := any(v).(int); ok {
		s := func(dp int) string {
			if x == math.MinInt {
				return "---"
			}
			return fmt.Sprintf("%d", x)
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*int); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			if *x == math.MinInt {
				return "---"
			}
			return fmt.Sprintf("%d", *x)
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(float64); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			if x == math.MaxFloat64 {
				return "---"
			} else {
				return strconv.FormatFloat(x, 'f', dp, 64)
			}
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*float64); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			if *x == math.MaxFloat64 {
				return "---"
			} else {
				return strconv.FormatFloat(*x, 'f', dp, 64)
			}
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(float32); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			if x == math.MaxFloat32 {
				return "---"
			} else {
				return strconv.FormatFloat(float64(x), 'f', dp, 32)
			}
		}
		return StringerValue(th, s, options...)
	}
	if x, ok := any(v).(*float32); ok {
		s := func(dp int) string {
			GuiLock.RLock()
			defer GuiLock.RUnlock()
			if *x == math.MaxFloat32 {
				return "---"
			} else {
				return strconv.FormatFloat(float64(*x), 'f', dp, 32)
			}
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
