# Extension of Gio

See [gioui.org](https://gioui.org).

This is a set of widgets made for my own use. They replaces (or complements) gioui.org/widget and
gioui.org/widget/material but the rest of the gio code is imported without modifications.

The code is copied extensively from the following sources:

* https://github.com/gioui/gio
* https://github.com/gioui/gio-x.git
* https://github.com/gioui/gio-example.git

THIS IS WORK IN PROGRESS - ANYTHING CAN CHANGE AT ANY TIME

# Features

## Material design

The design follows closely Google Material 3, where a few primary colors are used to generate all the other 
colors. Most other design elements can be tuned by modifying the default theme.

## Keyboard only operation

All widgets handle keyboard only operation. Focus is moved py TAB/SHIFT-TAB keys using standard gio

## Extended theming

The theme is very much extended, with default values for all colors and sizes. You can set up several themes for
different types of buttons etc, and use the themes when declaring the widgets.

Dark and Light mode are both supported, and can be easily selected from a widget. 

## Dynamic resizing

Everything scales with the text size, and the text size can be set automatically as a fraction of the window size. This
makes it easy to write programs that are maximized to fill the screen, or are operating mostly in full-screen mode

## Importing gio directly

The gio module itself does not need to be modified. The excellent work by Elias Naur and Chris Waldon is used without
modifications. My widgets are only a high-level extensions, replacing the material widgets in gio.

## Widget configuration by optional arguments

All widget functions can have any number of options, as ```func(options ...Option)``` . See example below

## Easy setup of forms

Here is an example from /examples/hello. The widgets can take a variable number of options for things like width and hints.
Otherwise, default fallbacks are used. The defaults are mostly defined in the theme.

```
package main

import (
	"gio-v/wid"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
)

func main() {
	go wid.Run(
		app.NewWindow(app.Title("Hello Gio demo"), app.Size(unit.Dp(900), unit.Dp(500))),
		wid.NewTheme(gofont.Collection(), 14),
		hello,
	)
	app.Main()
}

func hello(th *wid.Theme) layout.Widget {
	return wid.List(th, wid.Overlay,
		wid.Label(th, "Hello gio..", wid.Middle(), wid.Heading(), wid.Bold()),
		wid.Label(th, "A small demo program using 25 lines toal"),
	)
}
```

# Immediate mode?

This implementation does not follow the gio recomendations fully. This is actualy not an 
immediate mode design. The widgets are fully persistent, and callbacks and pointers are 
used extensively. This is done to make it much more user-friendly, and it is primarily intended for
desktop applications, where resources are plentiful.

Switches and edits modify the corresponding variables directly, via pointers. When the variable is
modified, the corresponding widget is emmediately uppdated without any acction from the program.
This is typically done from an other go-routine.

Note that the program is not yet protected from race conditions. 
The plan is to include a global lock.

# License

Dual MIT/Unlicense; same as Gio

# Demo

![Demo.go](https://github.com/jkvatne/gio-v/blob/main/demo.png)

# Grid demo

![Demo.go](https://github.com/jkvatne/gio-v/blob/main/grid.png)
