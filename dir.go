package main

import (
    "fmt"

    "github.com/marcusolsson/tui-go"
)

func showDir(ui tui.UI, dirStats DirInfo) {
    root, list, status := CreateListWindow()
	ui.SetWidget(root)

	status.SetText(
		fmt.Sprintf(
			"Apparent size: %v Items: %d",
			formatSize(dirStats.totalSize),
			dirStats.totalItemCount,
		))

	biggestItemSize := dirStats.items[0].size
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
	ui.SetKeybinding("PgDn", func() { list.Select(len(dirStats.items) - 1) })

    list.OnItemActivated(func(t *tui.Table){
        item := dirStats.items[t.Selected()]
        if item.isDir {
            showDir(ui, item.subDir)
        }
    })
}
