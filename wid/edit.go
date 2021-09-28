// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"image"
	"image/color"
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

func TextField(th *Theme, hint string) func(gtx C) D {
	c := new(EditDef)
	c.th = th
	c.shaper = th.Shaper
	c.hint = hint
	c.SetupTabs()
	return func(gtx C) D {
		return c.Layout(gtx)
	}
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return Disabled(c)
	}
	return c
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

func HandleMouseHover(gtx C, in *EditDef) {
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

func HandleMouseClick(gtx C, in *EditDef) {
	// Set pass-through mode so the underlying editor will recieve clicks
	stack := op.Save(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	// Handle clickable event handler
	in.Clickable.LayoutClickable(gtx)
	stack.Load()
}

func DeclareInputHandler(gtx C, in *EditDef) {
	stack := op.Save(gtx.Ops)
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	pointer.InputOp{
		Tag:   in,
		Types: pointer.Enter | pointer.Leave | pointer.Cancel,
	}.Add(gtx.Ops)
	stack.Load()
}

func (e *EditDef) LayoutBackground() func(gtx C) D {
	return func(gtx C) D {
		rr := Pxr(gtx, e.th.CornerRadius)
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		clip.UniformRRect(outline, rr).Add(gtx.Ops)
		switch {
		case e.Hovered() || e.Focused():
			paint.FillShape(gtx.Ops, Hovered(e.th.Palette.Background), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			PaintBorder(gtx, outline, Disabled(e.th.Palette.Primary), e.th.BorderThickness, e.th.CornerRadius)
		default:
			PaintBorder(gtx, outline, Disabled(e.th.Palette.Primary), e.th.BorderThickness, e.th.CornerRadius)
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
}

func (e *EditDef) Layout(gtx C) D {
	defer op.Save(gtx.Ops).Load()
	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				layout.Expanded(e.LayoutBackground()),
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
