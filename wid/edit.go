// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
)

type EditDef struct {
	Editor
	th        *Theme
	shaper    text.Shaper
	alignment layout.Alignment
	CharLimit uint
	font      text.Font
	hint      string
}

var prev Focuser

func Edit(th *Theme, hint string) func(gtx C) D {
	e := new(EditDef)
	e.th = th
	e.shaper = th.Shaper
	e.hint = hint
	e.SingleLine = true
	e.SetupTabs()
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
			paint.ColorOp{Color: e.th.Palette.OnBackground}.Add(gtx.Ops)
			e.Editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if !disabled {
			paint.ColorOp{Color: e.th.Palette.OnBackground}.Add(gtx.Ops)
			e.Editor.PaintCaret(gtx)
		}
		return dims
	}
}

func (e *EditDef) Layout(gtx C) D {
	defer op.Save(gtx.Ops).Load()
	min := gtx.Constraints.Min
	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Stacked(func(gtx C) D {
					gtx.Constraints.Min = min
					return e.th.LabelInset.Layout(gtx, func(gtx C) D {
						return e.LayoutEdit()(gtx)
					})
				}),
				layout.Expanded(func(gtx C) D {
					outline := f32.Rectangle{Max: f32.Point{
						X: float32(gtx.Constraints.Min.X),
						Y: float32(gtx.Constraints.Min.Y),
					}}
					if e.Focused() {
						PaintBorder(gtx, outline, MulAlpha(e.th.Palette.Primary, 255), e.th.BorderThicknessActive, e.th.CornerRadius)
					} else if e.Hovered() {
						PaintBorder(gtx, outline, MulAlpha(e.th.Palette.Primary, 100), e.th.BorderThickness, e.th.CornerRadius)
					} else {
						PaintBorder(gtx, outline, MulAlpha(e.th.Palette.Primary, 50), e.th.BorderThickness, e.th.CornerRadius)
					}
					return D{}
				}),
			)
		}),
	)
	return dims
}
