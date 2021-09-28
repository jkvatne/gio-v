package wid

import (
	"image"
	"image/color"
	"time"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

// Tooltip implements a material design tool tip as defined at:
// https://material.io/components/tooltips#specs
type Tooltip struct {
	th *Theme
	// MaxWidth is the maximum width of the tool-tip box. Should be less than form width.
	MaxWidth unit.Value
	// Text defines the content of the tooltip.
	Text LabelDef
	// Bg defines the color of the RRect background.
	Bg color.NRGBA
	VisibilityAnimation
	position  f32.Point
	Hover     InvalidateDeadline
	Press     InvalidateDeadline
	LongPress InvalidateDeadline
	init      bool
	// HoverDelay is the delay between the cursor entering the tip area
	// and the tooltip appearing.
	HoverDelay time.Duration
	// LongPressDelay is the required duration of a press in the area for
	// it to count as a long press.
	LongPressDelay time.Duration
	// LongPressDuration is the amount of time the tooltip should be displayed
	// after being triggered by a long press.
	LongPressDuration time.Duration
	// FadeDuration is the amount of time it takes the tooltip to fade in
	// and out.
	FadeDuration time.Duration
}

// MobileTooltip constructs a tooltip suitable for use on mobile devices.
func MobileTooltip(th *Theme, tips string) Tooltip {
	txt := CreateLabelDef(th, tips, text.Start, 0.8)
	txt.Color = th.Background
	return Tooltip{
		th:   th,
		Bg:   WithAlpha(th.OnBackground, 220),
		Text: txt,
	}
}

// DesktopTooltip constructs a tooltip suitable for use on desktop devices.
func DesktopTooltip(th *Theme, tips string,) Tooltip {
	txt := CreateLabelDef(th, tips, text.Start, 0.9)
	txt.Color = th.OnSecondary
	return Tooltip{
		th:       th,
		Bg:       WithAlpha(th.Secondary, 220),
		Text:     txt,
		MaxWidth: th.TooltipWidth,
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

func Interpolate(a, b color.NRGBA, progress float32) color.NRGBA {
	var out color.NRGBA
	out.R = uint8(int16(a.R) - int16(float32(int16(a.R)-int16(b.R))*progress))
	out.G = uint8(int16(a.G) - int16(float32(int16(a.G)-int16(b.G))*progress))
	out.B = uint8(int16(a.B) - int16(float32(int16(a.B)-int16(b.B))*progress))
	out.A = uint8(int16(a.A) - int16(float32(int16(a.A)-int16(b.A))*progress))
	return out
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

// TipArea holds the state information for displaying a tooltip. The zero
// value will choose sensible defaults for all fields.
type TipArea struct {
	VisibilityAnimation
	position  f32.Point
	Hover     InvalidateDeadline
	Press     InvalidateDeadline
	LongPress InvalidateDeadline
	init      bool
	// HoverDelay is the delay between the cursor entering the tip area
	// and the tooltip appearing.
	HoverDelay time.Duration
	// LongPressDelay is the required duration of a press in the area for
	// it to count as a long press.
	LongPressDelay time.Duration
	// LongPressDuration is the amount of time the tooltip should be displayed
	// after being triggered by a long press.
	LongPressDuration time.Duration
	// FadeDuration is the amount of time it takes the tooltip to fade in
	// and out.
	FadeDuration time.Duration
}

const (
	tipAreaHoverDelay        = time.Millisecond * 500
	tipAreaLongPressDuration = time.Millisecond * 1500
	tipAreaFadeDuration      = time.Millisecond * 750
	longPressTheshold        = time.Millisecond * 750
)

// Layout renders the provided widget with the provided tooltip. The tooltip
// will be summoned if the widget is hovered or long-pressed.
func (t *Tooltip) Layout(gtx C, tip Tooltip, w layout.Widget) D {
	if tip.Text.Text == "" {
		return w(gtx)
	}
	if !t.init {
		t.init = true
		t.VisibilityAnimation.State = Invisible
		if t.HoverDelay == time.Duration(0) {
			t.HoverDelay = tipAreaHoverDelay
		}
		if t.LongPressDelay == time.Duration(0) {
			t.LongPressDelay = longPressTheshold
		}
		if t.LongPressDuration == time.Duration(0) {
			t.LongPressDuration = tipAreaLongPressDuration
		}
		if t.FadeDuration == time.Duration(0) {
			t.FadeDuration = tipAreaFadeDuration
		}
		t.VisibilityAnimation.Duration = t.FadeDuration
	}
	for _, e := range gtx.Events(t) {
		e, ok := e.(pointer.Event)
		if !ok {
			continue
		}
		if !t.Visible() {
			t.position = e.Position
		}
		switch e.Type {
		case pointer.Enter:
			t.Hover.SetTarget(gtx.Now.Add(t.HoverDelay))
		case pointer.Leave:
			t.VisibilityAnimation.Disappear(gtx.Now)
			t.Hover.ClearTarget()
		case pointer.Press:
			t.Hover.ClearTarget()
			t.Press.SetTarget(gtx.Now.Add(t.LongPressDelay))
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
		t.LongPress.SetTarget(gtx.Now.Add(t.LongPressDuration))
	}
	if t.LongPress.Process(gtx) {
		t.VisibilityAnimation.Disappear(gtx.Now)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(w),
		layout.Expanded(func(gtx C) D {
			defer op.Save(gtx.Ops).Load()
			pointer.PassOp{Pass: true}.Add(gtx.Ops)
			pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			pointer.InputOp{
				Tag:   t,
				Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave | pointer.Move,
			}.Add(gtx.Ops)

			gtx.Constraints.Min = image.Point{}
			if t.Visible() {
				macro := op.Record(gtx.Ops)
				v := t.VisibilityAnimation.Revealed(gtx)
				tip.Bg = Interpolate(WithAlpha(tip.Bg, 0), tip.Bg, v)
				tip.Text.Color = Interpolate(WithAlpha(tip.Text.Color, 0), tip.Text.Color, v)
				gtx.Constraints.Max.X = gtx.Metric.Px(tip.MaxWidth)
				//dims := tip.Layout(gtx)
				// Layout renders the tooltip.
				dims := layout.Stack{}.Layout(
					gtx,
					layout.Expanded(func(gtx C) D {
						radius := float32(gtx.Px(t.th.TooltipCornerRadius))
						paint.FillShape(gtx.Ops, t.Bg, clip.RRect{
							Rect: f32.Rectangle{
								Max: layout.FPt(gtx.Constraints.Min),
							},
							NW: radius,
							NE: radius,
							SW: radius,
							SE: radius,
						}.Op(gtx.Ops))
						return D{}
					}),
					layout.Stacked(func(gtx C) D {
						return t.th.TooltipInset.Layout(gtx, t.Text.Layout)
					}),
				)


				if int(t.position.X)+dims.Size.X > gtx.Constraints.Max.X {
					t.position.X = float32(gtx.Constraints.Max.X - dims.Size.X)
				}
				if int(t.position.Y)+dims.Size.Y > gtx.Constraints.Max.Y {
					t.position.Y = float32(gtx.Constraints.Max.Y - dims.Size.Y)
				}
				call := macro.Stop()
				macro = op.Record(gtx.Ops)
				op.Offset(t.position.Add(f32.Pt(5, 5))).Add(gtx.Ops)
				call.Add(gtx.Ops)
				call = macro.Stop()
				op.Defer(gtx.Ops, call)
			}
			return D{}
		}),
	)
}
