// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/op"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// ShadowStyle defines a shadow cast by a rounded rectangle.
type ShadowStyle struct {
	CornerRadius int
	Elevation    int
}

var alpha = [7]byte{0, 82, 62, 42, 32, 14, 13}

func DrawShadow(gtx C, outline image.Rectangle, rr int, elevation int) {
	for i := 6; i > 0; i-- {
		ofs := elevation * i / 10
		rr := rr + ofs/2
		a := alpha[i]
		paint.FillShape(gtx.Ops, color.NRGBA{A: a}, RrOp(clip.UniformRRect(outline, rr), ofs, gtx.Ops))
	}
}

// RrOp returns the op for the rounded rectangle.
func RrOp(rr clip.RRect, d int, ops *op.Ops) clip.Op {
	return clip.Outline{Path: ShadowPath(rr, d, ops)}.Op()
}

type Rectangle struct {
	Min, Max f32.Point
}

func FPt(p image.Point) f32.Point {
	return f32.Point{
		X: float32(p.X), Y: float32(p.Y),
	}
}

func FRect(r image.Rectangle) Rectangle {
	return Rectangle{
		Min: FPt(r.Min), Max: FPt(r.Max),
	}
}

// ShadowPath returns the PathSpec for the shadow
// This is a border around a rounded rectangle with width d
func ShadowPath(rr clip.RRect, d int, ops *op.Ops) clip.PathSpec {
	var p clip.Path
	p.Begin(ops)
	const iq = 1 - 4*(math.Sqrt2-1)/3
	se, sw, nw, ne := float32(rr.SE), float32(rr.SW), float32(rr.NW), float32(rr.NE)
	rrf := FRect(rr.Rect)
	w, n, e, s := rrf.Min.X, rrf.Min.Y, rrf.Max.X, rrf.Max.Y

	p.MoveTo(f32.Point{X: w + nw, Y: n})
	p.LineTo(f32.Point{X: e - ne, Y: n}) // N
	p.CubeTo(                            // NE
		f32.Point{X: e - ne*iq, Y: n},
		f32.Point{X: e, Y: n + ne*iq},
		f32.Point{X: e, Y: n + ne})
	p.LineTo(f32.Point{X: e, Y: s - se}) // E
	p.CubeTo(                            // SE
		f32.Point{X: e, Y: s - se*iq},
		f32.Point{X: e - se*iq, Y: s},
		f32.Point{X: e - se, Y: s})
	p.LineTo(f32.Point{X: w + sw, Y: s}) // S
	p.CubeTo(                            // SW
		f32.Point{X: w + sw*iq, Y: s},
		f32.Point{X: w, Y: s - sw*iq},
		f32.Point{X: w, Y: s - sw})
	p.LineTo(f32.Point{X: w, Y: n + nw}) // W
	p.CubeTo(                            // NW
		f32.Point{X: w, Y: n + nw*iq},
		f32.Point{X: w + nw*iq, Y: n},
		f32.Point{X: w + nw, Y: n})

	df := float32(d)
	se += df
	sw += df
	nw += df
	ne += df
	w -= df
	n -= df
	e += df
	s += df

	p.LineTo(f32.Point{X: w + nw, Y: n}) // Start W
	p.CubeTo(                            // NW
		f32.Point{X: w + nw*iq, Y: n},
		f32.Point{X: w, Y: n + nw*iq},
		f32.Point{X: w, Y: n + nw})
	p.LineTo(f32.Point{X: w, Y: s - sw}) // W
	p.CubeTo(                            // SW
		f32.Point{X: w, Y: s - sw*iq},
		f32.Point{X: w + sw*iq, Y: s},
		f32.Point{X: w + sw, Y: s})
	p.LineTo(f32.Point{X: e - sw, Y: s}) // S
	p.CubeTo(                            // SE
		f32.Point{X: e - se*iq, Y: s},
		f32.Point{X: e, Y: s - se*iq},
		f32.Point{X: e, Y: s - se})
	p.LineTo(f32.Point{X: e, Y: n + ne}) // E
	p.CubeTo(                            // NE
		f32.Point{X: e, Y: n + ne*iq},
		f32.Point{X: e - ne*iq, Y: n},
		f32.Point{X: e - ne, Y: n})
	p.LineTo(f32.Point{X: nw, Y: n})     // N
	p.LineTo(f32.Point{X: w + nw, Y: n}) // To start

	return p.End()
}

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
