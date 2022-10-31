// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly and horizontaly and implements highlighting of rows.

import (
	"fmt"
	"gio-v/wid"
	"os"
	"sort"

	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/paint"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	upIcon       *widget.Icon
	downIcon     *widget.Icon
	currentTheme *wid.Theme  // the theme selected
	win          *app.Window // The main window
	form         layout.Widget
	page         string
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

func main() {
	upIcon, _ = widget.NewIcon(icons.HardwareKeyboardArrowUp)
	downIcon, _ = widget.NewIcon(icons.HardwareKeyboardArrowDown)
	makePersons(100)

	go func() {
		currentTheme = wid.NewTheme(gofont.Collection(), 24)
		win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
		setup()
		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					os.Exit(0)
				case system.FrameEvent:
					handleFrameEvents(e)
				}
			}
		}
	}()
	app.Main()
}

func handleFrameEvents(e system.FrameEvent) {
	var ops op.Ops
	gtx := layout.NewContext(&ops, e)
	// Set background color
	c := currentTheme.Bg(wid.Canvas)
	paint.Fill(gtx.Ops, c)
	// A hack to fetch mouse position and window size so we can avoid
	// tooltips going outside the main window area
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	wid.UpdateMousePos(gtx, win, e.Size)
	// Draw widgets
	form(gtx)
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

// Column widths are given in units of approximately one average character width (en).
var largeColWidth = []float32{2, 40, 40, 10, 10}
var smallColWidth = []float32{2, 20, 0.9, 6, 15}
var fracColWidth = []float32{2, 20.3, 0.3, 6, 0.14}

func setup() {
	if page == "Grid1" {
		form = Grid(currentTheme, wid.Occupy, data, largeColWidth)
	} else if page == "Grid2" {
		form = Grid(currentTheme, wid.Overlay, data, smallColWidth)
	} else if page == "Grid3" {
		form = Grid(currentTheme, wid.Overlay, data[:5], fracColWidth)
	} else {
		form = Grid(currentTheme, wid.Occupy, data, smallColWidth)
	}
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
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
	}
	dir = !dir
	sortCol = 2
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age >= data[j].Age })
	}
	dir = !dir
	sortCol = 3
}

// selectAll is not used, but is controlled from the heading checkbox.
// It could be used to check/uncheck all boxes in the table
var selectAll bool
var nameIcon *widget.Icon
var addressIcon *widget.Icon
var ageIcon *widget.Icon

func getIcon(colNo int) *widget.Icon {
	if sortCol == colNo {
		if dir {
			return upIcon
		}
		return downIcon
	}
	return nil
}

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, anchor wid.AnchorStrategy, data []person, colWidths []float32) layout.Widget {
	nameIcon = getIcon(1)
	addressIcon = getIcon(2)
	ageIcon = getIcon(3)
	// Setup theme for heading.
	thh := *th
	// thh.OnBackground = wid.WithAlpha(th.Primary, 210)
	// thh.Background = th.Surface
	thh.LabelPadding = layout.UniformInset(unit.Dp(th.TextSize * 0.05))
	// Setup theme for grid labels.
	thg := *th
	// thg.Background = th.Surface
	thg.LabelPadding = layout.Inset{0, 0, 0, 0}
	thg.DropDownPadding = layout.Inset{0, 0, 0, 0}
	// Configure a row with headings.
	c := th.Bg(wid.Primary)
	heading := wid.Row(th, &c, &selectAll, colWidths,
		wid.Checkbox(th, "X", wid.Bool(&selectAll)),
		wid.HeaderButton(th, "Name", wid.Do(onNameClick), wid.W(9999), wid.BtnIcon(nameIcon)),
		wid.HeaderButton(th, "Address", wid.Do(onAddressClick), wid.W(9999), wid.BtnIcon(addressIcon)),
		wid.HeaderButton(th, "Age", wid.Do(onAgeClick), wid.W(9999), wid.BtnIcon(ageIcon)),
		// wid.Label(th, "Gender", wid.Bold()),
	)
	var lines []layout.Widget

	lines = append(lines, heading)
	lines = append(lines, wid.Separator(th, unit.Dp(2.0), wid.W(9999)))
	for i := 0; i < len(data); i++ {
		col := wid.MulAlpha(wid.Blue, 20)
		if i%2 == 0 {
			col = wid.MulAlpha(wid.Red, 20)
		}
		w := wid.Row(&thg, &col, &data[i].Selected, colWidths,
			wid.Checkbox(&thg, "", wid.Bool(&data[i].Selected)),
			wid.Label(&thg, "Åg"),
			wid.Label(&thg, data[i].Name),
			wid.Label(&thg, data[i].Address),
			wid.Label(&thg, fmt.Sprintf("%d", data[i].Age)),
			// wid.DropDown(&thg, &data[i].Status, []string{"Male", "Female", "Other"}),
		)
		lines = append(lines, w, wid.Separator(th, unit.Dp(0.7), wid.W(9999)))
	}
	return wid.List(&thg, anchor, lines...)
}

func onCheck(b bool) {
	// Called when the header checkbox is clicked. It will set or clear all rows.
	for i := 0; i < len(data); i++ {
		data[i].Selected = b
	}
}
