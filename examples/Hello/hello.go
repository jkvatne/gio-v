// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"gio-v/wid"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/unit"
)

var th *wid.Theme

func main() {
	th = wid.NewTheme(gofont.Collection(), 14)
	form := hello(th)
	go wid.Run(app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500))), &form, th)
	app.Main()
}

func hello(th *wid.Theme) wid.Wid {
	return wid.List(th, wid.Overlay,
		wid.Label(th, "Hello gio..", wid.Heading(), wid.Bold()),
		wid.Label(th, "A small demo program using 28 lines total"),
	)
}
