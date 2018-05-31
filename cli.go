package main

import (
	"image"

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

	status := tui.NewStatusBar("")

	root := tui.NewVBox(
		tui.NewLabel("gdu ~ Use arrow keys to navigate, press ? for help"),
		list,
		tui.NewSpacer(),
		status,
	)
	return root, list, status
}

type SizeRatio struct {
	tui.WidgetBase

	part int
}

func NewSizeRatio(part int) *SizeRatio {
	return &SizeRatio{
		part: part,
	}
}
func (p *SizeRatio) Draw(painter *tui.Painter) {
	painter.DrawRune(0, 0, '[')
	for i := 0; i < 10; i++ {
		if p.part > i {
			painter.DrawRune(i+1, 0, '=')
		}
	}
	painter.DrawRune(11, 0, ']')
}
func (p *SizeRatio) MinSizeHint() image.Point {
	return image.Point{12, 1}
}
func (p *SizeRatio) SizeHint() image.Point {
	return image.Point{12, 1}
}

type MinSizeLabel struct {
	*tui.Label
}

func NewMinSizeLabel(text string) *MinSizeLabel {
	return &MinSizeLabel{
		Label: tui.NewLabel(text),
	}
}
func (l *MinSizeLabel) MinSizeHint() image.Point {
	return image.Point{len(l.Label.Text()), 1}
}
