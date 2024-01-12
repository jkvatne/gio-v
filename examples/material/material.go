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
	saveIcon     *wid.Icon
	otherPallete = false
	d            layout.Widget
)

func main() {
	homeIcon, _ = wid.NewIcon(icons.ActionHome)
	saveIcon, _ = wid.NewIcon(icons.ContentSave)
	theme = wid.NewTheme(gofont.Collection(), 16)
	d = wid.YesNoDialog(theme,
		"Save data?",
		"Click yes if you want to save the data to persitant memory or disk",
		"No", "Yes",
		onNo, onYes)
	theme.Scale = 1
	win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(500), unit.Dp(600)))
	form = demo(theme)
	go wid.Run(win, &form, theme)
	app.Main()
}

func onYes() {
	wid.Hide()
}

func onNo() {
	wid.Hide()
}

func myDialog(th *wid.Theme) layout.Widget {
	return wid.Container(th, wid.TransparentSurface, 0, wid.FlexInset, wid.NoInset,
		wid.Col(wid.SpaceDistribute,
			wid.Container(th, wid.PrimaryContainer, 20, layout.Inset{22, 22, 22, 22}, layout.Inset{62, 62, 62, 62},
				wid.Label(th, "Confirm", wid.Heading(), wid.Middle()),
				wid.Label(th, "Do you want to save data?", wid.Middle()),
				wid.Separator(th, 0, wid.Pads(10)),
				wid.Row(th, nil, wid.SpaceDistribute,
					wid.TextButton(th, "Yes", wid.Do(onYes), wid.Right(), wid.Margin(11)),
					wid.TextButton(th, "No", wid.Do(onNo), wid.Margin(11)),
				),
			),
		),
	)
}

func onSave() {
	wid.Show(d)
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
		wid.Row(th, nil, []float32{.5, .8, .5, .5},
			wid.Checkbox(th, "Dark mode", wid.Bool(&th.DarkMode), wid.Do(onSwitchMode), wid.Hint("Select light or dark mode")),
			wid.Checkbox(th, "Alt.pallete", wid.Bool(&otherPallete), wid.Do(onClick), wid.Hint("Select an alternative font")),
			wid.Checkbox(th, "Small font", wid.Bool(&SmallFont), wid.Do(onSwitchFontSize), wid.Hint("Select normal or small font size")),
			wid.Checkbox(th, "Fixed font", wid.Bool(&FixedFont), wid.Do(onSwitchFontSize), wid.Hint("Keep font size when resizing window height")),
		),
		wid.Separator(th, unit.Dp(1.0)),
		wid.Row(th, nil, []float32{0.3, 0.7},
			// Menu column
			wid.Container(th, wid.SurfaceContainerHigh, 15, th.DefaultPadding, th.DefaultMargin,
				wid.Col(wid.SpaceClose,
					wid.Label(th, "Items", wid.FontSize(1.5)),
					wid.TextButton(th, "Classic", wid.BtnIcon(homeIcon)),
					wid.TextButton(th, "Jazz", wid.BtnIcon(homeIcon)),
					wid.TextButton(th, "Rock", wid.BtnIcon(homeIcon)),
					wid.TextButton(th, "Hiphop", wid.BtnIcon(homeIcon)),
					wid.Space(9999),
				),
			),
			// Items
			wid.Col(wid.SpaceClose,
				wid.Container(th, wid.PrimaryContainer, 15, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Music", wid.FontSize(0.66), wid.Fg(th.PrimaryColor)),
					wid.Label(th, "What Buttons are Artists Pushing When They Perform Live", wid.FontSize(1.5), wid.PrimCont()),
					wid.Container(th, wid.PrimaryContainer, 15, layout.Inset{}, layout.Inset{0, 10, 0, 0},
						wid.ImageFromJpgFile("music.jpg", wid.Contain)),
					wid.Row(th, nil, wid.SpaceDistribute,
						wid.Label(th, "12 hrs ago", wid.FontSize(0.66), wid.Fg(th.PrimaryColor)),
						wid.Button(th, "Save", wid.Do(onSave), wid.BtnIcon(saveIcon), wid.RR(99), wid.Right()),
					),
				),
				wid.Container(th, wid.PrimaryContainer, 15, th.DefaultPadding, th.DefaultMargin,
					wid.Label(th, "Click Save button to test the confirmation dialog", wid.FontSize(1.0), wid.PrimCont()),
				),
			),
		),
	)
}
