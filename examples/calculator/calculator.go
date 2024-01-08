// SPDX-License-Identifier: Unlicense OR MIT

package main

// A Gio program that demonstrates gio-v widgets.
// See https://gioui.org for information on the gio
// gio-v is maintained by Jan Kåre Vatne (jkvatne@online.no)

import (
	"fmt"
	"gioui.org/font/gofont"
	"github.com/jkvatne/gio-v/wid"
	"image/color"
	"math"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
)

var (
	theme     *wid.Theme  // the theme selected
	win       *app.Window // The main window
	form      layout.Widget
	entry     float64
	operator  rune
	dpNo      int
	operand   float64
	dpPressed = false
)

func main() {
	theme = wid.NewTheme(gofont.Collection(), 12)
	theme.Scale = 2.0
	theme.DarkMode = false
	theme.SecondaryColor = color.NRGBA{100, 200, 100, 255}
	theme.UpdateColors()
	win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(450), unit.Dp(700)))
	form = demo(theme)
	go wid.Run(win, &form, theme)
	app.Main()
}

func NumDecPlaces(v float64) int {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i <= -1 {
		return 0
	}
	return len(s) - i - 1
}

func clearAcc() {
	entry = 0
	operand = 0
	dpNo = 0
	dpPressed = false
}

func addDigt(x float64) {
	if !dpPressed {
		entry = 10*entry + x
	} else {
		dpNo++
		entry = entry + float64(x)/math.Pow(10, float64(dpNo))
	}
}

func setOp(ch rune) {
	operator = ch
	operand = entry
	entry = 0
	dpNo = 0
	dpPressed = false
	fmt.Printf("entry=%0.3f  operand=%0.3f  op=%d\n", entry, operand, operator)
}

// Demo setup. Called from Setup(), only once - at start of showing it.
// Returns a widget - i.e. a function: func(gtx C) D
func demo(th *wid.Theme) layout.Widget {
	return wid.Col(wid.SpaceClose,
		wid.Container(th, wid.PrimaryContainer, 0, layout.Inset{}, layout.Inset{},
			wid.Label(th, "Calculator demo", wid.Middle(), wid.Heading(), wid.Bold(), wid.Role(wid.PrimaryContainer)),
		),
		wid.Col(wid.SpaceClose,
			wid.Edit(th, &entry, &dpNo, wid.FontSize(1.8)),
			wid.Row(th, nil, wid.SpaceClose,
				wid.Button(th, "AC", wid.RR(999), wid.FontSize(1.4), wid.Pads(12, 1), wid.Role(wid.PrimaryContainer),
					wid.Do(func() { clearAcc() })),
				wid.Button(th, "C", wid.RR(999), wid.FontSize(2), wid.Pads(7, 3), wid.Role(wid.SecondaryContainer)),
				wid.Button(th, "%", wid.RR(999), wid.FontSize(2), wid.Pads(8, 1), wid.Role(wid.SecondaryContainer),
					wid.Do(func() {
						entry = entry / 100
					})),
				wid.Button(th, "/", wid.RR(999), wid.FontSize(2), wid.Pads(8, 7), wid.Role(wid.SecondaryContainer),
					wid.Do(func() { setOp('/') })),
			),
			wid.Row(th, nil, wid.SpaceClose,
				wid.Button(th, "7", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(7) })),
				wid.Button(th, "8", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(8) })),
				wid.Button(th, "9", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(9) })),
				wid.Button(th, "x", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SecondaryContainer),
					wid.Do(func() { setOp('*') })),
			),
			wid.Row(th, nil, wid.SpaceClose,
				wid.Button(th, "4", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(4) })),
				wid.Button(th, "5", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(5) })),
				wid.Button(th, "6", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(6) })),
				wid.Button(th, "-", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SecondaryContainer),
					wid.Do(func() { setOp('-') })),
			),
			wid.Row(th, nil, wid.SpaceClose,
				wid.Button(th, "1", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(1) })),
				wid.Button(th, "2", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(2) })),
				wid.Button(th, "3", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(3) })),
				wid.Button(th, "+", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SecondaryContainer),
					wid.Do(func() { setOp('+') })),
			),
			wid.Row(th, nil, wid.SpaceClose,
				wid.Button(th, "0", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { addDigt(0) })),
				wid.Button(th, "\u2022", wid.RR(999), wid.FontSize(2), wid.Pads(8, 7), wid.Role(wid.SurfaceContainer),
					wid.Do(func() { dpNo = 0; dpPressed = true })),
				wid.Button(th, "    =    ", wid.RR(999), wid.FontSize(2), wid.Pads(8, 5), wid.Role(wid.TertiaryContainer),
					wid.Do(func() {
						fmt.Printf("entry=%0.3f  operand=%0.3f  op=%d\n", entry, operand, operator)
						switch operator {
						case '+':
							entry = entry + operand
						case '-':
							entry = entry - operand
						case '*':
							entry = entry * operand
						case '/':
							entry = operand / entry
						default:
							panic("Invalid operand")
						}
						operand = 0
						dpNo = NumDecPlaces(entry)
						dpPressed = false
						fmt.Printf("entry=%0.3f  operand=%0.3f\n", entry, operand)
					})),
			),
		),
	)
}
