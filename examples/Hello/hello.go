// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"github.com/jkvatne/gio-v/wid"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/unit"
)

var (
	theme *wid.Theme
	form  wid.Wid
	win   app.Window // The main window
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = hello(theme)
	win.Option(app.Title("Gio-v demo"), app.Size(unit.Dp(300), unit.Dp(100)))
	go wid.Run(&win, &form, theme)
	app.Main()
}

func hello(th *wid.Theme) wid.Wid {
	return wid.List(th, wid.Overlay,
		wid.Label(th, "Hello gio..", wid.Heading(), wid.Bold()),
		wid.Label(th, "A small demo program using 28 lines total"),
	)
}
