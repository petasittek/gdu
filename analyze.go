package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"github.com/marcusolsson/tui-go"
)

type ItemInfo struct {
	size   int64
	isDir  bool
	subDir map[string]*ItemInfo
}

type CurrentProgress struct {
	currentItemName string
	itemCount       int
	totalSize       int64
	done            bool
}

func processDir(dir string, itemCount int, totalSize int64, statusChannel chan<- CurrentProgress) (map[string]*ItemInfo, int, int64) {
	dirStats := make(map[string]*ItemInfo)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dirStats, itemCount, totalSize
	}

	for _, f := range files {
		dirStats[f.Name()] = &ItemInfo{
			isDir: f.IsDir(),
		}
		if f.IsDir() {
			subDirStats, subDirItemCount, subDirSize := processDir(
				path.Join(dir, f.Name()),
				itemCount,
				totalSize,
				statusChannel,
			)
			dirStats[f.Name()].size = subDirSize
			dirStats[f.Name()].subDir = subDirStats
			itemCount = subDirItemCount
			totalSize = subDirSize

			select {
			case statusChannel <- CurrentProgress{
				currentItemName: f.Name(),
				itemCount:       itemCount,
				totalSize:       totalSize,
			}:
			default:
			}
		} else {
			dirStats[f.Name()].size = f.Size()
			totalSize += f.Size()
		}
		itemCount += 1
	}
	return dirStats, itemCount, totalSize
}

func updateCurrentProgress(ui tui.UI, currentItemLabel *tui.Label, statsLabel *tui.Label, statusChannel <-chan CurrentProgress) {
	for {
		progress := <-statusChannel

		if progress.done {
			return
		}

		ui.Update(func() {
			currentItemLabel.SetText("Current item: " + progress.currentItemName)
			statsLabel.SetText("Total items: " + fmt.Sprintf("%d", progress.itemCount) + " Size: " + formatSize(progress.totalSize))
		})

		time.Sleep(100 * time.Millisecond)
	}
}

func processTopDir(dir string, ui tui.UI, currentItemLabel *tui.Label, statsLabel *tui.Label) {
	statusChannel := make(chan CurrentProgress)

	go updateCurrentProgress(ui, currentItemLabel, statsLabel, statusChannel)

	dirStats, totalItemCount, totalDirSize := processDir(
		dir,
		0,
		0,
		statusChannel,
	)

	statusChannel <- CurrentProgress{done: true}

	ui.Update(func() {
		root, list, status := CreateListWindow()
		ui.SetWidget(root)

		status.SetText(
			fmt.Sprintf(
				"Apparent size: %d Items: %d",
				totalDirSize,
				totalItemCount,
			))

		for name, item := range dirStats {
			list.AppendRow(
				tui.NewLabel(name),
				tui.NewLabel(formatSize(item.size)),
			)
		}
	})
}
