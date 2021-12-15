package main

import (
	"gio-v/wid"
	"image"
	"testing"

	"gioui.org/widget/material"

	"gioui.org/font/gofont"
	"gioui.org/io/router"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

func DoLayout(b *testing.B, show string, n int) {
	currentTheme = wid.NewTheme(gofont.Collection(), 18, wid.MaterialDesignLight)
	page = show
	setup()
	makePersons(10000)
	data = data[0:n]
	th = material.NewTheme(gofont.Collection())
	b.ResetTimer()
	b.ReportAllocs()

	var ops op.Ops
	var r router.Router
	for i := 0; i < b.N; i++ {
		gtx := layout.NewContext(&ops, system.FrameEvent{
			Size: image.Point{
				X: 1024,
				Y: 1024,
			},
			Queue: &r,
		})
		if page == "KitchenX" {
			kitchenX(gtx, th)
		} else {
			wid.Root(gtx)
		}
	}
}

func BenchmarkKitchenX(b *testing.B) {
	DoLayout(b, "KitchenX", 1)
}

func BenchmarkKitchenV(b *testing.B) {
	DoLayout(b, "KitchenV", 1)
}

/*
// BenchmarkGrid1 tests with 1 person in the data table
func Benchmark1(b *testing.B) {
	DoLayout(b, "Grid1", 1)
}

// Benchmark35 tests with 35 person in the data table - enough to fill the screen
func Benchmark35(b *testing.B) {
	DoLayout(b, "Grid2", 35)
}

// Benchmark10000 tests with 10000 person in the data table.
func Benchmark10000(b *testing.B) {
	DoLayout(b, "Grid2", 10000)
}

func BenchmarkButtons(b *testing.B) {
	DoLayout(b, "Buttons", 1)
}

func BenchmarkDropdown(b *testing.B) {
	DoLayout(b, "DropDown", 1)
}
*/
