// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"

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
	Base
	widget.Editor
	hovered         bool
	outlineColor    color.NRGBA
	selectionColor  color.NRGBA
	CharLimit       uint
	label           string
	value           *string
	LabelSize       unit.Sp
	borderThickness unit.Dp
	wasFocused      bool
}

// Edit will return a widget (layout function) for a text editor
func Edit(th *Theme, options ...Option) func(gtx C) D {
	e := new(EditDef)
	e.th = th
	e.Font = &th.DefaultFont
	e.LabelSize = th.TextSize * 8
	e.SingleLine = true
	e.borderThickness = th.BorderThickness
	e.width = unit.Dp(5000)    // Default to max width that is possible
	e.padding = layout.Inset{} // th.EditPadding
	e.outlineColor = th.Fg(Outline)
	e.selectionColor = MulAlpha(th.Bg(Primary), 60)
	e.role = Canvas
	// Read in options to change from default values to something else.
	for _, option := range options {
		option.apply(e)
	}
	if e.borderThickness > 0 {
		e.padding = th.EditPadding
	}
	if e.value != nil {
		e.Editor.SetText(*e.value)
	}
	return func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Constraints.Min.X
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

// Var is an option parameter to set the variable uptdated
func Var(s *string) EditOption {
	return func(w *EditDef) {
		w.value = s
	}
}

func (e *EditDef) setBorder(w unit.Dp) {
	e.borderThickness = w
}

func (e EditOption) apply(cfg interface{}) {
	if o, ok := cfg.(*EditDef); ok {
		e(o)
	}
}

func (e *EditDef) setLabel(s string) {
	e.label = s
}

func rr(gtx C, radius unit.Dp, height int) int {
	rr := gtx.Dp(radius)
	if rr > (height-1)/2 {
		return (height - 1) / 2
	}
	return rr
}

func (e *EditDef) layoutEditBackground() func(gtx C) D {
	return func(gtx C) D {
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Max.X,
			Y: gtx.Constraints.Max.Y,
		}}
		rr := rr(gtx, e.th.BorderCornerRadius, outline.Max.Y)
		if e.Focused() {
			paint.FillShape(gtx.Ops, e.Bg(), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		}
		return D{}
	}
}

func paintBorder(gtx C, outline image.Rectangle, col color.NRGBA, width unit.Dp, rr int) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, rr).Path(gtx.Ops),
			Width: Pxr(gtx, width),
		}.Op(),
	)
}

// LayoutBorder will draw a border around the widget
func LayoutBorder(e *EditDef, th *Theme) func(gtx C) D {
	return func(gtx C) D {
		if !e.Focused() && e.value != nil {
			current := e.Text()
			if e.wasFocused {
				// When the edit is loosing focus, we must update the underlying variable
				GuiLock.Lock()
				*e.value = current
				GuiLock.Unlock()
			} else {
				// When the underlying variable changes, update the edit buffer
				GuiLock.RLock()
				s := *e.value
				GuiLock.RUnlock()
				if s != current {
					e.SetText(s)
				}
			}
		}

		e.wasFocused = e.Focused()
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X,
			Y: gtx.Constraints.Min.Y,
		}}
		r := gtx.Dp(th.BorderCornerRadius)
		if r > outline.Max.Y/2 {
			r = outline.Max.Y / 2
		}
		if e.borderThickness > 0 {
			if e.Focused() {
				paintBorder(gtx, outline, e.outlineColor, th.BorderThicknessActive, r)
			} else if e.hovered {
				paintBorder(gtx, outline, e.outlineColor, (th.BorderThickness+th.BorderThicknessActive)/2, r)
			} else {
				paintBorder(gtx, outline, e.outlineColor, th.BorderThickness, r)
			}
		}
		eventArea := clip.Rect(outline).Push(gtx.Ops)
		for _, ev := range gtx.Events(&e.hovered) {
			if ev, ok := ev.(pointer.Event); ok {
				switch ev.Type {
				case pointer.Leave:
					e.hovered = false
				case pointer.Enter:
					e.hovered = true
				}
			}
		}
		pointer.InputOp{
			Tag:   &e.hovered,
			Types: pointer.Enter | pointer.Leave,
		}.Add(gtx.Ops)
		eventArea.Pop()

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
			paint.ColorOp{Color: e.th.Fg(Canvas)}.Add(gtx.Ops)
			w := widget.Label{Alignment: text.End}.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize, e.label)
			return w
		})
	}
}

func (e *EditDef) layoutEdit() func(gtx C) D {
	return func(gtx C) D {
		macro := op.Record(gtx.Ops)
		paint.ColorOp{Color: MulAlpha(e.Fg(), 110)}.Add(gtx.Ops)
		var maxLines int
		if e.Editor.SingleLine {
			maxLines = 1
		}
		tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: maxLines}
		dims := tl.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize, e.hint)
		call := macro.Stop()
		if w := dims.Size.X; gtx.Constraints.Min.X < w {
			gtx.Constraints.Min.X = w
		}
		if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
			gtx.Constraints.Min.Y = h
		}
		dims = e.Editor.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize, nil)
		disabled := gtx.Queue == nil
		if e.Editor.Len() > 0 {
			paint.ColorOp{Color: e.selectionColor}.Add(gtx.Ops)
			e.Editor.PaintSelection(gtx)
			paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
			e.Editor.PaintText(gtx)
		} else {
			call.Add(gtx.Ops)
		}
		if !disabled && e.Editor.Len() > 0 {
			paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
			e.Editor.PaintCaret(gtx)
		}
		return dims
	}
}
