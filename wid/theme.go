// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"

	"gioui.org/layout"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/text"
	"gioui.org/unit"
)

// Palette is the basic colors (according to Material Design), from where most other collors are derived
type Palette struct {
	Primary      color.NRGBA
	OnPrimary    color.NRGBA
	OnBackground color.NRGBA
	Background   color.NRGBA
	Surface      color.NRGBA
	OnSurface    color.NRGBA
	Error        color.NRGBA
	OnError      color.NRGBA
}

// Theme contains color/layout settings for all widgets
type Theme struct {
	Palette
	Shaper                text.Shaper
	TextSize              unit.Value
	CheckBoxChecked       *Icon
	CheckBoxUnchecked     *Icon
	RadioChecked          *Icon
	RadioUnchecked        *Icon
	FingerSize            unit.Value // FingerSize is the minimum touch target size.
	HintColor             color.NRGBA
	SelectionColor        color.NRGBA
	BorderThicknessActive unit.Value
	BorderThickness       unit.Value
	BorderColor           color.NRGBA
	BorderColorHovered    color.NRGBA
	BorderColorActive     color.NRGBA
	CornerRadius          unit.Value
	TooltipInset          layout.Inset
	TooltipCornerRadius   unit.Value
	TooltipWidth          unit.Value
	TooltipBackground     color.NRGBA
	TooltipOnBackground   color.NRGBA
	LabelPadding          layout.Inset
	EditPadding           layout.Inset
	DropDownPadding       layout.Inset
	IconInset             layout.Inset
	ListInset             layout.Inset
	ButtonPadding         layout.Inset
	ButtonLabelPadding    layout.Inset
	// Elevation is the shadow width
	Elevation unit.Value
	// SashColor is the color of the movable divider
	SashColor  color.NRGBA
	SashWidth  unit.Value
	TrackColor color.NRGBA
}

type (
	// C is a shortcut for layout.Context
	C = layout.Context
	// D is a shortcut for layout.Dimensions
	D = layout.Dimensions
)

// MaterialDesignLight is the baseline palette for material design.
// https://material.io/design/color/the-color-system.html#color-theme-creation
var MaterialDesignLight = Palette{
	Primary:      RGB(0x6200EE),
	Background:   RGB(0xeeeeee),
	Surface:      RGB(0xffffff),
	Error:        RGB(0xB00020),
	OnPrimary:    RGB(0xFFFFFF),
	OnBackground: RGB(0x000000),
	OnSurface:    RGB(0x000000),
	OnError:      RGB(0xFFFFFF),
}

// MaterialDesignDark is the baseline palette for material design dark mode
var MaterialDesignDark = Palette{
	Primary:      RGB(0xbb86fc),
	Background:   RGB(0x303030),
	Surface:      RGB(0x404040),
	Error:        RGB(0xcf6679),
	OnPrimary:    RGB(0x000000),
	OnBackground: RGB(0xffffff),
	OnSurface:    RGB(0xffffff),
	OnError:      RGB(0x000000),
}

// NewTheme creates a new theme with given FontFace and FontSize, based on the theme t
func NewTheme(fontCollection []text.FontFace, fontSize float32, p Palette) *Theme {
	t := new(Theme)
	t.Palette = p
	t.Shaper = text.NewCache(fontCollection)
	t.TextSize = unit.Sp(fontSize)
	v := t.TextSize.Scale(0.2)
	// Icons
	t.CheckBoxChecked = mustIcon(NewIcon(icons.ToggleCheckBox))
	t.CheckBoxUnchecked = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.RadioChecked = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.RadioUnchecked = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))
	t.IconInset = layout.Inset{Top: v, Right: v, Bottom: v, Left: v}
	t.FingerSize = unit.Dp(38)
	// Borders
	t.BorderThickness = t.TextSize.Scale(0.13)
	t.BorderThicknessActive = t.TextSize.Scale(0.18)
	t.BorderColor = WithAlpha(t.OnBackground, 200)
	t.BorderColorHovered = WithAlpha(t.OnBackground, 231)
	t.BorderColorActive = t.Primary
	t.CornerRadius = t.TextSize.Scale(0.2)
	// Shadow
	t.Elevation = t.TextSize.Scale(0.5)
	// Text
	t.LabelPadding = layout.Inset{Top: v, Right: v.Scale(2.0), Bottom: v, Left: v.Scale(2.0)}
	t.DropDownPadding = t.LabelPadding
	t.HintColor = DeEmphasis(t.OnBackground, 45)
	t.SelectionColor = MulAlpha(t.Primary, 0x60)
	t.EditPadding = layout.Inset{Top: v.Scale(2.0), Right: v.Scale(2.0), Bottom: v, Left: v.Scale(2.0)}
	t.ButtonPadding = t.LabelPadding
	t.ButtonLabelPadding = t.LabelPadding
	// Tooltip
	t.TooltipInset = layout.UniformInset(unit.Dp(10))
	t.TooltipCornerRadius = unit.Dp(6.0)
	t.TooltipWidth = t.TextSize.Scale(20)
	t.TooltipBackground = Interpolate(t.OnSurface, t.Surface, 0.9)
	t.TooltipOnBackground = t.OnSurface
	// List
	t.ListInset = layout.Inset{
		Top:    t.TextSize.Scale(0.2),
		Right:  t.TextSize.Scale(0.3),
		Bottom: t.TextSize.Scale(0.2),
		Left:   t.TextSize.Scale(0.3),
	}
	// Resizer
	t.SashColor = WithAlpha(t.OnSurface, 0x80)
	t.SashWidth = t.TextSize.Scale(0.2)
	t.TrackColor = WithAlpha(t.OnSurface, 0x40)
	return t
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
