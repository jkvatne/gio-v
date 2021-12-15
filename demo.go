// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"flag"
	"fmt"
	"gio-v/wid"
	"image"
	"image/color"
	"log"
	"os"
	"runtime"
	"time"

	"gioui.org/widget"

	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var mode = "maximized"
var fontSize = "medium"
var oldMode string
var oldFontSize string
var green = false           // the state variable for the button color
var currentTheme *wid.Theme // the theme selected
var darkMode = false
var oldWindowSize image.Point // the current window size, used to detect changes
var win *app.Window           // The main window
var thb *wid.Theme            // Secondary theme used for the color-shifting button
var progress float32
var sliderValue float32
var dummy bool
var th *material.Theme
var icon *widget.Icon
var addIcon *wid.Icon
var homeIcon *wid.Icon
var checkIcon *wid.Icon
var upIcon *wid.Icon
var downIcon *wid.Icon
var count float64
var startTime time.Time

func main() {
	flag.StringVar(&mode, "mode", "default", "Select window as fullscreen, maximized, centered or default")
	flag.StringVar(&fontSize, "fontsize", "large", "Select font size medium,small,large")
	flag.Parse()
	addIcon, _ = wid.NewIcon(icons.ContentAdd)
	checkIcon, _ = wid.NewIcon(icons.ActionCheckCircle)
	upIcon, _ = wid.NewIcon(icons.HardwareKeyboardArrowUp)
	downIcon, _ = wid.NewIcon(icons.HardwareKeyboardArrowDown)

	makePersons(100)
	ic, err := widget.NewIcon(icons.ContentAdd)
	if err != nil {
		log.Fatal(err)
	}
	icon = ic
	progressIncrementer := make(chan float32)
	go func() {
		for {
			time.Sleep(time.Millisecond * 1)
			progressIncrementer <- 0.001
		}
	}()
	go func() {
		th = material.NewTheme(gofont.Collection())
		currentTheme = wid.NewTheme(gofont.Collection(), 14, wid.MaterialDesignLight)
		win = app.NewWindow(app.Title("Gio-v demo"), modeFromString(mode).Option(), app.Size(unit.Dp(900), unit.Dp(500)))
		updateMode()
		setup()
		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					os.Exit(0)
				case system.FrameEvent:
					handleFrameEvents(e)
				}
			case pg := <-progressIncrementer:
				progress += pg
				if progress > 1 {
					progress = 0
				}
				win.Invalidate()
			}
		}
	}()
	app.Main()
}

func onClick() {
	green = !green
	if green {
		thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		thb.Primary = color.NRGBA{A: 0xff, R: 0x00, G: 0x00, B: 0xff}
	}
}

func update() {
	onSwitchMode(darkMode)
}

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
	setup()
}

func modeFromString(s string) app.WindowMode {
	switch {
	case s == "fullscreen":
		// A full-screen window
		return app.Fullscreen
	case s == "default":
		// Default positioned window with size given
		return app.Windowed
	}
	return app.Windowed
}

func updateMode() {
	if mode != oldMode {
		win.Option(modeFromString(mode).Option())
		oldMode = mode
	}
}

func onMaximize() {
	win.Maximize()
}

func onCenter() {
	win.Center()
}

func column1(th *wid.Theme) layout.Widget {
	return wid.MakeList(
		th, wid.Occupy, nil,
		wid.Label(th, "Scrollable list of fields with labels", wid.Middle()),
		wid.Edit(th, wid.Lbl("Value 1")),
		wid.Edit(th, wid.Lbl("Value 2")),
		wid.Edit(th, wid.Lbl("Value 3")),
		wid.Edit(th, wid.Lbl("Value 4")),
		wid.Edit(th, wid.Lbl("Value 5")),
		wid.Edit(th, wid.Lbl("Value 6")),
		wid.Edit(th, wid.Lbl("Value 7")))
}

func column2(th *wid.Theme) layout.Widget {
	return wid.MakeList(th, wid.Occupy, nil,
		wid.Label(th, "Scrollable list of fields without labels", wid.Middle()),
		wid.Edit(th, wid.Hint("Value 1")),
		wid.Edit(th, wid.Hint("Value 2")),
		wid.Edit(th, wid.Hint("Value 3")),
		wid.Edit(th, wid.Hint("Value 4")),
		wid.Edit(th, wid.Hint("Value 5")),
		wid.Edit(th, wid.Hint("Value 6")),
		wid.Edit(th, wid.Hint("Value 7")))
}

func demo(th *wid.Theme) layout.Widget {
	thb = th
	y := startTime.Year()
	if y == 1 {
		startTime = time.Now()
	}
	return wid.Col(
		wid.Label(th, "Demo page", wid.Middle(), wid.Large(), wid.Bold()),
		wid.Row(th, nil, nil,
			wid.ProgressBar(th, &progress),
			wid.Value(th, func() string { return fmt.Sprintf(" %0.1f", count/time.Since(startTime).Seconds()) }),
		),
		wid.Separator(th, unit.Dp(2), wid.Color(th.SashColor)),
		wid.SplitVertical(th, 0.25,
			wid.SplitHorizontal(th, 0.5, column1(th), column2(th)),
			wid.MakeList(
				th, wid.Occupy, nil,
				wid.Col(
					wid.Row(th, nil, nil,
						wid.RadioButton(th, &mode, "windowed", "windowed"),
						wid.RadioButton(th, &mode, "fullscreen", "fullscreen"),
						wid.OutlineButton(th, "Maximize", wid.Handler(onMaximize)),
						wid.OutlineButton(th, "Center", wid.Handler(onCenter)),
					),
					wid.Row(th, nil, nil,
						wid.RadioButton(th, &fontSize, "small", "small"),
						wid.RadioButton(th, &fontSize, "medium", "medium"),
						wid.RadioButton(th, &fontSize, "large", "large"),
					),
					wid.Row(th, nil, nil,
						wid.Label(th, "A switch"),
						wid.Switch(th, &dummy, nil),
					),
					wid.Checkbox(th, "Checkbox to select dark mode", &darkMode, onSwitchMode),
					// Three separators to test layout algorithm. Should give three thin lines
					wid.Separator(th, unit.Px(5), wid.Color(wid.RGB(0xFF6666)), wid.Pads(5, 20, 5, 20)),
					wid.Separator(th, unit.Px(1)),
					wid.Separator(th, unit.Px(1), wid.Pads(1)),
					wid.Separator(th, unit.Px(1)),
					wid.Row(th, nil, []float32{0.3, 0.7},
						wid.Label(th, "A slider that can be key operated:"),
						wid.Slider(th, &sliderValue, 0, 100),
					),
					wid.Label(th, "A fixed width button at the middle of the screen:"),
					wid.Row(th, nil, nil,
						wid.Button(th, "WIDE CENTERED BUTTON",
							wid.W(500),
							wid.Hint("This is a dummy button - it has no function except displaying this text, testing long help texts, breaking it into several lines"),
						),
					),
					wid.Label(th, "Two widgets at the left side of the screen:"),
					wid.Row(th, nil, []float32{0.05, 0.9},
						wid.RoundButton(th, addIcon,
							wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
						wid.RoundButton(th, checkIcon,
							wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
					),
					// Note that buttons default to their minimum size, unless set differently, here aligned to the middle
					wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Color(wid.RGB(0xffff00))),
					wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Disable(&darkMode), wid.Color(wid.RGB(0x228822))),
					wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.W(150), wid.Color(wid.RGB(0xffff00))),
					wid.Button(thb, "Change color", wid.Handler(onClick), wid.W(150)),
					wid.TextButton(th, "Text button"),
					wid.OutlineButton(th, "Outline button"),
				),
				// Row with all buttons at minimum size, spread evenly
				wid.Row(th, nil, nil,
					wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Disable(&darkMode), wid.Min()),
					wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.Min()),
					wid.Button(thb, "Change color", wid.Handler(onClick), wid.Min()),
					wid.TextButton(th, "Text button", wid.Min()),
					wid.OutlineButton(th, "Outline button", wid.Min()),
				),
				// Fixed size in Dp
				wid.Edit(th, wid.Hint("Value 1"), wid.W(300)),
				// Relative size
				wid.Edit(th, wid.Hint("Value 2"), wid.W(0.5)),
				// The edit's default to their max size so they each get 1/5 of the row size. The MakeFlex spacing parameter will have no effect.
				wid.Row(th, nil, nil,
					wid.Edit(th, wid.Hint("Value 3")),
					wid.Edit(th, wid.Hint("Value 4")),
					wid.Edit(th, wid.Hint("Value 5")),
					wid.Edit(th, wid.Hint("Value 6")),
					wid.Edit(th, wid.Hint("Value 7")),
				),
				wid.Row(th, nil, nil,
					wid.Label(th, "Name", wid.End()),
					wid.Edit(th, wid.Hint("")),
				),
				wid.Row(th, nil, nil,
					wid.Label(th, "Address", wid.End()),
					wid.Edit(th, wid.Hint("")),
				),
				wid.Separator(th, unit.Dp(2.0)),
				wid.ImageFromJpgFile("gopher.jpg")),
		),
	)
}

func dropDownDemo(th *wid.Theme) layout.Widget {
	var longList []string
	for i := 1; i < 100; i++ {
		longList = append(longList, fmt.Sprintf("Option %d", i))
	}
	return wid.Pad(topRowPadding,
		wid.Col(
			wid.Row(th, nil, nil,
				wid.DropDown(th, 1, []string{"Option 1 with very long text", "Option 2", "Option 3"}),
				wid.DropDown(th, 2, []string{"Option 1", "Option 2", "Option 3"}),
				wid.DropDown(th, 0, []string{"Option A", "Option B", "Option C"}),
				wid.DropDown(th, 0, []string{"Option A", "Option B", "Option C"}),
			),
			// DropDown defaults to max size, here filling a complete row across the form.
			wid.DropDown(th, 0, []string{"Option X", "Option Y", "Option Z"}),
			wid.Separator(th, unit.Dp(2.0), wid.Pads(20, 0)),
			wid.Label(th, "A very long list with scrolling, with fixed width 250"),
			wid.DropDown(th, 0, longList, wid.W(250)),
		))
}

var page = "Buttons"

var topRowPadding = layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}

// Column widths are given in units of approximately one average character width (en).
var largeColWidth = []float32{2, 40, 40, 40, 40}
var smallColWidth = []float32{2, 20, 0.9, 6, 15}
var fracColWidth = []float32{2, 20.3, 0.3, 6, 0.14}

func setup() {
	th := currentTheme
	var currentPage layout.Widget
	if page == "Grid1" {
		currentPage = Grid(th, wid.Occupy, data, largeColWidth)
	} else if page == "Grid2" {
		currentPage = Grid(th, wid.Overlay, data, smallColWidth)
	} else if page == "Grid3" {
		currentPage = Grid(th, wid.Overlay, data[:5], fracColWidth)
	} else if page == "DropDown" {
		currentPage = dropDownDemo(th)
	} else if page == "Buttons" {
		currentPage = demo(th)
	} else if page == "KitchenV" {
		currentPage = kitchenV(th)
	}
	wid.Init()
	if page == "KitchenX" || page == "KitchenV" {
		wid.Setup(currentPage)
	} else {
		wid.Setup(wid.Col(
			wid.Pad(topRowPadding, wid.Row(th, nil, nil,
				wid.RadioButton(th, &page, "Grid1", "Grid1", wid.Do(update)),
				wid.RadioButton(th, &page, "Grid2", "Grid2", wid.Do(update)),
				wid.RadioButton(th, &page, "Grid3", "Grid3", wid.Do(update)),
				wid.RadioButton(th, &page, "Buttons", "Buttons", wid.Do(update)),
				wid.RadioButton(th, &page, "DropDown", "DropDowns", wid.Do(update)),
				wid.RadioButton(th, &page, "KitchenX", "KitchenX", wid.Do(update)),
				wid.RadioButton(th, &page, "KitchenV", "KitchenV", wid.Do(update)),
				wid.Checkbox(th, "Dark mode", &darkMode, onSwitchMode),
			)),
			wid.Separator(th, unit.Dp(2.0)),
			currentPage,
		))
	}
}

func handleFrameEvents(e system.FrameEvent) {
	if oldWindowSize.X != e.Size.X || oldWindowSize.Y != e.Size.Y || mode != oldMode || fontSize != oldFontSize {
		switch fontSize {
		case "medium", "Medium":
			currentTheme.TextSize = unit.Dp(float32(e.Size.Y) / 80)
		case "large", "Large":
			currentTheme.TextSize = unit.Dp(float32(e.Size.Y) / 60)
		case "small", "Small":
			currentTheme.TextSize = unit.Dp(float32(e.Size.Y) / 100)
		}
		oldFontSize = fontSize
		oldWindowSize = e.Size
		updateMode()
		setup()
	}
	var ops op.Ops
	gtx := layout.NewContext(&ops, e)
	// Set background color
	paint.Fill(gtx.Ops, currentTheme.Background)
	// Traverse the widget tree and generate drawing operations
	if page == "KitchenX" {
		kitchenX(gtx, th)
	} else {
		count++
		wid.Root(gtx)
	}
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

var prevAlloc uint64
var prevGc uint32

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage(txt string) {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("%s\tDeltaAlloc = %0.3f MiB", txt, (float64(m.Alloc)-float64(prevAlloc))/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC-prevGc)
	prevGc = m.NumGC
	prevAlloc = m.Alloc
}
