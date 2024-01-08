// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/unit"
	"github.com/jkvatne/gio-v/wid"
)

var theme *wid.Theme
var form wid.Wid
var UserName, Password string

func main() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = hello(theme)
	go wid.Run(app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(500), unit.Dp(170))), &form, theme)
	app.Main()
}

func onLogin() {}

func onCancel() {}

func hello(th *wid.Theme) wid.Wid {
	return wid.Col(wid.SpaceClose,
		wid.Label(th, "Enter user name and password", wid.Heading(), wid.Bold()),
		wid.Edit(th, &UserName, wid.Ls(0.2), wid.Lbl("User name")),
		wid.Edit(th, &Password, wid.Ls(0.2), wid.Lbl("Password")),
		wid.Row(th, nil, []float32{25, 0, 0},
			wid.Space(1),
			wid.Button(theme, "Log in", wid.Do(onLogin), wid.W(20)),
			wid.Button(theme, "Cancel", wid.Do(onCancel), wid.W(20)),
		),
	)
}
