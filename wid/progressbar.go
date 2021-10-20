// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ProgressBarStyle defines the progress bar
type ProgressBarStyle struct {
	Color        color.NRGBA
	TrackColor   color.NRGBA
	Progress     *float32
	Width        unit.Value
	CornerRadius unit.Value
}

// ProgressBar returns a widget for a progress bar
func ProgressBar(th *Theme, progress *float32) func(gtx C) D {
	p := &ProgressBarStyle{
		Progress:     progress,
		Width:        unit.Dp(10),
		CornerRadius: unit.Dp(5),
		Color:        th.Primary,
		TrackColor:   MulAlpha(th.OnBackground, 0x88),
	}
	return func(gtx C) D {
		return p.layout(gtx)
	}
}

func (p ProgressBarStyle) layout(gtx C) D {
	shader := func(width float32, color color.NRGBA) D {
		rr := float32(gtx.Px(unit.Dp(2)))
		d := image.Point{X: int(width), Y: gtx.Px(p.Width)}
		height := float32(gtx.Px(p.Width))
		defer clip.UniformRRect(f32.Rectangle{Max: f32.Pt(width, height)}, rr).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return D{Size: d}
	}
	progressBarWidth := float32(gtx.Constraints.Max.X)
	return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
		return layout.Stack{Alignment: layout.W}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return shader(progressBarWidth, p.TrackColor)
			}),
			layout.Stacked(func(gtx C) D {
				fillWidth := progressBarWidth * clamp1(*p.Progress)
				fillColor := p.Color
				if gtx.Queue == nil {
					fillColor = Disabled(fillColor)
				}
				return shader(fillWidth, fillColor)
			},
			),
		)
	},
	)
}

// clamp1 limits v to range [0..1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	} else {
		return v
	}
}
