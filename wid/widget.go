package wid

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

type Widget struct {
	th      *Theme
	hint    string
	padding layout.Inset
	width   unit.Value
}

type WidgetIf interface {
	setWidth(width float32)
	setHint(hint string)
	setPadding(padding layout.Inset)
}

type Option interface {
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

func (wid *Widget) setPadding(padding layout.Inset) {
	wid.padding = padding
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

func Pad(pads ...float32) WidgetOption {
	return func(w WidgetIf) {
		switch len(pads) {
		case 0:
			w.setPadding(layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(4), Right: unit.Dp(4)})
		case 1:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[0]), Right: unit.Dp(pads[0])})
		case 2:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[1]), Bottom: unit.Dp(pads[1]), Left: unit.Dp(pads[0]), Right: unit.Dp(pads[0])})
		case 3:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[3]), Bottom: unit.Dp(pads[3]), Left: unit.Dp(pads[1]), Right: unit.Dp(pads[20])})
		case 4:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[1]), Bottom: unit.Dp(pads[3]), Left: unit.Dp(pads[0]), Right: unit.Dp(pads[02])})
		}
	}
}
