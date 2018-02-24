package main

import (
	"log"
	"os"

	"github.com/marcusolsson/tui-go"
)

func main() {
	var dir string
	if len(os.Args) > 1 {
		dir = os.Args[1]
	} else {
		dir = "."
	}

	root, currentItemLabel, statsLabel := CreateAnalysisWindow()

	ui := tui.New(root)

	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })

	go processTopDir(dir, ui, currentItemLabel, statsLabel)

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
