// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gio-v/wid"
	"image/color"
	"os"

	"gioui.org/op/paint"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

var (
	currentTheme *wid.Theme  // the theme selected
	win          *app.Window // The main window
	form         layout.Widget
)

func main() {
	go func() {
		currentTheme = wid.NewTheme(gofont.Collection(), 14)
		win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(600)))
		form = demo(currentTheme)
		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					os.Exit(0)
				case system.FrameEvent:
					handleFrameEvents(e)
				}
			}
		}
	}()
	app.Main()
}

func handleFrameEvents(e system.FrameEvent) {
	var ops op.Ops
	gtx := layout.NewContext(&ops, e)
	// Set background color
	c := currentTheme.Bg(wid.PrimaryContainer)
	paint.Fill(gtx.Ops, c)
	// Draw widgets
	form(gtx)
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

func showTones(th *wid.Theme, c color.NRGBA) layout.Widget {
	return wid.Row(th, nil, nil,
		wid.Label(th, "00", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 00))),
		wid.Label(th, "10", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 10))),
		wid.Label(th, "20", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 20))),
		wid.Label(th, "30", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 30))),
		wid.Label(th, "40", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 40))),
		wid.Label(th, "50", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 50))),
		wid.Label(th, "60", wid.Large(), wid.Fg(wid.White), wid.Bg(wid.Tone(c, 60))),
		wid.Label(th, "70", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(c, 70))),
		wid.Label(th, "80", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(c, 80))),
		wid.Label(th, "90", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(c, 90))),
		wid.Label(th, "95", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(c, 95))),
		wid.Label(th, "99", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(c, 99))),
		wid.Label(th, "100", wid.Large(), wid.Fg(wid.Black), wid.Bg(wid.Tone(c, 100))),
	)
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Overlay,
		wid.Label(th, "Demo page", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.Label(th, "Colors and tones according to Google Materials 3"),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Primary", wid.Large(), wid.Role(wid.Primary)),
		wid.Label(th, "PrimaryContainer", wid.Large(), wid.Role(wid.PrimaryContainer)),
		showTones(th, th.PrimaryColor),
		wid.Label(th, "Secondary", wid.Large(), wid.Role(wid.Secondary)),
		wid.Label(th, "SecondaryContainer", wid.Large(), wid.Role(wid.SecondaryContainer)),
		showTones(th, th.SecondaryColor),
		wid.Label(th, "Tertiary", wid.Large(), wid.Role(wid.Tertiary)),
		wid.Label(th, "TertiaryContainer", wid.Large(), wid.Role(wid.TertiaryContainer)),
		showTones(th, th.TertiaryColor),
		wid.Label(th, "Error", wid.Large(), wid.Role(wid.Error)),
		wid.Label(th, "ErrorContainer", wid.Large(), wid.Role(wid.ErrorContainer)),
		showTones(th, th.ErrorColor),
		wid.Label(th, "Surface", wid.Large(), wid.Role(wid.Surface)),
		wid.Label(th, "Canvas", wid.Large(), wid.Role(wid.Canvas)),
		showTones(th, th.NeutralColor),
	)
}
