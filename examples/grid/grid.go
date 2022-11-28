// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly and horizontaly and implements highlighting of rows.

import (
	"gio-v/wid"
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
	Alternative = "Fractional"
	// Column widths are given in units of approximately one average character width (en).
	// A witdth of zero means the widget's natural size should be used (f.ex. checkboxes)
	wideColWidth  = []float32{0, 30, 30, 10, 20}
	smallColWidth = []float32{0, 13, 13, 12, 12}
	fracColWidth  = []float32{0, 0.3, 0.3, .2, .2}
	selectAll     bool
	doOccupy      bool
	nameIcon      *wid.Icon
	addressIcon   *wid.Icon
	ageIcon       *wid.Icon
	dir           bool
	// button        = new(widget.Clickable)
	// mth           = material.NewTheme(gofont.Collection())
	// test          = 1
)

type person struct {
	Selected bool
	Name     string
	Age      float64
	Address  string
	Status   int
}

var data = []person{
	{Name: "Oleg Karlsen", Age: 21.333333, Address: "Storgata 3", Status: 1},
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
	makePersons(12)
	theme = wid.NewTheme(gofont.Collection(), 24)
	onWinChange()
	go wid.Run(app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500))), &form, theme)
	app.Main()
}

func onWinChange() {
	var f layout.Widget
	if Alternative == "Wide" {
		f = Grid(theme, data, wideColWidth)
	} else if Alternative == "Narrow" {
		f = Grid(theme, data, smallColWidth)
	} else if Alternative == "Fractional" {
		f = Grid(theme, data, fracColWidth)
	} else if Alternative == "Equal" {
		f = Grid(theme, data, wid.SpaceDistribute)
	} else if Alternative == "Native" {
		f = Grid(theme, data, wid.SpaceClose)
	} else {
		f = Grid(theme, data, wid.SpaceDistribute)
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

// gw is the grid line width
const gw = 2.0 / 1.75

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, data []person, colWidths []float32) layout.Widget {
	anchor := wid.Overlay
	if doOccupy {
		anchor = wid.Occupy
	}
	bgColor := th.Bg(wid.Primary)

	nameIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)
	addressIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)
	ageIcon, _ = wid.NewIcon(icons.NavigationUnfoldMore)

	// Configure a grid with headings and several rows
	var gridLines []layout.Widget
	header := wid.GridRow(th, &bgColor, gw, colWidths,
		wid.Checkbox(th, "", wid.Bool(&selectAll), wid.Do(onCheck)),
		wid.HeaderButton(th, "Name", wid.Do(onNameClick), wid.Prim(), wid.BtnIcon(nameIcon), wid.Pads(0)),
		wid.HeaderButton(th, "Address", wid.Do(onAddressClick), wid.Prim(), wid.BtnIcon(addressIcon), wid.Pads(0)),
		wid.HeaderButton(th, "Age", wid.Do(onAgeClick), wid.Prim(), wid.BtnIcon(ageIcon), wid.Pads(0)),
		wid.Label(th, "Gender", wid.Prim(), wid.Pads(0)),
	)

	for i := 0; i < len(data); i++ {
		bgColor := wid.MulAlpha(th.Bg(wid.PrimaryContainer), 50)
		if i%2 == 0 {
			bgColor = wid.MulAlpha(th.Bg(wid.SecondaryContainer), 50)
		}
		gridLines = append(gridLines,
			wid.GridRow(th, &bgColor, gw, colWidths,
				wid.Checkbox(th, "", wid.Bool(&data[i].Selected)),
				wid.Label(th, &data[i].Name, wid.Pads(0)),
				wid.Edit(th, wid.Var(&data[i].Address), wid.Border(0), wid.Pads(0)),
				wid.Label(th, &data[i].Age, wid.Dp(2), wid.Right(), wid.Pads(0)),
				wid.DropDown(th, &data[i].Status, []string{"Male", "Female", "Other"}, wid.Pads(0), wid.Border(0)),
			))

	}
	var lines = []layout.Widget{
		wid.Label(th, "Grid demo", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Row(th, nil, nil,
			wid.RadioButton(th, &Alternative, "Wide", "Wide Table", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Narrow", "Narrow Table", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Fractional", "Fractional", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Equal", "Equal", wid.Do(onWinChange)),
			wid.RadioButton(th, &Alternative, "Native", "Native", wid.Do(onWinChange)),
		),
		wid.Row(th, nil, nil,
			wid.Checkbox(th, "Dark mode", wid.Bool(&th.DarkMode), wid.Do(onWinChange)),
			wid.Checkbox(th, "Scroll-bar occupy", wid.Bool(&doOccupy), wid.Do(onWinChange)),
			wid.Label(th, ""),
			// wid.RadioButton(th, &fontSize, "Large", "Large", wid.Do(onFontChange)),
			// wid.RadioButton(th, &fontSize, "Medium", "Medium", wid.Do(onFontChange)),
			// wid.RadioButton(th, &fontSize, "Small", "Small", wid.Do(onFontChange)),
		),
		wid.Edit(th, wid.Hint("Line editor")),
		wid.List(th, anchor, header, gridLines...),
		wid.Separator(th, 2),
		wid.Row(th, nil, []float32{1.0, 1.0, 1.0},
			wid.Space(1),
			wid.Button(th, "Update"),
			wid.Space(1),
		),
	}

	// return wid.List(th, wid.Occupy, f32.Point{1.0, 1.0}, lines...)
	return func(gtx wid.C) wid.D {
		bgColor := th.Bg(wid.Canvas)
		paint.Fill(gtx.Ops, bgColor)
		return wid.Col([]float32{0, 0, 0, 0, 1, 0}, lines...)(gtx)
		// return wid.List(th, anchor, f32.Pt(1, 1), lines...)(gtx)
	}
}
