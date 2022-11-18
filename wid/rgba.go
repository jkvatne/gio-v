// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"fmt"
	"image/color"
	"math"
)

// Some default colors
var (
	Red    = RGB(0xFF0000)
	Yellow = RGB(0xFFFF00)
	Green  = RGB(0x00FF00)
	Blue   = RGB(0x0000FF)
	White  = RGB(0xFFFFFF)
	Black  = RGB(0x000000)
)

const inf = 5000

// DeEmphasis will change a color to a less prominent color
// In light mode, colors will be lighter, in dark mode, colors will be darker
// The amount of darkening is greater than the amount of lightening
func DeEmphasis(c color.NRGBA, amount uint8) color.NRGBA {
	if Luminance(c) < 128 {
		return MulAlpha(c, 255-amount)
	}
	return MulAlpha(c, amount)
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

// Interpolate returns a color in between given colors a and b, depending on progress from 0.0 to 1.0
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

// Internal implementation converting RGB to HSL, HSV, or HSI.
// Basically a direct implementation of this: https://en.wikipedia.org/wiki/HSL_and_HSV#General_approach
func Rgb2hsl(c color.NRGBA) (float64, float64, float64) {
	var h, s, lvi float64
	var huePrime float64
	r := float64(c.R) / 256.0
	g := float64(c.G) / 256.0
	b := float64(c.B) / 256.0
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	chroma := (max - min)
	if chroma == 0 {
		h = 0
	} else {
		if r == max {
			huePrime = math.Mod(((g - b) / chroma), 6)
		} else if g == max {
			huePrime = ((b - r) / chroma) + 2

		} else if b == max {
			huePrime = ((r - g) / chroma) + 4

		}

		h = huePrime * 60
	}
	if r == g && g == b {
		lvi = r
	} else {
		lvi = (max + min) / 2
	}
	if lvi == 1 {
		s = 0
	} else {
		s = (chroma / (1 - math.Abs(2*lvi-1)))
	}

	if math.IsNaN(s) {
		s = 0
	}

	if h < 0 {
		h = 360 + h
	}

	return h, s, lvi
}

// Internal HSV->RGB function for doing conversions using float inputs (saturation, value) and
// outputs (for R, G, and B).
// Basically a direct implementation of this: https://en.wikipedia.org/wiki/HSL_and_HSV#Converting_to_RGB
func Hsl2rgb(hueDegrees float64, saturation float64, light float64) color.NRGBA {
	var r, g, b float64
	hueDegrees = math.Mod(hueDegrees, 360)
	if saturation == 0 {
		r = light
		g = light
		b = light
	} else {
		var chroma float64
		var m float64
		chroma = (1 - math.Abs((2*light)-1)) * saturation
		hueSector := hueDegrees / 60
		intermediate := chroma * (1 - math.Abs(math.Mod(hueSector, 2)-1))
		switch {
		case hueSector >= 0 && hueSector <= 1:
			r = chroma
			g = intermediate
			b = 0
		case hueSector > 1 && hueSector <= 2:
			r = intermediate
			g = chroma
			b = 0
		case hueSector > 2 && hueSector <= 3:
			r = 0
			g = chroma
			b = intermediate
		case hueSector > 3 && hueSector <= 4:
			r = 0
			g = intermediate
			b = chroma
		case hueSector > 4 && hueSector <= 5:
			r = intermediate
			g = 0
			b = chroma
		case hueSector > 5 && hueSector <= 6:
			r = chroma
			g = 0
			b = intermediate
		default:
			panic(fmt.Errorf("hue input %v yielded sector %v", hueDegrees, hueSector))
		}
		m = light - (chroma / 2)
		r += m
		g += m
		b += m
	}
	return color.NRGBA{R: uint8(r*255 + 0.4), G: uint8(g*255 + 0.4), B: uint8(b*255 + 0.4), A: 255}
}
