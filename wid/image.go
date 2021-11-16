// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	_ "image/jpeg"
	"os"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// )

// ImageDef is a widget that displays an image.
type ImageDef struct {
	// Src is the image to display.
	Src paint.ImageOp
	// Fit specifies how to scale the image to the constraints.
	// By default, it does not do any scaling.
	Fit Fit
	// Position specifies where to position the image within
	// the constraints.
	Position layout.Direction
	// Scale is the ratio of image pixels to
	// dps. If Scale is zero Image falls back to
	// a scale that match a standard 72 DPI.
	Scale float32
}

func ImageFromJpgFile(filename string) func(gtx C) D {
	f, err := os.Open(filename)
	if err != nil {
		// Handle error
	}
	defer f.Close()

	pict, _, err := image.Decode(f)
	if err != nil {
		// Handle error
	}
	return Image(pict)
}

func Image(img image.Image) func(gtx C) D {
	src := paint.NewImageOp(img)
	im := ImageDef{}
	im.Fit = Contain
	im.Src = src
	return func(gtx C) D {
		return im.Layout(gtx)
	}
}

func (im ImageDef) Layout(gtx layout.Context) layout.Dimensions {
	scale := im.Scale
	if scale == 0 {
		scale = float32(160.0 / 72.0)
	}

	size := im.Src.Size()
	wf, hf := float32(size.X), float32(size.Y)
	w, h := gtx.Px(unit.Dp(wf*scale)), gtx.Px(unit.Dp(hf*scale))

	dims, trans := im.Fit.scale(gtx.Constraints, im.Position, layout.Dimensions{Size: image.Pt(w, h)})
	defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	pixelScale := scale * gtx.Metric.PxPerDp
	trans = trans.Mul(f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(pixelScale, pixelScale)))
	defer op.Affine(trans).Push(gtx.Ops).Pop()

	im.Src.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return dims
}
