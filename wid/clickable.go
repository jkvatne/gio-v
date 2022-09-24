// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"time"

	"gioui.org/op/clip"

	"gioui.org/layout"

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
	index        *int
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

// Prev is a global value used to save the previous widget that was tab-able
var Prev Focuser

// First will store the initial focused widget
var First Focuser

// Root is the root widget (usually a list), and is the root of the widget tree
var Root layout.Widget

// Init will initialize the widget tree
func Init() {
	// Drop all old data
	Root = nil
	First = nil
	Prev = nil
}

// Setup will initialze root
func Setup(w layout.Widget) {
	Root = w
	// Focus the component to be focused at startup.
	if First != nil {
		First.Focus()
	}
}

// Disabled returns true if the widget is disabled
func (c *Clickable) Disabled() bool {
	return c.disabled
}

// SetupTabs is used to set next/previous widget for tabbing between widgets
func (c *Clickable) SetupTabs() {
	if First == nil {
		First = c
	}
	if Prev != nil {
		c.SetPrev(Prev)
		Prev.SetNext(c)
	}
	Prev = c
}

// HandleClick will call the callback function
func (c *Clickable) HandleClick() {
	for c.Clicked() {
		if c.handler != nil {
			c.handler(true)
		}
	}
}

// HandleToggle will support toggling widgets
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

// SetNext set pointer for next widget
func (c *Clickable) SetNext(f Focuser) {
	c.next = &f
}

// SetPrev set pointer for previous widget
func (c *Clickable) SetPrev(f Focuser) {
	c.prev = &f
}

// Next returns pointer to next widget
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

// Prev returns pointer to previous widget
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

// Focus will request focus for this widget
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

// HasClicks is true if there are more clicks available
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

// LayoutClickable and update the button state
func (c *Clickable) LayoutClickable(gtx C) D {
	// Flush clicks from before the last update.
	n := copy(c.clicks, c.clicks[c.prevClicks:])
	c.clicks = c.clicks[:n]
	c.prevClicks = n
	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops).Pop()

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

// HandleKeys updates the button state by processing events.
func (c *Clickable) HandleKeys(gtx C) bool {
	var newKey bool
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
				newKey = true
			case key.NameUpArrow, key.NameLeftArrow:
				if !ke.Modifiers.Contain(key.ModCtrl) {
					*c.index--
				} else {
					*c.index -= 10
				}
				newKey = true
			case key.NameHome:
				*c.index = 0
				newKey = true
			case key.NameEnd:
				*c.index = 100
				newKey = true
			case key.NameDownArrow, key.NameRightArrow:
				if !ke.Modifiers.Contain(key.ModCtrl) {
					*c.index++
				} else {
					*c.index += 10
				}
				newKey = true
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
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	key.InputOp{Tag: &c.eventKey, Hint: 0}.Add(gtx.Ops)
	if c.requestFocus {
		key.FocusOp{Tag: &c.eventKey}.Add(gtx.Ops)
		key.SoftKeyboardOp{Show: false}.Add(gtx.Ops)
	}
	c.requestFocus = false
	return newKey
}

// HandleClicks will handle gestures
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

// LayoutBorder will draw a border around the widget
func LayoutBorder(e *Clickable, th *Theme) func(gtx C) D {
	return func(gtx C) D {
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X,
			Y: gtx.Constraints.Min.Y,
		}}
		if e.Focused() {
			paintBorder(gtx, outline, MulAlpha(th.Primary, 255), th.BorderThicknessActive, th.CornerRadius)
		} else if e.Hovered() {
			paintBorder(gtx, outline, MulAlpha(th.Primary, 140), th.BorderThickness, th.CornerRadius)
		} else {
			paintBorder(gtx, outline, MulAlpha(th.Primary, 50), th.BorderThickness, th.CornerRadius)
		}
		return D{}
	}
}
