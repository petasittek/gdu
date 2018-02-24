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
	subDir DirInfo
}

type DirInfo struct {
	items          []ItemInfo
	parentDir      *DirInfo
	totalSize      int64
	totalItemCount int
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

func processDir(dir string, parentDir DirInfo, statusChannel chan<- CurrentProgress) DirInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return DirInfo{}
	}

	dirStats := DirInfo{
		items: make([]ItemInfo, len(files)),
		totalSize: parentDir.totalSize,
		totalItemCount: parentDir.totalItemCount,
	}
	if len(parentDir.items) > 0 {
		dirStats.parentDir = &parentDir
	}

	for i, f := range files {
		info := ItemInfo{
			name: f.Name(),
			isDir: f.IsDir(),
		}
		if f.IsDir() {
			subDirStats := processDir(
				path.Join(dir, f.Name()),
				dirStats,
				statusChannel,
			)
			info.size = subDirStats.totalSize
			info.subDir = subDirStats
			dirStats.totalItemCount = subDirStats.totalItemCount
			dirStats.totalSize = subDirStats.totalSize

			select {
			case statusChannel <- CurrentProgress{
				currentItemName: f.Name(),
				itemCount:       dirStats.totalItemCount,
				totalSize:       dirStats.totalSize,
			}:
			default:
			}
		} else {
			info.size = f.Size()
			dirStats.totalSize += f.Size()
		}
		dirStats.items[i] = info
		dirStats.totalItemCount += 1
	}
	sort.Sort(sort.Reverse(BySize(dirStats.items)))
	return dirStats
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

func processTopDir(dir string, ui tui.UI, currentItemLabel *tui.Label, statsLabel *tui.Label) DirInfo {
	statusChannel := make(chan CurrentProgress)

	go updateCurrentProgress(ui, currentItemLabel, statsLabel, statusChannel)

	dirStats := processDir(
		dir,
		DirInfo{},
		statusChannel,
	)

	statusChannel <- CurrentProgress{done: true}

	return dirStats
}
