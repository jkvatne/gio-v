// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"
	"time"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
)

// UIRole describes the type of UI element
// There are two colors for each UIRole, one for text/icon and one for background
// Typicaly you specify a UIRole for each user element (button, checkbox etc).
// Default and Zero value is Canvas which gives black text/borders on white background.
type UIRole uint8

const (
	// Canvas is white/black. Used in edits, dropdowns etc. to standout
	Canvas UIRole = iota
	// Surface is the default surface for windows.
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
)

// Tone is the Google material tone implementation
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
		case Canvas: // Black
			return Tone(th.NeutralColor, 0)
		case Surface: // Black
			return Tone(th.NeutralColor, 0)
		case SurfaceVariant: // Black
			return Tone(th.NeutralVariantColor, 0)
		case Outline:
			return Tone(th.NeutralColor, 50)
		case Primary:
			return Tone(th.PrimaryColor, 100)
		case Secondary:
			return Tone(th.SecondaryColor, 100)
		case Tertiary:
			return Tone(th.TertiaryColor, 100)
		case Error:
			return Tone(th.TertiaryColor, 100)
		case PrimaryContainer:
			return Tone(th.PrimaryColor, 10)
		case SecondaryContainer:
			return Tone(th.SecondaryColor, 10)
		case TertiaryContainer:
			return Tone(th.TertiaryColor, 10)
		case ErrorContainer:
			return Tone(th.ErrorColor, 10)
		default:
			return Tone(th.NeutralColor, 10)
		}
	} else {
		switch kind {
		case Canvas: // White
			return Tone(th.NeutralColor, 80)
		case Surface: // Light silver
			return Tone(th.NeutralColor, 100)
		case SurfaceVariant: // Some other very light color
			return Tone(th.NeutralVariantColor, 90)
		case Outline:
			return Tone(th.NeutralColor, 60)
		case Primary:
			return Tone(th.PrimaryColor, 20)
		case Secondary:
			return Tone(th.SecondaryColor, 20)
		case Tertiary:
			return Tone(th.TertiaryColor, 20)
		case Error:
			return Tone(th.ErrorColor, 20)
		case PrimaryContainer:
			return Tone(th.PrimaryColor, 90)
		case SecondaryContainer:
			return Tone(th.SecondaryColor, 90)
		case TertiaryContainer:
			return Tone(th.TertiaryColor, 90)
		case ErrorContainer:
			return Tone(th.TertiaryColor, 90)
		default:
			return Tone(th.NeutralColor, 90)
		}
	}
}

// Bg returns the background color used to fill the element.
func (th *Theme) Bg(kind UIRole) color.NRGBA {
	if !th.DarkMode {
		switch kind {
		case Canvas: // White background
			return Tone(th.NeutralColor, 100)
		case Surface: // Light silver background
			return Tone(th.NeutralColor, 99)
		case SurfaceVariant: // Some other light background
			return Tone(th.NeutralVariantColor, 90)
		case Primary:
			return Tone(th.PrimaryColor, 40)
		case Secondary:
			return Tone(th.SecondaryColor, 40)
		case Tertiary:
			return Tone(th.TertiaryColor, 40)
		case Error:
			return Tone(th.ErrorColor, 40)
		case PrimaryContainer:
			return Tone(th.PrimaryColor, 90)
		case SecondaryContainer:
			return Tone(th.SecondaryColor, 90)
		case TertiaryContainer:
			return Tone(th.TertiaryColor, 90)
		case ErrorContainer:
			return Tone(th.ErrorColor, 90)
		default:
			return Tone(th.NeutralColor, 99)
		}
	} else {
		switch kind {
		case Canvas: // Black background
			return Tone(th.NeutralColor, 0)
		case Surface: // Dark gray background
			return Tone(th.NeutralColor, 10)
		case SurfaceVariant: // Another very dark background
			return Tone(th.NeutralVariantColor, 20)
		case Primary:
			return Tone(th.PrimaryColor, 80)
		case Secondary:
			return Tone(th.SecondaryColor, 80)
		case Tertiary:
			return Tone(th.TertiaryColor, 80)
		case Error:
			return Tone(th.ErrorColor, 80)
		case PrimaryContainer:
			return Tone(th.PrimaryColor, 30)
		case SecondaryContainer:
			return Tone(th.SecondaryColor, 30)
		case TertiaryContainer:
			return Tone(th.TertiaryColor, 30)
		case ErrorContainer:
			return Tone(th.TertiaryColor, 30)
		default:
			return Tone(th.NeutralColor, 10)
		}
	}
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
	Scale               float32 // Used for scaling forms
	DarkMode            bool
	Shaper              *text.Shaper
	TextSize            unit.Sp
	DefaultFont         font.Font
	CheckBoxChecked     *Icon
	CheckBoxUnchecked   *Icon
	RadioChecked        *Icon
	RadioUnchecked      *Icon
	FingerSize          unit.Dp // FingerSize is the minimum touch target size.
	SelectionColor      color.NRGBA
	BorderThickness     unit.Dp
	BorderColor         color.NRGBA
	BorderColorHovered  color.NRGBA
	BorderColorActive   color.NRGBA
	BorderCornerRadius  unit.Dp
	TooltipInset        layout.Inset
	TooltipCornerRadius unit.Dp
	TooltipWidth        unit.Dp
	TooltipBackground   color.NRGBA
	TooltipOnBackground color.NRGBA
	OutsidePadding      layout.Inset
	InsidePadding       layout.Inset
	IconInset           layout.Inset
	ListInset           layout.Inset
	ButtonPadding       layout.Inset
	ButtonLabelPadding  layout.Inset
	ButtonCornerRadius  unit.Dp
	IconSize            unit.Dp
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
	RowPadTop    unit.Dp
	RowPadBtm    unit.Dp
	// Scroll bar size
	ScrollMajorPadding unit.Dp
	ScrollMinorPadding unit.Dp
	ScrollMajorMinLen  unit.Dp
	ScrollMinorWidth   unit.Dp
	ScrollCornerRadius unit.Dp
	// Default split between edit label and edit field
	LabelSplit float32
	// If > 0, the font size will by form height divided by LinesPrForm
	LinesPrForm float64
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func uniformPadding(p unit.Dp) layout.Inset {
	return layout.Inset{Top: p, Bottom: p, Left: p, Right: p}
}

func (th *Theme) SetLinesPrForm(x float64) {
	th.LinesPrForm = x
	// Force recalculation of font size
	OldWinY = 0
}

func (th *Theme) Dp(x unit.Dp) unit.Dp {
	return x
}

type GuiUnit interface{ unit.Dp | unit.Sp }

// Px will convert a size given in either Dp or Sp to pixels
// It applies the theme's scaling factor in addition to
// the gtx metric's PixelPrSp and PixelPrDp
func (th *Theme) Px(gtx C, dp interface{}) int {
	if u, ok := dp.(unit.Dp); ok {
		return gtx.Dp(u * unit.Dp(th.Scale))
	}
	if u, ok := dp.(unit.Sp); ok {
		return gtx.Sp(u * unit.Sp(th.Scale))
	}
	return 0
}

func (th *Theme) UpdateColors() {
	// Borders around edit fields
	th.BorderColor = th.Fg(Outline)
	th.BorderColorHovered = th.Fg(Primary)
	th.BorderColorActive = th.Fg(Primary)
	th.SelectionColor = MulAlpha(th.Bg(Primary), 0x60)
	// Tooltip
	th.TooltipBackground = th.Bg(SecondaryContainer)
	th.TooltipOnBackground = th.Fg(SecondaryContainer)
	// Resizer
	th.SashColor = WithAlpha(th.Fg(Surface), 0x40)
	// Switch
	th.TrackColor = th.NeutralColor
	th.DotColor = th.Fg(Primary)
}

// NewTheme creates a new theme with given font size and pallete
// The pallet can be left out, to use the defaults - or include as many colors you like.
func NewTheme(fontCollection []text.FontFace, fontSize unit.Sp, colors ...color.NRGBA) *Theme {
	th := new(Theme)
	th.TextSize = fontSize
	// Set up the default pallete
	th.PrimaryColor = RGB(0x45682A)
	th.SecondaryColor = RGB(0x57624E)
	th.TertiaryColor = RGB(0x336669)
	th.ErrorColor = RGB(0xAF2525)
	th.NeutralColor = RGB(0x5D5D5D)
	th.NeutralVariantColor = RGB(0x756057)
	// Then replace the optional colors in the argument list
	if len(colors) >= 1 {
		th.PrimaryColor = colors[0]
	}
	if len(colors) >= 2 {
		th.SecondaryColor = colors[1]
	}
	if len(colors) >= 3 {
		th.TertiaryColor = colors[2]
	}
	if len(colors) >= 4 {
		th.ErrorColor = colors[3]
	}
	// Setup icons
	th.CheckBoxChecked = mustIcon(NewIcon(icons.ToggleCheckBox))
	th.CheckBoxUnchecked = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	th.RadioChecked = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	th.RadioUnchecked = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))
	// Setup font types
	// Old version (v0.1.0) : t.Shaper = text.NewShaper(fontCollection)
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(fontCollection))
	// Default to equal length for label and editor
	th.LabelSplit = 0.5
	th.FingerSize = unit.Dp(38)
	th.Scale = float32(unit.Dp(th.TextSize) / 14)
	th.IconInset = layout.Inset{Top: 1, Right: 1, Bottom: 1, Left: 1}
	th.BorderThickness = 1.0
	th.BorderCornerRadius = 3
	// Shadow
	th.Elevation = 0.5
	// Text
	th.OutsidePadding = uniformPadding(3.5)
	th.InsidePadding = uniformPadding(3.5)
	th.ButtonPadding = uniformPadding(3.5)
	th.ButtonCornerRadius = th.BorderCornerRadius
	th.ButtonLabelPadding = uniformPadding(5)
	th.IconSize = 20
	th.TooltipCornerRadius = th.BorderCornerRadius
	th.TooltipWidth = 250
	th.SashWidth = 8
	th.RowPadTop = 0.0
	th.RowPadBtm = 0.0
	th.ScrollMajorPadding = 0
	th.ScrollMinorPadding = 0
	th.ScrollMajorMinLen = 15.5
	th.ScrollMinorWidth = 15.5
	th.ScrollCornerRadius = 4
	th.TooltipInset = layout.UniformInset(1)
	// Update all colors from the pallete
	th.UpdateColors()
	return th
}
