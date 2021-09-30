# Extension of Gio
See [gioui.org](https://gioui.org).

This is a set of widgets made for my own use. They replaces (or complements) gioui.org/widget and gioui.org/widget/material but the rest of the gio code is imported without modifications.

The code is copied extensively from the following sources:

* https://github.com/gioui/gio
* https://github.com/gioui/gio-x.git
* https://github.com/gioui/gio-example.git

# Features

## Material design
Quite a lot of effort has gone into making it look like google's examples of material design forms. 
The shadows are noe almost identical, and hover/focus are at least similar. 

## Keyboard only operation
All (TODO) widgets handle keyboard only operation. Focus is moved py TAB/SHIF-TAB keys. 
The tab sequence is given by the declaration sequence automaticaly. No mouse is needed.

## Extended theming
The theme is very much extended, with default valuse for all colors and sizes. You can set up
several themes for different types of buttons etc, and use the thems when declaring the widgets.

## Dynamic resizing
Everything scales with the text size, and the text size is set automaticaly as a fraction of 
the window size. This makes it easy to write programs that are maximized to fill the screen, 
or are operating mostly in full-screen mode. (Creating maximized windows depends on a smallmodification
of gio itself, and is implemented in my gio fork)

## Importing gio directly
The gio module itself does not need to be modified. The excelent work by Eliass Naur and 
Chris Waldon is used without modifications. My widgets are only high-level extensions.

## Widget configuration by optional arguments
Using a similar techniqe to what is used for window setup, you can configure the widgets as needed.
Functions with optional parameters is not possible in Golang, but you can have a variable number of
interface parameters as in ```func(options ...Option)``` . See example below

## Easy setup of forms
Here is an example from demo.go. The widgets can take a variable number of options for things like 
width and hints. Othewise default fallbacks are used.

```
	root = wid.MakeList(
		th, layout.Vertical,
		wid.Label(th, "Demo page", text.Middle, 2.0),
		wid.Button(th, "WIDE BUTTON",
			wid.W(950),
			wid.Pad(30, 15, 15, 0),
			wid.Hint("This is a dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
		wid.MakeFlex(
			wid.Label(th, "Dark mode", text.Start, 1.0),
			wid.Switch(th, darkMode, onSwitchMode),
		),
		wid.Checkbox(th, "Checkbox", darkMode, nil),
		wid.MakeFlex(
			wid.RoundButton(th, icons.ContentAdd,
				wid.Hint("This is another dummy button - it has no function except displaying this text, testing long help texts. Perhaps breaking into several lines")),
			wid.Button(th, "Home", wid.BtnIcon(icons.ActionHome), wid.Disable(&darkMode)),
			wid.Button(th, "Check", wid.BtnIcon(icons.ActionCheckCircle)),
			wid.Button(&thb, "Change color", wid.Handler(onClick)),
			wid.TextButton(th, "Text button"),
			wid.OutlineButton(th, "Outline button"),
			wid.Label(th, "Disabled", text.End, 1.0),
			wid.Switch(th, false, doDisable),
		),
		wid.MakeFlex(
			wid.Combo(th, unit.Value{}, 0, []string{"Option A", "Option B", "Option C"}),
			wid.Combo(th, unit.Value{200, 0}, 1, []string{"Option 1", "Option 2", "Option 3"}),
			wid.Combo(th, unit.Value{300, 0}, 1, []string{"Option 1", "Option 2", "Option 3"}),
			wid.Combo(th, unit.Value{}, 0, []string{"Option A", "Option B", "Option C"}),
		),
		wid.Edit(th, wid.Hint("Value 1"), wid.W(950)),
		wid.Edit(th, wid.Hint("Value 2")),
		wid.Edit(th, wid.Hint("Value 3")),
```

# License
Dual MIT/Unlicense; same as Gio

# Demo
![Demo.go](https://github.com/jkvatne/gio-v/blob/main/demo.png)
