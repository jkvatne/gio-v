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
	win   *app.Window
	roles = true
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	theme.InsidePadding = layout.Inset{30, 30, 30, 30}
	theme.OutsidePadding = layout.Inset{20, 20, 30, 30}
	show()
	win = app.NewWindow(app.Title("Colors"), app.Size(1024, 600)) // , app.Maximized.Option())
	go wid.Run(win, &form, theme)
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
	theme = wid.NewTheme(gofont.Collection(), 14)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setPalette1() {
	theme.PrimaryColor = wid.RGB(0x57624E)
	theme.SecondaryColor = wid.RGB(0x57624E)
	theme.TertiaryColor = wid.RGB(0x336669)
	theme.ErrorColor = wid.RGB(0xAF2525)
	theme.NeutralColor = wid.RGB(0x1D5D7D)
	theme.NeutralVariantColor = wid.RGB(0x756057)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setPalette2() {
	theme.PrimaryColor = wid.RGB(0x17624E)
	theme.SecondaryColor = wid.RGB(0x17624E)
	theme.TertiaryColor = wid.RGB(0x136669)
	theme.ErrorColor = wid.RGB(0xAF2535)
	theme.NeutralColor = wid.RGB(0x1D4D7D)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	theme.NeutralVariantColor = wid.RGB(0x356057)
	show()
}

func setPalette3() {
	theme.PrimaryColor = wid.RGB(0x17329E)
	theme.SecondaryColor = wid.RGB(0x17624E)
	theme.TertiaryColor = wid.RGB(0x136669)
	theme.ErrorColor = wid.RGB(0xAF2535)
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
	theme.SetLinesPrForm(28)
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
			wid.Button(th, "Set default palette", wid.Do(setDefault)),
			wid.Button(th, "Set palette 1", wid.Do(setPalette1)),
			wid.Button(th, "Set palette 2", wid.Do(setPalette2)),
			wid.Button(th, "Set palette 3", wid.Do(setPalette3)),
			wid.Button(th, cr, wid.Do(setColorsRoles)),
			wid.Button(th, ld, wid.Do(setDarkLight)),
		),
		wid.Separator(th, unit.Dp(1.0), wid.Pads(3.0, 0)),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Col(wid.SpaceDistribute,
				wid.Label(th, "Primary", wid.Large(), wid.Role(wid.Primary)),
				wid.Label(th, "Secondary", wid.Large(), wid.Role(wid.Secondary)),
				wid.Label(th, "Tertiary", wid.Large(), wid.Role(wid.Tertiary)),
				wid.Label(th, "Error", wid.Large(), wid.Role(wid.Error)),
				wid.Label(th, "Primary Container", wid.Large(), wid.Role(wid.PrimaryContainer)),
				wid.Label(th, "Secondary Container", wid.Large(), wid.Role(wid.SecondaryContainer)),
				wid.Label(th, "Tertiary Container", wid.Large(), wid.Role(wid.TertiaryContainer)),
				wid.Label(th, "Error Containter", wid.Large(), wid.Role(wid.ErrorContainer))),
			wid.Col(wid.SpaceDistribute,
				wid.Label(th, "Surface Variant", wid.Large(), wid.Role(wid.SurfaceVariant)),
				wid.Label(th, "Surface Highest", wid.Large(), wid.Role(wid.SurfaceHighest)),
				wid.Label(th, "Surface High", wid.Large(), wid.Role(wid.SurfaceHigh)),
				wid.Label(th, "Surface", wid.Large(), wid.Role(wid.Surface)),
				wid.Label(th, "Surface Low", wid.Large(), wid.Role(wid.SurfaceLow)),
				wid.Label(th, "Surface Lowest", wid.Large(), wid.Role(wid.SurfaceLowest)),
				wid.Label(th, "Canvas", wid.Large(), wid.Role(wid.Canvas))),
		),
	)
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo1(th *wid.Theme) layout.Widget {
	var ld string
	var cr string
	theme.SetLinesPrForm(31)
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
	return wid.Col(wid.SpaceDistribute,
		wid.Label(th, "Show all tones for some palettes", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Label(th, "Also demonstrates a form that will fill the screen 100%", wid.Middle(), wid.Small()),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Set default palette", wid.Do(setDefault)),
			wid.Button(th, "Set palette 1", wid.Do(setPalette1)),
			wid.Button(th, "Set palette 2", wid.Do(setPalette2)),
			wid.Button(th, "Set palette 3", wid.Do(setPalette3)),
			wid.Button(th, cr, wid.Do(setColorsRoles)),
			wid.Button(th, ld, wid.Do(setDarkLight)),
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
