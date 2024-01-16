// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"os"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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

func ImageFromJpgFile(filename string, fit Fit) layout.Widget {
	f, err := os.Open(filename)
	defer func() {
		err := f.Close()
		if err != nil {
			panic("Jpg file close failed")
		}
	}()
	if err != nil {
		panic(fmt.Sprintf("Image '%s' not found", filename))
	}
	pict, _, err := image.Decode(f)
	if err != nil {
		panic(fmt.Sprintf("Image '%s' has unknown format", filename))
	}
	return Image(pict, fit)
}

func Image(img image.Image, fit Fit) layout.Widget {
	src := paint.NewImageOp(img)
	im := ImageDef{}
	im.Fit = fit
	im.Src = src
	return func(gtx C) D {
		return im.Layout(gtx)
	}
}

func (im ImageDef) Layout(gtx C) D {
	scale := im.Scale
	if scale == 0 {
		scale = gtx.Metric.PxPerDp
	}

	w := int(float32(im.Src.Size().X) * scale)
	h := int(float32(im.Src.Size().Y) * scale)

	dims, trans := im.Fit.scale(gtx.Constraints, im.Position, layout.Dimensions{Size: image.Pt(w, h)})
	defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	trans = trans.Mul(f32.Affine2D{}.Scale(f32.Point{}, f32.Pt(scale, scale)))
	defer op.Affine(trans).Push(gtx.Ops).Pop()

	im.Src.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return dims
}
