// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gioui.org/font"
	"gioui.org/font/gofont"
	"github.com/jkvatne/gio-v/wid"
	"image/color"
	"time"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	SmallFont      bool
	ScaleFont      bool
	theme          *wid.Theme  // the theme selected
	win            *app.Window // The main window
	form           layout.Widget
	name           = "Jan Kåre Vatne"
	age            = 35
	homeIcon       *wid.Icon
	checkIcon      *wid.Icon
	greenFlag              = false // the state variable for the button color
	dropDownValue1         = 1
	dropDownValue2         = 1
	dropDownValue3         = 1
	progress       float32 = 0.1
	sliderValue    float32 = 0.1
	WindowMode     string
	homeBg         = wid.RGB(0xF288F2)
	homeFg         = wid.RGB(0x0902200)
	list1          = []string{"Option 1 with very very very very very very very very very very very long text", "Option 2", "Option3"}
	list2          = []string{"Many options", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17"}
	list3          = []string{"Many options", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17"}
)

func main() {
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	homeIcon, _ = wid.NewIcon(icons.ActionHome)
	theme = wid.NewTheme(gofont.Collection(), 14)
	theme.DarkMode = false
	win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(700)))
	form = demo(theme)
	go wid.Run(win, &form, theme)
	go ticker()
	app.Main()
}

func ticker() {
	for {
		time.Sleep(time.Millisecond * 160)
		wid.GuiLock.Lock()
		progress = float32(int32((progress*1000)+5)%1000) / 1000.0
		wid.GuiLock.Unlock()
		wid.Invalidate()
	}
}

func onSwitchFontSize() {
	if ScaleFont && SmallFont {
		theme.SetLinesPrForm(60)
	} else if ScaleFont && !SmallFont {
		theme.SetLinesPrForm(40)
	} else if SmallFont {
		theme.SetLinesPrForm(0)
		theme.UpdateFontSize(11)
	} else {
		theme.SetLinesPrForm(0)
		theme.UpdateFontSize(14)
	}
}

func onSwitchMode() {
}

func onClick() {
	greenFlag = !greenFlag
	if greenFlag {
		theme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		theme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
	form = demo(theme)
}

func swColor() {
	if greenFlag {
		theme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x00, G: 0x9d, B: 0x00}
	} else {
		theme.PrimaryColor = color.NRGBA{A: 0xff, R: 0x10, G: 0x10, B: 0xff}
	}
}

func onWinChange() {
	switch WindowMode {
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
	ff := &font.Font{Typeface: "gomono"}
	return wid.Col(nil,
		wid.Label(th, "Demo page", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{.5, .5, .7, .5, .5},
			wid.Checkbox(th, "Dark mode", wid.Bool(&th.DarkMode), wid.Do(onSwitchMode)),
			wid.Checkbox(th, "Small font", wid.Bool(&SmallFont), wid.Do(onSwitchFontSize)),
			wid.Checkbox(th, "Scaleable", wid.Bool(&ScaleFont), wid.Do(onSwitchFontSize)),
			wid.RadioButton(th, &WindowMode, "windowed", "Windowed", wid.Do(onWinChange)),
			wid.RadioButton(th, &WindowMode, "fullscreen", "Fullscreen", wid.Do(onWinChange)),
			wid.RadioButton(th, &WindowMode, "maximized", "Maximized", wid.Do(onWinChange)),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Buttons with fixed length and large font, with and without icon"),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Change palette", wid.BtnIcon(checkIcon), wid.Do(onClick), wid.SecCont(), wid.W(450), wid.Large()),
		),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Button(th, "Change palette", wid.Do(onClick), wid.SecCont(), wid.W(450), wid.Large()),
			wid.Button(th, "Big Check", wid.BtnIcon(checkIcon), wid.FontSize(2), wid.Sec()),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Label(th, "Button spaced closely, left adjusted"),
		wid.Row(th, nil, wid.SpaceClose,
			wid.RoundButton(th, homeIcon, wid.Prim(),
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(homeIcon), wid.Bg(&homeBg), wid.Fg(&homeFg),
				wid.Hint("This is another hint")),
			wid.Button(th, "Check", wid.BtnIcon(checkIcon), wid.Role(wid.Secondary)),
			wid.Button(th, "Change color", wid.Do(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button", wid.Hint("An outlined button")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{1, 1, 1, 1, 1, 1},
			wid.Label(th, "Change color"),
			wid.Switch(th, &greenFlag, wid.Do(swColor)),
			wid.Label(th, " "),
			wid.Label(th, " "),
			wid.Label(th, "Tandem switch"),
			wid.Switch(th, &greenFlag, wid.Do(swColor)),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{1, 1, 1},
			wid.Edit(th, wid.Hint("Value 3"), wid.Var(&name)),
			wid.Edit(th, wid.Hint("Value 4"), wid.Var(&age)),
			wid.Edit(th, wid.Hint("Value 5"), wid.Var(&name)),
		),
		wid.Row(th, nil, []float32{1, 1, 1},
			wid.DropDown(th, &dropDownValue1, list1, wid.Hint("Value 3")),
			wid.DropDown(th, &dropDownValue2, list2, wid.Hint("Value 4")),
			wid.DropDown(th, &dropDownValue3, list3, wid.Hint("Value 5")),
		),
		wid.Row(th, nil, []float32{1, 1},
			wid.DropDown(th, &dropDownValue1, list1, wid.Lbl("Dropdown 1")),
			wid.DropDown(th, &dropDownValue2, list2, wid.Lbl("Dropdown 2")),
		),
		wid.Edit(th, wid.Lbl("Progress"), wid.Var(&progress), wid.Dp(4), wid.Ls(1/6.0)),
		wid.Slider(th, &sliderValue, 0, 100),
		wid.Row(th, nil, []float32{1, 1},
			wid.Col([]float32{},
				wid.Edit(th, wid.Hint("Hint 6"), wid.Lbl("Label 6"), wid.Ls(0.2)),
				wid.Edit(th, wid.Hint("Hint 7"), wid.Lbl("Label 7"), wid.Ls(0.2)),
			),
			wid.Col([]float32{},
				wid.Edit(th, wid.Lbl("Name"), wid.Var(&name), wid.Ls(0.5), wid.Font(ff)),
				wid.Edit(th, wid.Lbl("Age"), wid.Var(&age), wid.Ls(0.5)),
			),
		),
		wid.ProgressBar(th, &progress, wid.Pads(5.0), wid.W(12.0)),
		wid.Separator(th, 0, wid.Pads(5.0)),
		wid.ImageFromJpgFile("gopher.jpg", wid.Contain),
	)
}
