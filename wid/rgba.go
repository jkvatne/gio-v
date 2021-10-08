// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image/color"

	"gioui.org/unit"
)

// Zv is a zero unit.Value. Just saving a few keystrokes
var Zv = unit.Value{}

// DeEmphasis will change a color to a less prominent color
// In light mode, colors will be lighter, in dark mode, colors will be darker
// The amount of darkening is greater than the amount of lightening
func DeEmphasis(c color.NRGBA, amount uint8) color.NRGBA {
	if Luminance(c) < 128 {
		return MulAlpha(c, 255-amount)
	}
	return MulAlpha(c, amount)
}

// Pxr maps the value v to pixels, returning a float32
func Pxr(c C, v unit.Value) float32 {
	return float32(c.Metric.Px(v))
}

// Disabled blends color towards the luminance and multiplies alpha.
// Blending towards luminance will desaturate the color.
// Multiplying alpha blends the color together more with the background.
func Disabled(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := Luminance(c)
	return color.NRGBA{
		R: byte((int(c.R)*r + int(lum)*(256-r)) / 256),
		G: byte((int(c.G)*r + int(lum)*(256-r)) / 256),
		B: byte((int(c.B)*r + int(lum)*(256-r)) / 256),
		A: byte(int(c.A) * (128 + 32) / 256),
	}
}

// ColDisabled returns the disabled color of c, depending on the disabled flag.
func ColDisabled(c color.NRGBA, disabled bool) color.NRGBA {
	if disabled {
		return Disabled(c)
	}
	return c
}

// Hovered blends color towards a brighter color.
func Hovered(c color.NRGBA) (d color.NRGBA) {
	const r = 0x40 // lighten ratio
	return color.NRGBA{
		R: byte(255 - int(255-c.R)*(255-r)/256),
		G: byte(255 - int(255-c.G)*(255-r)/256),
		B: byte(255 - int(255-c.B)*(255-r)/256),
		A: c.A,
	}
}

// Interpolate returns a color in between given colors a and b, depending on progress
func Interpolate(a, b color.NRGBA, progress float32) color.NRGBA {
	var out color.NRGBA
	out.R = uint8(int16(a.R) - int16(float32(int16(a.R)-int16(b.R))*progress))
	out.G = uint8(int16(a.G) - int16(float32(int16(a.G)-int16(b.G))*progress))
	out.B = uint8(int16(a.B) - int16(float32(int16(a.B)-int16(b.B))*progress))
	out.A = uint8(int16(a.A) - int16(float32(int16(a.A)-int16(b.A))*progress))
	return out
}

// Gray returns a NRGBA color with the same luminance as the parameter
func Gray(c color.NRGBA) color.NRGBA {
	l := Luminance(c)
	return color.NRGBA{R: l, G: l, B: l, A: c.A}
}

// RGB creates a NRGBA color from its hex code, with alpha=255
func RGB(c uint32) color.NRGBA {
	return ARGB(0xff000000 | c)
}

// ARGB creates a NRGBA color from its hex code
func ARGB(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

// WithAlpha returns the input color with the new alpha value.
func WithAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = alpha
	return c
}

// MulAlpha applies the alpha to the color.
func MulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = uint8(uint32(c.A) * uint32(alpha) / 0xFF)
	return c
}

// Luminance is a fast approximate version of RGBA.Luminance.
func Luminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}
