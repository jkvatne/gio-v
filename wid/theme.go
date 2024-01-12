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
// Typicaly you specify a UIRole for each user element (button, checkbox etc.).
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
	OutlineVariant
	// SurfaceContainerHighest is the grayest surface
	SurfaceContainerHighest
	SurfaceContainerHigh
	SurfaceContainer
	SurfaceContainerLow
	// SurfaceContainerLowest is almost white/black
	SurfaceContainerLowest
	TransparentSurface
	RoleCount
)

// Tone is the Google material tone implementation
func Tone(c color.NRGBA, tone int) color.NRGBA {
	h, s, _ := Rgb2hsl(c)
	return Hsl2rgb(h, s, float64(tone)/100.0)
}

// Theme contains color/layout settings for all widgets
type Theme struct {
	PrimaryColor        color.NRGBA
	SecondaryColor      color.NRGBA
	TertiaryColor       color.NRGBA
	ErrorColor          color.NRGBA
	NeutralColor        color.NRGBA
	NeutralVariantColor color.NRGBA
	Bg                  [RoleCount]color.NRGBA
	Fg                  [RoleCount]color.NRGBA
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
	DefaultMargin       layout.Inset
	DefaultPadding      layout.Inset
	IconInset           layout.Inset
	ListInset           layout.Inset
	ButtonPadding       layout.Inset
	ButtonMargin        layout.Inset
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
	// Extra scaling of the Dp unit
	Scale           float32
	DialogPadding   layout.Inset
	DialogCorners   unit.Dp
	DialogTextWidth unit.Sp
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func uniformPadding(p float32) layout.Inset {
	pp := unit.Dp(p)
	return layout.Inset{Top: pp, Bottom: pp, Left: pp, Right: pp}
}

func (th *Theme) Dp(x unit.Dp) unit.Dp {
	return x
}

type GuiUnit interface{ unit.Dp | unit.Sp }

// Px will convert a size given in either Dp or Sp to pixels
// It applies the theme's scaling factor in addition to
// the gtx metric's PixelPrSp and PixelPrDp
func Px(gtx C, dp interface{}) int {
	if u, ok := dp.(unit.Dp); ok {
		return gtx.Dp(u)
	}
	if u, ok := dp.(unit.Sp); ok {
		return gtx.Sp(u)
	}
	panic("Px() called with illegal value")
}

// UpdateColors must be called after changing the pallete
// See https://m3.material.io/styles/color/static/baseline
func (th *Theme) UpdateColors() {
	if !th.DarkMode {

		// Light mode

		th.Fg[Canvas] = Tone(th.NeutralColor, 0)
		th.Bg[Canvas] = Tone(th.NeutralColor, 100)

		th.Fg[Primary] = Tone(th.PrimaryColor, 100)         // #FFFFFF
		th.Bg[Primary] = Tone(th.PrimaryColor, 48)          // #6750A4
		th.Fg[PrimaryContainer] = Tone(th.PrimaryColor, 10) // #21005D
		th.Bg[PrimaryContainer] = Tone(th.PrimaryColor, 90) // #EADDFF

		th.Fg[Secondary] = Tone(th.SecondaryColor, 100)
		th.Bg[Secondary] = Tone(th.SecondaryColor, 40)
		th.Fg[SecondaryContainer] = Tone(th.SecondaryColor, 10) // #1D192B
		th.Bg[SecondaryContainer] = Tone(th.SecondaryColor, 87) // #E8DEF8

		th.Fg[Tertiary] = Tone(th.TertiaryColor, 100)
		th.Bg[Tertiary] = Tone(th.TertiaryColor, 41)
		th.Fg[TertiaryContainer] = Tone(th.TertiaryColor, 10) // #1D192B
		th.Bg[TertiaryContainer] = Tone(th.TertiaryColor, 87) // #FFD8E4

		th.Fg[Error] = Tone(th.ErrorColor, 100)
		th.Bg[Error] = Tone(th.ErrorColor, 40)
		th.Fg[ErrorContainer] = Tone(th.ErrorColor, 30) // #410E0B
		th.Bg[ErrorContainer] = Tone(th.ErrorColor, 90) // #F9DEDC

		th.Fg[Outline] = Tone(th.NeutralVariantColor, 40)
		th.Bg[Outline] = Tone(th.NeutralVariantColor, 40)
		th.Fg[OutlineVariant] = Tone(th.NeutralVariantColor, 40)
		th.Bg[OutlineVariant] = Tone(th.NeutralVariantColor, 40)

		th.Fg[Surface] = Tone(th.NeutralColor, 10)               // #1D1B20
		th.Bg[Surface] = Tone(th.NeutralColor, 98)               // #FEF7FF
		th.Fg[SurfaceVariant] = Tone(th.NeutralVariantColor, 40) // #49454F
		th.Bg[SurfaceVariant] = Tone(th.NeutralVariantColor, 93) // #E7E0EC

		th.Fg[SurfaceContainerHighest] = Tone(th.NeutralColor, 10) // #1D1B20
		th.Bg[SurfaceContainerHighest] = Tone(th.NeutralColor, 90) // #E6E0E9
		th.Fg[SurfaceContainerHigh] = Tone(th.NeutralColor, 10)    // #1D1B20
		th.Bg[SurfaceContainerHigh] = Tone(th.NeutralColor, 92)    // #ECE6F0
		th.Fg[SurfaceContainer] = Tone(th.NeutralColor, 10)        // #1D1B20
		th.Bg[SurfaceContainer] = Tone(th.NeutralColor, 94)        // #F3EDF7
		th.Fg[SurfaceContainerLow] = Tone(th.NeutralColor, 10)     // #1D1B20
		th.Bg[SurfaceContainerLow] = Tone(th.NeutralColor, 96)     // #F7F2FA
		th.Fg[SurfaceContainerLowest] = Tone(th.NeutralColor, 10)  // #1D1B20
		th.Bg[SurfaceContainerLowest] = Tone(th.NeutralColor, 100) // #FFFFFF
		th.Bg[TransparentSurface] = MulAlpha(th.Fg[SurfaceContainer], 199)
		th.Fg[TransparentSurface] = MulAlpha(th.Fg[SurfaceContainer], 100)

	} else {

		// Dark mode

		th.Fg[Canvas] = Tone(th.NeutralColor, 100)
		th.Bg[Canvas] = Tone(th.NeutralColor, 0)

		th.Fg[Primary] = Tone(th.PrimaryColor, 20)
		th.Bg[Primary] = Tone(th.PrimaryColor, 80)
		th.Fg[PrimaryContainer] = Tone(th.PrimaryColor, 90)
		th.Bg[PrimaryContainer] = Tone(th.PrimaryColor, 30)

		th.Fg[Secondary] = Tone(th.SecondaryColor, 20)
		th.Bg[Secondary] = Tone(th.SecondaryColor, 80)
		th.Fg[SecondaryContainer] = Tone(th.SecondaryColor, 90)
		th.Bg[SecondaryContainer] = Tone(th.SecondaryColor, 30)

		th.Fg[Tertiary] = Tone(th.TertiaryColor, 20)
		th.Bg[Tertiary] = Tone(th.TertiaryColor, 80)
		th.Fg[TertiaryContainer] = Tone(th.TertiaryColor, 90)
		th.Bg[TertiaryContainer] = Tone(th.TertiaryColor, 30)

		th.Fg[Error] = Tone(th.ErrorColor, 20)
		th.Bg[Error] = Tone(th.ErrorColor, 80)
		th.Fg[ErrorContainer] = Tone(th.ErrorColor, 90)
		th.Bg[ErrorContainer] = Tone(th.ErrorColor, 30)

		th.Fg[Outline] = Tone(th.NeutralVariantColor, 60)
		th.Bg[Outline] = Tone(th.NeutralVariantColor, 60)
		th.Fg[OutlineVariant] = Tone(th.NeutralVariantColor, 30)
		th.Bg[OutlineVariant] = Tone(th.NeutralVariantColor, 30)

		th.Fg[Surface] = Tone(th.NeutralColor, 90)
		th.Bg[Surface] = Tone(th.NeutralColor, 12)
		th.Fg[SurfaceVariant] = Tone(th.NeutralVariantColor, 90)
		th.Bg[SurfaceVariant] = Tone(th.NeutralVariantColor, 30)

		th.Fg[SurfaceContainerHighest] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceContainerHighest] = Tone(th.NeutralColor, 22)
		th.Fg[SurfaceContainerHigh] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceContainerHigh] = Tone(th.NeutralColor, 17)
		th.Fg[SurfaceContainer] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceContainer] = Tone(th.NeutralColor, 13)
		th.Fg[SurfaceContainerLow] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceContainerLow] = Tone(th.NeutralColor, 9)
		th.Fg[SurfaceContainerLowest] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceContainerLowest] = Tone(th.NeutralColor, 4)
	}
	// Borders around edit fields
	th.BorderColor = th.Fg[Outline]
	th.BorderColorHovered = th.Fg[Primary]
	th.BorderColorActive = th.Fg[Primary]
	th.SelectionColor = MulAlpha(th.Bg[Primary], 0x60)
	// Tooltip
	th.TooltipBackground = th.Bg[TertiaryContainer]
	th.TooltipOnBackground = th.Fg[TertiaryContainer]
	// Resizer
	th.SashColor = WithAlpha(th.Fg[Surface], 0x40)
	// Switch
	th.TrackColor = th.Bg[Surface]
	th.DotColor = th.Fg[Primary]
}

// NewTheme creates a new theme with given font size and pallete
// The pallet can be left out, to use the defaults - or include as many colors you like.
func NewTheme(fontCollection []text.FontFace, fontSize unit.Sp, colors ...color.NRGBA) *Theme {
	th := new(Theme)
	th.Scale = 1.0
	th.TextSize = fontSize
	// Set up the default pallete
	th.PrimaryColor = RGB(0x6750A4)
	th.SecondaryColor = RGB(0x625B71)
	th.TertiaryColor = RGB(0x567E3E)
	th.ErrorColor = RGB(0xCF1010)
	th.NeutralColor = RGB(0x79747E)
	th.NeutralVariantColor = RGB(0x79747E)
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
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(fontCollection))
	// Default to equal length for label and editor
	th.LabelSplit = 0.5
	th.FingerSize = unit.Dp(38)
	th.IconInset = layout.Inset{Top: 1, Right: 1, Bottom: 1, Left: 1}
	th.BorderThickness = 1.0
	th.BorderCornerRadius = 4.0
	// Shadow
	th.Elevation = 0.5
	// Text
	th.DefaultMargin = uniformPadding(6.0)
	th.DefaultPadding = layout.Inset{4, 4, 4, 2}
	th.ButtonPadding = uniformPadding(6.0)
	th.ButtonCornerRadius = th.BorderCornerRadius
	th.ButtonMargin = uniformPadding(4.0)
	th.IconSize = 20.0
	th.TooltipCornerRadius = th.BorderCornerRadius
	th.TooltipWidth = 250.0
	th.SashWidth = 8.0
	th.RowPadTop = 0.0
	th.RowPadBtm = 0.0
	th.ScrollMajorPadding = 2
	th.ScrollMinorPadding = 2
	th.ScrollMajorMinLen = 15.5
	th.ScrollMinorWidth = 10
	th.ScrollCornerRadius = 4.0
	th.TooltipInset = layout.UniformInset(1)
	th.DialogPadding = layout.Inset{33, 13, 33, 33}
	th.DialogCorners = 20
	th.DialogTextWidth = th.TextSize * 20
	// Update all colors from the pallete
	th.UpdateColors()
	return th
}
