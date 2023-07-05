// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gio-v/wid"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	theme *wid.Theme // the theme selected
	form  layout.Widget
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = demo(theme)
	go wid.Run(app.NewWindow(app.Title("Colors"), app.Size(unit.Dp(900), unit.Dp(700))), &form, theme)
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
		wid.Label(th, "90", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 90))),
		wid.Label(th, "95", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 95))),
		wid.Label(th, "99", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 99))),
		wid.Label(th, "100", wid.Large(), wid.Fg(&wid.Black), wid.Bg(aTone(c, 100))),
	)
}

func setDefault() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = demo(theme)
}

func setPallete1() {
	theme.PrimaryColor = wid.RGB(0x57624E)
	theme.SecondaryColor = wid.RGB(0x57624E)
	theme.TertiaryColor = wid.RGB(0x336669)
	theme.ErrorColor = wid.RGB(0xAF2525)
	theme.NeutralColor = wid.RGB(0x1D5D7D)
	theme.NeutralVariantColor = wid.RGB(0x756057)
	theme.UpdateColors()
	form = demo(theme)
}

func setPallete2() {
	theme.PrimaryColor = wid.RGB(0x17624E)
	theme.SecondaryColor = wid.RGB(0x17624E)
	theme.TertiaryColor = wid.RGB(0x136669)
	theme.ErrorColor = wid.RGB(0xAF2535)
	theme.NeutralColor = wid.RGB(0x1D4D7D)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	theme.UpdateColors()
	form = demo(theme)
}
func setPallete3() {
	theme.PrimaryColor = wid.RGB(0x17329E)
	theme.SecondaryColor = wid.RGB(0x17624E)
	theme.TertiaryColor = wid.RGB(0x136669)
	theme.ErrorColor = wid.RGB(0xAF2535)
	theme.NeutralColor = wid.RGB(0x1D4D7D)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	theme.UpdateColors()
	form = demo(theme)
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Overlay,
		wid.Label(th, "Show all tones for some palletes", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Set default pallete", wid.Do(setDefault)),
			wid.Button(th, "Set pallete 1", wid.Do(setPallete1)),
			wid.Button(th, "Set pallete 2", wid.Do(setPallete2)),
			wid.Button(th, "Set pallete 3", wid.Do(setPallete3))),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Primary", wid.Large()),
		showTones(th, th.PrimaryColor),
		wid.Label(th, "Secondary", wid.Large()),
		showTones(th, th.SecondaryColor),
		wid.Label(th, "Tertiary", wid.Large()),
		showTones(th, th.TertiaryColor),
		wid.Label(th, "Error", wid.Large()),
		showTones(th, th.ErrorColor),
		wid.Label(th, "NeutralColor", wid.Large()),
		showTones(th, th.NeutralColor),
		wid.Label(th, "NeutralVariantColor", wid.Large()),
		showTones(th, th.NeutralVariantColor),
	)
}
