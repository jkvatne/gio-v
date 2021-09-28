// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gio-v/f32color"
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
	shaper text.Shaper
	// Alignment specifies where to anchor the text.
	Alignment layout.Alignment
	// Helper text to give additional context to a field.
	Helper string
	// CharLimit specifies the maximum number of characters the text input
	// will allow. Zero means "no limit".
	CharLimit uint
	border border
	Font     text.Font
	TextSize unit.Value
	// Color is the text color.
	Color color.NRGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.NRGBA
	// SelectionColor is the color of the background for selected text.
	SelectionColor color.NRGBA
}

type border struct {
	Thickness unit.Value
	Color     color.NRGBA
}

var prev Focuser

func TextField(th *Theme, hint string) func(gtx C) D {
	c := new(Edit)
	c.shaper = th.Shaper
	c.TextSize = th.TextSize
	c.Color = th.Palette.OnBackground
	c.HintColor = f32color.MulAlpha(th.Palette.OnBackground, 0xbb)
	c.SelectionColor = f32color.MulAlpha(th.Palette.Primary, 0x60)
	c.Hint = hint
	if prev != nil {
		c.SetPrev(prev)
		prev.SetNext(c)
	}
	prev = c
	return func(gtx C) D {
		return c.Layout(gtx, th, hint)
	}
}

func SetupBorder(in Clickable, th *Theme, disabled bool) (thick unit.Value, color color.NRGBA) {
	switch {
	case in.Focused() && !disabled:
		return th.BorderThicknessActive, th.BorderColorActive
	case in.Hovered() && !disabled:
		return th.BorderThickness, th.BorderColorHovered
	}
	return th.BorderThickness , th.BorderColor
}

func drawBorder(gtx C, in *Edit, th *Theme) (op.CallOp, D){
	macro := op.Record(gtx.Ops)
	dims := BorderDef{
		Color:        th.BorderColor,
		Width:        in.border.Thickness,
		CornerRadius: th.CornerRadius,
	}.Layout(
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

func getBorder(in *Edit, th *Theme) layout.StackChild {
	return layout.Expanded(func(gtx C) D {
		border, dims := drawBorder(gtx, in, th)
		border.Add(gtx.Ops)
		return dims
	})
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return f32color.Disabled(c)
	}
	return c
}

func (e *Edit) xLayout(gtx layout.Context) layout.Dimensions {
	defer op.Save(gtx.Ops).Load()
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	var maxlines int
	if e.Editor.SingleLine {
		maxlines = 1
	}
	tl := aLabel{Alignment: e.Editor.Alignment, MaxLines: maxlines}
	dims := tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	call := macro.Stop()
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = e.Editor.Layout(gtx, e.shaper, e.Font, e.TextSize)
	disabled := gtx.Queue == nil
	if e.Editor.Len() > 0 {
		paint.ColorOp{Color: blendDisabledColor(disabled, e.SelectionColor)}.Add(gtx.Ops)
		e.Editor.PaintSelection(gtx)
		paint.ColorOp{Color: blendDisabledColor(disabled, e.Color)}.Add(gtx.Ops)
		e.Editor.PaintText(gtx)
	} else {
		call.Add(gtx.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		e.Editor.PaintCaret(gtx)
	}
	return dims
}

func innerfunc(in *Edit, th *Theme, hint string) layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
			Spacing: func() layout.Spacing {
				switch in.Alignment {
				case layout.Middle:
					return layout.SpaceSides
				case layout.End:
					return layout.SpaceStart
				default: // layout.Start and all others
					return layout.SpaceEnd
				}
			}(),
		}.Layout(
			gtx,
			layout.Rigid(func(gtx C) D {
				return in.xLayout(gtx) //Editor(th, &in.Editor, hint).Layout(gtx)
			}),
		)
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


func (in *Edit) Layout(gtx C, th *Theme, hint string) D {
	in.border.Thickness, in.border.Color = SetupBorder(in.Clickable, th, gtx.Queue == nil)
	defer op.Save(gtx.Ops).Load()
	dims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(
		gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Stack{}.Layout(
				gtx,
				getBorder(in, th),
				layout.Stacked(func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, innerfunc(in, th, hint))
				}),
				layout.Expanded(func(gtx C) D {
					HandleMouseClick(gtx, in)
					HandleMouseHover(gtx, in)
					DeclareInputHandler(gtx, in)
					return D{Size: gtx.Constraints.Min}
				}),
			)
		}),
	)
	return dims
}
