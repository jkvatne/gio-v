// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/op/paint"
	"gioui.org/text"
	"golang.org/x/exp/constraints"
	"image"
	"image/color"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"

	"gioui.org/layout"
	"gioui.org/unit"
)

type (
	// C is a shortcut for layout.Context
	C = layout.Context
	// D is a shortcut for layout.Dimensions
	D   = layout.Dimensions
	Wid = layout.Widget
)

// UIState is the hovered/focused etc. state
type UIState uint8

var (
	mouseX        int
	mouseY        int
	WinX          int
	WinY          int
	startWinY     int
	FixedFontSize bool
	GuiLock       sync.RWMutex
	invalidate    chan struct{}
)

// Base is tha base structure for widgets. It contains variables that (almost) all widgets share
type Base struct {
	th           *Theme
	hint         string
	padding      layout.Inset
	margin       layout.Inset
	onUserChange func()
	disabler     *bool
	width        unit.Dp
	role         UIRole
	cornerRadius unit.Dp
	borderWidth  unit.Dp
	fgColor      *color.NRGBA
	bgColor      *color.NRGBA
	Font         *font.Font
	FontScale    float64
	DpNo         *int
	Alignment    text.Alignment
}

// Fg returns the foreground color of a widget, either from
// its role or from the widget specific fgColor field.
func (wid *Base) Fg() color.NRGBA {
	if wid.fgColor == nil {
		return wid.th.Fg[wid.role]
	} else {
		return *wid.fgColor
	}
}

// Bg returns the background color of a widget, either from
// its role or from the widget specific bgColor field.
func (wid *Base) Bg() color.NRGBA {
	if wid.bgColor == nil {
		return wid.th.Bg[wid.role]
	} else {
		return *wid.bgColor
	}
}

// CheckDisabler is used when a variable controls the disabling of a widget.
func (wid *Base) CheckDisabler(gtx C) {
	if wid.disabler != nil {
		GuiLock.RLock()
		if *wid.disabler {
			gtx = gtx.Disabled()
		}
		GuiLock.RUnlock()
	}
}

// UpdateMousePos must be called from the main program in order to get mouse
// position and window size. They are needed to avoid that the tooltip
// is outside the window frame
func UpdateMousePos(gtx C, win *app.Window) {
	eventArea := clip.Rect(image.Rect(0, 0, 99999, 99999)).Push(gtx.Ops)
	// pointer.InputOp{Kinds: pointer.Move,Tag:   win,}.Add(gtx.Ops)
	event.Op(gtx.Ops, win)
	eventArea.Pop()
	for {
		if e, ok := gtx.Event(pointer.Filter{Target: win, Kinds: pointer.Move}); ok {
			if ev, ok := e.(pointer.Event); ok {
				mouseX = int(ev.Position.X)
				mouseY = int(ev.Position.Y)
			}
		} else {
			break
		}
	}
}

// Invalidate will force a redraw of the current form
func Invalidate() {
	invalidate <- struct{}{}
}

// Run is the main event handler, called with "go run" from main()
func Run(win *app.Window, mainForm *layout.Widget, th *Theme) {
	invalidate = make(chan struct{})
	go func() {
		for range invalidate {
			win.Invalidate()
		}
	}()

	for {
		switch e := win.NextEvent().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			var ops op.Ops
			// Save window size for use by widgets. Must be done before drawing
			WinX = e.Size.X
			WinY = e.Size.Y
			gtx := app.NewContext(&ops, e)

			if startWinY == 0 {
				startWinY = WinY
			}
			// Font size is in units sp (like dp but for fonts) while WinY is in pixels
			// So we have to rescale using PxToSp
			if !FixedFontSize {
				scale := float32(WinY) / float32(startWinY) * th.Scale
				gtx.Metric.PxPerDp = scale * gtx.Metric.PxPerDp
				gtx.Metric.PxPerSp = scale * gtx.Metric.PxPerSp
			}
			paint.ColorOp{Color: th.Bg[Surface]}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Draw widgets
			GuiLock.Lock()
			GuiLock.Unlock()
			ctx := gtx
			if dialog != nil {
				ctx = gtx.Disabled()
			}
			(*mainForm)(ctx)
			if dialog != nil {
				dialog(gtx)
			}

			// A hack to fetch mouse position and window size, so we can avoid
			// tooltips going outside the main window area
			p := pointer.PassOp{}.Push(gtx.Ops)
			UpdateMousePos(gtx, win)
			p.Pop()

			// Apply the actual screen drawing
			e.Frame(gtx.Ops)
		}
	}
}

// Min is a generic minimum function. Can be removed when go includes it
func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max is a generic maximum function. Can be removed when go includes it
func Max[T constraints.Ordered](x, y T) T {
	if x >= y {
		return x
	}
	return y
}

// Clamp will return the first argument clamped between argument 2 and 3.
func Clamp[T constraints.Ordered](v T, lo T, hi T) T {
	if v < lo {
		return lo
	} else if v > hi {
		return hi
	}
	return v
}
