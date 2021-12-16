package main

import (
	"fmt"
	"gio-v/wid"
	"time"

	"gioui.org/layout"
)

func endKitchen() {
	page = "KitchenX"
	oldMode = "xx"
	PrintMemUsage("Gio-v")
	startTime = time.Now()
	count = 0
}

func kitchenV(th *wid.Theme) layout.Widget {
	thb = th
	return wid.Col(
		wid.Label(th, topLabel, wid.Middle(), wid.Size(2.1)),
		wid.Edit(th, wid.Hint("Value 1")),
		wid.Edit(th, wid.Hint("Value 2")),
		wid.Row(th, nil, []float32{35, 10, 20, 15, 20},
			wid.Button(thb, "Click me!", wid.W(500), wid.Handler(onClick)),
			wid.RoundButton(th, addIcon, wid.Hint("This is another dummy button")),
			wid.Button(th, "Icon", wid.BtnIcon(checkIcon), wid.Color(wid.RGB(0xffff00))),
			wid.Button(thb, "Blue", wid.Color(wid.Blue)),
			wid.TextButton(th, "Show other", wid.Handler(endKitchen)),
		),
		wid.Row(th, nil, nil,
			wid.ProgressBar(th, &progress),
			wid.Value(th, func() string { return fmt.Sprintf(" %0.1f frames/second", count/time.Since(startTime).Seconds()) }),
		),

		wid.Row(th, nil, nil,
			wid.RadioButton(th, &radioButtonValue, "RadioButton1", "RadioButton1"),
			wid.RadioButton(th, &radioButtonValue, "RadioButton2", "RadioButton2"),
			wid.RadioButton(th, &radioButtonValue, "RadioButton3", "RadioButton3"),
		),
		wid.Row(th, nil, []float32{0.9, 0.1},
			wid.Slider(th, &sliderValue, 0, 100),
			wid.Value(th, func() string { return fmt.Sprintf("  %0.2f", sliderValue) }),
		),
	)
}
