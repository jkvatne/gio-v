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

const test = false

var (
	currentTheme *wid.Theme  // the theme selected
	win          *app.Window // The main window
	form         layout.Widget
	Alternative  string
	// Column widths are given in units of approximately one average character width (en).
	largeColWidth = []float32{0, 40, 40, 20}
	smallColWidth = []float32{0, 9, 0.9, 12}
	fracColWidth  = []float32{0, 0.3, 0.3, .2}
	selectAll     bool
	nameIcon      *wid.Icon
	addressIcon   *wid.Icon
	ageIcon       *wid.Icon
	dir           bool
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
	makePersons(100)
	go func() {
		currentTheme = wid.NewTheme(gofont.Collection(), 24)
		win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
		onWinChange()
		for {
			e := <-win.Events()
			switch e := e.(type) {
			case system.DestroyEvent:
				os.Exit(0)
			case system.FrameEvent:
				handleFrameEvents(e)
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

func onWinChange() {
	if Alternative == "LargeColumns" {
		form = Grid(currentTheme, wid.Occupy, data, largeColWidth)
	} else if Alternative == "SmallColumns" {
		form = Grid(currentTheme, wid.Overlay, data, smallColWidth)
	} else if Alternative == "FractionalColumns" {
		form = Grid(currentTheme, wid.Overlay, data[:5], fracColWidth)
	} else {
		form = Grid(currentTheme, wid.Occupy, data, smallColWidth)
	}
}

// makePersons will create a list of n persons.
func makePersons(n int) {
	m := n - len(data)
	for i := 1; i < m; i++ {
		data[0].Age = data[0].Age + float64(i)
		data = append(data, data[0])
	}
	data = data[0:n]
}

func onNameClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Name >= data[j].Name })
		_ = nameIcon.Update(icons.HardwareKeyboardArrowDown)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Name < data[j].Name })
		_ = nameIcon.Update(icons.HardwareKeyboardArrowUp)
	}
	_ = addressIcon.Update(icons.ContentSort)
	_ = ageIcon.Update(icons.ContentSort)
	dir = !dir
}

func onAddressClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Address >= data[j].Address })
		_ = addressIcon.Update(icons.HardwareKeyboardArrowDown)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Address < data[j].Address })
		_ = addressIcon.Update(icons.HardwareKeyboardArrowUp)
	}
	_ = nameIcon.Update(icons.ContentSort)
	_ = ageIcon.Update(icons.ContentSort)
	dir = !dir
}

func onAgeClick() {
	if dir {
		sort.Slice(data, func(i, j int) bool { return data[i].Age >= data[j].Age })
		_ = ageIcon.Update(icons.HardwareKeyboardArrowDown)
	} else {
		sort.Slice(data, func(i, j int) bool { return data[i].Age < data[j].Age })
		_ = ageIcon.Update(icons.HardwareKeyboardArrowUp)
	}
	_ = nameIcon.Update(icons.ContentSort)
	_ = addressIcon.Update(icons.ContentSort)
	dir = !dir
}

// onCheck is called when the header checkbox is clicked. It will set or clear all rows.
func onCheck() {
	for i := 0; i < len(data); i++ {
		data[i].Selected = selectAll
	}
}

// Grid is a widget that lays out the grid. This is all that is needed.
func Grid(th *wid.Theme, anchor wid.AnchorStrategy, data []person, colWidths []float32) layout.Widget {
	if test {
		return wid.RadioButton(th, &Alternative, "LargeColumns", "LargeColumns", wid.Do(onWinChange))
	} else {
		nameIcon, _ = wid.NewIcon(icons.ContentSort)
		addressIcon, _ = wid.NewIcon(icons.ContentSort)
		ageIcon, _ = wid.NewIcon(icons.ContentSort)

		var lines []layout.Widget

		lines = append(lines,
			wid.GridRow(th, nil, 1.5, []float32{0, 9, 0.9, 8},
				wid.Checkbox(th, "", wid.Bool(&data[1].Selected)),
				wid.Label(th, &data[1].Address),
				wid.Label(th, &data[1].Name),
				wid.Label(th, &data[1].Age, wid.Dp(3), wid.Right())),
			// wid.DropDown(&thg, &data[i].Status, []string{"Male", "Female", "Other"}),
		)

		lines = append(lines,
			wid.GridRow(th, nil, 1.5, []float32{0, 0, 0},
				wid.RadioButton(th, &Alternative, "LargeColumns", "LargeColumns", wid.Do(onWinChange)),
				wid.RadioButton(th, &Alternative, "SmallColumns", "SmallColumns", wid.Do(onWinChange)),
				wid.RadioButton(th, &Alternative, "FractionalColumns", "FractionalColumns", wid.Do(onWinChange)),
			),
			wid.Space(40),
		)
		// Configure a row with headings.
		bgColor := th.Bg(wid.Primary)
		heading := wid.GridRow(th, &bgColor, 1.0, colWidths,
			wid.Checkbox(th, "", wid.Bool(&selectAll), wid.Do(onCheck), wid.Prim()),
			wid.HeaderButton(th, "Name", wid.Do(onNameClick), wid.Prim(), wid.BtnIcon(nameIcon)),
			wid.HeaderButton(th, "Address", wid.Do(onAddressClick), wid.Prim(), wid.BtnIcon(addressIcon)),
			wid.HeaderButton(th, "Age", wid.Do(onAgeClick), wid.Prim(), wid.BtnIcon(ageIcon)),
			// wid.Label(th, "Gender", wid.Bold()),
		)
		lines = append(lines, heading)

		// lines = append(lines, wid.Separator(th, unit.Dp(2.0), wid.W(9999)))
		for i := 0; i < len(data); i++ {
			bgColor := wid.MulAlpha(wid.Blue, 20)
			if i%2 == 0 {
				bgColor = wid.MulAlpha(wid.Red, 20)
			}
			lines = append(lines,
				wid.GridRow(th, &bgColor, 1.5, colWidths,
					wid.Checkbox(th, "", wid.Bool(&data[i].Selected)),
					wid.Label(th, &data[i].Name),
					wid.Label(th, &data[i].Address),
					wid.Label(th, &data[i].Age, wid.Dp(2), wid.Right()),
					// wid.DropDown(&thg, &data[i].Status, []string{"Male", "Female", "Other"}),
				))

		}
		return wid.List(th, anchor, lines...)
	}
}
