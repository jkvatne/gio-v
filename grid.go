// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly and horizontaly and implements highlighting of rows.

import (
	"fmt"
	"gio-v/wid"
	"sort"

	"gioui.org/layout"
	"gioui.org/unit"
)

type person struct {
	Selected bool
	Name     string
	Age      int
	Address  string
	Status   int
}

var data = []person{
	{Name: "Ole Karlsen", Age: 21, Address: "Storgata 3", Status: 1},
	{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 33", Selected: true, Status: 1},
	{Name: "Nils Aure", Age: 23, Address: "Brogata 34"},
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

// Make list of n persons.
func makePersons(n int) {
	m := n - len(data)
	for i := 1; i < m; i++ {
		data[0].Age = i
		data = append(data, data[0])
	}
	data = data[0:n]
}

var dir bool
var sortCol int

func onNameClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Name < data[j].Name })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Name >= data[j].Name })
	}
	dir = !dir
	sortCol = 1
	update()
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
	}
	dir = !dir
	sortCol = 2
	update()
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age >= data[j].Age })
	}
	dir = !dir
	sortCol = 3
	update()
}

// selectAll is not used, but is controlled from the heading checkbox.
// It could be used to check/uncheck all boxes in the table
var selectAll bool
var nameIcon *wid.Icon
var addressIcon *wid.Icon
var ageIcon *wid.Icon

func getIcon(colNo int) *wid.Icon {
	if sortCol == colNo {
		if dir {
			return upIcon
		}
		return downIcon
	}
	return nil
}

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, anchor wid.AnchorStrategy, data []person, colWidth []float32) layout.Widget {
	nameIcon = getIcon(1)
	addressIcon = getIcon(2)
	ageIcon = getIcon(3)
	// Setup theme for heading.
	thh := *th
	thh.OnBackground = wid.WithAlpha(th.Primary, 210)
	thh.Background = th.Surface
	thh.LabelPadding = layout.UniformInset(th.TextSize.Scale(0.4))
	// Setup theme for grid labels.
	thg := *th
	thg.Background = th.Surface
	thg.LabelPadding = layout.UniformInset(th.TextSize.Scale(0.35))
	// Configure a row with headings.
	heading := wid.Row(&thh, &selectAll, colWidth,
		wid.Checkbox(&thh, "", &selectAll, onCheck),
		wid.HeaderButton(&thh, "Name", wid.Handler(onNameClick), wid.W(9999), wid.BtnIcon(nameIcon)),
		wid.HeaderButton(&thh, "Address", wid.Handler(onAddressClick), wid.W(9999), wid.BtnIcon(addressIcon)),
		wid.HeaderButton(&thh, "Age", wid.Handler(onAgeClick), wid.W(9999), wid.BtnIcon(ageIcon)),
		wid.Label(&thh, "Gender", wid.Bold()),
	)
	var lines []layout.Widget
	lines = append(lines, heading)
	lines = append(lines, wid.Separator(th, unit.Dp(2.0), wid.W(9999)))
	for i := 0; i < len(data); i++ {
		w := wid.Row(&thg, &data[i].Selected, colWidth,
			wid.Checkbox(&thg, "", &data[i].Selected, nil),
			wid.Label(&thg, data[i].Name),
			wid.Label(&thg, data[i].Address),
			wid.Label(&thg, fmt.Sprintf("%d", data[i].Age)),
			wid.DropDown(&thg, data[i].Status, []string{"Male", "Female", "Other"}),
		)
		lines = append(lines, w, wid.Separator(th, unit.Dp(0.5), wid.W(9999)))
	}
	return wid.MakeList(&thg, anchor, colWidth, lines...)
}

func onCheck(b bool) {
	// Called when the header checkbox is clicked. It will set or clear all rows.
	for i := 0; i < len(data); i++ {
		data[i].Selected = b
	}
}
