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
	go func() {
		onSwitchMode(false)
		w := setupForm(currentTheme)
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
	paint.Fill(gtx.Ops, th.Background)
	// Traverse the widget tree and generate drawing operations
	root(gtx)
	// Apply the actual screen drawing
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
	s := float32(24.0)
	if currentTheme != nil {
		s = currentTheme.TextSize.V
	}
	if !darkMode {
		currentTheme = wid.NewTheme(gofont.Collection(), s, wid.MaterialDesignLight)
	} else {
		currentTheme = wid.NewTheme(gofont.Collection(), s, wid.MaterialDesignDark)
	}
	setupForm(currentTheme)
}

func doDisable(v bool) {
	wid.GlobalDisable = !wid.GlobalDisable
}

var thb wid.Theme
var alt int

func setupForm(th *wid.Theme) *app.Window {
	thb = *th
	var w *app.Window
	switch {
	case alt == 2:
		// A full-screen window
		w = app.NewWindow(app.Title("Gio-v demo"), app.Fullscreen.Option())
	case alt==3:
    	//Place at a given location.
		w = app.NewWindow(app.Title("Gio-v demo"),app.Size(unit.Px(960), unit.Px(540)), app.Pos(unit.Px(960),unit.Px(540)))
	case alt==4:
		//   A maximized window
		w = app.NewWindow(app.Title("Gio-v demo"),app.Size(unit.Px(1800), unit.Px(990)), app.Maximized.Option())
	case alt==5:
		// Place at center of monitor
		w = app.NewWindow(app.Title("Gio-v demo"),app.Size(unit.Px(960), unit.Px(540)), app.Center())
	default:
		w = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Px(1900), unit.Px(1000)))

	}

	// Test with gray as primary color
	// th.Primary = wid.RGB(0x555555)

	selected := ""

	root = wid.MakeList(
		th, layout.Vertical,

		wid.Label(th, "Demo page", text.Middle, 2.0),
		wid.Button(th, "WIDE BUTTON",
			wid.W(0.4),
			wid.Pad(30, 15, 15, 0),
			wid.Hint("This is a dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
		wid.MakeFlex(
			wid.Label(th, "Dark mode", text.Start, 1.0),
			wid.Switch(th, &darkMode, onSwitchMode),
		),
		wid.Checkbox(th, "Checkbox", &darkMode, onSwitchMode),
		wid.MakeFlex(
			wid.RoundButton(th, icons.ContentAdd,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(icons.ActionHome), wid.Disable(&darkMode)),
			wid.Button(th, "Check", wid.BtnIcon(icons.ActionCheckCircle)),
			wid.Button(&thb, "Change color", wid.Handler(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button"),
			wid.Label(th, "Disabled", text.End, 1.0),
			wid.Switch(th, &wid.GlobalDisable, nil),
		),
		wid.MakeFlex(
			wid.Combo(th, unit.Value{}, 0, []string{"Option A", "Option B", "Option C"}),
			wid.Combo(th, unit.Value{200, 0}, 1, []string{"Option 1", "Option 2", "Option 3"}),
			wid.Combo(th, unit.Value{300, 0}, 1, []string{"Option 1", "Option 2", "Option 3"}),
			wid.Combo(th, unit.Value{}, 0, []string{"Option A", "Option B", "Option C"}),
		),
		// Fixed size in Dp
		wid.Edit(th, wid.Hint("Value 1"), wid.W(300)),
		// Relative size
		wid.Edit(th, wid.Hint("Value 2"), wid.W(0.5)),
		wid.MakeFlex(
			wid.Edit(th, wid.Hint("Value 3"), wid.W(0.2)),
			wid.Edit(th, wid.Hint("Value 4"), wid.W(0.2)),
			wid.Edit(th, wid.Hint("Value 5"), wid.W(0.2)),
			wid.Edit(th, wid.Hint("Value 6"), wid.W(0.2)),
			wid.Edit(th, wid.Hint("Value 7"), wid.W(0.2)),
		),
		wid.MakeFlex(
			wid.RadioButton(th, &selected, "Option1", "Option1"),
			wid.RadioButton(th, &selected, "Option2", "Option2"),
			wid.RadioButton(th, &selected, "Option3", "Option3"),
		),
		/*
		wid.Edit(th, wid.Hint("Value 8")),
		wid.Edit(th, wid.Hint("Value 9")),
		wid.Edit(th, wid.Hint("Value 10")),
		wid.Edit(th, wid.Hint("Value 11")),
		wid.Edit(th, wid.Hint("Value 12")),
		 */
	)

	return w
}
