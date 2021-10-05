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
	Widget
	Editor
	shaper    text.Shaper
	alignment layout.Alignment
	CharLimit uint
	font      text.Font
	LabelSize unit.Value
}


func Edit(th *Theme, options ...Option) func(gtx C) D {
	e := new(EditDef)
	e.SetupTabs()
	// Set up default values
	e.th = th
	e.shaper = th.Shaper
	e.LabelSize = unit.Dp(150)
	e.MaxLines = 1
	e.width = unit.Dp(5000) // Default to max width that is possible
	e.padding = layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(2), Left: unit.Dp(5), Right: unit.Dp(1)}
	// Read in options to change from default values to something else.
	for _, option := range options {
		option.apply(e)
	}
	return func(gtx C) D {
			defer op.Save(gtx.Ops).Load()
			gtx.Constraints.Min.X = 0
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start, Spacing: layout.SpaceStart}.Layout(
				gtx,
				layout.Rigid(e.layLabel()),
				layout.Rigid(e.layEdit()),
			)
	}
}

func (e *EditDef)layEdit() layout.Widget{
	return func(gtx C) D {
		return e.padding.Layout(gtx, func (gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(func(gtx C) D {
					gtx.Constraints.Min.X = 5000
					return e.th.LabelPadding.Layout(gtx, func(gtx C) D {
						return e.LayoutEdit()(gtx)
					})
				}),
				layout.Expanded(e.LayoutBorder()),
			)})
	}
}

func (e *EditDef)layLabel() layout.Widget {
	return func(gtx C) D {
		p := e.padding
		p.Top = unit.Dp(p.Top.V + e.th.LabelPadding.Top.V)
		return p.Layout(gtx, func(gtx C) D {
			if e.hint=="" {
				return D{}
			}
			gtx.Constraints.Min.X = gtx.Metric.Px(e.LabelSize)
			paint.ColorOp{Color: e.th.OnBackground}.Add(gtx.Ops)
			return aLabel{Alignment: text.End}.Layout(gtx, e.shaper, e.font, e.th.TextSize, e.hint)
		})
	}
}

func (e *EditDef) LayoutEdit() func(gtx C) D {
	return func(gtx C) D {
		defer op.Save(gtx.Ops).Load()
		macro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.th.HintColor}.Add(gtx.Ops)
		var maxlines int
		if e.Editor.MaxLines <= 1 {
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
