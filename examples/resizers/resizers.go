package main

import (
	"github.com/jkvatne/gio-v/wid"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	theme *wid.Theme
	form  layout.Widget
	win   app.Window
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 14)
	form = demo(theme)
	win.Option(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
	wid.Run(&win, &form, theme)
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
