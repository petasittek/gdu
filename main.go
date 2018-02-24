package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		log.Fatal(err)
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
			statsLabel.SetText("Total items: " + fmt.Sprintf("%d", progress.itemCount) + " Size: " + fmt.Sprintf("%d", progress.totalSize))
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
		root, list, status := createListWindow()
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
				tui.NewLabel(fmt.Sprintf("%d", item.size)),
			)
		}
	})
}

func createAnalysisWindow() (tui.Widget, *tui.Label, *tui.Label) {
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

func createListWindow() (tui.Widget, *tui.Table, *tui.StatusBar) {
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

func main() {
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "."
	}

	root, currentItemLabel, statsLabel := createAnalysisWindow()

	ui := tui.New(root)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })

	go processTopDir(dir, ui, currentItemLabel, statsLabel)

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
