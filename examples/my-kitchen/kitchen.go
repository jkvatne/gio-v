package main

import (
	"github.com/jkvatne/gio-v/wid"
	"image"
	"image/color"
	"time"

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
	theme       *wid.Theme
	addIcon     *wid.Icon
	checkIcon   *wid.Icon
	group       string
	sliderValue float32 = 0.1
	win         app.Window
	progress    float32 = 0.1
	form        layout.Widget
	enabledText = "Disabled"
	enabled     bool
	blue        = true
	homeBg      = wid.RGB(0x1288F2)
	homeFg      = wid.RGB(0xFFFFFF)
	btnText     = "Blue"
	SomeText    = ""
)

func main() {
	checkIcon, _ = wid.NewIcon(icons.NavigationCheck)
	addIcon, _ = wid.NewIcon(icons.ContentAdd)
	theme = wid.NewTheme(gofont.Collection(), 14)
	win.Option(app.Title("Gio-v demo"), app.Size(unit.Dp(900), unit.Dp(650)))
	form = kitchen(theme)
	go ticker()
	wid.Run(&win, &form, theme)
	app.Main()
}

func ticker() {
	for {
		time.Sleep(time.Millisecond * 16)
		wid.GuiLock.Lock()
		progress = float32(int32((progress*1000)+5)%1000) / 1000.0
		wid.GuiLock.Unlock()
		wid.Invalidate()
	}
}

func onClick() {
	blue = !blue
	if blue {
		btnText = "Blue"
		homeBg = wid.RGB(0x1288F2)
		homeFg = wid.RGB(0xFFFFFF)
	} else {
		btnText = "Green"
		homeBg = wid.RGB(0x02F812)
		homeFg = wid.RGB(0x000000)
	}
}

func onDisable() {
	if enabled {
		enabledText = "Enabled"
	} else {
		enabledText = "Disabled"
	}
}

func colorBar(gtx wid.C) wid.D {
	gtx.Constraints.Min.Y = wid.Px(gtx, unit.Dp(50))
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
	return wid.Col(nil,
		wid.Label(th, topLabel, wid.Middle(), wid.FontSize(2.1)),

		wid.Label(th, longText),

		wid.Edit(th, &SomeText, wid.Hint("Value 1")),

		wid.Row(th, nil, wid.SpaceClose,
			wid.RoundButton(th, addIcon, wid.Hint("This is another dummy button"), wid.Role(wid.Primary)),
			wid.Button(th, "Icon", wid.BtnIcon(checkIcon), wid.Role(wid.Primary)),
			wid.Button(thb, "Click me!", wid.W(200), wid.Do(onClick), wid.Role(wid.Secondary)),
			wid.Button(thb, &btnText, wid.Fg(&homeFg), wid.Bg(&homeBg)),
			wid.TextButton(thb, "Flat"),
		),

		wid.ProgressBar(th, &progress),

		func(gtx wid.C) wid.D {
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, colorBar)
		},

		wid.Row(th, nil, wid.SpaceClose,
			wid.Switch(th, &enabled, wid.Do(onDisable)),
			wid.Button(th, &enabledText, wid.En(&enabled)),
		),

		wid.Row(th, nil, wid.SpaceDistribute,
			wid.RadioButton(th, &group, "RadioButton1", "RadioButton1"),
			wid.RadioButton(th, &group, "RadioButton2", "RadioButton2"),
			wid.RadioButton(th, &group, "RadioButton3", "RadioButton3"),
		),

		wid.Row(th, nil, []float32{0.9, 0.1},
			wid.Slider(th, &sliderValue, 0, 100),
		),
	)
}

const longText = `1. I learned from my grandfather, Verus, to use good manners, and to put restraint on anger. 
2. In the famous memory of my father I had a pattern of modesty and manliness. 
3. Of my mother I learned to be pious and generous; to keep myself not only from evil deeds, but even from evil thoughts; and to live with a simplicity which is far from customary among the rich. 
4. I owe it to my great-grandfather that I did not attend public lectures and discussions, but had good and able teachers at home; and I owe him also the knowledge that for things of this nature a man should count no expense too great.
5. My tutor taught me not to favour either green or blue at the chariot races, nor, in the contests of gladiators, to be a supporter either of light or heavy armed. He taught me also to endure labour; not to need many things; to serve myself without troubling others; not to intermeddle in the affairs of others, and not easily to listen to slanders against them.
`
