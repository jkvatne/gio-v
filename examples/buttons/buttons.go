// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gio-v/wid"
	"image/color"
	"os"

	"gioui.org/f32"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/widget"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var (
	currentTheme *wid.Theme  // the theme selected
	win          *app.Window // The main window
	form         layout.Widget
	name         string
	group        = new(widget.Enum)
	homeIcon     *wid.Icon
	checkIcon    *wid.Icon
	greenFlag            = false // the state variable for the button color
	darkMode             = false
	progress     float32 = 0.33
	homeBg               = wid.RGB(0xF288F2)
	homeFg               = wid.RGB(0x0902200)
)

func main() {
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	homeIcon, _ = wid.NewIcon(icons.ActionHome)
	go func() {
		currentTheme = wid.NewTheme(gofont.Collection(), 14)
		win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
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
	c := currentTheme.Bg(wid.Surface)
	paint.Fill(gtx.Ops, c)
	progress = progress + 0.01
	if progress > 1.0 {
		progress = 0
	}
	// Draw widgets
	form(gtx)
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

func onSwitchMode() {
	currentTheme.DarkMode = darkMode
	form = demo(currentTheme)
}

func onClick() {
	greenFlag = !greenFlag
	if greenFlag {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
	form = demo(currentTheme)
}

func swColor() {
	if greenFlag {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		currentTheme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
}

func onWinChange() {
	switch group.Value {
	case "windowed":
		win.Option(app.Windowed.Option())
	case "minimized":
		win.Option(app.Minimized.Option())
	case "fullscreen":
		win.Option(app.Fullscreen.Option())
	case "maximized":
		win.Option(app.Maximized.Option())
	}
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Overlay, f32.Point{},
		wid.Label(th, "Buttons demo page", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.Label(th, "Buttons with fixed length and large font"),
		wid.Button(th, "Change color", wid.Do(onClick), wid.W(450), wid.Large()),
		wid.Label(th, "Buttons with large font using primary container"),
		wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.FontSize(1.4), wid.Role(wid.PrimaryContainer)),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Button spaced closely, left adjusted"),
		wid.Row(th, nil, wid.SpaceClose,
			wid.RoundButton(th, homeIcon,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Bg(&homeBg), wid.Fg(&homeFg),
				wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Buttons distributed, equal space to each button"),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.RoundButton(th, homeIcon,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Bg(&homeBg), wid.Fg(&homeFg),
				wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Buttons with fixed spacing given by em sizes 7,20,20,20,20,20,"),
		wid.Row(th, nil, []float32{7, 20, 20, 20, 20, 20},
			wid.RoundButton(th, homeIcon,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Bg(&homeBg), wid.Fg(&homeFg),
				wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick), wid.W(150)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),

		wid.Label(th, "Buttons with relative spacing given by wieghts 0.2, 0.4, 0.4, 0.4, 0.4, 0.4"),
		wid.Row(th, nil, []float32{0.2, .4, .4, .4, .4, .4},
			wid.RoundButton(th, homeIcon,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Bg(&homeBg), wid.Fg(&homeFg),
				wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick), wid.W(150)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
	)
}
