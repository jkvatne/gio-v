package wid

import (
	"image"
	"image/color"
	"time"

	"gioui.org/font"

	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

const (
	tipAreaHoverDelay   = time.Millisecond * 900
	tipAreaFadeDuration = time.Millisecond * 500
	longPressDelay      = time.Millisecond * 500
	CursorSizeX         = 16
	CursorSizeY         = 32
)

// Tooltip implements a material design tool tip as defined at:
// https://material.io/components/tooltips#specs
type TooltipDef struct {
	VisibilityAnimation
	// MaxWidth is the maximum width of the tool-tip box. Should be less than form width.
	MaxWidth unit.Dp
	// Position of the last mouse pointer event
	position  image.Point
	Hover     InvalidateDeadline
	Fgc       color.NRGBA
	Bgc       color.NRGBA
	TooltipRR unit.Dp
	TextSize  unit.Sp
	init      bool
	font      font.Font
}

// DesktopTooltip constructs a tooltip suitable for use on desktop devices.
func Tooltip(th *Theme) TooltipDef {
	return TooltipDef{
		Fgc:       th.TooltipOnBackground,
		Bgc:       th.TooltipBackground,
		MaxWidth:  th.TooltipWidth,
		TooltipRR: th.TooltipCornerRadius,
		font:      font.Font{Weight: font.Medium},
		TextSize:  th.TextSize,
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
func (t *TooltipDef) Layout(gtx C, hint string, th *Theme) D {
	if hint == "" {
		return D{}
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
		if !t.Visible() {
			t.position.X = mouseX - int(e.Position.X)
			t.position.Y = mouseY - int(e.Position.Y)
		}
		switch e.Kind {
		case pointer.Enter:
			t.Hover.SetTarget(gtx.Now.Add(tipAreaHoverDelay))
		case pointer.Leave:
			t.VisibilityAnimation.Disappear(gtx.Now)
			t.Hover.ClearTarget()
		case pointer.Press:
			t.Hover.ClearTarget()
			t.Hover.SetTarget(gtx.Now.Add(longPressDelay))
		case pointer.Release:
		case pointer.Cancel:
			t.Hover.ClearTarget()
		default:
		}
	}
	if t.Hover.Process(gtx) {
		t.VisibilityAnimation.Appear(gtx.Now)
	}
	defer pointer.PassOp{}.Push(gtx.Ops).Pop()
	// Detect pointer movement within gtx.Constraints.Min
	r := clip.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Push(gtx.Ops)
	pointer.InputOp{
		Tag:   t,
		Kinds: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave | pointer.Move,
	}.Add(gtx.Ops)
	r.Pop()
	extHeight := gtx.Constraints.Min.Y
	if t.Visible() {
		// p is the inside padding of the tooltip
		p := Px(gtx, th.TextSize) / 2
		tooltipMacro := op.Record(gtx.Ops)
		// Calculate colors according to visibility
		v := t.VisibilityAnimation.Revealed(gtx)
		bg := WithAlpha(t.Bgc, uint8(v*255))
		t.Fgc = WithAlpha(t.Fgc, uint8(v*255))
		gtx.Constraints.Min = image.Point{}
		gtx.Constraints.Max = image.Point{gtx.Metric.Dp(t.MaxWidth), 99999}
		rr := Px(gtx, t.TooltipRR)
		textMacro := op.Record(gtx.Ops)
		textOffset := op.Offset(image.Pt(p, p)).Push(gtx.Ops)
		fgCol := op.Record(gtx.Ops)
		// Draw text
		paint.ColorOp{Color: t.Fgc}.Add(gtx.Ops)
		dims := widget.Label{}.Layout(gtx, th.Shaper, t.font, t.TextSize, hint, fgCol.Stop())
		textOffset.Pop()
		drawTextOp := textMacro.Stop()
		outline := image.Rectangle{Max: image.Pt(gtx.Metric.Dp(t.MaxWidth)+p, dims.Size.Y+p*2)}
		// Move the location to the left so it does not go outside the right edge of the form
		dx := Min(0, WinX-t.position.X-outline.Max.X-10)
		var dy int
		if WinY-t.position.Y > extHeight+outline.Max.Y {
			dy = Min(extHeight, WinY-t.position.Y)
		} else {
			dy = -outline.Max.Y - 5
		}
		op.Offset(image.Pt(dx, dy)).Add(gtx.Ops)
		// Fill background and border
		paint.FillShape(gtx.Ops, bg, clip.UniformRRect(outline, rr).Op(gtx.Ops))
		paintBorder(gtx, outline, t.Fgc, float32(Px(gtx, unit.Dp(0.75))), Px(gtx, t.TooltipRR))
		// Then actually draw the text
		drawTextOp.Add(gtx.Ops)
		// End the tooltipMacro and defer drawing so they appear on top of everything else
		op.Defer(gtx.Ops, tooltipMacro.Stop())
	}
	return D{}

}
