package main

import (
	"gio-v/wid"
	"image"
	"image/color"

	"gioui.org/op/clip"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/layout"
)

var (
	topLabel    = "Hello, Gio"
	thb         *wid.Theme
	th          *wid.Theme
	addIcon     *wid.Icon
	checkIcon   *wid.Icon
	group       string  = ""
	sliderValue float32 = 0.1
	win         *app.Window
	progress    float32
	disable     bool
	form        layout.Widget
)

func main() {
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	addIcon, _ = wid.NewIcon(icons.ContentAdd)
	th = wid.NewTheme(gofont.Collection(), 14)
	win = app.NewWindow(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(500)))
	form = kitchen(th)
	wid.Run(win, &form)
	app.Main()
}

func onClick() {

}

func colorBar(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(50))
	gtx.Constraints.Max.Y = gtx.Constraints.Min.Y

	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.LinearGradientOp{
		Stop1:  layout.FPt(dr.Min),
		Stop2:  layout.FPt(dr.Max),
		Color1: color.NRGBA{R: 0x10, G: 0xff, B: 0x10, A: 0xFF},
		Color2: color.NRGBA{R: 0x10, G: 0x10, B: 0xff, A: 0xFF},
	}.Add(gtx.Ops)
	defer clip.Rect(dr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func kitchen(th *wid.Theme) layout.Widget {
	thb = th
	return wid.List(th, wid.Occupy,
		wid.Label(th, topLabel, wid.Middle(), wid.FontSize(2.1)),
		wid.Edit(th, wid.Hint("Value 1")),
		wid.Edit(th, wid.Hint("Value 2")),
		wid.Row(th, nil, wid.SpaceClose,
			wid.RoundButton(th, addIcon, wid.Hint("This is another dummy button"), wid.Role(wid.Primary)),
			wid.Button(th, "Icon", wid.BtnIcon(checkIcon), wid.Role(wid.Primary)),
			wid.Button(thb, "Click me!", wid.W(200), wid.Do(onClick), wid.Role(wid.Secondary)),
			wid.Button(thb, "Blue", wid.Role(wid.Primary)),
			wid.TextButton(thb, "Flat"),
		),
		wid.ProgressBar(th, &progress),
		func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, colorBar)
		},
		wid.Row(th, nil, wid.SpaceClose,
			wid.Switch(th, &disable),
			wid.Button(th, "disabled", wid.En(&disable)),
		),

		wid.Row(th, nil, nil,
			wid.RadioButton(th, &group, "RadioButton1", "RadioButton1"),
			wid.RadioButton(th, &group, "RadioButton2", "RadioButton2"),
			wid.RadioButton(th, &group, "RadioButton3", "RadioButton3"),
		),
		wid.Row(th, nil, []float32{0.9, 0.1},
			wid.Slider(th, &sliderValue, 0, 100),
			wid.Label(th, &sliderValue, wid.Dp(2), wid.Pads(10)),
		),
	)
}

const longText = `1. I learned from my grandfather, Verus, to use good manners, and to
put restraint on anger. 2. In the famous memory of my father I had a
pattern of modesty and manliness. 3. Of my mother I learned to be
pious and generous; to keep myself not only from evil deeds, but even
from evil thoughts; and to live with a simplicity which is far from
customary among the rich. 4. I owe it to my great-grandfather that I
did not attend public lectures and discussions, but had good and able
teachers at home; and I owe him also the knowledge that for things of
this nature a man should count no expense too great.

5. My tutor taught me not to favour either green or blue at the
chariot races, nor, in the contests of gladiators, to be a supporter
either of light or heavy armed. He taught me also to endure labour;
not to need many things; to serve myself without troubling others; not
to intermeddle in the affairs of others, and not easily to listen to
slanders against them.
`
