package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"

	"github.com/marcusolsson/tui-go"
)

type ItemInfo struct {
	name   string
	path   string
	size   int64
	isDir  bool
	subDir *DirInfo
}

type DirInfo struct {
	items     []ItemInfo
	parentDir *DirInfo
	size      int64
	itemCount int

	// only for current progress stats
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

func (a BySize) Len() int      { return len(a) }
func (a BySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool {
	iSize := a[i].size
	jSize := a[j].size
	if a[i].isDir {
		iSize = a[i].subDir.size
	}
	if a[j].isDir {
		jSize = a[j].subDir.size
	}
	return iSize < jSize
}

func processDir(dir string, parentDir *DirInfo, statusChannel chan<- CurrentProgress) *DirInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return &DirInfo{}
	}

	dirStats := &DirInfo{
		items:          make([]ItemInfo, len(files)),
		size:           0,
		itemCount:      0,
		totalSize:      parentDir.totalSize,
		totalItemCount: parentDir.totalItemCount,
	}

	if len(parentDir.items) > 0 {
		dirStats.parentDir = parentDir
	}

	for i, f := range files {
		path, _ := filepath.Abs(filepath.Join(dir, f.Name()))

		info := ItemInfo{
			name:  f.Name(),
			isDir: f.IsDir(),
			path:  path,
		}
		if f.IsDir() {
			subDirStats := processDir(
				filepath.Join(dir, f.Name()),
				dirStats,
				statusChannel,
			)
			info.size = subDirStats.size
			info.subDir = subDirStats
			dirStats.size += subDirStats.size
			dirStats.itemCount += subDirStats.itemCount
			dirStats.totalItemCount = subDirStats.totalItemCount
			dirStats.totalSize = subDirStats.totalSize

			select {
			case statusChannel <- CurrentProgress{
				currentItemName: info.path,
				itemCount:       dirStats.totalItemCount,
				totalSize:       dirStats.totalSize,
			}:
			default:
			}
		} else {
			info.size = f.Size()
			dirStats.totalSize += f.Size()
			dirStats.size += f.Size()
		}
		dirStats.items[i] = info
		dirStats.itemCount += 1
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

func processTopDir(dir string, ui tui.UI, currentItemLabel *tui.Label, statsLabel *tui.Label) *DirInfo {
	statusChannel := make(chan CurrentProgress)

	go updateCurrentProgress(ui, currentItemLabel, statsLabel, statusChannel)

	topDir := &DirInfo{}

	dirStats := processDir(
		dir,
		topDir,
		statusChannel,
	)

	statusChannel <- CurrentProgress{done: true}

	return dirStats
}
