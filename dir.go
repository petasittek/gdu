package main

import (
	"fmt"

	"github.com/marcusolsson/tui-go"
)

var currentDir *DirInfo

func showDir(ui tui.UI, dirStats *DirInfo) {
	root, list, status := CreateListWindow()
	currentDir = dirStats
	ui.SetWidget(root)

	ui.ClearKeybindings()
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("h", func() { ShowHelpWindow(ui) })
	ui.SetKeybinding("?", func() { ShowHelpWindow(ui) })

	status.SetText(
		fmt.Sprintf(
			"Apparent size: %v Items: %d",
			formatSize(dirStats.size),
			dirStats.itemCount,
		))

	if dirStats.parentDir != nil {
		list.AppendRow(
			tui.NewLabel(""),
			tui.NewLabel(""),
			tui.NewLabel("/.."),
			tui.NewSpacer(),
		)
	}

	biggestItemSize := int64(0)
	if len(dirStats.items) > 0 {
		biggestItemSize = dirStats.items[0].size
	}

	for _, item := range dirStats.items {
		part := float64(item.size) / float64(biggestItemSize) * 10.0

		list.AppendRow(
			NewMinSizeLabel(
				fmt.Sprintf(
					"%10s",
					formatSize(item.size),
				),
			),
			NewSizeRatio(int(part)),
			NewMinSizeLabel(formatItemName(item)),
			tui.NewSpacer(),
		)
	}

	ui.SetKeybinding("PgUp", func() { list.Select(0) })
	ui.SetKeybinding("PgDn", func() {
		if dirStats.parentDir != nil {
			list.Select(len(dirStats.items))
		} else {
			list.Select(len(dirStats.items) - 1)
		}
	})

	list.OnItemActivated(func(t *tui.Table) {
		index := t.Selected()
		if dirStats.parentDir != nil {
			if index == 0 {
				showDir(ui, dirStats.parentDir)
				return
			} else {
				index--
			}
		}

		item := dirStats.items[index]
		if item.isDir {
			showDir(ui, item.subDir)
		}
	})
}
