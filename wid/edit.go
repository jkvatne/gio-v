// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"fmt"
	"gioui.org/text"
	"image"
	"image/color"
	"math"
	"strconv"

	"gioui.org/io/pointer"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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
	value           interface{}
	labelSize       float32
	borderThickness unit.Dp
	wasFocused      bool
}

func FloatToStr(x float64, dp int) string {
	if dp == 1 {
		return fmt.Sprintf("%0.1f", x)
	} else if dp == 2 {
		return fmt.Sprintf("%0.2f", x)
	} else if dp == 3 {
		return fmt.Sprintf("%0.3f", x)
	} else if dp == 4 {
		return fmt.Sprintf("%0.4f", x)
	} else if dp == 5 {
		return fmt.Sprintf("%0.3f", x)
	} else if dp == 6 {
		return fmt.Sprintf("%0.6f", x)
	} else if dp == 7 {
		return fmt.Sprintf("%0.7f", x)
	} else {
		return fmt.Sprintf("%0.0f", x)
	}

}

func ValueToString(v interface{}, dp int) string {
	if v == nil {
		return "nil"
	} else if x, ok := v.(*int); ok {
		if *x == math.MinInt {
			return "---"
		} else {
			return fmt.Sprintf("%d", *x)
		}
	} else if x, ok := v.(*float32); ok {
		if *x == math.MaxFloat32 {
			return "---"
		} else {
			return FloatToStr(float64(*x), dp)
		}
	} else if x, ok := v.(*float64); ok {
		if *x == math.MaxFloat64 {
			return "---"
		} else {
			return FloatToStr(*x, dp)
		}
	} else if x, ok := v.(*string); ok {
		return *x
	}
	return ""
}

func StringToValue(value interface{}, current string) {
	if _, ok := value.(*int); ok {
		x, err := strconv.Atoi(current)
		if err == nil {
			*value.(*int) = x
		}
	} else if _, ok := value.(float32); ok {
		f, err := strconv.ParseFloat(current, 32)
		if err == nil {
			*value.(*float32) = float32(f)
		}
	} else if _, ok := value.(float64); ok {
		f, err := strconv.ParseFloat(current, 64)
		if err == nil {
			*value.(*float64) = f
		}
	} else if _, ok := value.(*string); ok {
		*value.(*string) = current
	}
}

// Edit will return a widget (layout function) for a text editor
func Edit(th *Theme, options ...Option) func(gtx C) D {
	e := new(EditDef)
	e.th = th
	e.Font = &th.DefaultFont
	e.labelSize = th.LabelSplit // 1/3 of column width
	e.SingleLine = true
	e.borderThickness = th.BorderThickness
	e.width = unit.Dp(5000) // Default to max width that is possible
	e.padding = th.OutsidePadding
	e.outlineColor = th.Fg(Outline)
	e.selectionColor = MulAlpha(th.Bg(Primary), 60)
	e.value = ""
	// Read in options to change from default values to something else.
	for _, option := range options {
		option.apply(e)
	}
	if e.value != nil {
		e.Editor.SetText(ValueToString(e.value, e.Dp))
	}
	return func(gtx C) D {
		return e.Layout(gtx)
	}
}

func (e *EditDef) updateValue() {
	if !e.Focused() && e.value != nil {
		current := e.Text()
		if e.wasFocused {
			// When the edit is loosing focus, we must update the underlying variable
			GuiLock.Lock()
			StringToValue(e.value, current)
			GuiLock.Unlock()
		} else {
			// When the underlying variable changes, update the edit buffer
			GuiLock.RLock()
			s := e.value
			GuiLock.RUnlock()
			if s != current {
				e.SetText(ValueToString(e.value, e.Dp))
			}
		}
	}
	e.wasFocused = e.Focused()
}

func (e *EditDef) maxLines() int {
	if e.Editor.SingleLine {
		return 1
	}
	return 0
}

func (e *EditDef) Layout(gtx C) D {
	e.CheckDisable(gtx)
	// Choose colors.
	textColorMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
	textColor := textColorMacro.Stop()
	hintColorMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: MulAlpha(e.Fg(), 128)}.Add(gtx.Ops)
	hintColor := hintColorMacro.Stop()
	selectionColorMacro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.th.SelectionColor}.Add(gtx.Ops)
	selectionColor := selectionColorMacro.Stop()

	e.updateValue()

	// Move to offset the outside padding
	defer op.Offset(image.Pt(
		Px(gtx, e.padding.Left),
		Px(gtx, e.padding.Top))).Push(gtx.Ops).Pop()

	// And reduce the size to make space for the padding
	gtx.Constraints.Min.X -= Px(gtx, e.padding.Left+e.padding.Right+e.th.InsidePadding.Left+e.th.InsidePadding.Right)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X

	// Draw hint text with top/left padding offset
	macro := op.Record(gtx.Ops)
	o := op.Offset(image.Pt(Px(gtx, e.th.InsidePadding.Left), Px(gtx, e.th.InsidePadding.Top))).Push(gtx.Ops)
	paint.ColorOp{Color: MulAlpha(e.Fg(), 110)}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: e.maxLines()}
	LblDim := tl.Layout(gtx, e.th.Shaper, *e.Font, e.th.FontSp(), e.hint, hintColor)
	o.Pop()
	callHint := macro.Stop()

	// Add outside label to the left of the dropdown box
	if e.label != "" {
		o := op.Offset(image.Pt(0, Px(gtx, e.th.InsidePadding.Top))).Push(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		oldMaxX := gtx.Constraints.Max.X
		ofs := int(float32(oldMaxX) * e.labelSize)
		gtx.Constraints.Max.X = ofs - Px(gtx, e.th.InsidePadding.Left)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		ll := widget.Label{Alignment: text.End, MaxLines: 1}
		ll.Layout(gtx, e.th.Shaper, *e.Font, e.th.FontSp(), e.label, colMacro.Stop())
		o.Pop()
		gtx.Constraints.Max.X = oldMaxX - ofs
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()
	}
	// If a width is given, and it is within constraints, limit size
	if w := Px(gtx, e.width); w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
	}

	border := image.Rectangle{Max: image.Pt(
		gtx.Constraints.Max.X+Px(gtx, e.th.InsidePadding.Left+e.th.InsidePadding.Right),
		LblDim.Size.Y+Px(gtx, e.th.InsidePadding.Bottom+e.th.InsidePadding.Top))}

	r := Px(gtx, e.th.BorderCornerRadius)
	if r > border.Max.Y/2 {
		r = border.Max.Y / 2
	}
	if e.Focused() {
		paint.FillShape(gtx.Ops, e.th.Bg(Canvas), clip.UniformRRect(border, r).Op(gtx.Ops))
	}

	o = op.Offset(image.Pt(Px(gtx, e.th.InsidePadding.Left), Px(gtx, e.th.InsidePadding.Top))).Push(gtx.Ops)
	_ = e.Editor.Layout(gtx, e.th.Shaper, *e.Font, e.th.FontSp(), textColor, selectionColor)
	o.Pop()
	if e.Editor.Len() == 0 {
		callHint.Add(gtx.Ops)
	}
	if e.borderThickness > 0 {
		w := float32(Px(gtx, e.borderThickness))
		if e.Focused() {
			paintBorder(gtx, border, e.outlineColor, w*2, r)
		} else if e.hovered {
			paintBorder(gtx, border, e.outlineColor, w*3/2, r)
		} else {
			paintBorder(gtx, border, e.Fg(), w, r)
		}
	}

	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	eventArea := clip.Rect(border).Push(gtx.Ops)
	for _, ev := range gtx.Events(&e.hovered) {
		if ev, ok := ev.(pointer.Event); ok {
			switch ev.Kind {
			case pointer.Leave:
				e.hovered = false
			case pointer.Enter:
				e.hovered = true
			}
		}
	}

	pointer.InputOp{
		Tag:   &e.hovered,
		Kinds: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)
	eventArea.Pop()

	defer op.Offset(image.Pt(Px(gtx, e.th.InsidePadding.Left), 0)).Push(gtx.Ops).Pop()
	dim := image.Pt(gtx.Constraints.Max.X, border.Max.Y+Px(gtx, e.padding.Bottom+e.padding.Top))
	return D{Size: dim}
}

// EditOption is options specific to Edits
type EditOption func(w *EditDef)

// Var is an option parameter to set the variable to be updated
func Var[V Value](s *V) EditOption {
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

func (e *EditDef) setLabelSize(w float32) {
	e.labelSize = w
}

func rr(radius int, height int) int {
	if radius > (height-1)/2 {
		return (height - 1) / 2
	}
	return radius
}

func paintBorder(gtx C, outline image.Rectangle, col color.NRGBA, width float32, rr int) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, rr).Path(gtx.Ops),
			Width: width,
		}.Op(),
	)
}
