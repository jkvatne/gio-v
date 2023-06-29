// SPDX-License-Identifier: Unlicense OR MIT

// Package wid is an alternative implementation of gio's material widgets
package wid

import (
	"gioui.org/io/pointer"
	"image"
	"image/color"
	"math"

	"gioui.org/io/semantic"

	"gioui.org/unit"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
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
	// Header is used in tables to make them clickable
	Header
)

type StrValue interface {
	string | *string
}

// ButtonDef is the struct for buttons
type ButtonDef struct {
	Base
	Tooltip
	Clickable
	Text      *string
	Icon      *Icon
	Style     ButtonStyle
	internPad layout.Inset
}

// BtnOption is the options for buttons only
type BtnOption func(*ButtonDef)

// RoundButton is a shortcut to a round button
func RoundButton(th *Theme, d *Icon, options ...Option) layout.Widget {
	options = append([]Option{Role(Primary), BtnIcon(d), W(0), RR(99999)}, options...)
	return aButton(Round, th, "", options...).Layout
}

// TextButton is a shortcut to a text only button
func TextButton(th *Theme, label string, options ...Option) layout.Widget {
	options = append([]Option{Role(Canvas)}, options...)
	return aButton(Text, th, label, options...).Layout
}

// OutlineButton is a shortcut to an outlined button
func OutlineButton(th *Theme, label string, options ...Option) layout.Widget {
	options = append([]Option{Role(Canvas)}, options...)
	return aButton(Outlined, th, label, options...).Layout
}

// Button is the generic button selector. Defaults to primary
func Button[V StrValue](th *Theme, label V, options ...Option) layout.Widget {
	return aButton(Contained, th, label, options...).Layout
}

// HeaderButton is a shortcut to a text only button with left justified text and a given size
func HeaderButton(th *Theme, label string, options ...Option) layout.Widget {
	options = append([]Option{Role(Canvas)}, options...)
	b := aButton(Text, th, label, options...)
	b.cornerRadius = 0
	b.internPad = th.OutsidePadding
	b.Style = Header
	b.padding = layout.Inset{}
	return b.Layout
}

func aButton[V StrValue](style ButtonStyle, th *Theme, label V, options ...Option) *ButtonDef {
	b := ButtonDef{}
	// Setup default values
	b.th = th
	b.role = Primary
	if x, ok := any(label).(string); ok {
		b.Text = &x
	}
	if x, ok := any(label).(*string); ok {
		b.Text = x
	}
	b.Font = &th.DefaultFont
	b.shaper = th.Shaper
	b.Style = style
	b.internPad = th.ButtonLabelPadding
	// Apply standard padding on the outside of the button. Can be overridden by option function
	b.padding = th.ButtonPadding
	b.FontSize = 1.0
	b.cornerRadius = th.ButtonCornerRadius
	for _, option := range options {
		option.apply(&b)
	}
	b.Tooltip = PlatformTooltip(th)
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
	dims := b.padding.Layout(gtx,
		func(gtx C) D {
			if b.disabler != nil && !*b.disabler {
				gtx.Queue = nil
			}
			// Handle clickable pointer/keyboard inputs
			b.HandleEvents(gtx)
			b.HandleClick()
			dims := b.layout(gtx)
			b.SetupEventHandlers(gtx, dims.Size)
			dims = b.Tooltip.Layout(gtx, b.hint, func(gtx C) D {
				return dims
			})
			return dims
		})
	pointer.CursorPointer.Add(gtx.Ops)

	return dims
}

func (b *ButtonDef) layout(gtx C) D {
	// Make macro with text color
	recorder := op.Record(gtx.Ops)
	paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
	colorMacro := recorder.Stop()
	// Allow zero size for text.
	cgtx := gtx
	cgtx.Constraints.Min.X = 0
	cgtx.Constraints.Min.Y = 0
	// Render text to find button width
	recorder = op.Record(gtx.Ops)
	textDim := widget.Label{Alignment: text.Start}.Layout(cgtx, b.shaper, *b.Font, b.th.TextSize*unit.Sp(b.FontSize), *b.Text, colorMacro)
	textMacro := recorder.Stop()
	// Icon size is equal to label height *1.333
	iconSize := 0
	if b.Icon != nil {
		iconSize = textDim.Size.Y * 4 / 3
	}
	height := Min(textDim.Size.Y+gtx.Dp(b.internPad.Top)+gtx.Dp(b.internPad.Bottom), gtx.Constraints.Max.Y)
	// Default button width when width is not given has padding=1.2 char heights.
	contentWidth := textDim.Size.X + iconSize
	width := Min(gtx.Constraints.Max.X, Max(contentWidth+textDim.Size.Y, gtx.Dp(b.width)))
	if b.Style == Header {
		width = gtx.Constraints.Max.X
	}

	if b.Style == Round {
		width = height
	}
	// Limit corner radius
	rr := Min(gtx.Dp(b.cornerRadius), height/2)

	outline := image.Rect(0, 0, width, height)

	// Draw shadow if pressed. Must be done before clipping
	// because the shadow is outside the button
	if b.Clickable.Focused() {
		DrawShadow(gtx, outline, rr, 20)
	}
	defer clip.UniformRRect(outline, rr).Push(gtx.Ops).Pop()

	if b.Style == Outlined {
		paintBorder(gtx, outline, b.th.Fg(Outline), b.th.BorderThickness, rr)
	} else if b.Style != Text && gtx.Queue == nil {
		paint.Fill(gtx.Ops, Disabled(b.Bg()))
	} else if b.Style != Text && b.Style != Header {
		paint.Fill(gtx.Ops, b.Bg())
	}
	if b.Clickable.Focused() && b.Clickable.Hovered() {
		paint.Fill(gtx.Ops, MulAlpha(b.Fg(), 30))
	} else if b.Clickable.Focused() {
		paint.Fill(gtx.Ops, MulAlpha(b.Fg(), 20))
	} else if b.Clickable.Hovered() {
		paint.Fill(gtx.Ops, MulAlpha(b.Fg(), 15))
	}

	semantic.DisabledOp(gtx.Queue == nil).Add(gtx.Ops)

	// Icon context
	cgtx.Constraints.Min = image.Point{X: width, Y: height}
	for _, pressed := range b.Clickable.History() {
		drawInk(cgtx, pressed)
	}
	cgtx.Constraints.Min = image.Point{X: iconSize, Y: iconSize}

	// Calculate internal paddings and move
	dy := Max(0, (height-textDim.Size.Y)/2)
	dx := Max(0, (width-contentWidth)/2)
	if b.internPad.Left > 0 && dx < gtx.Dp(b.internPad.Left) {
		dx = gtx.Dp(b.internPad.Left)
	}
	if b.Style == Header {
		dx = gtx.Dp(b.internPad.Left)
	}

	if b.Icon != nil && *b.Text != "" {
		// Draw Icon at given offset
		ofs := op.Offset(image.Pt(dx, (height-iconSize)/2)).Push(gtx.Ops)
		_ = b.Icon.Layout(cgtx, b.Fg())
		ofs.Pop()
		// Draw text at given offset
		ofs = op.Offset(image.Pt(dx+iconSize, dy)).Push(gtx.Ops)
		paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
		textMacro.Add(gtx.Ops)
		ofs.Pop()
	} else if b.Icon != nil {
		// Icon only
		dx := (height - iconSize) / 2
		defer op.Offset(image.Pt(dx, dx)).Push(gtx.Ops).Pop()
		_ = b.Icon.Layout(cgtx, b.Fg())
	} else {
		// Text only
		defer op.Offset(image.Pt(dx, dy)).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: b.Fg()}.Add(gtx.Ops)
		textMacro.Add(gtx.Ops)
	}
	return D{Size: outline.Max}
}

func (b BtnOption) apply(cfg interface{}) {
	b(cfg.(*ButtonDef))
}

// BtnIcon sets button icon
func BtnIcon(i *Icon) BtnOption {
	return func(b *ButtonDef) {
		b.Icon = i
	}
}

// RR is the corner radius
func RR(rr unit.Dp) BtnOption {
	return func(b *ButtonDef) {
		b.cornerRadius = rr
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
