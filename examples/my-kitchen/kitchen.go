package main

import (
	"fmt"
	"gio-v/wid"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/widget"

	"gioui.org/layout"
)

var (
	topLabel         = "Hello, Gio"
	radioButtonValue bool
	thb              *wid.Theme
	th               *wid.Theme
	addIcon          *widget.Icon
	checkIcon        *widget.Icon
	group                        = new(widget.Enum)
	sliderValue      float32     = 0.1
	win              *app.Window // The main window
	form             layout.Widget
	progress         float32
)

func main() {
	c := wid.RGB(0x123456)
	d := wid.Hsl2rgb(wid.Rgb2hsl(c))
	fmt.Printf("D=%d", d)
	checkIcon, _ = widget.NewIcon(icons.NavigationCheck)
	addIcon, _ = widget.NewIcon(icons.ContentAdd)

	go func() {
		th = wid.NewTheme(gofont.Collection(), 14)
		win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
		form = kitchen(th)
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
	c := th.Bg(wid.Canvas)
	paint.Fill(gtx.Ops, c)
	// A hack to fetch mouse position and window size so we can avoid
	// tooltips going outside the main window area
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	wid.UpdateMousePos(gtx, win, e.Size)
	progress = progress + 0.01
	if progress > 1.0 {
		progress = 0
	}
	// Draw widgets
	form(gtx)
	// Apply the actual screen drawing
	e.Frame(gtx.Ops)
}

func onClick() {

}

func kitchen(th *wid.Theme) layout.Widget {
	thb = th
	return wid.Col(
		wid.Label(th, topLabel, wid.Middle(), wid.FontSize(2.1)),
		wid.Edit(th, wid.Hint("Value 1")),
		wid.Edit(th, wid.Hint("Value 2")),
		wid.Row(th, nil, []float32{25, 10, 20, 15, 20},
			wid.Button(thb, "Click me!", wid.W(100), wid.Do(onClick)),
			wid.RoundButton(th, addIcon, wid.Hint("This is another dummy button")),
			wid.Button(th, "Icon", wid.BtnIcon(checkIcon), wid.Fg(wid.RGB(0x00ffff))),
			wid.Button(thb, "Blue", wid.Role(wid.Primary), wid.Bg(wid.RGB(0x0000ff))),
		),
		// TODO wid.Row(th, nil, nil,
		// wid.ProgressBar(th, &progress),
		// wid.Value(th, func() string { return fmt.Sprintf(" %0.1f frames/second", count/time.Since(startTime).Seconds()) }),
		// ),

		wid.Row(th, nil, nil,
			wid.RadioButton(th, group, "RadioButton1", "RadioButton1"),
			wid.RadioButton(th, group, "RadioButton2", "RadioButton2"),
			wid.RadioButton(th, group, "RadioButton3", "RadioButton3"),
		),
		wid.Row(th, nil, []float32{0.9, 0.1},
			wid.Slider(th, &sliderValue, 0, 100),
			wid.Value(th, func() string { return fmt.Sprintf("  %0.2f", sliderValue) }),
		),
	)
}
