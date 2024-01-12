// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"golang.org/x/exp/shiny/iconvg"
)

// Icon is the definition of an icon
type Icon struct {
	src []byte
	// Cached values.
	op       paint.ImageOp
	imgSize  int
	imgColor color.NRGBA
}

// NewIcon returns a new Icon from IconVG data.
func NewIcon(data []byte) (*Icon, error) {
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil, err
	}
	return &Icon{src: data}, nil
}

// Layout displays the icon with its size set to the X minimum constraint.
func (ic *Icon) Layout(gtx C, color color.NRGBA) D {
	sz := gtx.Constraints.Min.X
	size := gtx.Constraints.Constrain(image.Pt(sz, sz))
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	ico := ic.image(size.X, color)
	ico.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return D{Size: ico.Size()}
}

func (ic *Icon) Update(data []byte) error {
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return err
	}
	ic.src = data
	ic.imgSize = 0
	ic.imgColor = color.NRGBA{}
	return nil
}

func (ic *Icon) image(sz int, c color.NRGBA) paint.ImageOp {
	if sz < 1 {
		sz = 1
	}
	if sz == ic.imgSize && c == ic.imgColor {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)

	// palette uses pre-multiplied RGBA colors. Apply pre-multiplication here.
	r, g, b, a := c.RGBA()
	m.Palette[0] = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}

	_ = iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	ic.imgColor = c
	return ic.op
}
