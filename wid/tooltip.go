package wid

import (
	"image"
	"image/color"
	"time"

	"gioui.org/font"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

const (
	tipAreaHoverDelay        = time.Millisecond * 1200
	tipAreaLongPressDuration = time.Millisecond * 1200
	tipAreaFadeDuration      = time.Millisecond * 500
	longPressDelay           = time.Millisecond * 1200
	CursorSizeX              = 10
	CursorSizeY              = 32
)

// Tooltip implements a material design tool tip as defined at:
// https://material.io/components/tooltips#specs
type Tooltip struct {
	VisibilityAnimation
	// MaxWidth is the maximum width of the tool-tip box. Should be less than form width.
	MaxWidth unit.Dp
	// Text defines the content of the tooltip.
	Text      widget.Label
	position  image.Point
	Hover     InvalidateDeadline
	Press     InvalidateDeadline
	LongPress InvalidateDeadline
	Fgc       color.NRGBA
	Bgc       color.NRGBA
	TooltipRR unit.Dp
	TextSize  unit.Sp
	init      bool
	shaper    *text.Shaper
	font      font.Font
}

// MobileTooltip constructs a tooltip suitable for use on mobile devices.
func MobileTooltip(th *Theme) Tooltip {
	return Tooltip{
		Fgc:      th.TooltipOnBackground,
		Bgc:      th.TooltipBackground,
		font:     font.Font{Weight: font.Medium},
		shaper:   th.Shaper,
		TextSize: th.TextSize * 0.9,
	}
}

// DesktopTooltip constructs a tooltip suitable for use on desktop devices.
func DesktopTooltip(th *Theme) Tooltip {
	return Tooltip{
		Fgc:       th.TooltipOnBackground,
		Bgc:       th.TooltipBackground,
		MaxWidth:  th.TooltipWidth,
		TooltipRR: th.TooltipCornerRadius,
		font:      font.Font{Weight: font.Medium},
		shaper:    th.Shaper,
		TextSize:  th.TextSize * 0.9,
	}

}

// InvalidateDeadline helps to ensure that a frame is generated at a specific
// point in time in the future. It does this by always requesting a future
// invalidation at its target time until it reaches its target time. This
// makes animating delays much cleaner.
type InvalidateDeadline struct {
	// The time at which a frame needs to be drawn.
	Target time.Time
	// Whether the deadline is active.
	Active bool
}

// SetTarget configures a specific time in the future at which a frame should
// be rendered.
func (i *InvalidateDeadline) SetTarget(t time.Time) {
	i.Active = true
	i.Target = t
}

// Process checks the current frame time and either requests a future invalidation
// or does nothing. It returns whether the current frame is the frame requested
// by the last call to SetTarget.
func (i *InvalidateDeadline) Process(gtx C) bool {
	if !i.Active {
		return false
	}
	if gtx.Now.Before(i.Target) {
		op.InvalidateOp{At: i.Target}.Add(gtx.Ops)
		return false
	}
	i.Active = false
	return true
}

// ClearTarget cancels a request to invalidate in the future.
func (i *InvalidateDeadline) ClearTarget() {
	i.Active = false
}

// Layout renders the provided widget with the provided tooltip. The tooltip
// will be summoned if the widget is hovered or long-pressed.
func (t *Tooltip) Layout(gtx C, hint string, w layout.Widget) D {
	if hint == "" {
		return w(gtx)
	}
	if !t.init {
		t.init = true
		t.VisibilityAnimation.State = Invisible
		t.VisibilityAnimation.Duration = tipAreaFadeDuration
	}
	for _, e := range gtx.Events(t) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		t.position.X = int(e.Position.X)
		t.position.Y = int(e.Position.Y)
		switch e.Type {
		case pointer.Enter:
			t.Hover.SetTarget(gtx.Now.Add(tipAreaHoverDelay))
		case pointer.Leave:
			t.VisibilityAnimation.Disappear(gtx.Now)
			t.Hover.ClearTarget()
		case pointer.Press:
			t.Hover.ClearTarget()
			t.Press.SetTarget(gtx.Now.Add(longPressDelay))
		case pointer.Release:
			t.Hover.ClearTarget()
			t.Press.ClearTarget()
		case pointer.Cancel:
			t.Hover.ClearTarget()
			t.Press.ClearTarget()
		}
	}
	if t.Hover.Process(gtx) {
		t.VisibilityAnimation.Appear(gtx.Now)
	}
	if t.Press.Process(gtx) {
		t.VisibilityAnimation.Appear(gtx.Now)
		t.LongPress.SetTarget(gtx.Now.Add(tipAreaLongPressDuration))
	}
	if t.LongPress.Process(gtx) {
		t.VisibilityAnimation.Disappear(gtx.Now)
	}
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(w),
		layout.Expanded(func(gtx C) D {
			defer pointer.PassOp{}.Push(gtx.Ops).Pop()
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			r := clip.Rect(rect).Push(gtx.Ops)
			pointer.InputOp{
				Tag:   t,
				Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave | pointer.Move,
			}.Add(gtx.Ops)
			r.Pop()
			gtx.Constraints.Min = image.Point{}
			if t.Visible() {
				macro := op.Record(gtx.Ops)
				v := t.VisibilityAnimation.Revealed(gtx)
				bg := WithAlpha(t.Bgc, uint8(v*255))
				t.Fgc = WithAlpha(t.Fgc, uint8(v*255))
				gtx.Constraints.Max.X = gtx.Metric.Dp(t.MaxWidth)
				p := unit.Dp(t.TextSize * 0.5)
				inset := layout.Inset{Top: p, Right: p, Bottom: p, Left: p}
				dims := layout.Stack{}.Layout(
					gtx,
					layout.Expanded(func(gtx C) D {
						rr := gtx.Dp(t.TooltipRR)
						outline := image.Rectangle{Max: gtx.Constraints.Min}
						paint.FillShape(gtx.Ops, bg, clip.UniformRRect(outline, rr).Op(gtx.Ops))
						paintBorder(gtx, outline, MulAlpha(t.Fgc, 128), unit.Dp(0.5), gtx.Dp(t.TooltipRR))
						return D{}
					}),
					layout.Stacked(func(gtx C) D {
						colMacro := op.Record(gtx.Ops)
						paint.ColorOp{Color: t.Fgc}.Add(gtx.Ops)
						return inset.Layout(gtx, func(gtx C) D {
							// paint.ColorOp{Color: t.Fgc}.Add(gtx.Ops)
							return t.Text.Layout(gtx, t.shaper, t.font, t.TextSize, hint, colMacro.Stop())
						})
					}),
				)
				dx := int(MouseX) + CursorSizeX + dims.Size.X - WinX
				if dx < 0 {
					dx = 0
				}
				dy := int(MouseY) + CursorSizeY + dims.Size.Y - WinY
				if dy < 0 {
					dy = 0
				}
				call := macro.Stop()
				macro = op.Record(gtx.Ops)
				op.Offset(t.position.Add(image.Pt(-dx+CursorSizeX, -dy+CursorSizeY))).Add(gtx.Ops)
				call.Add(gtx.Ops)
				call = macro.Stop()
				op.Defer(gtx.Ops, call)
			}
			return D{}
		}),
	)
}
