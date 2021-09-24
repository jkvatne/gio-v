//+build !android,!ios

package wid

import "gioui.org/unit"

// PlatformTooltip creates a tooltip styled to the current platform
// (desktop or mobile) by choosing based on the OS. This choice may
// not always be appropriate as it only uses the OS to decide.
func PlatformTooltip(th *Theme, text string, width unit.Value) Tooltip {
	return DesktopTooltip(th, text, width)
}
