// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gio-v/wid"
	"image/color"

	"gioui.org/font/gofont"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	theme        *wid.Theme  // the theme selected
	win          *app.Window // The main window
	form         layout.Widget
	greenFlag    = false // the state variable for the button color
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	homeIcon, _  = wid.NewIcon(icons.ActionHome)
)

func main() {

	win = app.NewWindow(app.Title("Rows demo"), app.Size(unit.Dp(900), unit.Dp(500)))
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = demo(theme)
	go wid.Run(win, &form, theme)
	app.Main()
}

func onClick() {
	greenFlag = !greenFlag
	if greenFlag {
		theme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		theme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
	form = demo(theme)
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Overlay,

		wid.Label(th, "Row examples", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.Label(th, "Button spaced closely, left adjusted"),
		wid.Row(th, nil, wid.SpaceClose,
			wid.RoundButton(th, homeIcon),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Buttons distributed, equal space for each button"),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.RoundButton(th, homeIcon),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Buttons with fixed spacing given by em sizes 7,20,20,20,20,20,"),
		wid.Row(th, nil, []float32{7, 20, 20, 20, 20, 20},
			wid.RoundButton(th, homeIcon),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick), wid.W(150)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Buttons with relative spacing given by wieghts 0.2, 0.4, 0.4, 0.4, 0.4, 0.4"),
		wid.Row(th, nil, []float32{0.2, .4, .4, .4, .4, .4},
			wid.RoundButton(th, homeIcon),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick), wid.W(150)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
	)
}
