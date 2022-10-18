// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// EditDef is the parameters for the text editor
type EditDef struct {
	Widget
	widget.Editor
	shaper    text.Shaper
	alignment layout.Alignment
	CharLimit uint
	font      text.Font
	label     string
	value     *string
	LabelSize unit.Sp
}

// Edit will return a widget (layout function) for a text editor
func Edit(th *Theme, options ...Option) func(gtx C) D {
	e := new(EditDef)
	e.th = th
	e.shaper = th.Shaper
	e.LabelSize = th.TextSize * 6
	e.SingleLine = true
	e.width = unit.Dp(5000) // Default to max width that is possible
	e.padding = th.EditPadding
	// Read in options to change from default values to something else.
	for _, option := range options {
		option.apply(e)
	}
	return func(gtx C) D {
		gtx.Constraints.Min.X = 0
		if e.label == "" {
			return e.layEdit()(gtx)
		}
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Start, Spacing: layout.SpaceStart}.Layout(
			gtx,
			layout.Rigid(e.layLabel()),
			layout.Rigid(e.layEdit()),
		)
	}
}

// EditOption is options specific to Edits
type EditOption func(w *EditDef)

// Lbl is an option parameter to set the widget label
func Lbl(s string) EditOption {
	return func(w *EditDef) {
		w.setLabel(s)
	}
}

func Var(s *string) EditOption {
	return func(w *EditDef) {
		w.value = s
	}
}

func (e EditOption) apply(cfg interface{}) {
	e(cfg.(*EditDef))
}

func (e *EditDef) setLabel(s string) {
	e.label = s
}

func (e *EditDef) layoutEditBackground() func(gtx C) D {
	return func(gtx C) D {
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X,
			Y: gtx.Constraints.Min.Y,
		}}
		rr := gtx.Dp(e.th.CornerRadius)
		color := e.th.Surface
		if e.Focused() {
			color = e.th.Background
		}
		paint.FillShape(gtx.Ops, color, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		return D{}
	}
}

func paintBorder(gtx C, outline image.Rectangle, col color.NRGBA, width unit.Dp, rr unit.Dp) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, gtx.Dp(rr)).Path(gtx.Ops),
			Width: Pxr(gtx, width),
		}.Op(),
	)
}

// LayoutBorder will draw a border around the widget
func LayoutBorder(e *EditDef, th *Theme) func(gtx C) D {
	return func(gtx C) D {
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X,
			Y: gtx.Constraints.Min.Y,
		}}
		if e.Focused() {
			paintBorder(gtx, outline, MulAlpha(th.Primary, 255), th.BorderThicknessActive, th.CornerRadius)
			// TODO } else if e.Hovered() {
			//	paintBorder(gtx, outline, MulAlpha(th.Primary, 140), th.BorderThickness, th.CornerRadius)
		} else {
			paintBorder(gtx, outline, MulAlpha(th.Primary, 50), th.BorderThickness, th.CornerRadius)
		}
		return D{}
	}
}

func (e *EditDef) layEdit() layout.Widget {
	return func(gtx C) D {
		return e.padding.Layout(gtx, func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(e.layoutEditBackground()),
				layout.Expanded(func(gtx C) D {
					gtx.Constraints.Min.X = 5000
					return e.th.LabelPadding.Layout(gtx, func(gtx C) D {
						return e.layoutEdit()(gtx)
					})
				}),
				layout.Expanded(LayoutBorder(e, e.th)),
			)
		})
	}
}

func (e *EditDef) layLabel() layout.Widget {
	return func(gtx C) D {
		p := e.padding
		p.Top = p.Top + e.th.LabelPadding.Top
		return p.Layout(gtx, func(gtx C) D {
			if e.label == "" {
				return D{}
			}
			gtx.Constraints.Min.X = gtx.Sp(e.LabelSize)
			paint.ColorOp{Color: e.th.OnBackground}.Add(gtx.Ops)
			w := widget.Label{Alignment: text.End}.Layout(gtx, e.shaper, e.font, e.th.TextSize, e.label)
			return w
		})
	}
}

func (e *EditDef) layoutEdit() func(gtx C) D {
	return func(gtx C) D {
		macro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.th.HintColor}.Add(gtx.Ops)
		var maxLines int
		if e.Editor.SingleLine {
			maxLines = 1
		}
		tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: maxLines}
		dims := tl.Layout(gtx, e.shaper, e.font, e.th.TextSize, e.hint)
		call := macro.Stop()
		if w := dims.Size.X; gtx.Constraints.Min.X < w {
			gtx.Constraints.Min.X = w
		}
		if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
			gtx.Constraints.Min.Y = h
		}
		dims = e.Editor.Layout(gtx, e.shaper, e.font, e.th.TextSize, nil)
		disabled := gtx.Queue == nil
		if e.Editor.Len() > 0 {
			paint.ColorOp{Color: e.th.SelectionColor}.Add(gtx.Ops)
			e.Editor.PaintSelection(gtx)
			paint.ColorOp{Color: e.th.OnBackground}.Add(gtx.Ops)
			e.Editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if !disabled && e.Editor.Len() > 0 {
			paint.ColorOp{Color: e.th.OnBackground}.Add(gtx.Ops)
			e.Editor.PaintCaret(gtx)
		}
		return dims
	}
}
