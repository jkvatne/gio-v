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
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ShadowStyle defines a shadow cast by a rounded rectangle.
//
// TODO(whereswaldon): make this support RRects that do not have
// uniform corner radii.
type ShadowStyle struct {
	// The radius of the corners of the rectangle casting the surface.
	// Non-rounded rectangles can just provide a zero.
	CornerRadius unit.Value
	// Elevation is how high the surface casting the shadow is above
	// the background, and therefore determines how diffuse and large
	// the shadow is.
	Elevation unit.Value
	// The colors of various components of the shadow. The Shadow()
	// constructor populates these with reasonable defaults.
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
func (s ShadowStyle) Layout(gtx layout.Context) layout.Dimensions {
	sz := gtx.Constraints.Min
	ofs := pxf(gtx.Metric, s.Elevation)/2
	rr := float32(gtx.Px(s.CornerRadius))
	penumbra := f32.Rect(-ofs, -ofs, float32(sz.X)+2*ofs, float32(sz.Y)+2*ofs)
	umbra := f32.Rect(0, 0, float32(sz.X)+ofs, float32(sz.Y)+ofs)
	paint.FillShape(gtx.Ops, s.UmbraColor, clip.RRect{
		Rect: umbra,
		SE:   rr, SW: rr, NW: rr, NE: rr,
	}.Op(gtx.Ops))
	paint.FillShape(gtx.Ops, s.PenumbraColor, clip.RRect{
		Rect: penumbra,
		SE:   rr, SW: rr, NW: rr, NE: rr,
	}.Op(gtx.Ops))
	return layout.Dimensions{Size: sz}
}

func imageRect(r f32.Rectangle) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(math.Round(float64(r.Min.X))),
			Y: int(math.Round(float64(r.Min.Y))),
		},
		Max: image.Point{
			X: int(math.Round(float64(r.Max.X))),
			Y: int(math.Round(float64(r.Max.Y))),
		},
	}
}

func round(r f32.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{
			X: float32(math.Round(float64(r.Min.X))),
			Y: float32(math.Round(float64(r.Min.Y))),
		},
		Max: f32.Point{
			X: float32(math.Round(float64(r.Max.X))),
			Y: float32(math.Round(float64(r.Max.Y))),
		},
	}
}

func outset(r f32.Rectangle, rr float32) f32.Rectangle {
	r.Min.X -= rr
	r.Min.Y -= rr
	r.Max.X += rr
	r.Max.Y += rr
	return r
}

func pxf(c unit.Metric, v unit.Value) float32 {
	switch v.U {
	case unit.UnitPx:
		return v.V
	case unit.UnitDp:
		s := c.PxPerDp
		if s == 0 {
			s = 1
		}
		return s * v.V
	case unit.UnitSp:
		s := c.PxPerSp
		if s == 0 {
			s = 1
		}
		return s * v.V
	default:
		panic("unknown unit")
	}
}

func topLeft(r image.Rectangle) image.Point     { return r.Min }
func topRight(r image.Rectangle) image.Point    { return image.Point{X: r.Max.X, Y: r.Min.Y} }
func bottomRight(r image.Rectangle) image.Point { return r.Max }
func bottomLeft(r image.Rectangle) image.Point  { return image.Point{X: r.Min.X, Y: r.Max.Y} }
