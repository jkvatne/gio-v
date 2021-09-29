package wid

import (
	"gioui.org/layout"
	"gioui.org/unit"

)

type Widget struct {
	th        *Theme
	hint      string
	padding   layout.Inset
	width     unit.Value
}

type WidgetIf interface {
	SetWidth(width float32)
	SetHint(hint string)
}

type Option interface  {
	Do(cfg interface{})
}

type WidgetOption func(WidgetIf)

func (wid *Widget) SetWidth(width float32) {
	wid.width = unit.Dp(width)
}

func (wid *Widget) SetHint(hint string) {
	wid.hint = hint
}

func (b WidgetOption) Do(cfg interface{}) {
	cc := cfg.(WidgetIf)
	b(cc)
}

func W(width float32) WidgetOption {
	return func(w WidgetIf) {
		w.SetWidth(width)
	}
}

func Hint(hint string) WidgetOption {
	return func(w WidgetIf) {
		w.SetHint(hint)
	}

}
