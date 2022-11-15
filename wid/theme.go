// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"
	"time"

	"gioui.org/layout"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/text"
	"gioui.org/unit"
)

// UIRole describes the type of UI element
// There are two colors for each UIRole, one for text/icon and one for background
// Typicaly you specify a UIRole for each user elemet (button, checkbox etc).
// Default and Zero value is Neutral.
type UIRole uint8

const (
	// Canvas is the default background
	Canvas UIRole = iota
	// Surface is usually the same as Canvas
	Surface
	// SurfaceVariant is for variation
	SurfaceVariant
	// Primary is for prominent buttons, active states etc
	Primary
	// PrimaryContainer is a light background tinted with Primary color.
	PrimaryContainer
	// Secondary is for less prominent components
	Secondary
	// SecondaryContainer is a light background tinted with Secondary color.
	SecondaryContainer
	// Tertiary is for contrasting elements
	Tertiary
	// TertiaryContainer is a light background tinted with Tertiary color.
	TertiaryContainer
	// Error is usualy red
	Error
	// ErrorContainer is usualy light red
	ErrorContainer
	// Outline is used for frames and buttons
	Outline
	// Undefined can be used to detect missing initializations
	Undefined
)

// Tone is the google material tone implementation
// See: https://m3.material.io/styles/color/the-color-system/key-colors-tones
func Tone(c color.NRGBA, tone int) color.NRGBA {
	h, s, _ := Rgb2hsl(c)
	switch {
	case tone < 5:
		return Black
	case tone < 15:
		return Hsl2rgb(h, s, 0.1)
	case tone < 25:
		return Hsl2rgb(h, s, 0.2)
	case tone < 35:
		return Hsl2rgb(h, s, 0.3)
	case tone < 45:
		return Hsl2rgb(h, s, 0.35)
	case tone < 55:
		return Hsl2rgb(h, s, 0.4)
	case tone < 65:
		return Hsl2rgb(h, s, 0.45)
	case tone < 75:
		return Hsl2rgb(h, s, 0.6)
	case tone < 85:
		return Hsl2rgb(h, s, 0.68)
	case tone < 94:
		return Hsl2rgb(h, s, 0.75)
	case tone < 96:
		return Hsl2rgb(h, s, 0.85)
	case tone < 100:
		return Hsl2rgb(h, s, 0.95)
	}
	return White
}

// Fg returns the text/icon color. This is the OnPrimary, OnBackground... colors
func (th *Theme) Fg(kind UIRole) color.NRGBA {
	if !th.DarkMode {
		switch kind {
		case Canvas:
			return Tone(th.Pallet.NeutralColor, 20)
		case Surface:
			return Tone(th.Pallet.NeutralColor, 20)
		case SurfaceVariant:
			return Tone(th.Pallet.NeutralVariantColor, 30)
		case Outline:
			return Tone(th.Pallet.NeutralColor, 50)
		case Primary:
			return Tone(th.Pallet.PrimaryColor, 100)
		case Secondary:
			return Tone(th.Pallet.SecondaryColor, 100)
		case Tertiary:
			return Tone(th.Pallet.TertiaryColor, 100)
		case Error:
			return Tone(th.Pallet.TertiaryColor, 100)
		case PrimaryContainer:
			return Tone(th.Pallet.PrimaryColor, 10)
		case SecondaryContainer:
			return Tone(th.Pallet.SecondaryColor, 10)
		case TertiaryContainer:
			return Tone(th.Pallet.TertiaryColor, 10)
		case ErrorContainer:
			return Tone(th.Pallet.ErrorColor, 10)
		default:
			return Tone(th.Pallet.NeutralColor, 10)
		}
	} else {
		switch kind {
		case Canvas:
			return Tone(th.Pallet.NeutralColor, 90)
		case Surface:
			return Tone(th.Pallet.NeutralColor, 90)
		case SurfaceVariant:
			return Tone(th.Pallet.NeutralVariantColor, 90)
		case Outline:
			return Tone(th.Pallet.NeutralColor, 60)
		case Primary:
			return Tone(th.Pallet.PrimaryColor, 20)
		case Secondary:
			return Tone(th.Pallet.SecondaryColor, 20)
		case Tertiary:
			return Tone(th.Pallet.TertiaryColor, 20)
		case Error:
			return Tone(th.Pallet.ErrorColor, 20)
		case PrimaryContainer:
			return Tone(th.Pallet.PrimaryColor, 90)
		case SecondaryContainer:
			return Tone(th.Pallet.SecondaryColor, 90)
		case TertiaryContainer:
			return Tone(th.Pallet.TertiaryColor, 90)
		case ErrorContainer:
			return Tone(th.Pallet.TertiaryColor, 90)
		default:
			return Tone(th.Pallet.NeutralColor, 90)
		}
	}
	return RGB(0x000000)
}

// Bg returns the background color used to fill the element.
func (th *Theme) Bg(kind UIRole) color.NRGBA {
	if !th.DarkMode {
		switch kind {
		case Canvas:
			return Tone(th.Pallet.NeutralColor, 99)
		case Surface:
			return Tone(th.Pallet.NeutralColor, 99)
		case SurfaceVariant:
			return Tone(th.Pallet.NeutralVariantColor, 90)
		case Primary:
			return Tone(th.Pallet.PrimaryColor, 40)
		case Secondary:
			return Tone(th.Pallet.SecondaryColor, 40)
		case Tertiary:
			return Tone(th.Pallet.TertiaryColor, 40)
		case Error:
			return Tone(th.Pallet.ErrorColor, 40)
		case PrimaryContainer:
			return Tone(th.Pallet.PrimaryColor, 90)
		case SecondaryContainer:
			return Tone(th.Pallet.SecondaryColor, 90)
		case TertiaryContainer:
			return Tone(th.Pallet.TertiaryColor, 90)
		case ErrorContainer:
			return Tone(th.Pallet.ErrorColor, 90)
		default:
			return Tone(th.Pallet.NeutralColor, 99)
		}
	} else {
		switch kind {
		case Canvas:
			return Tone(th.Pallet.NeutralColor, 10)
		case Surface:
			return Tone(th.Pallet.NeutralColor, 10)
		case SurfaceVariant:
			return Tone(th.Pallet.NeutralVariantColor, 30)
		case Primary:
			return Tone(th.Pallet.PrimaryColor, 80)
		case Secondary:
			return Tone(th.Pallet.SecondaryColor, 80)
		case Tertiary:
			return Tone(th.Pallet.TertiaryColor, 80)
		case Error:
			return Tone(th.Pallet.ErrorColor, 80)
		case PrimaryContainer:
			return Tone(th.Pallet.PrimaryColor, 30)
		case SecondaryContainer:
			return Tone(th.Pallet.SecondaryColor, 30)
		case TertiaryContainer:
			return Tone(th.Pallet.TertiaryColor, 30)
		case ErrorContainer:
			return Tone(th.Pallet.TertiaryColor, 30)
		default:
			return Tone(th.Pallet.NeutralColor, 10)
		}
	}
	return RGB(0xFFFFFFFF)
}

// Pallet is the key colors. All other colors are derived from them
type Pallet struct {
	PrimaryColor        color.NRGBA
	SecondaryColor      color.NRGBA
	TertiaryColor       color.NRGBA
	ErrorColor          color.NRGBA
	NeutralColor        color.NRGBA
	NeutralVariantColor color.NRGBA
}

// Theme contains color/layout settings for all widgets
type Theme struct {
	Pallet
	DarkMode              bool
	Shaper                text.Shaper
	TextSize              unit.Sp
	DefaultFont           text.Font
	CheckBoxChecked       *Icon
	CheckBoxUnchecked     *Icon
	RadioChecked          *Icon
	RadioUnchecked        *Icon
	FingerSize            unit.Dp // FingerSize is the minimum touch target size.
	SelectionColor        color.NRGBA
	BorderThicknessActive unit.Dp
	BorderThickness       unit.Dp
	BorderColor           color.NRGBA
	BorderColorHovered    color.NRGBA
	BorderColorActive     color.NRGBA
	BorderCornerRadius    unit.Dp
	TooltipInset          layout.Inset
	TooltipCornerRadius   unit.Dp
	TooltipWidth          unit.Dp
	TooltipBackground     color.NRGBA
	TooltipOnBackground   color.NRGBA
	LabelPadding          layout.Inset
	EditPadding           layout.Inset
	DropDownPadding       layout.Inset
	IconInset             layout.Inset
	ListInset             layout.Inset
	ButtonPadding         layout.Inset
	ButtonLabelPadding    layout.Inset
	ButtonCornerRadius    unit.Dp
	IconSize              unit.Dp
	// Elevation is the shadow width
	Elevation unit.Dp
	// SashColor is the color of the movable divider
	SashColor  color.NRGBA
	SashWidth  unit.Dp
	TrackColor color.NRGBA
	DotColor   color.NRGBA
	// Tooltip settings
	// HoverDelay is the delay between the cursor entering the tip area
	// and the tooltip appearing.
	HoverDelay time.Duration
	// LongPressDelay is the required duration of a press in the area for
	// it to count as a long press.
	LongPressDelay time.Duration
	// LongPressDuration is the amount of time the tooltip should be displayed
	// after being triggered by a long press.
	LongPressDuration time.Duration
	// FadeDuration is the amount of time it takes the tooltip to fade in
	// and out.
	FadeDuration time.Duration
	RowPadTop    unit.Sp
	RowPadBtm    unit.Sp
	// Scroll bar size
	ScrollMajorPadding unit.Sp
	ScrollMinorPadding unit.Sp
	ScrollMajorMinLen  unit.Sp
	ScrollMinorWidth   unit.Sp
	ScrollCornerRadius unit.Sp
}

type (
	// C is a shortcut for layout.Context
	C = layout.Context
	// D is a shortcut for layout.Dimensions
	D = layout.Dimensions
)

// NewTheme creates a new theme with given FontFace and FontSize, based on the theme t
func NewTheme(fontCollection []text.FontFace, fontSize unit.Sp, colors ...color.NRGBA) *Theme {
	t := new(Theme)
	t.PrimaryColor = RGB(0x45682A)
	t.SecondaryColor = RGB(0x57624E)
	t.TertiaryColor = RGB(0x336669)
	t.ErrorColor = RGB(0xAF2525)
	t.NeutralColor = RGB(0x5D5F5A)
	t.NeutralVariantColor = RGB(0x5C6057)
	if len(colors) >= 1 {
		t.PrimaryColor = colors[0]
	}
	if len(colors) >= 2 {
		t.SecondaryColor = colors[1]
	}
	if len(colors) >= 3 {
		t.TertiaryColor = colors[2]
	}
	if len(colors) >= 4 {
		t.ErrorColor = colors[3]
	}

	t.Shaper = text.NewCache(fontCollection)
	t.TextSize = fontSize
	v := unit.Dp(t.TextSize) * 0.4
	// Icons
	t.CheckBoxChecked = mustIcon(NewIcon(icons.ToggleCheckBox))
	t.CheckBoxUnchecked = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.RadioChecked = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.RadioUnchecked = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))
	t.IconInset = layout.Inset{Top: v, Right: v, Bottom: v, Left: v}
	t.FingerSize = unit.Dp(38)
	// Borders around edit fields
	t.BorderThickness = unit.Dp(t.TextSize) * 0.08
	t.BorderThicknessActive = unit.Dp(t.TextSize) * 0.18
	t.BorderColor = t.Fg(Outline)
	t.BorderColorHovered = t.Fg(Primary)
	t.BorderColorActive = t.Fg(Primary)
	t.BorderCornerRadius = unit.Dp(t.TextSize) * 0.2
	// Shadow
	t.Elevation = unit.Dp(t.TextSize) * 0.5
	// Text
	t.LabelPadding = layout.Inset{
		Top:    unit.Dp(t.TextSize) * 0.4,
		Right:  unit.Dp(t.TextSize) * 0.3,
		Bottom: unit.Dp(t.TextSize) * 0.4,
		Left:   unit.Dp(t.TextSize) * 0.6}
	t.DropDownPadding = t.LabelPadding
	t.SelectionColor = MulAlpha(t.Fg(Primary), 0x60)
	t.EditPadding = t.LabelPadding
	// Buttons
	// ButtonPadding is the margin outside a button, giving distance to other elements
	t.ButtonPadding = layout.Inset{Top: 8, Right: 4, Bottom: 8, Left: 4}
	t.ButtonCornerRadius = unit.Dp(t.TextSize) * 999 // Force rounded buttons
	t.ButtonLabelPadding = layout.Inset{Top: 2, Right: 4, Bottom: 2, Left: 4}
	t.IconSize = unit.Dp(t.TextSize) * 1.5
	// Tooltip
	t.TooltipInset = layout.UniformInset(unit.Dp(10))
	t.TooltipCornerRadius = unit.Dp(6.0)
	t.TooltipWidth = v * 50
	t.TooltipBackground = t.Bg(SecondaryContainer)
	t.TooltipOnBackground = t.Fg(SecondaryContainer)
	// List
	t.ListInset = layout.Inset{
		Top:    v * 0.5,
		Right:  v * 0.75,
		Bottom: v * 0.5,
		Left:   v * 0.75,
	}
	// Resizer
	t.SashColor = WithAlpha(t.Fg(Surface), 0x80)
	t.SashWidth = v * 0.5
	// Switch
	t.TrackColor = WithAlpha(t.Fg(Primary), 0x40)
	t.DotColor = t.Fg(Primary)
	t.RowPadTop = t.TextSize * 0.0
	t.RowPadBtm = t.TextSize * 0.0

	t.ScrollMajorPadding = 2
	t.ScrollMinorPadding = 2
	t.ScrollMajorMinLen = t.TextSize / 1.25
	t.ScrollMinorWidth = t.TextSize / 1.5
	t.ScrollCornerRadius = 3
	return t
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
