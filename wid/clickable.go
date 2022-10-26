package wid

import (
	"image"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
)

// Clickable represents a clickable area.
type Clickable struct {
	click  gesture.Click
	clicks []Click
	// prevClicks is the index into clicks that marks the clicks
	// from the most recent Layout call. prevClicks is used to keep
	// clicks bounded.
	prevClicks int
	history    []Press
	keyTag     struct{}
	focused    bool
	pressed    bool
	index      *int
}

// Click represents a click.
type Click struct {
	Modifiers key.Modifiers
	NumClicks int
}

// Press represents a past pointer press.
type Press struct {
	// Position of the press.
	Position image.Point
	// Start is when the press began.
	Start time.Time
	// End is when the press was ended by a release or cancel.
	// A zero End means it hasn't ended yet.
	End time.Time
	// Cancelled is true for cancelled presses.
	Cancelled bool
}

// Click executes a simple programmatic click
func (b *Clickable) Click() {
	b.clicks = append(b.clicks, Click{
		Modifiers: 0,
		NumClicks: 1,
	})
}

// Clicked reports whether there are pending clicks as would be
// reported by Clicks. If so, Clicked removes the earliest click.
func (b *Clickable) Clicked() bool {
	if len(b.clicks) == 0 {
		return false
	}
	n := copy(b.clicks, b.clicks[1:])
	b.clicks = b.clicks[:n]
	if b.prevClicks > 0 {
		b.prevClicks--
	}
	return true
}

// Hovered reports whether a pointer is over the element.
func (b *Clickable) Hovered() bool {
	return b.click.Hovered()
}

// Pressed reports whether a pointer is pressing the element.
func (b *Clickable) Pressed() bool {
	return b.click.Pressed()
}

// Focused reports whether b has focus.
func (b *Clickable) Focused() bool {
	return b.focused
}

// Clicks returns and clear the clicks since the last call to Clicks.
func (b *Clickable) Clicks() []Click {
	clicks := b.clicks
	b.clicks = nil
	b.prevClicks = 0
	return clicks
}

// History is the past pointer presses useful for drawing markers.
// History is retained for a short duration (about a second).
func (b *Clickable) History() []Press {
	return b.history
}

func (b *Clickable) SetupEventHandlers(gtx C, size image.Point) {
	defer clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops).Pop()
	disabled := gtx.Queue == nil
	semantic.DisabledOp(disabled).Add(gtx.Ops)
	b.click.Add(gtx.Ops)
	if !disabled {
		keys := key.Set("")
		if b.focused {
			keys = key.Set("⏎|Space|←|→|↑|↓")
		}
		key.InputOp{Tag: &b.keyTag, Keys: keys}.Add(gtx.Ops)
	} else {
		b.focused = false
	}
	for len(b.history) > 0 {
		c := b.history[0]
		if c.End.IsZero() || gtx.Now.Sub(c.End) < 1*time.Second {
			break
		}
		n := copy(b.history, b.history[1:])
		b.history = b.history[:n]
	}
}

// HandleEvents the button state by processing events.
func (b *Clickable) HandleEvents(gtx layout.Context) {
	// Flush clicks from before the last HandleEvents.
	n := copy(b.clicks, b.clicks[b.prevClicks:])
	b.clicks = b.clicks[:n]
	b.prevClicks = n

	for _, e := range b.click.Events(gtx) {
		switch e.Type {
		case gesture.TypeClick:
			b.clicks = append(b.clicks, Click{
				Modifiers: e.Modifiers,
				NumClicks: e.NumClicks,
			})
			if l := len(b.history); l > 0 {
				b.history[l-1].End = gtx.Now
			}
		case gesture.TypeCancel:
			for i := range b.history {
				b.history[i].Cancelled = true
				if b.history[i].End.IsZero() {
					b.history[i].End = gtx.Now
				}
			}
		case gesture.TypePress:
			if e.Source == pointer.Mouse {
				key.FocusOp{Tag: &b.keyTag}.Add(gtx.Ops)
			}
			b.history = append(b.history, Press{
				Position: e.Position,
				Start:    gtx.Now,
			})
		}
	}
	for _, e := range gtx.Events(&b.keyTag) {
		switch e := e.(type) {
		case key.FocusEvent:
			b.focused = e.Focus
		case key.Event:
			if (e.Name == key.NameSpace || e.Name == key.NameReturn) && b.focused {
				if e.State == key.Press && !b.pressed {
					b.history = append(b.history, Press{
						Position: image.Point{0, 0},
						Start:    gtx.Now,
					})
					b.pressed = true
				} else if e.State == key.Release {
					b.pressed = false
					b.clicks = append(b.clicks, Click{Modifiers: e.Modifiers, NumClicks: 1})
					if l := len(b.history); l > 0 {
						b.history[l-1].End = gtx.Now
					}
				}
			} else if b.index != nil && e.State == key.Release {
				if e.Name == key.NameDownArrow || e.Name == key.NameRightArrow {
					*b.index++
				} else if e.Name == key.NameUpArrow || e.Name == key.NameLeftArrow {
					*b.index--
				}
			}
		}
	}
}
