package main

import (
	"fmt"
)

func formatSize(size int64) string {
	if size > 1e9 {
		return fmt.Sprintf("%.1f GiB", float64(size)/float64(1e9))
	} else if size > 1e6 {
		return fmt.Sprintf("%.1f MiB", float64(size)/float64(1e6))
	} else if size > 1e3 {
		return fmt.Sprintf("%.1f KiB", float64(size)/float64(1e3))
	}
	return fmt.Sprintf("%d B", size)
}

func formatItemName(item ItemInfo) string {
	if item.isDir {
		return "/" + item.name
	}
	return " " + item.name
}
