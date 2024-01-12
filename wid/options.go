package wid

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"image/color"
)

// BaseIf is the interface functions for widgets, used by options to set parameters
type BaseIf interface {
	setWidth(width float32)
	setHint(hint string)
	setPadding(padding layout.Inset)
	setMargin(margin layout.Inset)
	setRole(role UIRole)
	setBgColor(c *color.NRGBA)
	setFgColor(c *color.NRGBA)
	setHandler(h func())
	setFont(f *font.Font)
	setDisabler(b *bool)
	getTheme() *Theme
	setFontSize(f float32)
	setBorder(b unit.Dp)
	setDp(dp *int)
	setAlignment(x text.Alignment)
}

// BaseOption is a type for optional parameters when creating widgets
type BaseOption func(BaseIf)

// Option is the interface for optional parameters
type Option interface {
	apply(cfg interface{})
}

// Apply will apply all optional parameters. This can only be used when the widget has no own options.
func (wid *Base) Apply(options ...Option) {
	for _, option := range options {
		option.apply(wid)
	}
}

func (wid BaseOption) apply(cfg interface{}) {
	cc := cfg.(BaseIf)
	wid(cc)
}

func (wid *Base) getTheme() *Theme {
	return wid.th
}

func (wid *Base) setWidth(width float32) {
	wid.width = unit.Dp(width)
}

func (wid *Base) setHint(hint string) {
	wid.hint = hint
}

func (wid *Base) setRole(role UIRole) {
	wid.role = role
}

func (wid *Base) setFont(font *font.Font) {
	wid.Font = font
}

func (wid *Base) setPadding(padding layout.Inset) {
	wid.padding = padding
}

func (wid *Base) setMargin(margin layout.Inset) {
	wid.margin = margin
}

func (wid *Base) setFgColor(c *color.NRGBA) {
	wid.fgColor = c
}

func (wid *Base) setBgColor(c *color.NRGBA) {
	wid.bgColor = c
}

func (wid *Base) setHandler(h func()) {
	wid.onUserChange = h
}

func (wid *Base) setFontSize(h float32) {
	wid.FontScale = float64(h)
}

func (wid *Base) setDisabler(b *bool) {
	wid.disabler = b
}

func (wid *Base) setDp(dp *int) {
	wid.DpNo = dp
}

func (wid *Base) setAlignment(x text.Alignment) {
	wid.Alignment = x
}

func (d *Base) setBorder(w unit.Dp) {
	d.borderWidth = w
}

func Border(b unit.Dp) BaseOption {
	return func(w BaseIf) {
		w.setBorder(b)
	}
}

func Dp[V int | *int](dp V) BaseOption {
	if d, ok := any(dp).(int); ok {
		return func(w BaseIf) {
			i := d
			w.setDp(&i)
		}
	} else if d, ok := any(dp).(*int); ok {
		return func(w BaseIf) {
			w.setDp(d)
		}
	} else {
		panic("Option Dp() must have int or *int as argument")
	}
}

func En(b *bool) BaseOption {
	return func(w BaseIf) {
		w.setDisabler(b)
	}

}

// Pad is used to set default widget paddings (outside of widget)
func (wid *Base) Pad(t, r, b, l float32) {
	wid.padding = layout.Inset{Top: unit.Dp(t), Bottom: unit.Dp(b), Left: unit.Dp(l), Right: unit.Dp(r)}
}

// Do is an optional parameter to set a callback when widget state changes
func Do(f func()) BaseOption {
	return func(w BaseIf) {
		w.setHandler(f)
	}
}

// Middle will align text in the middle.
func Middle() BaseOption {
	return func(d BaseIf) {
		d.setAlignment(text.Middle)
	}
}

// Right will align text to the end.
func Right() BaseOption {
	return func(d BaseIf) {
		d.setAlignment(text.End)
	}
}

// W is the option parameter for setting widget width
func W(width float32) BaseOption {
	return func(w BaseIf) {
		w.setWidth(width)
	}
}

// Hint is an option parameter to set the widget hint (tooltip)
func Hint(hint string) BaseOption {
	return func(w BaseIf) {
		w.setHint(hint)
	}
}

type Color interface {
	*color.NRGBA | color.NRGBA | UIRole
}

// Fg is an option parameter to set widget foreground color
func Fg[V Color](v V) BaseOption {
	if x, ok := any(v).(*color.NRGBA); ok {
		return func(w BaseIf) {
			w.setFgColor(x)
		}
	} else if x, ok := any(v).(color.NRGBA); ok {
		return func(w BaseIf) {
			w.setFgColor(&x)
		}
	} else if x, ok := any(v).(UIRole); ok {
		return func(w BaseIf) {
			c := w.getTheme().Fg[x]
			w.setFgColor(&c)
		}
	} else {
		return nil
	}
}

// Bg is an option parameter to set widget background color
func Bg(c *color.NRGBA) BaseOption {
	return func(w BaseIf) {
		w.setBgColor(c)
	}
}

// Role set the theme role for the widget (Primary, Secondary etc.)
func Role(r UIRole) BaseOption {
	return func(w BaseIf) {
		w.setRole(r)
	}
}

// Lbl is an option parameter to set the widget label
func Lbl(s string) BaseOption {
	return func(w BaseIf) {
		if o, ok := w.(*EditDef); ok {
			o.setLabel(s)
		}
		if o, ok := w.(*DropDownStyle); ok {
			o.setLabel(s)
		}
	}
}

// Ls is an option parameter to set the widget label size
func Ls(x float32) BaseOption {
	return func(w BaseIf) {
		if o, ok := w.(*EditDef); ok {
			o.setLabelSize(x)
		}
		if o, ok := w.(*DropDownStyle); ok {
			o.setLabelSize(x)
		}
	}
}

// Prim is a shortcut to set role=Primary
func Prim() BaseOption {
	return func(w BaseIf) {
		w.setRole(Primary)
	}
}

// PrimCont is a shortcut to set role=PrimaryContainer
func PrimCont() BaseOption {
	return func(w BaseIf) {
		w.setRole(PrimaryContainer)
	}
}

// Sec is a shortcut to set role=Secondary
func Sec() BaseOption {
	return func(w BaseIf) {
		w.setRole(Secondary)
	}
}

// SecCont is a shortcut to set role=SecondaryContainer
func SecCont() BaseOption {
	return func(w BaseIf) {
		w.setRole(SecondaryContainer)
	}
}

// Font set the font for text in the widget
func Font(v *font.Font) BaseOption {
	return func(w BaseIf) {
		w.setFont(v)
	}
}

// FontSize set the font size for text in the widget
func FontSize(v float32) BaseOption {
	return func(w BaseIf) {
		w.setFontSize(v)
	}
}

// Heading makes text 75% larger.
func Heading() BaseOption {
	return func(w BaseIf) {
		w.setFontSize(1.8)
	}
}

// Large makes text 40% larger.
func Large() BaseOption {
	return func(w BaseIf) {
		w.setFontSize(1.3)
	}
}

// Small makes text 20% smaller.
func Small() BaseOption {
	return func(w BaseIf) {
		w.setFontSize(0.8)
	}
}

func Pad(p layout.Inset) BaseOption {
	return func(w BaseIf) {
		w.setPadding(p)
	}
}

// Pads is an option parameter to set customized padding. Noe that 1,2,3 or 4 paddings can be specified.
// If 1 is supplied, it is used for left,right,top,bottom, all with the same padding
// If 2 is supplied, the first is used for top/bottom, and the second for left and right padding
// If 4 is supplied, it is used for top, right, bottom, left in that sequence.
// All values are in Dp (float32 device independent pixels)
func Pads(pads ...float32) BaseOption {
	return func(w BaseIf) {
		switch len(pads) {
		case 0:
			w.setPadding(layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(4), Right: unit.Dp(4)})
		case 1:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[0]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[0])})
		case 2:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[1])})
		case 3:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[2])})
		case 4:
			w.setPadding(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[2]), Left: unit.Dp(pads[3])})
		}
	}
}

// Pads is an option parameter to set customized padding. Noe that 1,2,3 or 4 paddings can be specified.
// If 1 is supplied, it is used for left,right,top,bottom, all with the same padding
// If 2 is supplied, the first is used for top/bottom, and the second for left and right padding
// If 4 is supplied, it is used for top, right, bottom, left in that sequence.
// All values are in Dp (float32 device independent pixels)
func Margin(pads ...float32) BaseOption {
	return func(w BaseIf) {
		switch len(pads) {
		case 0:
			w.setMargin(layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(4), Right: unit.Dp(4)})
		case 1:
			w.setMargin(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[0]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[0])})
		case 2:
			w.setMargin(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[1])})
		case 3:
			w.setMargin(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[0]), Left: unit.Dp(pads[2])})
		case 4:
			w.setMargin(layout.Inset{Top: unit.Dp(pads[0]), Right: unit.Dp(pads[1]), Bottom: unit.Dp(pads[2]), Left: unit.Dp(pads[3])})
		}
	}
}
