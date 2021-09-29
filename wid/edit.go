// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type EditDef struct {
	Editor
	Widget
	th        *Theme
	shaper    text.Shaper
	alignment layout.Alignment
	CharLimit uint
	font      text.Font
	hint      string
	padding   layout.Inset
	width     unit.Value
}

type EditOption func(*EditDef)

func (b *EditDef) ApplyOptions(options []Option) {
	//for _, option := range options {
		//option.Do(b)
	//}
}

var prev Focuser

func Edit(th *Theme, hint string, options ...Option) func(gtx C) D {
	e := new(EditDef)
	e.th = th
	e.shaper = th.Shaper
	e.hint = hint
	e.SingleLine = true
	e.SetupTabs()
	e.padding = layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(5), Right: unit.Dp(1)}
	e.ApplyOptions(options)
	return func(gtx C) D {
		return e.Layout(gtx)
	}
}

func (e *EditDef) LayoutEdit() func(gtx C) D {
	return func(gtx C) D {
		defer op.Save(gtx.Ops).Load()
		macro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.th.HintColor}.Add(gtx.Ops)
		var maxlines int
		if e.Editor.SingleLine {
			maxlines = 1
		}
		tl := aLabel{Alignment: e.Editor.Alignment, MaxLines: maxlines}
		dims := tl.Layout(gtx, e.shaper, e.font, e.th.TextSize, e.hint)
		call := macro.Stop()
		if w := dims.Size.X; gtx.Constraints.Min.X < w {
			gtx.Constraints.Min.X = w
		}
		if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
			gtx.Constraints.Min.Y = h
		}
		dims = e.Editor.Layout(gtx, e.shaper, e.font, e.th.TextSize)
		disabled := gtx.Queue == nil || GlobalDisable
		if e.Editor.Len() > 0 {
			paint.ColorOp{Color: e.th.SelectionColor}.Add(gtx.Ops)
			e.Editor.PaintSelection(gtx)
			paint.ColorOp{Color: e.th.OnBackground}.Add(gtx.Ops)
			e.Editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if !disabled {
			paint.ColorOp{Color: e.th.OnBackground}.Add(gtx.Ops)
			e.Editor.PaintCaret(gtx)
		}
		return dims
	}
}

func (e *EditDef) LayoutBorder() func(gtx C) D {
	return func(gtx C) D {
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		if e.Focused() {
			PaintBorder(gtx, outline, MulAlpha(e.th.Primary, 255), e.th.BorderThicknessActive, e.th.CornerRadius)
		} else if e.Hovered() {
			PaintBorder(gtx, outline, MulAlpha(e.th.Primary, 100), e.th.BorderThickness, e.th.CornerRadius)
		} else {
			PaintBorder(gtx, outline, MulAlpha(e.th.Primary, 50), e.th.BorderThickness, e.th.CornerRadius)
		}
		return D{}
	}
}

func (e *EditDef) Layout(gtx C) D {
	defer op.Save(gtx.Ops).Load()
	min := gtx.Constraints.Min
	return e.padding.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(
			gtx,
			layout.Expanded(func(gtx C) D {
				gtx.Constraints.Min = min
				return e.th.LabelInset.Layout(gtx, func(gtx C) D {
					return e.LayoutEdit()(gtx)
				})
			}),
			layout.Expanded(e.LayoutBorder()),
		)
	})
}
