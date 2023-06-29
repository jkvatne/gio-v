package main

import (
	"gio-v/wid"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	theme *wid.Theme
	form  layout.Widget
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = demo(theme)
	go wid.Run(app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500))), &form, theme)
	app.Main()
}

func demo(th *wid.Theme) layout.Widget {
	return wid.Col(nil,
		wid.Label(theme, "Resizer demo", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.SplitVertical(th, 0.5,
			wid.ImageFromJpgFile("gopher.jpg", wid.Contain),
			wid.SplitHorizontal(th, 0.5,
				wid.ImageFromJpgFile("gopher.jpg", wid.Contain),
				wid.ImageFromJpgFile("gopher.jpg", wid.Contain),
			),
		),
	)
}
