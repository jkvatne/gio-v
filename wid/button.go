// SPDX-License-Identifier: Unlicense OR MIT
// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
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
	Widget
	Clickable
	Tooltip
	shadow       ShadowStyle
	disabler     *bool
	Text         string
	ToolTipWidth unit.Value
	Font         text.Font
	shaper       text.Shaper
	Icon         *Icon
	Style        ButtonStyle
	fg           color.NRGBA
	bg           color.NRGBA
}

// BtnOption is the options for buttons only
type BtnOption func(*ButtonDef)

// RoundButton is a shortcut to a round button
func RoundButton(th *Theme, d []byte, options ...Option) func(gtx C) D {
	options = append(options, BtnIcon(d), W(0))
	return aButton(Round, th, "", options...)
}

// TextButton is a shortcut to a text only button
func TextButton(th *Theme, label string, options ...Option) func(gtx C) D {
	return aButton(Text, th, label, options...)
}

// OutlineButton is a shortcut to an outlined button
func OutlineButton(th *Theme, label string, options ...Option) func(gtx C) D {
	return aButton(Outlined, th, label, options...)
}

// Button is the generic button selector
func Button(th *Theme, label string, options ...Option) func(gtx C) D {
	return aButton(Contained, th, label, options...)
}

func aButton(style ButtonStyle, th *Theme, label string, options ...Option) func(gtx C) D {
	b := ButtonDef{}
	b.SetupTabs()
	// Setup default values
	b.th = th
	b.Text = label
	b.Font = text.Font{Weight: text.Medium}
	b.shaper = th.Shaper
	b.Style = style
	// Apply default padding. Can be overridden by option function
	b.Pad(5, 2, 2, 2)
	for _, option := range options {
		option.apply(&b)
	}
	b.Tooltip = PlatformTooltip(th, b.hint)
	return func(gtx C) D {
		b.fg = th.OnPrimary
		b.bg = th.Primary
		if b.Widget.fgColor.A != 0 {
			b.bg = b.Widget.fgColor
			if Luminance(b.Widget.fgColor) > 127 {
				b.fg = RGB(0x000000)
			} else {
				b.fg = RGB(0xFFFFFF)
			}
		}
		dims := b.layout(gtx)
		b.HandleClick()
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

func (b BtnOption) apply(cfg interface{}) {
	b(cfg.(*ButtonDef))
}

// BtnIcon sets button icon
func BtnIcon(d []byte) BtnOption {
	i, _ := NewIcon(d)
	return func(b *ButtonDef) {
		b.Icon = i
	}
}

// Handler is an optional parameter to set a callback when the button is clicked
func Handler(f func()) BtnOption {
	foo := func(b bool) { f() }
	return func(b *ButtonDef) {
		b.handler = foo
	}
}

// Disable is an optional parameter to set a bool variable to disable the button
func Disable(v *bool) BtnOption {
	return func(b *ButtonDef) {
		b.disabler = v
	}
}

func drawInk(gtx C, c Press) {
	now := gtx.Now
	t := now.Sub(c.Start).Seconds()
	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}
	endTime := end.Sub(c.Start).Seconds()
	// Compute the fade-in/out position in [0;1].
	var haste float64
	if c.Cancelled {
		// If the press was cancelled before the inkwell
		// was fully faded-in, fast-forward the animation
		// to match the fade-out.
		if h := 0.5 - endTime/0.9; h > 0 {
			haste = h
		}
	}
	// Fade in.
	half1 := math.Max(t/0.9+haste, 0.5)
	if half1 > 0.5 {
		half1 = 0.5
	}
	// Fade out.
	half2 := now.Sub(end).Seconds()/0.9 + haste
	if half2 > 0.5 {
		return
	}
	alpha := half1 + half2
	// Compute the expanded position in [0;1].
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		t = endTime
	}
	sizet := math.Min(t*2, 1.0)
	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if alpha > .5 {
		// Start fadeout after half the animation.
		alpha = 1.0 - alpha
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alpha * 2
	size := math.Max(float64(gtx.Constraints.Min.Y), float64(gtx.Constraints.Min.X))
	alpha = 0.7 * alpha * 2 * t2 * (3.0 - 3.0*alpha)
	ba, bc := byte(alpha*0xff), byte(0x80)
	rgba := MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := float32(size * math.Sqrt(2.0) * sizet * sizet * (3.0 - 2.0*sizet))
	op.Offset(c.Position.Add(f32.Point{
		X: -rr,
		Y: -rr,
	})).Add(gtx.Ops)
	defer clip.UniformRRect(f32.Rectangle{Max: f32.Pt(2*rr, 2*rr)}, rr).Push(gtx.Ops).Pop()
	paint.PaintOp{}.Add(gtx.Ops)
}

func paintBorder(gtx C, outline f32.Rectangle, col color.NRGBA, width unit.Value, rr unit.Value) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, Pxr(gtx, rr)).Path(gtx.Ops),
			Width: Pxr(gtx, width),
		}.Op(),
	)
}

func (b *ButtonDef) layoutBackground() func(gtx C) D {
	return func(gtx C) D {
		b.LayoutClickable(gtx)
		b.HandleClicks(gtx)
		b.HandleKeys(gtx)

		rr := Pxr(gtx, b.th.CornerRadius)
		if b.Style == Round {
			rr = float32(gtx.Constraints.Min.Y) / 2.0
		}
		if b.Focused() || b.Hovered() {
			Shadow(unit.Px(rr), b.th.Elevation).Layout(gtx)
		}
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		defer clip.UniformRRect(outline, rr).Push(gtx.Ops).Pop()

		switch {
		case b.Style == Text && gtx.Queue == nil:
			// Disabled Outlined button. Text and outline is grey when disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Text && (b.Hovered() || b.Focused()):
			// Outline button, hovered/focused
			paint.FillShape(gtx.Ops, Hovered(b.th.Background), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Text:
			// Outline button, not disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Outlined && gtx.Queue == nil:
			// Disabled Outlined button. Text and outline is grey when disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			paintBorder(gtx, outline, Disabled(b.fg), b.th.BorderThickness, b.th.CornerRadius)
		case b.Style == Outlined && (b.Hovered() || b.Focused()):
			// Outline button, hovered/focused
			paint.FillShape(gtx.Ops, Hovered(b.th.Background), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			paintBorder(gtx, outline, b.bg, b.th.BorderThickness, b.th.CornerRadius)
		case b.Style == Outlined:
			// Outline button, not disabled
			paint.FillShape(gtx.Ops, b.th.Background, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
			paintBorder(gtx, outline, b.bg, b.th.BorderThickness, b.th.CornerRadius)
		case (b.Style == Contained || b.Style == Round) && gtx.Queue == nil:
			// Disabled contained/round button.
			paint.FillShape(gtx.Ops, Disabled(b.bg), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case (b.Style == Contained || b.Style == Round) && (b.Hovered() || b.Focused()):
			// Hovered or focused contained/round button.
			paint.FillShape(gtx.Ops, Hovered(b.bg), clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		case b.Style == Contained || b.Style == Round:
			// Contained/round button, not disabled
			paint.FillShape(gtx.Ops, b.bg, clip.RRect{Rect: outline, SE: rr, SW: rr, NW: rr, NE: rr}.Op(gtx.Ops))
		}
		for _, c := range b.History() {
			drawInk(gtx, c)
		}
		return D{Size: gtx.Constraints.Min}
	}
}

func layLabel(b *ButtonDef) layout.Widget {
	if b.Text == "" {
		return func(gtx C) D { return D{} }
	}
	return func(gtx C) D {
		return b.th.LabelPadding.Layout(gtx, func(gtx C) D {
			switch {
			case (b.Style == Text || b.Style == Outlined) && gtx.Queue == nil:
				paint.ColorOp{Color: Disabled(b.bg)}.Add(gtx.Ops)
			case b.Style == Outlined && b.Hovered():
				paint.ColorOp{Color: Hovered(b.bg)}.Add(gtx.Ops)
			case b.Style == Contained:
				paint.ColorOp{Color: b.fg}.Add(gtx.Ops)
			case b.Style == Outlined || b.Style == Text:
				paint.ColorOp{Color: b.bg}.Add(gtx.Ops)
			}
			return aLabel{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.th.TextSize, b.Text)
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
				size := gtx.Px(b.th.TextSize.Scale(1.2)) //TODO: Move const to theme
				gtx.Constraints = layout.Exact(image.Pt(size, size))
				return b.Icon.Layout(gtx, b.fg)
			})
		}
	}
	return func(gtx C) D { return D{} }
}

func (b *ButtonDef) layout(gtx C) D {
	return b.padding.Layout(gtx, func(gtx C) D {
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
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
							gtx,
							layout.Rigid(layIcon(b)),
							layout.Rigid(layLabel(b)),
						)
					}),
			)
		})
	})
}
