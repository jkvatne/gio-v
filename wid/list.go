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
	"gioui.org/unit"
)

// fromListPosition converts a layout.Position into two floats representing
// the location of the viewport on the underlying content. It needs to know
// the number of elements in the list and the major-axis size of the list
// in order to do this. The returned values will be in the range [0,1], and
// start will be less than or equal to end.
func fromListPosition(lp layout.Position, elements int, majorAxisSize int) (start, end float32) {
	// Approximate the size of the scrollable content.
	lengthPx := float32(lp.Length)
	meanElementHeight := lengthPx / float32(elements)

	// Determine how much of the content is visible.
	listOffsetF := float32(lp.Offset)
	visiblePx := float32(majorAxisSize)
	visibleFraction := visiblePx / lengthPx

	// Compute the location of the beginning of the viewport.
	viewportStart := (float32(lp.First)*meanElementHeight + listOffsetF) / lengthPx

	return viewportStart, clamp1(viewportStart + visibleFraction)
}

// ScrollTrackStyle configures the presentation of a track for a scroll area.
type ScrollTrackStyle struct {
	// MajorPadding and MinorPadding along the major and minor axis of the
	// scrollbar's track. This is used to keep the scrollbar from touching
	// the edges of the content area.
	MajorPadding, MinorPadding unit.Dp
	// Color of the track background.
	Color color.NRGBA
}

// ScrollIndicatorStyle configures the presentation of a scroll indicator.
type ScrollIndicatorStyle struct {
	// MajorMinLen is the smallest that the scroll indicator is allowed to
	// be along the major axis.
	MajorMinLen unit.Dp
	// MinorWidth is the width of the scroll indicator across the minor axis.
	MinorWidth unit.Dp
	// Color and HoverColor are the normal and hovered colors of the scroll
	// indicator.
	Color, HoverColor color.NRGBA
	// CornerRadius is the corner radius of the rectangular indicator. 0
	// will produce square corners. 0.5*MinorWidth will produce perfectly
	// round corners.
	CornerRadius unit.Dp
}

// ScrollbarStyle configures the presentation of a scrollbar.
type ScrollbarStyle struct {
	Scrollbar *Scrollbar
	Track     ScrollTrackStyle
	Indicator ScrollIndicatorStyle
}

// MakeScrollbarStyle configures the presentation of a scrollbar using the provided
// theme and state.
func MakeScrollbarStyle(th *Theme) ScrollbarStyle {
	lightFg := th.Fg(Canvas)
	lightFg.A = 150
	darkFg := lightFg
	darkFg.A = 200
	return ScrollbarStyle{
		Scrollbar: &Scrollbar{},
		Track: ScrollTrackStyle{
			MajorPadding: unit.Dp(th.ScrollMajorPadding),
			MinorPadding: unit.Dp(th.ScrollMinorPadding),
			Color:        th.TrackColor,
		},
		Indicator: ScrollIndicatorStyle{
			MajorMinLen:  unit.Dp(th.ScrollMajorMinLen),
			MinorWidth:   unit.Dp(th.ScrollMinorWidth),
			CornerRadius: unit.Dp(th.ScrollCornerRadius),
			Color:        lightFg,
			HoverColor:   darkFg,
		},
	}
}

// Width returns the minor axis width of the scrollbar in its current
// configuration (taking padding for the scroll track into account).
func (s ScrollbarStyle) Width() unit.Dp {
	return s.Indicator.MinorWidth + s.Track.MinorPadding + s.Track.MinorPadding
}

// Layout the scrollbar.
func (s ScrollbarStyle) Layout(gtx C, axis layout.Axis, viewportStart, viewportEnd float32) D {
	// if !rangeIsScrollable(viewportStart, viewportEnd) {
	if viewportEnd-viewportStart >= 1 {
		return D{}
	}

	// Set minimum constraints in an axis-independent way, then convert to
	// the correct representation for the current axis.
	convert := axis.Convert
	maxMajorAxis := convert(gtx.Constraints.Max).X
	gtx.Constraints.Min.X = maxMajorAxis
	gtx.Constraints.Min.Y = gtx.Dp(s.Width())
	gtx.Constraints.Min = convert(gtx.Constraints.Min)
	gtx.Constraints.Max = gtx.Constraints.Min

	s.Scrollbar.Layout(gtx, axis, viewportStart, viewportEnd)

	// Darken indicator if hovered.
	if s.Scrollbar.IndicatorHovered() {
		s.Indicator.Color = s.Indicator.HoverColor
	}

	return s.layout(gtx, axis, viewportStart, viewportEnd)
}

// layout the scroll track and indicator.
func (s ScrollbarStyle) layout(gtx C, axis layout.Axis, viewportStart, viewportEnd float32) D {
	inset := layout.Inset{
		Top:    s.Track.MajorPadding,
		Bottom: s.Track.MajorPadding,
		Left:   s.Track.MinorPadding,
		Right:  s.Track.MinorPadding,
	}
	if axis == layout.Horizontal {
		inset.Top, inset.Bottom, inset.Left, inset.Right = inset.Left, inset.Right, inset.Top, inset.Bottom
	}
	// Capture the outer constraints because layout.Stack will reset
	// the minimum to zero.
	outerConstraints := gtx.Constraints

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// Lay out the draggable track underneath the scroll indicator.
			area := image.Rectangle{
				Max: gtx.Constraints.Min,
			}
			pointerArea := clip.Rect(area)
			defer pointerArea.Push(gtx.Ops).Pop()
			s.Scrollbar.AddDrag(gtx.Ops)

			// Stack a normal clickable area on top of the draggable area
			// to capture non-dragging clicks.
			defer pointer.PassOp{}.Push(gtx.Ops).Pop()
			defer pointerArea.Push(gtx.Ops).Pop()
			s.Scrollbar.AddTrack(gtx.Ops)

			paint.FillShape(gtx.Ops, s.Track.Color, clip.Rect(area).Op())
			return D{}
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints = outerConstraints
			return inset.Layout(gtx, func(gtx C) D {
				// Use axis-independent constraints.
				gtx.Constraints.Min = axis.Convert(gtx.Constraints.Min)
				gtx.Constraints.Max = axis.Convert(gtx.Constraints.Max)

				// Compute the pixel size and position of the scroll indicator within
				// the track.
				trackLen := gtx.Constraints.Min.X
				viewStart := int(math.Round(float64(viewportStart) * float64(trackLen)))
				viewEnd := int(math.Round(float64(viewportEnd) * float64(trackLen)))
				indicatorLen := Max(viewEnd-viewStart, gtx.Dp(s.Indicator.MajorMinLen))
				if viewStart+indicatorLen > trackLen {
					viewStart = trackLen - indicatorLen
				}
				indicatorDims := axis.Convert(image.Point{
					X: indicatorLen,
					Y: gtx.Dp(s.Indicator.MinorWidth),
				})
				radius := gtx.Dp(s.Indicator.CornerRadius)

				// Lay out the indicator.
				offset := axis.Convert(image.Pt(viewStart, 0))
				defer op.Offset(offset).Push(gtx.Ops).Pop()
				paint.FillShape(gtx.Ops, s.Indicator.Color, clip.RRect{
					Rect: image.Rectangle{
						Max: indicatorDims,
					},
					SW: radius,
					NW: radius,
					NE: radius,
					SE: radius,
				}.Op(gtx.Ops))

				// Add the indicator pointer hit area.
				area := clip.Rect(image.Rectangle{Max: indicatorDims})
				defer pointer.PassOp{}.Push(gtx.Ops).Pop()
				defer area.Push(gtx.Ops).Pop()
				s.Scrollbar.AddIndicator(gtx.Ops)
				return layout.Dimensions{Size: axis.Convert(gtx.Constraints.Min)}
			})
		}),
	)
}

// AnchorStrategy defines a means of attaching a scrollbar to content.
type AnchorStrategy uint8

const (
	// Occupy reserves space for the scrollbar, making the underlying
	// content region smaller on one axis.
	Occupy AnchorStrategy = iota
	// Overlay causes the scrollbar to float atop the content without
	// occupying any space. Content in the underlying area can be occluded
	// by the scrollbar.
	Overlay
)

// ListStyle configures the presentation of a layout.List with a scrollbar.
type ListStyle struct {
	list       *layout.List
	theme      *Theme
	Hpos       int
	VScrollBar ScrollbarStyle
	HScrollBar ScrollbarStyle
	AnchorStrategy
}

// List makes a vertical listp
func List(th *Theme, a AnchorStrategy, Size f32.Point, widgets ...layout.Widget) layout.Widget {
	node := makeNode(widgets)
	listStyle := ListStyle{
		list:           &layout.List{Axis: layout.Vertical},
		VScrollBar:     MakeScrollbarStyle(th),
		HScrollBar:     MakeScrollbarStyle(th),
		AnchorStrategy: a,
	}
	listStyle.theme = th
	return func(gtx C) D {
		if Size.X <= 1.0 && Size.X > 0.001 {
			gtx.Constraints.Max.X = int(float32(gtx.Constraints.Max.X) * Size.X)
		} else if Size.X > 1.0 {
			gtx.Constraints.Max.X = Min(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(Size.X)))
		}
		if Size.Y <= 1.0 && Size.Y > 0.001 {
			gtx.Constraints.Max.Y = int(float32(WinY) * Size.Y)
		} else if Size.Y > 1.0 {
			gtx.Constraints.Max.Y = Min(gtx.Constraints.Max.Y, gtx.Dp(unit.Dp(Size.Y)))
		}
		cl := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
		c := th.Bg(Canvas)
		paint.Fill(gtx.Ops, c)
		cl.Pop()

		var ch []layout.Widget
		for i := 0; i < len(node.children); i++ {
			ch = append(ch, *node.children[i].w)
		}
		return listStyle.Layout(
			gtx,
			len(ch),
			func(gtx C, i int) D {
				return ch[i](gtx)
			},
		)
	}
}

// Layout the list and its scrollbar.
func (l *ListStyle) Layout(gtx C, length int, w layout.ListElement) D {
	originalConstraints := gtx.Constraints
	// Determine how much space the scrollbar occupies.
	hBarWidth := gtx.Dp(l.HScrollBar.Width())
	vBarWidth := gtx.Dp(l.VScrollBar.Width())

	// Reserve space for the scrollbars using the gtx constraints.
	if l.AnchorStrategy == Occupy {
		gtx.Constraints.Max.X -= vBarWidth
		gtx.Constraints.Min.X -= vBarWidth
		gtx.Constraints.Max.Y -= hBarWidth
		gtx.Constraints.Min.Y -= hBarWidth
	}
	// Draw the list
	macro := op.Record(gtx.Ops)
	c := gtx
	// Must set Max.X to infinity to allow rows wider than the frame.
	c.Constraints.Max.X = inf
	listDims := l.list.Layout(c, length, w)
	call := macro.Stop()
	cl := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)

	gtx.Constraints = originalConstraints
	gtx.Constraints.Min = gtx.Constraints.Max

	pt := image.Pt(-l.Hpos, 0)
	trans := op.Offset(pt).Push(gtx.Ops)
	call.Add(gtx.Ops)
	trans.Pop()
	cl.Pop()
	totalWidth := listDims.Size.X

	// Draw the Vertical scrollbar.
	if vBarWidth > 0 {
		totalHeight := l.list.Position.Length
		// Get vertical scroll info.
		delta := l.VScrollBar.Scrollbar.ScrollDistance()
		if delta != 0 {
			l.list.Position.Offset += int(math.Round(float64(float32(totalHeight) * delta)))
		}

		c := gtx
		start, end := fromListPosition(l.list.Position, length, listDims.Size.Y)
		if l.AnchorStrategy == Overlay {
			c.Constraints.Max.Y -= hBarWidth
		} else {
			c.Constraints.Max.X += vBarWidth
		}
		c.Constraints.Min = c.Constraints.Max
		layout.E.Layout(c, func(gtx layout.Context) layout.Dimensions {
			return l.VScrollBar.Layout(gtx, layout.Vertical, start, end)
		})
	}

	// Draw the Horizontal scrollbar
	if hBarWidth > 0 {
		c := gtx
		start := float32(l.Hpos) / float32(totalWidth)
		end := start + float32(c.Constraints.Max.X)/float32(totalWidth)
		if l.AnchorStrategy == Occupy {
			// c.Constraints.Max.Y += hBarWidth
		}
		c.Constraints.Min = c.Constraints.Max
		layout.S.Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min = gtx.Constraints.Max
			return l.HScrollBar.Layout(gtx, layout.Horizontal, start, end)
		})
		delta := l.HScrollBar.Scrollbar.ScrollDistance()
		if delta != 0 {
			deltaPx := int(math.Round(float64(float32(totalWidth) * delta)))
			l.Hpos = limit(l.Hpos+deltaPx, 0, totalWidth-gtx.Constraints.Max.X)
		}
	}

	return listDims
}

func limit(x int, min int, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}
