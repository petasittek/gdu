package main

import (
    "github.com/marcusolsson/tui-go"
)

func CreateAnalysisWindow() (tui.Widget, *tui.Label, *tui.Label) {
	status := tui.NewStatusBar("")

	statsLabel := tui.NewLabel("Total items: 0 Size: 0")
	currentItemLabel := tui.NewLabel("Current item: ")

	window := tui.NewVBox(
		tui.NewPadder(10, 1, statsLabel),
		tui.NewPadder(12, 1, currentItemLabel),
	)
	window.SetSizePolicy(tui.Expanding, tui.Preferred)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	root := tui.NewVBox(
		tui.NewPadder(2, 0, wrapper),
		status,
	)

	return root, currentItemLabel, statsLabel
}

func CreateListWindow() (tui.Widget, *tui.Table, *tui.StatusBar) {
	list := tui.NewTable(0, 0)
	list.SetColumnStretch(0, 1)
	list.SetColumnStretch(1, 1)

	status := tui.NewStatusBar("Scan in progress... Press q or ESC to abort")

	root := tui.NewVBox(
		list,
		status,
	)
	return root, list, status
}
