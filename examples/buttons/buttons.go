// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/jkvatne/gio-v/wid"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image/color"
)

var (
	theme     *wid.Theme // the theme selected
	form      layout.Widget
	homeIcon  *wid.Icon
	checkIcon *wid.Icon
	greenFlag = false // the state variable for the button color
	homeBg    = wid.RGB(0xF288F2)
	homeFg    = wid.RGB(0x0902200)
)

func main() {
	homeIcon, _ = wid.NewIcon(icons.ActionHome)
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	theme = wid.NewTheme(gofont.Collection(), 14)
	onClick()
	go wid.Run(app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500))), &form, theme)
	app.Main()
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

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget -- i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Overlay,
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

		wid.Label(th, "Buttons with relative spacing given by weights 0.2, 0.4, 0.4, 0.4, 0.4, 0.4"),
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
