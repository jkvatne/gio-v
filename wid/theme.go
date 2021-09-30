// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/layout"
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/text"
	"gioui.org/unit"
)

type Theme struct {
	Shaper                text.Shaper
	TextSize              unit.Value
	CheckBoxChecked       *Icon
	CheckBoxUnchecked     *Icon
	RadioChecked          *Icon
	RadioUnchecked        *Icon
	Primary               color.NRGBA
	OnPrimary             color.NRGBA
	OnBackground          color.NRGBA
	Background            color.NRGBA
	Surface               color.NRGBA
	OnSurface             color.NRGBA
	Error                 color.NRGBA
	OnError               color.NRGBA
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
	TextTopInset          unit.Value
	LabelInset            layout.Inset
	IconInset             layout.Inset
	ListInset             layout.Inset
	// Elevation is the shadow width
	Elevation unit.Value
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


// MaterialDesign is the baseline palette for material design.
// https://material.io/design/color/the-color-system.html#color-theme-creation
var MaterialDesignLight = Theme{
	Primary:          RGB(0x6200EE),
	Background:       RGB(0xFFFFFF),
	Surface:          RGB(0xFFFFFF),
	Error:            RGB(0xB00020),
	OnPrimary:        RGB(0xFFFFFF),
	OnBackground:     RGB(0x000000),
	OnSurface:        RGB(0x000000),
	OnError:          RGB(0xFFFFFF),
}

var MaterialDesignDark = Theme {
	Primary:          RGB(0xbb86fc),
	Background:       RGB(0x303030),
	Surface:          RGB(0x303030),
	Error:            RGB(0xcf6679),
	OnPrimary:        RGB(0x000000),
	OnBackground:     RGB(0xffffff),
	OnSurface:        RGB(0xffffff),
	OnError:          RGB(0x000000),
}

func NewTheme(fontCollection []text.FontFace, fontSize float32, t Theme) *Theme {
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
	t.BorderThickness = t.TextSize.Scale(0.1)
	t.BorderThicknessActive = t.TextSize.Scale(0.12)
	t.BorderColor = WithAlpha(t.OnBackground, 128)
	t.BorderColorHovered = WithAlpha(t.OnBackground, 231)
	t.BorderColorActive = t.Primary
	t.CornerRadius = t.TextSize.Scale(0.2)
	// Shadow
	t.Elevation = t.TextSize.Scale(0.5)
	// Text
	t.LabelInset = layout.Inset{Top: v, Right: v.Scale(2.0), Bottom: v, Left: v.Scale(2.0)}
	t.HintColor = DeEmphasis(t.OnBackground, 15)
	t.SelectionColor = MulAlpha(t.Primary, 0x60)
	// Tooltip
	t.TooltipInset = layout.UniformInset(unit.Dp(10))
	t.TooltipCornerRadius = unit.Dp(6.0)
	t.TooltipWidth = t.TextSize.Scale(20)
	t.TooltipBackground = color.NRGBA{255, 255, 160, 233}
	t.TooltipOnBackground = color.NRGBA{A:255}
	// List
	t.ListInset = layout.Inset{
		Top:    t.TextSize.Scale(0.2),
		Right:  t.TextSize.Scale(0.3),
		Bottom: t.TextSize.Scale(0.2),
		Left:   t.TextSize.Scale(0.3),
	}
	return &t
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
