// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"image"
	"image/color"
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

func DefaultEditDef(th *Theme) EditDef {
	return EditDef{
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
}

// Edit will return a widget (layout function) for a text editor
func Edit(th *Theme, options ...any) layout.Widget {
	e := DefaultEditDef(th)
	// The first option should be the value. Will panic if no option is used.
	e.value = options[0]
	// Option 2 can be the number of decimals (if integer or pointer to integer)
	i := 0
	e.DpNo = &i
	if len(options) >= 2 {
		if v, ok := options[1].(int); ok {
			i := v
			e.DpNo = &i
		} else if v, ok := options[1].(*int); ok {
			e.DpNo = v
		}
	}

	// Read in all options to change from default values to something else.
	for _, option := range options {
		if v, ok := option.(UIRole); ok {
			e.role = v
		} else if v, ok := option.(layout.Inset); ok {
			e.padding = v
		} else if v, ok := option.(Option); ok {
			b := &e
			v.apply(b)
		} else if v, ok := option.(rune); ok {
			e.Mask = v
		}
	}

	return func(gtx C) D {
		return e.Layout(gtx)
	}
}

func (e *EditDef) updateValue(gtx C) {
	if !gtx.Focused(&e.Editor) && e.value != nil {
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
	e.wasFocused = gtx.Focused(&e.Editor)
}

func (e *EditDef) Layout(gtx C) D {
	e.CheckDisable(gtx)
	// Precalculate margin and pdding in pixels
	mt, mb, ml, mr := ScaleInset(gtx, e.margin)
	pt, pb, pl, pr := ScaleInset(gtx, e.padding)
	// Make macro for drawing text
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
	textColorOps := macro.Stop()
	// Make macro for drawing hint
	macro = op.Record(gtx.Ops)
	paint.ColorOp{Color: MulAlpha(e.Fg(), 128)}.Add(gtx.Ops)
	hintColorOps := macro.Stop()
	// Make macro for drawing selection
	macro = op.Record(gtx.Ops)
	paint.ColorOp{Color: e.th.SelectionColor}.Add(gtx.Ops)
	selectionColorOps := macro.Stop()
	// Update value
	e.updateValue(gtx)
	// Move to offset the outside margin
	defer op.Offset(image.Pt(pl, pt)).Push(gtx.Ops).Pop()
	// And reduce the size to make space for the padding and margin
	gtx.Constraints.Min.X -= pl + pr + ml + mr
	gtx.Constraints.Max.X = Max(100, gtx.Constraints.Min.X)
	// Draw hint text with top/left padding offset
	macro = op.Record(gtx.Ops)
	o := op.Offset(image.Pt(pl, pt)).Push(gtx.Ops)
	paint.ColorOp{Color: MulAlpha(e.Fg(), 110)}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Editor.Alignment, MaxLines: 1}
	LblDim := tl.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize*unit.Sp(e.FontScale), e.hint, hintColorOps)
	o.Pop()
	callHint := macro.Stop()
	// Add outside label to the left of the edit box
	if e.label != "" {
		o := op.Offset(image.Pt(0, pt)).Push(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		oldMaxX := gtx.Constraints.Max.X
		ofs := int(float32(oldMaxX) * e.labelSize)
		gtx.Constraints.Max.X = ofs - pl
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		colMacro := op.Record(gtx.Ops)
		paint.ColorOp{Color: e.Fg()}.Add(gtx.Ops)
		ll := widget.Label{Alignment: text.End, MaxLines: 1}
		ll.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize*unit.Sp(e.FontScale), e.label, colMacro.Stop())
		o.Pop()
		gtx.Constraints.Max.X = oldMaxX - ofs
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		defer op.Offset(image.Pt(ofs, 0)).Push(gtx.Ops).Pop()
	}
	// If a width is given, and it is within constraints, limit size
	if w := Px(gtx, e.width); w > gtx.Constraints.Min.X && w < gtx.Constraints.Max.X {
		gtx.Constraints.Min.X = w
	}
	// Calculate border size and fill it with white/black when focused
	border := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X+pl+pr, LblDim.Size.Y+pb+pt)}
	rr := Min(Px(gtx, e.th.BorderCornerRadius), border.Max.Y/2)
	if gtx.Focused(&e.Editor) {
		paint.FillShape(gtx.Ops, e.th.Bg[Canvas], clip.UniformRRect(border, rr).Op(gtx.Ops))
	}
	// Move to get the padding needed
	o = op.Offset(image.Pt(pl, pt)).Push(gtx.Ops)
	// Now layout the editor itself
	e.Editor.Layout(gtx, e.th.Shaper, *e.Font, e.th.TextSize*unit.Sp(e.FontScale), textColorOps, selectionColorOps)
	o.Pop()
	// If the editor is empty, we display the hint text
	if e.Editor.Len() == 0 {
		callHint.Add(gtx.Ops)
	}
	// Draw the border, if present
	if e.borderThickness > 0 {
		w := float32(Px(gtx, e.borderThickness))
		if gtx.Focused(&e.Editor) {
			paintBorder(gtx, border, e.outlineColor, w*2, rr)
		} else if e.hovered {
			paintBorder(gtx, border, e.outlineColor, w*3/2, rr)
		} else {
			paintBorder(gtx, border, e.Fg(), w, rr)
		}
	}
	// Setup the pointer event handling
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	eventArea := clip.Rect(border).Push(gtx.Ops)
	event.Op(gtx.Ops, e)
	eventArea.Pop()
	for {
		event, ok := gtx.Event(pointer.Filter{
			Target: e,
			Kinds:  pointer.Enter | pointer.Leave,
		})
		if !ok {
			break
		}
		ev, ok := event.(pointer.Event)
		if !ok {
			continue
		}
		switch ev.Kind {
		case pointer.Leave:
			e.hovered = false
		case pointer.Enter:
			e.hovered = true
		}
	}

	// Calculate size, including margins
	dim := image.Pt(gtx.Constraints.Max.X, border.Max.Y+mb+mt)
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
