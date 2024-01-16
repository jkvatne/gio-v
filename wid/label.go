// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
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
func Label[V Value](th *Theme, v V, options ...Option) layout.Widget {
	w := LabelDef{
		Base: Base{
			th:        th,
			role:      Surface,
			padding:   th.DefaultPadding,
			margin:    layout.Inset{-1, -1, -1, -1},
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
	if w.margin.Top != -1 {
		panic("Label does not use margin")
	}
	return w.Layout
}

func (w *LabelDef) Layout(gtx C) D {
	c := gtx
	if w.MaxLines == 1 {
		// This is a hack to avoid splitting the line when only one line is allowed
		c.Constraints.Max.X = inf
	}
	GuiLock.RLock()
	var str string
	if w.DpNo != nil {
		str = ValueToString(w.value, *w.DpNo)
	} else {
		str = ValueToString(w.value, 0)
	}
	GuiLock.RUnlock()
	defer op.Offset(image.Pt(Px(gtx, w.padding.Left), Px(gtx, w.padding.Top))).Push(gtx.Ops).Pop()
	tl := widget.Label{Alignment: w.Alignment, MaxLines: w.MaxLines}
	c.Constraints.Min.X = Max(c.Constraints.Min.X-Px(gtx, w.padding.Left+w.padding.Right), 0)
	c.Constraints.Max.X -= Px(gtx, w.padding.Left+w.padding.Right)
	c.Constraints.Min.Y = Max(0, c.Constraints.Min.Y-Px(gtx, w.padding.Top+w.padding.Bottom))
	// Fill background if bgColor is given
	if w.bgColor != nil && (*w.bgColor).A != 0 {
		paint.FillShape(gtx.Ops, *w.bgColor, clip.UniformRRect(image.Rectangle{Max: c.Constraints.Max}, 0).Op(gtx.Ops))
	}
	// Macro for the text drawing color
	colMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: w.Fg()}.Add(gtx.Ops)
	// Then lay out the text
	dims := tl.Layout(c, w.th.Shaper, w.Font, unit.Sp(w.FontScale)*w.th.TextSize, str, colMacro.Stop())
	dims.Size.X += Px(gtx, w.padding.Left+w.padding.Right)
	dims.Size.Y += Px(gtx, w.padding.Bottom+w.padding.Top)
	return dims
}
