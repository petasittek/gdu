package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"time"

	"github.com/marcusolsson/tui-go"
)

type ItemInfo struct {
	name   string
	size   int64
	isDir  bool
	subDir []ItemInfo
}

type CurrentProgress struct {
	currentItemName string
	itemCount       int
	totalSize       int64
	done            bool
}

type BySize []ItemInfo
func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].size < a[j].size }

func processDir(dir string, itemCount int, totalSize int64, statusChannel chan<- CurrentProgress) ([]ItemInfo, int, int64) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, itemCount, totalSize
	}

	dirStats := make([]ItemInfo, len(files))

	for i, f := range files {
		info := ItemInfo{
			name: f.Name(),
			isDir: f.IsDir(),
		}
		if f.IsDir() {
			subDirStats, subDirItemCount, subDirSize := processDir(
				path.Join(dir, f.Name()),
				itemCount,
				totalSize,
				statusChannel,
			)
			info.size = subDirSize
			info.subDir = subDirStats
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
			info.size = f.Size()
			totalSize += f.Size()
		}
		dirStats[i] = info
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
	sort.Sort(sort.Reverse(BySize(dirStats)))

	statusChannel <- CurrentProgress{done: true}

	ui.Update(func() {
		root, list, status := CreateListWindow()
		ui.SetWidget(root)

		status.SetText(
			fmt.Sprintf(
				"Apparent size: %v Items: %d",
				formatSize(totalDirSize),
				totalItemCount,
			))

		biggestItemSize := dirStats[0].size
		for _, item := range dirStats {
			part := float64(item.size) / float64(biggestItemSize) * 10.0

			list.AppendRow(
				NewMinSizeLabel(
					fmt.Sprintf(
						"%10s",
						formatSize(item.size),
					),
				),
				NewSizeRatio(int(part)),
				NewMinSizeLabel(item.name),
				tui.NewSpacer(),
			)
		}

		ui.SetKeybinding("PgUp", func() { list.Select(0) })
		ui.SetKeybinding("PgDn", func() { list.Select(len(dirStats) - 1) })
	})
}
