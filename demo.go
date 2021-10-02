// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio modules
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

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
	"time"
)

var selected string

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
	progressIncrementer := make(chan float32)
	go func() {
		for {
			time.Sleep(time.Millisecond * 50)
			progressIncrementer <- 0.01
		}
	}()
	go func() {
		//onSwitchMode(false)
		currentTheme = wid.NewTheme(gofont.Collection(), 14, wid.MaterialDesignLight)
		w := setupForm(currentTheme)
		for {
			select {
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					log.Fatal(e.Err)
				case system.FrameEvent:
					handleFrameEvents(currentTheme, e)
				}
			case pg := <-progressIncrementer:
				progress += pg
				if progress > 1 {
					progress = 0
				}
				w.Invalidate()
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

var darkMode = false

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
var progress float32

func setupForm(th *wid.Theme) *app.Window {
	thb = *th
	var w *app.Window
	switch {
	case alt == 2:
		// A full-screen window
		w = app.NewWindow(app.Title("Gio-v demo"), app.Fullscreen.Option())
	case alt == 3:
		//Place at a given location.
		w = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Px(960), unit.Px(540)), app.Pos(unit.Px(960), unit.Px(540)))
	case alt == 4:
		//   A maximized window
		w = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Px(1800), unit.Px(990)), app.Maximized.Option())
	case alt == 5:
		// Place at center of monitor
		w = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Px(960), unit.Px(540)), app.Center())
	default:
		w = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Px(1900), unit.Px(1000)))

	}

	// Test with gray as primary color
	// th.Primary = wid.RGB(0x555555)

	root = wid.MakeList(
		th, layout.Vertical,

		wid.Label(th, "Demo page", text.Middle, 2.0),
		wid.Label(th, "A fixed width button at the middle of the screen:", text.Start, 1.0),
		wid.MakeFlex(layout.Horizontal, layout.SpaceSides,
			wid.Button(th, "WIDE CENTERED BUTTON",
				wid.W(500),
				wid.Hint("This is a dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
		),
		wid.Label(th, "Two widgets at the right side of the screen:", text.Start, 1.0),
		wid.MakeFlex(layout.Horizontal, layout.SpaceStart,
			wid.Label(th, "Switch to select dark mode", text.Start, 1.0),
			wid.Switch(th, &darkMode, onSwitchMode),
		),
		wid.Checkbox(th, "Checkbox to select dark mode", &darkMode, onSwitchMode),
		wid.MakeFlex(layout.Horizontal, layout.SpaceEnd,
			wid.RoundButton(th, icons.ContentAdd,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Label(th, "Disabled", text.End, 1.0),
			wid.Switch(th, &wid.GlobalDisable, nil),
		),

		// Note that buttons default to their minimum size, unless set differently, here aligned to right side
		wid.MakeFlex(layout.Horizontal, layout.SpaceStart,
			wid.Button(th, "Home", wid.BtnIcon(icons.ActionHome), wid.Disable(&darkMode)),
			wid.Button(th, "Check", wid.BtnIcon(icons.ActionCheckCircle), wid.W(300)),
			wid.Button(&thb, "Change color", wid.Handler(onClick), wid.W(300)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button"),
		),
		// Row with all buttons at minimum size, spread evenly
		wid.MakeFlex(layout.Horizontal, layout.SpaceEvenly,
			wid.Button(th, "Home", wid.BtnIcon(icons.ActionHome), wid.Disable(&darkMode),wid.Min()),
			wid.Button(th, "Check", wid.BtnIcon(icons.ActionCheckCircle), wid.Min()),
			wid.Button(&thb, "Change color", wid.Handler(onClick), wid.Min()),
			wid.TextButton(th, "Text button", wid.Min()),
			wid.OutlineButton(th, "Outline button", wid.Min()),
		),
		wid.MakeFlex(layout.Horizontal, layout.SpaceEvenly,
			wid.DropDown(th, 0, []string{"Option A", "Option B", "Option C"}, wid.W(150)),
			wid.DropDown(th, 1, []string{"Option 1", "Option 2", "Option 3"}),
			wid.DropDown(th, 2, []string{"Option 1", "Option 2", "Option 3"}),
			wid.DropDown(th, 0, []string{"Option A", "Option B", "Option C"}),
			wid.DropDown(th, 0, []string{"Option A", "Option B", "Option C"}),
		),
		// DropDown defaults to max size, here filling a complete row across the form.
		wid.DropDown(th, 0, []string{"Option X", "Option Y", "Option Z"}),
		// Fixed size in Dp
		wid.Edit(th, wid.Hint("Value 1"), wid.W(300)),
		// Relative size
		wid.Edit(th, wid.Hint("Value 2"), wid.W(0.5)),
		wid.MakeFlex(layout.Horizontal, layout.SpaceEnd,
			wid.RadioButton(th, &selected, "RadioOption1", "Option1"),
			wid.RadioButton(th, &selected, "RadioOption2", "Option2"),
			wid.RadioButton(th, &selected, "RadioOption3", "Option3"),
		),
		// The edit's default to their max size so they each get 1/5 of the row size. The MakeFlex spacing parameter will have no effect.
		wid.Row(layout.SpaceStart,
			wid.Edit(th, wid.Hint("Value 3")),
			wid.Edit(th, wid.Hint("Value 4")),
			wid.Edit(th, wid.Hint("Value 5")),
			wid.Edit(th, wid.Hint("Value 6")),
			wid.Edit(th, wid.Hint("Value 7")),
		),

		wid.MakeFlex(layout.Horizontal, layout.SpaceEnd,
			wid.Label(th, "Name", text.End, 1.0),
			wid.Edit(th, wid.Hint("Enter name here")),
		),
		wid.MakeFlex(layout.Horizontal, layout.SpaceEnd,
			wid.Label(th, "Address", text.End, 1.0),
			wid.Edit(th, wid.Hint("Enter address here")),
		),

		wid.MakeFlex(layout.Horizontal, layout.SpaceEnd,
			wid.MakeFlex(layout.Vertical, layout.SpaceEnd,
				wid.Edit(th, wid.Hint("Value 8")),
				wid.Edit(th, wid.Hint("Value 9")),
				wid.Edit(th, wid.Hint("Value 10")),
				wid.Edit(th, wid.Hint("Value 11")),
				wid.Edit(th, wid.Hint("Value 12")),
			),
			wid.MakeFlex(layout.Vertical, layout.SpaceEnd,
				wid.Edit(th, wid.Hint("Value 13")),
				wid.Edit(th, wid.Hint("Value 14")),
				wid.Edit(th, wid.Hint("Value 15")),
				wid.Edit(th, wid.Hint("Value 16")),
				wid.Edit(th, wid.Hint("Value 17")),
			),
		),

		//wid.ProgressBar(th, &progress),
	)
	return w
}
