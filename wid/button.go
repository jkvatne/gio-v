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
	Button       widget.Clickable
	shadow       ShadowStyle
	disabled     bool
	disabler     *bool
	Text         string
	ToolTipWidth unit.Dp
	Font         text.Font
	shaper       text.Shaper
	Icon         *widget.Icon
	Style        ButtonStyle
	fg           color.NRGBA
	bg           color.NRGBA
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
	b := ButtonDef{} // b.SetupTabs()
	// Setup default values
	b.th = th
	b.Text = label
	b.Font = text.Font{Weight: text.Medium}
	b.shaper = th.Shaper
	b.Style = style
	b.align = layout.Middle
	// Apply standard padding on the outside of the button. Can be overridden by option function
	b.padding = th.ButtonPadding
	for _, option := range options {
		option.apply(&b)
	}
	b.Tooltip = PlatformTooltip(th, b.hint)
	return &b
}

// HandleClick will call the callback function
func (b *ButtonDef) HandleClick() {
	for b.Button.Clicked() {
		if b.handler != nil {
			b.handler()
		}
	}
}

// Layout will draw a button defined in b.
func (b *ButtonDef) Layout(gtx C) D {
	if b.Style == Contained || b.Style == Round {
		b.fg = b.th.OnPrimary
		b.bg = b.th.Primary
	} else {
		b.fg = b.th.OnBackground
		b.bg = color.NRGBA{}
	}
	if b.Base.fgColor.A != 0 {
		b.bg = b.Base.fgColor
		if Luminance(b.Base.fgColor) > 127 {
			b.fg = RGB(0x000000)
		} else {
			b.fg = RGB(0xFFFFFF)
		}
	}
	dims := b.layout(gtx)
	b.HandleClick()
	pointer.CursorPointer.Add(gtx.Ops)
	return dims
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

func drawInk(gtx C, c widget.Press) {
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
	{
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
	}

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
	// Cover the entire constraints min rectangle and
	// apply curve values to size and color.
	size = int(float32(size) * 2 * float32(math.Sqrt(2)) * sizeBezier)
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(col*0xff)
	rgba := MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := size / 2
	defer op.Offset(c.Position.Add(image.Point{
		X: -rr,
		Y: -rr,
	})).Push(gtx.Ops).Pop()
	defer clip.UniformRRect(image.Rectangle{Max: image.Pt(size, size)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func (b *ButtonDef) layoutBackground() func(gtx C) D {
	return func(gtx C) D {
		rr := gtx.Dp(b.th.ButtonCornerRadius)
		if rr > gtx.Constraints.Min.Y/2 {
			rr = gtx.Constraints.Min.Y / 2
		}

		if b.Style == Round {
			rr = gtx.Constraints.Min.Y / 2
		}
		if b.Button.Focused() || b.Button.Hovered() {
			e := b.th.Elevation
			Shadow(rr, gtx.Dp(e)).Layout(gtx)
		}
		outline := image.Rectangle{Max: image.Point{
			X: gtx.Constraints.Min.X,
			Y: gtx.Constraints.Min.Y,
		}}
		defer clip.UniformRRect(outline, rr).Push(gtx.Ops).Pop()

		switch {
		case b.Style == Text && gtx.Queue == nil:
			// Disabled Outlined button. Text and outline is grey when disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Text && (b.Button.Hovered() || b.Button.Focused()):
			// Outline button, hovered/focused
			paint.FillShape(gtx.Ops, Hovered(b.th.Background), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Text:
			// Outline button, not disabled, keep transparent.
		case b.Style == Outlined && gtx.Queue == nil:
			// Disabled Outlined button. Text and outline is grey when disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			paintBorder(gtx, outline, Disabled(b.fg), b.th.BorderThickness, rr)
		case b.Style == Outlined && (b.Button.Hovered() || b.Button.Focused()):
			// Outline button, hovered/focused
			paint.FillShape(gtx.Ops, Hovered(b.th.Background), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			paintBorder(gtx, outline, b.fg, b.th.BorderThickness, rr)
		case b.Style == Outlined:
			// Outline button, not disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			paintBorder(gtx, outline, b.fg, b.th.BorderThickness, rr)
		case (b.Style == Contained || b.Style == Round) && gtx.Queue == nil:
			// Disabled contained/round button.
			paint.FillShape(gtx.Ops, Disabled(b.bg), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case (b.Style == Contained || b.Style == Round) && (b.Button.Hovered() || b.Button.Focused()):
			// Hovered or focused contained/round button.
			paint.FillShape(gtx.Ops, Hovered(b.bg), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Contained || b.Style == Round:
			// Contained/round button, not disabled
			paint.FillShape(gtx.Ops, b.bg, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		}
		for _, pressed := range b.Button.History() {
			drawInk(gtx, pressed)
		}
		return D{Size: gtx.Constraints.Min}
	}
}

func layLabel(b *ButtonDef) layout.Widget {
	if b.Text == "" {
		return func(gtx C) D { return D{} }
	}
	return func(gtx C) D {
		return b.th.ButtonLabelPadding.Layout(gtx, func(gtx C) D {
			switch {
			case (b.Style == Text || b.Style == Outlined) && gtx.Queue == nil:
				paint.ColorOp{Color: Disabled(b.fg)}.Add(gtx.Ops)
			default:
				paint.ColorOp{Color: b.fg}.Add(gtx.Ops)
			}
			return widget.Label{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.th.TextSize, b.Text)
		})
	}
}

func layIcon(b *ButtonDef) layout.Widget {
	if b.Icon != nil {
		return func(gtx C) D {
			inset := b.th.IconInset
			if b.Icon != nil && b.Text != "" {
				// Avoid large gap between icon and text when both are present
				inset.Right = unit.Dp(0)
			}
			return inset.Layout(gtx, func(gtx C) D {
				size := gtx.Dp(b.th.IconSize)
				gtx.Constraints = layout.Exact(image.Pt(size, size))
				return b.Icon.Layout(gtx, b.fg)
			})
		}
	}
	return func(gtx C) D { return D{} }
}

func (b *ButtonDef) layout(gtx C) D {
	return b.padding.Layout(gtx, func(gtx C) D {
		return b.Button.Layout(gtx,
			func(gtx C) D {
				return b.Tooltip.Layout(gtx, b.hint, func(gtx C) D {
					b.disabled = false
					if b.disabler != nil && *b.disabler {
						gtx = gtx.Disabled()
						b.disabled = true
					}
					min := CalcMin(gtx, b.width)
					return layout.Stack{Alignment: layout.Center}.Layout(gtx,
						layout.Expanded(b.layoutBackground()),
						layout.Stacked(
							func(gtx C) D {
								gtx.Constraints.Min = min
								if gtx.Constraints.Max.X < min.X {
									gtx.Constraints.Max.X = min.X
								}
								if b.Icon != nil {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceEnd}.Layout(
										gtx,
										layout.Rigid(layIcon(b)),
										layout.Rigid(layLabel(b)),
									)
								}
								return layLabel(b)(gtx)
							}),
					)
				})
			})
	})
}
