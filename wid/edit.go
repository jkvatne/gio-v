// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"fmt"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"image"
	"image/color"
	"math"
	"strconv"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

type Value interface {
	int | float64 | float32 | string | *int | *float64 | *float32 | *string
}

// EditDef is the parameters for the text editor
type EditDef struct {
	Base
	widget.Editor
	hovered         bool
	outlineColor    color.NRGBA
	selectionColor  color.NRGBA
	label           string
	value           interface{}
	labelSize       float32
	borderThickness unit.Dp
	wasFocused      bool
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
	} else if x, ok := v.(int); ok {
		if x == math.MinInt {
			return "---"
		} else {
			return fmt.Sprintf("%d", x)
		}
	} else if x, ok := v.(*int32); ok {
		if *x == math.MinInt32 {
			return "---"
		} else {
			return fmt.Sprintf("%d", *x)
		}
	} else if x, ok := v.(int32); ok {
		if x == math.MinInt32 {
			return "---"
		} else {
			return fmt.Sprintf("%d", x)
		}
	} else if x, ok := v.(*int64); ok {
		if *x == math.MinInt64 {
			return "---"
		} else {
			return fmt.Sprintf("%d", *x)
		}
	} else if x, ok := v.(int64); ok {
		if x == math.MinInt {
			return "---"
		} else {
			return fmt.Sprintf("%d", x)
		}
	} else if x, ok := v.(*float32); ok {
		if *x == math.MaxFloat32 {
			return "---"
		} else {
			return fmt.Sprintf("%.*f", dp, *x)
		}
	} else if x, ok := v.(*float64); ok {
		if *x == math.MaxFloat64 {
			return "---"
		} else {
			return fmt.Sprintf("%.*f", dp, *x)
		}
	} else if x, ok := v.(*string); ok {
		return *x
	} else if x, ok := v.(string); ok {
		return x
	}
	return ""
}

func StringToValue(value interface{}, current string) {
	if _, ok := value.(*int); ok {
		x, err := strconv.Atoi(current)
		if err == nil {
			*value.(*int) = x
		}
	} else if _, ok := value.(*int32); ok {
		x, err := strconv.Atoi(current)
		if err == nil {
			*value.(*int) = x
		}
	} else if _, ok := value.(*int64); ok {
		x, err := strconv.ParseInt(current, 10, 64)
		if err == nil {
			*value.(*int64) = x
		}
	} else if _, ok := value.(*float32); ok {
		f, err := strconv.ParseFloat(current, 32)
		if err == nil {
			*value.(*float32) = float32(f)
		}
	} else if _, ok := value.(*float64); ok {
		f, err := strconv.ParseFloat(current, 64)
		if err == nil {
			*value.(*float64) = f
		}
	} else if _, ok := value.(*string); ok {
		*value.(*string) = current
	} else {
		panic("Edit value should be pointer to value")
	}
}

// Edit will return a widget (layout function) for a text editor
func Edit(th *Theme, options ...any) func(gtx C) D {
	e := EditDef{
		Base: Base{
			th:        th,
			Font:      &th.DefaultFont,
			FontScale: 1.0,
			margin:    th.DefaultMargin,
			padding:   th.DefaultPadding,
		},
		Editor: widget.Editor{
			SingleLine: true,
		},
		borderThickness: th.BorderThickness,
		labelSize:       th.LabelSplit,
		outlineColor:    th.Fg[Outline],
		selectionColor:  MulAlpha(th.Bg[Primary], 60),
	}

	// The first option should be the value. Will panic if no option is used.
	if v, ok := options[0].(*float64); ok {
		e.value = v
	} else if v, ok := options[0].(*float32); ok {
		e.value = v
	} else if v, ok := options[0].(*int32); ok {
		e.value = v
	} else if v, ok := options[0].(*int16); ok {
		e.value = v
	} else if v, ok := options[0].(*int); ok {
		e.value = v
	} else if v, ok := options[0].(*int64); ok {
		e.value = v
	} else if v, ok := options[0].(*string); ok {
		e.value = v
	}

	i := 0
	e.DpNo = &i
	if len(options) >= 2 {
		// Option 2 can be the number of decimals (if integer)
		if v, ok := options[1].(int); ok {
			i := v
			e.DpNo = &i
		} else if v, ok := options[1].(*int); ok {
			e.DpNo = v
		}
	}

	// Read in options to change from default values to something else.
	for _, option := range options {
		if v, ok := option.(UIRole); ok {
			e.role = v
		} else if v, ok := option.(layout.Inset); ok {
			e.padding = v
		} else if v, ok := option.(Option); ok {
			b := &e
			v.apply(b)
		}
	}

	// Verify that the input value is a pointer and not a value
	if e.value == nil {
		panic("Editor value should be pointer to an integer or float or string value")
	} else {
		GuiLock.Lock()
		s := ValueToString(e.value, *e.DpNo)
		e.Editor.SetText(s)
		GuiLock.Unlock()

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
			e.SetText(ValueToString(e.value, *e.DpNo))
			GuiLock.RUnlock()
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

	// Move to offset the outside margin
	defer op.Offset(image.Pt(
		Px(gtx, e.margin.Left),
		Px(gtx, e.margin.Top))).Push(gtx.Ops).Pop()

	// And reduce the size to make space for the padding and margin
	gtx.Constraints.Min.X -= Px(gtx, e.padding.Left+e.padding.Right+e.margin.Left+e.margin.Right)
	gtx.Constraints.Max.X = gtx.Constraints.Min.X

	// Draw hint text with top/left padding offset
	macro := op.Record(gtx.Ops)
	o := op.Offset(image.Pt(Px(gtx, e.padding.Left), Px(gtx, e.padding.Top))).Push(gtx.Ops)
	paint.ColorOp{Color: MulAlpha(e.Fg(), 110)}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: e.maxLines()}
	LblDim := tl.Layout(gtx, e.th.Shaper, *e.Font, e.th.FontSp()*unit.Sp(e.FontScale), e.hint, hintColor)
	o.Pop()
	callHint := macro.Stop()

	// Add outside label to the left of the edit box
	if e.label != "" {
		o := op.Offset(image.Pt(0, Px(gtx, e.padding.Top))).Push(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		oldMaxX := gtx.Constraints.Max.X
		ofs := int(float32(oldMaxX) * e.labelSize)
		gtx.Constraints.Max.X = ofs - Px(gtx, e.padding.Left)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		ll := widget.Label{Alignment: text.End, MaxLines: 1}
		ll.Layout(gtx, e.th.Shaper, *e.Font, e.th.FontSp()*unit.Sp(e.FontScale), e.label, colMacro.Stop())
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
		gtx.Constraints.Max.X+Px(gtx, e.padding.Left+e.padding.Right),
		LblDim.Size.Y+Px(gtx, e.padding.Bottom+e.padding.Top))}

	r := Px(gtx, e.th.BorderCornerRadius)
	if r > border.Max.Y/2 {
		r = border.Max.Y / 2
	}
	if e.Focused() {
		paint.FillShape(gtx.Ops, e.th.Bg[Canvas], clip.UniformRRect(border, r).Op(gtx.Ops))
	}

	o = op.Offset(image.Pt(Px(gtx, e.padding.Left), Px(gtx, e.padding.Top))).Push(gtx.Ops)
	_ = e.Editor.Layout(gtx, e.th.Shaper, *e.Font, e.th.FontSp()*unit.Sp(e.FontScale), textColor, selectionColor)
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
			default:
			}
		}
	}

	pointer.InputOp{
		Tag:   &e.hovered,
		Kinds: pointer.Enter | pointer.Leave,
	}.Add(gtx.Ops)
	eventArea.Pop()

	defer op.Offset(image.Pt(Px(gtx, e.padding.Left), 0)).Push(gtx.Ops).Pop()
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

func paintBorder(gtx C, outline image.Rectangle, col color.NRGBA, width float32, rr int) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, rr).Path(gtx.Ops),
			Width: width,
		}.Op(),
	)
}
