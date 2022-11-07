// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"gio-v/wid"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
)

var th *wid.Theme

func main() {
	th = wid.NewTheme(gofont.Collection(), 14)
	go wid.Run(
		app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500))),
		th,
		hello,
	)
	app.Main()
}

func hello(gtx layout.Context) layout.Dimensions {
	return wid.List(th, wid.Overlay,
		wid.Label(th, "Hello gio..", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Label(th, "A small demo program using 31 lines toal"),
	)(gtx)
}
