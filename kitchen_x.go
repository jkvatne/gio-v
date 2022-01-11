package main

import (
	"fmt"
	"gio-v/wid"
	"image"
	"image/color"
	"math"
	"time"

	"gioui.org/text"

	"gioui.org/io/pointer"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type (
	// D is a shortcut
	D = layout.Dimensions
	// C is a shortcut
	C = layout.Context
)

var (
	lineEditor1 = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	lineEditor2 = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	button            = new(widget.Clickable)
	greenButton       = new(widget.Clickable)
	iconTextButton    = new(widget.Clickable)
	iconButton        = new(widget.Clickable)
	flatBtn           = new(widget.Clickable)
	radioButtonsGroup = new(widget.Enum)
	radioButtonValue  string
	list              = &widget.List{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
	float    = new(widget.Float)
	topLabel = "Hello, Gio"
)

type iconAndTextButton struct {
	theme  *material.Theme
	button *widget.Clickable
	icon   *widget.Icon
	word   string
}

func (b iconAndTextButton) Layout(gtx layout.Context) layout.Dimensions {
	return material.ButtonLayout(b.theme, b.button).Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			textIconSpacer := unit.Dp(5)

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					var d D
					if b.icon != nil {
						size := gtx.Px(unit.Dp(56)) - 2*gtx.Px(unit.Dp(16))
						gtx.Constraints = layout.Exact(image.Pt(size, size))
						d = b.icon.Layout(gtx, b.theme.ContrastFg)
					}
					return d
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					l := material.Body1(b.theme, b.word)
					l.Color = b.theme.Palette.ContrastFg
					return l.Layout(gtx)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)
		})
	})
}

func kitchenX(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range lineEditor1.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			topLabel = e.Text
			lineEditor1.SetText("")
		}
	}
	in := layout.UniformInset(unit.Dp(3))
	tl := material.Label(th, th.TextSize.Scale(48.0/16.0), topLabel)
	tl.Alignment = text.Middle
	widgets := []layout.Widget{
		func(gtx C) D {
			return in.Layout(gtx, tl.Layout)
		},
		func(gtx C) D {
			return in.Layout(gtx, func(gtx C) D {
				e := material.Editor(th, lineEditor1, "Value 1")
				e.Font.Style = text.Italic
				border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
				return border.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, e.Layout)
				})
			})
		},
		func(gtx C) D {
			return in.Layout(gtx, func(gtx C) D {
				e := material.Editor(th, lineEditor2, "Value 2")
				e.Font.Style = text.Italic
				border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(2)}
				return border.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, e.Layout)
				})
			})
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, material.IconButton(th, iconButton, icon, "??").Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, iconAndTextButton{theme: th, icon: icon, word: "Icon", button: iconTextButton}.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						for button.Clicked() {
							green = !green
						}
						dims := material.Button(th, button, "Click me!").Layout(gtx)
						pointer.CursorNameOp{Name: pointer.CursorPointer}.Add(gtx.Ops)
						return dims
					})
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						l := "Green"
						if !green {
							l = "Blue"
						}
						btn := material.Button(th, greenButton, l)
						if green {
							btn.Background = color.NRGBA{A: 0xff, R: 0x9e, G: 0x9d, B: 0x24}
						}
						return btn.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return in.Layout(gtx, func(gtx C) D {
						return material.Clickable(gtx, flatBtn, func(gtx C) D {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {
								flatBtnText := material.Body1(th, "Show other")
								if gtx.Queue == nil {
									flatBtnText.Color.A = 150
								}
								for flatBtn.Clicked() {
									page = "KitchenV"
									oldMode = "xx"
									PrintMemUsage("Gio-X")
									startTime = time.Now()
									count = 0
								}
								return layout.Center.Layout(gtx, flatBtnText.Layout)
							})
						})
					})
				}),
			)
		},
		material.ProgressBar(th, progress).Layout,
		wid.Value(currentTheme, func() string { return fmt.Sprintf(" %0.1f frames/second", count/time.Since(startTime).Seconds()) }),

		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(material.RadioButton(th, radioButtonsGroup, "r1", "RadioButton1").Layout),
				layout.Rigid(material.RadioButton(th, radioButtonsGroup, "r2", "RadioButton2").Layout),
				layout.Rigid(material.RadioButton(th, radioButtonsGroup, "r3", "RadioButton3").Layout),
			)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, material.Slider(th, float, 0, 2*math.Pi).Layout),
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx,
						material.Body1(th, fmt.Sprintf("%.2f", float.Value)).Layout,
					)
				}),
			)
		},
	}
	return material.List(th, list).Layout(gtx, len(widgets), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(3)).Layout(gtx, widgets[i])
	})
}
