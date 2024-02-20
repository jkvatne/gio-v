package main

import (
	"github.com/jkvatne/gio-v/wid"
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"

	"gioui.org/font/gofont"
)

func TestColors(t *testing.T) {
	theme = wid.NewTheme(gofont.Collection(), 14)
	gtx := layout.Context{
		Ops: new(op.Ops),
		// Rigid constraints with both minimum and maximum set.
		Constraints: layout.Exact(image.Point{X: 500, Y: 400}),
	}
	form = demo2(theme)
	form(gtx)
}

func BenchmarkColors(b *testing.B) {
	theme = wid.NewTheme(gofont.Collection(), 14)
	b.ResetTimer()
	b.ReportAllocs()

	form = demo1(theme)
	for i := 0; i < b.N; i++ {
		gtx := layout.Context{
			Ops: new(op.Ops),
			// Rigid constraints with both minimum and maximum set.
			Constraints: layout.Exact(image.Point{X: 500, Y: 400}),
		}
		form(gtx)
	}
}
