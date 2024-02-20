// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/op"
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ProgressBarDef defines the progress bar
type ProgressBarDef struct {
	Base
	Progress  *float32
	Thickness unit.Dp
}

// ProgressBar returns a widget for a progress bar
func ProgressBar(th *Theme, progress *float32, options ...Option) layout.Widget {
	p := &ProgressBarDef{
		Base: Base{
			th:           th,
			cornerRadius: unit.Dp(10),
			role:         SurfaceVariant,
			padding:      layout.Inset{2, 2, 2, 2},
		},
		Progress:  progress,
		Thickness: 20,
	}
	// Read in options to change from default values to something else.
	for _, option := range options {
		option.apply(p)
	}
	return p.Layout
}

func ScaleInset(gtx C, ins layout.Inset) (pt, pb, pl, pr int) {
	pt = gtx.Metric.Dp(ins.Top)
	pb = gtx.Metric.Dp(ins.Bottom)
	pl = gtx.Metric.Dp(ins.Left)
	pr = gtx.Metric.Dp(ins.Right)
	return
}

func (p *ProgressBarDef) Layout(gtx C) D {
	pt, pb, pl, pr := ScaleInset(gtx, p.padding)
	progressBarWidth := gtx.Constraints.Min.X - pl - pr
	GuiLock.RLock()
	value := *p.Progress
	GuiLock.RUnlock()
	width := int(float32(progressBarWidth) * Clamp(value, 0, 1))
	thickness := Px(gtx, p.Thickness)
	color := p.Fg()
	gtx = gtx.Disabled()
	rr := Px(gtx, p.cornerRadius)
	if p.cornerRadius > (p.width-1)/2 {
		rr = Px(gtx, (p.width-1)/2)
	}
	dims := image.Point{X: progressBarWidth + pl + pr, Y: thickness + pt + pb}
	// Fill background if bgColor is given
	if p.bgColor != nil && (*p.bgColor).A != 0 {
		paint.FillShape(gtx.Ops, *p.bgColor, clip.UniformRRect(image.Rectangle{Max: dims}, 0).Op(gtx.Ops))
	}
	defer op.Offset(image.Pt(pl, pt)).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Point{X: width, Y: thickness}}, rr).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return D{Size: dims}
}

func (p *ProgressBarDef) setThickness(t unit.Dp) {
	p.Thickness = t
}
