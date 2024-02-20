package wid

import (
	"gioui.org/io/event"
	"image"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
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

	keyTag       struct{}
	requestFocus bool
	focused      bool

	requestClicks int
	pressed       bool
	pressedKey    key.Name
	// ClickMovesFocus can be set true if you want clicking on a button
	// to move focus. If false, only Tab will move focus.
	// Dropdowns must have this set to true
	ClickMovesFocus bool
	index           *int
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

// Focus requests the input focus for the element.
func (b *Clickable) Focus() {
	b.requestFocus = true
}

// History is the past pointer presses useful for drawing markers.
// History is retained for a short duration (about a second).
func (b *Clickable) History() []Press {
	return b.history
}

func (b *Clickable) SetupEventHandlers(gtx C, size image.Point) {
	defer clip.Rect(image.Rectangle{Max: size}).Push(gtx.Ops).Pop()
	b.click.Add(gtx.Ops)
	event.Op(gtx.Ops, b)
}

// HandleEvents the button state by processing events.
func (b *Clickable) HandleEvents(t event.Tag, gtx C) {
	for len(b.history) > 0 {
		c := b.history[0]
		if c.End.IsZero() || gtx.Now.Sub(c.End) < 1*time.Second {
			break
		}
		n := copy(b.history, b.history[1:])
		b.history = b.history[:n]
	}
	for {
		e, ok := b.click.Update(gtx.Source)
		if !ok {
			break
		}
		switch e.Kind {
		case gesture.KindClick:
			b.clicks = append(b.clicks, Click{
				Modifiers: e.Modifiers,
				NumClicks: e.NumClicks,
			})
			if l := len(b.history); l > 0 {
				b.history[l-1].End = gtx.Now
			}
			if b.ClickMovesFocus {
				b.Focus()
			}
		case gesture.KindCancel:
			for i := range b.history {
				b.history[i].Cancelled = true
				if b.history[i].End.IsZero() {
					b.history[i].End = gtx.Now
				}
			}
		case gesture.KindPress:
			if e.Source == pointer.Mouse {
				gtx.Execute(key.FocusCmd{Tag: t})
			}
			b.history = append(b.history, Press{
				Position: e.Position,
				Start:    gtx.Now,
			})
		}
	}

	for {
		e, ok := gtx.Event(
			key.FocusFilter{Target: t},
			key.Filter{Focus: t, Name: key.NameReturn},
			key.Filter{Focus: t, Name: key.NameSpace},
		)
		if !ok {
			break
		}
		switch e := e.(type) {
		case key.FocusEvent:
			if e.Focus {
				b.pressedKey = ""
			}
		case key.Event:
			if !gtx.Focused(t) {
				break
			}
			if e.Name != key.NameReturn && e.Name != key.NameSpace {
				break
			}
			switch e.State {
			case key.Press:
				if !b.pressed {
					b.pressedKey = e.Name
					b.history = append(b.history, Press{
						Start: gtx.Now,
					})
					b.pressed = true
				}
			case key.Release:
				if b.pressedKey != e.Name {
					break
				}
				b.pressed = false
				// only register a key as a click if the key was pressed and released while this button was focused
				b.pressedKey = ""
				if l := len(b.history); l > 0 {
					b.history[l-1].End = gtx.Now
				}
				b.clicks = append(b.clicks, Click{Modifiers: e.Modifiers, NumClicks: 1})
			}
		}
	}
}

func (b *Clickable) GetIndex(n int) int {
	GuiLock.RLock()
	idx := *b.index
	GuiLock.RUnlock()
	if idx < 0 {
		idx = 0
		GuiLock.Lock()
		*b.index = idx
		GuiLock.Unlock()
	} else if idx >= n {
		idx = n - 1
		GuiLock.Lock()
		*b.index = idx
		GuiLock.Unlock()
	}
	return idx
}
