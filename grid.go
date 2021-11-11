// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly only, but implements highlighting of rows.

import (
	"fmt"
	"gio-v/wid"

	"gioui.org/layout"
	"gioui.org/unit"
)

type person struct {
	Selected bool
	Name     string
	Age      int
	Address  string
}

var data = []person{
	{Name: "Ole", Age: 21, Address: "Storgata 3"},
	{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 33", Selected: true},
	{Name: "Nils", Age: 23, Address: "Brogata 34"},
	{Name: "Kai", Age: 28, Address: "Soleieveien 12"},
	{Name: "Gro", Age: 29, Address: "Blomsterveien 22"},
	{Name: "Ole", Age: 21, Address: "Blåklokkevikua 33"},
	{Name: "Per Pedersen", Age: 22, Address: "Gamleveien 35"},
	{Name: "Nils", Age: 23, Address: "Nygata 64"},
	{Name: "Sindre Gratangen", Age: 28, Address: "Brosundet 34"},
	{Name: "Gro", Age: 29, Address: "Blomsterveien 22"},
	{Name: "Petter Olsen", Age: 21, Address: "Katavågen 44"},
	{Name: "Per Pedersen", Age: 22, Address: "Nidelva 43"},
}

// Make a lot of persons...
func makePersons() {
	for i := 1; i < 1000; i++ {
		data[0].Age = i
		data = append(data, data[0])
	}
}

// selectAll is not used, but is controlled from the heading checkbox.
// It could be used to check/uncheck all boxes in the table
var selectAll bool

// v is the relative width of each column. Use like Flexed weight.
var v = []float32{0.1, 1.0, 0.3, 2.0}

// Grid is a widget that lays out the grid. This is all that is needed.
func grid(th *wid.Theme, data []person) layout.Widget {
	// Set up a new theme for the headings
	thh := *th
	thh.TextSize = th.TextSize.Scale(0.85)
	thh.OnBackground = wid.WithAlpha(th.OnSurface, 192)

	heading := wid.Row(th, &selectAll, v,
		wid.Checkbox(&thh, " ", &selectAll, nil),
		wid.Label(&thh, "Name", wid.Bold()),
		wid.Label(&thh, "Age", wid.Bold()),
		wid.Label(&thh, "Address", wid.Bold()))

	//var lines = []layout.Widget{wid.Row(th, &selectAll, v, names...), wid.Separator(th, unit.Dp(0.5))}
	var lines []layout.Widget
	for i := 0; i < len(data); i++ {
		w := wid.Row(th, &data[i].Selected, v,
			wid.Checkbox(th, " ", &data[i].Selected, onCheck),
			wid.Label(th, data[i].Name),
			wid.Label(th, fmt.Sprintf("%d", data[i].Age)),
			wid.Label(th, data[i].Address),
		)
		lines = append(lines, w, wid.Separator(th, unit.Dp(0.5)))
	}
	grid := wid.MakeList(th, layout.Vertical, lines...)
	return wid.Col(heading, wid.Separator(th, unit.Dp(0.5)), grid)
}

func setupGridDemo(th *wid.Theme) {
	// thb is the theme for highlighted rows.
	thb = th
	wid.Init()
	wid.Setup(wid.Col(
		wid.Row(th, nil, nil,
			wid.Checkbox(th, "Grid demo", &showGrid, onSwitchMode),
			wid.Checkbox(th, "Dark mode", &darkMode, onSwitchMode)),
		wid.Separator(th, unit.Dp(2.0)),
		grid(th, data),
	))
}

func onCheck(b bool) {
	// Called when a checkbox in a row is clicked. Not used yet.
}
