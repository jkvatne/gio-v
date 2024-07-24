// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"github.com/jkvatne/gio-v/wid"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	theme *wid.Theme // the theme selected
	form  layout.Widget
	win   app.Window
	roles = true
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 20)
	show()
	win.Option(app.Title("Colors"), app.Size(1024, 620)) // , app.Maximized.Option())
	go wid.Run(&win, &form, theme)
	app.Main()
}

func aTone(c color.NRGBA, n int) *color.NRGBA {
	col := wid.Tone(c, n)
	return &col
}

func showTones(th *wid.Theme, c color.NRGBA) layout.Widget {
	return wid.Row(th, nil, wid.SpaceDistribute,
		wid.Label(th, "00", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 00))),
		wid.Label(th, "10", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 10))),
		wid.Label(th, "20", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 20))),
		wid.Label(th, "30", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 30))),
		wid.Label(th, "40", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 40))),
		wid.Label(th, "50", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 50))),
		wid.Label(th, "60", wid.Large(), wid.Fg(&wid.White), wid.Bg(aTone(c, 60))),
		wid.Label(th, "70", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 70))),
		wid.Label(th, "80", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 80))),
		wid.Label(th, "90", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 88))),
		wid.Label(th, "94", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 93))),
		wid.Label(th, "98", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 97))),
		wid.Label(th, "100", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 100))),
	)
}

func setDefault() {
	theme = wid.NewTheme(gofont.Collection(), 20)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setPalett1() {
	theme.PrimaryColor = wid.RGB(0x67622E)
	theme.SecondaryColor = wid.RGB(0x27622E)
	theme.TertiaryColor = wid.RGB(0x316669)
	theme.ErrorColor = wid.RGB(0xAF1515)
	theme.NeutralColor = wid.RGB(0x1D5D7D)
	theme.NeutralVariantColor = wid.RGB(0x756057)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setPalett2() {
	theme.PrimaryColor = wid.RGB(0x17624E)
	theme.SecondaryColor = wid.RGB(0x27624E)
	theme.TertiaryColor = wid.RGB(0x136669)
	theme.ErrorColor = wid.RGB(0xAF1505)
	theme.NeutralColor = wid.RGB(0x1D4D7D)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setPalett3() {
	theme.PrimaryColor = wid.RGB(0x17329E)
	theme.SecondaryColor = wid.RGB(0x17624E)
	theme.TertiaryColor = wid.RGB(0x136669)
	theme.ErrorColor = wid.RGB(0xBF0000)
	theme.NeutralColor = wid.RGB(0x1D4D7D)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setColorsRoles() {
	roles = !roles
	show()
}

func setDarkLight() {
	theme.DarkMode = !theme.DarkMode
	show()
}

func show() {
	theme.UpdateColors()
	if roles == true {
		form = demo2(theme)
	} else {
		form = demo1(theme)
	}
}

func demo2(th *wid.Theme) layout.Widget {
	var ld string
	var cr string
	if theme.DarkMode {
		ld = "Set light"
	} else {
		ld = "Set dark"
	}
	if roles {
		cr = "Show Colors"
	} else {
		cr = "Show Roles"
	}
	return wid.Col(wid.SpaceClose,
		wid.Label(th, "Show all UI roles", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Set default", wid.Do(setDefault), wid.Hint("Set the default pallete on all widgets")),
			wid.Button(th, "Set palette 1", wid.Do(setPalett1), wid.Hint("Use a pallete 1")),
			wid.Button(th, "Set palette 2", wid.Do(setPalett2), wid.Hint("Select pallete nr 2")),
			wid.Button(th, "Set palette 3", wid.Do(setPalett3), wid.Hint("Select pallete nr. 3")),
			wid.Button(th, cr, wid.Do(setColorsRoles), wid.Hint("Change between showing color tones and role pallete")),
			wid.Button(th, ld, wid.Do(setDarkLight), wid.Hint("Select light or dark mode")),
		),
		wid.Separator(th, unit.Dp(1.0), wid.Pads(3.0, 0)),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Col(wid.SpaceDistribute,
				wid.Container(th, wid.Primary, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Primary", wid.Large(), wid.Role(wid.Primary))),
				wid.Container(th, wid.Secondary, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Secondary", wid.Large(), wid.Role(wid.Secondary))),
				wid.Container(th, wid.Tertiary, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Tertiary", wid.Large(), wid.Role(wid.Tertiary))),
				wid.Container(th, wid.Error, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Error", wid.Large(), wid.Role(wid.Error))),
				wid.Container(th, wid.PrimaryContainer, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "PrimaryContainer", wid.Large(), wid.Role(wid.PrimaryContainer))),
				wid.Container(th, wid.SecondaryContainer, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SecondaryContainer", wid.Large(), wid.Role(wid.SecondaryContainer))),
				wid.Container(th, wid.TertiaryContainer, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "TertiaryContainer", wid.Large(), wid.Role(wid.TertiaryContainer))),
				wid.Container(th, wid.ErrorContainer, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "ErrorContainer", wid.Large(), wid.Role(wid.ErrorContainer)))),
			wid.Col(wid.SpaceDistribute,
				wid.Container(th, wid.SurfaceContainerHighest, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SurfaceContainerHighest", wid.Large(), wid.Role(wid.SurfaceContainerHighest))),
				wid.Container(th, wid.SurfaceContainerHigh, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SurfaceContainerHigh", wid.Large(), wid.Role(wid.SurfaceContainerHigh))),
				wid.Container(th, wid.SurfaceContainer, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SurfaceContainer", wid.Large(), wid.Role(wid.SurfaceContainer))),
				wid.Container(th, wid.SurfaceContainerLow, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SurfaceContainerLow", wid.Large(), wid.Role(wid.SurfaceContainerLow))),
				wid.Container(th, wid.SurfaceContainerLowest, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SurfaceContainerLowest", wid.Large(), wid.Role(wid.SurfaceContainerLowest))),
				wid.Container(th, wid.Canvas, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Canvas", wid.Large(), wid.Role(wid.Canvas))),
				wid.Container(th, wid.Surface, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Surface", wid.Large(), wid.Role(wid.Surface))),
				wid.Container(th, wid.SurfaceVariant, 0, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "SurfaceVariant", wid.Large(), wid.Role(wid.SurfaceVariant)))),
		),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Set default", wid.Do(setDefault), wid.Hint("Set the default pallete on all widgets")),
			wid.Button(th, "Set palette 1", wid.Do(setPalett1), wid.Hint("Use a pallete 1")),
			wid.Button(th, "Set palette 2", wid.Do(setPalett2), wid.Hint("Select pallete nr 2")),
			wid.Button(th, "Set palette 3", wid.Do(setPalett3), wid.Hint("Select pallete nr. 3")),
			wid.Button(th, cr, wid.Do(setColorsRoles), wid.Hint("Change between showing color tones and role pallete")),
			wid.Button(th, ld, wid.Do(setDarkLight), wid.Hint("Select light or dark mode")),
		),
	)
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo1(th *wid.Theme) layout.Widget {
	var ld string
	var cr string
	if theme.DarkMode {
		ld = "Set light"
	} else {
		ld = "Set dark"
	}
	if roles {
		cr = "Show Colors"
	} else {
		cr = "Show Roles"
	}
	return wid.Col([]float32{0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1},
		wid.Label(th, "Show all tones for some palettes", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Label(th, "Also demonstrates a form that will fill the screen 100%", wid.Middle(), wid.Small()),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Set default", wid.Do(setDefault), wid.Hint("Set the default pallete on all widgets")),
			wid.Button(th, "Set palette 1", wid.Do(setPalett1), wid.Hint("Use a pallete 1")),
			wid.Button(th, "Set palette 2", wid.Do(setPalett2), wid.Hint("Select pallete nr 2")),
			wid.Button(th, "Set palette 3", wid.Do(setPalett3), wid.Hint("Select pallete nr. 3")),
			wid.Button(th, cr, wid.Do(setColorsRoles), wid.Hint("Change between showing color tones and role pallete")),
			wid.Button(th, ld, wid.Do(setDarkLight), wid.Hint("Select light or dark mode")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Primary"),
		showTones(th, th.PrimaryColor),
		wid.Label(th, "Secondary"),
		showTones(th, th.SecondaryColor),
		wid.Label(th, "Tertiary"),
		showTones(th, th.TertiaryColor),
		wid.Label(th, "Error"),
		showTones(th, th.ErrorColor),
		wid.Label(th, "NeutralColor"),
		showTones(th, th.NeutralColor),
		wid.Label(th, "NeutralVariantColor"),
		showTones(th, th.NeutralVariantColor),
	)
}
