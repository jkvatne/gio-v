// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"image"
	"image/color"
)

type Edit struct {
	// Editor contains the edit buffer.
	Editor
	th     *Theme
	shaper text.Shaper
	// Alignment specifies where to anchor the text.
	Alignment layout.Alignment
	// Helper text to give additional context to a field.
	//Helper string
	// CharLimit specifies the maximum number of characters the text input
	// will allow. Zero means "no limit".
	CharLimit uint
	border    border
	Font      text.Font
	// Hint contains the text displayed when the editor is empty.
	hint string
}

type border struct {
	Thickness unit.Value
	Color     color.NRGBA
}

var prev Focuser

func TextField(th *Theme, hint string) func(gtx C) D {
	c := new(Edit)
	c.th = th
	c.shaper = th.Shaper
	c.hint = hint
	c.SetupTabs()
	return func(gtx C) D {
		return c.Layout(gtx)
	}
}

func DrawBorder(gtx C, e *Edit) (op.CallOp, D) {
	macro := op.Record(gtx.Ops)
	w := e.th.BorderThickness
	c := e.th.BorderColor
	switch {
	case e.Focused() && !e.disabled:
		w = e.th.BorderThicknessActive
		c = e.th.BorderColorActive
	case e.Hovered() && !e.disabled:
		w = e.th.BorderThickness
		c = e.th.BorderColorHovered
	}
	dims := BorderDef{Color: c, Width: w, CornerRadius: e.th.CornerRadius}.Layout(
		gtx,
		func(gtx C) D {
			return D{Size: image.Point{
				X: gtx.Constraints.Max.X,
				Y: gtx.Constraints.Min.Y,
			}}
		},
	)
	return macro.Stop(), dims
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return Disabled(c)
	}
	return c
}

func (e *Edit) LayoutEdit() func(gtx C) D {
	return func(gtx C) D {
		defer op.Save(gtx.Ops).Load()
		macro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.th.HintColor}.Add(gtx.Ops)
		var maxlines int
		if e.Editor.SingleLine {
			maxlines = 1
		}
		tl := aLabel{Alignment: e.Editor.Alignment, MaxLines: maxlines}
		dims := tl.Layout(gtx, e.shaper, e.Font, e.th.TextSize, e.hint)
		call := macro.Stop()
		if w := dims.Size.X; gtx.Constraints.Min.X < w {
			gtx.Constraints.Min.X = w
		}
		if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
			gtx.Constraints.Min.Y = h
		}
		dims = e.Editor.Layout(gtx, e.shaper, e.Font, e.th.TextSize)
		disabled := gtx.Queue == nil
		if e.Editor.Len() > 0 {
			paint.ColorOp{Color: blendDisabledColor(disabled, e.th.SelectionColor)}.Add(gtx.Ops)
			e.Editor.PaintSelection(gtx)
			paint.ColorOp{Color: blendDisabledColor(disabled, e.th.Palette.OnBackground)}.Add(gtx.Ops)
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

func HandleMouseHover(gtx C, in *Edit) {
	for _, event := range gtx.Events(in) {
		if event, ok := event.(pointer.Event); ok {
			switch event.Type {
			case pointer.Enter:
				in.SetHovered(true)
			case pointer.Leave, pointer.Cancel:
				in.SetHovered(false)
			}
		}
	}
}

func HandleMouseClick(gtx C, in *Edit) {
	// Set pass-through mode so the underlying editor will recieve clicks
	stack := op.Save(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	// Handle clickable event handler
	in.Clickable.LayoutClickable(gtx)
	stack.Load()
}

func DeclareInputHandler(gtx C, in *Edit) {
	stack := op.Save(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   in,
		Types: pointer.Enter | pointer.Leave | pointer.Cancel,
	}.Add(gtx.Ops)
	stack.Load()
}

func (e *Edit) Layout(gtx C) D {
	//e.border.Thickness, e.border.Color = SetupBorder(e.Clickable, e.th, gtx.Queue == nil)
	defer op.Save(gtx.Ops).Load()
	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(func(gtx C) D {
					border, dims := DrawBorder(gtx, e)
					border.Add(gtx.Ops)
					return dims
				}),
				layout.Stacked(func(gtx C) D {
					return e.th.LabelInset.Layout(gtx, func(gtx C) D {
						return e.LayoutEdit()(gtx)
					})
				}),
				layout.Expanded(func(gtx C) D {
					HandleMouseClick(gtx, e)
					HandleMouseHover(gtx, e)
					DeclareInputHandler(gtx, e)
					return D{Size: gtx.Constraints.Min}
				}),
			)
		}),
	)
	return dims
}
