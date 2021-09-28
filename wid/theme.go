// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/text"
	"gioui.org/unit"
)

// Palette contains the minimal set of colors that a widget may need to
// draw itself.
type Palette struct {
	// Primary color displayed most frequently across screens and components.
	Primary        color.NRGBA
	PrimaryVariant color.NRGBA

	// Secondary color used sparingly to accent ui elements.
	Secondary        color.NRGBA
	SecondaryVariant color.NRGBA

	// Surface affects surfaces of components such as cards, sheets and menus.
	Surface color.NRGBA

	// Background appears behind scrollable content.
	Background color.NRGBA

	// Error indicates errors in components.
	Error color.NRGBA

	// On colors appear "on top" of the base color.
	// Choose contrasting colors.
	OnPrimary    color.NRGBA
	OnSecondary  color.NRGBA
	OnBackground color.NRGBA
	OnSurface    color.NRGBA
	OnError      color.NRGBA
}

type Theme struct {
	Shaper text.Shaper
	Palette
	TextSize unit.Value
	Icon     struct {
		CheckBoxChecked   *Icon
		CheckBoxUnchecked *Icon
		RadioChecked      *Icon
		RadioUnchecked    *Icon
	}
	// FingerSize is the minimum touch target size.
	FingerSize unit.Value
	HintColor color.NRGBA
	SelectionColor color.NRGBA
	BorderThicknessActive unit.Value
	BorderThickness unit.Value
	BorderColor color.NRGBA
	BorderColorHovered color.NRGBA
	BorderColorActive color.NRGBA
	CornerRadius unit.Value
	TooltipInset layout.Inset
	TooltipCornerRadius unit.Value
	TextTopInset  unit.Value
	LabelInset  layout.Inset
	IconInset  layout.Inset
	// Elevation is the shadow width
	Elevation      unit.Value
	// UmbraColor is the darkest shadow color
	UmbraColor color.NRGBA
	// PenumbraColor is the lightest shadow color
	PenumbraColor color.NRGBA
	// Text inset is the fracion of font height used for padding around text. Typically 0.2 to 0.6
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// WithAlpha returns the input color with the new alpha value.
func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	return color.NRGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: a,
	}
}

func brightness(c uint32) uint32 {
	return (c&0xFF + (c>>8)&0xFF + (c>>16)&0xFF)/3
}

// MaterialDesign is the baseline palette for material design.
// https://material.io/design/color/the-color-system.html#color-theme-creation
var MaterialDesignLight Palette = Palette{
	Primary:          RGB(0x6200EE),
	PrimaryVariant:   RGB(0x3700B3),
	Secondary:        RGB(0x03DAC6),
	SecondaryVariant: RGB(0x018786),
	Background:       RGB(0xFFFFFF),
	Surface:          RGB(0xFFFFFF),
	Error:            RGB(0xB00020),
	OnPrimary:        RGB(0xFFFFFF),
	OnSecondary:      RGB(0x000000),
	OnBackground:     RGB(0x000000),
	OnSurface:        RGB(0x000000),
	OnError:          RGB(0xFFFFFF),
}

var MaterialDesignDark Palette = Palette{
	Primary:          RGB(0xbb86fc),
	PrimaryVariant:   RGB(0x3700b3),
	Secondary:        RGB(0x03DAC6),
	SecondaryVariant: RGB(0x018786),
	Background:       RGB(0x303030),
	Surface:          RGB(0x303030),
	Error:            RGB(0xcf6679),
	OnPrimary:        RGB(0x000000),
	OnSecondary:      RGB(0x000000),
	OnBackground:     RGB(0xffffff),
	OnSurface:        RGB(0xffffff),
	OnError:          RGB(0x000000),
}

func NewTheme(fontCollection []text.FontFace, fontSize float32, p Palette) *Theme {
	t := &Theme{Shaper: text.NewCache(fontCollection)}
	t.Palette = p
	t.TextSize = unit.Sp(fontSize)
	v := t.TextSize.Scale(0.2)
	t.Icon.CheckBoxChecked = mustIcon(NewIcon(icons.ToggleCheckBox))
	t.Icon.CheckBoxUnchecked = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.Icon.RadioChecked = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.Icon.RadioUnchecked = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))
	t.FingerSize = unit.Dp(38)
	t.BorderThickness = t.TextSize.Scale(0.5)
	t.BorderThicknessActive = t.TextSize.Scale(0.5)
	t.BorderColor        = WithAlpha(t.Palette.OnBackground, 128)
	t.BorderColorHovered = WithAlpha(t.Palette.OnBackground, 231)
	t.BorderColorActive  = t.Palette.Primary
	t.CornerRadius = v
	t.Elevation = t.TextSize.Scale(0.5)
	t.LabelInset = layout.Inset{Top:   v, Right: v.Scale(2.0), Bottom: v, Left:   v.Scale(2.0)}
	t.IconInset = layout.Inset{Top:   v, Right: v, Bottom: v, Left:   v}
	t.TooltipInset = layout.UniformInset(unit.Dp(10))
	t.TooltipCornerRadius = unit.Dp(0)
	if approxLuminance(t.OnBackground)<128 {
		t.HintColor = MulAlpha(t.Palette.OnBackground, 0xc0)
	} else {
		t.HintColor = MulAlpha(t.Palette.OnBackground, 0x20)
	}
	t.SelectionColor = MulAlpha(t.Palette.Primary, 0x60)
	return t
}

func (t Theme) Corner(r float32) Theme {
	t.CornerRadius.V = r
	return t
}

func (t Theme) WithPalette(p Palette) Theme {
	t.Palette = p
	return t
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func RGB(c uint32) color.NRGBA {
	return ARGB(0xff000000 | c)
}

func ARGB(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
