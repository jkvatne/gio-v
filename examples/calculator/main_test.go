package main

import (
	"gioui.org/op"
	"github.com/jkvatne/gio-v/wid"
	"image"
	"testing"

	"gioui.org/font/gofont"
	"gioui.org/layout"
)

func TestButtons(t *testing.T) {
	theme = wid.NewTheme(gofont.Collection(), 14)
	gtx := layout.Context{
		Ops: new(op.Ops),
		// Rigid constraints with both minimum and maximum set.
		Constraints: layout.Exact(image.Point{X: 500, Y: 400}),
	}
	form = demo(theme)
	form(gtx)
}

func BenchmarkButtons(b *testing.B) {
	theme = wid.NewTheme(gofont.Collection(), 14)
	b.ResetTimer()
	b.ReportAllocs()

	form = demo(theme)
	for i := 0; i < b.N; i++ {
		gtx := layout.Context{
			Ops: new(op.Ops),
			// Rigid constraints with both minimum and maximum set.
			Constraints: layout.Exact(image.Point{X: 500, Y: 400}),
		}
		form(gtx)
	}
}
