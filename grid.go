// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly and horizontaly and implements highlighting of rows.

import (
	"fmt"
	"gio-v/wid"

	"gioui.org/op/clip"
	"gioui.org/op/paint"

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
	{Name: "Ole", Age: 21, Address: "Storgata 3", Status: 1},
	{Name: "Per Pedersen", Age: 22, Address: "Svenskveien 33", Selected: true, Status: 1},
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

// Make a lot of extra persons...
func makePersons() {
	for i := 1; i < 100; i++ {
		data[0].Age = i
		data = append(data, data[0])
	}
}

// selectAll is not used, but is controlled from the heading checkbox.
// It could be used to check/uncheck all boxes in the table
var selectAll bool

var colWidth = []float32{35, 300, 300, 300, 300}

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, data []person) layout.Widget {
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
		wid.Label(&thh, "Name", wid.Bold()),
		wid.Label(&thh, "Address", wid.Bold()),
		wid.Label(&thh, "Age", wid.Bold()),
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
	return wid.MakeList(&thg, wid.Occupy, int(totalWidth), lines...)
}

// InsetGrid is the grid with some padding
func InsetGrid(th *wid.Theme, grid layout.Widget) layout.Widget {
	outerPadding := layout.UniformInset(th.TextSize.Scale(1.0))
	innerPadding := layout.UniformInset(unit.Dp(2))
	return func(gtx wid.C) wid.D {
		//Outer padding. Drawn with background color.
		c := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
		paint.ColorOp{Color: th.Background}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		c.Pop()
		return outerPadding.Layout(gtx, func(gtx wid.C) wid.D {
			// Inner padding.Drawn with surface background
			c := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
			paint.ColorOp{Color: th.Surface}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			c.Pop()
			return innerPadding.Layout(gtx, grid)
		})
	}
}

func onCheck(b bool) {
	// Called when the header checkbox is clicked. It will set or clear all rows.
	for i := 0; i < len(data); i++ {
		data[i].Selected = selectAll
	}
}
