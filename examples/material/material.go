// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"gioui.org/font/gofont"
	"github.com/jkvatne/gio-v/wid"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	SmallFont    bool
	FixedFont    bool
	theme        *wid.Theme  // the theme selected
	win          *app.Window // The main window
	form         layout.Widget
	homeIcon     *wid.Icon
	checkIcon    *wid.Icon
	otherPallete = false
)

func main() {
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	homeIcon, _ = wid.NewIcon(icons.ActionHome)
	theme = wid.NewTheme(gofont.Collection(), 24)
	theme.DarkMode = false
	win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(1200), unit.Dp(800)))
	form = demo(theme)
	go wid.Run(win, &form, theme)
	app.Main()
}

func onClick() {
	if otherPallete {
		theme.PrimaryColor = wid.RGB(0x17624E)
		theme.SecondaryColor = wid.RGB(0x17624E)
		theme.TertiaryColor = wid.RGB(0x136669)
		theme.ErrorColor = wid.RGB(0xAF2535)
		theme.NeutralColor = wid.RGB(0x1D4D7D)
		theme.NeutralVariantColor = wid.RGB(0x356057)
		theme.NeutralVariantColor = wid.RGB(0x356057)
	} else {
		// Set up the default pallete
		theme.PrimaryColor = wid.RGB(0x45682A)
		theme.SecondaryColor = wid.RGB(0x57624E)
		theme.TertiaryColor = wid.RGB(0x336669)
		theme.ErrorColor = wid.RGB(0xAF2525)
		theme.NeutralColor = wid.RGB(0x5D5D5D)
	}
	theme.UpdateColors()
}

func onSwitchFontSize() {
	if SmallFont {
		theme.TextSize = 11
	} else {
		theme.TextSize = 20
	}
	wid.FixedFontSize = FixedFont
}

func onSwitchMode() {
	theme.UpdateColors()
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.Col(wid.SpaceClose,
		wid.Label(th, "Material demo", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{.5, .5, .8, .5, .5, .5},
			wid.Checkbox(th, "Dark mode", wid.Bool(&th.DarkMode), wid.Do(onSwitchMode)),
			wid.Checkbox(th, "Small font", wid.Bool(&SmallFont), wid.Do(onSwitchFontSize)),
			wid.Checkbox(th, "Fixed size", wid.Bool(&FixedFont), wid.Do(onSwitchFontSize)),
			wid.Checkbox(th, "Alt.pallete", wid.Bool(&otherPallete), wid.Do(onClick)),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{0.3, 0.7},
			// Menu column
			wid.Container(th, wid.SurfaceContainerHigh, 15,
				wid.Col(wid.SpaceClose,
					wid.Label(th, "Items", wid.FontSize(1.5)),
					wid.TextButton(th, "Freelance", wid.BtnIcon(homeIcon)),
					wid.Button(th, "Contracts", wid.BtnIcon(homeIcon)),
					wid.Space(9999),
				),
			),
			// Items
			wid.Col(wid.SpaceClose,
				wid.Container(th, wid.PrimaryContainer, 10,
					wid.Label(th, "Folders", wid.FontSize(1.0), wid.Role(wid.PrimaryContainer)),
					wid.Label(th, "Files", wid.FontSize(1.0), wid.Role(wid.PrimaryContainer)),
				),
				wid.Container(th, wid.PrimaryContainer, 10,
					wid.Label(th, "Folders", wid.FontSize(1.0), wid.Role(wid.PrimaryContainer)),
					wid.Label(th, "Files", wid.FontSize(1.0), wid.Role(wid.PrimaryContainer)),
				),
			),
		),
	)
}
