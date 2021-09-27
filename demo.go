// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates Gio widgets. See https://gioui.org for more information.

import (
	"flag"
	"gio-v/wid"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
	"log"
	"os"
)

// green is the state variable for the button color
var green = false

// currentTheme is the theme selected
var currentTheme *wid.Theme
// root is the root widget (usualy a list), and is the root of the widget tree
var root layout.Widget
// formSize is the current window size
var windowSize image.Point


func main() {
	flag.IntVar(&alt, "alt", 0, "Select windows placement/mode")
	flag.Parse()

	currentTheme = wid.NewTheme(gofont.Collection(), 24.0, wid.MaterialDesignDark)
	go func() {
		w := setupForm(currentTheme);
		for {
			select {
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					log.Fatal(e.Err)
					os.Exit(1)
				case system.FrameEvent:
					handleFrameEvents(currentTheme, e)
				}
			}
		}
	}()
	app.Main()
}


func updateWindowSize(th *wid.Theme, e system.FrameEvent) {
	if windowSize.X != e.Size.X || windowSize.Y != e.Size.Y {
		th.TextSize = unit.Dp(float32(e.Size.X) / 100)
		setupForm(th)
		windowSize = e.Size
	}
}

func handleFrameEvents(th *wid.Theme, e system.FrameEvent) {
	updateWindowSize(th, e)
	var ops op.Ops
	gtx := layout.NewContext(&ops, e)
	// Set background color
	paint.Fill(gtx.Ops, th.Palette.Background)
	// Traverse the widget tree and generate drawing operations
	root(gtx)
	// Do the actual screen drawing
	e.Frame(gtx.Ops)
}

func onClick() {
  	green = !green
	if green {
		thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x00, B: 0xff}
	}
}

var darkMode = true

func onSwitchMode(v bool) {
	darkMode = v
	if !darkMode {
		currentTheme = wid.NewTheme(gofont.Collection(), currentTheme.TextSize.V, wid.MaterialDesignLight)
	} else {
		currentTheme = wid.NewTheme(gofont.Collection(), currentTheme.TextSize.V, wid.MaterialDesignDark)
	}
	setupForm(currentTheme)
}

func doDisable(v bool) {
	wid.GlobalDisable = ! wid.GlobalDisable
}

var thb wid.Theme
var alt int

func setupForm(th *wid.Theme) *app.Window {
	thb = *th
	var w *app.Window
	switch {
	case alt==1:
		// A maximized window
		w = app.NewWindow(app.Title("Gio-v demo"), app.Maximized.Option())
	case alt==2:
		// Default placement of window with fixed size in upper left corner of screen
		w = app.NewWindow(app.Title("Gio-v demo"),app.Size(unit.Px(960), unit.Px(500)))
	case alt==3:
		// Place at a given location.
		w = app.NewWindow(app.Title("Gio-v demo"),app.Size(unit.Px(960), unit.Px(540)), app.Pos(unit.Px(960),unit.Px(540)))
	case alt==4:
		// A full-screen window
		w = app.NewWindow(app.Title("Gio-v demo"), app.Fullscreen.Option())
	default:
		// Place at center of monitor
		w = app.NewWindow(app.Title("Gio-v demo"),app.Size(unit.Px(960), unit.Px(540)), app.Center())

	}
	icon1, _ := wid.NewIcon(icons.ContentAdd)
	icon2, _ := wid.NewIcon(icons.ActionHome)
	icon3, _ := wid.NewIcon(icons.ActionCheckCircle)
	root = wid.MakeList(
		th, layout.Vertical,
		wid.Label(th, "Demo page", text.Middle, 2.0),
		wid.Button(wid.Contained, th, "WIDE BUTTON", wid.Width(1900),
			wid.Hint("This is a dummy button")), //  - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
		wid.MakeFlex(
			wid.Label(th, "Dark mode", text.Start, 1.0),
			wid.Switch(th, darkMode, onSwitchMode),
		),
		wid.TextField(th, "Value 1"),
		wid.TextField(th, "Value 2"),
		wid.Checkbox(th, "Checkbox", darkMode, nil),
		wid.MakeFlex(
			wid.Button(wid.Round, th, "", wid.BtnIcon(icon1)),
			wid.Button(wid.Contained, th, "Home", wid.BtnIcon(icon2), wid.Disable(&darkMode)),
			wid.Button(wid.Contained, th, "Check", wid.BtnIcon(icon3)),
			wid.Button(wid.Contained, &thb, "Change color", wid.Handler(onClick)),
			wid.Button(wid.Text, th, "Text button"),
			wid.Button(wid.Outlined, th, "Outline button"),
			wid.Label(th, "         Disabled", text.End, 1.0),
			wid.Switch(th, false, doDisable),
		),
		wid.MakeFlex(
			wid.Combo(th, 0, []string{"Option A", "Option B", "Option C"}),
			wid.Combo(th, 0, []string{"Option 1", "Option 2", "Option 3"}),
		),
		wid.TextField(th, "Value 4"),
		wid.TextField(th, "Value 5"),
		wid.TextField(th, "Value 6"),
		wid.TextField(th, "Value 7"),
		wid.TextField(th, "Value 8"),
		wid.TextField(th, "Value 9"),
		wid.TextField(th, "Value 10"),
		wid.TextField(th, "Value 11"),
		wid.TextField(th, "Value 12"),
	)

	return w
}
