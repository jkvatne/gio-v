// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gio-v/wid"
	"os"

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
	address      string
	group        = new(widget.Enum)
	homeIcon     *widget.Icon
	checkIcon    *widget.Icon
	green        = false // the state variable for the button color
)

func main() {
	checkIcon, _ = widget.NewIcon(icons.ActionCheckCircle)
	homeIcon, _ = widget.NewIcon(icons.ActionHome)

	go func() {
		currentTheme = wid.NewTheme(gofont.Collection(), 14, wid.MaterialDesignLight)
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
	paint.Fill(gtx.Ops, currentTheme.Background)
	// Draw widgets
	form(gtx)
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

func onClick() {
	green = !green
	if green {
		// thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		// thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x00, B: 0xff}
	}
}

func onModeChange(mode string) {
	switch mode {
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
	return wid.Col(
		wid.Label(th, "Demo page", wid.Middle(), wid.Large(), wid.Bold()),
		// The edit's default to their max size so they each get 1/5 of the row size. The MakeFlex spacing parameter will have no effect.
		wid.Row(th, nil, nil,
			wid.Edit(th, wid.Hint("Value 3")),
			wid.Edit(th, wid.Hint("Value 4")),
			wid.Edit(th, wid.Hint("Value 5")),
		),
		wid.Row(th, nil, nil,
			wid.Col(
				wid.Edit(th, wid.Hint("Value 6"), wid.Lbl("Value 76")),
				wid.Edit(th, wid.Hint("Value 7"), wid.Lbl("Value 7")),
			),
			wid.Col(
				wid.Edit(th, wid.Lbl("Name"), wid.Var(&name)),
				wid.Edit(th, wid.Lbl("Address"), wid.Var(&address)),
			),
		),
		wid.Row(th, nil, nil,
			wid.Edit(th, wid.Hint("")),
		),
		wid.Row(th, nil, nil,
			wid.RadioButton(th, group, "windowed", "Windowed", wid.Do(onModeChange)),
			wid.RadioButton(th, group, "fullscreen", "Fullscreen", wid.Do(onModeChange)),
			wid.RadioButton(th, group, "minimized", "Minimized", wid.Do(onModeChange)),
			wid.RadioButton(th, group, "maximized", "Maximized", wid.Do(onModeChange)),
		),
		wid.Row(th, nil, nil,
			wid.RoundButton(th, homeIcon, wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")).Layout,
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Fg(0x228822), wid.Hint("This is another hint")).Layout,
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Color(wid.RGB(0xffff00))).Layout,
			wid.Button(th, "Change color", wid.Handler(onClick), wid.W(150)).Layout,
			wid.TextButton(th, "Text button").Layout,
			wid.OutlineButton(th, "Outline button").Layout,
		),
	)
}
