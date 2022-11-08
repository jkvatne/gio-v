// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ProgressBarStyle defines the progress bar
type ProgressBarStyle struct {
	Base
	Progress *float32
}

// ProgressBar returns a widget for a progress bar
func ProgressBar(th *Theme, progress *float32, options ...Option) func(gtx C) D {
	p := &ProgressBarStyle{
		Progress: progress,
	}
	p.cornerRadius = unit.Dp(10)
	p.width = 10
	p.role = Primary
	p.Apply(options...)
	p.role = SurfaceVariant
	p.width = 10
	p.th = th
	return func(gtx C) D {
		return p.layout(gtx)
	}
}

func (p ProgressBarStyle) layout(gtx C) D {
	shader := func(width int, color color.NRGBA) D {
		rr := rr(gtx, p.cornerRadius, gtx.Dp(p.width))
		d := image.Point{X: width, Y: gtx.Dp(p.width)}
		height := p.width
		defer clip.UniformRRect(image.Rectangle{Max: image.Pt(width, gtx.Dp(height))}, rr).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return D{Size: d}
	}
	progressBarWidth := gtx.Constraints.Max.X - gtx.Dp(4)
	return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
		return layout.Stack{Alignment: layout.W}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return shader(progressBarWidth, p.Bg())
			}),
			layout.Stacked(func(gtx C) D {
				GuiLock.RLock()
				value := *p.Progress
				GuiLock.RUnlock()
				fillWidth := int(float32(progressBarWidth) * clamp1(value))
				fillColor := p.Fg()
				if gtx.Queue == nil {
					fillColor = Disabled(fillColor)
				}
				return shader(fillWidth, fillColor)
			}),
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
