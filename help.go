package main

import (
	"github.com/marcusolsson/tui-go"
)

type StyledLabel struct {
	Style string
	*tui.Label
}

func (s *StyledLabel) Draw(p *tui.Painter) {
	p.WithStyle(s.Style, func(p *tui.Painter) {
		s.Label.Draw(p)
	})
}

func ShowHelpWindow(ui tui.UI) {
	t := tui.NewTheme()
	bold := tui.Style{Bold: tui.DecorationOn}
	t.SetStyle("bold", bold)

	table := tui.NewGrid(0, 0)

	table.AppendRow(
		&StyledLabel{Label: tui.NewLabel("up"), Style: "bold"},
		tui.NewLabel("Move cursor up"),
	)
	table.AppendRow(
		&StyledLabel{Label: tui.NewLabel("down"), Style: "bold"},
		tui.NewLabel("Move cursor down"),
	)
	table.AppendRow(
		&StyledLabel{Label: tui.NewLabel("enter"), Style: "bold"},
		tui.NewLabel("Open selected directory"),
	)
	table.AppendRow(
		&StyledLabel{Label: tui.NewLabel("d"), Style: "bold"},
		tui.NewLabel("Delete selected file or directory"),
	)
	table.SetSizePolicy(tui.Expanding, tui.Minimum)

	window := tui.NewVBox(table)
	window.SetSizePolicy(tui.Expanding, tui.Preferred)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	root := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	ui.SetTheme(t)
	ui.SetWidget(root)

	ui.ClearKeybindings()
	showCurrentDir := func() {
		showDir(ui, currentDir)
	}
	ui.SetKeybinding("q", showCurrentDir)
	ui.SetKeybinding("Esc", showCurrentDir)
}
