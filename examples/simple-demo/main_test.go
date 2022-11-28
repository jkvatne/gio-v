package main

import (
	"gio-v/wid"
	"image"
	"testing"

	"gioui.org/io/router"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"gioui.org/font/gofont"
)

func TestDemo(t *testing.T) {
	var ops op.Ops
	var r router.Router
	theme = wid.NewTheme(gofont.Collection(), 14)
	gtx := layout.NewContext(&ops, system.FrameEvent{
		Size: image.Point{
			X: 500,
			Y: 400,
		},
		Queue: &r,
	})
	form = demo(theme)
	form(gtx)
}

func BenchmarkDemo(b *testing.B) {
	theme = wid.NewTheme(gofont.Collection(), 14)
	b.ResetTimer()
	b.ReportAllocs()

	var ops op.Ops
	var r router.Router
	form = demo(theme)
	for i := 0; i < b.N; i++ {
		gtx := layout.NewContext(&ops, system.FrameEvent{
			Size: image.Point{
				X: 500,
				Y: 400,
			},
			Queue: &r,
		})
		form(gtx)
	}
}
