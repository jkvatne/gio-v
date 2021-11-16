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

// Make a lot of extra persons...
func makePersons() {
	for i := 1; i < 100; i++ {
		data[0].Age = i
		data = append(data, data[0])
	}
}

var dir bool

func onNameClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Name < data[j].Name })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Name >= data[j].Name })
	}
	dir = !dir
	update()
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
	}
	dir = !dir
	update()
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
	}
	dir = !dir
	update()
}

// selectAll is not used, but is controlled from the heading checkbox.
// It could be used to check/uncheck all boxes in the table
var selectAll bool

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, anchor wid.AnchorStrategy, data []person, colWidth []float32) layout.Widget {
	totalWidth := float32(0)
	for i := 0; i < len(colWidth); i++ {
		totalWidth += colWidth[i]
	}
	// Set up a new theme for the headings and rows
	thh := *th
	thg := *th
	thh.OnBackground = wid.WithAlpha(th.OnSurface, 192)
	thh.Background = th.Surface
	thg.Background = th.Surface
	heading := wid.Row(&thh, &selectAll, colWidth,
		wid.Checkbox(&thh, "", &selectAll, onCheck),
		wid.TextButton(&thh, "Name", wid.Handler(onNameClick)),
		wid.TextButton(&thh, "Address", wid.Handler(onAddressClick)),
		wid.TextButton(&thh, "Age", wid.Handler(onAgeClick)),
		wid.Label(&thh, "Gender", wid.Bold()),
	)
	var lines []layout.Widget
	lines = append(lines, heading)
	lines = append(lines, wid.Separator(th, unit.Dp(2.0), wid.W(totalWidth)))
	for i := 0; i < len(data); i++ {
		w := wid.Row(&thg, &data[i].Selected, colWidth,
			wid.Checkbox(&thg, "", &data[i].Selected, nil),
			wid.Label(&thg, data[i].Name),
			wid.Label(&thg, data[i].Address),
			wid.Label(&thg, fmt.Sprintf("%d", data[i].Age)),
			wid.DropDown(&thg, data[i].Status, []string{"Male", "Female", "Other"}),
		)
		lines = append(lines, w, wid.Separator(th, unit.Dp(0.5), wid.W(totalWidth)))
	}
	return wid.MakeList(&thg, anchor, int(totalWidth), lines...)
}

func onCheck(b bool) {
	// Called when the header checkbox is clicked. It will set or clear all rows.
	for i := 0; i < len(data); i++ {
		data[i].Selected = b
	}
}
