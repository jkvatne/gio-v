// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly only, but implements highlighting of rows.

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
var v = []float32{0, 24, 8, 30, 10}

// Grid is a widget that lays out the grid. This is all that is needed.
func grid(th *wid.Theme, data []person) layout.Widget {
	// Set up a new theme for the headings and rows
	thh := *th
	thg := *th
	thh.OnBackground = wid.WithAlpha(th.OnSurface, 192)
	thh.Background = th.Surface
	thg.Background = th.Surface
	heading := wid.Row(&thh, &selectAll, v,
		wid.Checkbox(&thh, "", &selectAll, onCheck),
		wid.Label(&thh, "Name", wid.Bold()),
		wid.Label(&thh, "Age", wid.Bold()),
		wid.Label(&thh, "Address", wid.Bold()),
		wid.Label(&thh, "Status"))

	//var lines = []layout.Widget{wid.Row(th, &selectAll, v, names...), wid.Separator(th, unit.Dp(0.5))}
	var lines []layout.Widget
	for i := 0; i < 2; i++ { //OBS len(data); i++ {
		w := wid.Row(&thg, &data[i].Selected, v,
			wid.Checkbox(&thg, "", &data[i].Selected, nil),
			wid.Label(&thg, data[i].Name),
			wid.Label(&thg, fmt.Sprintf("%d", data[i].Age)),
			wid.Label(&thg, data[i].Address),
			wid.DropDown(&thg, data[i].Status, []string{"Male", "Female", "Other"}),
		)
		lines = append(lines, w, wid.Separator(th, unit.Dp(0.5)))
	}
	grid := wid.MakeList(&thg, layout.Vertical, lines...)
	return wid.Col(heading, wid.Separator(th, unit.Dp(0.5)), grid)
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

func setupGridDemo(th *wid.Theme) {
	wid.Init()
	wid.Setup(wid.Col(
		wid.Row(th, nil, nil,
			wid.Checkbox(th, "Grid demo", &showGrid, onSwitchMode),
			wid.Checkbox(th, "Dark mode", &darkMode, onSwitchMode)),
		wid.Separator(th, unit.Dp(2.0)),
		InsetGrid(th, grid(th, data)),
	))
}

func onCheck(b bool) {
	// Called when a checkbox in a row is clicked. Not used yet.
	for i := 0; i < len(data); i++ {
		data[i].Selected = selectAll
	}
}
