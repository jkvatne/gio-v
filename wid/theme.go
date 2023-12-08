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
	// OutlineHighest is the grayest surface
	SurfaceHighest
	SurfaceHigh
	SurfaceLow
	// SurfaceLowest is almost white/black
	SurfaceLowest
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
	// If > 0, the font size will be form height divided by LinesPrForm
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
func Px(gtx C, dp interface{}) int {
	if u, ok := dp.(unit.Dp); ok {
		return gtx.Dp(u * unit.Dp(Scale))
	}
	if u, ok := dp.(unit.Sp); ok {
		return gtx.Sp(u * unit.Sp(Scale))
	}
	panic("Px() called with illegal value")
}

func (th *Theme) FontSp() unit.Sp {
	return th.TextSize * unit.Sp(Scale)
}

func (th *Theme) UpdateColors() {
	if !th.DarkMode {
		th.Fg[Canvas] = Tone(th.NeutralColor, 0)
		th.Bg[Canvas] = Tone(th.NeutralColor, 100)

		th.Fg[Primary] = Tone(th.PrimaryColor, 100)
		th.Bg[Primary] = Tone(th.PrimaryColor, 40)
		th.Fg[PrimaryContainer] = Tone(th.PrimaryColor, 10)
		th.Bg[PrimaryContainer] = Tone(th.PrimaryColor, 90)

		th.Fg[Secondary] = Tone(th.SecondaryColor, 100)
		th.Bg[Secondary] = Tone(th.SecondaryColor, 40)
		th.Fg[SecondaryContainer] = Tone(th.SecondaryColor, 10)
		th.Bg[SecondaryContainer] = Tone(th.SecondaryColor, 90)

		th.Fg[Tertiary] = Tone(th.TertiaryColor, 100)
		th.Bg[Tertiary] = Tone(th.TertiaryColor, 40)
		th.Fg[TertiaryContainer] = Tone(th.TertiaryColor, 10)
		th.Bg[TertiaryContainer] = Tone(th.TertiaryColor, 90)

		th.Fg[Error] = Tone(th.ErrorColor, 100)
		th.Bg[Error] = Tone(th.ErrorColor, 40)
		th.Fg[ErrorContainer] = Tone(th.ErrorColor, 10)
		th.Bg[ErrorContainer] = Tone(th.ErrorColor, 90)

		th.Fg[Outline] = Tone(th.NeutralVariantColor, 40)
		th.Bg[Outline] = Tone(th.NeutralVariantColor, 40)
		th.Fg[OutlineVariant] = Tone(th.NeutralVariantColor, 40)
		th.Bg[OutlineVariant] = Tone(th.NeutralVariantColor, 40)

		th.Fg[SurfaceVariant] = Tone(th.NeutralVariantColor, 40)
		th.Bg[SurfaceVariant] = Tone(th.NeutralVariantColor, 93)
		th.Fg[SurfaceHighest] = Tone(th.NeutralColor, 10)
		th.Bg[SurfaceHighest] = Tone(th.NeutralColor, 92)
		th.Fg[SurfaceHigh] = Tone(th.NeutralColor, 10)
		th.Bg[SurfaceHigh] = Tone(th.NeutralColor, 94)
		th.Fg[Surface] = Tone(th.NeutralColor, 10)
		th.Bg[Surface] = Tone(th.NeutralColor, 96)
		th.Fg[SurfaceLow] = Tone(th.NeutralColor, 10)
		th.Bg[SurfaceLow] = Tone(th.NeutralColor, 98)
		th.Fg[SurfaceLowest] = Tone(th.NeutralColor, 10)
		th.Bg[SurfaceLowest] = Tone(th.NeutralColor, 100)
	} else {
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

		th.Fg[SurfaceVariant] = Tone(th.NeutralVariantColor, 90)
		th.Bg[SurfaceVariant] = Tone(th.NeutralVariantColor, 30)

		th.Fg[SurfaceHighest] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceHighest] = Tone(th.NeutralColor, 22)
		th.Fg[SurfaceHigh] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceHigh] = Tone(th.NeutralColor, 17)
		th.Fg[Surface] = Tone(th.NeutralColor, 90)
		th.Bg[Surface] = Tone(th.NeutralColor, 12)
		th.Fg[SurfaceLow] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceLow] = Tone(th.NeutralColor, 10)
		th.Fg[SurfaceLowest] = Tone(th.NeutralColor, 90)
		th.Bg[SurfaceLowest] = Tone(th.NeutralColor, 4)
	}
	// Borders around edit fields
	th.BorderColor = th.Fg[Outline]
	th.BorderColorHovered = th.Fg[Primary]
	th.BorderColorActive = th.Fg[Primary]
	th.SelectionColor = MulAlpha(th.Bg[Primary], 0x60)
	// Tooltip
	th.TooltipBackground = th.Bg[SecondaryContainer]
	th.TooltipOnBackground = th.Fg[SecondaryContainer]
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
	th.IconInset = layout.Inset{Top: 1, Right: 1, Bottom: 1, Left: 1}
	th.BorderThickness = 1.0
	th.BorderCornerRadius = 5
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
