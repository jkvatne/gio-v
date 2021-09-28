// SPDX-License-Identifier: Unlicense OR MIT

package wid

/*
This file is derived from work by Egon Elbre in his gio experiments
repository available here:

https://github.com/egonelbre/expgio/tree/master/box-shadows

He generously licensed it under the Unlicense, and thus is is
reproduced here under the same terms.
*/

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image/color"
)

// ShadowStyle defines a shadow cast by a rounded rectangle.
//
// TODO(whereswaldon): make this support RRects that do not have
// uniform corner radii.
type ShadowStyle struct {
	CornerRadius unit.Value
	Elevation unit.Value
	AmbientColor, PenumbraColor, UmbraColor color.NRGBA
}

// Shadow defines a shadow cast by a rounded rectangle with the given
// corner radius and elevation. It sets reasonable defaults for the
// shadow colors.
func Shadow(radius, elevation unit.Value) ShadowStyle {
	return ShadowStyle{
		CornerRadius:  radius,
		Elevation:     elevation,
		UmbraColor:    color.NRGBA{A: 0x50},
		PenumbraColor: color.NRGBA{A: 0x28},
	}
}

// Layout renders the shadow into the gtx. The shadow's size will assume
// that the rectangle casting the shadow is of size gtx.Constraints.Min.
var alpha = [7]byte{0, 82, 62, 42, 32, 14, 13}

func (s ShadowStyle) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Constraints.Min
	for i:=6; i>0; i-- {
		ofs := Pxr(gtx, s.Elevation)*float32(i)/10
		rr := Pxr(gtx, s.CornerRadius)+ofs/2
		a := alpha[i]
		paint.FillShape(gtx.Ops, color.NRGBA{A: a}, clip.RRect{
			Rect: f32.Rect(-ofs/2, -ofs/4, float32(sz.X)+ofs/2, float32(sz.Y)+ofs),
			SE:   rr, SW: rr, NW: rr, NE: rr,
		}.Op(gtx.Ops))
	}
	return layout.Dimensions{Size: sz}
}
