package main

import (
    "fmt"
)

func formatSize(size int64) string {
    if size > 10e9 {
        return fmt.Sprintf("%.1f GiB", float64(size) / float64(10e9))
    } else if size > 10e6 {
        return fmt.Sprintf("%.1f MiB", float64(size) / float64(10e6))
    } else if size > 10e3 {
        return fmt.Sprintf("%.1f KiB", float64(size) / float64(10e3))
    }
    return fmt.Sprintf("%d B", size)
}
