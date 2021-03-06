// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// BorderDef lays out a widget and draws a border inside it.
type BorderDef struct {
	Color        color.NRGBA
	CornerRadius unit.Value
	Width        unit.Value
}

// Layout will draw the border
func (b BorderDef) Layout(gtx C, w layout.Widget) D {
	dims := w(gtx)
	sz := layout.FPt(dims.Size)
	rr := float32(gtx.Px(b.CornerRadius))
	width := float32(gtx.Px(b.Width))
	sz.X -= width
	sz.Y -= width

	r := f32.Rectangle{Max: sz}
	r = r.Add(f32.Point{X: width * 0.5, Y: width * 0.5})

	paint.FillShape(gtx.Ops,
		b.Color,
		clip.Stroke{
			Path:  clip.UniformRRect(r, rr).Path(gtx.Ops),
			Width: width,
		}.Op(),
	)

	return dims
}
