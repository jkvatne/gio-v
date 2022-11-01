// SPDX-License-Identifier: Unlicense OR MIT

package main

// This file demonstrates a simple grid, trying to follow https://material.io/components/data-tables
// It scrolls verticaly and horizontaly and implements highlighting of rows.

import (
	"gio-v/wid"
	"os"
	"sort"

	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/paint"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	// upIcon       *wid.Icon
	// downIcon     *wid.Icon
	// sortIcon     *wid.Icon
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
	// upIcon, _ = wid.NewIcon(icons.HardwareKeyboardArrowUp)
	// downIcon, _ = wid.NewIcon(icons.HardwareKeyboardArrowDown)
	// sortIcon, _ = wid.NewIcon(icons.ContentSort)
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
var largeColWidth = []float32{20, 40, 40, 10, 10}
var smallColWidth = []float32{20, 20, 0.9, 6, 15}
var fracColWidth = []float32{20, 20.3, 0.3, 6, 0.14}

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
		nameIcon.Update(icons.HardwareKeyboardArrowDown)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Name >= data[j].Name })
		nameIcon.Update(icons.HardwareKeyboardArrowUp)
	}
	addressIcon.Update(icons.ContentSort)
	ageIcon.Update(icons.ContentSort)
	dir = !dir
	sortCol = 1
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
		addressIcon.Update(icons.HardwareKeyboardArrowDown)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
		addressIcon.Update(icons.HardwareKeyboardArrowUp)
	}
	nameIcon.Update(icons.ContentSort)
	ageIcon.Update(icons.ContentSort)
	dir = !dir
	sortCol = 2
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
		ageIcon.Update(icons.HardwareKeyboardArrowDown)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age >= data[j].Age })
		ageIcon.Update(icons.HardwareKeyboardArrowUp)
	}
	nameIcon.Update(icons.ContentSort)
	addressIcon.Update(icons.ContentSort)
	dir = !dir
	sortCol = 3
}

// selectAll is not used, but is controlled from the heading checkbox.
// It could be used to check/uncheck all boxes in the table
var selectAll bool
var nameIcon *wid.Icon
var addressIcon *wid.Icon
var ageIcon *wid.Icon

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, anchor wid.AnchorStrategy, data []person, colWidths []float32) layout.Widget {
	nameIcon, _ = wid.NewIcon(icons.ContentSort)
	addressIcon, _ = wid.NewIcon(icons.ContentSort)
	ageIcon, _ = wid.NewIcon(icons.ContentSort)
	// Setup theme for heading.
	thh := *th
	thg := *th
	// Configure a row with headings.
	bgColor := th.Bg(wid.Primary)
	heading := wid.Row(&thh, &bgColor, colWidths,
		wid.Checkbox(th, "", wid.Bool(&selectAll), wid.Role(wid.Primary)),
		wid.HeaderButton(&thh, "Name", wid.Do(onNameClick), wid.BtnIcon(nameIcon), wid.P()),
		wid.HeaderButton(&thh, "Address", wid.Do(onAddressClick), wid.BtnIcon(addressIcon), wid.P()),
		wid.HeaderButton(&thh, "Age", wid.Do(onAgeClick), wid.BtnIcon(ageIcon), wid.P()),
		// wid.Label(th, "Gender", wid.Bold()),
	)
	var lines []layout.Widget

	lines = append(lines, heading)
	lines = append(lines, wid.Separator(th, unit.Dp(2.0), wid.W(9999)))
	for i := 0; i < len(data); i++ {
		bgColor := wid.MulAlpha(wid.Blue, 20)
		if i%2 == 0 {
			bgColor = wid.MulAlpha(wid.Red, 20)
		}
		w := wid.Row(&thg, &bgColor, colWidths,
			wid.Checkbox(&thg, "", wid.Bool(&data[i].Selected)),
			wid.Label(&thg, &data[i].Name),
			wid.Label(&thg, &data[i].Address),
			// wid.Label(&thg, fmt.Sprintf("%d", data[i].Age)),
			wid.DropDown(&thg, &data[i].Status, []string{"Male", "Female", "Other"}),
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
