// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"time"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
)

// Focuser implements the set/move focus
type Focuser interface {
	Next() Focuser
	SetNext(w Focuser)
	Prev() Focuser
	SetPrev(w Focuser)
	Focus()
	Disabled() bool
}

// Clickable represents a clickable area.
type Clickable struct {
	click    gesture.Click
	handler  func(v bool)
	eventKey int
	clicks   []Click
	// prevClicks is the index into clicks that marks the clicks
	// from the most recent Layout call. prevClicks is used to keep
	// clicks bounded.
	prevClicks int
	history    []Press
	// All clickable widgets should be able to have focus. This is needed for
	// usability reasons - all widgets should support keyboard only operation
	focused      bool
	disabled     bool
	requestFocus bool
	next         *Focuser
	prev         *Focuser
	index        int
}

// Click represents a click.
type Click struct {
	Modifiers key.Modifiers
	NumClicks int
}

// Press represents a past pointer press.
type Press struct {
	// Position of the press.
	Position f32.Point
	// Start is when the press began.
	Start time.Time
	// End is when the press was ended by a release or cancel.
	// A zero End means it hasn't ended yet.
	End time.Time
	// Cancelled is true for cancelled presses.
	Cancelled bool
}

// Global value used to save the previous widget that was tab-able
var prev Focuser


func (c *Clickable) Disabled() bool {
	return c.disabled
}

func (c *Clickable) SetupTabs() {
	if prev != nil {
		c.SetPrev(prev)
		prev.SetNext(c)
	}
	prev = c
}

func makeClickable() *Clickable {
	c := new(Clickable)
	if prev != nil {
		c.SetPrev(prev)
		prev.SetNext(c)
	}
	prev = c
	return c
}

func (c *Clickable) HandleClick() {
	for c.Clicked() {
		if c.handler != nil {
			c.handler(true)
		}
	}
}

func (c *Clickable) HandleToggle(value *bool, changed *bool) {
	for c.Clicked() {
		if value != nil {
			*value = !*value
			if c.handler != nil {
				c.handler(*value)
			}
		}
		if changed != nil {
			*changed = true
		}
	}
}

func (c *Clickable) SetNext(f Focuser) {
	c.next = &f
}

func (c *Clickable) SetPrev(f Focuser) {
	c.prev = &f
}

func (c *Clickable) Next() Focuser {
	if &c == nil || c == nil || c.next == nil {
		return nil
	}
	b := *c.next
	if b.Disabled() {
		return b.Next()
	}
	return b
}

func (c *Clickable) Prev() Focuser {
	if &c == nil || c == nil || c.prev == nil {
		return nil
	}
	b := *c.prev
	if b.Disabled() {
		return b.Prev()
	}
	return b
}

func (c *Clickable) Focus() {
	c.requestFocus = true
}

// Focused returns whether the editor is focused or not.
func (c *Clickable) Focused() bool {
	return c.focused
}

// Click executes a simple programmatic click
func (c *Clickable) Click() {
	c.clicks = append(c.clicks, Click{
		Modifiers: 0,
		NumClicks: 1,
	})
}

// Clicked reports whether there are pending clicks as would be
// reported by Clicks. If so, Clicked removes the earliest click.
func (c *Clickable) Clicked() bool {
	if len(c.clicks) == 0 {
		return false
	}
	n := copy(c.clicks, c.clicks[1:])
	c.clicks = c.clicks[:n]
	if c.prevClicks > 0 {
		c.prevClicks--
	}
	return true
}

func (c *Clickable) HasClicks() bool {
	return len(c.clicks) > 0
}

// Hovered returns whether pointer is over the element.
func (c *Clickable) Hovered() bool {
	return c.click.Hovered()
}

// Pressed returns whether pointer is pressing the element.
func (c *Clickable) Pressed() bool {
	return c.click.Pressed()
}

// Clicks returns and clear the clicks since the last call to Clicks.
func (c *Clickable) Clicks() []Click {
	clicks := c.clicks
	c.clicks = nil
	c.prevClicks = 0
	return clicks
}

// History is the past pointer presses useful for drawing markers.
// History is retained for a short duration (about a second).
func (c *Clickable) History() []Press {
	return c.history
}

// Layout and update the button state
func (c *Clickable) LayoutClickable(gtx C) D {
	// Flush clicks from before the last update.
	n := copy(c.clicks, c.clicks[c.prevClicks:])
	c.clicks = c.clicks[:n]
	c.prevClicks = n

	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	c.click.Add(gtx.Ops)
	if c.HasClicks() {
		c.Focus()
	}
	for len(c.history) > 0 {
		h := c.history[0]
		if h.End.IsZero() || gtx.Now.Sub(h.End) < 1*time.Second {
			break
		}
		n := copy(c.history, c.history[1:])
		c.history = c.history[:n]
	}
	return D{Size: gtx.Constraints.Min}
}

// update the button state by processing events.
func (c *Clickable) HandleKeys(gtx C) D {
	for _, ev := range gtx.Events(&c.eventKey) {
		switch ke := ev.(type) {
		case key.FocusEvent:
			c.focused = ke.Focus
		case key.Event:
			if !c.focused || ke.State != key.Press {
				break
			}
			switch ke.Name {
			case key.NameEnter, key.NameReturn, key.NameSpace, key.NameEscape:
				c.clicks = append(c.clicks, Click{
					Modifiers: 0,
					NumClicks: 1,
				})
				if l := len(c.history); l > 0 {
					c.history[l-1].End = gtx.Now
				}
			case key.NameUpArrow:
				c.index--
			case key.NameDownArrow:
				c.index++
			case key.NameTab:
				if !ke.Modifiers.Contain(key.ModShift) {
					if c.Next() != nil {
						c.Next().Focus()
					}
				} else {
					if c.Prev() != nil {
						c.Prev().Focus()
					}
				}
			}
		}
	}
	pointer.PassOp{Pass: true}.Add(gtx.Ops)
	key.InputOp{Tag: &c.eventKey, Hint: 0}.Add(gtx.Ops)
	if c.requestFocus {
		key.FocusOp{Tag: &c.eventKey}.Add(gtx.Ops)
		key.SoftKeyboardOp{Show: false}.Add(gtx.Ops)
	}
	c.requestFocus = false
	return D{}
}

func (c *Clickable) HandleClicks(gtx C) D {
	for _, e := range c.click.Events(gtx) {
		switch e.Type {
		case gesture.TypeClick:
			c.requestFocus = true
			c.clicks = append(c.clicks, Click{
				Modifiers: e.Modifiers,
				NumClicks: e.NumClicks,
			})
			if l := len(c.history); l > 0 {
				c.history[l-1].End = gtx.Now
			}
		case gesture.TypeCancel:
			for i := range c.history {
				c.history[i].Cancelled = true
				if c.history[i].End.IsZero() {
					c.history[i].End = gtx.Now
				}
			}
		case gesture.TypePress:
			c.history = append(c.history, Press{
				Position: e.Position,
				Start:    gtx.Now,
			})
		}
	}
	return D{}
}
