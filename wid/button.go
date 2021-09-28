// SPDX-License-Identifier: Unlicense OR MIT

package wid

import (
	"gio-v/f32color"
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

var GlobalDisable = false

type ButtonStyle int

const (
	Contained ButtonStyle = iota
	Text
	Outlined
	Round
)

type ButtonDef struct {
	Clickable
	Tooltip
	tipArea      *TipArea
	th           *Theme
	shadow       ShadowStyle
	disabler     *bool
	Text         string
	helptext     string
	ToolTipWidth unit.Value
	Font         text.Font
	shaper       text.Shaper
	Icon         *Icon
	Width        unit.Value
	Style        ButtonStyle
}

type BtnOption func(*ButtonDef)

func Width(w float32) BtnOption {
	return func(b *ButtonDef) {
		b.Width = unit.Dp(w)
	}
}

// BtnIcon sets button icon
func BtnIcon(i *Icon) BtnOption {
	return func(b *ButtonDef) {
		b.Icon = i
	}
}

func Handler(f func()) BtnOption {
	foo := func(b bool) {f()}
	return func(b *ButtonDef) {
		b.handler = foo
	}
}

func Disable(v *bool) BtnOption {
	return func(b *ButtonDef) {
		b.disabler = v
	}
}

func Hint(s string) BtnOption {
	return func(b *ButtonDef){
		b.helptext = s
	}
}

func (b *ButtonDef) apply(options []BtnOption) {
	for _, option := range options {
		option(b)
	}
}

func Button(style ButtonStyle, th *Theme, label string, options ...BtnOption) func(gtx C) D {
	b := ButtonDef{}
	b.SetupTabs()
	b.th = th
	b.tipArea = &TipArea{}
	b.Text = label
	b.Font = text.Font{Weight: text.Medium}
	b.shadow = Shadow(th.CornerRadius, th.Elevation)
	b.shaper = th.Shaper
	b.Style = style
	if b.ToolTipWidth.V==0 {
		b.ToolTipWidth.V = th.TextSize.V*20
	}
	b.apply(options)
	if b.helptext != "" {
		b.Tooltip = PlatformTooltip(th, b.helptext, b.ToolTipWidth)
	}
	if style == Round {
		b.shadow = Shadow(b.th.CornerRadius, b.th.Elevation)
	}

	return func(gtx C) D {
		dims := b.Layout(gtx)
		b.HandleClick()
		pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
		return dims
	}
}

func drawInk(gtx layout.Context, c Press) {
	now := gtx.Now
	t := float64(now.Sub(c.Start).Seconds())
	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}
	endt := float64(end.Sub(c.Start).Seconds())
	// Compute the fade-in/out position in [0;1].
	var haste float64
	if c.Cancelled {
		// If the press was cancelled before the inkwell
		// was fully faded in, fast forward the animation
		// to match the fade-out.
		if h := 0.5 - endt/0.9; h > 0 {
			haste = h
		}
	}
	// Fade in.
	half1 := math.Max(t/0.9 + haste, 0.5)
	if half1 > 0.5 {
		half1 = 0.5
	}
	// Fade out.
	half2 := now.Sub(end).Seconds()/0.9 + haste
	if half2 > 0.5 {
		return
	}
	alphat := half1 + half2
	// Compute the expand position in [0;1].
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		t = endt
	}
	sizet := math.Min(t*2, 1.0)
	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}
	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	size := math.Max(float64(gtx.Constraints.Min.Y), float64(gtx.Constraints.Min.X))
	alpha := 0.7 * alphat * 2 * t2 * (3.0 - 3.0*alphat)
	ba, bc := byte(alpha*0xff), byte(0x80)
	defer op.Save(gtx.Ops).Load()
	rgba := f32color.MulAlpha(color.NRGBA{A: 0xff, R: bc, G: bc, B: bc}, ba)
	ink := paint.ColorOp{Color: rgba}
	ink.Add(gtx.Ops)
	rr := float32( size*math.Sqrt(2.0) *sizet * sizet * (3.0 - 2.0*sizet))
	op.Offset(c.Position.Add(f32.Point{
		X: -rr,
		Y: -rr,
	})).Add(gtx.Ops)
	clip.UniformRRect(f32.Rectangle{Max: f32.Pt(2*rr, 2*rr)}, rr).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func paintBorder(gtx layout.Context, outline f32.Rectangle, col color.NRGBA, width float32, rr float32) {
	paint.FillShape(gtx.Ops,
		col,
		clip.Stroke{
			Path:  clip.UniformRRect(outline, rr).Path(gtx.Ops),
			Style: clip.StrokeStyle{Width: width},
		}.Op(),
	)
}

func (b *ButtonDef) LayoutBackground() func(gtx C) D {
	return func(gtx C) D {
		if b.Focused() || b.Hovered() {
			b.shadow.Layout(gtx)
		}
		rr := gtx.Pxr(b.th.CornerRadius)
		if b.Style==Round {
			rr = float32(gtx.Constraints.Min.Y) / 2.0
		}
		outline := f32.Rectangle{Max: f32.Point{
			X: float32(gtx.Constraints.Min.X),
			Y: float32(gtx.Constraints.Min.Y),
		}}
		clip.UniformRRect(outline, rr).Add(gtx.Ops)

		switch {
		case b.Style == Outlined && gtx.Queue == nil:
			paintBorder(gtx, outline, f32color.Disabled(b.th.Palette.Primary), gtx.Pxr(b.th.BorderThickness), rr)
		case b.Style == Outlined:
			paintBorder(gtx, outline, b.th.Palette.Primary, gtx.Pxr(b.th.BorderThickness), rr)
		case gtx.Queue == nil && (b.Style == Contained || b.Style == Round):
			paint.Fill(gtx.Ops, f32color.Disabled(b.th.Palette.Primary))
		case gtx.Queue == nil:
			// Transparent background when disabled
		case (b.Hovered() || b.Focused()) && (b.Style == Contained || b.Style == Round):
			paint.Fill(gtx.Ops, f32color.Hovered(b.th.Palette.Primary))
		case b.Style == Contained || b.Style == Round:
			paint.Fill(gtx.Ops, b.th.Palette.Primary)
		}
		for _, c := range b.History() {
			drawInk(gtx, c)
		}
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
}

func layLabel(b *ButtonDef) layout.Widget {
	if b.Text != "" {
		return func(gtx C) D {
			return b.th.LabelInset.Layout(gtx, func(gtx C) D {
				switch {
				case (b.Style == Text || b.Style == Outlined) && gtx.Queue == nil:
					paint.ColorOp{Color: f32color.Disabled(b.th.Palette.Primary)}.Add(gtx.Ops)
				case b.Style == Outlined && b.Hovered():
					paint.ColorOp{Color: f32color.Hovered(b.th.Palette.Primary)}.Add(gtx.Ops)
				case b.Style == Contained:
					paint.ColorOp{Color: b.th.Palette.OnPrimary}.Add(gtx.Ops)
				case b.Style == Outlined || b.Style == Text:
					paint.ColorOp{Color: b.th.Palette.Primary}.Add(gtx.Ops)
				}
				return aLabel{Alignment: text.Middle}.Layout(gtx, b.shaper, b.Font, b.th.TextSize, b.Text)
			})
		}
	}
	return func(gtx C) D { return D{} }
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
				size := gtx.Px(b.th.TextSize.Scale(1.2))
				gtx.Constraints = layout.Exact(image.Pt(size, size))
				return b.Icon.Layout(gtx, b.th.Palette.OnPrimary)
			})
		}
	}
	return func(gtx C) D { return D{} }
}

func (b *ButtonDef) Layout(gtx layout.Context) layout.Dimensions {
	return b.tipArea.Layout(gtx, b.Tooltip, func(gtx C) D {
		b.disabled = false
		if b.disabler!= nil && *b.disabler || GlobalDisable {
			gtx = gtx.Disabled()
			b.disabled = true
		}
		min := gtx.Constraints.Min
		if b.Width.V <= 1.0 {
			min.X = gtx.Px(b.Width.Scale(float32(gtx.Constraints.Max.X)))
		} else if min.X < gtx.Px(b.Width) {
			min.X = gtx.Px(b.Width)
		}
		if min.X>gtx.Constraints.Max.X {
			min.X = gtx.Constraints.Max.X
		}
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(b.LayoutBackground()),
			layout.Stacked(
				func(gtx C) D {
					gtx.Constraints.Min = min
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle, Spacing: layout.SpaceSides}.Layout(
						gtx,
						layout.Rigid(layIcon(b)),
						layout.Rigid(layLabel(b)),
					)
				}),
			layout.Expanded(b.LayoutClickable),
		)
	})
}
