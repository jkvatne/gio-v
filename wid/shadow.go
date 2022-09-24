// SPDX-License-Identifier: Unlicense OR MIT

package wid

/*
This file is derived from work by Egon Elbre in his gio experiments
repository available here:

https://github.com/egonelbre/expgio/tree/master/box-shadows

He generously licensed it under the Unlicense, and thus it is
reproduced here under the same terms.
*/

import (
	"image"
	"image/color"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// ShadowStyle defines a shadow cast by a rounded rectangle.
//
// TODO(whereswaldon): make this support RRects that do not have
// uniform corner radii.
type ShadowStyle struct {
	CornerRadius int
	Elevation    int
}

// Shadow defines a shadow cast by a rounded rectangle with the given
// corner radius and elevation. It sets reasonable defaults for the
// shadow colors.
func Shadow(radius int, elevation int) ShadowStyle {
	return ShadowStyle{
		CornerRadius: radius,
		Elevation:    elevation,
	}
}

var alpha = [7]byte{0, 82, 62, 42, 32, 14, 13}

// Layout renders the shadow into the gtx. The shadow's size will assume
// that the rectangle casting the shadow is of size gtx.Constraints.Min.
func (s ShadowStyle) Layout(gtx C) D {
	sz := gtx.Constraints.Min
	for i := 6; i > 0; i-- {
		ofs := s.Elevation * i / 10
		rr := s.CornerRadius + ofs/2
		a := alpha[i]
		paint.FillShape(gtx.Ops, color.NRGBA{A: a}, clip.RRect{
			Rect: image.Rect(-ofs/2, -ofs/4, sz.X+ofs/2, sz.Y+ofs),
			SE:   rr, SW: rr, NW: rr, NE: rr,
		}.Op(gtx.Ops))
	}
	return D{Size: sz}
}
