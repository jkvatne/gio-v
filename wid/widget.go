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
	setWidth(width float32)
	setHint(hint string)
}

type Option interface  {
	apply(cfg interface{})
}

func (b WidgetOption) apply(cfg interface{}) {
	cc := cfg.(WidgetIf)
	b(cc)
}

type WidgetOption func(WidgetIf)

func (wid *Widget) setWidth(width float32) {
	wid.width = unit.Dp(width)
}

func (wid *Widget) setHint(hint string) {
	wid.hint = hint
}

func W(width float32) WidgetOption {
	return func(w WidgetIf) {
		w.setWidth(width)
	}
}

func Hint(hint string) WidgetOption {
	return func(w WidgetIf) {
		w.setHint(hint)
	}

}
