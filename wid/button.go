// SPDX-License-Identifier: Unlicense OR MIT

// Package wid is an alternative implementation of gio's material widgets
package wid

import (
	"image"
	"image/color"
	"math"

	"gioui.org/io/pointer"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

// ButtonStyle indicates a Contained, Text, Outline or round button
type ButtonStyle int

const (
	// Contained is a solid, colored button
	Contained ButtonStyle = iota
	// Text is a button without outline or color. Just text
	Text
	// Outlined is a text button with outline
	Outlined
	// Round is a round button, usually with icon only
	Round
)

// ButtonDef is the struct for buttons
type ButtonDef struct {
	Base
	Tooltip
	Clickable
	shadow       ShadowStyle
	Text         string
	ToolTipWidth unit.Dp
	Font         text.Font
	shaper       text.Shaper
	Icon         *widget.Icon
	Style        ButtonStyle
	align        layout.Alignment
}

// BtnOption is the options for buttons only
type BtnOption func(*ButtonDef)

// RoundButton is a shortcut to a round button
func RoundButton(th *Theme, d *widget.Icon, options ...Option) layout.Widget {
	options = append(options, BtnIcon(d), W(0))
	return aButton(Round, th, "", options...).Layout
}

// TextButton is a shortcut to a text only button
func TextButton(th *Theme, label string, options ...Option) layout.Widget {
	return aButton(Text, th, label, options...).Layout
}

// OutlineButton is a shortcut to an outlined button
func OutlineButton(th *Theme, label string, options ...Option) layout.Widget {
	return aButton(Outlined, th, label, options...).Layout
}

// Button is the generic button selector
func Button(th *Theme, label string, options ...Option) layout.Widget {
	return aButton(Contained, th, label, options...).Layout
}

// HeaderButton is a shortcut to a text only button with left justified text and a given size
func HeaderButton(th *Theme, label string, options ...Option) layout.Widget {
	options = append(options, AlignLeft())
	return aButton(Text, th, label, options...).Layout
}

func aButton(style ButtonStyle, th *Theme, label string, options ...Option) *ButtonDef {
	b := ButtonDef{}
	// Setup default values
	b.th = th
	b.Text = label
	b.Font = text.Font{Weight: text.Medium}
	b.shaper = th.Shaper
	b.Style = style
	b.align = layout.Middle
	b.role = Undefined
	// Apply standard padding on the outside of the button. Can be overridden by option function
	b.padding = th.ButtonPadding
	for _, option := range options {
		option.apply(&b)
	}
	if b.Style == Outlined || b.Style == Text {
		b.fgColor = th.Fg(b.role)
		b.bgColor = color.NRGBA{}
	} else if (b.fgColor == color.NRGBA{}) && (b.bgColor == color.NRGBA{}) {
		if (b.role == Undefined) && (b.fgColor == color.NRGBA{}) {
			b.role = Primary
		}
		b.fgColor = th.Fg(b.role)
		b.bgColor = th.Bg(b.role)
	}
	b.Tooltip = PlatformTooltip(th, b.hint)
	return &b
}

// HandleClick will call the callback function
func (b *ButtonDef) HandleClick() {
	for b.Clickable.Clicked() {
		if b.onUserChange != nil {
			b.onUserChange()
		}
	}
}

// Layout will draw a button defined in b.
func (b *ButtonDef) Layout(gtx C) D {
	// Add an outer padding outside the button
	dims := b.padding.Layout(gtx, func(gtx C) D {
		// Handle clickable pointer/keyboard inputs
		b.HandleEvents(gtx)
		dims := b.Tooltip.Layout(gtx, b.hint, func(gtx C) D {
			return b.layout(gtx)
		})
		b.SetupEventHandlers(gtx, dims.Size)
		return dims
	})
	b.HandleClick()
	pointer.CursorPointer.Add(gtx.Ops)
	return dims
}

func (b *ButtonDef) layout(gtx C) D {
	// Render text to find button width
	macro := op.Record(gtx.Ops)
	cgtx := gtx
	cgtx.Constraints.Min.X = 0
	dims := widget.Label{Alignment: text.Start}.Layout(cgtx, b.shaper, b.Font, b.th.TextSize, b.Text)
	call := macro.Stop()
	height := 2 * dims.Size.Y
	width := height + dims.Size.X

	rr := gtx.Dp(b.th.ButtonCornerRadius)
	if rr > height/2 {
		rr = height / 2
	}
	dx := height / 2
	dy := dx / 2
	if b.Style == Round {
		rr = height / 2
		width = height
		dx = height / 4
		dy = height / 4
	}

	outline := image.Rect(0, 0, width, height)

	if b.Icon != nil && b.Text != "" {
		outline.Max.X += dims.Size.Y * 2
	}
	// Draw shadow if pressed. Must be done before cliping
	// because the shadow is outside the button
	if b.Clickable.Focused() {
		DrawShadow(gtx, outline, rr, 20)
	}
	defer clip.UniformRRect(outline, rr).Push(gtx.Ops).Pop()

	if b.Style == Outlined {
		paintBorder(gtx, outline, b.th.Fg(Outline), b.th.BorderThickness, rr)
	} else if b.Style != Text {
		paint.Fill(gtx.Ops, b.bgColor)
	}
	if b.Clickable.Focused() && b.Clickable.Hovered() {
		paint.Fill(gtx.Ops, MulAlpha(b.fgColor, 30))
	} else if b.Clickable.Focused() {
		paint.Fill(gtx.Ops, MulAlpha(b.fgColor, 20))
	} else if b.Clickable.Hovered() {
		paint.Fill(gtx.Ops, MulAlpha(b.fgColor, 15))
	}
	cgtx.Constraints.Min = image.Point{X: height / 2, Y: height / 2}
	for _, pressed := range b.Clickable.History() {
		drawInk(cgtx, pressed)
	}

	defer op.Offset(image.Pt(dx, dy)).Push(gtx.Ops).Pop()
	if b.Icon != nil {
		im := b.Icon.Layout(cgtx, b.fgColor)
		defer op.Offset(image.Pt(height/4+im.Size.X, 0)).Push(gtx.Ops).Pop()
		width += im.Size.X + height*2
	}
	paint.ColorOp{Color: b.fgColor}.Add(gtx.Ops)
	call.Add(gtx.Ops)
	return D{Size: outline.Max}
}

func (b BtnOption) apply(cfg interface{}) {
	b(cfg.(*ButtonDef))
}

// AlignLeft will align text to the left. Used for Text buttons.
func AlignLeft() BtnOption {
	return func(b *ButtonDef) {
		b.align = layout.Start
	}
}

// BtnIcon sets button icon
func BtnIcon(i *widget.Icon) BtnOption {
	return func(b *ButtonDef) {
		b.Icon = i
	}
}

// Disable is an optional parameter to set a bool variable to disable the button
func Disable(v *bool) BtnOption {
	return func(b *ButtonDef) {
		b.disabler = v
	}
}

func drawInk(gtx C, c Press) {
	// duration is the number of seconds for the completed animation:
	// expand while fading in, then out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)
	now := gtx.Now
	t := float32(now.Sub(c.Start).Seconds())
	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}
	endt := float32(end.Sub(c.Start).Seconds())
	// Compute the fade-in/out position in [0;1].
	var alphat float32
	var haste float32
	if c.Cancelled {
		// If the press was cancelled before the inkwell
		// was fully faded in, fast-forward the animation
		// to match the fade-out.
		if h := 0.5 - endt/fadeDuration; h > 0 {
			haste = h
		}
	}
	// Fade in.
	half1 := t/fadeDuration + haste
	if half1 > 0.5 {
		half1 = 0.5
	}
	// Fade out.
	half2 := float32(now.Sub(end).Seconds())
	half2 /= fadeDuration
	half2 += haste
	if half2 > 0.5 {
		// Too old.
		return
	}
	alphat = half1 + half2
	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration
	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := gtx.Constraints.Min.X
	if h := gtx.Constraints.Min.Y; h > size {
		size = h
	}
	size = 10 * int(float32(size)*2*float32(math.Sqrt(2))*sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	rgba := MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op.Offset(c.Position.Add(image.Point{X: -rr, Y: -rr})).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}
