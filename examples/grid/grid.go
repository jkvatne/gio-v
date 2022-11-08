// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly and horizontaly and implements highlighting of rows.

import (
	"gio-v/wid"
	"sort"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/layout"
	"gioui.org/unit"
)

const test = 0

var (
	form        layout.Widget
	theme       wid.Theme
	Alternative = "Wide"
	// Column widths are given in units of approximately one average character width (en).
	// A witdth of zero means the widget's natural size should be used (f.ex. checkboxes)
	wideColWidth  = []float32{0, 40, 40, 20, 20}
	smallColWidth = []float32{0, 13, 13, 12, 12}
	fracColWidth  = []float32{0, 0.3, 0.3, .2, .2}
	selectAll     bool
	nameIcon      *wid.Icon
	addressIcon   *wid.Icon
	ageIcon       *wid.Icon
	dir           bool
	fontSize      = "Large"
)

type person struct {
	Selected bool
	Name     string
	Age      float64
	Address  string
	Status   int
}

var data = []person{
	{Name: "Ole Karlsen", Age: 21.333333, Address: "Storgata 3", Status: 1},
	{Name: "Per Pedersen", Age: 22.111111111, Address: "Svenskveien 33", Selected: true, Status: 1},
	{Name: "Nils Aure", Age: 23.4, Address: "Brogata 34"},
	{Name: "Kai Oppdal", Age: 28, Address: "Soleieveien 12"},
	{Name: "Gro Arneberg", Age: 29, Address: "Blomsterveien 22"},
	{Name: "Ole Kolås", Age: 21, Address: "Blåklokkevikua 33"},
	{Name: "Per Pedersen", Age: 22, Address: "Gamleveien 35"},
	{Name: "Nils Vukubråten", Age: 23, Address: "Nygata 64"},
	{Name: "Sindre Gratangen", Age: 28, Address: "Brosundet 34"},
	{Name: "Gro Nilsasveen", Age: 29, Address: "Blomsterveien 22"},
	{Name: "Petter Olsen", Age: 21, Address: "Katavågen 44"},
	{Name: "Per Pedersen", Age: 22, Address: "Nidelva 43"},
}

func main() {
	makePersons(20)
	theme = *wid.NewTheme(gofont.Collection(), 24)
	onWinChange()
	go wid.Run(app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500))), &form)
	app.Main()
}

func onWinChange() {
	var f layout.Widget
	if Alternative == "Wide" {
		f = Grid(&theme, wid.Overlay, data, wideColWidth)
	} else if Alternative == "Narrow" {
		f = Grid(&theme, wid.Overlay, data[:5], smallColWidth)
	} else if Alternative == "Fractional" {
		f = Grid(&theme, wid.Overlay, data, fracColWidth)
	} else if Alternative == "Equal" {
		f = Grid(&theme, wid.Occupy, data, wid.SpaceDistribute)
	} else if Alternative == "Native" {
		f = Grid(&theme, wid.Occupy, data, wid.SpaceClose)
	} else {
		f = Grid(&theme, wid.Occupy, data, wid.SpaceDistribute)
	}
	wid.GuiLock.Lock()
	form = f
	defer wid.GuiLock.Unlock()
}

func onFontChange() {
	if fontSize == "Large" {
		theme = *wid.NewTheme(gofont.Collection(), 24)
	} else if fontSize == "Small" {
		theme = *wid.NewTheme(gofont.Collection(), 10)
	} else if fontSize == "Medium" {
		theme = *wid.NewTheme(gofont.Collection(), 14)
	}
	onWinChange()
}

// makePersons will create a list of n persons.
func makePersons(n int) {
	m := n - len(data)
	for i := 1; i < m; i++ {
		data[0].Age = data[0].Age + float64(i)
		data = append(data, data[0])
	}
}

func onNameClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Name >= data[j].Name })
		_ = nameIcon.Update(icons.NavigationArrowDownward)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Name < data[j].Name })
		_ = nameIcon.Update(icons.NavigationArrowUpward)
	}
	_ = addressIcon.Update(icons.NavigationUnfoldMore)
	_ = ageIcon.Update(icons.NavigationUnfoldMore)
	dir = !dir
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
		_ = addressIcon.Update(icons.NavigationArrowDownward)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
		_ = addressIcon.Update(icons.NavigationArrowUpward)
	}
	_ = nameIcon.Update(icons.NavigationUnfoldMore)
	_ = ageIcon.Update(icons.NavigationUnfoldMore)
	dir = !dir
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age >= data[j].Age })
		_ = ageIcon.Update(icons.NavigationArrowDownward)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
		_ = ageIcon.Update(icons.NavigationArrowUpward)
	}
	_ = nameIcon.Update(icons.NavigationUnfoldMore)
	_ = addressIcon.Update(icons.NavigationUnfoldMore)
	dir = !dir
}

// onCheck is called when the header checkbox is clicked. It will set or clear all rows.
func onCheck() {
	for i := 0; i < len(data); i++ {
		data[i].Selected = selectAll
	}
}

const gw = 2.0 / 1.75

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, anchor wid.AnchorStrategy, data []person, colWidths []float32) layout.Widget {
	// Set background color according to theme
	if test == 1 {
		return wid.GridRow(th, nil, gw, []float32{0, 0.9},
			wid.Checkbox(th, "", wid.Bool(&data[1].Selected)),
			wid.Label(th, &data[1].Name))
	} else if test == 2 {
		return wid.GridRow(th, nil, gw, []float32{0, 40, 40, 20, 20},
			wid.Checkbox(th, "", wid.Bool(&data[1].Selected)),
			wid.Label(th, &data[1].Address),
			wid.Label(th, &data[1].Name),
			wid.Label(th, &data[1].Age, wid.Dp(3), wid.Right()))

	} else if test == 3 {
		return wid.RadioButton(th, &Alternative, "WideColumns", "WideColumns", wid.Do(onWinChange))
	} else {
		nameIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)
		addressIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)
		ageIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)

		var lines []layout.Widget
		lines = append(lines,
			wid.Label(th, "Grid demo", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
			wid.Label(th, "Different wighting and size of columns"),
			wid.Row(th, nil, nil,
				wid.RadioButton(th, &Alternative, "Wide", "Wide", wid.Do(onWinChange)),
				wid.RadioButton(th, &Alternative, "Narrow", "Narrow", wid.Do(onWinChange)),
				wid.RadioButton(th, &Alternative, "Fractional", "Fractional", wid.Do(onWinChange)),
				wid.RadioButton(th, &Alternative, "Equal", "Equal", wid.Do(onWinChange)),
				wid.RadioButton(th, &Alternative, "Native", "Native", wid.Do(onWinChange)),
			),
			wid.Space(5),
			wid.Label(th, "Select font size"),
			wid.Row(th, nil, nil,
				wid.Checkbox(th, "Dark mode", wid.Bool(&th.DarkMode), wid.Do(onWinChange)),
				wid.Label(th, ""),
				wid.RadioButton(th, &fontSize, "Large", "Large", wid.Do(onFontChange)),
				wid.RadioButton(th, &fontSize, "Medium", "Medium", wid.Do(onFontChange)),
				wid.RadioButton(th, &fontSize, "Small", "Small", wid.Do(onFontChange)),
			),
			wid.Space(20),
		)
		// Configure a row with headings.
		bgColor := th.Bg(wid.PrimaryContainer)
		heading := wid.GridRow(th, &bgColor, gw, colWidths,
			wid.Checkbox(th, "", wid.Bool(&selectAll), wid.Do(onCheck)),
			wid.HeaderButton(th, "Name", wid.Do(onNameClick), wid.Prim(), wid.BtnIcon(nameIcon)),
			wid.HeaderButton(th, "Address", wid.Do(onAddressClick), wid.Prim(), wid.BtnIcon(addressIcon)),
			wid.HeaderButton(th, "Age", wid.Do(onAgeClick), wid.Prim(), wid.BtnIcon(ageIcon)),
			wid.Label(th, "Gender", wid.Bold()),
		)
		lines = append(lines, heading)

		for i := 0; i < len(data); i++ {
			bgColor := wid.MulAlpha(th.Bg(wid.PrimaryContainer), 50)
			if i%2 == 0 {
				bgColor = wid.MulAlpha(th.Bg(wid.SecondaryContainer), 50)
			}
			lines = append(lines,
				wid.GridRow(th, &bgColor, gw, colWidths,
					wid.Checkbox(th, "", wid.Bool(&data[i].Selected)),
					wid.Label(th, &data[i].Name),
					wid.Label(th, &data[i].Address),
					wid.Label(th, &data[i].Age, wid.Dp(2), wid.Right()),
					wid.DropDown(th, &data[i].Status, []string{"Male", "Female", "Other"}, wid.Border(0)),
				))

		}
		return wid.List(th, anchor, lines...)
	}
}
