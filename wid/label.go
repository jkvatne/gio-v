// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/font"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"image"
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
	value    interface{}
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

func (e LabelOption) apply(cfg interface{}) {
	e(cfg.(*LabelDef))
}

// StringerValue returns a widget for a value given by stringer function
func Label[V Value](th *Theme, v V, options ...Option) func(gtx C) D {
	w := LabelDef{
		Base: Base{
			th:        th,
			role:      Surface,
			padding:   uniformPadding(3.5),
			FontScale: 1.0,
			Alignment: text.Start,
		},
		Font:     th.DefaultFont,
		MaxLines: 0,
		value:    v,
	}
	// Apply options after initialization of LabelDef
	for _, option := range options {
		option.apply(&w)
	}

	return func(gtx C) D {
		c := gtx
		if w.MaxLines == 1 {
			// This is a hack to avoid splitting the line when only one line is allowed
			c.Constraints.Max.X = inf
		}
		GuiLock.RLock()
		var str string
		if w.DpNo != nil {
			str = ValueToString(v, *w.DpNo)
		} else {
			str = ValueToString(v, 0)
		}
		GuiLock.RUnlock()
		o := op.Offset(image.Pt(Px(gtx, w.padding.Left), Px(gtx, w.padding.Top))).Push(gtx.Ops)
		tl := widget.Label{Alignment: w.Alignment, MaxLines: w.MaxLines}
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: w.Fg()}.Add(gtx.Ops)
		c.Constraints.Min.X -= Px(gtx, w.padding.Left+w.padding.Right)
		c.Constraints.Max.X -= Px(gtx, w.padding.Left+w.padding.Right)
		c.Constraints.Min.Y = Max(0, c.Constraints.Min.Y-Px(gtx, w.padding.Top+w.padding.Bottom))
		c.Constraints.Max.Y -= Px(gtx, w.padding.Top+w.padding.Bottom)
		dims := tl.Layout(c, w.th.Shaper, w.Font, unit.Sp(w.FontScale)*w.th.FontSp(), str, colMacro.Stop())
		o.Pop()
		dims.Size.X += Px(gtx, w.padding.Left+w.padding.Right)
		dims.Size.Y += Px(gtx, w.padding.Bottom+w.padding.Top)
		return dims
	}
}
