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
	prevClicks    int
	history       []Press
	focused       bool
	requestClicks int
	pressed       bool
	pressedKey    key.Name
	// ClickMovesFocus can be set true if you want clicking on a button
	// to move focus. If false, only Tab will move focus.
	// Dropdowns must have this set to true
	ClickMovesFocus bool
	index           *int
	maxIndex        int
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
func (b *Clickable) HandleEvents(gtx C) {
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
		case gesture.KindCancel:
			for i := range b.history {
				b.history[i].Cancelled = true
				if b.history[i].End.IsZero() {
					b.history[i].End = gtx.Now
				}
			}
		case gesture.KindPress:
			if e.Source == pointer.Mouse {
				gtx.Execute(key.FocusCmd{Tag: b})
			}
			b.history = append(b.history, Press{
				Position: e.Position,
				Start:    gtx.Now,
			})
		}
	}

	for {
		e, ok := gtx.Event(
			key.FocusFilter{Target: b},
			key.Filter{Focus: b, Name: key.NameReturn},
			key.Filter{Focus: b, Name: key.NameSpace},
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
			if !gtx.Focused(b) {
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
				// Only handle release from same key as was pressed
				if b.pressedKey != e.Name {
					break
				}
				b.pressed = false
				b.pressedKey = ""
				// Clicking via keyboard
				if e.Name == key.NameSpace || e.Name == key.NameReturn {
					if l := len(b.history); l > 0 {
						b.history[l-1].End = gtx.Now
					}
					b.clicks = append(b.clicks, Click{Modifiers: e.Modifiers, NumClicks: 1})
				} else if e.Name == key.NameDownArrow || e.Name == key.NameRightArrow {
					GuiLock.Lock()
					*b.index++
					GuiLock.Unlock()
				} else if e.Name == key.NameUpArrow || e.Name == key.NameLeftArrow {
					GuiLock.Lock()
					*b.index--
					if *b.index < 0 {
						*b.index = 0
					}
					GuiLock.Unlock()
				}
			}
		}
	}
}
