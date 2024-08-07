// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls vertically and horizontally and implements highlighting of rows.

import (
	"github.com/jkvatne/gio-v/wid"
	"sort"

	"gioui.org/op/paint"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	form        layout.Widget
	theme       *wid.Theme
	win         app.Window
	Alternative = "Fractional"
	fontSize    = "Medium"
	// Column widths are given in units of approximately one average character width (en).
	// A width of zero means the widget's natural size should be used (f.ex. checkboxes)
	wideColWidth  = []float32{0, 60, 60, 10, 30}
	smallColWidth = []float32{0, 13, 13, 12, 12}
	fracColWidth  = []float32{0, 0.3, 0.3, .2, .2}
	selectAll     bool
	doOccupy      bool
	withoutHeader bool = false
	nameIcon      *wid.Icon
	addressIcon   *wid.Icon
	ageIcon       *wid.Icon
	dir           bool
	line          string
)

type person struct {
	Selected bool
	Name     string
	Age      float64
	Address  string
	Status   int
}

var data = []person{
	{Name: "Oleg Karlsen", Age: 21, Address: "Storgata 1", Status: 1},
	{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 2", Selected: true, Status: 1},
	{Name: "Nils Aure", Age: 23, Address: "Brogata 3"},
	{Name: "Kai Oppdal", Age: 24, Address: "Soleieveien 4"},
	{Name: "Gro Arneberg", Age: 25, Address: "Blomsterveien 5"},
	{Name: "Ole Kolås", Age: 26, Address: "Blåklokkevikua 6"},
	{Name: "Per Pedersen", Age: 27, Address: "Gamleveien 7"},
	{Name: "Nils Vukubråten", Age: 28, Address: "Nygata 8"},
	{Name: "Sindre Gratangen", Age: 29, Address: "Brosundet 9"},
	{Name: "Gro Nilsasveen", Age: 30, Address: "Blomsterveien 10"},
	{Name: "Petter Olsen", Age: 31, Address: "Katavågen 11"},
	{Name: "Per Pedersen", Age: 32, Address: "Nidelva 12"},
}

func main() {
	makePersons(12)
	theme = wid.NewTheme(gofont.Collection(), 16)
	onWinChange()
	win.Option(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(300)))
	wid.Run(&win, &form, theme)
	app.Main()
}

func onWinChange() {
	var f layout.Widget
	theme.UpdateColors()
	if Alternative == "Wide" {
		f = GridDemo(theme, data, wideColWidth)
	} else if Alternative == "Narrow" {
		f = GridDemo(theme, data, smallColWidth)
	} else if Alternative == "Fractional" {
		f = GridDemo(theme, data, fracColWidth)
	} else if Alternative == "Equal" {
		f = GridDemo(theme, data, wid.SpaceDistribute)
	} else if Alternative == "Native" {
		f = GridDemo(theme, data, wid.SpaceClose)
	} else {
		f = GridDemo(theme, data, wid.SpaceDistribute)
	}
	wid.GuiLock.Lock()
	form = f
	defer wid.GuiLock.Unlock()
}

// makePersons will create a list of n persons.
func makePersons(n int) {
	m := n - len(data)
	for i := 1; i < m; i++ {
		data[0].Age = data[0].Age + float64(i)
		data = append(data, data[0])
	}
	data = data[:n]
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

func onFontChange() {
	if fontSize == "Medium" {
		theme.TextSize = 16
	} else if fontSize == "Large" {
		theme.TextSize = 26
	} else if fontSize == "Small" {
		theme.TextSize = 10
	}
	onWinChange()
}

// gw is the grid line width
const gw = unit.Dp(2.0 / 1.75)

// GridDemo is a widget that lays out the grid. This is all that is needed.
func GridDemo(th *wid.Theme, data []person, colWidths []float32) layout.Widget {
	anchor := wid.Overlay
	if doOccupy {
		anchor = wid.Occupy
	}
	bgColor := th.Bg[wid.PrimaryContainer]

	nameIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)
	addressIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)
	ageIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)

	// Configure a grid with headings and several rows
	var gridLines []layout.Widget
	header := wid.Row(th, &bgColor, gw, colWidths,
		wid.Checkbox(th, "", wid.Bool(&selectAll), wid.Do(onCheck)),
		wid.HeaderButton(th, "Name", wid.Do(onNameClick), wid.PrimCont(), wid.BtnIcon(nameIcon), wid.Pads(0)),
		wid.HeaderButton(th, "Address", wid.Do(onAddressClick), wid.PrimCont(), wid.BtnIcon(addressIcon), wid.Pads(0)),
		wid.HeaderButton(th, "Age", wid.Do(onAgeClick), wid.PrimCont(), wid.BtnIcon(ageIcon), wid.Pads(0)),
		// When using a label, padding has to be added. It should be equal to the default button padding.
		wid.Label(th, "Gender", wid.PrimCont()),
	)
	if withoutHeader {
		header = nil
	}

	for i := 0; i < len(data); i++ {
		bgColor := wid.MulAlpha(th.Bg[wid.PrimaryContainer], 50)
		if i%2 == 0 {
			bgColor = wid.MulAlpha(th.Bg[wid.SecondaryContainer], 50)
		}
		gridLines = append(gridLines,
			wid.Row(th, &bgColor, gw, colWidths,
				// One row of the grid is defined here, Name can not be edited
				wid.Checkbox(th, "", wid.Bool(&data[i].Selected)),
				wid.Label(th, &data[i].Name),
				wid.Edit(th, wid.Var(&data[i].Address), wid.Border(0), wid.Margin(0)),
				wid.Edit(th, wid.Var(&data[i].Age), wid.Border(0), wid.Margin(0)),
				wid.DropDown(th, &data[i].Status, []string{"Male", "Female", "Other"}, wid.Margin(0), wid.Border(0)),
			))

	}
	var lines = []layout.Widget{
		wid.Label(th, "GridDemo demo", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.RadioButton(th, &Alternative, "Wide", "Wide Table", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Narrow", "Narrow Table", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Fractional", "Fractional", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Equal", "Equal", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Native", "Native", wid.Do(onWinChange)),
		),
		wid.Row(th, nil, wid.SpaceDistribute,
			wid.Checkbox(th, "Dark mode", wid.Bool(&th.DarkMode), wid.Do(onWinChange)),
			wid.Checkbox(th, "Scroll-bar occupy", wid.Bool(&doOccupy), wid.Do(onWinChange)),
			wid.Checkbox(th, "No header", wid.Bool(&withoutHeader), wid.Do(onWinChange)),
			wid.Label(th, ""),
			wid.RadioButton(th, &fontSize, "Large", "Large", wid.Do(onFontChange)),
			wid.RadioButton(th, &fontSize, "Medium", "Medium", wid.Do(onFontChange)),
			wid.RadioButton(th, &fontSize, "Small", "Small", wid.Do(onFontChange)),
		),
		wid.Edit(th, &line, wid.Hint("Line editor")),
		wid.Table(th, anchor, header, gridLines...),
		wid.Separator(th, 2),
		// Center button that is <10 em wide. The width should be close to the native width, or the
		// button will not be centered.
		wid.Row(th, nil, []float32{1.0, 0.0, 1.0},
			wid.Space(1),
			wid.Button(th, "Update", wid.Hint("Click to update variables")),
			wid.Space(1),
		),
	}

	return func(gtx wid.C) wid.D {
		bgColor := th.Bg[wid.Surface]
		paint.Fill(gtx.Ops, bgColor)
		// Use flexible row heights. Set 1 for the grid, so it will use all available space.
		return wid.Col([]float32{0, 0, 0, 0, 1, 0, 0}, lines...)(gtx)
	}
}
