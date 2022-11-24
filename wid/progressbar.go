// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"

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
	progressBarWidth := gtx.Constraints.Min.X - gtx.Dp(4)
	return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx C) D {
		GuiLock.RLock()
		value := *p.Progress
		GuiLock.RUnlock()
		width := int(float32(progressBarWidth) * Clamp(value, 0, 1))
		color := p.Fg()
		if gtx.Queue == nil {
			color = Disabled(color)
		}
		rr := rr(gtx, p.cornerRadius, gtx.Dp(p.width))
		d := image.Point{X: width, Y: gtx.Dp(p.width)}
		height := p.width
		defer clip.UniformRRect(image.Rectangle{Max: image.Pt(width, gtx.Dp(height))}, rr).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return D{Size: d}
	})
}
